package com.yonotech.ppob.presentation.staff

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Spacer
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
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.unit.dp

@Composable
fun AddStaffScreen(
    viewModel: com.yonotech.ppob.presentation.staff.AddStaffViewModel = androidx.lifecycle.viewmodel.compose.viewModel(),
    onSuccess: () -> Unit,
    onCancel: () -> Unit
) {
    val uiState by viewModel.uiState.collectAsState()

    var phone by remember { mutableStateOf("") }
    var name by remember { mutableStateOf("") }
    var password by remember { mutableStateOf("") }
    var pin by remember { mutableStateOf("") }
    var marginScheme by remember { mutableStateOf("FixedAllowance") }
    var marginValue by remember { mutableStateOf("10000") }

    LaunchedEffect(uiState.isSuccess) {
        if (uiState.isSuccess) {
            onSuccess()
        }
    }

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Tambah Staff") },
                navigationIcon = {
                    IconButton(onClick = onCancel) {
                        androidx.compose.material.icons.Icons.Filled.ArrowBack.let { /* Back */ }
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
            OutlinedTextField(
                value = phone,
                onValueChange = { phone = it },
                label = { Text("Nomor Telepon Staff") },
                placeholder = { Text("+628xxxxxxxxxx") },
                keyboardOptions = androidx.compose.foundation.text.KeyboardOptions(keyboardType = KeyboardType.Phone),
                modifier = Modifier.fillMaxWidth()
            )

            OutlinedTextField(
                value = name,
                onValueChange = { name = it },
                label = { Text("Nama Lengkap") },
                modifier = Modifier.fillMaxWidth()
            )

            OutlinedTextField(
                value = password,
                onValueChange = { password = it },
                label = { Text("Password") },
                visualTransformation = androidx.compose.ui.text.input.PasswordVisualTransformation(),
                keyboardOptions = androidx.compose.foundation.text.KeyboardOptions(keyboardType = KeyboardType.Password),
                modifier = Modifier.fillMaxWidth()
            )

            OutlinedTextField(
                value = pin,
                onValueChange = { pin = it },
                label = { Text("PIN (6 digit)") },
                placeholder = { Text("123456") },
                visualTransformation = androidx.compose.ui.text.input.PasswordVisualTransformation(),
                keyboardOptions = androidx.compose.foundation.text.KeyboardOptions(keyboardType = KeyboardType.NumberPassword),
                modifier = Modifier.fillMaxWidth()
            )

            Text(
                text = "Skema Margin",
                style = MaterialTheme.typography.titleMedium
            )

            // Margin Scheme Selection
            androidx.compose.material3.RadioButton(
                selected = marginScheme == "FixedAllowance",
                onClick = { marginScheme = "FixedAllowance" }
            )
            Text("Fixed Allowance (Rp/transaksi)")

            androidx.compose.material3.RadioButton(
                selected = marginScheme == "MarginShare",
                onClick = { marginScheme = "MarginShare" }
            )
            Text("Margin Share (%)")

            if (marginScheme == "FixedAllowance") {
                OutlinedTextField(
                    value = marginValue,
                    onValueChange = { marginValue = it },
                    label = { Text("Allowance per transaksi (Rp)") },
                    keyboardOptions = androidx.compose.foundation.text.KeyboardOptions(keyboardType = KeyboardType.Number),
                    modifier = Modifier.fillMaxWidth()
                )
            } else {
                OutlinedTextField(
                    value = marginValue,
                    onValueChange = { marginValue = it },
                    label = { Text("Persentase Margin (%)") },
                    keyboardOptions = androidx.compose.foundation.text.KeyboardOptions(keyboardType = KeyboardType.Number),
                    modifier = Modifier.fillMaxWidth()
                )
            }

            uiState.error?.let { error ->
                Text(
                    text = error,
                    color = MaterialTheme.colorScheme.error
                )
            }

            Spacer(modifier = Modifier.padding(8.dp))

            Button(
                onClick = {
                    val token = "" // TODO: Get from DataStore
                    viewModel.addStaff(
                        phone, name, password, pin,
                        marginScheme,
                        marginValue.toDoubleOrNull() ?: 10000.0
                    )
                },
                enabled = phone.isNotEmpty() && name.isNotEmpty() && pin.length == 6 && !uiState.isLoading,
                modifier = Modifier.fillMaxWidth()
            ) {
                Text(if (uiState.isLoading) "Menambahkan..." else "Tambah Staff")
            }
        }
    }
}