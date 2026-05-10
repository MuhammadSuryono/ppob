package com.yonotech.ppob.data.local.entity

import androidx.room.ColumnInfo
import androidx.room.Entity
import androidx.room.PrimaryKey

@Entity(tableName = "user_profile")
data class UserEntity(
    @PrimaryKey
    @ColumnInfo(name = "user_id") val userId: String,
    @ColumnInfo(name = "phone_number") val phoneNumber: String,
    @ColumnInfo(name = "name") val name: String,
    @ColumnInfo(name = "active_role") val activeRole: String,
    @ColumnInfo(name = "wallet_id") val walletId: String,
    @ColumnInfo(name = "access_token") val accessToken: String,
    @ColumnInfo(name = "refresh_token") val refreshToken: String,
    @ColumnInfo(name = "token_expires_at") val tokenExpiresAt: Long,
    @ColumnInfo(name = "updated_at") val updatedAt: Long = System.currentTimeMillis()
)