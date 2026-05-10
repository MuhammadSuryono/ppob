package com.yonotech.ppob.services

import android.util.Log
import com.google.firebase.messaging.FirebaseMessagingService
import com.google.firebase.messaging.RemoteMessage
import com.yonotech.ppob.R

class FcmMessagingService : FirebaseMessagingService() {

    override fun onNewToken(token: String) {
        super.onNewToken(token)
        Log.d("FCM", "New token: $token")
        // Send token to your server
        sendTokenToServer(token)
    }

    override fun onMessageReceived(remoteMessage: RemoteMessage) {
        Log.d("FCM", "From: ${remoteMessage.from}")

        // Check if message contains a data payload
        remoteMessage.data.isNotEmpty().let {
            Log.d("FCM", "Message data payload: ${remoteMessage.data}")

            val type = remoteMessage.data["type"]
            when (type) {
                "transaction_success" -> handleTransactionSuccess(remoteMessage)
                "transaction_failed" -> handleTransactionFailed(remoteMessage)
                "transaction_pending" -> handleTransactionPending(remoteMessage)
                "low_balance_warning" -> handleLowBalanceWarning(remoteMessage)
                "staff_notification" -> handleStaffNotification(remoteMessage)
                else -> Log.d("FCM", "Unknown message type: $type")
            }
        }

        // Check if message contains a notification payload
        remoteMessage.notification?.let {
            Log.d("FCM", "Message notification body: ${it.body}")
        }
    }

    private fun handleTransactionSuccess(remoteMessage: RemoteMessage) {
        val amount = remoteMessage.data["amount"]?.toDoubleOrNull() ?: 0.0
        val product = remoteMessage.data["product_name"].orEmpty()
        val deepLink = remoteMessage.data["deep_link"].orEmpty()

        NotificationHelper.showTransactionNotification(
            this,
            title = "Transaksi Berhasil",
            body = "Rp ${amount.toInt()} — $product",
            deepLink
        )
    }

    private fun handleTransactionFailed(remoteMessage: RemoteMessage) {
        val amount = remoteMessage.data["amount"]?.toDoubleOrNull() ?: 0.0
        val product = remoteMessage.data["product_name"].orEmpty()
        val errorMessage = remoteMessage.data["error_message"].orEmpty()

        NotificationHelper.showTransactionNotification(
            this,
            title = "Transaksi Gagal",
            body = "Rp ${amount.toInt()} — $product — $errorMessage",
            deepLink = ""
        )
    }

    private fun handleTransactionPending(remoteMessage: RemoteMessage) {
        val amount = remoteMessage.data["amount"]?.toDoubleOrNull() ?: 0.0
        val product = remoteMessage.data["product_name"].orEmpty()

        NotificationHelper.showTransactionNotification(
            this,
            title = "Transaksi Diproses",
            body = "Rp ${amount.toInt()} — $product — Sedang diproses",
            deepLink = ""
        )
    }

    private fun handleLowBalanceWarning(remoteMessage: RemoteMessage) {
        NotificationHelper.showBalanceNotification(
            this,
            title = "Saldo Rendah",
            body = "Saldo Anda kurang dari Rp 10.000. Segera top up!"
        )
    }

    private fun handleStaffNotification(remoteMessage: RemoteMessage) {
        val message = remoteMessage.data["message"].orEmpty()

        NotificationHelper.showTransactionNotification(
            this,
            title = "Staff",
            body = message,
            deepLink = ""
        )
    }

    private fun sendTokenToServer(token: String) {
        // TODO: Send FCM token to backend server
        // Use auth token from DataStore and call UserService.updateFcmToken
        Log.d("FCM", "Sending token to server: $token")
    }
}