package com.yonotech.ppob.mobile.ui.auth

import androidx.compose.animation.AnimatedContent
import androidx.compose.foundation.layout.*
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.ArrowBack
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.platform.LocalContext
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.input.PasswordVisualTransformation
import androidx.compose.ui.text.input.VisualTransformation
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.hilt.navigation.compose.hiltViewModel
import com.yonotech.ppob.mobile.data.remote.dto.AuthResponse
import com.yonotech.ppob.mobile.ui.components.PinDots
import com.yonotech.ppob.mobile.ui.components.PinPad
import com.yonotech.ppob.mobile.ui.components.PpoButton
import com.yonotech.ppob.mobile.ui.components.PpoTextField
import com.yonotech.ppob.mobile.ui.theme.PpoMobileTheme
import com.yonotech.ppob.mobile.utils.DeviceUtils
import com.yonotech.ppob.mobile.utils.Resource
import com.yonotech.ppob.mobile.viewmodels.auth.AuthViewModel

@Composable
fun SetPasswordPinScreen(
    phone: String,
    requestId: String,
    onRegisterSuccess: (Int) -> Unit, // userId
    viewModel: AuthViewModel = hiltViewModel()
) {
    val authState by viewModel.authState.collectAsState()
    val context = LocalContext.current

    LaunchedEffect(authState) {
        if (authState is Resource.Success) {
            val data = (authState as Resource.Success).data
            onRegisterSuccess(data.userId ?: 0)
            viewModel.resetState()
        }
    }

    SetPasswordPinContent(
        phone = phone,
        requestId = requestId,
        authState = authState,
        onRegister = { password, pin ->
            viewModel.register(
                email = "",
                phone = phone,
                fullName = "",
                password = password,
                pin = pin,
                deviceId = DeviceUtils.getDeviceId(context),
                requestId = requestId
            )
        }
    )
}

