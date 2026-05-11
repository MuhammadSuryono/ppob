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
    onNavigateToOtp: (String, String, String) -> Unit, // requestId, phone, type
    onNavigateToPinLogin: (String) -> Unit,
    viewModel: AuthViewModel = hiltViewModel()
) {
    var phone by remember { mutableStateOf("") }
    val initiateState by viewModel.initiateState.collectAsState()
    val sendOtpState by viewModel.sendOtpState.collectAsState()
    val context = LocalContext.current

    LaunchedEffect(initiateState) {
        if (initiateState is Resource.Success) {
            val data = (initiateState as Resource.Success).data
            if (data.isTrusted) {
                // Trusted device: direct PIN login
                onNavigateToPinLogin(phone)
                viewModel.resetState()
            } else {
                // Requires OTP: send based on registration status
                val otpType = if (data.isRegistered) "login" else "register"
                viewModel.sendOtp(phone, otpType)
            }
        }
    }

    LaunchedEffect(sendOtpState) {
        if (sendOtpState is Resource.Success) {
            val data = (sendOtpState as Resource.Success).data
            val initiateData = (initiateState as? Resource.Success)?.data
            val type = if (initiateData?.isRegistered == false) "register" else "login"
            onNavigateToOtp(data.requestId, phone, type)
            viewModel.resetState()
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
            text = "Masukkan nomor telepon Anda untuk melanjutkan",
            fontSize = 16.sp,
            color = Color.Gray,
            modifier = Modifier.padding(top = 8.dp, bottom = 32.dp),
            textAlign = androidx.compose.ui.text.style.TextAlign.Center
        )

        PpoTextField(
            value = phone,
            onValueChange = { 
                if (it.length <= 15 && it.all { char -> char.isDigit() || char == '+' }) phone = it 
            },
            label = "Nomor Telepon",
            keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Phone),
            placeholder = "+62xxxxxxxxxx"
        )

        Spacer(modifier = Modifier.height(32.dp))

        if (sendOtpState is Resource.Error) {
            Text(
                text = (sendOtpState as Resource.Error).message,
                color = MaterialTheme.colorScheme.error,
                modifier = Modifier.padding(bottom = 16.dp)
            )
        }

        if (initiateState is Resource.Error) {
            Text(
                text = (initiateState as Resource.Error).message,
                color = MaterialTheme.colorScheme.error,
                modifier = Modifier.padding(bottom = 16.dp)
            )
        }

        PpoButton(
            label = "Lanjutkan",
            onClick = { 
                viewModel.initiateAuth(phone, DeviceUtils.getDeviceId(context)) 
            },
            isLoading = initiateState is Resource.Loading || sendOtpState is Resource.Loading,
            enabled = phone.length >= 10
        )
    }
}
