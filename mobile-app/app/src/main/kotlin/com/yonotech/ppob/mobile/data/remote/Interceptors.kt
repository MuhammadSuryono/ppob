package com.yonotech.ppob.mobile.data.remote

import com.yonotech.ppob.mobile.data.local.TokenManager
import com.yonotech.ppob.mobile.data.remote.dto.RefreshTokenRequest
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.runBlocking
import okhttp3.Authenticator
import okhttp3.Interceptor
import okhttp3.Response
import okhttp3.Route
import javax.inject.Inject
import javax.inject.Provider

class AuthInterceptor @Inject constructor(
    private val tokenManager: TokenManager
) : Interceptor {
    override fun intercept(chain: Interceptor.Chain): Response {
        val originalRequest = chain.request()
        val requestBuilder = originalRequest.newBuilder()
            .header("Content-Type", "application/json")
            .header("Accept", "application/json")

        // List of public endpoints that don't require an access token
        val publicEndpoints = listOf(
            "/auth/initiate",
            "/auth/register",
            "/auth/login",
            "/auth/send-otp",
            "/auth/verify-otp",
            "/auth/verify-password",
            "/auth/verify-pin",
            "/auth/verify-credential",
            "/auth/refresh"
        )

        val path = originalRequest.url.encodedPath
        val isPublic = publicEndpoints.any { path.contains(it) }

        if (isPublic) {
            println("AuthInterceptor: Skipping token for public endpoint ${originalRequest.url}")
            return chain.proceed(requestBuilder.build())
        }

        val token = runBlocking {
            tokenManager.accessToken.first()
        }

        if (!token.isNullOrBlank()) {
            println("AuthInterceptor: Adding Token for ${originalRequest.url}")
            requestBuilder.header("Authorization", "Bearer $token")
        } else {
            println("AuthInterceptor: TOKEN IS MISSING for ${originalRequest.url}")
        }

        val request = requestBuilder.build()
        return chain.proceed(request)
    }
}

class AuthAuthenticator @Inject constructor(
    private val tokenManager: TokenManager,
    private val authServiceProvider: Provider<AuthService>
) : Authenticator {
    override fun authenticate(route: Route?, response: okhttp3.Response): okhttp3.Request? {
        println("AuthAuthenticator: 401 intercepted for ${response.request.url}")

        // Skip refresh if the request was to an auth endpoint
        if (response.request.url.encodedPath.contains("/auth/")) {
            println("AuthAuthenticator: 401 on auth endpoint, skipping refresh")
            return null
        }

        val refreshToken = runBlocking {
            tokenManager.refreshToken.first()
        }

        if (refreshToken.isNullOrBlank()) {
            println("AuthAuthenticator: Refresh token is null or blank")
            return null
        }

        synchronized(this) {
            val currentToken = runBlocking { tokenManager.accessToken.first() }
            val authHeader = response.request.header("Authorization")

            if (authHeader != null && authHeader != "Bearer $currentToken") {
                println("AuthAuthenticator: Token already refreshed by another request, retrying with new token")
                return response.request.newBuilder()
                    .header("Authorization", "Bearer $currentToken")
                    .build()
            }

            println("AuthAuthenticator: Attempting to refresh token...")
            val refreshResponse = try {
                runBlocking {
                    authServiceProvider.get().refreshToken(RefreshTokenRequest(refreshToken))
                }
            } catch (e: Exception) {
                println("AuthAuthenticator: Exception during refresh: ${e.message}")
                null
            }

            if (refreshResponse != null && refreshResponse.isSuccessful) {
                val newAuthData = refreshResponse.body()
                if (newAuthData?.accessToken != null) {
                    println("AuthAuthenticator: Token refresh successful")
                    runBlocking {
                        tokenManager.saveTokens(
                            newAuthData.accessToken,
                            newAuthData.refreshToken ?: refreshToken
                        )
                    }
                    return response.request.newBuilder()
                        .header("Authorization", "Bearer ${newAuthData.accessToken}")
                        .build()
                } else {
                    println("AuthAuthenticator: Token refresh response successful but access token is null")
                }
            } else {
                println("AuthAuthenticator: Token refresh failed: code=${refreshResponse?.code()}")
                runBlocking {
                    tokenManager.clearTokens()
                }
            }
        }

        return null
    }
}

class RetryInterceptor(private val maxRetries: Int = 3) : Interceptor {
    override fun intercept(chain: Interceptor.Chain): Response {
        val request = chain.request()
        var response = chain.proceed(request)
        var retryCount = 0

        while (!response.isSuccessful && retryCount < maxRetries) {
            val status = response.code
            if (status == 502 || status == 503 || status == 504) {
                retryCount++
                val delay = (1000L shl retryCount) // exponential backoff
                Thread.sleep(delay)
                response = chain.proceed(request)
            } else {
                break
            }
        }
        return response
    }
}
