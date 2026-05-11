package com.yonotech.ppob.mobile.ui.navigation

sealed class Screen(val route: String) {
    object Login : Screen("login")
    object Register : Screen("register")
    object Otp : Screen("otp/{identifier}/{type}") {
        fun createRoute(identifier: String, type: String) = "otp/$identifier/$type"
    }
    object Home : Screen("home")
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
    object Wallet : Screen("wallet")
    object TransactionHistory : Screen("transaction/history")
    object StaffList : Screen("staff/list")
    object StaffTopUp : Screen("staff/topup/{staffId}") {
        fun createRoute(staffId: String) = "staff/topup/$staffId"
    }
}
