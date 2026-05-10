package com.yonotech.ppob.widgets

import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.OutlinedTextField
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.input.PasswordVisualTransformation
import androidx.compose.ui.text.input.VisualTransformation

enum class PPOBInputVariant {
    TEXT, PASSWORD, PIN, PHONE
}

@Composable
fun PPOBInput(
    value: String,
    onValueChange: (String) -> Unit,
    label: String,
    variant: PPOBInputVariant = PPOBInputVariant.TEXT,
    errorText: String? = null,
    isEnabled: Boolean = true,
    modifier: Modifier = Modifier
) {
    val visualTransformation = when (variant) {
        PPOBInputVariant.PASSWORD -> PasswordVisualTransformation()
        PPOBInputVariant.PIN -> PasswordVisualTransformation()
        else -> VisualTransformation.None
    }

    val keyboardType = when (variant) {
        PPOBInputVariant.PHONE -> androidx.compose.foundation.text.KeyboardType.Phone
        PPOBInputVariant.PIN -> androidx.compose.foundation.text.KeyboardType.NumberPassword
        PPOBInputVariant.PASSWORD -> androidx.compose.foundation.text.KeyboardType.Password
        else -> androidx.compose.foundation.text.KeyboardType.Text
    }

    OutlinedTextField(
        value = value,
        onValueChange = onValueChange,
        label = { Text(label) },
        visualTransformation = visualTransformation,
        keyboardOptions = androidx.compose.foundation.text.KeyboardOptions(keyboardType = keyboardType),
        isError = errorText != null,
        supportingText = {
            errorText?.let {
                Text(
                    text = it,
                    color = Color.Red
                )
            }
        },
        modifier = modifier.fillMaxWidth(),
        enabled = isEnabled,
        shape = RoundedCornerShape(12.dp)
    )
}