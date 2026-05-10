package com.yonotech.ppob.presentation.profile

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.yonotech.ppob.domain.repository.AuthRepository
import com.yonotech.ppob.domain.repository.UserRepository
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

data class ProfileUiState(
    val userName: String = "",
    val phoneNumber: String = "",
    val activeRole: String = "",
    val walletBalance: Double = 0.0,
    val balanceHeld: Double = 0.0,
    val isLoading: Boolean = false,
    val error: String? = null,
    val isLoggedIn: Boolean = false
)

@HiltViewModel
class ProfileViewModel @Inject constructor(
    private val authRepository: AuthRepository,
    private val userRepository: UserRepository
) : ViewModel() {

    private val _uiState = MutableStateFlow(ProfileUiState())
    val uiState: StateFlow<ProfileUiState> = _uiState.asStateFlow()

    private var authToken: String? = null

    fun setAuthToken(token: String) {
        authToken = token
        loadProfile()
    }

    fun loadProfile() {
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isLoading = true, error = null)

            try {
                val token = authToken ?: run {
                    _uiState.value = ProfileUiState(isLoading = false, error = "Token tidak tersedia")
                    return@launch
                }

                val result = userRepository.getProfile(token)

                if (result.isSuccess) {
                    val profile = result.getOrNull()
                    if (profile != null) {
                        _uiState.value = ProfileUiState(
                            userName = profile.name,
                            phoneNumber = profile.phoneNumber,
                            activeRole = profile.activeRole.roleName,
                            walletBalance = profile.wallet.balanceAvailable,
                            balanceHeld = profile.wallet.balanceHeld,
                            isLoggedIn = true,
                            isLoading = false
                        )
                    }
                } else {
                    _uiState.value = ProfileUiState(
                        isLoading = false,
                        error = result.exceptionOrNull()?.message ?: "Gagal memuat profil"
                    )
                }
            } catch (e: Exception) {
                _uiState.value = ProfileUiState(
                    isLoading = false,
                    error = e.message ?: "Gagal memuat profil"
                )
            }
        }
    }

    fun logout() {
        viewModelScope.launch {
            val token = authToken
            if (token != null) {
                authRepository.logout(token)
            }
            authRepository.clearTokens()
            _uiState.value = ProfileUiState()
        }
    }

    fun setBiometricEnabled(enabled: Boolean) {
        viewModelScope.launch {
            // TODO: Call user settings API or update local preference
        }
    }

    fun clearError() {
        _uiState.value = _uiState.value.copy(error = null)
    }
}

data class SettingsUiState(
    val biometricEnabled: Boolean = false,
    val trustedDevice: Boolean = false,
    val isLoading: Boolean = false,
    val error: String? = null
)

class SettingsViewModel {
    val uiState = MutableStateFlow(SettingsUiState())

    fun loadSettings() {
        // Load from DataStore preferences
    }

    fun toggleBiometric(enabled: Boolean) {
        uiState.value = uiState.value.copy(biometricEnabled = enabled)
    }

    fun toggleTrustedDevice(trusted: Boolean) {
        uiState.value = uiState.value.copy(trustedDevice = trusted)
    }
}

data class DeviceUiState(
    val devices: List<DeviceItem> = emptyList(),
    val isLoading: Boolean = false,
    val error: String? = null
)

data class DeviceItem(
    val deviceId: String,
    val userAgent: String,
    val lastSeen: String,
    val isTrusted: Boolean
)

class DeviceManagementViewModel {
    val uiState = MutableStateFlow(DeviceUiState())

    fun loadDevices() {
        // Load trusted devices list
    }

    fun revokeDevice(deviceId: String) {
        // Call API to revoke device trust
    }
}

data class ChangePinUiState(
    val currentPin: String = "",
    val newPin: String = "",
    val confirmPin: String = "",
    val isLoading: Boolean = false,
    val error: String? = null,
    val success: Boolean = false
)

class ChangePinViewModel @Inject constructor(
    private val authRepository: AuthRepository
) : ViewModel() {

    private val _uiState = MutableStateFlow(ChangePinUiState())
    val uiState: StateFlow<ChangePinUiState> = _uiState.asStateFlow()

    private var authToken: String? = null

    fun setAuthToken(token: String) {
        authToken = token
    }

    fun onCurrentPinChange(pin: String) {
        _uiState.value = _uiState.value.copy(currentPin = pin, error = null, success = false)
    }

    fun onNewPinChange(pin: String) {
        _uiState.value = _uiState.value.copy(newPin = pin, error = null, success = false)
    }

    fun onConfirmPinChange(pin: String) {
        _uiState.value = _uiState.value.copy(confirmPin = pin, error = null, success = false)
    }

    fun changePin() {
        val currentPin = _uiState.value.currentPin
        val newPin = _uiState.value.newPin
        val confirmPin = _uiState.value.confirmPin

        when {
            newPin != confirmPin -> {
                _uiState.value = _uiState.value.copy(error = "PIN baru tidak cocok")
                return
            }
            newPin.length != 6 || !newPin.all { it.isDigit() } -> {
                _uiState.value = _uiState.value.copy(error = "PIN harus 6 digit angka")
                return
            }
        }

        val token = authToken ?: run {
            _uiState.value = _uiState.value.copy(error = "Token tidak tersedia")
            return
        }

        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isLoading = true, error = null)

            try {
                val request = com.yonotech.ppob.data.remote.model.ChangePinRequest(
                    currentPin = currentPin,
                    newPin = newPin,
                    confirmPin = confirmPin
                )

                val result = authRepository.changePin(token, request)

                if (result.isSuccess) {
                    _uiState.value = ChangePinUiState(
                        isLoading = false,
                        success = true
                    )
                } else {
                    _uiState.value = _uiState.value.copy(
                        isLoading = false,
                        error = result.exceptionOrNull()?.message ?: "Gagal mengubah PIN"
                    )
                }
            } catch (e: Exception) {
                _uiState.value = _uiState.value.copy(
                    isLoading = false,
                    error = e.message ?: "Gagal mengubah PIN"
                )
            }
        }
    }

    fun resetState() {
        _uiState.value = ChangePinUiState()
    }
}