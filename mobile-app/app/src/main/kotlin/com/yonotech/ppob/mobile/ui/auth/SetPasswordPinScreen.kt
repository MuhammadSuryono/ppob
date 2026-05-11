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
import androidx.compose.ui.text.input.VisualTransformation
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
fun SetPasswordPinScreen(
    phone: String,
    requestId: String,
    onRegisterSuccess: (Int) -> Unit, // userId
    viewModel: AuthViewModel = hiltViewModel()
) {
    var email by remember { mutableStateOf("") }
    var name by remember { mutableStateOf("") }
    var password by remember { mutableStateOf("") }
    var confirmPassword by remember { mutableStateOf("") }
    var pin by remember { mutableStateOf("") }
    var showPassword by remember { mutableStateOf(false) }
    
    val authState by viewModel.authState.collectAsState()
    val context = LocalContext.current

    LaunchedEffect(authState) {
        if (authState is Resource.Success) {
            val data = (authState as Resource.Success).data
            onRegisterSuccess(data.userId ?: 0)
            viewModel.resetState()
        }
    }

    Column(
        modifier = Modifier
            .fillMaxSize()
            .padding(24.dp),
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        Text(
            text = "Daftar Akun Baru",
            fontSize = 28.sp,
            fontWeight = FontWeight.Bold,
            color = MaterialTheme.colorScheme.primary
        )
        Text(
            text = "Lengkapi data diri Anda untuk membuat akun",
            fontSize = 16.sp,
            color = Color.Gray,
            modifier = Modifier.padding(top = 8.dp, bottom = 24.dp),
            textAlign = androidx.compose.ui.text.style.TextAlign.Center
        )

        PpoTextField(value = name, onValueChange = { name = it }, label = "Nama Lengkap")
        Spacer(modifier = Modifier.height(8.dp))
        PpoTextField(value = email, onValueChange = { email = it }, label = "Email", keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Email))
        Spacer(modifier = Modifier.height(8.dp))

        PpoTextField(
            value = password,
            onValueChange = { password = it },
            label = "Password",
            visualTransformation = if (showPassword) VisualTransformation.None else PasswordVisualTransformation(),
            keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Password),
            trailingIcon = {
                TextButton(onClick = { showPassword = !showPassword }) {
                    Text(if (showPassword) "Lihat" else "Tutup")
                }
            }
        )

        Spacer(modifier = Modifier.height(16.dp))

        Text(text = "Buat 6-Digit PIN", fontSize = 16.sp, fontWeight = FontWeight.Medium)
        Spacer(modifier = Modifier.height(8.dp))
        PinDots(pin = pin)

        Spacer(modifier = Modifier.weight(1f))

        PinPad(
            onDigit = { digit -> if (pin.length < 6) pin += digit },
            onBackspace = { if (pin.isNotEmpty()) pin = pin.dropLast(1) }
        )

        Spacer(modifier = Modifier.height(24.dp))

        if (authState is Resource.Error) {
            Text(
                text = (authState as Resource.Error).message,
                color = MaterialTheme.colorScheme.error,
                modifier = Modifier.padding(bottom = 16.dp)
            )
        }

        PpoButton(
            label = "Daftar Sekarang",
            onClick = { 
                viewModel.register(email, phone, name, password, pin, DeviceUtils.getDeviceId(context), requestId) 
            },
            isLoading = authState is Resource.Loading,
            enabled = name.isNotEmpty() && email.isNotEmpty() && password.isNotEmpty() && pin.length == 6
        )
    }
}
