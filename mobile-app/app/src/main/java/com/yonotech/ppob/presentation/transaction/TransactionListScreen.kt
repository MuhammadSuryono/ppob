package com.yonotech.ppob.presentation.transaction

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.padding
import androidx.compose.material3.CircularProgressIndicator
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Scaffold
import androidx.compose.material3.Tab
import androidx.compose.material3.TabRow
import androidx.compose.material3.Text
import androidx.compose.material3.TopAppBar
import androidx.compose.runtime.Composable
import androidx.compose.runtime.collectAsState
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.remember
import androidx.compose.ui.Modifier
import com.yonotech.ppob.domain.model.Transaction
import com.yonotech.ppob.presentation.components.TransactionList

@Composable
fun TransactionListScreen(
    viewModel: com.yonotech.ppob.presentation.transactiondetail.TransactionHistoryViewModel = androidx.lifecycle.viewmodel.compose.viewModel(),
    onTransactionClick: (String) -> Unit,
    onFilterClicked: () -> Unit
) {
    val tabs = listOf("Semua", "Berhasil", "Pending", "Gagal")
    var selectedTabIndex by remember { mutableStateOf(0) }
    val uiState by viewModel.uiState.collectAsState()

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Riwayat Transaksi") },
                actions = {
                    // Filter button
                }
            )
        }
    ) { paddingValues ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues)
        ) {
            TabRow(selectedTabIndex = selectedTabIndex) {
                tabs.forEachIndexed { index, title ->
                    Tab(
                        selected = selectedTabIndex == index,
                        onClick = { selectedTabIndex = index },
                        text = { Text(title) }
                    )
                }
            }

            if (uiState.isLoading) {
                Box(
                    modifier = Modifier.fillMaxSize(),
                    contentAlignment = androidx.compose.ui.Alignment.Center
                ) {
                    CircularProgressIndicator()
                }
            } else {
                // Mock data for demonstration
                val mockTransactions = listOf(
                    Transaction(
                        id = "1",
                        refId = "REF001",
                        productName = "Pulsa XL 10GB",
                        customerNumber = "081234567890",
                        status = "Success",
                        sellingPrice = 50000.0,
                        platformPrice = 47000.0,
                        marginAmount = 3000.0,
                        commissionAmount = 1800.0,
                        createdAt = System.currentTimeMillis(),
                        completedAt = System.currentTimeMillis()
                    ),
                    Transaction(
                        id = "2",
                        refId = "REF002",
                        productName = "PLN Token 50kWh",
                        customerNumber = "1234567890",
                        status = "Pending",
                        sellingPrice = 55000.0,
                        platformPrice = 52000.0,
                        marginAmount = 3000.0,
                        commissionAmount = 1800.0,
                        createdAt = System.currentTimeMillis() - 3600000,
                        completedAt = null
                    ),
                    Transaction(
                        id = "3",
                        refId = "REF003",
                        productName = "Pulsa Telkomsel 5GB",
                        customerNumber = "081298765432",
                        status = "Success",
                        sellingPrice = 30000.0,
                        platformPrice = 29000.0,
                        marginAmount = 1000.0,
                        commissionAmount = 600.0,
                        createdAt = System.currentTimeMillis() - 7200000,
                        completedAt = System.currentTimeMillis() - 7200000
                    )
                )

                val filtered = when (selectedTabIndex) {
                    0 -> mockTransactions
                    1 -> mockTransactions.filter { it.status == "Success" }
                    2 -> mockTransactions.filter { it.status == "Pending" }
                    3 -> mockTransactions.filter { it.status == "Failed" }
                    else -> mockTransactions
                }

                TransactionList(
                    transactions = filtered,
                    onTransactionClick = onTransactionClick,
                    modifier = Modifier.fillMaxSize()
                )
            }
        }
    }
}

// Default viewModel stub
class TransactionHistoryViewModel {
    val uiState = androidx.compose.runtime.mutableStateOf(
        com.yonotech.ppob.presentation.transactiondetail.TransactionDetailUiState(isLoading = false)
    )
}