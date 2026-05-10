package com.yonotech.ppob.data.local.database

import android.content.Context
import androidx.room.Database
import androidx.room.Room
import androidx.room.RoomDatabase
import com.yonotech.ppob.data.local.dao.CategoryDao
import com.yonotech.ppob.data.local.dao.PendingSyncDao
import com.yonotech.ppob.data.local.dao.ProductDao
import com.yonotech.ppob.data.local.dao.TransactionDao
import com.yonotech.ppob.data.local.dao.UserDao
import com.yonotech.ppob.data.local.dao.WalletDao
import com.yonotech.ppob.data.local.entity.CategoryEntity
import com.yonotech.ppob.data.local.entity.PendingSyncItem
import com.yonotech.ppob.data.local.entity.ProductEntity
import com.yonotech.ppob.data.local.entity.TransactionEntity
import com.yonotech.ppob.data.local.entity.UserEntity
import com.yonotech.ppob.data.local.entity.WalletEntity

@Database(
    entities = [
        UserEntity::class,
        WalletEntity::class,
        TransactionEntity::class,
        ProductEntity::class,
        CategoryEntity::class,
        PendingSyncItem::class
    ],
    version = 1,
    exportSchema = true
)
abstract class PpoDatabase : RoomDatabase() {

    abstract fun userDao(): UserDao
    abstract fun walletDao(): WalletDao
    abstract fun transactionDao(): TransactionDao
    abstract fun productDao(): ProductDao
    abstract fun categoryDao(): CategoryDao
    abstract fun pendingSyncDao(): PendingSyncDao

    companion object {
        @Volatile
        private var INSTANCE: PpoDatabase? = null

        fun getInstance(context: Context): PpoDatabase {
            return INSTANCE ?: synchronized(this) {
                Room.databaseBuilder(
                    context.applicationContext,
                    PpoDatabase::class.java,
                    "ppob_database"
                )
                    .fallbackToDestructiveMigration()
                    .build()
                    .also { INSTANCE = it }
            }
        }
    }
}