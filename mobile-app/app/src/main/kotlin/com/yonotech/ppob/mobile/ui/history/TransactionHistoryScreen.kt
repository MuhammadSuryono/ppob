package com.yonotech.ppob.mobile.ui.history

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.ArrowBack
import androidx.compose.material.icons.filled.CheckCircle
import androidx.compose.material.icons.filled.Error
import androidx.compose.material.icons.filled.HourglassEmpty
import androidx.compose.material.icons.filled.Help
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.hilt.navigation.compose.hiltViewModel
import com.yonotech.ppob.mobile.data.remote.dto.TransactionHistoryResponse
import com.yonotech.ppob.mobile.ui.theme.*
import com.yonotech.ppob.mobile.utils.Resource
import com.yonotech.ppob.mobile.viewmodels.staff.StaffViewModel
import java.text.NumberFormat
import java.util.Locale

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun TransactionHistoryScreen(
    onBackClick: () -> Unit,
    viewModel: StaffViewModel = hiltViewModel()
) {
    val historyState by viewModel.transactionHistoryState.collectAsState()
    val currencyFormat = NumberFormat.getCurrencyInstance(Locale("in", "ID")).apply { maximumFractionDigits = 0 }

    var selectedStatus by remember { mutableStateOf("All") }
    val statuses = listOf("All", "SUCCESS", "PENDING", "FAILED")

    LaunchedEffect(Unit) {
        viewModel.getTransactionHistory()
    }

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Riwayat Transaksi") },
                navigationIcon = {
                    IconButton(onClick = onBackClick) {
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
        Column(modifier = Modifier.fillMaxSize().padding(padding)) {
            ScrollableTabRow(selectedTabIndex = statuses.indexOf(selectedStatus)) {
                statuses.forEach { status ->
                    Tab(
                        selected = selectedStatus == status,
                        onClick = { selectedStatus = status },
                        text = { Text(status) }
                    )
                }
            }

            when (historyState) {
                is Resource.Loading -> {
                    Box(modifier = Modifier.fillMaxSize(), contentAlignment = Alignment.Center) {
                        CircularProgressIndicator()
                    }
                }
                is Resource.Success -> {
                    val transactions = (historyState as Resource.Success).data
                    if (transactions.isEmpty()) {
                        Box(modifier = Modifier.fillMaxSize(), contentAlignment = Alignment.Center) {
                            Text("Tidak ada transaksi")
                        }
                    } else {
                        LazyColumn(
                            contentPadding = PaddingValues(16.dp),
                            verticalArrangement = Arrangement.spacedBy(8.dp)
                        ) {
                            items(transactions) { transaction ->
                                TransactionHistoryItem(transaction, currencyFormat)
                            }
                        }
                    }
                }
                is Resource.Error -> {
                    Box(modifier = Modifier.fillMaxSize(), contentAlignment = Alignment.Center) {
                        Text(
                            text = (historyState as Resource.Error).message,
                            color = MaterialTheme.colorScheme.error
                        )
                    }
                }
                else -> {}
            }
        }
    }
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun TransactionHistoryItem(
    transaction: TransactionHistoryResponse,
    currencyFormat: NumberFormat
) {
    Card(
        modifier = Modifier.fillMaxWidth(),
        shape = MaterialTheme.shapes.medium,
        colors = CardDefaults.cardColors(containerColor = MaterialTheme.colorScheme.surface)
    ) {
        Row(
            modifier = Modifier.padding(16.dp),
            verticalAlignment = Alignment.CenterVertically
        ) {
            Surface(
                shape = MaterialTheme.shapes.small,
                color = getStatusColor(transaction.status).copy(alpha = 0.2f),
                modifier = Modifier.size(40.dp)
            ) {
                Box(contentAlignment = Alignment.Center) {
                    Icon(
                        imageVector = getStatusIcon(transaction.status),
                        contentDescription = null,
                        tint = getStatusColor(transaction.status)
                    )
                }
            }

            Spacer(modifier = Modifier.width(12.dp))

            Column(modifier = Modifier.weight(1f)) {
                Text(text = transaction.productName, fontWeight = FontWeight.Bold, maxLines = 1)
                Text(
                    text = transaction.customerNumber,
                    fontSize = 12.sp,
                    color = MaterialTheme.colorScheme.onSurfaceVariant
                )
            }

            Column(horizontalAlignment = Alignment.End) {
                Text(
                    text = currencyFormat.format(transaction.sellingPrice),
                    fontWeight = FontWeight.Bold,
                    color = MaterialTheme.colorScheme.primary
                )
                StatusChip(status = transaction.status)
            }
        }
    }
}

@Composable
fun StatusChip(status: String) {
    val color = getStatusColor(status)
    Surface(
        shape = MaterialTheme.shapes.small,
        color = color.copy(alpha = 0.1f),
        modifier = Modifier.padding(top = 4.dp)
    ) {
        Text(
            text = status,
            color = color,
            fontSize = 10.sp,
            fontWeight = FontWeight.Medium,
            modifier = Modifier.padding(horizontal = 8.dp, vertical = 2.dp)
        )
    }
}

fun getStatusColor(status: String): Color {
    return when (status.uppercase()) {
        "SUCCESS" -> SuccessGreen
        "PENDING" -> WarningOrange
        "FAILED" -> ErrorRed
        else -> TextLight
    }
}

fun getStatusIcon(status: String) = when (status.uppercase()) {
    "SUCCESS" -> Icons.Default.CheckCircle
    "PENDING" -> Icons.Default.HourglassEmpty
    "FAILED" -> Icons.Default.Error
    else -> Icons.Default.Help
}