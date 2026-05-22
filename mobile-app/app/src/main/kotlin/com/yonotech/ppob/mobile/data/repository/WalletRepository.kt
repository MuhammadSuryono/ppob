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
    suspend fun getBalance(walletId: String): retrofit2.Response<WalletResponse> {
        return walletService.getBalance(walletId)
    }
}