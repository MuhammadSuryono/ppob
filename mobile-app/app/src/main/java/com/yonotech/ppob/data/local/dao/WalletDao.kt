package com.yonotech.ppob.data.local.dao

import androidx.room.Dao
import androidx.room.Insert
import androidx.room.OnConflictStrategy
import androidx.room.Query
import androidx.room.Update
import com.yonotech.ppob.data.local.entity.WalletEntity
import kotlinx.coroutines.flow.Flow

@Dao
interface WalletDao {

    @Query("SELECT * FROM wallet LIMIT 1")
    fun getWallet(): Flow<WalletEntity?>

    @Insert(onConflict = OnConflictStrategy.REPLACE)
    suspend fun insert(wallet: WalletEntity)

    @Update
    suspend fun update(wallet: WalletEntity)

    @Query("UPDATE wallet SET balance_available = :available, balance_held = :held, updated_at = :updatedAt")
    suspend fun updateBalance(available: Double, held: Double, updatedAt: Long = System.currentTimeMillis())
}