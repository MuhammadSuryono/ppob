package com.yonotech.ppob.navigation

sealed class Screen(val route: String) {
    // Auth Screens
    object PhoneInput : Screen("phone_input")
    object OtpVerify : Screen("otp_verify")
    object SetCredentials : Screen("set_credentials")
    object PinLogin : Screen("pin_login")

    // Main Screens (Bottom Nav)
    object Home : Screen("home")
    object Transactions : Screen("transactions")
    object Wallet : Screen("wallet")
    object Profile : Screen("profile")

    // Transaction Flow
    object CategorySelection : Screen("category_selection")
    object ProductSelection : Screen("product_selection/{category_id}") {
        fun createRoute(categoryId: String) = "product_selection/$categoryId"
    }

    object TransactionConfirmation : Screen("transaction_confirmation/{product_id}") {
        fun createRoute(productId: String) = "transaction_confirmation/$productId"
    }

    object PinAuthorization : Screen("pin_authorization/{transaction_id}") {
        fun createRoute(transactionId: String) = "pin_authorization/$transactionId"
    }

    object TransactionResult : Screen("transaction_result/{transaction_id}") {
        fun createRoute(transactionId: String) = "transaction_result/$transactionId"
    }

    // Transaction Detail
    object TransactionDetail : Screen("transaction_detail/{transaction_id}") {
        fun createRoute(transactionId: String) = "transaction_detail/$transactionId"
    }

    // Staff Management (Mitra only)
    object StaffList : Screen("staff_list")
    object AddStaff : Screen("add_staff")
    object StaffDetail : Screen("staff_detail/{staff_id}") {
        fun createRoute(staffId: String) = "staff_detail/$staffId"
    }

    // Profile
    object Settings : Screen("settings")
    object DeviceManagement : Screen("device_management")
    object GantiPin : Screen("ganti_pin")
    object Bantuan : Screen("bantuan")
}