package com.yonotech.ppob.domain.model

data class Transaction(
    val id: String,
    val refId: String?,
    val productName: String,
    val customerNumber: String,
    val status: String,
    val sellingPrice: Double,
    val platformPrice: Double,
    val marginAmount: Double?,
    val commissionAmount: Double?,
    val createdAt: Long,
    val completedAt: Long? = null
)

data class Product(
    val productId: String,
    val buyerSkuCode: String,
    val productName: String,
    val basePrice: Double,
    val platformPrice: Double,
    val isPrepaid: Boolean,
    val mitraSellingPrice: Double? = null,
    val categoryId: String,
    val iconUrl: String? = null
)

data class Category(
    val categoryId: String,
    val categoryName: String,
    val iconUrl: String,
    val displayOrder: Int
)

data class WalletBalance(
    val walletId: String,
    val balanceAvailable: Double,
    val balanceHeld: Double,
    val currency: String = "IDR"
)

data class Staff(
    val userId: String,
    val name: String,
    val phoneNumber: String,
    val walletBalance: Double,
    val dailyTxnCount: Int,
    val dailyTxnAmount: Double,
    val marginScheme: String,
    val marginValue: Double,
    val isActive: Boolean
)

enum class TransactionStatus {
    SUCCESS, PENDING, FAILED, CANCELLED
}

enum class MarginScheme {
    FIXED_ALLOWANCE, MARGIN_SHARE
}