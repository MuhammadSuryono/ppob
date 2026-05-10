package com.yonotech.ppob.presentation.wallet

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.yonotech.ppob.domain.model.WalletBalance
import com.yonotech.ppob.domain.usecase.GetBalanceUseCase
import com.yonotech.ppob.domain.usecase.TopUpStaffUseCase
import com.yonotech.ppob.domain.repository.WalletRepository
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import java.math.BigDecimal
import javax.inject.Inject

data class WalletUiState(
    val balanceAvailable: Double = 0.0,
    val balanceHeld: Double = 0.0,
    val balanceTotal: Double = 0.0,
    val currency: String = "IDR",
    val isLoading: Boolean = true,
    val error: String? = null,
    val isRefreshing: Boolean = false,
    val lastUpdatedAt: Long = 0L
)

@HiltViewModel
class WalletViewModel @Inject constructor(
    private val getBalanceUseCase: GetBalanceUseCase,
    private val topUpStaffUseCase: TopUpStaffUseCase,
    private val walletRepository: WalletRepository
) : ViewModel() {

    private val _uiState = MutableStateFlow(WalletUiState())
    val uiState: StateFlow<WalletUiState> = _uiState.asStateFlow()

    private var authToken: String? = null

    fun setAuthToken(token: String) {
        authToken = token
        loadBalance()
    }

    fun loadBalance() {
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isLoading = true, error = null)
            try {
                val token = authToken
                if (token == null) {
                    // Try cached
                    val cached = walletRepository.getCachedBalance()
                    if (cached != null) {
                        _uiState.value = WalletUiState(
                            balanceAvailable = cached.balanceAvailable,
                            balanceHeld = cached.balanceHeld,
                            balanceTotal = cached.balanceAvailable + cached.balanceHeld,
                            currency = cached.currency,
                            isLoading = false,
                            lastUpdatedAt = System.currentTimeMillis()
                        )
                        return@launch
                    }
                    _uiState.value = _uiState.value.copy(isLoading = false, error = "Belum login")
                    return@launch
                }

                val result = getBalanceUseCase(token)
                if (result.isSuccess) {
                    val walletBalance = result.getOrNull()
                    if (walletBalance != null) {
                        _uiState.value = WalletUiState(
                            balanceAvailable = walletBalance.balanceAvailable,
                            balanceHeld = walletBalance.balanceHeld,
                            balanceTotal = walletBalance.balanceAvailable + walletBalance.balanceHeld,
                            currency = walletBalance.currency,
                            isLoading = false,
                            lastUpdatedAt = System.currentTimeMillis()
                        )
                    }
                } else {
                    // Fall back to cached balance
                    val cached = walletRepository.getCachedBalance()
                    if (cached != null) {
                        _uiState.value = WalletUiState(
                            balanceAvailable = cached.balanceAvailable,
                            balanceHeld = cached.balanceHeld,
                            balanceTotal = cached.balanceAvailable + cached.balanceHeld,
                            currency = cached.currency,
                            isLoading = false,
                            error = result.exceptionOrNull()?.message,
                            lastUpdatedAt = System.currentTimeMillis()
                        )
                    } else {
                        _uiState.value = _uiState.value.copy(
                            isLoading = false,
                            error = result.exceptionOrNull()?.message ?: "Gagal memuat saldo"
                        )
                    }
                }
            } catch (e: Exception) {
                _uiState.value = _uiState.value.copy(
                    isLoading = false,
                    error = e.message ?: "Gagal memuat saldo"
                )
            }
        }
    }

    fun refreshBalance() {
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isRefreshing = true)
            loadBalance()
            _uiState.value = _uiState.value.copy(isRefreshing = false)
        }
    }

    fun topUpStaff(staffId: String, amount: BigDecimal) {
        viewModelScope.launch {
            val token = authToken ?: return@launch
            _uiState.value = _uiState.value.copy(isLoading = true)

            try {
                val result = topUpStaffUseCase(token, staffId, amount)
                if (result.isSuccess) {
                    // Refresh balance after success
                    loadBalance()
                } else {
                    _uiState.value = _uiState.value.copy(
                        isLoading = false,
                        error = result.exceptionOrNull()?.message ?: "Top up gagal"
                    )
                }
            } catch (e: Exception) {
                _uiState.value = _uiState.value.copy(
                    isLoading = false,
                    error = e.message ?: "Top up gagal"
                )
            }
        }
    }

    fun clearError() {
        _uiState.value = _uiState.value.copy(error = null)
    }
}

data class WalletTransaction(
    val id: String,
    val type: String, // "topup", "debit", "credit", "commission"
    val amount: Double,
    val description: String,
    val balanceAfter: Double,
    val timestamp: Long
)