package com.yonotech.ppob

import android.app.Application
import android.util.Log
import com.google.firebase.FirebaseApp
import com.google.firebase.messaging.FirebaseMessaging
import com.yonotech.ppob.data.local.datastore.AuthPreferencesDataStore
import dagger.hilt.android.HiltAndroidApp
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.SupervisorJob
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.launch
import javax.inject.Inject

@HiltAndroidApp
class PPOBApp : Application() {

    private val appScope = CoroutineScope(SupervisorJob() + Dispatchers.IO)

    @Inject
    lateinit var authPreferencesDataStore: AuthPreferencesDataStore

    @Inject
    lateinit var syncWorkerHelper: com.yonotech.ppob.data.local.datastore.SyncWorkerHelper

    override fun onCreate() {
        super.onCreate()

        // Initialize Firebase
        FirebaseApp.initializeApp(this)

        // Request FCM token and store it
        requestFcmToken()

        // Register notification channels
        com.yonotech.ppob.services.NotificationHelper.createNotificationChannels(this)

        // Enqueue periodic sync work on app start
        appScope.launch {
            // Check if user is logged in before scheduling sync
            try {
                val prefs = authPreferencesDataStore.authPreferences.first()
                if (prefs.accessToken.isNotEmpty()) {
                    com.yonotech.ppob.data.local.datastore.SyncWorkerHelper.enqueuePeriodicSync(this@PPOBApp)
                }
            } catch (e: Exception) {
                Log.d("PPOBApp", "No auth prefs found, skipping sync scheduling")
            }
        }
    }

    private fun requestFcmToken() {
        FirebaseMessaging.getInstance().token.addOnCompleteListener { task ->
            if (task.isSuccessful) {
                val token = task.result
                Log.d("PPOBApp", "FCM Token: $token")
                // Store token locally and send to backend
                appScope.launch {
                    try {
                        authPreferencesDataStore.saveFcmToken(token)
                    } catch (e: Exception) {
                        Log.e("PPOBApp", "Failed to save FCM token", e)
                    }
                }
            } else {
                Log.w("PPOBApp", "Fetching FCM registration token failed", task.exception)
            }
        }
    }
}