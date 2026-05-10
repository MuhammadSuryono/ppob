package com.yonotech.ppob.presentation.staff

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.Card
import androidx.compose.material3.CardDefaults
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Switch
import androidx.compose.material3.Text
import androidx.compose.material3.TopAppBar
import androidx.compose.runtime.Composable
import androidx.compose.runtime.LaunchedEffect
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.yonotech.ppob.domain.model.Staff

@Composable
fun StaffDetailScreen(
    viewModel: com.yonotech.ppob.presentation.staff.StaffDetailViewModel = androidx.lifecycle.viewmodel.compose.viewModel(),
    staffId: String,
    onTopUp: () -> Unit,
    onBack: () -> Unit
) {
    val uiState by viewModel.uiState.collectAsState()

    LaunchedEffect(staffId) {
        viewModel.loadStaffDetail(staffId)
    }

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Detail Staff") },
                navigationIcon = {
                    androidx.compose.material.icons.Icons.Filled.ArrowBack.let { /* Back */ }
                }
            )
        }
    ) { paddingValues ->
        Column(
            modifier = Modifier
                .fillMaxWidth()
                .padding(paddingValues)
                .padding(24.dp)
        ) {
            if (uiState.isLoading) {
                Text("Memuat...")
                return@Column
            }

            uiState.staff?.let { staff ->
                StaffDetailCard(staff = staff)

                Spacer(modifier = Modifier.padding(vertical = 16.dp))

                // Active Toggle
                Row(
                    modifier = Modifier.fillMaxWidth(),
                    horizontalArrangement = Arrangement.SpaceBetween,
                    verticalAlignment = androidx.compose.ui.Alignment.CenterVertically
                ) {
                    Text(
                        text = "Status Aktif",
                        style = MaterialTheme.typography.bodyLarge
                    )
                    Switch(
                        checked = staff.isActive,
                        onCheckedChange = {
                            viewModel.updateStaff(
                                staffId = staffId,
                                isActive = it
                            )
                        }
                    )
                }

                Spacer(modifier = Modifier.padding(vertical = 8.dp))

                // Top Up Button
                Card(
                    modifier = Modifier.fillMaxWidth(),
                    shape = androidx.compose.foundation.shape.RoundedCornerShape(12.dp),
                    colors = CardDefaults.cardColors(
                        containerColor = MaterialTheme.colorScheme.primary
                    ),
                    onClick = onTopUp
                ) {
                    Column(
                        modifier = Modifier.padding(16.dp),
                        horizontalAlignment = androidx.compose.ui.Alignment.CenterHorizontally
                    ) {
                        Text(
                            text = "Top Up Staff",
                            style = MaterialTheme.typography.titleMedium,
                            color = androidx.compose.ui.graphics.Color.White
                        )
                    }
                }
            }
        }
    }
}

@Composable
fun StaffDetailCard(staff: Staff) {
    Card(
        modifier = Modifier.fillMaxWidth(),
        shape = androidx.compose.foundation.shape.RoundedCornerShape(16.dp),
        elevation = CardDefaults.cardElevation(defaultElevation = 4.dp)
    ) {
        Column(modifier = Modifier.padding(20.dp)) {
            Text(
                text = staff.name,
                style = MaterialTheme.typography.headlineSmall,
                fontWeight = androidx.compose.ui.text.font.FontWeight.Bold
            )

            Spacer(modifier = Modifier.padding(top = 4.dp))

            Text(
                text = staff.phoneNumber,
                style = MaterialTheme.typography.bodyMedium
            )

            Spacer(modifier = Modifier.padding(top = 8.dp))

            Row(modifier = Modifier.fillMaxWidth()) {
                Column(modifier = androidx.compose.foundation.layout.weight(1f)) {
                    Text(
                        text = "Saldo",
                        style = MaterialTheme.typography.bodySmall
                    )
                    Text(
                        text = "Rp %,d".format(staff.walletBalance.toInt()),
                        style = MaterialTheme.typography.titleLarge,
                        fontWeight = androidx.compose.ui.text.font.FontWeight.Bold
                    )
                }

                Column(horizontalAlignment = androidx.compose.ui.Alignment.End) {
                    Text(
                        text = "Transaksi Hari Ini",
                        style = MaterialTheme.typography.bodySmall
                    )
                    Text(
                        text = "${staff.dailyTxnCount} / Rp %,d".format(staff.dailyTxnAmount.toInt()),
                        style = MaterialTheme.typography.bodyMedium
                    )
                }
            }

            Spacer(modifier = Modifier.padding(top = 8.dp))

            Text(
                text = "Margin: ${staff.marginScheme} (${staff.marginValue})",
                style = MaterialTheme.typography.bodySmall
            )
        }
    }
}