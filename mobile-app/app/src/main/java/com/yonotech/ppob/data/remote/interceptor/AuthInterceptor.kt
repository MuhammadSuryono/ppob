package com.yonotech.ppob.data.remote.interceptor

import com.yonotech.ppob.data.local.datastore.AuthPreferencesDataStore
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.runBlocking
import okhttp3.Interceptor
import okhttp3.Response
import java.io.IOException

class AuthInterceptor(
    private val authPreferencesDataStore: AuthPreferencesDataStore
) : Interceptor {

    @Throws(IOException::class)
    override fun intercept(chain: Interceptor.Chain): Response {
        val original = chain.request()

        // Only add auth header for API calls (not auth endpoints themselves)
        val url = original.url.toString()
        val isAuthEndpoint = url.contains("/auth/")

        val requestBuilder = original.newBuilder()
            .header("Content-Type", "application/json")
            .header("Accept", "application/json")

        if (!isAuthEndpoint) {
            val accessToken = runBlocking {
                try {
                    authPreferencesDataStore.authPreferences.first().accessToken
                } catch (e: Exception) {
                    ""
                }
            }
            if (accessToken.isNotEmpty()) {
                requestBuilder.header("Authorization", "Bearer $accessToken")
            }
        }

        return chain.proceed(requestBuilder.build())
    }
}
