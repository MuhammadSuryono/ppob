package com.yonotech.ppob.mobile.ui.components

import androidx.compose.animation.*
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.grid.GridCells
import androidx.compose.foundation.lazy.grid.LazyVerticalGrid
import androidx.compose.foundation.lazy.grid.items
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.automirrored.filled.Backspace
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.draw.clip
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import com.yonotech.ppob.mobile.ui.theme.PpoMobileTheme

@Composable
fun PpoButton(
    label: String,
    onClick: () -> Unit,
    modifier: Modifier = Modifier,
    enabled: Boolean = true,
    isLoading: Boolean = false,
    containerColor: Color = MaterialTheme.colorScheme.primary,
    contentColor: Color = MaterialTheme.colorScheme.onPrimary
) {
    Button(
        onClick = onClick,
        modifier = modifier
            .fillMaxWidth()
            .height(56.dp),
        enabled = enabled && !isLoading,
        shape = RoundedCornerShape(16.dp),
        colors = ButtonDefaults.buttonColors(
            containerColor = containerColor,
            contentColor = contentColor
        )
    ) {
        if (isLoading) {
            CircularProgressIndicator(
                modifier = Modifier.size(24.dp),
                color = contentColor,
                strokeWidth = 2.dp
            )
        } else {
            Text(
                text = label, 
                style = MaterialTheme.typography.titleMedium,
                fontWeight = FontWeight.ExtraBold,
                letterSpacing = 1.sp
            )
        }
    }
}

@Composable
fun PpoCard(
    modifier: Modifier = Modifier,
    content: @Composable () -> Unit
) {
    Card(
        modifier = modifier,
        shape = RoundedCornerShape(12.dp),
        elevation = CardDefaults.cardElevation(defaultElevation = 2.dp),
        colors = CardDefaults.cardColors(containerColor = MaterialTheme.colorScheme.surface)
    ) {
        content()
    }
}

@Composable
fun CheckoutSummary(
    title: String = "Ringkasan Pembayaran",
    items: List<SummaryItem>,
    totalLabel: String = "Total Pembayaran",
    totalValue: String,
    buttonLabel: String = "BAYAR SEKARANG",
    isLoading: Boolean = false,
    onButtonClick: () -> Unit,
    modifier: Modifier = Modifier
) {
    Column(
        modifier = modifier
            .fillMaxWidth()
            .padding(24.dp)
    ) {
        Text(
            text = title,
            style = MaterialTheme.typography.titleMedium,
            fontWeight = FontWeight.ExtraBold,
            modifier = Modifier.padding(bottom = 20.dp)
        )
        
        items.forEach { item ->
            Row(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(vertical = 4.dp),
                horizontalArrangement = Arrangement.SpaceBetween
            ) {
                Text(item.label, color = Color.Gray, fontSize = 14.sp)
                Text(
                    item.value, 
                    color = item.valueColor, 
                    fontWeight = FontWeight.Bold, 
                    fontSize = 14.sp
                )
            }
        }
        
        Spacer(modifier = Modifier.height(16.dp))
        HorizontalDivider(thickness = 0.5.dp, color = Color(0xFFEEEEEE))
        Spacer(modifier = Modifier.height(16.dp))
        
        Row(
            modifier = Modifier.fillMaxWidth(),
            horizontalArrangement = Arrangement.SpaceBetween,
            verticalAlignment = Alignment.CenterVertically
        ) {
            Text(totalLabel, fontWeight = FontWeight.Bold, fontSize = 16.sp)
            Text(
                totalValue,
                style = MaterialTheme.typography.headlineSmall,
                fontWeight = FontWeight.Black,
                color = MaterialTheme.colorScheme.primary
            )
        }
        
        Spacer(modifier = Modifier.height(24.dp))
        
        PpoButton(
            label = buttonLabel,
            onClick = onButtonClick,
            isLoading = isLoading
        )
    }
}

data class SummaryItem(
    val label: String,
    val value: String,
    val valueColor: Color = Color.Black
)

@Composable
fun PinAuthContent(
    onPinComplete: (String) -> Unit,
    title: String = "Otorisasi Pembayaran",
    subtitle: String = "Masukkan 6-digit PIN Anda untuk melanjutkan",
    isLoading: Boolean = false,
    modifier: Modifier = Modifier
) {
    var pin by remember { mutableStateOf("") }

    if (isLoading) {
        // Reset PIN if loading to allow retry on error
        LaunchedEffect(Unit) {
            // Optional: reset pin or handle state
        }
    }

    Column(
        modifier = modifier
            .fillMaxWidth()
            .padding(24.dp),
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        Text(
            text = title,
            style = MaterialTheme.typography.titleLarge,
            fontWeight = FontWeight.ExtraBold
        )
        Text(
            text = subtitle,
            fontSize = 14.sp,
            color = Color.Gray,
            modifier = Modifier.padding(top = 8.dp, bottom = 32.dp)
        )

        Box(contentAlignment = Alignment.Center) {
            PinDots(pin = pin)
            if (isLoading) {
                CircularProgressIndicator(
                    modifier = Modifier.size(48.dp),
                    strokeWidth = 3.dp,
                    color = MaterialTheme.colorScheme.primary.copy(alpha = 0.5f)
                )
            }
        }

        Spacer(modifier = Modifier.height(48.dp))

        PinPad(
            onDigit = { digit -> 
                if (pin.length < 6 && !isLoading) {
                    pin += digit
                    if (pin.length == 6) {
                        onPinComplete(pin)
                        // Clear pin after complete so it's ready for next attempt if fails
                        pin = ""
                    }
                }
            },
            onBackspace = { if (pin.isNotEmpty() && !isLoading) pin = pin.dropLast(1) }
        )
        
        Spacer(modifier = Modifier.height(24.dp))
    }
}

