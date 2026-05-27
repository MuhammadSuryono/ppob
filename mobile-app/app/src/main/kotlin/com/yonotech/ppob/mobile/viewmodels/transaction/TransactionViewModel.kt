package com.yonotech.ppob.mobile.viewmodels.transaction

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.yonotech.ppob.mobile.data.local.TokenManager
import com.yonotech.ppob.mobile.data.remote.dto.AuthorizeRequest
import com.yonotech.ppob.mobile.data.remote.dto.WalletResponse
import com.yonotech.ppob.mobile.data.remote.dto.transaction.InitiateTransactionRequest
import com.yonotech.ppob.mobile.data.remote.dto.transaction.InquiryRequest
import com.yonotech.ppob.mobile.data.remote.dto.transaction.InquiryResponse
import com.yonotech.ppob.mobile.data.remote.dto.transaction.TransactionResponse
import com.yonotech.ppob.mobile.data.repository.AuthRepository
import com.yonotech.ppob.mobile.data.repository.TransactionRepository
import com.yonotech.ppob.mobile.data.repository.WalletRepository
import com.yonotech.ppob.mobile.utils.Resource
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.flow.first
import kotlinx.coroutines.launch
import javax.inject.Inject

@HiltViewModel
class TransactionViewModel @Inject constructor(
    private val transactionRepository: TransactionRepository,
    private val authRepository: AuthRepository,
    private val walletRepository: WalletRepository,
    private val tokenManager: TokenManager
) : ViewModel() {

    private val _transactionState = MutableStateFlow<Resource<TransactionResponse>>(Resource.Idle)
    val transactionState = _transactionState.asStateFlow()

    private val _inquiryState = MutableStateFlow<Resource<InquiryResponse>>(Resource.Idle)
    val inquiryState = _inquiryState.asStateFlow()

    private val _walletState = MutableStateFlow<Resource<WalletResponse>>(Resource.Idle)
    val walletState = _walletState.asStateFlow()

    // State for the multi-step flow
    var selectedProductCode: String = ""
    var customerNo: String = ""
    var amount: Double = 0.0

    fun checkBalance(requiredAmount: Double) {
        amount = requiredAmount
        viewModelScope.launch {
            _walletState.value = Resource.Loading
            try {
                // Using 'me' as id for current user wallet
                val response = walletRepository.getBalance("me")
                if (response.isSuccessful && response.body() != null) {
                    val wallet = response.body()!!
                    if (wallet.balanceAvailable >= requiredAmount) {
                        _walletState.value = Resource.Success(wallet)
                    } else {
                        _walletState.value = Resource.Error("Saldo tidak mencukupi. Saldo Anda: Rp ${wallet.balanceAvailable}")
                    }
                } else {
                    _walletState.value = Resource.Error("Gagal memeriksa saldo: ${response.message()}")
                }
            } catch (e: Exception) {
                _walletState.value = Resource.Error(e.message ?: "Terjadi kesalahan sistem saat cek saldo")
            }
        }
    }

    fun performInquiry(categoryId: Long, brand: String, customerNo: String) {
        this.customerNo = customerNo
        viewModelScope.launch {
            _inquiryState.value = Resource.Loading
            try {
                val response = transactionRepository.inquiry(
                    InquiryRequest(
                        categoryId,
                        brand,
                        customerNo
                    )
                )
                if (response.isSuccessful && response.body() != null) {
                    _inquiryState.value = Resource.Success(response.body()!!)
                } else {
                    _inquiryState.value = Resource.Error("Gagal cek data pelanggan: ${response.message()}")
                }
            } catch (e: Exception) {
                _inquiryState.value = Resource.Error(e.message ?: "Terjadi kesalahan saat inquiry")
            }
        }
    }

    fun initiateTransaction(pin: String) {
        viewModelScope.launch {
            _transactionState.value = Resource.Loading
            try {
                // Get the login token from DataStore
                val token = tokenManager.accessToken.first()
                
                if (token.isNullOrEmpty()) {
                    _transactionState.value = Resource.Error("Sesi berakhir, silakan login kembali.")
                    return@launch
                }

                // Step 1: Authorize PIN to get authorize_id
                val authResponse = authRepository.authorize(AuthorizeRequest(pin))
                
                if (authResponse.isSuccessful && authResponse.body() != null) {
                    val authorizeId = authResponse.body()!!.authorizeId
                    
                    // Step 2: Initiate Transaction with authorize_id
                    val request = InitiateTransactionRequest(
                        productCode = selectedProductCode,
                        customerNumber = customerNo,
                        amount = amount,
                        authorizeId = authorizeId
                    )
                    val response = transactionRepository.initiate(request)
                    if (response.isSuccessful) {
                        _transactionState.value = Resource.Success(response.body()!!)
                    } else {
                        _transactionState.value = Resource.Error("Transaksi gagal: ${response.message()}")
                    }
                } else {
                    _transactionState.value = Resource.Error("PIN salah atau otorisasi gagal.")
                }
            } catch (e: Exception) {
                _transactionState.value = Resource.Error(e.message ?: "Terjadi kesalahan sistem")
            }
        }
    }

    fun fetchTransactionStatus(txId: String) {
        viewModelScope.launch {
            _transactionState.value = Resource.Loading
            try {
                val response = transactionRepository.getStatus(txId)
                if (response.isSuccessful && response.body() != null) {
                    _transactionState.value = Resource.Success(response.body()!!)
                } else {
                    _transactionState.value = Resource.Error("Gagal mengambil status transaksi")
                }
            } catch (e: Exception) {
                _transactionState.value = Resource.Error(e.message ?: "Terjadi kesalahan")
            }
        }
    }

    fun startPolling(txId: String) {
        viewModelScope.launch {
            while (true) {
                try {
                    val response = transactionRepository.getStatus(txId)
                    if (response.isSuccessful && response.body() != null) {
                        val transaction = response.body()!!
                        _transactionState.value = Resource.Success(transaction)
                        
                        // Stop polling if status is terminal
                        val status = transaction.status.lowercase()
                        if (status == "success" || status == "failed") {
                            break
                        }
                    }
                } catch (e: Exception) {
                    // Silently ignore polling errors to keep trying
                }
                kotlinx.coroutines.delay(3000) // Poll every 3 seconds
            }
        }
    }

    fun resetState() {
        _transactionState.value = Resource.Idle
    }
}