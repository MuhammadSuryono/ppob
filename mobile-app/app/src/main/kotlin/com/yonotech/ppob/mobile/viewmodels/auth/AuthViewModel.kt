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

    private val _initiateState = MutableStateFlow<Resource<InitiateAuthResponse>>(Resource.Idle)
    val initiateState = _initiateState.asStateFlow()

    private val _sendOtpState = MutableStateFlow<Resource<SendOtpResponse>>(Resource.Idle)
    val sendOtpState = _sendOtpState.asStateFlow()

    private val _verifyOtpState = MutableStateFlow<Resource<VerifyOtpResponse>>(Resource.Idle)
    val verifyOtpState = _verifyOtpState.asStateFlow()

    fun initiateAuth(phone: String, deviceId: String, fingerPrint: String? = null) {
        viewModelScope.launch {
            _initiateState.value = Resource.Loading
            try {
                val response = repository.initiateAuth(InitiateAuthRequest(phone, deviceId, fingerPrint))
                if (response.isSuccessful) {
                    val data = response.body()
                    if (data != null) {
                        _initiateState.value = Resource.Success(data)
                    } else {
                        _initiateState.value = Resource.Error("Response body is empty")
                    }
                } else {
                    _initiateState.value = Resource.Error("Failed to initiate auth: ${response.message()}")
                }
            } catch (e: Exception) {
                _initiateState.value = Resource.Error(e.message ?: "An unknown error occurred")
            }
        }
    }

    fun sendOtp(phone: String, type: String) {
        viewModelScope.launch {
            _sendOtpState.value = Resource.Loading
            try {
                val response = repository.sendOtp(SendOtpRequest(phone, type))
                if (response.isSuccessful) {
                    val data = response.body()
                    if (data != null) {
                        _sendOtpState.value = Resource.Success(data)
                    } else {
                        _sendOtpState.value = Resource.Error("Response body is empty")
                    }
                } else {
                    _sendOtpState.value = Resource.Error("Failed to send OTP: ${response.message()}")
                }
            } catch (e: Exception) {
                _sendOtpState.value = Resource.Error(e.message ?: "An unknown error occurred")
            }
        }
    }

    fun login(identifier: String, password: String, deviceId: String) {
        viewModelScope.launch {
            _authState.value = Resource.Loading
            try {
                val response = repository.login(LoginRequest(phone = identifier, password = password, deviceId = deviceId))
                if (response.isSuccessful) {
                    val authData = response.body()
                    if (authData != null) {
                        if (authData.accessToken != null) {
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

    fun verifyPin(phone: String, pin: String, deviceId: String) {
        viewModelScope.launch {
            _authState.value = Resource.Loading
            try {
                val response = repository.verifyPin(LoginRequest(phone = phone, pin = pin, deviceId = deviceId))
                if (response.isSuccessful) {
                    val authData = response.body()
                    if (authData != null) {
                        if (authData.accessToken != null) {
                            tokenManager.saveTokens(authData.accessToken, authData.refreshToken)
                        }
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

    fun register(email: String, phone: String, fullName: String, password: String, pin: String, deviceId: String, requestId: String) {
        viewModelScope.launch {
            _authState.value = Resource.Loading
            try {
                val response = repository.register(RegisterRequest(email, phone, phone, password, pin, deviceId, requestId))
                if (response.isSuccessful) {
                    val authData = response.body()
                    if (authData != null) {
                        if (authData.accessToken != null) {
                            tokenManager.saveTokens(authData.accessToken, authData.refreshToken)
                        }
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

    fun verifyOtp(requestId: String, phone: String, code: String, type: String) {
        viewModelScope.launch {
            _verifyOtpState.value = Resource.Loading
            try {
                val response = repository.verifyOtp(VerifyOtpRequest(requestId, phone, code, type))
                if (response.isSuccessful) {
                    val data = response.body()
                    if (data != null && data.isVerified) {
                        _verifyOtpState.value = Resource.Success(data)
                    } else {
                        _verifyOtpState.value = Resource.Error("OTP verification failed")
                    }
                } else {
                    _verifyOtpState.value = Resource.Error("OTP verification failed: ${response.message()}")
                }
            } catch (e: Exception) {
                _verifyOtpState.value = Resource.Error(e.message ?: "An unknown error occurred")
            }
        }
    }

    fun verifyPassword(phone: String, password: String, deviceId: String, requestId: String) {
        viewModelScope.launch {
            _authState.value = Resource.Loading
            try {
                val response = repository.verifyPassword(VerifyPasswordRequest(phone, password, deviceId, requestId))
                if (response.isSuccessful) {
                    val authData = response.body()
                    if (authData != null) {
                        if (authData.accessToken != null) {
                            tokenManager.saveTokens(authData.accessToken, authData.refreshToken)
                        }
                        _authState.value = Resource.Success(authData)
                    } else {
                        _authState.value = Resource.Error("Response body is empty")
                    }
                } else {
                    _authState.value = Resource.Error("Password verification failed: ${response.message()}")
                }
            } catch (e: Exception) {
                _authState.value = Resource.Error(e.message ?: "An unknown error occurred")
            }
        }
    }

    fun resetState() {
        _authState.value = Resource.Idle
        _initiateState.value = Resource.Idle
        _sendOtpState.value = Resource.Idle
        _verifyOtpState.value = Resource.Idle
    }
}