package com.yonotech.ppob.mobile.ui.product

import androidx.compose.animation.*
import androidx.compose.foundation.background
import androidx.compose.foundation.BorderStroke
import androidx.compose.foundation.border
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.grid.GridCells
import androidx.compose.foundation.lazy.grid.LazyVerticalGrid
import androidx.compose.foundation.lazy.grid.items
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.ArrowBack
import androidx.compose.material.icons.filled.Check
import androidx.compose.material.icons.filled.Person
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.hilt.navigation.compose.hiltViewModel
import com.yonotech.ppob.mobile.data.remote.dto.ProductCollection
import com.yonotech.ppob.mobile.data.remote.dto.ProductDto
import com.yonotech.ppob.mobile.data.remote.dto.WalletResponse
import com.yonotech.ppob.mobile.data.remote.dto.transaction.TransactionResponse
import com.yonotech.ppob.mobile.ui.components.*
import com.yonotech.ppob.mobile.ui.theme.PpoMobileTheme
import com.yonotech.ppob.mobile.utils.Resource
import com.yonotech.ppob.mobile.viewmodels.product.ProductViewModel
import com.yonotech.ppob.mobile.viewmodels.transaction.TransactionViewModel
import java.text.NumberFormat
import java.util.*

private enum class CheckoutStep {
    INQUIRY, SUMMARY, PIN
}

