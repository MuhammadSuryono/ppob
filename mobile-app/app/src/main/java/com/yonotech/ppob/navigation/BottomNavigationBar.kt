package com.yonotech.ppob.navigation

import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.AccountCircle
import androidx.compose.material.icons.filled.Home
import androidx.compose.material.icons.filled.Receipt
import androidx.compose.material.icons.filled.Wallet
import androidx.compose.material3.BadgedBox
import androidx.compose.material3.Icon
import androidx.compose.material3.NavigationBar
import androidx.compose.material3.NavigationBarItem
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.navigation.NavController
import androidx.navigation.NavDestination.Companion.hierarchy
import androidx.navigation.compose.currentBackStackEntryAsState

data class BottomNavItem(
    val route: String,
    val icon: ImageVector,
    val label: String,
    var badgeCount: Int? = null
)

@Composable
fun BottomNavigationBar(navController: NavController) {
    val navBackStackEntry by navController.currentBackStackEntryAsState()
    val currentDestination = navBackStackEntry?.destination

    val items = listOf(
        BottomNavItem(
            route = "home",
            icon = Icons.Filled.Home,
            label = "Beranda"
        ),
        BottomNavItem(
            route = "transactions",
            icon = Icons.Filled.Receipt,
            label = "Transaksi"
        ),
        BottomNavItem(
            route = "wallet",
            icon = Icons.Filled.Wallet,
            label = "Dompet"
        ),
        BottomNavItem(
            route = "profile",
            icon = Icons.Filled.AccountCircle,
            label = "Profil"
        )
    )

    NavigationBar {
        items.forEach { item ->
            val isSelected = currentDestination?.hierarchy?.any { it.route == item.route } == true
            NavigationBarItem(
                selected = isSelected,
                onClick = {
                    if (!isSelected) {
                        navController.navigate(item.route) {
                            popUpTo(navController.graph.startDestinationId) {
                                saveState = true
                            }
                            launchSingleTop = true
                            restoreState = true
                        }
                    }
                },
                icon = {
                    BadgedBox(
                        badge = {
                            item.badgeCount?.let { count ->
                                // TODO: Implement BadgeBox when available
                            }
                        }
                    ) {
                        Icon(
                            imageVector = item.icon,
                            contentDescription = item.label
                        )
                    }
                },
                label = { Text(text = item.label) }
            )
        }
    }
}
