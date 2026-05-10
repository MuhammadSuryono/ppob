package com.yonotech.ppob.presentation.home

import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.lazy.LazyRow
import androidx.compose.foundation.lazy.grid.GridCells
import androidx.compose.foundation.lazy.grid.LazyVerticalGrid
import androidx.compose.foundation.lazy.grid.items
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.AccountCircle
import androidx.compose.material.icons.filled.ArrowForward
import androidx.compose.material.icons.filled.Notifications
import androidx.compose.material3.Button
import androidx.compose.material3.ButtonDefaults
import androidx.compose.material3.Card
import androidx.compose.material3.CardDefaults
import androidx.compose.material3.HorizontalDivider
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
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Brush
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import coil.compose.AsyncImage
import com.yonotech.ppob.domain.model.Category

@Composable
fun HomeScreen(
    viewModel: HomeViewModel = androidx.lifecycle.viewmodel.compose.viewModel(),
    onNavigateToCategory: () -> Unit,
    onNavigateToTransactions: () -> Unit,
    onNavigateToWallet: () -> Unit,
    onNavigateToProfile: () -> Unit
) {
    val uiState by viewModel.uiState.collectAsState()

    LaunchedEffect(Unit) {
        viewModel.refreshBalance()
    }

    Scaffold(
        topBar = {
            TopAppBar(
                 title = {
                     Text(
                         "Halo, ${uiState.categories.getOrNull(0)?.categoryName ?: "User"}",
                         style = MaterialTheme.typography.titleLarge
                     )
                 },
                actions = {
                    IconButton(onClick = { /* Notifications */ }) {
                        Icon(Icons.Filled.Notifications, contentDescription = "Notifikasi")
                    }
                }
            )
        }
    ) { paddingValues ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues)
        ) {
            // Balance Card
            BalanceCard(
                balance = uiState.balance,
                balanceHeld = uiState.balanceHeld,
                onClick = onNavigateToWallet
            )

            // Quick Actions Grid
            HorizontalSection(title = "Layanan Cepat", showMore = true, onMoreClick = onNavigateToCategory)

             LazyVerticalGrid(
                 columns = GridCells.Fixed(4),
                 contentPadding = padding(16.dp)
             ) {
                items(uiState.categories.take(8)) { category ->
                    CategoryQuickCard(
                        category = category,
                        onClick = { onNavigateToCategory() }
                    )
                }
            }

            // Recent Transactions
            HorizontalSection(
                title = "Transaksi Terakhir",
                showMore = true,
                onMoreClick = onNavigateToTransactions
            )

            if (uiState.recentTransactions.isEmpty()) {
                Text(
                    text = "Belum ada transaksi",
                    style = MaterialTheme.typography.bodyMedium,
                    modifier = Modifier.padding(horizontal = 16.dp, vertical = 8.dp)
                )
            } else {
                LazyRow(
                    contentPadding = androidx.compose.foundation.layout.padding(horizontal = 16.dp),
                    horizontalArrangement = Arrangement.spacedBy(12.dp)
                ) {
                    items(uiState.recentTransactions.take(5)) { transaction ->
                        RecentTransactionCard(transaction = transaction)
                    }
                }
            }
        }
    }
}

