package com.yonotech.ppob.mobile.data.repository

import com.yonotech.ppob.mobile.data.remote.TransactionService
import com.yonotech.ppob.mobile.data.remote.dto.transaction.InitiateTransactionRequest
import com.yonotech.ppob.mobile.data.remote.dto.transaction.InquiryRequest
import com.yonotech.ppob.mobile.data.remote.dto.transaction.InquiryResponse
import com.yonotech.ppob.mobile.data.remote.dto.transaction.TransactionResponse
import retrofit2.Response
import java.util.UUID
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class TransactionRepository @Inject constructor(
    private val transactionService: TransactionService
) {
    suspend fun initiate(request: InitiateTransactionRequest): Response<TransactionResponse> {
        val idempotencyKey = UUID.randomUUID().toString()
        return transactionService.initiate(idempotencyKey, request)
    }

    suspend fun inquiry(request: InquiryRequest): Response<InquiryResponse> {
        return transactionService.inquiry(request)
    }

    suspend fun getStatus(id: String): Response<TransactionResponse> {
        return transactionService.getStatus(id)
    }
}