@Composable
fun GenericProductScreen(
    categoryId: String,
    categoryCode: String,
    categoryName: String,
    onBackClick: () -> Unit,
    onTransactionSuccess: (String) -> Unit,
    productViewModel: ProductViewModel = hiltViewModel(),
    transactionViewModel: TransactionViewModel = hiltViewModel()
) {
    var customerId by remember { mutableStateOf("") }
    var selectedBrand by remember { mutableStateOf<String?>(null) }
    val categoriesState by productViewModel.categories.collectAsState()
    val productsState by productViewModel.products.collectAsState()
    val transactionState by transactionViewModel.transactionState.collectAsState()
    val walletState by transactionViewModel.walletState.collectAsState()
    val inquiryState by transactionViewModel.inquiryState.collectAsState()
    var selectedProduct by remember { mutableStateOf<ProductDto?>(null) }

    // Find the current category metadata from the loaded categories
    val currentCategory = remember(categoriesState) {
        if (categoriesState is Resource.Success) {
            (categoriesState as Resource.Success).data.categories.find { it.id == categoryId }
        } else null
    }

    // Logic for UI behavior
    val showOperator = categoryCode in listOf("pulsa", "data", "masa_aktif", "paket_sms_and_telpon")
    
    // Use metadata from API if available, otherwise fallback to hardcoded defaults
    val label = currentCategory?.inputLabel ?: if (showOperator || categoryCode == "e-money") "NOMOR HANDPHONE" else "ID PELANGGAN / TUJUAN"
    val placeholder = currentCategory?.placeholder ?: if (showOperator || categoryCode == "e-money") "081..." else "Masukkan ID..."
    val keyboardType = when {
        currentCategory?.inputType == "NUMBER" -> KeyboardType.Number
        showOperator || categoryCode == "e-money" || categoryCode == "pln" -> KeyboardType.Phone
        else -> KeyboardType.Text
    }
    
    val operatorInfo = remember(customerId) {
        if (showOperator) detectOperator(customerId) else null
    }

    // Call getProducts based on logic
    LaunchedEffect(categoryId, operatorInfo, categoryCode, customerId, selectedBrand) {
        if (showOperator) {
            if (operatorInfo != null && operatorInfo.first != "Lainnya" && customerId.length <= 4) {
                productViewModel.getProducts(categoryId = categoryId, brand = operatorInfo.first)
            }
        } else if (selectedBrand != null) {
            productViewModel.getProducts(categoryId = categoryId, brand = selectedBrand)
        } else if (currentCategory?.needsInquiry != true) {
            // For non-operator and non-inquiry categories, load all products
            productViewModel.getProducts(categoryId = categoryId)
        }
    }

    LaunchedEffect(transactionState) {
        if (transactionState is Resource.Success) {
            val data = (transactionState as Resource.Success).data
            onTransactionSuccess(data.transactionId)
            transactionViewModel.resetState()
        }
    }

    GenericProductContent(
        categoryName = categoryName,
        categoryCode = categoryCode,
        customerId = customerId,
        onCustomerIdChange = { newValue -> 
            // Apply validation regex if provided
            val regex = currentCategory?.validationRegex
            if (regex == null || newValue.isEmpty() || Regex(regex).matches(newValue) || newValue.length < customerId.length) {
                if (newValue.length <= 20) customerId = newValue
            }
        },
        inputLabel = label,
        inputPlaceholder = placeholder,
        keyboardType = keyboardType,
        showOperator = showOperator,
        operatorInfo = operatorInfo,
        selectedBrand = selectedBrand,
        onBrandSelect = { selectedBrand = if (it.isEmpty()) null else it },
        productsState = productsState,
        selectedProduct = selectedProduct,
        transactionState = transactionState,
        walletState = walletState,
        inquiryState = inquiryState,
        onProductSelect = { selectedProduct = it },
        onBackClick = onBackClick,
        onLanjutPembayaran = { product ->
            if (currentCategory?.needsInquiry == true) {
                transactionViewModel.performInquiry(categoryId.toLong(), product.brand, customerId)
            } else {
                transactionViewModel.checkBalance(product.price)
            }
        },
        onConfirmPayment = { product, pin -> 
            transactionViewModel.selectedProductCode = product.code
            transactionViewModel.customerNo = customerId
            transactionViewModel.initiateTransaction(pin)
        }
    )
}

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun GenericProductContent(
    categoryName: String,
    categoryCode: String,
    customerId: String,
    onCustomerIdChange: (String) -> Unit,
    inputLabel: String,
    inputPlaceholder: String,
    keyboardType: KeyboardType,
    showOperator: Boolean,
    operatorInfo: Pair<String, Color>?,
    selectedBrand: String?,
    onBrandSelect: (String) -> Unit,
    productsState: Resource<ProductCollection>,
    selectedProduct: ProductDto?,
    transactionState: Resource<TransactionResponse>,
    walletState: Resource<WalletResponse>,
    inquiryState: Resource<InquiryResponse>,
    onProductSelect: (ProductDto) -> Unit,
    onBackClick: () -> Unit,
    onLanjutPembayaran: (ProductDto) -> Unit,
    onConfirmPayment: (ProductDto, String) -> Unit
) {
    val sheetState = rememberModalBottomSheetState(skipPartiallyExpanded = true)
    var showCheckoutSheet by remember { mutableStateOf(false) }
    var currentStep by remember { mutableStateOf(CheckoutStep.SUMMARY) }

    // Logic to open sheet and navigate steps
    LaunchedEffect(inquiryState) {
        if (inquiryState is Resource.Success && !showCheckoutSheet) {
            currentStep = CheckoutStep.INQUIRY
            showCheckoutSheet = true
        }
    }

    LaunchedEffect(walletState) {
        if (walletState is Resource.Success && showCheckoutSheet) {
            if (currentStep == CheckoutStep.SUMMARY || currentStep == CheckoutStep.INQUIRY) {
                currentStep = CheckoutStep.PIN
            }
        }
    }

    if (showCheckoutSheet && selectedProduct != null) {
        ModalBottomSheet(
            onDismissRequest = { 
                showCheckoutSheet = false
                currentStep = CheckoutStep.SUMMARY
            },
            sheetState = sheetState,
            containerColor = Color.White,
            dragHandle = { BottomSheetDefaults.DragHandle() }
        ) {
            val currencyFormat = NumberFormat.getCurrencyInstance(Locale("in", "ID"))
            currencyFormat.maximumFractionDigits = 0

            when (currentStep) {
                CheckoutStep.INQUIRY -> {
                    val inquiryData = (inquiryState as? Resource.Success)?.data
                    Column {
                        if (walletState is Resource.Error) {
                            Text(
                                text = walletState.message,
                                color = MaterialTheme.colorScheme.error,
                                modifier = Modifier.padding(horizontal = 24.dp, vertical = 8.dp),
                                fontSize = 12.sp,
                                fontWeight = FontWeight.Bold
                            )
                        }

                        CheckoutSummary(
                            title = "Detail Pelanggan",
                            items = listOf(
                                SummaryItem("Nama Pelanggan", inquiryData?.customerName ?: "-"),
                                SummaryItem("ID Pelanggan", customerId),
                                SummaryItem("Tagihan", currencyFormat.format(inquiryData?.billAmount ?: 0.0)),
                                SummaryItem("Admin Fee", currencyFormat.format(inquiryData?.adminFee ?: 0.0))
                            ),
                            totalValue = currencyFormat.format(inquiryData?.totalAmount ?: selectedProduct.price),
                            buttonLabel = "KONFIRMASI & BAYAR",
                            isLoading = walletState is Resource.Loading,
                            onButtonClick = { onLanjutPembayaran(selectedProduct) }
                        )
                    }
                }
                CheckoutStep.SUMMARY -> {
                    // Pricing breakdown (simulated for UI)
                    val basePrice = selectedProduct.price + 1500.0
                    val discount = 1500.0
                    val total = selectedProduct.price

                    Column {
                        if (walletState is Resource.Error) {
                            Surface(
                                color = MaterialTheme.colorScheme.errorContainer,
                                modifier = Modifier.fillMaxWidth().padding(horizontal = 24.dp, vertical = 8.dp),
                                shape = RoundedCornerShape(8.dp)
                            ) {
                                Text(
                                    text = walletState.message,
                                    color = MaterialTheme.colorScheme.onErrorContainer,
                                    modifier = Modifier.padding(12.dp),
                                    fontSize = 12.sp,
                                    fontWeight = FontWeight.Bold
                                )
                            }
                        }

                        CheckoutSummary(
                            title = "Ringkasan Checkout",
                            items = listOf(
                                SummaryItem("Harga Produk", currencyFormat.format(basePrice)),
                                SummaryItem("Diskon", "- ${currencyFormat.format(discount)}", valueColor = Color(0xFF4CAF50)),
                                SummaryItem("Subtotal", currencyFormat.format(total))
                            ),
                            totalValue = currencyFormat.format(total),
                            buttonLabel = "BAYAR",
                            isLoading = walletState is Resource.Loading,
                            onButtonClick = { onLanjutPembayaran(selectedProduct) }
                        )
                    }
                }
                CheckoutStep.PIN -> {
                    Column {
                        if (transactionState is Resource.Error) {
                            Text(
                                text = transactionState.message,
                                color = MaterialTheme.colorScheme.error,
                                modifier = Modifier.padding(horizontal = 24.dp, vertical = 8.dp),
                                fontSize = 12.sp,
                                fontWeight = FontWeight.Bold
                            )
                        }
                        
                        PinAuthContent(
                            onPinComplete = { pin ->
                                onConfirmPayment(selectedProduct, pin)
                            },
                            isLoading = transactionState is Resource.Loading
                        )
                    }
                }
            }
        }
    }

    Scaffold(
        topBar = {
            TopAppBar(
                title = { 
                    Text(
                        categoryName.uppercase(), 
                        style = MaterialTheme.typography.titleMedium.copy(
                            fontWeight = FontWeight.ExtraBold,
                            letterSpacing = 1.sp
                        )
                    ) 
                },
                navigationIcon = {
                    IconButton(onClick = onBackClick) {
                        Icon(Icons.AutoMirrored.Filled.ArrowBack, contentDescription = "Back")
                    }
                },
                actions = {
                    Surface(
                        shape = CircleShape,
                        color = MaterialTheme.colorScheme.primary.copy(alpha = 0.1f),
                        modifier = Modifier.padding(end = 12.dp).size(32.dp)
                    ) {
                        Box(contentAlignment = Alignment.Center) {
                            Icon(
                                Icons.Default.Person, 
                                contentDescription = "Profile",
                                modifier = Modifier.size(18.dp),
                                tint = MaterialTheme.colorScheme.primary
                            )
                        }
                    }
                },
                colors = TopAppBarDefaults.topAppBarColors(containerColor = Color.White)
            )
        },
        bottomBar = {
            AnimatedVisibility(
                visible = selectedProduct != null,
                enter = slideInVertically { it } + fadeIn(),
                exit = slideOutVertically { it } + fadeOut()
            ) {
                Surface(
                    modifier = Modifier.fillMaxWidth(),
                    shadowElevation = 8.dp,
                    color = Color.White
                ) {
                    PpoButton(
                        label = "LANJUT PEMBAYARAN",
                        onClick = { 
                            currentStep = CheckoutStep.SUMMARY
                            showCheckoutSheet = true 
                        },
                        modifier = Modifier.padding(16.dp)
                    )
                }
            }
        },
        containerColor = Color(0xFFF8F9FA)
    ) { padding ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(padding)
                .padding(16.dp)
        ) {
            if ((categoryCode == "pln" || categoryCode == "e-money") && selectedBrand == null) {
                Text(
                    text = "PILIH SUB-KATEGORI",
                    fontSize = 12.sp,
                    fontWeight = FontWeight.Bold,
                    color = Color.Gray,
                    modifier = Modifier.padding(bottom = 16.dp)
                )
                
                val brands = if (categoryCode == "pln") {
                    listOf("PLN Token", "PLN Pasca", "PLN Non-Taglis")
                } else {
                    listOf("DANA", "GO PAY", "OVO", "SHOPEE PAY", "LinkAja")
                }
                
                brands.chunked(2).forEach { rowBrands ->
                    Row(modifier = Modifier.fillMaxWidth(), horizontalArrangement = Arrangement.spacedBy(12.dp)) {
                        rowBrands.forEach { brand ->
                            Card(
                                onClick = { onBrandSelect(brand) },
                                modifier = Modifier.weight(1f).padding(bottom = 12.dp),
                                colors = CardDefaults.cardColors(containerColor = Color.White)
                            ) {
                                Box(modifier = Modifier.fillMaxWidth().padding(16.dp), contentAlignment = Alignment.Center) {
                                    Text(text = brand, fontWeight = FontWeight.Bold)
                                }
                            }
                        }
                        if (rowBrands.size < 2) Spacer(modifier = Modifier.weight(1f))
                    }
                }
            } else {
                // Input Card
                Surface(
                modifier = Modifier.fillMaxWidth(),
                shape = RoundedCornerShape(20.dp),
                color = Color.White,
                shadowElevation = 1.dp
            ) {
                Column(modifier = Modifier.padding(20.dp)) {
                    Row(modifier = Modifier.fillMaxWidth(), horizontalArrangement = Arrangement.SpaceBetween, verticalAlignment = Alignment.CenterVertically) {
                      Text(
                          inputLabel.uppercase(),
                          fontSize = 12.sp,
                          fontWeight = FontWeight.Bold,
                          color = Color.Gray,
                          letterSpacing = 0.5.sp
                      )

                      if (selectedBrand != null) {
                          Text(
                              text = "GANTI ($selectedBrand)",
                              fontSize = 10.sp,
                              color = MaterialTheme.colorScheme.primary,
                              fontWeight = FontWeight.Bold,
                              modifier = Modifier.clickable { onBrandSelect("") } // Reset brand
                          )
                      }
                    }

                    Box(modifier = Modifier.fillMaxWidth(), contentAlignment = Alignment.CenterEnd) {                        TextField(
                            value = customerId,
                            onValueChange = { onCustomerIdChange(it) },
                            placeholder = { Text(inputPlaceholder, color = Color.LightGray) },
                            modifier = Modifier.fillMaxWidth(),
                            colors = TextFieldDefaults.colors(
                                focusedContainerColor = Color.Transparent,
                                unfocusedContainerColor = Color.Transparent,
                                disabledContainerColor = Color.Transparent,
                                focusedIndicatorColor = MaterialTheme.colorScheme.primary,
                                unfocusedIndicatorColor = Color(0xFFF5F5F5)
                            ),
                            textStyle = MaterialTheme.typography.headlineSmall.copy(fontWeight = FontWeight.Bold),
                            keyboardOptions = KeyboardOptions(keyboardType = keyboardType),
                            singleLine = true
                        )
                        
                        if (showOperator) {
                            operatorInfo?.let { (name, color) ->
                                Surface(
                                    color = color.copy(alpha = 0.1f),
                                    shape = RoundedCornerShape(8.dp),
                                    modifier = Modifier.padding(bottom = 8.dp),
                                    border = BorderStroke(1.dp, color.copy(alpha = 0.2f))
                                ) {
                                    Row(
                                        modifier = Modifier.padding(horizontal = 10.dp, vertical = 6.dp),
                                        verticalAlignment = Alignment.CenterVertically
                                    ) {
                                        Text(
                                            text = name,
                                            fontSize = 10.sp,
                                            fontWeight = FontWeight.ExtraBold,
                                            color = color
                                        )
                                    }
                                }
                            }
                        }
                    }
                    
                    if (showOperator) {
                        Text(
                            "Masukkan nomor tujuan untuk mendeteksi operator otomatis.",
                            fontSize = 10.sp,
                            color = Color.LightGray,
                            modifier = Modifier.padding(top = 8.dp)
                        )
                    }
                }
            }

            Spacer(modifier = Modifier.height(24.dp))

            Text(
                "PILIH PRODUK",
                fontSize = 12.sp,
                fontWeight = FontWeight.ExtraBold,
                color = Color(0xFF424242),
                letterSpacing = 1.5.sp,
                modifier = Modifier.padding(start = 4.dp, bottom = 16.dp)
            )

            when (productsState) {
                is Resource.Loading -> {
                    Box(modifier = Modifier.fillMaxWidth().height(200.dp), contentAlignment = Alignment.Center) {
                        CircularProgressIndicator(color = MaterialTheme.colorScheme.primary)
                    }
                }
                is Resource.Success -> {
                    val productCollection = productsState.data
                    if (productCollection.products.isEmpty()) {
                        Box(modifier = Modifier.fillMaxWidth().height(100.dp), contentAlignment = Alignment.Center) {
                            Text("Tidak ada produk tersedia", color = Color.Gray)
                        }
                    } else {
                        LazyVerticalGrid(
                            columns = GridCells.Fixed(2),
                            horizontalArrangement = Arrangement.spacedBy(12.dp),
                            verticalArrangement = Arrangement.spacedBy(12.dp),
                            modifier = Modifier.fillMaxSize()
                        ) {
                            items(productCollection.products) { product ->
                                DenomCard(
                                    product = product,
                                    categoryName = categoryName,
                                    isSelected = selectedProduct?.id == product.id,
                                    onClick = { onProductSelect(product) }
                                )
                            }
                        }
                    }
                }
                is Resource.Error -> {
                    Text(text = productsState.message, color = MaterialTheme.colorScheme.error)
                }
                else -> {}
            }
        }
    }
}

