package com.yonotech.ppob.mobile.ui.auth

import androidx.compose.foundation.border
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.foundation.text.BasicTextField
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.SolidColor
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.text.style.TextAlign
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.hilt.navigation.compose.hiltViewModel
import com.yonotech.ppob.mobile.data.remote.dto.VerifyOtpResponse
import com.yonotech.ppob.mobile.ui.theme.PpoMobileTheme
import com.yonotech.ppob.mobile.utils.Resource
import com.yonotech.ppob.mobile.viewmodels.auth.AuthViewModel

@Composable
fun OtpScreen(
    requestId: String,
    phone: String,
    type: String, // "register" or "login"
    onOtpSuccess: (String) -> Unit,
    viewModel: AuthViewModel = hiltViewModel()
) {
    var otpCode by remember { mutableStateOf("") }
    val verifyOtpState by viewModel.verifyOtpState.collectAsState()

    LaunchedEffect(verifyOtpState) {
        if (verifyOtpState is Resource.Success) {
            val success = verifyOtpState as Resource.Success<VerifyOtpResponse>
            if (success.data.isVerified) {
                onOtpSuccess(phone)
            }
            viewModel.resetState()
        }
    }

    OtpContent(
        phone = phone,
        otpCode = otpCode,
        onOtpCodeChange = { if (it.length <= 6) otpCode = it },
        verifyOtpState = verifyOtpState,
        onVerifyClick = { viewModel.verifyOtp(requestId, phone, otpCode, type) },
        onResendClick = { /* Implement Resend OTP */ },
        modifier = Modifier.fillMaxSize()
    )
}

@Composable
fun OtpContent(
    phone: String,
    otpCode: String,
    onOtpCodeChange: (String) -> Unit,
    verifyOtpState: Resource<VerifyOtpResponse>,
    onVerifyClick: () -> Unit,
    onResendClick: () -> Unit,
    modifier: Modifier = Modifier
) {
    // Automatically trigger verification when 6 digits are entered
    LaunchedEffect(otpCode) {
        if (otpCode.length == 6) {
            onVerifyClick()
        }
    }

    Column(
        modifier = modifier.padding(24.dp),
        horizontalAlignment = Alignment.CenterHorizontally,
        verticalArrangement = Arrangement.Center
    ) {
        Text(
            text = "Verifikasi OTP",
            fontSize = 28.sp,
            fontWeight = FontWeight.Bold,
            color = MaterialTheme.colorScheme.primary
        )
        Text(
            text = "Masukkan kode OTP yang dikirim ke $phone",
            fontSize = 16.sp,
            color = MaterialTheme.colorScheme.onSurfaceVariant,
            modifier = Modifier.padding(top = 8.dp, bottom = 32.dp),
            textAlign = TextAlign.Center
        )

        OtpInputField(
            otpCode = otpCode,
            onOtpCodeChange = onOtpCodeChange
        )

        Spacer(modifier = Modifier.height(32.dp))

        if (verifyOtpState is Resource.Loading) {
            CircularProgressIndicator()
            Spacer(modifier = Modifier.height(32.dp))
        }

        if (verifyOtpState is Resource.Error) {
            Text(
                text = (verifyOtpState as Resource.Error).message,
                color = MaterialTheme.colorScheme.error,
                modifier = Modifier.padding(bottom = 16.dp),
                textAlign = TextAlign.Center
            )
        }

        TextButton(onClick = onResendClick) {
            Text(text = "Tidak menerima kode? Kirim ulang")
        }
    }
}

@Composable
fun OtpInputField(
    otpCode: String,
    onOtpCodeChange: (String) -> Unit
) {
    BasicTextField(
        value = otpCode,
        onValueChange = {
            if (it.length <= 6 && it.all { char -> char.isDigit() }) {
                onOtpCodeChange(it)
            }
        },
        keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Number),
        cursorBrush = SolidColor(Color.Transparent),
        decorationBox = {
            Row(
                horizontalArrangement = Arrangement.spacedBy(8.dp),
                verticalAlignment = Alignment.CenterVertically
            ) {
                repeat(6) { index ->
                    val char = otpCode.getOrNull(index)?.toString() ?: ""
                    val isFocused = otpCode.length == index
                    Box(
                        modifier = Modifier
                            .size(width = 45.dp, height = 55.dp)
                            .border(
                                width = if (isFocused) 2.dp else 1.dp,
                                color = if (isFocused) MaterialTheme.colorScheme.primary else MaterialTheme.colorScheme.outline,
                                shape = RoundedCornerShape(8.dp)
                            ),
                        contentAlignment = Alignment.Center
                    ) {
                        Text(
                            text = char,
                            fontSize = 20.sp,
                            fontWeight = FontWeight.Bold,
                            textAlign = TextAlign.Center,
                            color = MaterialTheme.colorScheme.onSurface
                        )
                    }
                }
            }
        }
    )
}

@Preview(showBackground = true)
@Composable
fun OtpScreenPreview() {
    PpoMobileTheme {
        OtpContent(
            phone = "081234567890",
            otpCode = "123456",
            onOtpCodeChange = {},
            verifyOtpState = Resource.Idle,
            onVerifyClick = {},
            onResendClick = {}
        )
    }
}
