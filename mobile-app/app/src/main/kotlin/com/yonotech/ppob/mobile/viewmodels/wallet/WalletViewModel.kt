package com.yonotech.ppob.mobile.viewmodels.wallet

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.yonotech.ppob.mobile.data.remote.dto.WalletResponse
import com.yonotech.ppob.mobile.data.repository.WalletRepository
import com.yonotech.ppob.mobile.utils.Resource
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

@HiltViewModel
class WalletViewModel @Inject constructor(
    private val repository: WalletRepository
) : ViewModel() {

    private val _balanceState = MutableStateFlow<Resource<WalletResponse>>(Resource.Idle)
    val balanceState = _balanceState.asStateFlow()

    fun getBalance(walletId: String = "me") {
        viewModelScope.launch {
            _balanceState.value = Resource.Loading
            try {
                val response = repository.getBalance(walletId)
                if (response.isSuccessful && response.body() != null) {
                    _balanceState.value = Resource.Success(response.body()!!)
                } else {
                    _balanceState.value = Resource.Error("Gagal mengambil saldo: ${response.message()}")
                }
            } catch (e: Exception) {
                _balanceState.value = Resource.Error(e.message ?: "Terjadi kesalahan sistem")
            }
        }
    }

    fun refreshBalance() {
        getBalance()
    }
}