package com.yonotech.ppob.mobile.ui.auth

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.text.input.PasswordVisualTransformation
import androidx.compose.ui.text.input.VisualTransformation
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.hilt.navigation.compose.hiltViewModel
import com.yonotech.ppob.mobile.ui.components.PpoButton
import com.yonotech.ppob.mobile.ui.components.PpoTextField
import com.yonotech.ppob.mobile.ui.theme.PpoMobileTheme
import com.yonotech.ppob.mobile.utils.DeviceUtils
import com.yonotech.ppob.mobile.utils.Resource
import com.yonotech.ppob.mobile.viewmodels.auth.AuthViewModel

@Composable
fun PasswordLoginScreen(
    phone: String,
    requestId: String,
    onLoginSuccess: () -> Unit,
    viewModel: AuthViewModel = hiltViewModel()
) {
    val authState by viewModel.authState.collectAsState()
    val context = LocalContext.current

    LaunchedEffect(authState) {
        if (authState is Resource.Success) {
            onLoginSuccess()
            viewModel.resetState()
        }
    }

    PasswordLoginContent(
        phone = phone,
        authState = authState,
        onVerifyPassword = { password ->
            viewModel.verifyPassword(phone, password, DeviceUtils.getDeviceId(context), requestId)
        }
    )
}

@Composable
fun PasswordLoginContent(
    phone: String,
    authState: Resource<Any>,
    onVerifyPassword: (String) -> Unit
) {
    var password by remember { mutableStateOf("") }
    var showPassword by remember { mutableStateOf(false) }

    Column(
        modifier = Modifier
            .fillMaxSize()
            .padding(24.dp),
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.Center
    ) {
        Text(
            text = "Masukkan Password",
            fontSize = 28.sp,
            fontWeight = FontWeight.Bold,
            color = MaterialTheme.colorScheme.primary
        )
        Text(
            text = "Silakan masukkan password untuk nomor $phone",
            fontSize = 16.sp,
            color = MaterialTheme.colorScheme.onSurfaceVariant,
            modifier = Modifier.padding(top = 8.dp, bottom = 32.dp),
            textAlign = androidx.compose.ui.text.style.TextAlign.Center
        )

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

        Spacer(modifier = Modifier.height(32.dp))

        if (authState is Resource.Error) {
            Text(
                text = (authState as Resource.Error).message,
                color = MaterialTheme.colorScheme.error,
                modifier = Modifier.padding(bottom = 16.dp)
            )
        }

        PpoButton(
            label = "Masuk",
            onClick = { onVerifyPassword(password) },
            isLoading = authState is Resource.Loading,
            enabled = password.isNotEmpty()
        )
    }
}

@Preview(showBackground = true)
@Composable
fun PasswordLoginScreenPreview() {
    PpoMobileTheme {
        PasswordLoginContent(
            phone = "08123456789",
            authState = Resource.Idle,
            onVerifyPassword = {}
        )
    }
}
