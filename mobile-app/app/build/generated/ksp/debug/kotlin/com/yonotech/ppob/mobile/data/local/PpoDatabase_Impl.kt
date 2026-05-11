package com.yonotech.ppob.mobile.`data`.local

import androidx.room.InvalidationTracker
import androidx.room.RoomOpenDelegate
import androidx.room.migration.AutoMigrationSpec
import androidx.room.migration.Migration
import androidx.room.util.TableInfo
import androidx.room.util.TableInfo.Companion.read
import androidx.room.util.dropFtsSyncTriggers
import androidx.sqlite.SQLiteConnection
import androidx.sqlite.execSQL
import com.yonotech.ppob.mobile.`data`.local.dao.CategoryDao
import com.yonotech.ppob.mobile.`data`.local.dao.CategoryDao_Impl
import com.yonotech.ppob.mobile.`data`.local.dao.PendingSyncDao
import com.yonotech.ppob.mobile.`data`.local.dao.PendingSyncDao_Impl
import com.yonotech.ppob.mobile.`data`.local.dao.ProductDao
import com.yonotech.ppob.mobile.`data`.local.dao.ProductDao_Impl
import com.yonotech.ppob.mobile.`data`.local.dao.TransactionDao
import com.yonotech.ppob.mobile.`data`.local.dao.TransactionDao_Impl
import javax.`annotation`.processing.Generated
import kotlin.Lazy
import kotlin.String
import kotlin.Suppress
import kotlin.collections.List
import kotlin.collections.Map
import kotlin.collections.MutableList
import kotlin.collections.MutableMap
import kotlin.collections.MutableSet
import kotlin.collections.Set
import kotlin.collections.mutableListOf
import kotlin.collections.mutableMapOf
import kotlin.collections.mutableSetOf
import kotlin.reflect.KClass

@Generated(value = ["androidx.room.RoomProcessor"])
@Suppress(names = ["UNCHECKED_CAST", "DEPRECATION", "REDUNDANT_PROJECTION", "REMOVAL"])
public class PpoDatabase_Impl : PpoDatabase() {
  private val _categoryDao: Lazy<CategoryDao> = lazy {
    CategoryDao_Impl(this)
  }

  private val _productDao: Lazy<ProductDao> = lazy {
    ProductDao_Impl(this)
  }

  private val _transactionDao: Lazy<TransactionDao> = lazy {
    TransactionDao_Impl(this)
  }

  private val _pendingSyncDao: Lazy<PendingSyncDao> = lazy {
    PendingSyncDao_Impl(this)
  }

