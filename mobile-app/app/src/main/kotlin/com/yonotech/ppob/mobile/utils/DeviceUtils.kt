package com.yonotech.ppob.mobile.utils

import android.annotation.SuppressLint
import android.content.Context
import android.os.Build
import android.provider.Settings
import java.util.*

object DeviceUtils {
    @SuppressLint("HardwareIds")
    fun getDeviceId(context: Context): String {
        return Settings.Secure.getString(context.contentResolver, Settings.Secure.ANDROID_ID) ?: UUID.randomUUID().toString()
    }

    fun getDeviceName(): String {
        val manufacturer = Build.MANUFACTURER
        val model = Build.MODEL
        return if (model.startsWith(manufacturer)) {
            model.replaceFirstChar { if (it.isLowerCase()) it.titlecase(Locale.getDefault()) else it.toString() }
        } else {
            "${manufacturer.replaceFirstChar { if (it.isLowerCase()) it.titlecase(Locale.getDefault()) else it.toString() }} $model"
        }
    }

    fun getOsVersion(): String {
        return "Android ${Build.VERSION.RELEASE} (API ${Build.VERSION.SDK_INT})"
    }

    fun getFingerprint(): String {
        return Build.FINGERPRINT
    }
}
