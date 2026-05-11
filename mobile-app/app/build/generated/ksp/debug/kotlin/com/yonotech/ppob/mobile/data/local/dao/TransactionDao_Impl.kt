package com.yonotech.ppob.mobile.`data`.local.dao

import androidx.room.EntityInsertAdapter
import androidx.room.RoomDatabase
import androidx.room.coroutines.createFlow
import androidx.room.util.getColumnIndexOrThrow
import androidx.room.util.performSuspending
import androidx.sqlite.SQLiteStatement
import com.yonotech.ppob.mobile.`data`.local.entities.TransactionEntity
import javax.`annotation`.processing.Generated
import kotlin.Double
import kotlin.Int
import kotlin.Long
import kotlin.String
import kotlin.Suppress
import kotlin.Unit
import kotlin.collections.List
import kotlin.collections.MutableList
import kotlin.collections.mutableListOf
import kotlin.reflect.KClass
import kotlinx.coroutines.flow.Flow

@Generated(value = ["androidx.room.RoomProcessor"])
@Suppress(names = ["UNCHECKED_CAST", "DEPRECATION", "REDUNDANT_PROJECTION", "REMOVAL"])
public class TransactionDao_Impl(
  __db: RoomDatabase,
) : TransactionDao {
  private val __db: RoomDatabase

  private val __insertAdapterOfTransactionEntity: EntityInsertAdapter<TransactionEntity>
  init {
    this.__db = __db
    this.__insertAdapterOfTransactionEntity = object : EntityInsertAdapter<TransactionEntity>() {
      protected override fun createQuery(): String =
          "INSERT OR REPLACE INTO `transactions` (`id`,`transactionId`,`productName`,`sellingPrice`,`platformPrice`,`customerNumber`,`status`,`createdAt`,`brand`,`categoryId`) VALUES (?,?,?,?,?,?,?,?,?,?)"

      protected override fun bind(statement: SQLiteStatement, entity: TransactionEntity) {
        statement.bindText(1, entity.id)
        statement.bindText(2, entity.transactionId)
        statement.bindText(3, entity.productName)
        statement.bindDouble(4, entity.sellingPrice)
        statement.bindDouble(5, entity.platformPrice)
        statement.bindText(6, entity.customerNumber)
        statement.bindText(7, entity.status)
        statement.bindLong(8, entity.createdAt)
        val _tmpBrand: String? = entity.brand
        if (_tmpBrand == null) {
          statement.bindNull(9)
        } else {
          statement.bindText(9, _tmpBrand)
        }
        val _tmpCategoryId: String? = entity.categoryId
        if (_tmpCategoryId == null) {
          statement.bindNull(10)
        } else {
          statement.bindText(10, _tmpCategoryId)
        }
      }
    }
  }

  public override suspend fun insert(transaction: TransactionEntity): Unit = performSuspending(__db,
      false, true) { _connection ->
    __insertAdapterOfTransactionEntity.insert(_connection, transaction)
  }

  public override fun getAll(): Flow<List<TransactionEntity>> {
    val _sql: String = "SELECT * FROM transactions ORDER BY createdAt DESC"
    return createFlow(__db, false, arrayOf("transactions")) { _connection ->
      val _stmt: SQLiteStatement = _connection.prepare(_sql)
      try {
        val _columnIndexOfId: Int = getColumnIndexOrThrow(_stmt, "id")
        val _columnIndexOfTransactionId: Int = getColumnIndexOrThrow(_stmt, "transactionId")
        val _columnIndexOfProductName: Int = getColumnIndexOrThrow(_stmt, "productName")
        val _columnIndexOfSellingPrice: Int = getColumnIndexOrThrow(_stmt, "sellingPrice")
        val _columnIndexOfPlatformPrice: Int = getColumnIndexOrThrow(_stmt, "platformPrice")
        val _columnIndexOfCustomerNumber: Int = getColumnIndexOrThrow(_stmt, "customerNumber")
        val _columnIndexOfStatus: Int = getColumnIndexOrThrow(_stmt, "status")
        val _columnIndexOfCreatedAt: Int = getColumnIndexOrThrow(_stmt, "createdAt")
        val _columnIndexOfBrand: Int = getColumnIndexOrThrow(_stmt, "brand")
        val _columnIndexOfCategoryId: Int = getColumnIndexOrThrow(_stmt, "categoryId")
        val _result: MutableList<TransactionEntity> = mutableListOf()
        while (_stmt.step()) {
          val _item: TransactionEntity
          val _tmpId: String
          _tmpId = _stmt.getText(_columnIndexOfId)
          val _tmpTransactionId: String
          _tmpTransactionId = _stmt.getText(_columnIndexOfTransactionId)
          val _tmpProductName: String
          _tmpProductName = _stmt.getText(_columnIndexOfProductName)
          val _tmpSellingPrice: Double
          _tmpSellingPrice = _stmt.getDouble(_columnIndexOfSellingPrice)
          val _tmpPlatformPrice: Double
          _tmpPlatformPrice = _stmt.getDouble(_columnIndexOfPlatformPrice)
          val _tmpCustomerNumber: String
          _tmpCustomerNumber = _stmt.getText(_columnIndexOfCustomerNumber)
          val _tmpStatus: String
          _tmpStatus = _stmt.getText(_columnIndexOfStatus)
          val _tmpCreatedAt: Long
          _tmpCreatedAt = _stmt.getLong(_columnIndexOfCreatedAt)
          val _tmpBrand: String?
          if (_stmt.isNull(_columnIndexOfBrand)) {
            _tmpBrand = null
          } else {
            _tmpBrand = _stmt.getText(_columnIndexOfBrand)
          }
          val _tmpCategoryId: String?
          if (_stmt.isNull(_columnIndexOfCategoryId)) {
            _tmpCategoryId = null
          } else {
            _tmpCategoryId = _stmt.getText(_columnIndexOfCategoryId)
          }
          _item =
              TransactionEntity(_tmpId,_tmpTransactionId,_tmpProductName,_tmpSellingPrice,_tmpPlatformPrice,_tmpCustomerNumber,_tmpStatus,_tmpCreatedAt,_tmpBrand,_tmpCategoryId)
          _result.add(_item)
        }
        _result
      } finally {
        _stmt.close()
      }
    }
  }

  public override suspend fun getById(transactionId: String): TransactionEntity? {
    val _sql: String = "SELECT * FROM transactions WHERE id = ?"
    return performSuspending(__db, true, false) { _connection ->
      val _stmt: SQLiteStatement = _connection.prepare(_sql)
      try {
        var _argIndex: Int = 1
        _stmt.bindText(_argIndex, transactionId)
        val _columnIndexOfId: Int = getColumnIndexOrThrow(_stmt, "id")
        val _columnIndexOfTransactionId: Int = getColumnIndexOrThrow(_stmt, "transactionId")
        val _columnIndexOfProductName: Int = getColumnIndexOrThrow(_stmt, "productName")
        val _columnIndexOfSellingPrice: Int = getColumnIndexOrThrow(_stmt, "sellingPrice")
        val _columnIndexOfPlatformPrice: Int = getColumnIndexOrThrow(_stmt, "platformPrice")
        val _columnIndexOfCustomerNumber: Int = getColumnIndexOrThrow(_stmt, "customerNumber")
        val _columnIndexOfStatus: Int = getColumnIndexOrThrow(_stmt, "status")
        val _columnIndexOfCreatedAt: Int = getColumnIndexOrThrow(_stmt, "createdAt")
        val _columnIndexOfBrand: Int = getColumnIndexOrThrow(_stmt, "brand")
        val _columnIndexOfCategoryId: Int = getColumnIndexOrThrow(_stmt, "categoryId")
        val _result: TransactionEntity?
        if (_stmt.step()) {
          val _tmpId: String
          _tmpId = _stmt.getText(_columnIndexOfId)
          val _tmpTransactionId: String
          _tmpTransactionId = _stmt.getText(_columnIndexOfTransactionId)
          val _tmpProductName: String
          _tmpProductName = _stmt.getText(_columnIndexOfProductName)
          val _tmpSellingPrice: Double
          _tmpSellingPrice = _stmt.getDouble(_columnIndexOfSellingPrice)
          val _tmpPlatformPrice: Double
          _tmpPlatformPrice = _stmt.getDouble(_columnIndexOfPlatformPrice)
          val _tmpCustomerNumber: String
          _tmpCustomerNumber = _stmt.getText(_columnIndexOfCustomerNumber)
          val _tmpStatus: String
          _tmpStatus = _stmt.getText(_columnIndexOfStatus)
          val _tmpCreatedAt: Long
          _tmpCreatedAt = _stmt.getLong(_columnIndexOfCreatedAt)
          val _tmpBrand: String?
          if (_stmt.isNull(_columnIndexOfBrand)) {
            _tmpBrand = null
          } else {
            _tmpBrand = _stmt.getText(_columnIndexOfBrand)
          }
          val _tmpCategoryId: String?
          if (_stmt.isNull(_columnIndexOfCategoryId)) {
            _tmpCategoryId = null
          } else {
            _tmpCategoryId = _stmt.getText(_columnIndexOfCategoryId)
          }
          _result =
              TransactionEntity(_tmpId,_tmpTransactionId,_tmpProductName,_tmpSellingPrice,_tmpPlatformPrice,_tmpCustomerNumber,_tmpStatus,_tmpCreatedAt,_tmpBrand,_tmpCategoryId)
        } else {
          _result = null
        }
        _result
      } finally {
        _stmt.close()
      }
    }
  }

  public override suspend fun deleteOlderThan(cutoffTime: Long) {
    val _sql: String = "DELETE FROM transactions WHERE createdAt < ?"
    return performSuspending(__db, false, true) { _connection ->
      val _stmt: SQLiteStatement = _connection.prepare(_sql)
      try {
        var _argIndex: Int = 1
        _stmt.bindLong(_argIndex, cutoffTime)
        _stmt.step()
      } finally {
        _stmt.close()
      }
    }
  }

  public companion object {
    public fun getRequiredConverters(): List<KClass<*>> = emptyList()
  }
}
