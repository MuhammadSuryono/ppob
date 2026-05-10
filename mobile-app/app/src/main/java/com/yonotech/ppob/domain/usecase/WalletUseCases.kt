package com.yonotech.ppob.domain.usecase

import com.yonotech.ppob.data.remote.model.TopUpResponse
import com.yonotech.ppob.domain.model.WalletBalance
import com.yonotech.ppob.domain.repository.WalletRepository
import java.math.BigDecimal

class GetBalanceUseCase(private val repository: WalletRepository) {
    suspend operator fun invoke(token: String): Result<WalletBalance> {
        return repository.getBalance(token)
    }

    suspend fun getCachedBalance(): WalletBalance? {
        return repository.getCachedBalance()
    }
}

class TopUpStaffUseCase(private val repository: WalletRepository) {
    suspend operator fun invoke(token: String, staffId: String, amount: BigDecimal): Result<TopUpResponse> {
        return repository.topUpStaff(token, staffId, amount.toDouble())
    }
}

class ValidateBalanceUseCase(private val repository: WalletRepository) {
    suspend operator fun invoke(balance: Double, required: Double): Boolean {
        return balance >= required
    }
}
