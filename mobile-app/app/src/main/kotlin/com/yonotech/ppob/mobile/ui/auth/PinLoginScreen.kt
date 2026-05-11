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
import androidx.compose.ui.text.input.PasswordVisualTransformation
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.hilt.navigation.compose.hiltViewModel
import com.yonotech.ppob.mobile.ui.components.PinDots
import com.yonotech.ppob.mobile.ui.components.PinPad
import com.yonotech.ppob.mobile.ui.components.PpoButton
import com.yonotech.ppob.mobile.ui.components.PpoTextField
import com.yonotech.ppob.mobile.utils.DeviceUtils
import com.yonotech.ppob.mobile.utils.Resource
import com.yonotech.ppob.mobile.viewmodels.auth.AuthViewModel

@Composable
fun PinLoginScreen(
    phoneArg: String = "",
    onLoginSuccess: () -> Unit,
    onPasswordLogin: () -> Unit,
    onNavigateToPhoneInput: () -> Unit,
    viewModel: AuthViewModel = hiltViewModel()
) {
    var phone by remember { mutableStateOf(phoneArg) }
    var pin by remember { mutableStateOf("") }
    val authState by viewModel.authState.collectAsState()
    val context = LocalContext.current

    LaunchedEffect(authState) {
        if (authState is Resource.Success) {
            onLoginSuccess()
            viewModel.resetState()
        }
    }

    Column(
        modifier = Modifier
            .fillMaxSize()
            .padding(24.dp),
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        Spacer(modifier = Modifier.height(32.dp))
        Text(
            text = "Masuk dengan PIN",
            fontSize = 28.sp,
            fontWeight = FontWeight.Bold,
            color = MaterialTheme.colorScheme.primary
        )
        Text(
            text = "Masukkan nomor telepon dan PIN Anda",
            fontSize = 16.sp,
            color = Color.Gray,
            modifier = Modifier.padding(top = 8.dp, bottom = 32.dp),
            textAlign = androidx.compose.ui.text.style.TextAlign.Center
        )

        PpoTextField(
            value = phone,
            onValueChange = { 
                if (phoneArg.isEmpty() && it.length <= 15 && it.all { char -> char.isDigit() || char == '+' }) phone = it 
            },
            label = "Nomor Telepon",
            keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Phone),
            placeholder = "+62xxxxxxxxxx"
        )

        Spacer(modifier = Modifier.height(32.dp))

        Text(text = "Masukkan 6-Digit PIN", fontSize = 16.sp, fontWeight = FontWeight.Medium)

        Spacer(modifier = Modifier.height(16.dp))

        PinDots(pin = pin)

        Spacer(modifier = Modifier.height(32.dp))

        if (authState is Resource.Error) {
            Text(
                text = (authState as Resource.Error).message,
                color = MaterialTheme.colorScheme.error,
                modifier = Modifier.padding(bottom = 16.dp)
            )
        }

        Spacer(modifier = Modifier.weight(1f))

        PinPad(
            onDigit = { digit -> if (pin.length < 6) pin += digit },
            onBackspace = { if (pin.isNotEmpty()) pin = pin.dropLast(1) }
        )

        Spacer(modifier = Modifier.height(24.dp))

        PpoButton(
            label = "Masuk",
            onClick = { viewModel.verifyPin(phone, pin, DeviceUtils.getDeviceId(context)) },
            isLoading = authState is Resource.Loading,
            enabled = phone.isNotEmpty() && pin.length == 6
        )

        Spacer(modifier = Modifier.height(8.dp))

        TextButton(onClick = onPasswordLogin) {
            Text(text = "Masuk dengan Password")
        }
    }
}