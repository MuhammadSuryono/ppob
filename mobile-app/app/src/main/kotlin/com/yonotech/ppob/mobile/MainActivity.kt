package com.yonotech.ppob.mobile

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Surface
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.remember
import androidx.compose.ui.Modifier
import androidx.compose.ui.tooling.preview.Preview
import androidx.hilt.navigation.compose.hiltViewModel
import androidx.navigation.NavType
import androidx.navigation.compose.NavHost
import androidx.navigation.compose.composable
import androidx.navigation.compose.rememberNavController
import androidx.navigation.navArgument
import com.yonotech.ppob.mobile.ui.auth.LoginScreen
import com.yonotech.ppob.mobile.ui.auth.OtpScreen
import com.yonotech.ppob.mobile.ui.auth.RegisterScreen
import com.yonotech.ppob.mobile.ui.history.TransactionHistoryScreen
import com.yonotech.ppob.mobile.ui.home.HomeScreen
import com.yonotech.ppob.mobile.ui.navigation.Screen
import com.yonotech.ppob.mobile.ui.product.ProductListScreen
import com.yonotech.ppob.mobile.ui.staff.StaffListScreen
import com.yonotech.ppob.mobile.ui.staff.StaffTopUpScreen
import com.yonotech.ppob.mobile.ui.transaction.TransactionConfirmScreen
import com.yonotech.ppob.mobile.ui.transaction.TransactionInitScreen
import com.yonotech.ppob.mobile.ui.transaction.TransactionResultScreen
import com.yonotech.ppob.mobile.ui.wallet.WalletScreen
import com.yonotech.ppob.mobile.ui.theme.PpoMobileTheme
import com.yonotech.ppob.mobile.viewmodels.transaction.TransactionViewModel
import com.yonotech.ppob.mobile.data.remote.dto.StaffDto
import dagger.hilt.android.AndroidEntryPoint

@AndroidEntryPoint
class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContent {
            PpoMobileTheme {
                val navController = rememberNavController()
                
                Surface(
                    modifier = Modifier.fillMaxSize(),
                    color = MaterialTheme.colorScheme.background
                ) {
                    NavHost(
                        navController = navController,
                        startDestination = Screen.Login.route
                    ) {
                        composable(Screen.Login.route) {
                            LoginScreen(
                                onLoginSuccess = {
                                    navController.navigate(Screen.Home.route) {
                                        popUpTo(Screen.Login.route) { inclusive = true }
                                    }
                                },
                                onRegisterClick = {
                                    navController.navigate(Screen.Register.route)
                                },
                                onNavigateToOtp = { identifier ->
                                    navController.navigate(Screen.Otp.createRoute(identifier, "login"))
                                }
                            )
                        }
                        composable(Screen.Register.route) {
                            RegisterScreen(
                                onRegisterSuccess = { phone ->
                                    navController.navigate(Screen.Otp.createRoute(phone, "registration"))
                                },
                                onLoginClick = {
                                    navController.popBackStack()
                                }
                            )
                        }
                        composable(Screen.Otp.route) { backStackEntry ->
                            val identifier = backStackEntry.arguments?.getString("identifier") ?: ""
                            val type = backStackEntry.arguments?.getString("type") ?: ""
                            OtpScreen(
                                identifier = identifier,
                                type = type,
                                onOtpSuccess = {
                                    navController.navigate(Screen.Home.route) {
                                        popUpTo(Screen.Login.route) { inclusive = true }
                                    }
                                }
                            )
                        }
                        composable(Screen.Home.route) {
                            HomeScreen(
                                onCategoryClick = { category ->
                                    navController.navigate(Screen.ProductList.createRoute(category.id))
                                },
                                onWalletClick = {
                                    navController.navigate(Screen.Wallet.route)
                                }
                            )
                        }
                        composable(Screen.ProductList.route) { backStackEntry ->
                            val categoryId = backStackEntry.arguments?.getString("categoryId") ?: ""
                            ProductListScreen(
                                categoryId = categoryId,
                                onProductClick = { product ->
                                    navController.navigate(Screen.TransactionInit.createRoute(product.id))
                                },
                                onBackClick = { navController.popBackStack() }
                            )
                        }
                        
                        // Transaction Flow with shared ViewModel
                        composable(Screen.TransactionInit.route) { backStackEntry ->
                            val productId = backStackEntry.arguments?.getString("productId") ?: ""
                            val viewModel: TransactionViewModel = hiltViewModel(backStackEntry)
                            TransactionInitScreen(
                                productId = productId,
                                onNext = { navController.navigate(Screen.TransactionConfirm.route) },
                                onBack = { navController.popBackStack() },
                                viewModel = viewModel
                            )
                        }
                        composable(Screen.TransactionConfirm.route) {
                            // Find the TransactionInit backstack entry to get the same ViewModel
                            val parentEntry = remember(it) {
                                navController.getBackStackEntry(Screen.TransactionInit.route)
                            }
                            val viewModel: TransactionViewModel = hiltViewModel(parentEntry)
                            TransactionConfirmScreen(
                                onSuccess = { txId ->
                                    navController.navigate(Screen.TransactionResult.createRoute(txId)) {
                                        popUpTo(Screen.Home.route) { inclusive = false }
                                    }
                                },
                                onBack = { navController.popBackStack() },
                                viewModel = viewModel
                            )
                        }
                        composable(Screen.TransactionResult.route) { backStackEntry ->
                            val txId = backStackEntry.arguments?.getString("txId") ?: ""
                            TransactionResultScreen(
                                txId = txId,
                                onFinish = {
                                    navController.navigate(Screen.Home.route) {
                                        popUpTo(Screen.Home.route) { inclusive = true }
                                    }
                                }
                            )
                        }
                        composable(Screen.Wallet.route) {
                            WalletScreen(
                                onBackClick = { navController.popBackStack() },
                                onTransactionHistoryClick = {
                                    navController.navigate(Screen.TransactionHistory.route)
                                },
                                onStaffClick = {
                                    navController.navigate(Screen.StaffList.route)
                                }
                            )
                        }
                        composable(Screen.TransactionHistory.route) {
                            TransactionHistoryScreen(
                                onBackClick = { navController.popBackStack() }
                            )
                        }
                        composable(Screen.StaffList.route) {
                            StaffListScreen(
                                onBackClick = { navController.popBackStack() },
                                onTopUpClick = { staff ->
                                    navController.navigate(Screen.StaffTopUp.createRoute(staff.id))
                                }
                            )
                        }
                        composable(Screen.StaffTopUp.route) { backStackEntry ->
                            val staffId = backStackEntry.arguments?.getString("staffId") ?: ""
                            // For now, using a placeholder staff - in real app, pass staff object via savedstatehandle
                            StaffTopUpScreen(
                                staff = StaffDto(
                                    id = staffId,
                                    name = "Staff",
                                    phone = "",
                                    email = "",
                                    balance = 0.0,
                                    dailyLimit = 0.0,
                                    dailyUsed = 0.0,
                                    marginScheme = "FixedAllowance",
                                    marginValue = 0.0,
                                    isActive = true
                                ),
                                onBackClick = { navController.popBackStack() }
                            )
                        }
                    }
                }
            }
        }
    }
}

@Composable
fun Greeting(name: String, modifier: Modifier = Modifier) {
    Text(
        text = "Welcome to $name!",
        modifier = modifier
    )
}

@Preview(showBackground = true)
@Composable
fun GreetingPreview() {
    PpoMobileTheme {
        Greeting("Android")
    }
}
