package com.yonotech.ppob.presentation.auth

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.text.KeyboardOptions
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
import androidx.compose.ui.text.input.PasswordVisualTransformation
import androidx.compose.ui.unit.dp

@Composable
fun PinLoginScreen(
    viewModel: AuthViewModel = androidx.lifecycle.viewmodel.compose.viewModel(),
    onNavigateToHome: () -> Unit,
    onNavigateToPhoneInput: () -> Unit,
    onBack: () -> Unit
) {
    val uiState by viewModel.uiState.collectAsState()

    LaunchedEffect(uiState.currentStep) {
        if (uiState.currentStep == AuthStep.COMPLETE) {
            onNavigateToHome()
        }
    }

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Masuk dengan PIN") },
                navigationIcon = {
                    androidx.compose.material.icons.Icons.Default.ArrowBack.let { /* Back */ }
                }
            )
        }
    ) { paddingValues ->
        Column(
            modifier = Modifier
                .fillMaxWidth()
                .padding(paddingValues)
                .padding(24.dp),
            verticalArrangement = Arrangement.spacedBy(16.dp)
        ) {
            Text(
                text = "Masukkan PIN Anda",
                style = MaterialTheme.typography.headlineMedium
            )

            OutlinedTextField(
                value = uiState.pin,
                onValueChange = { viewModel.onPinChange(it) },
                label = { Text("PIN") },
                placeholder = { Text("123456") },
                visualTransformation = PasswordVisualTransformation(),
                isError = !uiState.isPinValid && uiState.pin.isNotEmpty(),
                keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.NumberPassword),
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
                onClick = { viewModel.login() },
                enabled = uiState.isPinValid && !uiState.isLoading,
                modifier = Modifier.fillMaxWidth()
            ) {
                Text(if (uiState.isLoading) "Masuk..." else "Masuk")
            }

            Text(
                text = "Lupa PIN? Login dengan password",
                style = MaterialTheme.typography.labelMedium,
                modifier = Modifier.padding(top = 8.dp)
                // TODO: Navigate to phone input for password login
            )
        }
    }
}