  protected override fun createOpenDelegate(): RoomOpenDelegate {
    val _openDelegate: RoomOpenDelegate = object : RoomOpenDelegate(1,
        "b247b74c89bfd7892a57e55dbba6715e", "930bc0b28cc4969a3f6ba7be68326f5c") {
      public override fun createAllTables(connection: SQLiteConnection) {
        connection.execSQL("CREATE TABLE IF NOT EXISTS `categories` (`id` TEXT NOT NULL, `name` TEXT NOT NULL, `iconUrl` TEXT, PRIMARY KEY(`id`))")
        connection.execSQL("CREATE TABLE IF NOT EXISTS `products` (`id` TEXT NOT NULL, `name` TEXT NOT NULL, `code` TEXT NOT NULL, `categoryId` TEXT NOT NULL, `brand` TEXT NOT NULL, `price` REAL NOT NULL, `description` TEXT, `status` TEXT NOT NULL, `lastSync` INTEGER NOT NULL, PRIMARY KEY(`id`))")
        connection.execSQL("CREATE TABLE IF NOT EXISTS `transactions` (`id` TEXT NOT NULL, `transactionId` TEXT NOT NULL, `productName` TEXT NOT NULL, `sellingPrice` REAL NOT NULL, `platformPrice` REAL NOT NULL, `customerNumber` TEXT NOT NULL, `status` TEXT NOT NULL, `createdAt` INTEGER NOT NULL, `brand` TEXT, `categoryId` TEXT, PRIMARY KEY(`id`))")
        connection.execSQL("CREATE TABLE IF NOT EXISTS `pending_sync_queue` (`id` TEXT NOT NULL, `type` TEXT NOT NULL, `payload` TEXT NOT NULL, `retryCount` INTEGER NOT NULL, `lastAttempt` INTEGER NOT NULL, `status` TEXT NOT NULL, PRIMARY KEY(`id`))")
        connection.execSQL("CREATE TABLE IF NOT EXISTS room_master_table (id INTEGER PRIMARY KEY,identity_hash TEXT)")
        connection.execSQL("INSERT OR REPLACE INTO room_master_table (id,identity_hash) VALUES(42, 'b247b74c89bfd7892a57e55dbba6715e')")
      }

      public override fun dropAllTables(connection: SQLiteConnection) {
        connection.execSQL("DROP TABLE IF EXISTS `categories`")
        connection.execSQL("DROP TABLE IF EXISTS `products`")
        connection.execSQL("DROP TABLE IF EXISTS `transactions`")
        connection.execSQL("DROP TABLE IF EXISTS `pending_sync_queue`")
      }

      public override fun onCreate(connection: SQLiteConnection) {
      }

      public override fun onOpen(connection: SQLiteConnection) {
        internalInitInvalidationTracker(connection)
      }

      public override fun onPreMigrate(connection: SQLiteConnection) {
        dropFtsSyncTriggers(connection)
      }

      public override fun onPostMigrate(connection: SQLiteConnection) {
      }

      public override fun onValidateSchema(connection: SQLiteConnection):
          RoomOpenDelegate.ValidationResult {
        val _columnsCategories: MutableMap<String, TableInfo.Column> = mutableMapOf()
        _columnsCategories.put("id", TableInfo.Column("id", "TEXT", true, 1, null,
            TableInfo.CREATED_FROM_ENTITY))
        _columnsCategories.put("name", TableInfo.Column("name", "TEXT", true, 0, null,
            TableInfo.CREATED_FROM_ENTITY))
        _columnsCategories.put("iconUrl", TableInfo.Column("iconUrl", "TEXT", false, 0, null,
            TableInfo.CREATED_FROM_ENTITY))
        val _foreignKeysCategories: MutableSet<TableInfo.ForeignKey> = mutableSetOf()
        val _indicesCategories: MutableSet<TableInfo.Index> = mutableSetOf()
        val _infoCategories: TableInfo = TableInfo("categories", _columnsCategories,
            _foreignKeysCategories, _indicesCategories)
        val _existingCategories: TableInfo = read(connection, "categories")
        if (!_infoCategories.equals(_existingCategories)) {
          return RoomOpenDelegate.ValidationResult(false, """
              |categories(com.yonotech.ppob.mobile.data.local.entities.CategoryEntity).
              | Expected:
              |""".trimMargin() + _infoCategories + """
              |
              | Found:
              |""".trimMargin() + _existingCategories)
        }
        val _columnsProducts: MutableMap<String, TableInfo.Column> = mutableMapOf()
        _columnsProducts.put("id", TableInfo.Column("id", "TEXT", true, 1, null,
            TableInfo.CREATED_FROM_ENTITY))
        _columnsProducts.put("name", TableInfo.Column("name", "TEXT", true, 0, null,
            TableInfo.CREATED_FROM_ENTITY))
        _columnsProducts.put("code", TableInfo.Column("code", "TEXT", true, 0, null,
            TableInfo.CREATED_FROM_ENTITY))
        _columnsProducts.put("categoryId", TableInfo.Column("categoryId", "TEXT", true, 0, null,
            TableInfo.CREATED_FROM_ENTITY))
        _columnsProducts.put("brand", TableInfo.Column("brand", "TEXT", true, 0, null,
            TableInfo.CREATED_FROM_ENTITY))
        _columnsProducts.put("price", TableInfo.Column("price", "REAL", true, 0, null,
            TableInfo.CREATED_FROM_ENTITY))
        _columnsProducts.put("description", TableInfo.Column("description", "TEXT", false, 0, null,
            TableInfo.CREATED_FROM_ENTITY))
        _columnsProducts.put("status", TableInfo.Column("status", "TEXT", true, 0, null,
            TableInfo.CREATED_FROM_ENTITY))
        _columnsProducts.put("lastSync", TableInfo.Column("lastSync", "INTEGER", true, 0, null,
            TableInfo.CREATED_FROM_ENTITY))
        val _foreignKeysProducts: MutableSet<TableInfo.ForeignKey> = mutableSetOf()
        val _indicesProducts: MutableSet<TableInfo.Index> = mutableSetOf()
        val _infoProducts: TableInfo = TableInfo("products", _columnsProducts, _foreignKeysProducts,
            _indicesProducts)
        val _existingProducts: TableInfo = read(connection, "products")
        if (!_infoProducts.equals(_existingProducts)) {
          return RoomOpenDelegate.ValidationResult(false, """
              |products(com.yonotech.ppob.mobile.data.local.entities.ProductEntity).
              | Expected:
              |""".trimMargin() + _infoProducts + """
              |
              | Found:
              |""".trimMargin() + _existingProducts)
        }
        val _columnsTransactions: MutableMap<String, TableInfo.Column> = mutableMapOf()
        _columnsTransactions.put("id", TableInfo.Column("id", "TEXT", true, 1, null,
            TableInfo.CREATED_FROM_ENTITY))
        _columnsTransactions.put("transactionId", TableInfo.Column("transactionId", "TEXT", true, 0,
            null, TableInfo.CREATED_FROM_ENTITY))
        _columnsTransactions.put("productName", TableInfo.Column("productName", "TEXT", true, 0,
            null, TableInfo.CREATED_FROM_ENTITY))
        _columnsTransactions.put("sellingPrice", TableInfo.Column("sellingPrice", "REAL", true, 0,
            null, TableInfo.CREATED_FROM_ENTITY))
        _columnsTransactions.put("platformPrice", TableInfo.Column("platformPrice", "REAL", true, 0,
            null, TableInfo.CREATED_FROM_ENTITY))
        _columnsTransactions.put("customerNumber", TableInfo.Column("customerNumber", "TEXT", true,
            0, null, TableInfo.CREATED_FROM_ENTITY))
        _columnsTransactions.put("status", TableInfo.Column("status", "TEXT", true, 0, null,
            TableInfo.CREATED_FROM_ENTITY))
        _columnsTransactions.put("createdAt", TableInfo.Column("createdAt", "INTEGER", true, 0,
            null, TableInfo.CREATED_FROM_ENTITY))
        _columnsTransactions.put("brand", TableInfo.Column("brand", "TEXT", false, 0, null,
            TableInfo.CREATED_FROM_ENTITY))
        _columnsTransactions.put("categoryId", TableInfo.Column("categoryId", "TEXT", false, 0,
            null, TableInfo.CREATED_FROM_ENTITY))
        val _foreignKeysTransactions: MutableSet<TableInfo.ForeignKey> = mutableSetOf()
        val _indicesTransactions: MutableSet<TableInfo.Index> = mutableSetOf()
        val _infoTransactions: TableInfo = TableInfo("transactions", _columnsTransactions,
            _foreignKeysTransactions, _indicesTransactions)
        val _existingTransactions: TableInfo = read(connection, "transactions")
        if (!_infoTransactions.equals(_existingTransactions)) {
          return RoomOpenDelegate.ValidationResult(false, """
              |transactions(com.yonotech.ppob.mobile.data.local.entities.TransactionEntity).
              | Expected:
              |""".trimMargin() + _infoTransactions + """
              |
              | Found:
              |""".trimMargin() + _existingTransactions)
        }
        val _columnsPendingSyncQueue: MutableMap<String, TableInfo.Column> = mutableMapOf()
        _columnsPendingSyncQueue.put("id", TableInfo.Column("id", "TEXT", true, 1, null,
            TableInfo.CREATED_FROM_ENTITY))
        _columnsPendingSyncQueue.put("type", TableInfo.Column("type", "TEXT", true, 0, null,
            TableInfo.CREATED_FROM_ENTITY))
        _columnsPendingSyncQueue.put("payload", TableInfo.Column("payload", "TEXT", true, 0, null,
            TableInfo.CREATED_FROM_ENTITY))
        _columnsPendingSyncQueue.put("retryCount", TableInfo.Column("retryCount", "INTEGER", true,
            0, null, TableInfo.CREATED_FROM_ENTITY))
        _columnsPendingSyncQueue.put("lastAttempt", TableInfo.Column("lastAttempt", "INTEGER", true,
            0, null, TableInfo.CREATED_FROM_ENTITY))
        _columnsPendingSyncQueue.put("status", TableInfo.Column("status", "TEXT", true, 0, null,
            TableInfo.CREATED_FROM_ENTITY))
        val _foreignKeysPendingSyncQueue: MutableSet<TableInfo.ForeignKey> = mutableSetOf()
        val _indicesPendingSyncQueue: MutableSet<TableInfo.Index> = mutableSetOf()
        val _infoPendingSyncQueue: TableInfo = TableInfo("pending_sync_queue",
            _columnsPendingSyncQueue, _foreignKeysPendingSyncQueue, _indicesPendingSyncQueue)
        val _existingPendingSyncQueue: TableInfo = read(connection, "pending_sync_queue")
        if (!_infoPendingSyncQueue.equals(_existingPendingSyncQueue)) {
          return RoomOpenDelegate.ValidationResult(false, """
              |pending_sync_queue(com.yonotech.ppob.mobile.data.local.entities.PendingSyncItem).
              | Expected:
              |""".trimMargin() + _infoPendingSyncQueue + """
              |
              | Found:
              |""".trimMargin() + _existingPendingSyncQueue)
        }
        return RoomOpenDelegate.ValidationResult(true, null)
      }
    }
    return _openDelegate
  }

