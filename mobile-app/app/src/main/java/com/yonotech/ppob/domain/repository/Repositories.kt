package com.yonotech.ppob.domain.repository

import com.yonotech.ppob.data.remote.model.*
import com.yonotech.ppob.domain.model.Category as DomainCategory
import com.yonotech.ppob.domain.model.Product as DomainProduct
import com.yonotech.ppob.domain.model.Transaction as DomainTransaction
import com.yonotech.ppob.domain.model.WalletBalance

interface AuthRepository {
    suspend fun register(request: RegisterRequest): Result<AuthResponse>
    suspend fun login(request: LoginRequest): Result<AuthResponse>
    suspend fun verifyOtp(request: VerifyOtpRequest): Result<AuthResponse>
    suspend fun refreshToken(refreshToken: String): Result<AuthResponse>
    suspend fun logout(token: String): Result<Unit>
    suspend fun changePin(token: String, request: ChangePinRequest): Result<Unit>
    suspend fun getStoredToken(): String?
    suspend fun clearTokens()
}

interface UserRepository {
    suspend fun getProfile(token: String): Result<UserProfileResponse>
    suspend fun switchRole(token: String, roleId: String): Result<SwitchRoleResponse>
    suspend fun getStaff(token: String, limit: Int, offset: Int): Result<StaffListResponse>
    suspend fun addStaff(token: String, request: AddStaffRequest): Result<StaffDetailResponse>
    suspend fun updateStaff(token: String, staffId: String, request: UpdateStaffRequest): Result<StaffDetailResponse>
    suspend fun getStaffDetail(token: String, staffId: String): Result<StaffDetailResponse>
}

interface WalletRepository {
    suspend fun getBalance(token: String): Result<WalletBalance>
    suspend fun getCachedBalance(): WalletBalance?
    suspend fun topUpStaff(token: String, staffId: String, amount: Double): Result<TopUpResponse>
}

interface ProductRepository {
    suspend fun getCategories(): Result<List<DomainCategory>>
    suspend fun getProducts(categoryId: String, limit: Int, cursor: String?): Result<ProductsResponse>
    suspend fun getCachedProducts(categoryId: String): List<DomainProduct>
    suspend fun syncProducts(token: String): Result<Unit>
}

interface TransactionRepository {
    suspend fun initiateTransaction(
        token: String,
        idempotencyKey: String,
        productId: String,
        customerNo: String,
        sellingPrice: Double,
        pin: String
    ): Result<TransactionResponse>

    suspend fun getTransactionHistory(
        token: String,
        limit: Int,
        cursor: String?,
        status: String?
    ): Result<TransactionHistoryResponse>

    suspend fun getTransactionStatus(token: String, transactionId: String): Result<TransactionDetailResponse>
    suspend fun getCachedTransactions(): List<DomainTransaction>
    suspend fun queuePendingTransaction(payload: String): String
    suspend fun updateTransactionStatus(id: String, status: String)
}
