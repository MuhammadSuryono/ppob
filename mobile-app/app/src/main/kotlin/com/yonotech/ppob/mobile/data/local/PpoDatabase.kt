package com.yonotech.ppob.mobile.data.local

import androidx.room.Database
import androidx.room.RoomDatabase
import androidx.room.TypeConverters
import com.yonotech.ppob.mobile.data.local.dao.CategoryDao
import com.yonotech.ppob.mobile.data.local.dao.ProductDao
import com.yonotech.ppob.mobile.data.local.dao.TransactionDao
import com.yonotech.ppob.mobile.data.local.dao.PendingSyncDao
import com.yonotech.ppob.mobile.data.local.entities.CategoryEntity
import com.yonotech.ppob.mobile.data.local.entities.ProductEntity
import com.yonotech.ppob.mobile.data.local.entities.TransactionEntity
import com.yonotech.ppob.mobile.data.local.entities.PendingSyncItem

@Database(
    entities = [
        CategoryEntity::class,
        ProductEntity::class,
        TransactionEntity::class,
        PendingSyncItem::class
    ],
    version = 1,
    exportSchema = false
)
@TypeConverters(Converters::class)
abstract class PpoDatabase : RoomDatabase() {
    abstract fun categoryDao(): CategoryDao
    abstract fun productDao(): ProductDao
    abstract fun transactionDao(): TransactionDao
    abstract fun pendingSyncDao(): PendingSyncDao
}
