package com.yonotech.ppob.widgets

import androidx.compose.foundation.BorderStroke
import androidx.compose.foundation.layout.PaddingValues
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material3.Button
import androidx.compose.material3.ButtonDefaults
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.OutlinedButton
import androidx.compose.material3.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.unit.dp

enum class PPOBButtonVariant {
    PRIMARY, SECONDARY, TEXT, DANGER
}

@Composable
fun PPOBButton(
    text: String,
    onClick: () -> Unit,
    variant: PPOBButtonVariant = PPOBButtonVariant.PRIMARY,
    modifier: Modifier = Modifier,
    isLoading: Boolean = false,
    isEnabled: Boolean = true,
    paddingValues: PaddingValues = PaddingValues(horizontal = 24.dp, vertical = 14.dp)
) {
    val (bgColor, textColor, borderColor) = when (variant) {
        PPOBButtonVariant.PRIMARY -> Triple(
            MaterialTheme.colorScheme.primary,
            Color.White,
            null
        )
        PPOBButtonVariant.SECONDARY -> Triple(
            MaterialTheme.colorScheme.secondaryContainer,
            MaterialTheme.colorScheme.onSecondaryContainer,
            BorderStroke(1.dp, MaterialTheme.colorScheme.secondaryContainer)
        )
        PPOBButtonVariant.TEXT -> Triple(
            Color.Transparent,
            MaterialTheme.colorScheme.primary,
            null
        )
        PPOBButtonVariant.DANGER -> Triple(
            MaterialTheme.colorScheme.errorContainer,
            MaterialTheme.colorScheme.onErrorContainer,
            null
        )
    }

    if (borderColor != null) {
        OutlinedButton(
            onClick = onClick,
            modifier = modifier.fillMaxWidth(),
            enabled = isEnabled && !isLoading,
            colors = ButtonDefaults.outlinedButtonColors(
                contentColor = textColor
            ),
            border = borderColor,
            shape = RoundedCornerShape(12.dp),
            contentPadding = paddingValues
        ) {
            if (isLoading) {
                androidx.compose.material3.CircularProgressIndicator(
                    modifier = Modifier.size(20.dp),
                    color = textColor,
                    strokeWidth = 2.dp
                )
                Spacer(modifier = androidx.compose.foundation.layout.width(8.dp))
            }
            Text(text)
        }
    } else {
        Button(
            onClick = onClick,
            modifier = modifier.fillMaxWidth(),
            enabled = isEnabled && !isLoading,
            colors = ButtonDefaults.buttonColors(
                containerColor = bgColor,
                contentColor = textColor
            ),
            shape = RoundedCornerShape(12.dp),
            contentPadding = paddingValues
        ) {
            if (isLoading) {
                androidx.compose.material3.CircularProgressIndicator(
                    modifier = Modifier.size(20.dp),
                    color = textColor,
                    strokeWidth = 2.dp
                )
                Spacer(modifier = androidx.compose.foundation.layout.width(8.dp))
            }
            Text(text)
        }
    }
}