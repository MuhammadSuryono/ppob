package com.yonotech.ppob.navigation

import androidx.compose.animation.*
import androidx.compose.animation.core.tween
import androidx.compose.runtime.Composable
import androidx.lifecycle.viewmodel.compose.viewModel
import androidx.navigation.NavHostController
import androidx.navigation.NavType
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.navArgument
import com.yonotech.ppob.presentation.auth.*
import com.yonotech.ppob.presentation.components.EmptyStateScreen
import com.yonotech.ppob.presentation.home.HomeScreen
import com.yonotech.ppob.presentation.profile.*
import com.yonotech.ppob.presentation.staff.*
import com.yonotech.ppob.presentation.transaction.*
import com.yonotech.ppob.presentation.transactiondetail.TransactionDetailScreen
import com.yonotech.ppob.presentation.wallet.WalletScreen

@Composable
fun PPOBNavHost(
    navController: NavHostController,
    modifier: androidx.compose.ui.Modifier = androidx.compose.ui.Modifier
) {
    NavHost(
        navController = navController,
        startDestination = Screen.PhoneInput.route,
        modifier = modifier,
        enterTransition = { slideInHorizontally(initialOffsetX = { it }) },
        exitTransition = { slideOutHorizontally(targetOffsetX = { -it }) },
        popEnterTransition = { slideInHorizontally(initialOffsetX = { -it }) },
        popExitTransition = { slideOutHorizontally(targetOffsetX = { it }) }
    ) {
        // ========== AUTH FLOW ==========
        composable(Screen.PhoneInput.route) {
            val viewModel: AuthViewModel = viewModel()
            PhoneInputScreen(
                viewModel = viewModel,
                onNavigateToOtp = { phone ->
                    navController.navigate("${Screen.OtpVerify.route}?phone=${phone}")
                }
            )
        }

        composable(
            route = "${Screen.OtpVerify.route}?phone={phone}",
            arguments = listOf(navArgument("phone") { type = NavType.StringType; defaultValue = "" })
        ) { backStackEntry ->
            val phone = backStackEntry.arguments?.getString("phone") ?: ""
            val viewModel: AuthViewModel = viewModel()
            OtpVerifyScreen(
                viewModel = viewModel,
                phoneNumber = phone,
                onNavigateToSetCredentials = {
                    navController.navigate(Screen.SetCredentials.route)
                },
                onNavigateToPinLogin = {
                    navController.navigate(Screen.PinLogin.route)
                },
                onBack = { navController.popBackStack() }
            )
        }

        composable(Screen.SetCredentials.route) {
            val viewModel: AuthViewModel = viewModel()
            SetCredentialsScreen(
                viewModel = viewModel,
                onNavigateToHome = {
                    navController.navigate(Screen.Home.route) {
                        popUpTo(Screen.PhoneInput.route) { inclusive = true }
                    }
                },
                onBack = { navController.popBackStack() }
            )
        }

        composable(Screen.PinLogin.route) {
            val viewModel: AuthViewModel = viewModel()
            PinLoginScreen(
                viewModel = viewModel,
                onNavigateToHome = {
                    navController.navigate(Screen.Home.route) {
                        popUpTo(Screen.PhoneInput.route) { inclusive = true }
                    }
                },
                onNavigateToPhoneInput = {
                    navController.popBackStack(Screen.PhoneInput.route, inclusive = false)
                },
                onBack = { navController.popBackStack() }
            )
        }

        // ========== MAIN FLOW ==========
        composable(Screen.Home.route) {
            val viewModel: HomeViewModel = viewModel()
            HomeScreen(
                viewModel = viewModel,
                onNavigateToCategory = {
                    navController.navigate(Screen.CategorySelection.route)
                },
                onNavigateToTransactions = {
                    navController.navigate(Screen.Transactions.route)
                },
                onNavigateToWallet = {
                    navController.navigate(Screen.Wallet.route)
                },
                onNavigateToProfile = {
                    navController.navigate(Screen.Profile.route)
                }
            )
        }

        composable(Screen.Transactions.route) {
            val viewModel: TransactionHistoryViewModel = viewModel()
            TransactionListScreen(
                viewModel = viewModel,
                onTransactionClick = { txId ->
                    navController.navigate(Screen.TransactionDetail.createRoute(txId))
                },
                onFilterClicked = { /* TODO: Implement date/status filter */ }
            )
        }

        composable(Screen.Wallet.route) {
            val viewModel: WalletViewModel = viewModel()
            WalletScreen(
                viewModel = viewModel,
                onTopUpClicked = {
                    // Navigate to staff top-up for Mitra
                    navController.navigate(Screen.StaffList.route)
                },
                onBack = { navController.popBackStack() }
            )
        }

        composable(Screen.Profile.route) {
            val viewModel: ProfileViewModel = viewModel()
            ProfileScreen(
                viewModel = viewModel,
                onSettingsClick = { navController.navigate(Screen.Settings.route) },
                onDeviceManagementClick = { navController.navigate(Screen.DeviceManagement.route) },
                onGantiPinClick = { navController.navigate(Screen.GantiPin.route) },
                onBantuanClick = { navController.navigate(Screen.Bantuan.route) },
                onLogout = {
                    navController.navigate(Screen.PhoneInput.route) {
                        popUpTo(Screen.PhoneInput.route) { inclusive = true }
                    }
                }
            )
        }

        // ========== TRANSACTION FLOW ==========
        composable(Screen.CategorySelection.route) {
            val viewModel: CategoryViewModel = viewModel()
            CategorySelectionScreen(
                viewModel = viewModel,
                onCategoryClick = { categoryId ->
                    navController.navigate(Screen.ProductSelection.createRoute(categoryId))
                },
                onBack = { navController.popBackStack() }
            )
        }

        composable(
            route = Screen.ProductSelection.route,
            arguments = listOf(navArgument("category_id") { type = NavType.StringType })
        ) { backStackEntry ->
            val categoryId = backStackEntry.arguments?.getString("category_id") ?: ""
            val viewModel: ProductCatalogViewModel = viewModel()
            ProductSelectionScreen(
                viewModel = viewModel,
                categoryId = categoryId,
                onProductClick = { productId ->
                    navController.navigate(Screen.TransactionConfirmation.createRoute(productId))
                },
                onBack = { navController.popBackStack() }
            )
        }

        composable(
            route = Screen.TransactionConfirmation.route,
            arguments = listOf(navArgument("product_id") { type = NavType.StringType })
        ) { backStackEntry ->
            val productId = backStackEntry.arguments?.getString("product_id") ?: ""
            val viewModel: TransactionInitViewModel = viewModel()
            TransactionConfirmationScreen(
                viewModel = viewModel,
                productId = productId,
                onConfirm = {
                    navController.navigate(Screen.PinAuthorization.createRoute(it))
                },
                onBack = { navController.popBackStack() }
            )
        }

        composable(
            route = Screen.PinAuthorization.route,
            arguments = listOf(navArgument("transaction_id") { type = NavType.StringType })
        ) { backStackEntry ->
            val transactionId = backStackEntry.arguments?.getString("transaction_id") ?: ""
            val viewModel: TransactionInitViewModel = viewModel()
            PinAuthorizationScreen(
                viewModel = viewModel,
                transactionId = transactionId,
                onPinEntered = {
                    navController.navigate(Screen.TransactionResult.createRoute(transactionId))
                },
                onBack = { navController.popBackStack() }
            )
        }

        composable(
            route = Screen.TransactionResult.route,
            arguments = listOf(navArgument("transaction_id") { type = NavType.StringType })
        ) { backStackEntry ->
            val transactionId = backStackEntry.arguments?.getString("transaction_id") ?: ""
            TransactionResultScreen(
                transactionId = transactionId,
                onDone = {
                    navController.navigate(Screen.Home.route) {
                        popUpTo(Screen.Home.route) { inclusive = false }
                    }
                },
                onViewDetail = {
                    navController.navigate(Screen.TransactionDetail.createRoute(transactionId))
                }
            )
        }

        // ========== TRANSACTION DETAIL ==========
        composable(
            route = Screen.TransactionDetail.route,
            arguments = listOf(navArgument("transaction_id") { type = NavType.StringType })
        ) { backStackEntry ->
            val transactionId = backStackEntry.arguments?.getString("transaction_id") ?: ""
            val viewModel: TransactionDetailViewModel = viewModel()
            TransactionDetailScreen(
                viewModel = viewModel,
                transactionId = transactionId,
                onBack = { navController.popBackStack() }
            )
        }

        // ========== STAFF MANAGEMENT (MITRA) ==========
        composable(Screen.StaffList.route) {
            val viewModel: StaffListViewModel = viewModel()
            StaffListScreen(
                viewModel = viewModel,
                onStaffClick = { staffId ->
                    navController.navigate(Screen.StaffDetail.createRoute(staffId))
                },
                onAddStaffClick = {
                    navController.navigate(Screen.AddStaff.route)
                },
                onBack = { navController.popBackStack() }
            )
        }

        composable(Screen.AddStaff.route) {
            val viewModel: AddStaffViewModel = viewModel()
            AddStaffScreen(
                viewModel = viewModel,
                onSuccess = {
                    navController.popBackStack()
                },
                onCancel = { navController.popBackStack() }
            )
        }

        composable(
            route = Screen.StaffDetail.route,
            arguments = listOf(navArgument("staff_id") { type = NavType.StringType })
        ) { backStackEntry ->
            val staffId = backStackEntry.arguments?.getString("staff_id") ?: ""
            val viewModel: StaffDetailViewModel = viewModel()
            StaffDetailScreen(
                viewModel = viewModel,
                staffId = staffId,
                onTopUp = {
                    // Navigate to top-up modal (handled within StaffDetailScreen)
                },
                onBack = { navController.popBackStack() }
            )
        }

        // ========== PROFILE ==========
        composable(Screen.Settings.route) {
            SettingsScreen(
                onBack = { navController.popBackStack() }
            )
        }

        composable(Screen.DeviceManagement.route) {
            DeviceManagementScreen(
                onBack = { navController.popBackStack() }
            )
        }

        composable(Screen.GantiPin.route) {
            ChangePinScreen(
                onBack = { navController.popBackStack() }
            )
        }

        composable(Screen.Bantuan.route) {
            EmptyStateScreen(
                icon = null,
                message = "Halaman bantuan akan segera tersedia",
                cta = null
            )
        }
    }
}