package com.yonotech.ppob.data.local.dao

import androidx.room.Dao
import androidx.room.Insert
import androidx.room.OnConflictStrategy
import androidx.room.Query
import androidx.room.Update
import com.yonotech.ppob.data.local.entity.TransactionEntity
import kotlinx.coroutines.flow.Flow

@Dao
interface TransactionDao {

    @Query("SELECT * FROM transactions ORDER BY created_at DESC LIMIT :limit OFFSET :offset")
    fun getTransactions(limit: Int = 20, offset: Int = 0): Flow<List<TransactionEntity>>

    @Query("SELECT * FROM transactions WHERE id = :transactionId")
    suspend fun getById(transactionId: String): TransactionEntity?

    @Query("SELECT * FROM transactions ORDER BY created_at DESC")
    fun getAll(): Flow<List<TransactionEntity>>

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun insert(transaction: TransactionEntity)

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun insertAll(transactions: List<TransactionEntity>)

    @Query("DELETE FROM transactions WHERE created_at < :cutoff")
    suspend fun deleteOld(cutoff: Long = System.currentTimeMillis() - 90L * 86400000)
}