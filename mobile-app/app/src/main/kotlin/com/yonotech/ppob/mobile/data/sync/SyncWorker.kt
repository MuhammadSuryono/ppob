package com.yonotech.ppob.mobile.data.sync

import android.content.Context
import androidx.hilt.work.HiltWorker
import androidx.work.CoroutineWorker
import androidx.work.WorkerParameters
import com.yonotech.ppob.mobile.data.local.dao.PendingSyncDao
import com.yonotech.ppob.mobile.data.remote.TransactionService
import com.squareup.moshi.Moshi
import dagger.assisted.Assisted
import dagger.assisted.AssistedInject

@HiltWorker
class SyncWorker @AssistedInject constructor(
    @Assisted ctx: Context,
    @Assisted params: WorkerParameters,
    private val pendingSyncDao: PendingSyncDao,
    private val transactionService: TransactionService
) : CoroutineWorker(ctx, params) {

    private val moshi = Moshi.Builder().build()

    override suspend fun doWork(): Result {
        return try {
            val pendingItems = pendingSyncDao.getPendingItems()

            pendingItems.forEach { item ->
                try {
                    when (item.type) {
                        "TRANSACTION" -> {
                            val jsonAdapter = moshi.adapter(Map::class.java)
                            val payload = jsonAdapter.fromJson(item.payload)
                            pendingSyncDao.updateStatus(item.id, "COMPLETED", System.currentTimeMillis())
                        }
                    }
                } catch (e: Exception) {
                    val newRetryCount = item.retryCount + 1
                    if (newRetryCount >= 3) {
                        pendingSyncDao.updateStatus(item.id, "FAILED", System.currentTimeMillis())
                    } else {
                        pendingSyncDao.updateStatus(item.id, "PENDING", System.currentTimeMillis())
                    }
                }
            }

            pendingSyncDao.clearCompleted()
            Result.success()
        } catch (e: Exception) {
            Result.retry()
        }
    }
}