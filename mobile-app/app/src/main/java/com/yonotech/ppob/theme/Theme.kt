package com.yonotech.ppob.theme

import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.graphics.Color
// ============== Shapes ==============
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.ui.unit.dp
// ============== Typography ==============
import androidx.compose.ui.text.TextStyle
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.sp

// ============== Color System ==============

val PrimaryGreen = Color(0xFF4CAF50)
val PrimaryGreenLight = Color(0xFF81C784)
val PrimaryGreenDark = Color(0xFF1B5E20)
val SecondaryBlue = Color(0xFF2196F3)
val AccentOrange = Color(0xFFFF9800)
val BackgroundLight = Color(0xFFF5F5F5)
val SurfaceWhite = Color(0xFFFFFFFF)
val ErrorRed = Color(0xFFF44336)
val WarningYellow = Color(0xFFFFC107)
val SuccessGreen = Color(0xFF4CAF50)
val PendingAmber = Color(0xFFFFA726)
val FailedRed = Color(0xFFE53935)
val TextPrimary = Color(0xFF212121)
val TextSecondary = Color(0xFF757575)
val TextHint = Color(0xFFBDBDBD)
val DividerLight = Color(0xFFE0E0E0)

// Light Color Scheme
val PPOBColorScheme = lightColorScheme(
    primary = PrimaryGreen,
    onPrimary = SurfaceWhite,
    primaryContainer = PrimaryGreenLight,
    onPrimaryContainer = PrimaryGreenDark,
    secondary = SecondaryBlue,
    onSecondary = SurfaceWhite,
    tertiary = AccentOrange,
    onTertiary = SurfaceWhite,
    background = BackgroundLight,
    onBackground = TextPrimary,
    surface = SurfaceWhite,
    onSurface = TextPrimary,
    surfaceVariant = Color(0xFFF0F0F0),
    onSurfaceVariant = TextSecondary,
    error = ErrorRed,
    onError = SurfaceWhite,
    outline = DividerLight,
    outlineVariant = Color(0xFFE0E0E0),
    inverseSurface = Color(0xFF1A1A2E),
    inverseOnSurface = Color(0xFFFFFFFF),
    inversePrimary = PrimaryGreenLight
)

val PPOBTypography = Typography(
    displayLarge = TextStyle(
        fontWeight = FontWeight.Bold,
        fontSize = 32.sp,
        lineHeight = 40.sp
    ),
    displayMedium = TextStyle(
        fontWeight = FontWeight.SemiBold,
        fontSize = 28.sp,
        lineHeight = 36.sp
    ),
    headlineLarge = TextStyle(
        fontWeight = FontWeight.Bold,
        fontSize = 24.sp,
        lineHeight = 32.sp
    ),
    headlineMedium = TextStyle(
        fontWeight = FontWeight.Medium,
        fontSize = 20.sp,
        lineHeight = 28.sp
    ),
    titleLarge = TextStyle(
        fontWeight = FontWeight.Medium,
        fontSize = 18.sp,
        lineHeight = 24.sp
    ),
    titleMedium = TextStyle(
        fontWeight = FontWeight.Medium,
        fontSize = 16.sp,
        lineHeight = 24.sp
    ),
    bodyLarge = TextStyle(
        fontWeight = FontWeight.Normal,
        fontSize = 16.sp,
        lineHeight = 24.sp
    ),
    bodyMedium = TextStyle(
        fontWeight = FontWeight.Normal,
        fontSize = 14.sp,
        lineHeight = 20.sp
    ),
    bodySmall = TextStyle(
        fontWeight = FontWeight.Normal,
        fontSize = 12.sp,
        lineHeight = 16.sp
    ),
    labelLarge = TextStyle(
        fontWeight = FontWeight.Medium,
        fontSize = 14.sp,
        lineHeight = 20.sp
    ),
    labelMedium = TextStyle(
        fontWeight = FontWeight.Medium,
        fontSize = 12.sp,
        lineHeight = 16.sp
    ),
    labelSmall = TextStyle(
        fontWeight = FontWeight.Medium,
        fontSize = 11.sp,
        lineHeight = 16.sp
    )
)

val PPOBShapes = Shapes(
    small = RoundedCornerShape(4.dp),
    medium = RoundedCornerShape(8.dp),
    large = RoundedCornerShape(16.dp),
    extraLarge = RoundedCornerShape(24.dp)
)

// ============== Theme ==============

@Composable
fun PPOBTheme(content: @Composable () -> Unit) {
    MaterialTheme(
        colorScheme = PPOBColorScheme,
        typography = PPOBTypography,
        shapes = PPOBShapes,
        content = content
    )
}