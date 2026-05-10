package com.yonotech.ppob.presentation.transaction

import android.util.Log
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.yonotech.ppob.data.remote.model.InitiateTransactionRequest
import com.yonotech.ppob.data.remote.model.TransactionResponse
import com.yonotech.ppob.domain.usecase.InitiateTransactionUseCase
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

sealed class TransactionState {
    object Idle : TransactionState()
    object Loading : TransactionState()
    data class Success(val transactionId: String, val status: String, val message: String?) : TransactionState()
    data class Pending(val transactionId: String, val message: String?) : TransactionState()
    data class Error(val message: String) : TransactionState()
    data class InsufficientBalance(val required: Double, val current: Double) : TransactionState()
}

@HiltViewModel
class TransactionInitViewModel @Inject constructor(
    private val initiateTransactionUseCase: InitiateTransactionUseCase
) : ViewModel() {

    private val _transactionState = MutableStateFlow<TransactionState>(TransactionState.Idle)
    val transactionState: StateFlow<TransactionState> = _transactionState.asStateFlow()

    private val _customerNumber = MutableStateFlow("")
    val customerNumber = _customerNumber.asStateFlow()

    private val _sellingPrice = MutableStateFlow(0.0)
    val sellingPrice = _sellingPrice.asStateFlow()

    private var currentProductId: String = ""

    fun setProduct(productId: String, sellingPrice: Double) {
        currentProductId = productId
        _sellingPrice.value = sellingPrice
    }

    fun setCustomerNumber(number: String) {
        _customerNumber.value = number
    }

    fun setSellingPrice(price: Double) {
        _sellingPrice.value = price
    }

    fun initiateTransaction(pin: String, authToken: String) {
        if (currentProductId.isEmpty()) {
            _transactionState.value = TransactionState.Error("Produk tidak valid")
            return
        }
        if (_customerNumber.value.isEmpty()) {
            _transactionState.value = TransactionState.Error("Nomor pelanggan harus diisi")
            return
        }

        viewModelScope.launch {
            _transactionState.value = TransactionState.Loading

            try {
                val request = InitiateTransactionRequest(
                    productId = currentProductId,
                    customerNo = _customerNumber.value,
                    sellingPrice = _sellingPrice.value,
                    pin = pin
                )

                val result = initiateTransactionUseCase(authToken, request)

                if (result.isSuccess) {
                    val response = result.getOrNull()
                    if (response != null) {
                        when (response.status) {
                            "Success" -> {
                                _transactionState.value = TransactionState.Success(
                                    transactionId = response.transactionId,
                                    status = response.status,
                                    message = response.message
                                )
                            }
                            "Pending" -> {
                                _transactionState.value = TransactionState.Pending(
                                    transactionId = response.transactionId,
                                    message = response.message
                                )
                            }
                            else -> {
                                _transactionState.value = TransactionState.Error(
                                    response.message ?: "Transaksi gagal"
                                )
                            }
                        }
                    } else {
                        _transactionState.value = TransactionState.Error("Respons tidak valid")
                    }
                } else {
                    val exception = result.exceptionOrNull()
                    val message = exception?.message ?: "Transaksi gagal"

                    when {
                        message.contains("INSUFFICIENT_BALANCE", ignoreCase = true) ||
                        message.contains("Saldo tidak mencukupi", ignoreCase = true) -> {
                            _transactionState.value = TransactionState.InsufficientBalance(
                                required = _sellingPrice.value,
                                current = 0.0 // TODO: Get actual balance
                            )
                        }
                        message.contains("DAILY_LIMIT", ignoreCase = true) -> {
                            _transactionState.value = TransactionState.Error("Limit harian tercapai")
                        }
                        else -> {
                            _transactionState.value = TransactionState.Error(message)
                        }
                    }
                }
            } catch (e: Exception) {
                Log.e("TransactionInitVM", "Transaction initiation error", e)
                _transactionState.value = TransactionState.Error(
                    e.message ?: "Terjadi kesalahan tidak terduga"
                )
            }
        }
    }

    fun resetState() {
        _transactionState.value = TransactionState.Idle
        _customerNumber.value = ""
        _sellingPrice.value = 0.0
        currentProductId = ""
    }
}