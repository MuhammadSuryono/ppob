package com.yonotech.ppob.mobile.ui.home

import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.verticalScroll
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.AccountBalanceWallet
import androidx.compose.material.icons.filled.Notifications
import androidx.compose.material.icons.filled.Search
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Brush
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.hilt.navigation.compose.hiltViewModel
import com.yonotech.ppob.mobile.data.remote.dto.CategoryCollection
import com.yonotech.ppob.mobile.data.remote.dto.CategoryDto
import com.yonotech.ppob.mobile.data.remote.dto.WalletResponse
import com.yonotech.ppob.mobile.ui.theme.PpoMobileTheme
import com.yonotech.ppob.mobile.utils.Resource
import com.yonotech.ppob.mobile.viewmodels.product.ProductViewModel
import com.yonotech.ppob.mobile.viewmodels.wallet.WalletViewModel
import java.text.NumberFormat
import java.util.Locale

@Composable
fun HomeScreen(
    onCategoryClick: (CategoryDto) -> Unit,
    onWalletClick: () -> Unit,
    productViewModel: ProductViewModel = hiltViewModel(),
    walletViewModel: WalletViewModel = hiltViewModel()
) {
    val categoriesState by productViewModel.categories.collectAsState()
    val walletState by walletViewModel.balanceState.collectAsState()

    LaunchedEffect(Unit) {
        productViewModel.getCategories()
        walletViewModel.getBalance()
    }

    HomeContent(
        categoriesState = categoriesState,
        walletState = walletState,
        onCategoryClick = onCategoryClick,
        onWalletClick = onWalletClick
    )
}

@Composable
fun HomeContent(
    categoriesState: Resource<CategoryCollection>,
    walletState: Resource<WalletResponse>,
    onCategoryClick: (CategoryDto) -> Unit,
    onWalletClick: () -> Unit
) {
    Scaffold(
        topBar = {
            TopNavigationBar()
        },
        containerColor = Color(0xFFF8F9FA)
    ) { padding ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(padding)
                .verticalScroll(rememberScrollState())
        ) {
            Spacer(modifier = Modifier.height(24.dp))
            
            WalletCard(walletState = walletState, onTopUpClick = onWalletClick)
            
            Spacer(modifier = Modifier.height(20.dp))

            Surface(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(horizontal = 16.dp),
                shape = RoundedCornerShape(16.dp),
                color = Color.White,
                shadowElevation = 1.dp
            ) {
                Column(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(vertical = 24.dp, horizontal = 12.dp)
                ) {
                    when (categoriesState) {
                        is Resource.Loading -> {
                            Box(
                                modifier = Modifier
                                    .fillMaxWidth()
                                    .height(200.dp),
                                contentAlignment = Alignment.Center
                            ) {
                                CircularProgressIndicator(color = MaterialTheme.colorScheme.primary)
                            }
                        }
                        is Resource.Success -> {
                            ServiceGrid(
                                categories = categoriesState.data.categories,
                                onCategoryClick = onCategoryClick
                            )
                        }
                        is Resource.Error -> {
                            Text(
                                text = categoriesState.message,
                                color = MaterialTheme.colorScheme.error,
                                modifier = Modifier.padding(16.dp),
                                textAlign = TextAlign.Center
                            )
                        }
                        else -> {}
                    }
                }
            }
            
            Spacer(modifier = Modifier.height(32.dp))
        }
    }
}

@Composable
fun TopNavigationBar() {
    Surface(
        modifier = Modifier
            .fillMaxWidth()
            .clip(RoundedCornerShape(bottomStart = 28.dp, bottomEnd = 28.dp)),
        color = MaterialTheme.colorScheme.primary,
        shadowElevation = 8.dp
    ) {
        Column(
            modifier = Modifier
                .statusBarsPadding()
                .padding(horizontal = 16.dp, vertical = 20.dp)
        ) {
            Row(
                modifier = Modifier.fillMaxWidth(),
                verticalAlignment = Alignment.CenterVertically
            ) {
                Text(
                    text = "PPOB Mobile",
                    color = Color.White,
                    fontSize = 22.sp,
                    fontWeight = FontWeight.ExtraBold
                )
                Spacer(modifier = Modifier.weight(1f))
                IconButton(onClick = { }) {
                    Icon(
                        Icons.Default.Notifications,
                        contentDescription = "Notifikasi",
                        tint = Color.White
                    )
                }
            }
            
            Spacer(modifier = Modifier.height(16.dp))
            
            Surface(
                modifier = Modifier
                    .fillMaxWidth()
                    .height(48.dp)
                    .clickable { },
                shape = RoundedCornerShape(12.dp),
                color = Color.White
            ) {
                Row(
                    modifier = Modifier
                        .fillMaxSize()
                        .padding(horizontal = 16.dp),
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Icon(
                        Icons.Default.Search,
                        contentDescription = null,
                        tint = Color.Gray,
                        modifier = Modifier.size(20.dp)
                    )
                    Spacer(modifier = Modifier.width(12.dp))
                    Text(
                        text = "Cari pulsa, paket data, atau lainnya...",
                        color = Color.Gray,
                        fontSize = 14.sp
                    )
                }
            }
            Spacer(modifier = Modifier.height(8.dp))
        }
    }
}