@Composable
fun DenomCard(
    product: ProductDto,
    categoryName: String,
    isSelected: Boolean,
    onClick: () -> Unit
) {
    val currencyFormat = NumberFormat.getCurrencyInstance(Locale("in", "ID"))
    currencyFormat.maximumFractionDigits = 0
    
    // Simulate some promo logic
    val priceThreshold = product.name.filter { it.isDigit() }.toDoubleOrNull()?.let { it * 1.1 }
    val isPromo = priceThreshold?.let { product.price < it } ?: false

    Box(
        modifier = Modifier
            .fillMaxWidth()
            .clip(RoundedCornerShape(20.dp))
            .background(Color.White)
            .border(
                width = if (isSelected) 2.dp else 1.dp,
                color = if (isSelected) MaterialTheme.colorScheme.primary else Color(0xFFF5F5F5),
                shape = RoundedCornerShape(20.dp)
            )
            .clickable { onClick() }
            .padding(16.dp)
    ) {
        if (isSelected) {
            Surface(
                color = Color(0xFF4CAF50),
                shape = CircleShape,
                modifier = Modifier.size(20.dp).align(Alignment.TopEnd)
            ) {
                Icon(
                    Icons.Default.Check, 
                    contentDescription = null, 
                    tint = Color.White,
                    modifier = Modifier.size(12.dp).padding(4.dp)
                )
            }
        } else if (isPromo) {
            Surface(
                color = Color(0xFFFF9800),
                modifier = Modifier
                    .align(Alignment.TopEnd)
                    .offset(x = 16.dp, y = (-8).dp)
                    .width(60.dp)
                    .height(20.dp),
                shape = RoundedCornerShape(bottomStart = 8.dp)
            ) {
                Box(contentAlignment = Alignment.Center) {
                    Text("PROMO", color = Color.White, fontSize = 8.sp, fontWeight = FontWeight.Bold)
                }
            }
        }

        Column {
            Text(
                categoryName.uppercase(), 
                fontSize = 10.sp, 
                fontWeight = FontWeight.Bold, 
                color = MaterialTheme.colorScheme.primary.copy(alpha = 0.7f),
                maxLines = 1
            )
            Text(
                text = product.name,
                fontSize = 16.sp,
                fontWeight = FontWeight.Black,
                color = Color(0xFF424242),
                modifier = Modifier.padding(vertical = 4.dp),
                maxLines = 2
            )
            Text(
                text = "Harga: ${currencyFormat.format(product.price)}",
                fontSize = 11.sp,
                fontWeight = FontWeight.SemiBold,
                color = Color.Gray
            )
        }
    }
}

