package com.yonotech.ppob.presentation.profile

import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.material3.TopAppBar
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp

@Composable
fun SettingsScreen(onBack: () -> Unit) {
    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Pengaturan") },
                navigationIcon = {
                    IconButton(onClick = onBack) {
                        androidx.compose.material.icons.Icons.Filled.ArrowBack.let { /* Back */ }
                    }
                }
            )
        }
    ) { paddingValues ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues)
                .padding(24.dp)
        ) {
            Text(text = "Pengaturan aplikasi akan ditampilkan di sini")
        }
    }
}

@Composable
fun DeviceManagementScreen(onBack: () -> Unit) {
    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Perangkat Terpercaya") },
                navigationIcon = {
                    androidx.compose.material.icons.Icons.Filled.ArrowBack.let { /* Back */ }
                }
            )
        }
    ) { paddingValues ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues)
                .padding(24.dp)
        ) {
            Text(text = "Daftar perangkat terpercaya")
        }
    }
}

@Composable
fun ChangePinScreen(onBack: () -> Unit) {
    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Ganti PIN") },
                navigationIcon = {
                    androidx.compose.material.icons.Icons.Filled.ArrowBack.let { /* Back */ }
                }
            )
        }
    ) { paddingValues ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues)
                .padding(24.dp)
        ) {
            Text(text = "Ubah PIN transaksi Anda")
        }
    }
}

@Composable
fun BantuanScreen(onBack: () -> Unit) {
    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Bantuan") },
                navigationIcon = {
                    androidx.compose.material.icons.Icons.Filled.ArrowBack.let { /* Back */ }
                }
            )
        }
    ) { paddingValues ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues)
                .padding(24.dp)
        ) {
            Text(text = "Halaman bantuan akan segera tersedia")
        }
    }
}