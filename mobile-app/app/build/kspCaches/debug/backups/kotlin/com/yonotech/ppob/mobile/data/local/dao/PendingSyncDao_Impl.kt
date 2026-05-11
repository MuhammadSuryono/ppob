package com.yonotech.ppob.mobile.`data`.local.dao

import androidx.room.EntityInsertAdapter
import androidx.room.RoomDatabase
import androidx.room.util.getColumnIndexOrThrow
import androidx.room.util.performSuspending
import androidx.sqlite.SQLiteStatement
import com.yonotech.ppob.mobile.`data`.local.entities.PendingSyncItem
import javax.`annotation`.processing.Generated
import kotlin.Int
import kotlin.Long
import kotlin.String
import kotlin.Suppress
import kotlin.Unit
import kotlin.collections.List
import kotlin.collections.MutableList
import kotlin.collections.mutableListOf
import kotlin.reflect.KClass

@Generated(value = ["androidx.room.RoomProcessor"])
@Suppress(names = ["UNCHECKED_CAST", "DEPRECATION", "REDUNDANT_PROJECTION", "REMOVAL"])
public class PendingSyncDao_Impl(
  __db: RoomDatabase,
) : PendingSyncDao {
  private val __db: RoomDatabase

  private val __insertAdapterOfPendingSyncItem: EntityInsertAdapter<PendingSyncItem>
  init {
    this.__db = __db
    this.__insertAdapterOfPendingSyncItem = object : EntityInsertAdapter<PendingSyncItem>() {
      protected override fun createQuery(): String =
          "INSERT OR REPLACE INTO `pending_sync_queue` (`id`,`type`,`payload`,`retryCount`,`lastAttempt`,`status`) VALUES (?,?,?,?,?,?)"

      protected override fun bind(statement: SQLiteStatement, entity: PendingSyncItem) {
        statement.bindText(1, entity.id)
        statement.bindText(2, entity.type)
        statement.bindText(3, entity.payload)
        statement.bindLong(4, entity.retryCount.toLong())
        statement.bindLong(5, entity.lastAttempt)
        statement.bindText(6, entity.status)
      }
    }
  }

  public override suspend fun insert(item: PendingSyncItem): Unit = performSuspending(__db, false,
      true) { _connection ->
    __insertAdapterOfPendingSyncItem.insert(_connection, item)
  }

  public override suspend fun getPendingItems(): List<PendingSyncItem> {
    val _sql: String =
        "SELECT * FROM pending_sync_queue WHERE status = 'PENDING' ORDER BY lastAttempt ASC"
    return performSuspending(__db, true, false) { _connection ->
      val _stmt: SQLiteStatement = _connection.prepare(_sql)
      try {
        val _columnIndexOfId: Int = getColumnIndexOrThrow(_stmt, "id")
        val _columnIndexOfType: Int = getColumnIndexOrThrow(_stmt, "type")
        val _columnIndexOfPayload: Int = getColumnIndexOrThrow(_stmt, "payload")
        val _columnIndexOfRetryCount: Int = getColumnIndexOrThrow(_stmt, "retryCount")
        val _columnIndexOfLastAttempt: Int = getColumnIndexOrThrow(_stmt, "lastAttempt")
        val _columnIndexOfStatus: Int = getColumnIndexOrThrow(_stmt, "status")
        val _result: MutableList<PendingSyncItem> = mutableListOf()
        while (_stmt.step()) {
          val _item: PendingSyncItem
          val _tmpId: String
          _tmpId = _stmt.getText(_columnIndexOfId)
          val _tmpType: String
          _tmpType = _stmt.getText(_columnIndexOfType)
          val _tmpPayload: String
          _tmpPayload = _stmt.getText(_columnIndexOfPayload)
          val _tmpRetryCount: Int
          _tmpRetryCount = _stmt.getLong(_columnIndexOfRetryCount).toInt()
          val _tmpLastAttempt: Long
          _tmpLastAttempt = _stmt.getLong(_columnIndexOfLastAttempt)
          val _tmpStatus: String
          _tmpStatus = _stmt.getText(_columnIndexOfStatus)
          _item =
              PendingSyncItem(_tmpId,_tmpType,_tmpPayload,_tmpRetryCount,_tmpLastAttempt,_tmpStatus)
          _result.add(_item)
        }
        _result
      } finally {
        _stmt.close()
      }
    }
  }

  public override suspend fun updateStatus(
    id: String,
    status: String,
    timestamp: Long,
  ) {
    val _sql: String =
        "UPDATE pending_sync_queue SET status = ?, lastAttempt = ?, retryCount = retryCount + 1 WHERE id = ?"
    return performSuspending(__db, false, true) { _connection ->
      val _stmt: SQLiteStatement = _connection.prepare(_sql)
      try {
        var _argIndex: Int = 1
        _stmt.bindText(_argIndex, status)
        _argIndex = 2
        _stmt.bindLong(_argIndex, timestamp)
        _argIndex = 3
        _stmt.bindText(_argIndex, id)
        _stmt.step()
      } finally {
        _stmt.close()
      }
    }
  }

  public override suspend fun clearCompleted() {
    val _sql: String = "DELETE FROM pending_sync_queue WHERE status = 'COMPLETED'"
    return performSuspending(__db, false, true) { _connection ->
      val _stmt: SQLiteStatement = _connection.prepare(_sql)
      try {
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
