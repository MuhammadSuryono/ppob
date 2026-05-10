package com.yonotech.ppob.presentation.transactiondetail

import android.util.Log
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.yonotech.ppob.domain.usecase.GetTransactionStatusUseCase
import com.yonotech.ppob.domain.usecase.GetTransactionHistoryUseCase
import com.yonotech.ppob.data.remote.model.TransactionDetailResponse
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

data class TransactionDetailUiState(
    val transactionId: String = "",
    val productName: String = "",
    val customerNumber: String = "",
    val status: String = "",
    val sellingPrice: Double = 0.0,
    val platformPrice: Double = 0.0,
    val marginAmount: Double? = null,
    val commissionAmount: Double? = null,
    val createdAt: String = "",
    val refId: String? = null,
    val serialNumber: String? = null,
    val message: String? = null,
    val isLoading: Boolean = false,
    val error: String? = null
)

@HiltViewModel
class TransactionDetailViewModel @Inject constructor(
    private val getTransactionStatusUseCase: GetTransactionStatusUseCase
) : ViewModel() {

    private val _uiState = MutableStateFlow(TransactionDetailUiState())
    val uiState: StateFlow<TransactionDetailUiState> = _uiState.asStateFlow()

    fun loadTransactionDetail(authToken: String, transactionId: String) {
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(
                transactionId = transactionId,
                isLoading = true,
                error = null
            )

            try {
                val result = getTransactionStatusUseCase(authToken, transactionId)
                if (result.isSuccess) {
                    val response = result.getOrNull()
                    if (response != null && response.success == true && response.data != null) {
                        with(response.data) {
                            _uiState.value = _uiState.value.copy(
                                productName = details?.productName ?: productName.orEmpty(),
                                customerNumber = details?.customerNumber ?: customerNumber.orEmpty(),
                                status = status.orEmpty(),
                                sellingPrice = details?.sellingPrice ?: 0.0,
                                platformPrice = details?.amount ?: 0.0,
                                marginAmount = details?.marginAmount,
                                commissionAmount = details?.commissionAmount,
                                createdAt = details?.createdAt ?: "",
                                refId = refId,
                                serialNumber = serialNumber,
                                message = message,
                                isLoading = false
                            )
                        }
                    } else {
                        _uiState.value = _uiState.value.copy(
                            isLoading = false,
                            error = response?.message ?: "Data tidak ditemukan"
                        )
                    }
                } else {
                    _uiState.value = _uiState.value.copy(
                        isLoading = false,
                        error = result.exceptionOrNull()?.message ?: "Gagal memuat detail"
                    )
                }
            } catch (e: Exception) {
                Log.e("TxDetailVM", "Error loading transaction detail", e)
                _uiState.value = _uiState.value.copy(
                    isLoading = false,
                    error = e.message ?: "Gagal memuat detail transaksi"
                )
            }
        }
    }

    fun clearError() {
        _uiState.value = _uiState.value.copy(error = null)
    }
}

// For TransactionHistory listing
class TransactionHistoryViewModel {
    // TODO: Implement with TransactionRepository for cached transactions
    val transactions = mutableListOf<com.yonotech.ppob.domain.model.Transaction>()
}