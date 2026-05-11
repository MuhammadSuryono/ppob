package com.yonotech.ppob.mobile.viewmodels.staff

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.yonotech.ppob.mobile.data.remote.dto.CreateStaffRequest
import com.yonotech.ppob.mobile.data.remote.dto.StaffDto
import com.yonotech.ppob.mobile.data.remote.dto.TransactionHistoryResponse
import com.yonotech.ppob.mobile.data.repository.StaffRepository
import com.yonotech.ppob.mobile.utils.Resource
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

@HiltViewModel
class StaffViewModel @Inject constructor(
    private val repository: StaffRepository
) : ViewModel() {

    private val _staffListState = MutableStateFlow<Resource<List<StaffDto>>>(Resource.Idle)
    val staffListState = _staffListState.asStateFlow()

    private val _createStaffState = MutableStateFlow<Resource<StaffDto>>(Resource.Idle)
    val createStaffState = _createStaffState.asStateFlow()

    private val _topUpState = MutableStateFlow<Resource<TransactionHistoryResponse>>(Resource.Idle)
    val topUpState = _topUpState.asStateFlow()

    private val _transactionHistoryState = MutableStateFlow<Resource<List<TransactionHistoryResponse>>>(Resource.Idle)
    val transactionHistoryState = _transactionHistoryState.asStateFlow()

    fun getStaffList() {
        viewModelScope.launch {
            _staffListState.value = Resource.Loading
            _staffListState.value = repository.getStaffList()
        }
    }

    fun createStaff(request: CreateStaffRequest) {
        viewModelScope.launch {
            _createStaffState.value = Resource.Loading
            _createStaffState.value = repository.createStaff(request)
        }
    }

    fun topUpStaff(staffId: String, amount: Double, pin: String) {
        viewModelScope.launch {
            _topUpState.value = Resource.Loading
            _topUpState.value = repository.topUpStaff(staffId, amount, pin)
        }
    }

    fun getTransactionHistory() {
        viewModelScope.launch {
            _transactionHistoryState.value = Resource.Loading
            _transactionHistoryState.value = repository.getTransactionHistory()
        }
    }

    fun resetState() {
        _createStaffState.value = Resource.Idle
        _topUpState.value = Resource.Idle
    }
}