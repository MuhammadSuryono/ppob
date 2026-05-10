package com.yonotech.ppob.presentation.profile

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.ArrowBack
import androidx.compose.material.icons.filled.Login
import androidx.compose.material3.Button
import androidx.compose.material3.Card
import androidx.compose.material3.CardDefaults
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.Icon
import androidx.compose.material3.IconButton
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Text
import androidx.compose.material3.TopAppBar
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp

@Composable
fun ProfileScreen(
    viewModel: com.yonotech.ppob.presentation.profile.ProfileViewModel = androidx.lifecycle.viewmodel.compose.viewModel(),
    onSettingsClick: () -> Unit,
    onDeviceManagementClick: () -> Unit,
    onGantiPinClick: () -> Unit,
    onBantuanClick: () -> Unit,
    onLogout: () -> Unit
) {
    val uiState by viewModel.uiState.collectAsState()

    LaunchedEffect(Unit) {
        val token = "" // TODO: Get from DataStore
        if (token.isNotEmpty()) {
            viewModel.setAuthToken(token)
            viewModel.loadProfile()
        }
    }

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Profil") },
                navigationIcon = {
                    IconButton(onClick = onLogout) {
                        Icon(Icons.Filled.ArrowBack, contentDescription = "Kembali")
                    }
                }
            )
        }
    ) { paddingValues ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues)
                .padding(24.dp),
            horizontalAlignment = Alignment.CenterHorizontally
        ) {
            if (uiState.isLoading) {
                CircularProgressIndicator()
                return@Column
            }

            // Profile Card
            Card(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(vertical = 8.dp),
                shape = RoundedCornerShape(16.dp),
                elevation = CardDefaults.cardElevation(defaultElevation = 4.dp)
            ) {
                Column(
                    modifier = Modifier.padding(24.dp),
                    horizontalAlignment = Alignment.CenterHorizontally
                ) {
                    Icon(
                        imageVector = Icons.Filled.AccountCircle,
                        contentDescription = "Avatar",
                        modifier = Modifier.size(64.dp),
                        tint = MaterialTheme.colorScheme.primary
                    )

                    Spacer(modifier = Modifier.padding(top = 12.dp))

                    Text(
                        text = uiState.userName,
                        style = MaterialTheme.typography.headlineMedium,
                        fontWeight = androidx.compose.ui.text.font.FontWeight.Bold
                    )

                    Text(
                        text = uiState.phoneNumber,
                        style = MaterialTheme.typography.bodyMedium,
                        color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.6f)
                    )

                    Text(
                        text = uiState.activeRole,
                        style = MaterialTheme.typography.bodySmall,
                        color = MaterialTheme.colorScheme.primary
                    )
                }
            }

            Spacer(modifier = Modifier.padding(vertical = 16.dp))

            // Menu Items
            ProfileMenuItem(
                icon = { Icon(Icons.Filled.Login, contentDescription = "Settings") },
                title = "Pengaturan",
                onClick = onSettingsClick
            )

            ProfileMenuItem(
                icon = { Icon(Icons.Filled.AccountCircle, contentDescription = "Devices") },
                title = "Perangkat Terpercaya",
                onClick = onDeviceManagementClick
            )

            ProfileMenuItem(
                icon = { Icon(Icons.Filled.Login, contentDescription = "PIN") },
                title = "Ganti PIN",
                onClick = onGantiPinClick
            )

            ProfileMenuItem(
                icon = { Icon(Icons.Filled.Login, contentDescription = "Help") },
                title = "Bantuan & Panduan",
                onClick = onBantuanClick
            )

            Spacer(modifier = Modifier.padding(vertical = 16.dp))

            Button(
                onClick = onLogout,
                modifier = Modifier.fillMaxWidth(),
                colors = ButtonDefaults.buttonColors(
                    containerColor = MaterialTheme.colorScheme.errorContainer
                )
            ) {
                Text("Keluar")
            }
        }
    }
}

@Composable
fun ProfileMenuItem(
    icon: @Composable () -> Unit,
    title: String,
    subtitle: String = "",
    onClick: () -> Unit
) {
    Card(
        modifier = Modifier
            .fillMaxWidth()
            .padding(vertical = 2.dp),
        shape = RoundedCornerShape(12.dp),
        onClick = onClick
    ) {
        Row(
            modifier = Modifier.padding(16.dp),
            horizontalArrangement = Arrangement.SpaceBetween,
            verticalAlignment = Alignment.CenterVertically
        ) {
            Row(verticalAlignment = Alignment.CenterVertically) {
                icon()
                Spacer(modifier = Modifier.padding(start = 16.dp))
                Column {
                    Text(text = title, style = MaterialTheme.typography.bodyLarge)
                    if (subtitle.isNotEmpty()) {
                        Text(
                            text = subtitle,
                            style = MaterialTheme.typography.bodySmall,
                            color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.6f)
                        )
                    }
                }
            }

            Icon(
                imageVector = Icons.Filled.ArrowForward,
                contentDescription = null,
                tint = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.4f)
            )
        }
    }
}