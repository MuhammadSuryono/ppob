package com.yonotech.ppob.presentation.transaction

import android.util.Patterns
import androidx.compose.foundation.Canvas
import androidx.compose.foundation.background
import androidx.compose.foundation.border
import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.layout.size
import androidx.compose.foundation.layout.width
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.Button
import androidx.compose.material3.CircularProgressIndicator
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
import androidx.compose.ui.geometry.Offset
import androidx.compose.ui.geometry.Size
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.StrokeCap
import androidx.compose.ui.graphics.drawscope.Stroke
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.unit.dp
import com.yonotech.ppob.presentation.transaction.TransactionState

@Composable
fun PinAuthorizationScreen(
    viewModel: com.yonotech.ppob.presentation.transaction.TransactionInitViewModel = androidx.lifecycle.viewmodel.compose.viewModel(),
    transactionId: String,
    onPinEntered: () -> Unit,
    onBack: () -> Unit
) {
    val txState by viewModel.transactionState.collectAsState()

    LaunchedEffect(txState) {
        if (txState is TransactionState.Success || txState is TransactionState.Pending) {
            onPinEntered()
        }
    }

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Otorisasi PIN") },
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
                .fillMaxSize()
                .padding(paddingValues)
                .padding(horizontal = 24.dp),
            horizontalAlignment = Alignment.CenterHorizontally,
            verticalArrangement = Arrangement.Center
        ) {
            // PIN Pad Visual Indicator
            PinDotsIndicator(
                pinLength = 6,
                filled = viewModel.uiState.value.pin.length,
                modifier = Modifier.padding(bottom = 48.dp)
            )

            // Number Pad
            NumberPad(
                onNumberClick = { number ->
                    if (viewModel.uiState.value.pin.length < 6) {
                        viewModel.onPinChange(viewModel.uiState.value.pin + number)
                    }
                },
                onDeleteClick = {
                    if (viewModel.uiState.value.pin.isNotEmpty()) {
                        viewModel.onPinChange(viewModel.uiState.value.pin.dropLast(1))
                    }
                },
                onSubmitClick = {
                    if (viewModel.uiState.value.pin.length == 6) {
                        // Submit transaction with PIN
                        viewModel.uiState.value.let { state ->
                            viewModel.initiateTransaction(
                                pin = state.pin,
                                authToken = "", // TODO: Get from DataStore
                                customerNo = state.customerNumber.value,
                                productId = "", // Should be set
                                sellingPrice = state.sellingPrice.value
                            )
                        }
                    }
                }
            )

            // Error display
            if (txState is TransactionState.Error) {
                Text(
                    text = (txState as TransactionState.Error).message,
                    color = MaterialTheme.colorScheme.error,
                    style = MaterialTheme.typography.bodyMedium,
                    modifier = Modifier.padding(top = 16.dp)
                )
            }

            if (txState is TransactionState.Loading) {
                CircularProgressIndicator(
                    modifier = Modifier.padding(top = 16.dp)
                )
            }
        }
    }
}

@Composable
fun PinDotsIndicator(
    pinLength: Int = 6,
    filled: Int = 0,
    modifier: Modifier = Modifier
) {
    Row(
        modifier = modifier,
        horizontalArrangement = Arrangement.spacedBy(12.dp)
    ) {
        for (i in 0 until pinLength) {
            val isFilled = i < filled
            Box(
                modifier = Modifier
                    .size(16.dp)
                    .then(
                        if (isFilled) {
                            Modifier
                                .background(
                                    color = MaterialTheme.colorScheme.primary,
                                    shape = RoundedCornerShape(8.dp)
                                )
                        } else {
                            Modifier
                                .border(
                                    width = 2.dp,
                                    color = MaterialTheme.colorScheme.outline,
                                    shape = RoundedCornerShape(8.dp)
                                )
                        }
                    )
            )
        }
    }
}

@Composable
fun NumberPad(
    onNumberClick: (String) -> Unit,
    onDeleteClick: () -> Unit,
    onSubmitClick: () -> Unit
) {
    val numbers = listOf("1", "2", "3", "4", "5", "6", "7", "8", "9")

    Column(
        modifier = Modifier.padding(vertical = 16.dp),
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.spacedBy(8.dp)
    ) {
        // Rows 1-3
        for (row in 0..2) {
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.SpaceEvenly
            ) {
                for (col in 0..2) {
                    val index = row * 3 + col
                    if (index < numbers.size) {
                        PinNumberButton(
                            number = numbers[index],
                            onClick = { onNumberClick(numbers[index]) }
                        )
                    }
                }
            }
        }

        // Row 4: Empty, 0, Delete
        Row(
            modifier = Modifier.fillMaxWidth()
                .padding(top = 8.dp),
            horizontalArrangement = Arrangement.SpaceEvenly,
            verticalAlignment = Alignment.CenterVertically
        ) {
            // Empty placeholder
            Box(
                modifier = Modifier
                    .size(72.dp)
                    .padding(8.dp)
            )

            PinNumberButton(
                number = "0",
                onClick = { onNumberClick("0") }
            )

            androidx.compose.material3.IconButton(
                onClick = onDeleteClick,
                modifier = Modifier
                    .size(72.dp)
                    .padding(8.dp)
            ) {
                androidx.compose.material.icons.Icons.Filled.Backspace
            }
        }

        Spacer(modifier = Modifier.height(16.dp))

        Button(
            onClick = onSubmitClick,
            modifier = Modifier
                .fillMaxWidth()
                .height(56.dp),
            shape = RoundedCornerShape(12.dp)
        ) {
            Text(
                text = "Konfirmasi PIN",
                style = MaterialTheme.typography.titleMedium
            )
        }
    }
}

@Composable
fun PinNumberButton(
    number: String,
    onClick: () -> Unit
) {
    androidx.compose.material3.Button(
        onClick = onClick,
        modifier = Modifier
            .size(72.dp)
            .padding(4.dp),
        shape = CircleShape,
        colors = androidx.compose.material3.ButtonDefaults.buttonColors(
            containerColor = MaterialTheme.colorScheme.surfaceVariant
        )
    ) {
        Text(
            text = number,
            style = MaterialTheme.typography.headlineMedium,
            fontWeight = FontWeight.Medium
        )
    }
}