  protected override fun createInvalidationTracker(): InvalidationTracker {
    val _shadowTablesMap: MutableMap<String, String> = mutableMapOf()
    val _viewTables: MutableMap<String, Set<String>> = mutableMapOf()
    return InvalidationTracker(this, _shadowTablesMap, _viewTables, "categories", "products",
        "transactions", "pending_sync_queue")
  }

  public override fun clearAllTables() {
    super.performClear(false, "categories", "products", "transactions", "pending_sync_queue")
  }

  protected override fun getRequiredTypeConverterClasses(): Map<KClass<*>, List<KClass<*>>> {
    val _typeConvertersMap: MutableMap<KClass<*>, List<KClass<*>>> = mutableMapOf()
    _typeConvertersMap.put(CategoryDao::class, CategoryDao_Impl.getRequiredConverters())
    _typeConvertersMap.put(ProductDao::class, ProductDao_Impl.getRequiredConverters())
    _typeConvertersMap.put(TransactionDao::class, TransactionDao_Impl.getRequiredConverters())
    _typeConvertersMap.put(PendingSyncDao::class, PendingSyncDao_Impl.getRequiredConverters())
    return _typeConvertersMap
  }

  public override fun getRequiredAutoMigrationSpecClasses(): Set<KClass<out AutoMigrationSpec>> {
    val _autoMigrationSpecsSet: MutableSet<KClass<out AutoMigrationSpec>> = mutableSetOf()
    return _autoMigrationSpecsSet
  }

  public override
      fun createAutoMigrations(autoMigrationSpecs: Map<KClass<out AutoMigrationSpec>, AutoMigrationSpec>):
      List<Migration> {
    val _autoMigrations: MutableList<Migration> = mutableListOf()
    return _autoMigrations
  }

  public override fun categoryDao(): CategoryDao = _categoryDao.value

  public override fun productDao(): ProductDao = _productDao.value

  public override fun transactionDao(): TransactionDao = _transactionDao.value

  public override fun pendingSyncDao(): PendingSyncDao = _pendingSyncDao.value
}
