package com.yonotech.ppob.mobile.ui.transaction

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.CheckCircle
import androidx.compose.material.icons.filled.Error
import androidx.compose.material.icons.filled.HourglassEmpty
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.hilt.navigation.compose.hiltViewModel
import com.yonotech.ppob.mobile.utils.Resource
import com.yonotech.ppob.mobile.viewmodels.transaction.TransactionViewModel

@Composable
fun TransactionResultScreen(
    txId: String,
    onFinish: () -> Unit,
    viewModel: TransactionViewModel = hiltViewModel()
) {
    val transactionState by viewModel.transactionState.collectAsState()

    // Status mapping
    val status = when (val state = transactionState) {
        is Resource.Success -> state.data.status
        else -> "pending"
    }
    
    val transactionData = (transactionState as? Resource.Success)?.data

    LaunchedEffect(txId) {
        viewModel.fetchTransactionStatus(txId)
        viewModel.startPolling(txId)
    }

    val statusConfig = when (status.lowercase()) {
        "success" -> Triple(Icons.Default.CheckCircle, Color(0xFF4CAF50), "Transaksi Berhasil")
        "failed" -> Triple(Icons.Default.Error, Color(0xFFF44336), "Transaksi Gagal")
        else -> Triple(Icons.Default.HourglassEmpty, Color(0xFFFF9800), "Sedang Diproses")
    }

    Column(
        modifier = Modifier
            .fillMaxSize()
            .padding(24.dp)
            .verticalScroll(rememberScrollState()),
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        Spacer(modifier = Modifier.height(32.dp))
        
        Icon(
            imageVector = statusConfig.first,
            contentDescription = null,
            modifier = Modifier.size(100.dp),
            tint = statusConfig.second
        )
        
        Spacer(modifier = Modifier.height(16.dp))
        
        Text(
            text = statusConfig.third,
            fontSize = 24.sp,
            fontWeight = FontWeight.Bold
        )

        if (transactionData?.message != null && transactionData.message.isNotEmpty()) {
            Text(
                text = transactionData.message,
                fontSize = 14.sp,
                textAlign = androidx.compose.ui.text.style.TextAlign.Center,
                modifier = Modifier.padding(horizontal = 16.dp).padding(top = 8.dp),
                color = Color.Gray
            )
        }
        
        Spacer(modifier = Modifier.height(32.dp))

        // Transaction Details Card
        Card(
            modifier = Modifier.fillMaxWidth(),
            colors = CardDefaults.cardColors(containerColor = MaterialTheme.colorScheme.surfaceVariant.copy(alpha = 0.5f))
        ) {
            Column(modifier = Modifier.padding(16.dp)) {
                Text(
                    text = "Detail Transaksi",
                    fontWeight = FontWeight.Bold,
                    fontSize = 16.sp,
                    modifier = Modifier.padding(bottom = 12.dp)
                )
                
                DetailRow("ID Transaksi", txId)
                DetailRow("Waktu", transactionData?.createdAt?.replace("T", " ")?.take(19) ?: "-")
                
                HorizontalDivider(modifier = Modifier.padding(vertical = 12.dp))
                
                DetailRow("Produk", transactionData?.productCode ?: "-")
                DetailRow("Nomor Tujuan", transactionData?.customerNumber ?: "-")
                
                if (transactionData?.amount != null && transactionData.amount > 0) {
                    DetailRow("Nominal", "Rp ${transactionData.amount}")
                }
                
                DetailRow("Total Bayar", "Rp ${transactionData?.price ?: 0}")
            }
        }

        Spacer(modifier = Modifier.height(48.dp))

        Button(
            onClick = onFinish,
            modifier = Modifier.fillMaxWidth(),
            shape = MaterialTheme.shapes.medium
        ) {
            Text("Selesai")
        }
    }
}

@Composable
fun DetailRow(label: String, value: String) {
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .padding(vertical = 4.dp),
        horizontalArrangement = Arrangement.SpaceBetween
    ) {
        Text(text = label, color = MaterialTheme.colorScheme.onSurfaceVariant, fontSize = 14.sp)
        Text(
            text = value, 
            fontWeight = FontWeight.Medium, 
            fontSize = 14.sp,
            color = MaterialTheme.colorScheme.onSurface
        )
    }
}
