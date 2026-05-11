package com.yonotech.ppob.mobile.data.remote.dto.transaction

import com.squareup.moshi.Json
import com.squareup.moshi.JsonClass

@JsonClass(generateAdapter = true)
data class InitiateTransactionRequest(
    @Json(name = "product_id") val productId: String,
    @Json(name = "customer_no") val customerNo: String,
    @Json(name = "selling_price") val sellingPrice: Double? = null,
    @Json(name = "pin") val pin: String
)

@JsonClass(generateAdapter = true)
data class TransactionResponse(
    @Json(name = "id") val id: String,
    @Json(name = "transaction_id") val transactionId: String,
    @Json(name = "status") val status: String,
    @Json(name = "message") val message: String? = null
)