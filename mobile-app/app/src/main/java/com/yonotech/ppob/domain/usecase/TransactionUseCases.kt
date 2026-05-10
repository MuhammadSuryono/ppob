package com.yonotech.ppob.domain.usecase

import com.yonotech.ppob.data.remote.model.ApiError
import com.yonotech.ppob.domain.model.Product
import com.yonotech.ppob.domain.model.WalletBalance
import com.yonotech.ppob.domain.repository.AuthRepository
import com.yonotech.ppob.domain.repository.TransactionRepository
import kotlinx.coroutines.flow.Flow

// ========== Auth Use Cases ==========

class RegisterUseCase(private val repository: AuthRepository) {
    suspend operator fun invoke(phone: String, name: String): Result<Unit> {
        return repository.register(
            com.yonotech.ppob.data.remote.model.RegisterRequest(phone, name)
        ).mapSuccess { }
    }
}

class LoginUseCase(private val repository: AuthRepository) {
    suspend operator fun invoke(
        phone: String, password: String, pin: String,
        deviceFingerprint: com.yonotech.ppob.data.remote.model.DeviceFingerprint? = null,
        otpCode: String? = null
    ): Result<Unit> {
        return repository.login(
            com.yonotech.ppob.data.remote.model.LoginRequest(phone, password, pin, deviceFingerprint, otpCode)
        ).mapSuccess { }
    }
}

class VerifyOtpUseCase(private val repository: AuthRepository) {
    suspend operator fun invoke(phone: String, otp: String, password: String, pin: String): Result<Unit> {
        return repository.verifyOtp(
            com.yonotech.ppob.data.remote.model.VerifyOtpRequest(phone, otp, password, pin)
        ).mapSuccess { }
    }
}

class RefreshTokenUseCase(private val repository: AuthRepository) {
    suspend operator fun invoke(refreshToken: String): Result<Unit> {
        return repository.refreshToken(refreshToken).mapSuccess { }
    }
}

class LogoutUseCase(private val repository: AuthRepository) {
    suspend operator fun invoke(token: String): Result<Unit> {
        return repository.logout(token)
    }
}

class ChangePinUseCase(private val repository: AuthRepository) {
    suspend operator fun invoke(
        token: String,
        currentPin: String? = null,
        password: String? = null,
        otpCode: String? = null,
        newPin: String,
        confirmPin: String
    ): Result<Unit> {
        val request = com.yonotech.ppob.data.remote.model.ChangePinRequest(
            currentPin, password, otpCode, newPin, confirmPin
        )
        return repository.changePin(token, request)
    }
}

// ========== Transaction Use Cases ==========

class InitiateTransactionUseCase(private val repository: TransactionRepository) {
    suspend operator fun invoke(
        token: String,
        productId: String,
        customerNo: String,
        sellingPrice: Double,
        pin: String
    ): Result<com.yonotech.ppob.data.remote.model.TransactionResponse> {
        val key = java.util.UUID.randomUUID().toString()
        return repository.initiateTransaction(token, key, productId, customerNo, sellingPrice, pin)
    }
}

class GetTransactionHistoryUseCase(private val repository: TransactionRepository) {
    suspend operator fun invoke(
        token: String,
        limit: Int = 20,
        cursor: String? = null,
        status: String? = null
    ): Result<com.yonotech.ppob.data.remote.model.TransactionHistoryResponse> {
        return repository.getTransactionHistory(token, limit, cursor, status)
    }

    fun getCachedTransactions(): Flow<List<com.yonotech.ppob.domain.model.Transaction>> {
        // TODO: Convert to Flow
        return kotlinx.coroutines.flow.flow {
            emit(repository.getCachedTransactions())
        }
    }
}

class GetTransactionStatusUseCase(private val repository: TransactionRepository) {
    suspend operator fun invoke(token: String, transactionId: String): Result<com.yonotech.ppob.data.remote.model.TransactionDetailResponse> {
        return repository.getTransactionStatus(token, transactionId)
    }
}

// ========== Common ==========

fun <T> Result<T>.mapSuccess(mapper: (T) -> Unit): Result<Unit> {
    return this.mapCatching { mapper(it) }
}