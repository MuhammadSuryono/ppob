package com.yonotech.ppob.data.remote.interceptor

import android.util.Log
import okhttp3.Interceptor
import okhttp3.Response
import java.io.IOException

class LoggingInterceptor : Interceptor {

    @Throws(IOException::class)
    override fun intercept(chain: Interceptor.Chain): Response {
        val request = chain.request()
        val t1 = System.nanoTime()

        Log.d("OkHttp", "Sending request: ${request.url} on ${chain.connection()}")
        Log.d("OkHttp", "Request headers: ${request.headers}")

        val response = chain.proceed(request)
        val t2 = System.nanoTime()

        Log.d(
            "OkHttp",
            "Received response for ${response.request.url} in ${(t2 - t1) / 1e6}ms"
        )

        return response
    }
}