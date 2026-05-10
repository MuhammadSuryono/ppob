package com.yonotech.ppob.presentation.transaction

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
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
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.unit.dp
import com.yonotech.ppob.data.remote.model.Product
import com.yonotech.ppob.domain.model.Product as DomainProduct

@Composable
fun TransactionConfirmationScreen(
    viewModel: com.yonotech.ppob.presentation.transaction.TransactionInitViewModel = androidx.lifecycle.viewmodel.compose.viewModel(),
    productId: String,
    onConfirm: (String) -> Unit,
    onBack: () -> Unit
) {
    val txState by viewModel.transactionState.collectAsState()
    val productUiState by viewModel.uiState.collectAsState()

    // Find the product details
    val product = productUiState.products.find { it.productId == productId }
        ?: DomainProduct(
            productId = productId,
            buyerSkuCode = "",
            productName = "Produk",
            basePrice = 0.0,
            platformPrice = 0.0,
            isPrepaid = true,
            categoryId = ""
        )

    var sellingPrice by remember { mutableStateOf(product.mitraSellingPrice ?: product.platformPrice) }
    var customerNumber by remember { mutableStateOf("") }

    LaunchedEffect(txState) {
        if (txState is com.yonotech.ppob.presentation.transaction.TransactionState.Success ||
            txState is com.yonotech.ppob.presentation.transaction.TransactionState.Pending
        ) {
            val txId = when (val state = txState) {
                is com.yonotech.ppob.presentation.transaction.TransactionState.Success -> state.transactionId
                is com.yonotech.ppob.presentation.transaction.TransactionState.Pending -> state.transactionId
                else -> ""
            }
            onConfirm(txId)
        }
    }

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Konfirmasi Transaksi") },
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
                .fillMaxWidth()
                .padding(paddingValues)
                .padding(24.dp)
                .verticalScroll(rememberScrollState()),
            verticalArrangement = Arrangement.spacedBy(16.dp)
        ) {
            Text(
                text = product.productName,
                style = MaterialTheme.typography.headlineMedium,
                fontWeight = androidx.compose.ui.text.font.FontWeight.Bold
            )

            Text(
                text = product.buyerSkuCode,
                style = MaterialTheme.typography.bodyMedium,
                color = MaterialTheme.colorScheme.onSurface.copy(alpha = 0.6f)
            )

            Spacer(modifier = Modifier.padding(8.dp))

            // Price Breakdown
            Card(
                modifier = Modifier.fillMaxWidth(),
                shape = androidx.compose.foundation.shape.RoundedCornerShape(12.dp)
            ) {
                Column(modifier = Modifier.padding(16.dp)) {
                    Row(
                        modifier = Modifier.fillMaxWidth(),
                        horizontalArrangement = Arrangement.SpaceBetween
                    ) {
                        Text(
                            text = "Harga Platform",
                            style = MaterialTheme.typography.bodyMedium
                        )
                        Text(
                            text = formatRupiah(product.platformPrice),
                            style = MaterialTheme.typography.bodyMedium
                        )
                    }

                    Spacer(modifier = Modifier.padding(8.dp))

                    Row(
                        modifier = Modifier.fillMaxWidth(),
                        horizontalArrangement = Arrangement.SpaceBetween
                    ) {
                        Text(
                            text = "Margin",
                            style = MaterialTheme.typography.bodyMedium
                        )
                        Text(
                            text = formatRupiah(sellingPrice - product.platformPrice),
                            style = MaterialTheme.typography.bodyMedium,
                            color = MaterialTheme.colorScheme.primary
                        )
                    }

                    HorizontalDivider(modifier = Modifier.padding(vertical = 8.dp))

                    Row(
                        modifier = Modifier.fillMaxWidth(),
                        horizontalArrangement = Arrangement.SpaceBetween
                    ) {
                        Text(
                            text = "Total Pembayaran",
                            style = MaterialTheme.typography.titleMedium,
                            fontWeight = androidx.compose.ui.text.font.FontWeight.Bold
                        )
                        Text(
                            text = formatRupiah(sellingPrice),
                            style = MaterialTheme.typography.titleMedium,
                            fontWeight = androidx.compose.ui.text.font.FontWeight.Bold,
                            color = MaterialTheme.colorScheme.primary
                        )
                    }
                }
            }

            OutlinedTextField(
                value = customerNumber,
                onValueChange = { customerNumber = it },
                label = { Text("Nomor Pelanggan") },
                placeholder = { Text("Masukkan nomor pelanggan") },
                keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Phone),
                modifier = Modifier.fillMaxWidth()
            )

            when (txState) {
                is com.yonotech.ppob.presentation.transaction.TransactionState.Error -> {
                    Text(
                        text = (txState as com.yonotech.ppob.presentation.transaction.TransactionState.Error).message,
                        color = MaterialTheme.colorScheme.error,
                        style = MaterialTheme.typography.bodyMedium
                    )
                }
                is com.yonotech.ppob.presentation.transaction.TransactionState.InsufficientBalance -> {
                    Text(
                        text = "Saldo tidak mencukupi",
                        color = MaterialTheme.colorScheme.error,
                        style = MaterialTheme.typography.bodyMedium
                    )
                }
                else -> {}
            }

            Button(
                onClick = {
                    viewModel.initiateTransaction(
                        pin = "", // PIN should be entered separately
                        authToken = "", // TODO: Get from DataStore
                        customerNo = customerNumber,
                        productId = productId,
                        sellingPrice = sellingPrice
                    )
                },
                enabled = customerNumber.isNotEmpty() && txState is com.yonotech.ppob.presentation.transaction.TransactionState.Idle,
                modifier = Modifier.fillMaxWidth()
            ) {
                Text(
                    if (txState is com.yonotech.ppob.presentation.transaction.TransactionState.Loading)
                        "Memproses..."
                    else "Lanjutkan ke Pembayaran"
                )
            }
        }
    }
}

fun formatRupiah(amount: Double): String {
    return "Rp %,d".format(amount.toInt())
}