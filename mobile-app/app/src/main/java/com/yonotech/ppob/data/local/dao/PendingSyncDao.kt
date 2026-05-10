package com.yonotech.ppob.data.local.dao

import androidx.room.Dao
import androidx.room.Insert
import androidx.room.OnConflictStrategy
import androidx.room.Query
import androidx.room.Update
import com.yonotech.ppob.data.local.entity.PendingSyncItem
import kotlinx.coroutines.flow.Flow

@Dao
interface PendingSyncDao {

    @Query("SELECT * FROM pending_sync_queue WHERE status = 'PENDING' OR status = 'RETRYING' ORDER BY id ASC")
    fun getAllPending(): Flow<List<PendingSyncItem>>

    @Query("SELECT * FROM pending_sync_queue WHERE id = :id")
    suspend fun getById(id: String): PendingSyncItem?

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun insert(item: PendingSyncItem)

    @Update
    suspend fun update(item: PendingSyncItem)

    @Query("DELETE FROM pending_sync_queue WHERE id = :id")
    suspend fun delete(id: String)

    @Query("DELETE FROM pending_sync_queue WHERE status = 'COMPLETED' OR status = 'FAILED'")
    suspend fun cleanupCompleted()
}