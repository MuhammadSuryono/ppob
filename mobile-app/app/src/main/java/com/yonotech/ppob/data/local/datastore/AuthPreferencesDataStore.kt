package com.yonotech.ppob.data.local.datastore

import android.content.Context
import androidx.datastore.preferences.core.Preferences
import androidx.datastore.preferences.core.booleanPreferencesKey
import androidx.datastore.preferences.core.edit
import androidx.datastore.preferences.core.longPreferencesKey
import androidx.datastore.preferences.core.stringPreferencesKey
import androidx.datastore.preferences.preferencesDataStore
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.map

// DataStore for auth preferences
val Context.authDataStore by preferencesDataStore(name = "auth_preferences")
val Context.appDataStore by preferencesDataStore(name = "app_preferences")

object AuthPreferencesKeys {
    val ACCESS_TOKEN = stringPreferencesKey("access_token")
    val REFRESH_TOKEN = stringPreferencesKey("refresh_token")
    val USER_ID = stringPreferencesKey("user_id")
    val PHONE_NUMBER = stringPreferencesKey("phone_number")
    val USER_NAME = stringPreferencesKey("user_name")
    val ACTIVE_ROLE = stringPreferencesKey("active_role")
    val WALLET_ID = stringPreferencesKey("wallet_id")
    val FCM_TOKEN = stringPreferencesKey("fcm_token")
    val DEVICE_ID = stringPreferencesKey("device_id")
    val IS_FIRST_LAUNCH = booleanPreferencesKey("is_first_launch")
    val PIN_SET = booleanPreferencesKey("pin_set")
    val BIOMETRIC_ENABLED = booleanPreferencesKey("biometric_enabled")
    val TRUSTED_DEVICE = booleanPreferencesKey("trusted_device")
}

object AppPreferenceKeys {
    val THEME = stringPreferencesKey("theme")
    val LANGUAGE = stringPreferencesKey("language")
    val LAST_PRODUCT_SYNC = stringPreferencesKey("last_product_sync")
    val LAST_SYNC_TIMESTAMP = longPreferencesKey("last_sync_timestamp")
}

data class AuthPreferences(
    val accessToken: String = "",
    val refreshToken: String = "",
    val userId: String = "",
    val phoneNumber: String = "",
    val userName: String = "",
    val activeRole: String = "",
    val walletId: String = "",
    val fcmToken: String = "",
    val deviceId: String = "",
    val isFirstLaunch: Boolean = true,
    val pinSet: Boolean = false,
    val biometricEnabled: Boolean = false,
    val trustedDevice: Boolean = false
)

class AuthPreferencesDataStore(private val context: Context) {

    val authPreferences: Flow<AuthPreferences> = context.authDataStore.data.map { preferences ->
        AuthPreferences(
            accessToken = preferences[AuthPreferencesKeys.ACCESS_TOKEN] ?: "",
            refreshToken = preferences[AuthPreferencesKeys.REFRESH_TOKEN] ?: "",
            userId = preferences[AuthPreferencesKeys.USER_ID] ?: "",
            phoneNumber = preferences[AuthPreferencesKeys.PHONE_NUMBER] ?: "",
            userName = preferences[AuthPreferencesKeys.USER_NAME] ?: "",
            activeRole = preferences[AuthPreferencesKeys.ACTIVE_ROLE] ?: "",
            walletId = preferences[AuthPreferencesKeys.WALLET_ID] ?: "",
            fcmToken = preferences[AuthPreferencesKeys.FCM_TOKEN] ?: "",
            deviceId = preferences[AuthPreferencesKeys.DEVICE_ID] ?: "",
            isFirstLaunch = preferences[AuthPreferencesKeys.IS_FIRST_LAUNCH] ?: true,
            pinSet = preferences[AuthPreferencesKeys.PIN_SET] ?: false,
            biometricEnabled = preferences[AuthPreferencesKeys.BIOMETRIC_ENABLED] ?: false,
            trustedDevice = preferences[AuthPreferencesKeys.TRUSTED_DEVICE] ?: false
        )
    }

    suspend fun saveAccessToken(token: String) {
        context.authDataStore.edit { it[AuthPreferencesKeys.ACCESS_TOKEN] = token }
    }

    suspend fun saveRefreshToken(token: String) {
        context.authDataStore.edit { it[AuthPreferencesKeys.REFRESH_TOKEN] = token }
    }

