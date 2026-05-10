package com.yonotech.ppob.data.local.entity

import androidx.room.ColumnInfo
import androidx.room.Entity
import androidx.room.PrimaryKey

@Entity(tableName = "pending_sync_queue")
data class PendingSyncItem(
    @PrimaryKey
    @ColumnInfo(name = "id") val id: String = java.util.UUID.randomUUID().toString(),
    @ColumnInfo(name = "type") val type: String, // "transaction_initiate", "top_up", etc.
    @ColumnInfo(name = "payload") val payload: String, // JSON payload
    @ColumnInfo(name = "retry_count") var retryCount: Int = 0,
    @ColumnInfo(name = "last_attempt") var lastAttempt: Long = 0L,
    @ColumnInfo(name = "status") val status: SyncStatus = SyncStatus.PENDING
)

enum class SyncStatus {
    PENDING,
    RETRYING,
    FAILED,
    COMPLETED
}