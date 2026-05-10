package com.yonotech.ppob.data.repository

import android.util.Log
import com.yonotech.ppob.data.local.dao.PendingSyncDao
import com.yonotech.ppob.data.local.dao.TransactionDao
import com.yonotech.ppob.data.local.entity.TransactionEntity
import com.yonotech.ppob.data.remote.api.TransactionApiService
import com.yonotech.ppob.data.remote.model.InitiateTransactionRequest
import com.yonotech.ppob.data.remote.model.TransactionHistoryResponse
import com.yonotech.ppob.data.remote.model.TransactionSummary
import com.yonotech.ppob.domain.model.Transaction
import com.yonotech.ppob.domain.repository.TransactionRepository
import kotlinx.coroutines.flow.first
import java.text.SimpleDateFormat
import java.util.Locale

class TransactionRepositoryImpl(
    private val apiService: TransactionApiService,
    private val transactionDao: TransactionDao,
    private val pendingSyncDao: PendingSyncDao
) : TransactionRepository {

    private val dateFormat = SimpleDateFormat("yyyy-MM-dd'T'HH:mm:ss'Z'", Locale.getDefault())

    override suspend fun initiateTransaction(
        token: String,
        idempotencyKey: String,
        productId: String,
        customerNo: String,
        sellingPrice: Double,
        pin: String
    ): Result<com.yonotech.ppob.data.remote.model.TransactionResponse> {
        return try {
            val request = InitiateTransactionRequest(productId, customerNo, sellingPrice, pin)
            val response = apiService.initiateTransaction(
                "Bearer $token",
                idempotencyKey,
                request
            )
            if (response.isSuccessful) {
                val body = response.body()
                if ((body?.success == true || body?.data != null) && body?.data != null) {
                    Result.success(body.data)
                } else {
                    Result.failure(Exception(body?.message ?: "Transaction initiation failed"))
                }
            } else {
                Result.failure(Exception("HTTP ${response.code()}"))
            }
        } catch (e: Exception) {
            Log.e("TransactionRepository", "Initiate transaction error", e)
            Result.failure(e)
        }
    }

    override suspend fun getTransactionHistory(
        token: String,
        limit: Int,
        cursor: String?,
        status: String?
    ): Result<TransactionHistoryResponse> {
        return try {
            val response = apiService.getTransactionHistory("Bearer $token", limit, cursor, status)
            if (response.isSuccessful) {
                val body = response.body()
                if (body?.success == true && body.data != null) {
                    // Cache transactions locally
                    cacheTransactions(body.data.transactions)
                    Result.success(body.data)
                } else {
                    Result.failure(Exception(body?.message ?: "Failed to get history"))
                }
            } else {
                // Return cached on error
                val cached = getCachedTransactions()
                if (cached.isNotEmpty()) {
                    val history = TransactionHistoryResponse(
                        transactions = cached.map {
                            TransactionSummary(
                                id = it.id,
                                refId = it.refId,
                                productName = it.productName,
                                customerNumber = it.customerNumber,
                                status = it.status,
                                sellingPrice = it.sellingPrice,
                                platformPrice = it.platformPrice,
                                marginAmount = it.marginAmount,
                                commissionAmount = it.commissionAmount,
                                createdAt = dateFormat.format(it.createdAt),
                                completedAt = if (it.completedAt != null) dateFormat.format(it.completedAt) else null
                            )
                        },
                        pagination = null
                    )
                    Result.success(history)
                } else {
                    Result.failure(Exception("HTTP ${response.code()}"))
                }
            }
        } catch (e: Exception) {
            Log.e("TransactionRepository", "Get transaction history error", e)
            val cached = getCachedTransactions()
            if (cached.isNotEmpty()) {
                val history = TransactionHistoryResponse(
                    transactions = cached.map {
                        TransactionSummary(
                            id = it.id,
                            refId = it.refId,
                            productName = it.productName,
                            customerNumber = it.customerNumber,
                            status = it.status,
                            sellingPrice = it.sellingPrice,
                            platformPrice = it.platformPrice,
                            marginAmount = it.marginAmount,
                            commissionAmount = it.commissionAmount,
                            createdAt = dateFormat.format(it.createdAt),
                            completedAt = if (it.completedAt != null) dateFormat.format(it.completedAt) else null
                        )
                    },
                    pagination = null
                )
                Result.success(history)
            } else {
                Result.failure(e)
            }
        }
    }

    override suspend fun getTransactionStatus(token: String, transactionId: String): Result<com.yonotech.ppob.data.remote.model.TransactionDetailResponse> {
        return try {
            val response = apiService.getTransactionStatus("Bearer $token", transactionId)
            if (response.isSuccessful) {
                val body = response.body()
                if (body?.success == true && body.data != null) {
                    Result.success(body.data)
                } else {
                    Result.failure(Exception(body?.message ?: "Failed to get transaction status"))
                }
            } else {
                Result.failure(Exception("HTTP ${response.code()}"))
            }
        } catch (e: Exception) {
            Log.e("TransactionRepository", "Get transaction status error", e)
            Result.failure(e)
        }
    }

    override suspend fun getCachedTransactions(): List<Transaction> {
        return try {
            transactionDao.getAll().first().map { entity ->
                Transaction(
                    id = entity.id,
                    refId = entity.refId,
                    productName = entity.productName,
                    customerNumber = entity.customerNumber,
                    status = entity.status,
                    sellingPrice = entity.sellingPrice,
                    platformPrice = entity.platformPrice,
                    marginAmount = entity.marginAmount,
                    commissionAmount = entity.commissionAmount,
                    createdAt = entity.createdAt,
                    completedAt = entity.completedAt
                )
            }
        } catch (e: Exception) {
            Log.e("TransactionRepository", "Get cached transactions error", e)
            emptyList()
        }
    }

    override suspend fun queuePendingTransaction(payload: String): String {
        val item = com.yonotech.ppob.data.local.entity.PendingSyncItem(
            type = "transaction_initiate",
            payload = payload
        )
        pendingSyncDao.insert(item)
        return item.id
    }

    override suspend fun updateTransactionStatus(id: String, status: String) {
        try {
            transactionDao.getById(id)?.let { entity ->
                val updated = entity.copy(status = status, completedAt = System.currentTimeMillis())
                transactionDao.insert(updated)
            }
        } catch (e: Exception) {
            Log.e("TransactionRepository", "Update transaction status error", e)
        }
    }

    private suspend fun cacheTransactions(summaries: List<TransactionSummary>) {
        try {
            val entities = summaries.map {
                TransactionEntity(
                    id = it.id,
                    refId = it.refId,
                    productName = it.productName,
                    customerNumber = it.customerNumber,
                    status = it.status,
                    sellingPrice = it.sellingPrice,
                    platformPrice = it.platformPrice,
                    marginAmount = it.marginAmount,
                    commissionAmount = it.commissionAmount,
                    createdAt = parseDate(it.createdAt),
                    completedAt = it.completedAt?.let { c -> parseDate(c) }
                )
            }
            transactionDao.insertAll(entities)
        } catch (e: Exception) {
            Log.e("TransactionRepository", "Cache transactions error", e)
        }
    }

    private fun parseDate(dateStr: String?): Long {
        if (dateStr.isNullOrEmpty()) return System.currentTimeMillis()
        return try {
            dateFormat.parse(dateStr)?.time ?: System.currentTimeMillis()
        } catch (e: Exception) {
            System.currentTimeMillis()
        }
    }
}
