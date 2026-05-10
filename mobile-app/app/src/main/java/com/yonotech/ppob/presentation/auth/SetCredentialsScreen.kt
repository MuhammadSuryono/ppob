package com.yonotech.ppob.presentation.auth

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.ArrowBack
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.text.input.PasswordVisualTransformation
import androidx.compose.ui.text.input.VisualTransformation
import androidx.compose.ui.unit.dp
import androidx.lifecycle.compose.collectAsStateWithLifecycle

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun SetCredentialsScreen(
    viewModel: AuthViewModel,
    onNavigateToHome: () -> Unit,
    onBack: () -> Unit
) {
    val uiState by viewModel.uiState.collectAsStateWithLifecycle()

    var passwordVisible by remember { mutableStateOf(false) }
    var pinVisible by remember { mutableStateOf(false) }

    LaunchedEffect(uiState.currentStep) {
        if (uiState.currentStep == AuthStep.COMPLETE) {
            onNavigateToHome()
        }
    }

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Buat Akun") },
                navigationIcon = {
                    IconButton(onClick = onBack) {
                        Icon(Icons.Default.ArrowBack, contentDescription = "Kembali")
                    }
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
                text = "Buat Password & PIN",
                style = MaterialTheme.typography.headlineMedium
            )

            Text(
                text = "Silakan buat password dan PIN transaksi untuk akun Anda",
                style = MaterialTheme.typography.bodyMedium
            )

            OutlinedTextField(
                value = uiState.password,
                onValueChange = { viewModel.onPasswordChange(it) },
                label = { Text("Password") },
                placeholder = { Text("Minimal 8 karakter") },
                visualTransformation = if (passwordVisible) VisualTransformation.None else PasswordVisualTransformation(),
                isError = !uiState.isPasswordValid && uiState.password.isNotEmpty(),
                supportingText = {
                    if (!uiState.isPasswordValid && uiState.password.isNotEmpty()) {
                        Text("Min 8 karakter, huruf besar, huruf kecil, dan angka")
                    }
                },
                keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Password),
                modifier = Modifier.fillMaxWidth()
            )

            OutlinedTextField(
                value = uiState.confirmPassword,
                onValueChange = { viewModel.onConfirmPasswordChange(it) },
                label = { Text("Konfirmasi Password") },
                visualTransformation = if (passwordVisible) VisualTransformation.None else PasswordVisualTransformation(),
                isError = uiState.confirmPassword.isNotEmpty() && uiState.confirmPassword != uiState.password,
                keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Password),
                modifier = Modifier.fillMaxWidth()
            )

            OutlinedTextField(
                value = uiState.pin,
                onValueChange = { viewModel.onPinChange(it) },
                label = { Text("PIN (6 digit)") },
                placeholder = { Text("123456") },
                visualTransformation = if (pinVisible) VisualTransformation.None else PasswordVisualTransformation(),
                isError = !uiState.isPinValid && uiState.pin.isNotEmpty(),
                supportingText = {
                    if (!uiState.isPinValid && uiState.pin.isNotEmpty()) {
                        Text("PIN harus 6 digit, tidak boleh urut atau sama")
                    }
                },
                keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.NumberPassword),
                modifier = Modifier.fillMaxWidth()
            )

            OutlinedTextField(
                value = uiState.confirmPin,
                onValueChange = { viewModel.onConfirmPinChange(it) },
                label = { Text("Konfirmasi PIN") },
                visualTransformation = if (pinVisible) VisualTransformation.None else PasswordVisualTransformation(),
                isError = uiState.confirmPin.isNotEmpty() && uiState.confirmPin != uiState.pin,
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
                onClick = { viewModel.register() },
                enabled = uiState.isPasswordValid &&
                        uiState.password == uiState.confirmPassword &&
                        uiState.isPinValid &&
                        uiState.pin == uiState.confirmPin &&
                        !uiState.isLoading,
                modifier = Modifier.fillMaxWidth()
            ) {
                Text(if (uiState.isLoading) "Membuat Akun..." else "Buat Akun")
            }
        }
    }
}
