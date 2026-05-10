package com.yonotech.ppob.data.local.entity

import androidx.room.ColumnInfo
import androidx.room.Entity
import androidx.room.PrimaryKey

@Entity(tableName = "category")
data class CategoryEntity(
    @PrimaryKey
    @ColumnInfo(name = "category_id") val categoryId: String,
    @ColumnInfo(name = "category_name") val categoryName: String,
    @ColumnInfo(name = "icon_url") val iconUrl: String,
    @ColumnInfo(name = "display_order") val displayOrder: Int,
    @ColumnInfo(name = "synced_at") val syncedAt: Long = System.currentTimeMillis()
)