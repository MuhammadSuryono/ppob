package com.yonotech.ppob.presentation.transaction

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.CheckCircle
import androidx.compose.material.icons.filled.Error
import androidx.compose.material.icons.filled.HourglassEmpty
import androidx.compose.material3.Button
import androidx.compose.material3.Card
import androidx.compose.material3.CardDefaults
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.Icon
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
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import com.yonotech.ppob.domain.model.Transaction

@Composable
fun TransactionResultScreen(
    transactionId: String,
    onDone: () -> Unit,
    onViewDetail: () -> Unit
) {
    // In a real implementation, this would fetch the transaction status
    // For now, we simulate based on the transactionId

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Hasil Transaksi") },
                navigationIcon = {
                    IconButton(onClick = onDone) {
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
                .padding(24.dp),
            horizontalAlignment = Alignment.CenterHorizontally,
            verticalArrangement = Arrangement.Center
        ) {
            // Result Icon
            Icon(
                imageVector = Icons.Filled.CheckCircle,
                contentDescription = "Berhasil",
                modifier = Modifier.size(80.dp),
                tint = MaterialTheme.colorScheme.primary
            )

            Spacer(modifier = Modifier.size(24.dp))

            Text(
                text = "Transaksi Berhasil!",
                style = MaterialTheme.typography.headlineMedium,
                fontWeight = FontWeight.Bold
            )

            Spacer(modifier = Modifier.size(8.dp))

            Text(
                text = "Rp 50.000",
                style = MaterialTheme.typography.displayLarge,
                fontWeight = FontWeight.Bold,
                color = MaterialTheme.colorScheme.primary
            )

            Spacer(modifier = Modifier.size(4.dp))

            Text(
                text = "Ref: TXN-${System.currentTimeMillis()}",
                style = MaterialTheme.typography.bodyMedium,
                color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.6f)
            )

            Spacer(modifier = Modifier.size(24.dp))

            Card(
                modifier = Modifier.fillMaxWidth(),
                shape = RoundedCornerShape(12.dp),
                colors = CardDefaults.cardColors(
                    containerColor = MaterialTheme.colorScheme.surfaceVariant
                )
            ) {
                Column(modifier = Modifier.padding(16.dp)) {
                    DetailRow("Produk", "Pulsa XL 10GB")
                    DetailRow("Nomor", "0812-XXXX-5678")
                    DetailRow("Harga Admin", "Rp 50.000")
                    DetailRow("Margin", "Rp 3.000")
                    HorizontalDivider(modifier = Modifier.padding(vertical = 8.dp))
                    DetailRow("Total", "Rp 50.000", bold = true)
                }
            }

            Spacer(modifier = Modifier.size(24.dp))

            Button(
                onClick = onDone,
                modifier = Modifier.fillMaxWidth()
            ) {
                Text("Selesai")
            }

            Spacer(modifier = Modifier.size(8.dp))

            Button(
                onClick = onViewDetail,
                modifier = Modifier.fillMaxWidth(),
                colors = ButtonDefaults.buttonColors(
                    containerColor = MaterialTheme.colorScheme.secondaryContainer
                )
            ) {
                Text("Lihat Detail")
            }
        }
    }
}

@Composable
fun DetailRow(label: String, value: String, bold: Boolean = false) {
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .padding(vertical = 4.dp),
        horizontalArrangement = Arrangement.SpaceBetween
    ) {
        Text(
            text = label,
            style = MaterialTheme.typography.bodyMedium,
            color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.6f)
        )
        Text(
            text = value,
            style = if (bold) MaterialTheme.typography.titleMedium else MaterialTheme.typography.bodyMedium,
            fontWeight = if (bold) FontWeight.Bold else FontWeight.Normal
        )
    }
}

fun formatRupiah(amount: Double): String {
    return "Rp %,d".format(amount.toInt())
}