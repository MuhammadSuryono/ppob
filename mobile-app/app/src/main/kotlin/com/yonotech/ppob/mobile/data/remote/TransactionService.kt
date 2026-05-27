package com.yonotech.ppob.mobile.data.remote

import com.yonotech.ppob.mobile.data.remote.dto.transaction.InitiateTransactionRequest
import com.yonotech.ppob.mobile.data.remote.dto.transaction.InquiryRequest
import com.yonotech.ppob.mobile.data.remote.dto.transaction.InquiryResponse
import com.yonotech.ppob.mobile.data.remote.dto.transaction.TransactionResponse
import retrofit2.Response
import retrofit2.http.Body
import retrofit2.http.GET
import retrofit2.http.Header
import retrofit2.http.POST
import retrofit2.http.Path
import retrofit2.http.Query

interface TransactionService {
    @POST("transactions/initiate")
    suspend fun initiate(
        @Header("Idempotency-Key") idempotencyKey: String,
        @Body request: InitiateTransactionRequest
    ): Response<TransactionResponse>

    @POST("transactions/inquiry")
    suspend fun inquiry(
        @Body request: InquiryRequest
    ): Response<InquiryResponse>

    @GET("transactions/by-id/{id}")
    suspend fun getStatus(@Path("id") id: String): Response<TransactionResponse>

    @GET("transactions/history")
    suspend fun getHistory(
        @Query("limit") limit: Int = 20,
        @Query("offset") offset: Int = 0
    ): Response<List<TransactionResponse>>
}
