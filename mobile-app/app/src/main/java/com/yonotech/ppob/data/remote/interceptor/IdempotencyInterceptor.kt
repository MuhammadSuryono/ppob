package com.yonotech.ppob.data.remote.interceptor

import okhttp3.Interceptor
import okhttp3.Response
import java.io.IOException
import java.util.UUID

class IdempotencyInterceptor : Interceptor {

    @Throws(IOException::class)
    override fun intercept(chain: Interceptor.Chain): Response {
        val original = chain.request()
        val url = original.url.toString()

        // Add idempotency key for transaction initiation and other write operations
        val writeEndpoints = listOf("/transaction/", "/wallet/")
        val needsIdempotency = writeEndpoints.any { url.contains(it) }

        if (needsIdempotency && original.header("Idempotency-Key") == null) {
            val requestBuilder = original.newBuilder()
                .header("Idempotency-Key", UUID.randomUUID().toString())
            return chain.proceed(requestBuilder.build())
        }

        return chain.proceed(original)
    }
}