package com.yonotech.ppob.presentation.auth

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.foundation.verticalScroll
import androidx.compose.material3.Button
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.material3.TopAppBar
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.unit.dp

@Composable
fun OtpVerifyScreen(
    viewModel: AuthViewModel = androidx.lifecycle.viewmodel.compose.viewModel(),
    phoneNumber: String,
    onNavigateToSetCredentials: () -> Unit,
    onNavigateToPinLogin: () -> Unit,
    onBack: () -> Unit
) {
    val uiState by viewModel.uiState.collectAsState()

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Verifikasi OTP") },
                navigationIcon = {
                    androidx.compose.material.icons.Icons.Default.ArrowBack.let { /* Back icon */ }
                }
            )
        }
    ) { paddingValues ->
        Column(
            modifier = Modifier
                .fillMaxWidth()
                .padding(paddingValues)
                .padding(24.dp)
                .verticalScroll(rememberScrollState()),
            verticalArrangement = Arrangement.spacedBy(16.dp)
        ) {
            Text(
                text = "Kode OTP telah dikirim ke $phoneNumber",
                style = MaterialTheme.typography.bodyLarge
            )

            OutlinedTextField(
                value = uiState.otpCode,
                onValueChange = { viewModel.onOtpChange(it) },
                label = { Text("Kode OTP") },
                placeholder = { Text("123456") },
                keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.NumberPassword),
                isError = !uiState.isOtpValid,
                modifier = Modifier.fillMaxWidth()
            )

            uiState.error?.let { error ->
                Text(
                    text = error,
                    color = MaterialTheme.colorScheme.error,
                    style = MaterialTheme.typography.bodyMedium
                )
            }

            Button(
                onClick = {
                    if (uiState.isOtpValid) {
                        viewModel.verifyOtp()
                    }
                },
                enabled = uiState.isOtpValid && !uiState.isLoading,
                modifier = Modifier.fillMaxWidth()
            ) {
                Text(if (uiState.isLoading) "Memverifikasi..." else "Verifikasi")
            }

            Text(
                text = "Kirim ulang OTP",
                style = MaterialTheme.typography.labelMedium,
                modifier = Modifier.padding(top = 8.dp)
                // TODO: Add click handler for resend
            )
        }
    }

    LaunchedEffect(uiState.currentStep) {
        when (uiState.currentStep) {
            AuthStep.SET_CREDENTIALS -> onNavigateToSetCredentials()
            AuthStep.PIN_LOGIN -> onNavigateToPinLogin()
            AuthStep.COMPLETE -> onNavigateToPinLogin()
            else -> {}
        }
    }
}