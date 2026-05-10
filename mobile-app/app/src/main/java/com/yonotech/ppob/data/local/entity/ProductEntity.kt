package com.yonotech.ppob.data.local.entity

import androidx.room.ColumnInfo
import androidx.room.Entity
import androidx.room.PrimaryKey

@Entity(tableName = "product")
data class ProductEntity(
    @PrimaryKey
    @ColumnInfo(name = "product_id") val productId: String,
    @ColumnInfo(name = "buyer_sku_code") val buyerSkuCode: String,
    @ColumnInfo(name = "product_name") val productName: String,
    @ColumnInfo(name = "base_price") val basePrice: Double,
    @ColumnInfo(name = "platform_price") val platformPrice: Double,
    @ColumnInfo(name = "is_prepaid") val isPrepaid: Boolean,
    @ColumnInfo(name = "mitra_selling_price") val mitraSellingPrice: Double? = null,
    @ColumnInfo(name = "category_id") val categoryId: String,
    @ColumnInfo(name = "icon_url") val iconUrl: String? = null,
    @ColumnInfo(name = "synced_at") val syncedAt: Long = System.currentTimeMillis()
)