@Composable
fun WalletCard(walletState: Resource<WalletResponse>, onTopUpClick: () -> Unit) {
    val currencyFormat = NumberFormat.getCurrencyInstance(Locale("in", "ID"))
    currencyFormat.maximumFractionDigits = 0

    val balanceText = when (walletState) {
        is Resource.Success -> currencyFormat.format(walletState.data.balanceAvailable)
        is Resource.Loading -> "Memuat..."
        is Resource.Error -> "Rp ---"
        else -> "Rp 0"
    }

    Card(
        modifier = Modifier
            .fillMaxWidth()
            .padding(horizontal = 16.dp),
        shape = RoundedCornerShape(12.dp),
        elevation = CardDefaults.cardElevation(defaultElevation = 2.dp)
    ) {
        Box(
            modifier = Modifier
                .fillMaxWidth()
                .background(
                    Brush.horizontalGradient(
                        colors = listOf(
                            MaterialTheme.colorScheme.primary,
                            MaterialTheme.colorScheme.primary.copy(alpha = 0.85f)
                        )
                    )
                )
                .padding(16.dp)
        ) {
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.SpaceBetween,
                verticalAlignment = Alignment.CenterVertically
            ) {
                Row(verticalAlignment = Alignment.CenterVertically) {
                    Surface(
                        color = Color.White.copy(alpha = 0.2f),
                        shape = RoundedCornerShape(10.dp),
                        modifier = Modifier.size(44.dp)
                    ) {
                        Box(contentAlignment = Alignment.Center) {
                            Icon(
                                Icons.Default.AccountBalanceWallet,
                                contentDescription = null,
                                tint = Color.White,
                                modifier = Modifier.size(24.dp)
                            )
                        }
                    }
                    Spacer(modifier = Modifier.width(12.dp))
                    Column {
                        Text(
                            "Saldo PPOB",
                            color = Color.White.copy(alpha = 0.9f),
                            fontSize = 13.sp,
                            fontWeight = FontWeight.Medium
                        )
                        Text(
                            balanceText,
                            color = Color.White,
                            fontSize = 20.sp,
                            fontWeight = FontWeight.ExtraBold
                        )
                    }
                }
                
                Button(
                    onClick = onTopUpClick,
                    colors = ButtonDefaults.buttonColors(
                        containerColor = Color(0xFFFFD54F) 
                    ),
                    shape = RoundedCornerShape(20.dp),
                    contentPadding = PaddingValues(horizontal = 16.dp, vertical = 6.dp),
                    modifier = Modifier.height(36.dp)
                ) {
                    Text(
                        "ISI SALDO",
                        color = Color(0xFF5D4037),
                        fontWeight = FontWeight.Bold,
                        fontSize = 11.sp
                    )
                }
            }
        }
    }
}

@Composable
fun ServiceGrid(
    categories: List<CategoryDto>,
    onCategoryClick: (CategoryDto) -> Unit
) {
    val chunks = categories.chunked(4)
    Column(verticalArrangement = Arrangement.spacedBy(24.dp)) {
        chunks.forEach { rowCategories ->
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.SpaceEvenly
            ) {
                rowCategories.forEach { category ->
                    CategoryItem(
                        category = category,
                        modifier = Modifier.weight(1f),
                        onClick = { onCategoryClick(category) }
                    )
                }
                if (rowCategories.size < 4) {
                    repeat(4 - rowCategories.size) {
                        Spacer(modifier = Modifier.weight(1f))
                    }
                }
            }
        }
    }
}

