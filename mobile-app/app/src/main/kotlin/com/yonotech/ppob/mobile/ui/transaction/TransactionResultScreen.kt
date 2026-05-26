package com.yonotech.ppob.mobile.ui.transaction

import androidx.compose.foundation.layout.*
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

    var status by remember { mutableStateOf("initiated") }
    var message by remember { mutableStateOf("") }

    LaunchedEffect(transactionState) {
        if (transactionState is Resource.Success) {
            val data = (transactionState as Resource.Success).data
            status = data.status
            message = data.message ?: ""
        }
    }

    val statusConfig = when (status.lowercase()) {
        "success" -> Triple(Icons.Default.CheckCircle, Color(0xFF4CAF50), "Transaksi Berhasil")
        "failed" -> Triple(Icons.Default.Error, Color(0xFFF44336), "Transaksi Gagal")
        else -> Triple(Icons.Default.HourglassEmpty, Color(0xFFFF9800), "Sedang Diproses")
    }

    Column(
        modifier = Modifier
            .fillMaxSize()
            .padding(24.dp),
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.Center
    ) {
        Icon(
            imageVector = statusConfig.first,
            contentDescription = null,
            modifier = Modifier.size(120.dp),
            tint = statusConfig.second
        )
        
        Spacer(modifier = Modifier.height(24.dp))
        
        Text(
            text = statusConfig.third,
            fontSize = 24.sp,
            fontWeight = FontWeight.Bold
        )

        if (message.isNotEmpty()) {
            Text(
                text = message,
                fontSize = 14.sp,
                textAlign = androidx.compose.ui.text.style.TextAlign.Center,
                modifier = Modifier.padding(horizontal = 16.dp).padding(top = 8.dp),
                color = Color.Gray
            )
        }
        
        Text(
            text = "ID Transaksi: $txId",
            fontSize = 14.sp,
            color = MaterialTheme.colorScheme.onSurfaceVariant,
            modifier = Modifier.padding(top = 16.dp)
        )

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
