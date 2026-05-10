package com.yonotech.ppob.data.local.entity

import androidx.room.ColumnInfo
import androidx.room.Entity
import androidx.room.PrimaryKey
import androidx.room.Index

@Entity(
    tableName = "transactions",
    indices = [
        Index(value = ["created_at"]),
        Index(value = ["status"])
    ]
)
data class TransactionEntity(
    @PrimaryKey
    @ColumnInfo(name = "id") val id: String,
    @ColumnInfo(name = "ref_id") val refId: String? = null,
    @ColumnInfo(name = "product_name") val productName: String,
    @ColumnInfo(name = "customer_number") val customerNumber: String,
    @ColumnInfo(name = "status") val status: String, // Success, Pending, Failed, Cancelled
    @ColumnInfo(name = "selling_price") val sellingPrice: Double,
    @ColumnInfo(name = "platform_price") val platformPrice: Double,
    @ColumnInfo(name = "margin_amount") val marginAmount: Double? = null,
    @ColumnInfo(name = "commission_amount") val commissionAmount: Double? = null,
    @ColumnInfo(name = "created_at") val createdAt: Long,
    @ColumnInfo(name = "completed_at") val completedAt: Long? = null
)