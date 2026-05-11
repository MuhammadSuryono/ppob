package com.yonotech.ppob.mobile

import android.app.Application
import android.app.NotificationChannel
import android.app.NotificationManager
import android.os.Build
import dagger.hilt.android.HiltAndroidApp

@HiltAndroidApp
class PpoApplication : Application() {
    override fun onCreate() {
        super.onCreate()
        createNotificationChannels()
    }

    private fun createNotificationChannels() {
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
            val channels = listOf(
                NotificationChannel(
                    "transactions",
                    "Transaksi",
                    NotificationManager.IMPORTANCE_HIGH
                ),
                NotificationChannel(
                    "alerts",
                    "Keamanan",
                    NotificationManager.IMPORTANCE_HIGH
                ),
                NotificationChannel(
                    "staff",
                    "Staff",
                    NotificationManager.IMPORTANCE_DEFAULT
                ),
                NotificationChannel(
                    "marketing",
                    "Promo",
                    NotificationManager.IMPORTANCE_LOW
                )
            )
            val manager = getSystemService(NotificationManager::class.java)
            channels.forEach { manager.createNotificationChannel(it) }
        }
    }
}
