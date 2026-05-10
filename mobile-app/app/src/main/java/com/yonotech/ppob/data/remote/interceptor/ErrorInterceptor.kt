package com.yonotech.ppob.data.remote.interceptor

import android.util.Log
import com.yonotech.ppob.data.local.datastore.AuthPreferencesDataStore
import com.yonotech.ppob.data.remote.model.ApiError
import com.yonotech.ppob.data.remote.model.ApiErrorResponse
import com.yonotech.ppob.data.remote.model.AuthResponse
import com.squareup.moshi.Moshi
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.runBlocking
import okhttp3.Interceptor
import okhttp3.Response
import java.io.IOException
import javax.inject.Inject

class ErrorInterceptor @Inject constructor(
    private val authPreferencesDataStore: AuthPreferencesDataStore
) : Interceptor {

    private val moshi = Moshi.Builder().build()
    private val errorAdapter = moshi.adapter(ApiErrorResponse::class.java)

    @Throws(IOException::class)
    override fun intercept(chain: Interceptor.Chain): Response {
        val response = chain.proceed(chain.request())

        if (!response.isSuccessful) {
            val errorBody = response.body?.string()

            when (response.code) {
                401 -> {
                    // Token might be expired, try to refresh
                    Log.d("ErrorInterceptor", "401 received, attempting token refresh")
                    val newResponse = handleTokenRefresh(chain)
                    if (newResponse != null) {
                        return newResponse
                    }
                }
                429 -> {
                    // Rate limited - could add retry-after logic
                    Log.w("ErrorInterceptor", "Rate limited: $errorBody")
                }
                502, 503 -> {
                    Log.e("ErrorInterceptor", "Server error: $errorBody")
                }
            }
        }

        return response
    }

    private fun handleTokenRefresh(chain: Interceptor.Chain): Response? {
        return try {
            val currentRefreshToken = runBlocking {
                authPreferencesDataStore.authPreferences.first().refreshToken
            }

            if (currentRefreshToken.isEmpty()) {
                Log.d("ErrorInterceptor", "No refresh token available")
                return null
            }

            // Build a new request with the new token
            val newRequest = chain.request().newBuilder()
                .removeHeader("Authorization")
                .build()

            // Note: Actual token refresh should be done via AuthRepository
            // This interceptor just rebuilds the request; the ViewModel handles refresh logic
            Log.d("ErrorInterceptor", "Token refresh delegated to ViewModel layer")
            null

        } catch (e: Exception) {
            Log.e("ErrorInterceptor", "Error during token refresh", e)
            null
        }
    }
}
