package com.yonotech.ppob.mobile.data.remote.dto.transaction

import com.squareup.moshi.Json
import com.squareup.moshi.JsonClass

@JsonClass(generateAdapter = true)
data class InitiateTransactionRequest(
    @Json(name = "product_code") val productCode: String,
    @Json(name = "customer_number") val customerNumber: String,
    @Json(name = "amount") val amount: Double,
    @Json(name = "authorize_id") val authorizeId: String
)

@JsonClass(generateAdapter = true)
data class TransactionResponse(
    @Json(name = "id") val id: Long? = null,
    @Json(name = "transaction_id") val transactionId: String,
    @Json(name = "product_code") val productCode: String? = null,
    @Json(name = "customer_number") val customerNumber: String? = null,
    @Json(name = "amount") val amount: Double? = null,
    @Json(name = "price") val price: Double? = null,
    @Json(name = "selling_price") val sellingPrice: Double? = null,
    @Json(name = "status") val status: String,
    @Json(name = "message") val message: String? = null,
    @Json(name = "created_at") val createdAt: String? = null
)

@JsonClass(generateAdapter = true)
data class InquiryRequest(
    @Json(name = "category_id") val categoryId: Long,
    @Json(name = "brand") val brand: String,
    @Json(name = "customer_number") val customerNumber: String,
    @Json(name = "product_code") val productCode: String? = null
)

@JsonClass(generateAdapter = true)
data class InquiryResponse(
    @Json(name = "inquiry_id") val inquiryId: String,
    @Json(name = "customer_number") val customerNumber: String,
    @Json(name = "customer_name") val customerName: String,
    @Json(name = "bill_amount") val billAmount: Double = 0.0,
    @Json(name = "admin_fee") val adminFee: Double = 0.0,
    @Json(name = "total_amount") val totalAmount: Double = 0.0,
    @Json(name = "description") val description: String? = null,
    @Json(name = "is_postpaid") val isPostpaid: Boolean = false,
    @Json(name = "product_code") val productCode: String? = null
)