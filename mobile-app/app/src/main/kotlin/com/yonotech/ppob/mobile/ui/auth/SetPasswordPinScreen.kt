package com.yonotech.ppob.mobile.ui.auth

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.text.KeyboardOptions
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.input.KeyboardType
import androidx.compose.ui.text.input.PasswordVisualTransformation
import androidx.compose.ui.text.input.VisualTransformation
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import androidx.hilt.navigation.compose.hiltViewModel
import com.yonotech.ppob.mobile.ui.components.PpoButton
import com.yonotech.ppob.mobile.ui.components.PpoTextField
import com.yonotech.ppob.mobile.utils.Resource
import com.yonotech.ppob.mobile.viewmodels.auth.AuthViewModel

@Composable
fun SetPasswordPinScreen(
    phone: String,
    onAuthSuccess: () -> Unit,
    viewModel: AuthViewModel = hiltViewModel()
) {
    var password by remember { mutableStateOf("") }
    var confirmPassword by remember { mutableStateOf("") }
    var pin by remember { mutableStateOf("") }
    var showPassword by remember { mutableStateOf(false) }
    var showConfirmPassword by remember { mutableStateOf(false) }
    
    val authState by viewModel.authState.collectAsState()

    LaunchedEffect(authState) {
        if (authState is Resource.Success) {
            onAuthSuccess()
            viewModel.resetState()
        }
    }

    Column(
        modifier = Modifier
            .fillMaxSize()
            .padding(24.dp),
        horizontalAlignment = Alignment.CenterHorizontally
    ) {
        Spacer(modifier = Modifier.height(32.dp))
        Text(
            text = "Buat Akun",
            fontSize = 28.sp,
            fontWeight = FontWeight.Bold,
            color = MaterialTheme.colorScheme.primary
        )
        Text(
            text = "Buat password dan PIN untuk melanjutkan",
            fontSize = 16.sp,
            color = Color.Gray,
            modifier = Modifier.padding(top = 8.dp, bottom = 32.dp)
        )

        PpoTextField(
            value = password,
            onValueChange = { password = it },
            label = "Password",
            visualTransformation = if (showPassword) VisualTransformation.None else PasswordVisualTransformation(),
            keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Password),
            trailingIcon = {
                TextButton(onClick = { showPassword = !showPassword }) {
                    Text(if (showPassword) "Sembunyikan" else "Lihat")
                }
            }
        )

        Spacer(modifier = Modifier.height(16.dp))

        PpoTextField(
            value = confirmPassword,
            onValueChange = { confirmPassword = it },
            label = "Konfirmasi Password",
            visualTransformation = if (showConfirmPassword) VisualTransformation.None else PasswordVisualTransformation(),
            keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.Password),
            trailingIcon = {
                TextButton(onClick = { showConfirmPassword = !showConfirmPassword }) {
                    Text(if (showConfirmPassword) "Sembunyikan" else "Lihat")
                }
            }
        )

        Spacer(modifier = Modifier.height(16.dp))

        PpoTextField(
            value = pin,
            onValueChange = { if (it.length <= 6 && it.all { char -> char.isDigit() }) pin = it },
            label = "6-Digit PIN",
            keyboardOptions = KeyboardOptions(keyboardType = KeyboardType.NumberPassword),
            placeholder = "123456"
        )

        Spacer(modifier = Modifier.height(32.dp))

        if (authState is Resource.Error) {
            Text(
                text = (authState as Resource.Error).message,
                color = MaterialTheme.colorScheme.error,
                modifier = Modifier.padding(bottom = 16.dp)
            )
        }

        PpoButton(
            label = "Buat Akun",
            onClick = { 
                viewModel.setPasswordPin(phone, password, pin) 
            },
            isLoading = authState is Resource.Loading,
            enabled = password.isNotEmpty() && 
                     password == confirmPassword && 
                     pin.length == 6
        )

        Spacer(modifier = Modifier.height(32.dp))
    }
}