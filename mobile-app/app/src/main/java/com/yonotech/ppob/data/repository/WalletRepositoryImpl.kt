package com.yonotech.ppob.data.repository

import android.util.Log
import com.yonotech.ppob.data.local.dao.WalletDao
import com.yonotech.ppob.data.local.entity.WalletEntity
import com.yonotech.ppob.data.remote.api.WalletApiService
import com.yonotech.ppob.data.remote.model.TopUpRequest
import com.yonotech.ppob.data.remote.model.TopUpResponse
import com.yonotech.ppob.data.remote.model.WalletBalanceResponse
import com.yonotech.ppob.domain.model.WalletBalance
import com.yonotech.ppob.domain.repository.WalletRepository
import kotlinx.coroutines.flow.first

class WalletRepositoryImpl(
    private val apiService: WalletApiService,
    private val walletDao: WalletDao
) : WalletRepository {

    override suspend fun getBalance(token: String): Result<WalletBalance> {
        return try {
            val response = apiService.getBalance("Bearer $token")
            if (response.isSuccessful) {
                val body = response.body()
                if (body?.success == true && body.data != null) {
                    val balance = body.data
                    // Cache locally
                    cacheBalance(balance)
                    Result.success(
                        WalletBalance(
                            walletId = balance.walletId,
                            balanceAvailable = balance.balanceAvailable,
                            balanceHeld = balance.balanceHeld,
                            currency = balance.currency
                        )
                    )
                } else {
                    // Try cached
                    val cached = getCachedBalance()
                    if (cached != null) {
                        Result.success(cached)
                    } else {
                        Result.failure(Exception(body?.message ?: "Failed to get balance"))
                    }
                }
            } else {
                // Try cached on error
                val cached = getCachedBalance()
                if (cached != null) {
                    Result.success(cached)
                } else {
                    Result.failure(Exception("HTTP ${response.code()}"))
                }
            }
        } catch (e: Exception) {
            Log.e("WalletRepository", "Get balance error", e)
            val cached = getCachedBalance()
            if (cached != null) {
                Result.success(cached)
            } else {
                Result.failure(e)
            }
        }
    }

    override suspend fun getCachedBalance(): WalletBalance? {
        return try {
            val entity = walletDao.getWallet().first()
            entity?.let {
                WalletBalance(
                    walletId = it.walletId,
                    balanceAvailable = it.balanceAvailable,
                    balanceHeld = it.balanceHeld,
                    currency = it.currency
                )
            }
        } catch (e: Exception) {
            Log.e("WalletRepository", "Get cached balance error", e)
            null
        }
    }

    override suspend fun topUpStaff(token: String, staffId: String, amount: Double): Result<TopUpResponse> {
        return try {
            val response = apiService.topUpStaff("Bearer $token", staffId, TopUpRequest(amount))
            if (response.isSuccessful) {
                val body = response.body()
                if (body?.success == true && body.data != null) {
                    // Update cached balances
                    body.data.mitraWallet?.let { cacheBalance(it) }
                    body.data.staffWallet?.let { cacheBalance(it) }
                    Result.success(body.data)
                } else {
                    Result.failure(Exception(body?.message ?: "Top-up failed"))
                }
            } else {
                Result.failure(Exception("HTTP ${response.code()}"))
            }
        } catch (e: Exception) {
            Log.e("WalletRepository", "Top-up error", e)
            Result.failure(e)
        }
    }

    private suspend fun cacheBalance(response: WalletBalanceResponse) {
        try {
            val entity = WalletEntity(
                walletId = response.walletId,
                balanceAvailable = response.balanceAvailable,
                balanceHeld = response.balanceHeld,
                currency = response.currency,
                updatedAt = System.currentTimeMillis()
            )
            walletDao.insert(entity)
        } catch (e: Exception) {
            Log.e("WalletRepository", "Cache balance error", e)
        }
    }
}
