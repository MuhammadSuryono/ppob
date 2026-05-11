package com.yonotech.ppob.mobile.data.remote.dto

import com.squareup.moshi.Json
import com.squareup.moshi.JsonClass

@JsonClass(generateAdapter = true)
data class CategoryDto(
    @Json(name = "id") val id: String,
    @Json(name = "name") val name: String,
    @Json(name = "icon_url") val iconUrl: String? = null
)

@JsonClass(generateAdapter = true)
data class ProductDto(
    @Json(name = "id") val id: String,
    @Json(name = "name") val name: String,
    @Json(name = "code") val code: String,
    @Json(name = "category_id") val categoryId: String,
    @Json(name = "brand") val brand: String,
    @Json(name = "price") val price: Double,
    @Json(name = "description") val description: String? = null,
    @Json(name = "status") val status: String
)
