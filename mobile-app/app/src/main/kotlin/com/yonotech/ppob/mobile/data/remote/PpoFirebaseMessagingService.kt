package com.yonotech.ppob.mobile.data.remote

import android.app.NotificationChannel
import android.app.NotificationManager
import android.content.Context
import android.os.Build
import androidx.core.app.NotificationCompat
import com.google.firebase.messaging.FirebaseMessagingService
import com.google.firebase.messaging.RemoteMessage
import com.yonotech.ppob.mobile.R
import dagger.hilt.android.AndroidEntryPoint
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.launch
import javax.inject.Inject

@AndroidEntryPoint
class PpoFirebaseMessagingService : FirebaseMessagingService() {

    override fun onCreate() {
        super.onCreate()
        createNotificationChannels()
    }

    override fun onNewToken(token: String) {
        super.onNewToken(token)
        // Send token to backend
    }

    override fun onMessageReceived(message: RemoteMessage) {
        super.onMessageReceived(message)

        val title = message.notification?.title ?: "PPOB Notification"
        val body = message.notification?.body ?: ""
        val txId = message.data["transaction_id"]

        showNotification(title, body, txId)
    }

    private fun createNotificationChannels() {
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            val successChannel = NotificationChannel(
                "success_channel",
                "Transaksi Berhasil",
                NotificationManager.IMPORTANCE_HIGH
            )
            val failureChannel = NotificationChannel(
                "failure_channel",
                "Transaksi Gagal",
                NotificationManager.IMPORTANCE_HIGH
            )
            val generalChannel = NotificationChannel(
                "general_channel",
                "Notifikasi Umum",
                NotificationManager.IMPORTANCE_DEFAULT
            )

            val manager = getSystemService(NotificationManager::class.java)
            manager.createNotificationChannel(successChannel)
            manager.createNotificationChannel(failureChannel)
            manager.createNotificationChannel(generalChannel)
        }
    }

    private fun showNotification(title: String, body: String, txId: String?) {
        val channelId = when {
            title.contains("Berhasil", ignoreCase = true) -> "success_channel"
            title.contains("Gagal", ignoreCase = true) -> "failure_channel"
            else -> "general_channel"
        }

        val notification = NotificationCompat.Builder(this, channelId)
            .setContentTitle(title)
            .setContentText(body)
            .setSmallIcon(R.drawable.ic_notification)
            .setAutoCancel(true)
            .build()

        val manager = getSystemService(Context.NOTIFICATION_SERVICE) as NotificationManager
        manager.notify(System.currentTimeMillis().toInt(), notification)
    }
}