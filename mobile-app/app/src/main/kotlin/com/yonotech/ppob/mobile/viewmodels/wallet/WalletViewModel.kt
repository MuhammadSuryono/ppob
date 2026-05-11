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

    fun getBalance() {
        viewModelScope.launch {
            _balanceState.value = Resource.Loading
            _balanceState.value = repository.getBalance()
        }
    }

    fun refreshBalance() {
        getBalance()
    }
}