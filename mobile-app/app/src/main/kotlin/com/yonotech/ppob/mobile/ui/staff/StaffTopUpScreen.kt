package com.yonotech.ppob.mobile.ui.staff

import androidx.compose.foundation.layout.*
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.ArrowBack
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.hilt.navigation.compose.hiltViewModel
import com.yonotech.ppob.mobile.data.remote.dto.StaffDto
import com.yonotech.ppob.mobile.ui.components.PpoButton
import com.yonotech.ppob.mobile.ui.components.PpoTextField
import com.yonotech.ppob.mobile.utils.Resource
import com.yonotech.ppob.mobile.viewmodels.staff.StaffViewModel
import java.text.NumberFormat
import java.util.Locale

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun StaffTopUpScreen(
    staff: StaffDto,
    onBackClick: () -> Unit,
    viewModel: StaffViewModel = hiltViewModel()
) {
    var amount by remember { mutableStateOf("") }
    var pin by remember { mutableStateOf("") }
    val topUpState by viewModel.topUpState.collectAsState()
    val currencyFormat = NumberFormat.getCurrencyInstance(Locale("in", "ID")).apply { maximumFractionDigits = 0 }

    LaunchedEffect(topUpState) {
        if (topUpState is Resource.Success) {
            onBackClick()
            viewModel.resetState()
        }
    }

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Top Up Staff") },
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
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(padding)
                .padding(24.dp)
        ) {
            Text(
                text = "Staff: ${staff.name}",
                fontSize = 18.sp,
                fontWeight = FontWeight.Bold,
                modifier = Modifier.padding(bottom = 4.dp)
            )
            Text(
                text = "Saldo: ${currencyFormat.format(staff.balance)}",
                color = MaterialTheme.colorScheme.onSurfaceVariant,
                modifier = Modifier.padding(bottom = 24.dp)
            )

            PpoTextField(
                value = amount,
                onValueChange = { amount = it },
                label = "Jumlah Top Up (Rp)",
                keyboardOptions = androidx.compose.foundation.text.KeyboardOptions(
                    keyboardType = androidx.compose.ui.text.input.KeyboardType.Number
                )
            )

            Spacer(modifier = Modifier.height(16.dp))

            PpoTextField(
                value = pin,
                onValueChange = { if (it.length <= 6) pin = it },
                label = "6-Digit PIN",
                keyboardOptions = androidx.compose.foundation.text.KeyboardOptions(
                    keyboardType = androidx.compose.ui.text.input.KeyboardType.NumberPassword
                )
            )

            Spacer(modifier = Modifier.height(32.dp))

            if (topUpState is Resource.Error) {
                Text(
                    text = (topUpState as Resource.Error).message,
                    color = MaterialTheme.colorScheme.error,
                    modifier = Modifier.padding(bottom = 16.dp)
                )
            }

            PpoButton(
                label = "Top Up Sekarang",
                onClick = {
                    viewModel.topUpStaff(staff.id, amount.toDoubleOrNull() ?: 0.0, pin)
                },
                isLoading = topUpState is Resource.Loading,
                enabled = amount.isNotEmpty() && pin.length == 6
            )
        }
    }
}