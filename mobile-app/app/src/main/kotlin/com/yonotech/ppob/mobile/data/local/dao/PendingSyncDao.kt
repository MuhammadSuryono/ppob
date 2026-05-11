package com.yonotech.ppob.mobile.data.local.dao

import androidx.room.Dao
import androidx.room.Insert
import androidx.room.OnConflictStrategy
import androidx.room.Query
import com.yonotech.ppob.mobile.data.local.entities.PendingSyncItem

@Dao
interface PendingSyncDao {
    @Query("SELECT * FROM pending_sync_queue WHERE status = 'PENDING' ORDER BY lastAttempt ASC")
    suspend fun getPendingItems(): List<PendingSyncItem>

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun insert(item: PendingSyncItem)

    @Query("UPDATE pending_sync_queue SET status = :status, lastAttempt = :timestamp, retryCount = retryCount + 1 WHERE id = :id")
    suspend fun updateStatus(id: String, status: String, timestamp: Long)

    @Query("DELETE FROM pending_sync_queue WHERE status = 'COMPLETED'")
    suspend fun clearCompleted()
}