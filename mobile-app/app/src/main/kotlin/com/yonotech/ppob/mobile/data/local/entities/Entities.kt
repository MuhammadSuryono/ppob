package com.yonotech.ppob.mobile.data.local.entities

import androidx.room.Entity
import androidx.room.PrimaryKey

@Entity(tableName = "categories")
data class CategoryEntity(
    @PrimaryKey val id: String,
    val name: String,
    val iconUrl: String? = null
)

@Entity(tableName = "products")
data class ProductEntity(
    @PrimaryKey val id: String,
    val name: String,
    val code: String,
    val categoryId: String,
    val brand: String,
    val price: Double,
    val description: String? = null,
    val status: String,
    val lastSync: Long
)

@Entity(tableName = "transactions")
data class TransactionEntity(
    @PrimaryKey val id: String,
    val transactionId: String,
    val productName: String,
    val sellingPrice: Double,
    val platformPrice: Double,
    val customerNumber: String,
    val status: String,
    val createdAt: Long,
    val brand: String? = null,
    val categoryId: String? = null
)

@Entity(tableName = "pending_sync_queue")
data class PendingSyncItem(
    @PrimaryKey val id: String,
    val type: String,
    val payload: String,
    var retryCount: Int = 0,
    var lastAttempt: Long = 0L,
    var status: String = "PENDING"
)