private fun normalizedCategoryKey(name: String): String {
    // API kadang kirim newline/multi-space di nama kategori, normalisasi untuk pencocokan yang stabil.
    return name
        .lowercase()
        .replace("\n", " ")
        .replace(Regex("\\s+"), " ")
        .trim()
}

@Composable
fun CategoryItem(
    category: CategoryDto,
    modifier: Modifier = Modifier,
    onClick: () -> Unit
) {
    val categoryKey = normalizedCategoryKey(category.name)

    val iconColor = when (categoryKey) {
        "aktivasi_perdana" -> Color(0xFF90CAF9) // Blue
        "aktivasi_voucher" -> Color(0xFF80DEEA) // Cyan
        "data" -> Color(0xFFA5D6A7) // Mint Green
        "e-money" -> Color(0xFF80CBC4) // Teal
        "games" -> Color(0xFFCE93D8) // Pastel Purple
        "gas" -> Color(0xFFFFAB91) // Soft Coral
        "masa_aktif" -> Color(0xFFB0BEC5) // Neutral Gray
        "pln" -> Color(0xFFFFCC80) // Soft Orange (Yellowish)
        "paket_sms_and_telpon" -> Color(0xFFC5E1A5) // Hijau mint terang
        "pulsa" -> Color(0xFF81D4FA) // Biru muda
        "tv" -> Color(0xFFB39DDB) // Lavender
        "voucher" -> Color(0xFFF48FB1) // Soft Pink
        else -> MaterialTheme.colorScheme.primary.copy(alpha = 0.6f)
    }

    val iconText = when (categoryKey) {
        "aktivasi_perdana" -> "📶"
        "aktivasi_voucher" -> "✅"
        "data" -> "📡"
        "e-money" -> "💳"
        "games" -> "🎮"
        "gas" -> "🔥"
        "masa_aktif" -> "⏰"
        "pln" -> "⚡"
        "paket_sms_and_telpon" -> "📞"
        "pulsa" -> "📱"
        "tv" -> "📺"
        "voucher" -> "🎟️"
        else -> category.name.take(1).uppercase()
    }

    Column(
        horizontalAlignment = Alignment.CenterHorizontally,
        modifier = modifier.clickable { onClick() }
    ) {
        Surface(
            shape = RoundedCornerShape(16.dp),
            color = iconColor,
            modifier = Modifier.size(56.dp)
        ) {
            Box(contentAlignment = Alignment.Center) {
                Text(
                    text = iconText,
                    fontSize = if (iconText.length > 2) 18.sp else 24.sp,
                    color = Color.White
                )
            }
        }
        
        Spacer(modifier = Modifier.height(8.dp))
        
        Text(
            text = category.name,
            fontSize = 11.sp,
            fontWeight = FontWeight.Medium,
            color = Color(0xFF424242),
            textAlign = TextAlign.Center,
            maxLines = 2,
            lineHeight = 13.sp,
            modifier = Modifier.padding(horizontal = 4.dp)
        )
    }
}

@Preview(showBackground = true)
@Composable
fun HomeScreenPreview() {
    val sampleCategories = listOf(
        CategoryDto("1", "Aktivasi Perdana", "aktivasi_perdana"),
        CategoryDto("2", "Aktivasi Voucher", "aktivasi_voucher"),
        CategoryDto("3", "Data", "data"),
        CategoryDto("4", "E-Money", "e-money"),
        CategoryDto("5", "Games", "games"),
        CategoryDto("6", "Gas", "gas"),
        CategoryDto("7", "Masa Aktif", "masa_aktif"),
        CategoryDto("8", "PLN", "pln"),
        CategoryDto("9", "Paket SMS & Telpon", "paket_sms_and_telpon"),
        CategoryDto("10", "Pulsa", "pulsa"),
        CategoryDto("11", "TV", "tv"),
        CategoryDto("12", "Voucher", "voucher")
    )
    val sampleWallet = WalletResponse("1", 200000.0, 0.0, 200000.0)
    
    PpoMobileTheme {
        HomeContent(
            categoriesState = Resource.Success(CategoryCollection(sampleCategories)),
            walletState = Resource.Success(sampleWallet),
            onCategoryClick = {},
            onWalletClick = {}
        )
    }
}