@Composable
fun PpoTextField(
    value: String,
    onValueChange: (String) -> Unit,
    label: String,
    modifier: Modifier = Modifier,
    isError: Boolean = false,
    errorText: String? = null,
    visualTransformation: androidx.compose.ui.text.input.VisualTransformation = androidx.compose.ui.text.input.VisualTransformation.None,
    keyboardOptions: androidx.compose.foundation.text.KeyboardOptions = androidx.compose.foundation.text.KeyboardOptions.Default,
    trailingIcon: @Composable (() -> Unit)? = null,
    placeholder: String? = null
) {
    OutlinedTextField(
        value = value,
        onValueChange = onValueChange,
        label = { Text(label) },
        placeholder = { if (placeholder != null) Text(placeholder) },
        modifier = modifier.fillMaxWidth(),
        isError = isError,
        supportingText = {
            if (isError && errorText != null) {
                Text(text = errorText, color = MaterialTheme.colorScheme.error)
            }
        },
        shape = RoundedCornerShape(12.dp),
        visualTransformation = visualTransformation,
        keyboardOptions = keyboardOptions,
        singleLine = true,
        trailingIcon = trailingIcon
    )
}

@Composable
fun PinDots(
    pin: String,
    modifier: Modifier = Modifier,
    length: Int = 6
) {
    Row(
        modifier = modifier,
        horizontalArrangement = Arrangement.spacedBy(12.dp)
    ) {
        repeat(length) { index ->
            val isFilled = index < pin.length
            Surface(
                modifier = Modifier.size(16.dp),
                shape = CircleShape,
                color = if (isFilled) MaterialTheme.colorScheme.primary else MaterialTheme.colorScheme.outlineVariant
            ) {}
        }
    }
}

@Composable
fun PinPad(
    onDigit: (String) -> Unit,
    onBackspace: () -> Unit,
    modifier: Modifier = Modifier
) {
    val buttons = listOf("1", "2", "3", "4", "5", "6", "7", "8", "9", "", "0", "back")
    
    LazyVerticalGrid(
        columns = GridCells.Fixed(3),
        modifier = modifier.width(280.dp),
        horizontalArrangement = Arrangement.spacedBy(16.dp),
        verticalArrangement = Arrangement.spacedBy(16.dp)
    ) {
        items(buttons) { btn ->
            when (btn) {
                "" -> Spacer(modifier = Modifier.size(64.dp))
                "back" -> {
                    IconButton(
                        onClick = onBackspace,
                        modifier = Modifier
                            .size(64.dp)
                            .clip(CircleShape)
                            .background(
                                color = Color.Transparent
                            )
                    ) {
                        Icon(Icons.AutoMirrored.Filled.Backspace, contentDescription = "backspace", modifier = Modifier.size(32.dp))
                    }
                }
                else -> {
                    TextButton(
                        onClick = { onDigit(btn) },
                        modifier = Modifier.size(85.dp),
                        shape = CircleShape,
                        colors = ButtonDefaults.textButtonColors(
                            containerColor = MaterialTheme.colorScheme.surfaceVariant,
                            contentColor = MaterialTheme.colorScheme.onSurfaceVariant
                        ),
                        contentPadding = PaddingValues(0.dp)
                    ) {
                        Box(contentAlignment = Alignment.Center) {
                            Text(text = btn, fontSize = 24.sp, fontWeight = FontWeight.Bold)
                        }
                    }
                }
            }
        }
    }
}

@Preview(showBackground = true)
@Composable
fun CheckoutSummaryPreview() {
    PpoMobileTheme {
        CheckoutSummary(
            totalValue = "Rp10.000",
            items = listOf(
                SummaryItem("Harga Produk", "Rp11.500"),
                SummaryItem("Diskon", "- Rp1.500", valueColor = Color(0xFF4CAF50)),
                SummaryItem("Subtotal", "Rp10.000")
            ),
            onButtonClick = {}
        )
    }
}

@Preview(showBackground = true)
@Composable
fun PinAuthContentPreview() {
    PpoMobileTheme {
        PinAuthContent(onPinComplete = {})
    }
}