@Composable
fun SetPasswordPinContent(
    phone: String,
    requestId: String,
    authState: Resource<AuthResponse>,
    onRegister: (String, String) -> Unit,
    modifier: Modifier = Modifier
) {
    var password by remember { mutableStateOf("") }
    var confirmPassword by remember { mutableStateOf("") }
    var pin by remember { mutableStateOf("") }
    var confirmPin by remember { mutableStateOf("") }
    var currentStep by remember { mutableIntStateOf(1) } // 1: PIN, 2: Confirm PIN, 3: Password
    var showPassword by remember { mutableStateOf(false) }

    Scaffold(
        modifier = modifier,
        topBar = {
            if (currentStep > 1) {
                IconButton(onClick = { 
                    if (currentStep == 2) confirmPin = ""
                    currentStep-- 
                }) {
                    Icon(Icons.Default.ArrowBack, contentDescription = "Back")
                }
            }
        }
    ) { padding ->
        Box(
            modifier = Modifier
                .fillMaxSize()
                .padding(padding)
                .padding(horizontal = 24.dp),
            contentAlignment = Alignment.Center
        ) {
            AnimatedContent(targetState = currentStep, label = "RegStep") { step ->
                Column(
                    modifier = Modifier.fillMaxHeight(),
                    horizontalAlignment = Alignment.CenterHorizontally,
                    verticalArrangement = Arrangement.Center
                ) {
                    when (step) {
                        1 -> {
                            Text(
                                text = "Buat PIN Baru",
                                fontSize = 24.sp,
                                fontWeight = FontWeight.Bold,
                                color = MaterialTheme.colorScheme.primary
                            )
                            Text(
                                text = "PIN akan digunakan untuk setiap transaksi",
                                fontSize = 14.sp,
                                color = MaterialTheme.colorScheme.onSurfaceVariant,
                                modifier = Modifier.padding(top = 8.dp, bottom = 32.dp)
                            )

                            PinDots(pin = pin)
                            
                            Spacer(modifier = Modifier.height(48.dp))
                            
                            PinPad(
                                onDigit = { digit ->
                                    if (pin.length < 6) {
                                        pin += digit
                                        if (pin.length == 6) {
                                            currentStep = 2
                                        }
                                    }
                                },
                                onBackspace = { if (pin.isNotEmpty()) pin = pin.dropLast(1) }
                            )
                        }
                        2 -> {
                            Text(
                                text = "Konfirmasi PIN",
                                fontSize = 24.sp,
                                fontWeight = FontWeight.Bold,
                                color = MaterialTheme.colorScheme.primary
                            )
                            Text(
                                text = "Masukkan kembali PIN yang telah Anda buat",
                                fontSize = 14.sp,
                                color = MaterialTheme.colorScheme.onSurfaceVariant,
                                modifier = Modifier.padding(top = 8.dp, bottom = 32.dp)
                            )

                            PinDots(pin = confirmPin)
                            
                            if (confirmPin.length == 6 && confirmPin != pin) {
                                Text(
                                    text = "PIN tidak cocok",
                                    color = MaterialTheme.colorScheme.error,
                                    modifier = Modifier.padding(top = 8.dp)
                                )
                            }
                            
                            Spacer(modifier = Modifier.height(48.dp))
                            
                            PinPad(
                                onDigit = { digit ->
                                    if (confirmPin.length < 6) {
                                        confirmPin += digit
                                        if (confirmPin.length == 6 && confirmPin == pin) {
                                            currentStep = 3
                                        }
                                    }
                                },
                                onBackspace = { if (confirmPin.isNotEmpty()) confirmPin = confirmPin.dropLast(1) }
                            )
                        }
                        3 -> {
                            Text(
                                text = "Buat Password",
                                fontSize = 24.sp,
                                fontWeight = FontWeight.Bold,
                                color = MaterialTheme.colorScheme.primary
                            )
                            Text(
                                text = "Password digunakan untuk masuk ke akun",
                                fontSize = 14.sp,
                                color = MaterialTheme.colorScheme.onSurfaceVariant,
                                modifier = Modifier.padding(top = 8.dp, bottom = 32.dp)
                            )

                            PpoTextField(
                                value = password,
                                onValueChange = { password = it },
                                label = "Password Baru",
                                visualTransformation = if (showPassword) VisualTransformation.None else PasswordVisualTransformation(),
                                trailingIcon = {
                                    TextButton(onClick = { showPassword = !showPassword }) {
                                        Text(if (showPassword) "Tutup" else "Lihat")
                                    }
                                }
                            )
                            
                            Spacer(modifier = Modifier.height(16.dp))

                            PpoTextField(
                                value = confirmPassword,
                                onValueChange = { confirmPassword = it },
                                label = "Konfirmasi Password",
                                visualTransformation = if (showPassword) VisualTransformation.None else PasswordVisualTransformation(),
                                isError = confirmPassword.isNotEmpty() && confirmPassword != password,
                                errorText = "Password tidak cocok"
                            )

                            if (authState is Resource.Error) {
                                Text(
                                    text = (authState as Resource.Error).message,
                                    color = MaterialTheme.colorScheme.error,
                                    modifier = Modifier.padding(vertical = 16.dp)
                                )
                            }

                            Spacer(modifier = Modifier.height(32.dp))

                            PpoButton(
                                label = "Daftar Sekarang",
                                onClick = { 
                                    onRegister(password, pin)
                                },
                                isLoading = authState is Resource.Loading,
                                enabled = password.isNotEmpty() && password == confirmPassword
                            )
                        }
                    }
                }
            }
        }
    }
}

@Preview(showBackground = true)
@Composable
fun SetPasswordPinScreenPreview() {
    PpoMobileTheme {
        SetPasswordPinContent(
            phone = "08123456789",
            requestId = "req-123",
            authState = Resource.Idle,
            onRegister = { _, _ -> }
        )
    }
}
