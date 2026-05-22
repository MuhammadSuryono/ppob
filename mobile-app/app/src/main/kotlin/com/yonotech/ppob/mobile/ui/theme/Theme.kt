package com.yonotech.ppob.mobile.ui.theme

import android.app.Activity
import android.os.Build
import androidx.compose.foundation.isSystemInDarkTheme
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.darkColorScheme
import androidx.compose.material3.dynamicDarkColorScheme
import androidx.compose.material3.dynamicLightColorScheme
import androidx.compose.material3.lightColorScheme
import androidx.compose.runtime.Composable
import androidx.compose.runtime.SideEffect
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.toArgb
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.platform.LocalView
import androidx.core.view.WindowCompat

private val LightColorScheme = lightColorScheme(
    primary = PrimaryGreen,
    onPrimary = Color.White,
    primaryContainer = Green200,
    onPrimaryContainer = Green900,
    secondary = SecondaryBlue,
    onSecondary = Color.White,
    tertiary = AccentOrange,
    background = BackgroundLight,
    surface = SurfaceWhite,
    onSurface = TextDark,
    onSurfaceVariant = TextMedium,
    outline = TextLight,
    error = ErrorRed,
    onError = Color.White
)

// Dark scheme defined but will be bypassed as per "jangan dark" requirement
private val DarkColorScheme = darkColorScheme(
    primary = Green200,
    onPrimary = Green900,
    secondary = SecondaryBlue,
    background = Color(0xFF121212),
    surface = Color(0xFF1E1E1E),
    onSurface = Color.White
)

@Composable
fun PpoMobileTheme(
    darkTheme: Boolean = isSystemInDarkTheme(),
    // Dynamic color is available on Android 12+
    dynamicColor: Boolean = false, // Set to false to strictly follow design system colors
    content: @Composable () -> Unit
) {
    // Force LightColorScheme as per "jangan dark" requirement
    val colorScheme = if (darkTheme) {
        // Even if system is dark, we follow the "Modern & Bright" design principle
        LightColorScheme 
    } else {
        LightColorScheme
    }
    
    val view = LocalView.current
    if (!view.isInEditMode) {
        SideEffect {
            val window = (view.context as Activity).window
            window.statusBarColor = colorScheme.primary.toArgb()
            // Set status bar icons to light (since we use a green primary) or dark based on needs
            WindowCompat.getInsetsController(window, view).isAppearanceLightStatusBars = false
        }
    }

    MaterialTheme(
        colorScheme = colorScheme,
        typography = Typography,
        content = content
    )
}
