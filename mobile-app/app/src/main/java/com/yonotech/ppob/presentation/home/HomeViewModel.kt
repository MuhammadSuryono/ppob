package com.yonotech.ppob.presentation.home

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.yonotech.ppob.data.remote.model.ApiResponse
import com.yonotech.ppob.data.remote.model.Product
import com.yonotech.ppob.domain.model.Category
import com.yonotech.ppob.domain.model.Product as DomainProduct
import com.yonotech.ppob.domain.usecase.GetCategoriesUseCase
import com.yonotech.ppob.domain.usecase.GetProductsUseCase
import com.yonotech.ppob.domain.usecase.GetBalanceUseCase
import com.yonotech.ppob.domain.repository.AuthRepository
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.*
import kotlinx.coroutines.launch
import javax.inject.Inject

data class HomeUiState(
    val categories: List<Category> = emptyList(),
    val featuredProducts: List<DomainProduct> = emptyList(),
    val balance: Double = 0.0,
    val balanceHeld: Double = 0.0,
    val recentTransactions: List<com.yonotech.ppob.domain.model.Transaction> = emptyList(),
    val isLoading: Boolean = true,
    val error: String? = null,
    val unreadNotifications: Int = 0
)

@HiltViewModel
class HomeViewModel @Inject constructor(
    private val getCategoriesUseCase: GetCategoriesUseCase,
    private val getProductsUseCase: GetProductsUseCase,
    private val getBalanceUseCase: GetBalanceUseCase,
    private val authRepository: AuthRepository
) : ViewModel() {

    private val _uiState = MutableStateFlow(HomeUiState())
    val uiState: StateFlow<HomeUiState> = _uiState.asStateFlow()

    init {
        loadHomeData()
    }

    private fun loadHomeData() {
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isLoading = true)
            try {
                // Load categories
                val categoriesResult = getCategoriesUseCase()
                if (categoriesResult.isSuccess) {
                    _uiState.value = _uiState.value.copy(
                        categories = categoriesResult.getOrDefault(emptyList())
                    )
                }

                // Load balance
                val token = authRepository.getStoredToken()
                if (token != null) {
                    val balanceResult = getBalanceUseCase(token)
                    if (balanceResult.isSuccess) {
                        balanceResult.getOrNull()?.let { walletBalance ->
                            _uiState.value = _uiState.value.copy(
                                balance = walletBalance.balanceAvailable,
                                balanceHeld = walletBalance.balanceHeld
                            )
                        }
                    }
                }

                _uiState.value = _uiState.value.copy(isLoading = false)
            } catch (e: Exception) {
                _uiState.value = _uiState.value.copy(
                    isLoading = false,
                    error = e.message ?: "Gagal memuat data"
                )
            }
        }
    }

    fun refresh() {
        loadHomeData()
    }

    fun refreshBalance() {
        viewModelScope.launch {
            val token = authRepository.getStoredToken() ?: return@launch
            try {
                val balanceResult = getBalanceUseCase(token)
                if (balanceResult.isSuccess) {
                    balanceResult.getOrNull()?.let { walletBalance ->
                        _uiState.value = _uiState.value.copy(
                            balance = walletBalance.balanceAvailable,
                            balanceHeld = walletBalance.balanceHeld
                        )
                    }
                }
            } catch (e: Exception) {
                // Use cached balance
            }
        }
    }

    fun clearError() {
        _uiState.value = _uiState.value.copy(error = null)
    }
}

data class CategoryItem(
    val categoryId: String,
    val name: String,
    val iconUrl: String,
    val iconResId: Int? = null
)