package com.yonotech.ppob.mobile.ui.components

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.grid.GridCells
import androidx.compose.foundation.lazy.grid.LazyVerticalGrid
import androidx.compose.foundation.lazy.grid.items
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Backspace
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp

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
        shape = RoundedCornerShape(12.dp),
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
            Text(text = label, style = MaterialTheme.typography.titleMedium)
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
                        modifier = Modifier.size(64.dp)
                    ) {
                        Icon(Icons.Default.Backspace, contentDescription = "backspace", modifier = Modifier.size(32.dp))
                    }
                }
                else -> {
                    TextButton(
                        onClick = { onDigit(btn) },
                        modifier = Modifier.size(64.dp),
                        shape = CircleShape,
                        colors = ButtonDefaults.textButtonColors(containerColor = MaterialTheme.colorScheme.surfaceVariant)
                    ) {
                        Text(text = btn, fontSize = 24.sp, fontWeight = FontWeight.Bold)
                    }
                }
            }
        }
    }
}