@Composable
fun BalanceCard(
    balance: Double,
    balanceHeld: Double,
    onClick: () -> Unit
) {
    Card(
        modifier = Modifier
            .fillMaxWidth()
            .padding(16.dp)
            .height(160.dp),
        shape = RoundedCornerShape(16.dp),
        colors = CardDefaults.cardColors(
            containerColor = Color(0xFF4CAF50)
        ),
        elevation = CardDefaults.cardElevation(defaultElevation = 4.dp)
    ) {
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(20.dp),
            verticalArrangement = Arrangement.SpaceBetween
        ) {
            Text(
                text = "Saldo Tersedia",
                color = Color.White.copy(alpha = 0.8f),
                style = MaterialTheme.typography.bodyMedium
            )

            Text(
                text = formatRupiah(balance),
                style = MaterialTheme.typography.displayLarge.copy(
                    fontWeight = FontWeight.Bold,
                    color = Color.White
                )
            )

            if (balanceHeld > 0) {
                Row(verticalAlignment = Alignment.CenterVertically) {
                    Text(
                        text = "+ ${formatRupiah(balanceHeld)} dikunci",
                        color = Color.White.copy(alpha = 0.7f),
                        style = MaterialTheme.typography.bodySmall
                    )
                }
            }

            Button(
                onClick = onClick,
                colors = ButtonDefaults.buttonColors(containerColor = Color.White.copy(alpha = 0.2f)),
                shape = RoundedCornerShape(8.dp)
            ) {
                Text("Top Up", color = Color.White)
            }
        }
    }
}

@Composable
fun CategoryQuickCard(
    category: Category,
    onClick: () -> Unit
) {
    Column(
        modifier = Modifier
            .padding(8.dp)
            .size(64.dp),
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.Center
    ) {
        Card(
            shape = RoundedCornerShape(12.dp),
            modifier = Modifier
                .size(48.dp)
                .clip(RoundedCornerShape(12.dp))
        ) {
            Box(
                modifier = Modifier
                    .fillMaxSize()
                    .background(MaterialTheme.colorScheme.primaryContainer),
                contentAlignment = Alignment.Center
            ) {
                if (!category.iconUrl.isNullOrEmpty()) {
                    AsyncImage(
                        model = category.iconUrl,
                        contentDescription = category.categoryName,
                        modifier = Modifier.size(28.dp)
                    )
                } else {
                    Text(
                        text = category.categoryName.take(1),
                        style = MaterialTheme.typography.bodySmall
                    )
                }
            }
        }

        Spacer(modifier = Modifier.height(4.dp))

        Text(
            text = category.categoryName,
            style = MaterialTheme.typography.labelSmall,
            maxLines = 1
        )
    }
}

@Composable
fun RecentTransactionCard(transaction: com.yonotech.ppob.domain.model.Transaction) {
    Card(
        modifier = Modifier.width(160.dp),
        shape = RoundedCornerShape(12.dp),
        elevation = CardDefaults.cardElevation(defaultElevation = 2.dp)
    ) {
        Column(
            modifier = Modifier.padding(12.dp),
            verticalArrangement = Arrangement.spacedBy(4.dp)
        ) {
            Text(
                text = transaction.productName,
                style = MaterialTheme.typography.bodySmall,
                fontWeight = FontWeight.Medium,
                maxLines = 1
            )

            Text(
                text = formatRupiah(transaction.sellingPrice),
                style = MaterialTheme.typography.bodyLarge,
                fontWeight = FontWeight.Bold
            )

            Text(
                text = transaction.status,
                style = MaterialTheme.typography.labelSmall,
                color = when (transaction.status) {
                    "Success" -> Color(0xFF4CAF50)
                    "Pending" -> Color(0xFFFFA726)
                    else -> Color(0xFFF44336)
                }
            )
        }
    }
}

@Composable
fun HorizontalSection(
    title: String,
    showMore: Boolean = false,
    onMoreClick: () -> Unit = {}
) {
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .padding(horizontal = 16.dp, vertical = 8.dp),
        horizontalArrangement = Arrangement.SpaceBetween,
        verticalAlignment = Alignment.CenterVertically
    ) {
        Text(
            text = title,
            style = MaterialTheme.typography.titleMedium,
            fontWeight = FontWeight.SemiBold
        )

        if (showMore) {
            Text(
                text = "Lihat semua",
                style = MaterialTheme.typography.labelMedium,
                color = MaterialTheme.colorScheme.primary
                // TODO: Add click handler
            )
        }
    }
}

fun formatRupiah(amount: Double): String {
    return "Rp %,d".format(amount.toInt())
}