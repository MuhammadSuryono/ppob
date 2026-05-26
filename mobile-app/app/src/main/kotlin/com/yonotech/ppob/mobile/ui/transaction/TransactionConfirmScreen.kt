package com.yonotech.ppob.mobile.ui.transaction

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.grid.GridCells
import androidx.compose.foundation.lazy.grid.LazyVerticalGrid
import androidx.compose.foundation.lazy.grid.items
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.ArrowBack
import androidx.compose.material.icons.filled.Backspace
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.yonotech.ppob.mobile.ui.components.PinDots
import com.yonotech.ppob.mobile.ui.components.PinPad
import com.yonotech.ppob.mobile.ui.components.PpoButton
import com.yonotech.ppob.mobile.utils.Resource
import com.yonotech.ppob.mobile.viewmodels.transaction.TransactionViewModel
import java.text.NumberFormat
import java.util.*

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun TransactionConfirmScreen(
    onSuccess: (String) -> Unit,
    onBack: () -> Unit,
    viewModel: TransactionViewModel
) {
    var pin by remember { mutableStateOf("") }
    val transactionState by viewModel.transactionState.collectAsState()

    LaunchedEffect(transactionState) {
        if (transactionState is Resource.Success) {
            val data = (transactionState as Resource.Success).data
            onSuccess(data.id ?: data.transactionId)
            viewModel.resetState()
        }
    }

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Konfirmasi & PIN") },
                navigationIcon = {
                    IconButton(onClick = onBack) {
                        Icon(Icons.Default.ArrowBack, contentDescription = "Back")
                    }
                },
                colors = TopAppBarDefaults.topAppBarColors(
                    containerColor = MaterialTheme.colorScheme.primary,
                    titleContentColor = MaterialTheme.colorScheme.onPrimary,
                    navigationIconContentColor = MaterialTheme.colorScheme.onPrimary
                )
            )
        }
    ) { padding ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(padding)
                .padding(24.dp),
            horizontalAlignment = Alignment.CenterHorizontally
        ) {
            Card(
                modifier = Modifier.fillMaxWidth(),
                colors = CardDefaults.cardColors(containerColor = MaterialTheme.colorScheme.secondaryContainer)
            ) {
                Column(modifier = Modifier.padding(16.dp)) {
                    Row(modifier = Modifier.fillMaxWidth(), horizontalArrangement = Arrangement.SpaceBetween) {
                        Text("Pelanggan:")
                        Text(viewModel.customerNo, fontWeight = FontWeight.Bold)
                    }
                }
            }

            Spacer(modifier = Modifier.height(32.dp))

            Text("Masukkan 6-Digit PIN", fontSize = 16.sp, fontWeight = FontWeight.Medium)

            Spacer(modifier = Modifier.height(16.dp))

            PinDots(pin = pin)

            Spacer(modifier = Modifier.height(32.dp))

            if (transactionState is Resource.Error) {
                Text(text = (transactionState as Resource.Error).message, color = MaterialTheme.colorScheme.error)
            }

            Spacer(modifier = Modifier.weight(1f))

            PinPad(
                onDigit = { digit -> if (pin.length < 6) pin += digit },
                onBackspace = { if (pin.isNotEmpty()) pin = pin.dropLast(1) }
            )

            Spacer(modifier = Modifier.height(24.dp))

            PpoButton(
                label = "Konfirmasi Transaksi",
                onClick = { viewModel.initiateTransaction(pin) },
                isLoading = transactionState is Resource.Loading,
                enabled = pin.length == 6
            )
        }
    }
}