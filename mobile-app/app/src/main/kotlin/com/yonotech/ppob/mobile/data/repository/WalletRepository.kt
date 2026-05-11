package com.yonotech.ppob.mobile.data.repository

import com.yonotech.ppob.mobile.data.remote.WalletService
import com.yonotech.ppob.mobile.data.remote.dto.WalletResponse
import com.yonotech.ppob.mobile.utils.Resource
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class WalletRepository @Inject constructor(
    private val walletService: WalletService
) {
    suspend fun getBalance(): Resource<WalletResponse> {
        return try {
            val response = walletService.getBalance()
            if (response.isSuccessful) {
                Resource.Success(response.body()!!)
            } else {
                Resource.Error("Failed to get balance: ${response.message()}")
            }
        } catch (e: Exception) {
            Resource.Error(e.message ?: "Unknown error")
        }
    }
}