package com.yonotech.ppob.data.remote.interceptor

import okhttp3.Interceptor
import okhttp3.Response
import java.io.IOException

class RetryInterceptor(
    private val maxRetries: Int = 3
) : Interceptor {

    @Throws(IOException::class)
    override fun intercept(chain: Interceptor.Chain): Response {
        val request = chain.request()
        var response = chain.proceed(request)
        var retryCount = 0

        while (!response.isSuccessful && retryCount < maxRetries) {
            val status = response.code
            if (status == 502 || status == 503 || status == 504) {
                retryCount++
                val delay = (1000L shl retryCount).coerceAtMost(10000L)
                Thread.sleep(delay)
                response.close()
                response = chain.proceed(request)
            } else {
                break
            }
        }
        return response
    }
}