package com.yonotech.ppob.mobile.ui.splash

import androidx.compose.animation.core.tween
import androidx.compose.foundation.Image
import androidx.compose.foundation.layout.*
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.runtime.setValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.alpha
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.res.painterResource
import androidx.compose.ui.unit.dp
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Text
import androidx.compose.runtime.collectAsState
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.sp
import androidx.hilt.navigation.compose.hiltViewModel
import com.yonotech.ppob.mobile.R
import com.yonotech.ppob.mobile.viewmodels.SplashViewModel
import kotlinx.coroutines.delay

@Composable
fun SplashScreen(
    onNavigateToHome: () -> Unit,
    onNavigateToLogin: () -> Unit,
    viewModel: SplashViewModel = hiltViewModel()
) {
    val isLoggedIn by viewModel.isLoggedIn.collectAsState()

    LaunchedEffect(isLoggedIn) {
        delay(2000) // Branding delay
        when (isLoggedIn) {
            true -> onNavigateToHome()
            false -> onNavigateToLogin()
            else -> {} // Wait for state
        }
    }

    Box(
        modifier = Modifier.fillMaxSize(),
        contentAlignment = Alignment.Center
    ) {
        Column(horizontalAlignment = Alignment.CenterHorizontally) {
            Image(
                painter = painterResource(id = R.drawable.ic_logo_app_background),
                contentDescription = "Logo",
                modifier = Modifier.size(120.dp)
            )
            Spacer(modifier = Modifier.height(16.dp))
            Text(
                text = "PPOB Mobile",
                fontSize = 24.sp,
                fontWeight = FontWeight.Bold,
                color = MaterialTheme.colorScheme.primary
            )
        }
    }
}