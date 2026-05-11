package com.yonotech.ppob.mobile.ui.navigation

sealed class Screen(val route: String) {
    // Auth screens
    object PhoneInput : Screen("phone_input")
    object Otp : Screen("otp/{requestId}/{phone}/{type}") {
        fun createRoute(requestId: String, phone: String, type: String) = "otp/$requestId/$phone/$type"
    }
    object PasswordLogin : Screen("password_login/{phone}/{requestId}") {
        fun createRoute(phone: String, requestId: String) = "password_login/$phone/$requestId"
    }
    object SetPasswordPin : Screen("set_password_pin/{phone}/{requestId}") {
        fun createRoute(phone: String, requestId: String) = "set_password_pin/$phone/$requestId"
    }
    object PinLogin : Screen("pin_login/{phone}") {
        fun createRoute(phone: String) = "pin_login/$phone"
    }

    // Main screens (bottom nav)
    object Home : Screen("home")
    object Transactions : Screen("transactions")
    object Wallet : Screen("wallet")
    object Staff : Screen("staff")
    object Profile : Screen("profile")

    // Drawer menu
    object Settings : Screen("settings")
    object ChangePin : Screen("change_pin")
    object TrustedDevices : Screen("trusted_devices")
    object Help : Screen("help")

    // Other screens
    object ProductList : Screen("products/{categoryId}") {
        fun createRoute(categoryId: String) = "products/$categoryId"
    }
    object TransactionInit : Screen("transaction/init/{productId}") {
        fun createRoute(productId: String) = "transaction/init/$productId"
    }
    object TransactionConfirm : Screen("transaction/confirm")
    object TransactionResult : Screen("transaction/result/{txId}") {
        fun createRoute(txId: String) = "transaction/result/$txId"
    }
    object StaffAddEdit : Screen("staff/add_edit/{staffId}") {
        fun createRoute(staffId: String = "") = "staff/add_edit/$staffId"
    }
    object StaffTopUp : Screen("staff/topup/{staffId}") {
        fun createRoute(staffId: String) = "staff/topup/$staffId"
    }
}