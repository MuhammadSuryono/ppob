package com.yonotech.ppob.mobile.viewmodels.transaction

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.yonotech.ppob.mobile.data.remote.dto.transaction.InitiateTransactionRequest
import com.yonotech.ppob.mobile.data.remote.dto.transaction.TransactionResponse
import com.yonotech.ppob.mobile.data.repository.TransactionRepository
import com.yonotech.ppob.mobile.utils.Resource
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

@HiltViewModel
class TransactionViewModel @Inject constructor(
    private val repository: TransactionRepository
) : ViewModel() {

    private val _transactionState = MutableStateFlow<Resource<TransactionResponse>>(Resource.Idle)
    val transactionState = _transactionState.asStateFlow()

    // State for the multi-step flow
    var selectedProductId: String = ""
    var customerNo: String = ""
    var sellingPrice: Double = 0.0

    fun initiateTransaction(pin: String) {
        viewModelScope.launch {
            _transactionState.value = Resource.Loading
            try {
                val request = InitiateTransactionRequest(
                    productId = selectedProductId,
                    customerNo = customerNo,
                    sellingPrice = sellingPrice,
                    pin = pin
                )
                val response = repository.initiate(request)
                if (response.isSuccessful) {
                    _transactionState.value = Resource.Success(response.body()!!)
                } else {
                    _transactionState.value = Resource.Error("Transaction failed: ${response.message()}")
                }
            } catch (e: Exception) {
                _transactionState.value = Resource.Error(e.message ?: "An unknown error occurred")
            }
        }
    }

    fun resetState() {
        _transactionState.value = Resource.Idle
    }
}
