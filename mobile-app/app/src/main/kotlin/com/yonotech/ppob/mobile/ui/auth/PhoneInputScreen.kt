package com.yonotech.ppob.mobile.ui.auth

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.hilt.navigation.compose.hiltViewModel
import com.yonotech.ppob.mobile.ui.components.PpoButton
import com.yonotech.ppob.mobile.ui.components.PpoTextField
import com.yonotech.ppob.mobile.utils.DeviceUtils
import com.yonotech.ppob.mobile.utils.Resource
import com.yonotech.ppob.mobile.viewmodels.auth.AuthViewModel

@Composable
fun PhoneInputScreen(
    onNavigateToOtp: (String) -> Unit,
    viewModel: AuthViewModel = hiltViewModel()
) {
    var phone by remember { mutableStateOf("") }
    val authState by viewModel.authState.collectAsState()
    val context = LocalContext.current

    LaunchedEffect(authState) {
        when (authState) {
            is Resource.Success -> {
                val data = (authState as Resource.Success).data
                if (data.requiresOtp) {
                    onNavigateToOtp(phone)
                }
                viewModel.resetState()
            }
            else -> {}
        }
    }

    Column(
        modifier = Modifier
            .fillMaxSize()
            .padding(24.dp),
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.Center
    ) {
        Text(
            text = "Masuk / Daftar",
            fontSize = 28.sp,
            fontWeight = FontWeight.Bold,
            color = MaterialTheme.colorScheme.primary
        )
        Text(
            text = "Masukkan nomor telepon Anda",
            fontSize = 16.sp,
            color = Color.Gray,
            modifier = Modifier.padding(top = 8.dp, bottom = 32.dp)
        )

        PpoTextField(
            value = phone,
            onValueChange = { 
                // Only allow digits, limit to reasonable length
                if (it.length <= 15 && it.all { char -> char.isDigit() || char == '+' }) phone = it 
            },
            label = "Nomor Telepon",
            keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Phone),
            placeholder = "+62xxxxxxxxxx"
        )

        Spacer(modifier = Modifier.height(32.dp))

        if (authState is Resource.Error) {
            Text(
                text = (authState as Resource.Error).message,
                color = MaterialTheme.colorScheme.error,
                modifier = Modifier.padding(bottom = 16.dp)
            )
        }

        PpoButton(
            label = "Lanjutkan",
            onClick = { 
                viewModel.sendOtp(phone, DeviceUtils.getDeviceId(context)) 
            },
            isLoading = authState is Resource.Loading,
            enabled = phone.length >= 10
        )

        Spacer(modifier = Modifier.height(16.dp))

        TextButton(onClick = { /* Handle existing account login */ }) {
            Text(text = "Sudah punya akun? Masuk")
        }
    }
}