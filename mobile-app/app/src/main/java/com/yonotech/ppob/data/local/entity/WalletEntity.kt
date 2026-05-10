package com.yonotech.ppob.data.local.entity

import androidx.room.ColumnInfo
import androidx.room.Entity
import androidx.room.PrimaryKey

@Entity(tableName = "wallet")
data class WalletEntity(
    @PrimaryKey
    @ColumnInfo(name = "wallet_id") val walletId: String,
    @ColumnInfo(name = "balance_available") val balanceAvailable: Double,
    @ColumnInfo(name = "balance_held") val balanceHeld: Double,
    @ColumnInfo(name = "currency") val currency: String = "IDR",
    @ColumnInfo(name = "updated_at") val updatedAt: Long = System.currentTimeMillis()
)