private fun detectOperator(phone: String): Pair<String, Color>? {
    if (phone.length < 4) return null
    val prefix = phone.take(4)
    return when {
        listOf("0811", "0812", "0813", "0821", "0822", "0823", "0851", "0852", "0853").contains(prefix) ->
            "Telkomsel" to Color(0xFFE53935)
        listOf("0817", "0818", "0819", "0859", "0877", "0878").contains(prefix) ->
            "XL Axiata" to Color(0xFF1E88E5)
        listOf("0814", "0815", "0816", "0855", "0856", "0857", "0858").contains(prefix) ->
            "Indosat IM3" to Color(0xFFFBC02D)
        listOf("0895", "0896", "0897", "0898", "0899").contains(prefix) ->
            "Tri (3)" to Color(0xFFF50057)
        listOf("0881", "0882", "0883", "0884", "0885", "0886", "0887", "0888", "0889").contains(prefix) ->
            "Smartfren" to Color(0xFFE91E63)
        listOf("0831", "0832", "0833", "0838").contains(prefix) ->
            "Axis" to Color(0xFF8E24AA)
        else -> "Lainnya" to Color.Gray
    }
}

@Preview(showBackground = true)
@Composable
fun GenericProductScreenPreview() {
    val sampleProducts = listOf(
        ProductDto("1", "PULSA 5.000", "p5", "1", "Telkomsel", 5500.0, "Pulsa 5rb", "active"),
        ProductDto("2", "PULSA 10.000", "p10", "1", "Telkomsel", 10500.0, "Pulsa 10rb", "active")
    )
    PpoMobileTheme {
        GenericProductContent(
            categoryName = "Pulsa",
            categoryCode = "pulsa",
            customerId = "08123456789",
            onCustomerIdChange = {},
            inputLabel = "Nomor HP",
            inputPlaceholder = "081...",
            keyboardType = KeyboardType.Phone,
            showOperator = true,
            operatorInfo = "Telkomsel" to Color.Red,
            selectedBrand = null,
            onBrandSelect = {},
            productsState = Resource.Success(ProductCollection(sampleProducts)),
            selectedProduct = sampleProducts[1],
            transactionState = Resource.Idle,
            walletState = Resource.Idle,
            inquiryState = Resource.Idle,
            onProductSelect = {},
            onBackClick = {},
            onLanjutPembayaran = {},
            onConfirmPayment = { _, _ -> }
        )
    }
}
