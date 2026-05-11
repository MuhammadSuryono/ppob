package com.yonotech.ppob.mobile.viewmodels.auth

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.yonotech.ppob.mobile.data.local.TokenManager
import com.yonotech.ppob.mobile.data.remote.dto.*
import com.yonotech.ppob.mobile.data.repository.AuthRepository
import com.yonotech.ppob.mobile.utils.Resource
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

@HiltViewModel
class AuthViewModel @Inject constructor(
    private val repository: AuthRepository,
    private val tokenManager: TokenManager
) : ViewModel() {

    private val _authState = MutableStateFlow<Resource<AuthResponse>>(Resource.Idle)
    val authState = _authState.asStateFlow()

    fun sendOtp(phone: String, deviceId: String) {
        viewModelScope.launch {
            _authState.value = Resource.Loading
            try {
                val response = repository.sendOtp(SendOtpRequest(phone, deviceId))
                if (response.isSuccessful) {
                    val authData = response.body()
                    if (authData != null) {
                        _authState.value = Resource.Success(authData)
                    } else {
                        _authState.value = Resource.Error("Response body is empty")
                    }
                } else {
                    _authState.value = Resource.Error("Failed to send OTP: ${response.message()}")
                }
            } catch (e: Exception) {
                _authState.value = Resource.Error(e.message ?: "An unknown error occurred")
            }
        }
    }

    fun login(identifier: String, password: String, deviceId: String) {
        viewModelScope.launch {
            _authState.value = Resource.Loading
            try {
                val response = repository.login(LoginRequest(identifier, password, deviceId))
                if (response.isSuccessful) {
                    val authData = response.body()
                    if (authData != null) {
                        if (!authData.requiresOtp && authData.accessToken != null && authData.refreshToken != null) {
                            tokenManager.saveTokens(authData.accessToken, authData.refreshToken)
                        }
                        _authState.value = Resource.Success(authData)
                    } else {
                        _authState.value = Resource.Error("Response body is empty")
                    }
                } else {
                    _authState.value = Resource.Error("Login failed: ${response.message()}")
                }
            } catch (e: Exception) {
                _authState.value = Resource.Error(e.message ?: "An unknown error occurred")
            }
        }
    }

    fun pinLogin(phone: String, pin: String, deviceId: String) {
        viewModelScope.launch {
            _authState.value = Resource.Loading
            try {
                val response = repository.pinLogin(PinLoginRequest(phone, pin, deviceId))
                if (response.isSuccessful) {
                    val authData = response.body()
                    if (authData != null && authData.accessToken != null && authData.refreshToken != null) {
                        tokenManager.saveTokens(authData.accessToken, authData.refreshToken)
                        _authState.value = Resource.Success(authData)
                    } else if (authData != null) {
                        _authState.value = Resource.Success(authData)
                    } else {
                        _authState.value = Resource.Error("Response body is empty")
                    }
                } else {
                    _authState.value = Resource.Error("PIN login failed: ${response.message()}")
                }
            } catch (e: Exception) {
                _authState.value = Resource.Error(e.message ?: "An unknown error occurred")
            }
        }
    }

    fun register(email: String, phone: String, name: String, password: String, pin: String) {
        viewModelScope.launch {
            _authState.value = Resource.Loading
            try {
                val response = repository.register(RegisterRequest(email, phone, name, password, pin))
                if (response.isSuccessful) {
                    val authData = response.body()
                    if (authData != null) {
                        _authState.value = Resource.Success(authData)
                    } else {
                        _authState.value = Resource.Error("Response body is empty")
                    }
                } else {
                    _authState.value = Resource.Error("Registration failed: ${response.message()}")
                }
            } catch (e: Exception) {
                _authState.value = Resource.Error(e.message ?: "An unknown error occurred")
            }
        }
    }

    fun verifyOtp(identifier: String, otpCode: String, type: String) {
        viewModelScope.launch {
            _authState.value = Resource.Loading
            try {
                val response = repository.verifyOtp(VerifyOtpRequest(identifier, otpCode, type))
                if (response.isSuccessful) {
                    val authData = response.body()
                    if (authData != null) {
                        if (authData.accessToken != null && authData.refreshToken != null) {
                            tokenManager.saveTokens(authData.accessToken, authData.refreshToken)
                        }
                        _authState.value = Resource.Success(authData)
                    } else {
                        _authState.value = Resource.Error("Response body is empty")
                    }
                } else {
                    _authState.value = Resource.Error("OTP verification failed: ${response.message()}")
                }
            } catch (e: Exception) {
                _authState.value = Resource.Error(e.message ?: "An unknown error occurred")
            }
        }
    }

    fun setPasswordPin(phone: String, password: String, pin: String) {
        viewModelScope.launch {
            _authState.value = Resource.Loading
            try {
                val response = repository.setPasswordPin(SetPasswordPinRequest(phone, password, pin))
                if (response.isSuccessful) {
                    val authData = response.body()
                    if (authData != null) {
                        if (authData.accessToken != null && authData.refreshToken != null) {
                            tokenManager.saveTokens(authData.accessToken, authData.refreshToken)
                        }
                        _authState.value = Resource.Success(authData)
                    } else {
                        _authState.value = Resource.Error("Response body is empty")
                    }
                } else {
                    _authState.value = Resource.Error("Failed to set password: ${response.message()}")
                }
            } catch (e: Exception) {
                _authState.value = Resource.Error(e.message ?: "An unknown error occurred")
            }
        }
    }

    fun resetState() {
        _authState.value = Resource.Idle
    }
}