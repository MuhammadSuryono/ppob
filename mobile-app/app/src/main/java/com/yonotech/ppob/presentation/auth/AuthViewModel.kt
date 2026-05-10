package com.yonotech.ppob.presentation.auth

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.yonotech.ppob.domain.repository.AuthRepository
import com.yonotech.ppob.domain.usecase.LoginUseCase
import com.yonotech.ppob.domain.usecase.RegisterUseCase
import com.yonotech.ppob.domain.usecase.VerifyOtpUseCase
import com.yonotech.ppob.data.remote.model.DeviceFingerprint
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

data class AuthUiState(
    val phoneNumber: String = "",
    val otpCode: String = "",
    val password: String = "",
    val confirmPassword: String = "",
    val pin: String = "",
    val confirmPin: String = "",
    val isLoading: Boolean = false,
    val error: String? = null,
    val isPhoneValid: Boolean = true,
    val isOtpValid: Boolean = true,
    val isPasswordValid: Boolean = true,
    val isPinValid: Boolean = true,
    val currentStep: AuthStep = AuthStep.PHONE_INPUT,
    val userData: User? = null
)

enum class AuthStep {
    PHONE_INPUT, OTP_VERIFY, SET_CREDENTIALS, PIN_LOGIN, COMPLETE
}

data class User(
    val userId: String,
    val phoneNumber: String,
    val name: String,
    val roles: List<String>,
    val activeRole: String,
    val walletId: String
)

@HiltViewModel
class AuthViewModel @Inject constructor(
    private val authRepository: AuthRepository,
    private val registerUseCase: RegisterUseCase,
    private val verifyOtpUseCase: VerifyOtpUseCase,
    private val loginUseCase: LoginUseCase
) : ViewModel() {

    private val _uiState = MutableStateFlow(AuthUiState())
    val uiState: StateFlow<AuthUiState> = _uiState

    fun onPhoneNumberChange(phone: String) {
        _uiState.value = _uiState.value.copy(
            phoneNumber = phone,
            isPhoneValid = validatePhone(phone),
            error = null
        )
    }

    fun onOtpChange(otp: String) {
        _uiState.value = _uiState.value.copy(
            otpCode = otp,
            isOtpValid = otp.length == 6,
            error = null
        )
    }

    fun onPasswordChange(password: String) {
        _uiState.value = _uiState.value.copy(
            password = password,
            isPasswordValid = validatePassword(password),
            error = null
        )
    }

    fun onConfirmPasswordChange(password: String) {
        _uiState.value = _uiState.value.copy(
            confirmPassword = password,
            error = if (_uiState.value.password != password) "Password tidak cocok" else null
        )
    }

    fun onPinChange(pin: String) {
        _uiState.value = _uiState.value.copy(
            pin = pin,
            isPinValid = validatePin(pin),
            error = null
        )
    }

    fun onConfirmPinChange(pin: String) {
        _uiState.value = _uiState.value.copy(
            confirmPin = pin,
            error = if (_uiState.value.pin != pin) "PIN tidak cocok" else null
        )
    }

    fun sendOtp() {
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isLoading = true, error = null)
            try {
                val result = registerUseCase(_uiState.value.phoneNumber, "User")
                result.onSuccess {
                    _uiState.value = _uiState.value.copy(
                        currentStep = AuthStep.OTP_VERIFY,
                        isLoading = false
                    )
                }.onFailure {
                    _uiState.value = _uiState.value.copy(
                        error = it.message ?: "Gagal mengirim OTP",
                        isLoading = false
                    )
                }
            } catch (e: Exception) {
                _uiState.value = _uiState.value.copy(
                    error = e.message ?: "Gagal mengirim OTP",
                    isLoading = false
                )
            }
        }
    }

    fun verifyOtp() {
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isLoading = true, error = null)
            try {
                val result = verifyOtpUseCase(
                    _uiState.value.phoneNumber,
                    _uiState.value.otpCode,
                    _uiState.value.password,
                    _uiState.value.pin
                )
                result.onSuccess {
                    _uiState.value = _uiState.value.copy(
                        currentStep = AuthStep.COMPLETE,
                        isLoading = false
                    )
                }.onFailure {
                    _uiState.value = _uiState.value.copy(
                        error = it.message ?: "Verifikasi OTP gagal",
                        isLoading = false
                    )
                }
            } catch (e: Exception) {
                _uiState.value = _uiState.value.copy(
                    error = e.message ?: "Verifikasi OTP gagal",
                    isLoading = false
                )
            }
        }
    }

    fun login() {
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isLoading = true, error = null)
            try {
                val deviceFingerprint = DeviceFingerprint(
                    deviceId = getDeviceId(),
                    userAgent = "Android/14 PPOB/1.0.0",
                    appVersion = "1.0.0",
                    installTs = System.currentTimeMillis(),
                    lastLoginTs = null
                )

                val result = loginUseCase(
                    _uiState.value.phoneNumber,
                    _uiState.value.password,
                    _uiState.value.pin,
                    deviceFingerprint
                )

                result.onSuccess {
                    _uiState.value = _uiState.value.copy(
                        currentStep = AuthStep.COMPLETE,
                        isLoading = false
                    )
                }.onFailure {
                    _uiState.value = _uiState.value.copy(
                        error = it.message ?: "Login gagal",
                        isLoading = false
                    )
                }
            } catch (e: Exception) {
                _uiState.value = _uiState.value.copy(
                    error = e.message ?: "Login gagal",
                    isLoading = false
                )
            }
        }
    }

    fun resendOtp() {
        sendOtp() // Re-send OTP
    }

    fun resetForm() {
        _uiState.value = AuthUiState()
    }

    fun moveToStep(step: AuthStep) {
        _uiState.value = _uiState.value.copy(currentStep = step, error = null)
    }

    private fun validatePhone(phone: String): Boolean {
        return phone.startsWith("+62") && phone.length >= 10
    }

    private fun validatePassword(password: String): Boolean {
        if (password.length < 8) return false
        if (!password.any { it.isUpperCase() }) return false
        if (!password.any { it.isLowerCase() }) return false
        if (!password.any { it.isDigit() }) return false
        return true
    }

    private fun validatePin(pin: String): Boolean {
        if (pin.length != 6) return false
        if (!pin.all { it.isDigit() }) return false
        if (pin.all { it == pin[0] }) return false // Not all same
        if (pin == "123456" || pin == "654321") return false // Not sequential
        return true
    }

    private fun getDeviceId(): String {
        return java.util.UUID.randomUUID().toString()
    }
}