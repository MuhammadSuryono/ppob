package com.yonotech.ppob.mobile.data.remote.dto

import com.squareup.moshi.Json
import com.squareup.moshi.JsonClass

@JsonClass(generateAdapter = true)
data class WalletResponse(
    @Json(name = "id") val id: String,
    @Json(name = "balance_available") val balanceAvailable: Double,
    @Json(name = "balance_held") val balanceHeld: Double,
    @Json(name = "balance_total") val balanceTotal: Double,
    @Json(name = "updated_at") val updatedAt: String? = null
)

@JsonClass(generateAdapter = true)
data class TopUpRequest(
    @Json(name = "staff_id") val staffId: String,
    @Json(name = "amount") val amount: Double,
    @Json(name = "pin") val pin: String
)

@JsonClass(generateAdapter = true)
data class StaffDto(
    @Json(name = "id") val id: String,
    @Json(name = "name") val name: String,
    @Json(name = "phone") val phone: String,
    @Json(name = "email") val email: String,
    @Json(name = "balance") val balance: Double,
    @Json(name = "daily_limit") val dailyLimit: Double,
    @Json(name = "daily_used") val dailyUsed: Double,
    @Json(name = "margin_scheme") val marginScheme: String,
    @Json(name = "margin_value") val marginValue: Double,
    @Json(name = "is_active") val isActive: Boolean
)

@JsonClass(generateAdapter = true)
data class CreateStaffRequest(
    @Json(name = "name") val name: String,
    @Json(name = "phone") val phone: String,
    @Json(name = "email") val email: String,
    @Json(name = "margin_scheme") val marginScheme: String,
    @Json(name = "margin_value") val marginValue: Double,
    @Json(name = "daily_limit") val dailyLimit: Double
)

@JsonClass(generateAdapter = true)
data class TransactionHistoryResponse(
    @Json(name = "id") val id: String,
    @Json(name = "product_name") val productName: String,
    @Json(name = "customer_number") val customerNumber: String,
    @Json(name = "selling_price") val sellingPrice: Double,
    @Json(name = "platform_price") val platformPrice: Double,
    @Json(name = "status") val status: String,
    @Json(name = "created_at") val createdAt: String,
    @Json(name = "brand") val brand: String? = null
)