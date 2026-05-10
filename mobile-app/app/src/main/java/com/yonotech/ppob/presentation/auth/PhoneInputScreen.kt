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
fun PhoneInputScreen(
    viewModel: AuthViewModel = androidx.lifecycle.viewmodel.compose.viewModel(),
    onNavigateToOtp: (String) -> Unit
) {
    val uiState by viewModel.uiState.collectAsState()

    Scaffold(
        topBar = {
            TopAppBar(title = { Text("Masuk / Daftar") })
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
                text = "Selamat datang di PPOB",
                style = MaterialTheme.typography.displaySmall
            )

            Text(
                text = "Masukkan nomor telepon Anda untuk memulai",
                style = MaterialTheme.typography.bodyLarge
            )

            OutlinedTextField(
                value = uiState.phoneNumber,
                onValueChange = { viewModel.onPhoneNumberChange(it) },
                label = { Text("Nomor Telepon") },
                placeholder = { Text("+628xxxxxxxxxx") },
                keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Phone),
                isError = !uiState.isPhoneValid,
                supportingText = {
                    if (!uiState.isPhoneValid && uiState.phoneNumber.isNotEmpty()) {
                        Text("Masukkan nomor telepon yang valid")
                    }
                },
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
                    if (uiState.isPhoneValid) {
                        viewModel.sendOtp()
                    }
                },
                enabled = uiState.isPhoneValid && !uiState.isLoading,
                modifier = Modifier.fillMaxWidth()
            ) {
                Text(if (uiState.isLoading) "Mengirim..." else "Lanjutkan")
            }
        }
    }

    LaunchedEffect(uiState.currentStep) {
        if (uiState.currentStep == AuthStep.OTP_VERIFY) {
            onNavigateToOtp(uiState.phoneNumber)
        }
    }
}