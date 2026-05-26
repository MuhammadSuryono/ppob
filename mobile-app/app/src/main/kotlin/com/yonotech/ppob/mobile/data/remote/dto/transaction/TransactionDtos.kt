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