    suspend fun saveUserCredentials(
        accessToken: String,
        refreshToken: String,
        userId: String,
        phoneNumber: String,
        userName: String,
        activeRole: String,
        walletId: String
    ) {
        context.authDataStore.edit { prefs ->
            prefs[AuthPreferencesKeys.ACCESS_TOKEN] = accessToken
            prefs[AuthPreferencesKeys.REFRESH_TOKEN] = refreshToken
            prefs[AuthPreferencesKeys.USER_ID] = userId
            prefs[AuthPreferencesKeys.PHONE_NUMBER] = phoneNumber
            prefs[AuthPreferencesKeys.USER_NAME] = userName
            prefs[AuthPreferencesKeys.ACTIVE_ROLE] = activeRole
            prefs[AuthPreferencesKeys.WALLET_ID] = walletId
            prefs[AuthPreferencesKeys.IS_FIRST_LAUNCH] = false
        }
    }

    suspend fun savePin(pin: String) {
        // NOTE: In production, store a hash, not plaintext
        context.authDataStore.edit { it[AuthPreferencesKeys.ACCESS_TOKEN] = pin }
    }

    suspend fun saveFcmToken(token: String) {
        context.authDataStore.edit { it[AuthPreferencesKeys.FCM_TOKEN] = token }
    }

    suspend fun saveDeviceId(deviceId: String) {
        context.authDataStore.edit { it[AuthPreferencesKeys.DEVICE_ID] = deviceId }
    }

    suspend fun setTrustedDevice(isTrusted: Boolean) {
        context.authDataStore.edit { it[AuthPreferencesKeys.TRUSTED_DEVICE] = isTrusted }
    }

    suspend fun enableBiometric(enabled: Boolean) {
        context.authDataStore.edit { it[AuthPreferencesKeys.BIOMETRIC_ENABLED] = enabled }
    }

    suspend fun clearAuthData() {
        context.authDataStore.edit { prefs ->
            prefs.remove(AuthPreferencesKeys.ACCESS_TOKEN)
            prefs.remove(AuthPreferencesKeys.REFRESH_TOKEN)
            prefs.remove(AuthPreferencesKeys.USER_ID)
            prefs.remove(AuthPreferencesKeys.PHONE_NUMBER)
            prefs.remove(AuthPreferencesKeys.USER_NAME)
            prefs.remove(AuthPreferencesKeys.ACTIVE_ROLE)
            prefs.remove(AuthPreferencesKeys.WALLET_ID)
            prefs.remove(AuthPreferencesKeys.FCM_TOKEN)
        }
    }

    suspend fun clear() {
        context.authDataStore.edit { it.clear() }
    }
}

data class AppPreferences(
    val theme: String = "light",
    val language: String = "id",
    val lastProductSync: String = "",
    val lastSyncTimestamp: Long = 0L
)

class AppPreferencesDataStore(private val context: Context) {

    val appPreferences: Flow<AppPreferences> = context.appDataStore.data.map { preferences ->
        AppPreferences(
            theme = preferences[AppPreferenceKeys.THEME] ?: "light",
            language = preferences[AppPreferenceKeys.LANGUAGE] ?: "id",
            lastProductSync = preferences[AppPreferenceKeys.LAST_PRODUCT_SYNC] ?: "",
            lastSyncTimestamp = preferences[AppPreferenceKeys.LAST_SYNC_TIMESTAMP] ?: 0L
        )
    }

    suspend fun setTheme(theme: String) {
        context.appDataStore.edit { it[AppPreferenceKeys.THEME] = theme }
    }

    suspend fun setLanguage(language: String) {
        context.appDataStore.edit { it[AppPreferenceKeys.LANGUAGE] = language }
    }

    suspend fun setLastProductSync(timestamp: String) {
        context.appDataStore.edit { it[AppPreferenceKeys.LAST_PRODUCT_SYNC] = timestamp }
    }

    suspend fun setLastSyncTimestamp(timestamp: Long) {
        context.appDataStore.edit { it[AppPreferenceKeys.LAST_SYNC_TIMESTAMP] = timestamp }
    }
}

// Sync Worker Helper
object SyncWorkerHelper {
    fun enqueuePeriodicSync(context: Context) {
        // TODO: Implement WorkManager periodic sync request
        // val syncRequest = PeriodicWorkRequestBuilder<SyncWorker>(
        //     15, TimeUnit.MINUTES,
        //     5, TimeUnit.MINUTES
        // ).build()
        // WorkManager.getInstance(context).enqueueUniquePeriodicWork(
        //     "sync_pending_transactions",
        //     ExistingPeriodicWorkPolicy.KEEP,
        //     syncRequest
        // )
    }
}
