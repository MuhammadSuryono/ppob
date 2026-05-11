package com.yonotech.ppob.mobile

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.padding
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.sp
import androidx.compose.ui.unit.dp
import androidx.navigation.NavController
import androidx.navigation.compose.*
import com.yonotech.ppob.mobile.ui.auth.*
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
                        startDestination = Screen.PhoneInput.route
                    ) {
                        // Auth Flow
                        composable(Screen.PhoneInput.route) {
                            PhoneInputScreen(
                                onNavigateToOtp = { identifier ->
                                    navController.navigate(Screen.Otp.createRoute(identifier, "login"))
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
                                        popUpTo(Screen.PhoneInput.route) { inclusive = true }
                                    }
                                },
                                onNewUserSetup = { phone ->
                                    navController.navigate(Screen.SetPasswordPin.createRoute(phone))
                                }
                            )
                        }
                        composable(Screen.SetPasswordPin.route) { backStackEntry ->
                            val phone = backStackEntry.arguments?.getString("phone") ?: ""
                            SetPasswordPinScreen(
                                phone = phone,
                                onAuthSuccess = {
                                    navController.navigate(Screen.Home.route) {
                                        popUpTo(Screen.PhoneInput.route) { inclusive = true }
                                    }
                                }
                            )
                        }
                        composable(Screen.PinLogin.route) {
                            PinLoginScreen(
                                onLoginSuccess = {
                                    navController.navigate(Screen.Home.route) {
                                        popUpTo(Screen.PhoneInput.route) { inclusive = true }
                                    }
                                },
                                onPasswordLogin = {
                                    // Navigate to login screen (would need to be added)
                                },
                                onNavigateToPhoneInput = {
                                    navController.popBackStack(Screen.PhoneInput.route, false)
                                }
                            )
                        }

                        // Main App with Bottom Navigation
                        composable(Screen.Home.route) {
                            MainScreen(navController) {
                                HomeScreen(
                                    onCategoryClick = { category ->
                                        navController.navigate(Screen.ProductList.createRoute(category.id))
                                    },
                                    onWalletClick = {
                                        navController.navigate(Screen.Wallet.route)
                                    }
                                )
                            }
                        }
                        composable(Screen.Transactions.route) {
                            MainScreen(navController) {
                                TransactionHistoryScreen(
                                    onBackClick = { navController.popBackStack() }
                                )
                            }
                        }
                        composable(Screen.Wallet.route) {
                            MainScreen(navController) {
                                WalletScreen(
                                    onBackClick = { navController.popBackStack() },
                                    onTransactionHistoryClick = {
                                        navController.navigate(Screen.Transactions.route)
                                    },
                                    onStaffClick = {
                                        navController.navigate(Screen.Staff.route)
                                    }
                                )
                            }
                        }
                        composable(Screen.Staff.route) {
                            MainScreen(navController) {
                                StaffListScreen(
                                    onBackClick = { navController.popBackStack() },
                                    onTopUpClick = { staff ->
                                        navController.navigate(Screen.StaffTopUp.createRoute(staff.id))
                                    }
                                )
                            }
                        }
                        composable(Screen.Profile.route) {
                            MainScreen(navController) {
                                ProfileScreen(
                                    onBackClick = { navController.popBackStack() }
                                )
                            }
                        }

                        // Other screens
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
                        composable(Screen.StaffTopUp.route) { backStackEntry ->
                            val staffId = backStackEntry.arguments?.getString("staffId") ?: ""
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
fun MainScreen(navController: NavController, content: @Composable () -> Unit) {
    val items = listOf(
        BottomNavItem("Beranda", Icons.Default.Home, Screen.Home),
        BottomNavItem("Transaksi", Icons.Default.ReceiptLong, Screen.Transactions),
        BottomNavItem("Wallet", Icons.Default.AccountBalanceWallet, Screen.Wallet),
        BottomNavItem("Staff", Icons.Default.People, Screen.Staff),
        BottomNavItem("Profil", Icons.Default.Person, Screen.Profile)
    )
    
    var selectedItem by remember { mutableStateOf(0) }
    
    Scaffold(
        bottomBar = {
            NavigationBar {
                items.forEachIndexed { index, item ->
                    NavigationBarItem(
                        icon = { Icon(item.icon, contentDescription = null) },
                        label = { Text(item.label) },
                        selected = selectedItem == index,
                        onClick = {
                            selectedItem = index
                            navController.navigate(item.screen.route) {
                                popUpTo(navController.graph.startDestinationId) {
                                    saveState = true
                                }
                                launchSingleTop = true
                                restoreState = true
                            }
                        }
                    )
                }
            }
        }
    ) { padding ->
        Box(modifier = Modifier.padding(padding)) {
            content()
        }
    }
}

data class BottomNavItem(val label: String, val icon: ImageVector, val screen: Screen)

@Composable
fun ProfileScreen(onBackClick: () -> Unit) {
    Column(
        modifier = Modifier
            .fillMaxSize()
            .padding(24.dp)
    ) {
        Text("Profil", fontSize = 24.sp, fontWeight = FontWeight.Bold)
        // Profile content would go here
    }
}

@Composable
fun Greeting(name: String, modifier: Modifier = Modifier) {
    Text(
        text = "Welcome to $name!",
        modifier = modifier
    )
}