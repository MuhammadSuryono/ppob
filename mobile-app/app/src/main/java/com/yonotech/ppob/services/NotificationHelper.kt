package com.yonotech.ppob.services

import android.app.NotificationChannel
import android.app.NotificationManager
import android.content.Context
import android.os.Build
import androidx.core.app.NotificationCompat
import com.yonotech.ppob.R

object NotificationHelper {

    const val CHANNEL_TRANSACTIONS = "transactions"
    const val CHANNEL_ALERTS = "alerts"
    const val CHANNEL_STAFF = "staff"
    const val CHANNEL_MARKETING = "marketing"

    fun createNotificationChannels(context: Context) {
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            val channels = listOf(
                NotificationChannel(
                    CHANNEL_TRANSACTIONS,
                    "Transaksi",
                    NotificationManager.IMPORTANCE_HIGH
                ).apply {
                    description = "Notifikasi transaksi berhasil, gagal, dan saldo rendah"
                },
                NotificationChannel(
                    CHANNEL_ALERTS,
                    "Keamanan",
                    NotificationManager.IMPORTANCE_HIGH
                ).apply {
                    description = "Notifikasi keamanan akun"
                },
                NotificationChannel(
                    CHANNEL_STAFF,
                    "Staff",
                    NotificationManager.IMPORTANCE_DEFAULT
                ).apply {
                    description = "Notifikasi staff ditambahkan, top-up diterima"
                },
                NotificationChannel(
                    CHANNEL_MARKETING,
                    "Promo",
                    NotificationManager.IMPORTANCE_LOW
                ).apply {
                    description = "Notifikasi promosi (opsional)"
                }
            )

            val notificationManager =
                context.getSystemService(Context.NOTIFICATION_SERVICE) as NotificationManager
            channels.forEach { notificationManager.createNotificationChannel(it) }
        }
    }

    fun showTransactionNotification(context: Context, title: String, body: String, deepLink: String) {
        val notificationManager =
            context.getSystemService(Context.NOTIFICATION_SERVICE) as NotificationManager

        val notification = NotificationCompat.Builder(context, CHANNEL_TRANSACTIONS)
            .setSmallIcon(R.drawable.ic_notification)
            .setContentTitle(title)
            .setContentText(body)
            .setPriority(NotificationCompat.PRIORITY_HIGH)
            .setAutoCancel(true)
            .build()

        notificationManager.notify(System.currentTimeMillis().toInt(), notification)
    }

    fun showBalanceNotification(context: Context, title: String, body: String) {
        val notificationManager =
            context.getSystemService(Context.NOTIFICATION_SERVICE) as NotificationManager

        val notification = NotificationCompat.Builder(context, CHANNEL_ALERTS)
            .setSmallIcon(R.drawable.ic_notification)
            .setContentTitle(title)
            .setContentText(body)
            .setPriority(NotificationCompat.PRIORITY_HIGH)
            .setAutoCancel(true)
            .build()

        notificationManager.notify(System.currentTimeMillis().toInt(), notification)
    }
}