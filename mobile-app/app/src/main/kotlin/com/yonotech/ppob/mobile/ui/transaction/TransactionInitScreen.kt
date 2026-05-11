package com.yonotech.ppob.mobile.ui.transaction

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.ArrowBack
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.yonotech.ppob.mobile.ui.components.PpoButton
import com.yonotech.ppob.mobile.ui.components.PpoTextField
import com.yonotech.ppob.mobile.viewmodels.transaction.TransactionViewModel

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun TransactionInitScreen(
    productId: String,
    onNext: () -> Unit,
    onBack: () -> Unit,
    viewModel: TransactionViewModel
) {
    var customerNo by remember { mutableStateOf("") }

    viewModel.selectedProductId = productId

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Detail Transaksi") },
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
                .padding(24.dp)
        ) {
            Text(
                text = "Informasi Pelanggan",
                fontSize = 18.sp,
                fontWeight = FontWeight.Bold,
                modifier = Modifier.padding(bottom = 16.dp)
            )

            PpoTextField(
                value = customerNo,
                onValueChange = { customerNo = it },
                label = "Nomor Pelanggan / HP",
                keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Phone)
            )

            Spacer(modifier = Modifier.weight(1f))

            PpoButton(
                label = "Lanjutkan",
                onClick = {
                    viewModel.customerNo = customerNo
                    onNext()
                },
                enabled = customerNo.isNotEmpty()
            )
        }
    }
}