package com.yonotech.ppob

import android.annotation.SuppressLint
import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.compose.animation.AnimatedVisibility
import androidx.compose.animation.fadeIn
import androidx.compose.animation.fadeOut
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Scaffold
import androidx.compose.runtime.Composable
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.unit.dp
import androidx.navigation.compose.currentBackStackEntryAsState
import androidx.navigation.compose.rememberNavController
import com.yonotech.ppob.navigation.PPOBNavHost
import com.yonotech.ppob.theme.PPOBTheme
import dagger.hilt.android.AndroidEntryPoint

@AndroidEntryPoint
class MainActivity : ComponentActivity() {

    @SuppressLint("UnusedMaterial3ScaffoldPaddingParameter")
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        enableEdgeToEdge()

        setContent {
            PPOBMainApp()
        }
    }
}

@Composable
fun PPOBMainApp() {
    var showSplash by remember { mutableStateOf(true) }

    PPOBTheme {
        if (showSplash) {
            SplashScreen(onSplashFinished = { showSplash = false })
        } else {
            MainNavigation()
        }
    }
}

@Composable
fun SplashScreen(onSplashFinished: () -> Unit) {
    Box(
        modifier = Modifier
            .fillMaxSize()
            .background(Color(0xFF1B5E20))
    ) {
        // Logo or branding here
        androidx.compose.material3.Text(
            text = "PPOB",
            color = Color.White,
            style = androidx.compose.material3.MaterialTheme.typography.displayLarge,
            modifier = Modifier.padding(16.dp)
        )
    }

    androidx.compose.runtime.LaunchedEffect(Unit) {
        kotlinx.coroutines.delay(2000)
        onSplashFinished()
    }
}

@Composable
fun MainNavigation() {
    val navController = rememberNavController()
    val currentBackStack by navController.currentBackStackEntryAsState()
    val currentDestination = currentBackStack?.destination?.route

    val showBottomBar = when (currentDestination) {
        "home", "transactions", "wallet", "staff", "profile" -> true
        else -> false
    }

    Scaffold(
        modifier = Modifier.fillMaxSize(),
        bottomBar = {
            AnimatedVisibility(
                visible = showBottomBar,
                enter = fadeIn(),
                exit = fadeOut()
            ) {
                com.yonotech.ppob.navigation.BottomNavigationBar(navController = navController)
            }
        }
    ) { innerPadding ->
        PPOBNavHost(
            navController = navController,
            modifier = Modifier.padding(innerPadding)
        )
    }
}
