package com.yonotech.ppob.presentation.staff

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.yonotech.ppob.data.remote.model.StaffListResponse
import com.yonotech.ppob.data.remote.model.StaffDetailResponse
import com.yonotech.ppob.domain.model.Staff
import com.yonotech.ppob.domain.repository.UserRepository
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

data class StaffListUiState(
    val staffList: List<Staff> = emptyList(),
    val isLoading: Boolean = true,
    val error: String? = null,
    val hasMore: Boolean = false,
    val nextOffset: Int = 20,
    val searchQuery: String = ""
)

data class StaffDetailUiState(
    val staff: Staff? = null,
    val isLoading: Boolean = true,
    val error: String? = null
)

@HiltViewModel
class StaffListViewModel @Inject constructor(
    private val userRepository: UserRepository
) : ViewModel() {

    private val _uiState = MutableStateFlow(StaffListUiState())
    val uiState: StateFlow<StaffListUiState> = _uiState.asStateFlow()

    private var authToken: String? = null
    private var currentPage: Int = 0
    private var isLastPage: Boolean = false

    fun setAuthToken(token: String) {
        authToken = token
        loadStaff(reset = true)
    }

    fun loadStaff(reset: Boolean = false) {
        if (reset) {
            currentPage = 0
            isLastPage = false
        }
        if (isLastPage) return

        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isLoading = true, error = null)

            try {
                val token = authToken ?: run {
                    _uiState.value = _uiState.value.copy(
                        isLoading = false,
                        error = "Token tidak tersedia"
                    )
                    return@launch
                }

                val result = userRepository.getStaff(token, 20, currentPage * 20)

                if (result.isSuccess) {
                    val response = result.getOrNull()?.staff ?: emptyList()
                    val staffList = if (reset) response else (_uiState.value.staffList + response)

                    isLastPage = response.size < 20
                    currentPage++

                    _uiState.value = StaffListUiState(
                        staffList = staffList,
                        isLoading = false,
                        error = null,
                        hasMore = !isLastPage,
                        nextOffset = currentPage * 20,
                        searchQuery = _uiState.value.searchQuery
                    )
                } else {
                    _uiState.value = _uiState.value.copy(
                        isLoading = false,
                        error = result.exceptionOrNull()?.message ?: "Gagal memuat staff"
                    )
                }
            } catch (e: Exception) {
                _uiState.value = _uiState.value.copy(
                    isLoading = false,
                    error = e.message ?: "Gagal memuat data staff"
                )
            }
        }
    }

    fun searchStaff(query: String) {
        _uiState.value = _uiState.value.copy(searchQuery = query)
        // Filter staff list locally
        val filtered = _uiState.value.staffList.filter {
            it.name.contains(query, ignoreCase = true) ||
            it.phoneNumber.contains(query, ignoreCase = true)
        }
        _uiState.value = _uiState.value.copy(
            staffList = if (query.isEmpty()) {
                // Reload original
                _uiState.value.staffList
            } else {
                filtered
            }
        )
    }

    fun refresh() {
        loadStaff(reset = true)
    }

    fun clearError() {
        _uiState.value = _uiState.value.copy(error = null)
    }
}

@HiltViewModel
class AddStaffViewModel @Inject constructor(
    private val userRepository: UserRepository
) : ViewModel() {

    private val _uiState = MutableStateFlow(AddStaffUiState())
    val uiState: StateFlow<AddStaffUiState> = _uiState.asStateFlow()

    private var authToken: String? = null

    fun setAuthToken(token: String) {
        authToken = token
    }

    fun addStaff(
        phone: String,
        name: String,
        password: String,
        pin: String,
        marginScheme: String = "FixedAllowance",
        marginValue: Double = 10000.0,
        dailyTxnLimit: Int = 50,
        dailyAmountLimit: Double = 5000000.0
    ) {
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isLoading = true, error = null)

            try {
                val token = authToken ?: run {
                    _uiState.value = _uiState.value.copy(
                        isLoading = false,
                        error = "Token tidak tersedia"
                    )
                    return@launch
                }

                val request = com.yonotech.ppob.data.remote.model.AddStaffRequest(
                    phone, name, password, pin, marginScheme, marginValue,
                    dailyTxnLimit, dailyAmountLimit
                )

                val result = userRepository.addStaff(token, request)

                if (result.isSuccess) {
                    _uiState.value = AddStaffUiState(
                        isLoading = false,
                        isSuccess = true,
                        message = "Staff berhasil ditambahkan"
                    )
                } else {
                    _uiState.value = AddStaffUiState(
                        isLoading = false,
                        error = result.exceptionOrNull()?.message ?: "Gagal menambahkan staff"
                    )
                }
            } catch (e: Exception) {
                _uiState.value = AddStaffUiState(
                    isLoading = false,
                    error = e.message ?: "Gagal menambahkan staff"
                )
            }
        }
    }

    fun clearError() {
        _uiState.value = _uiState.value.copy(error = null)
    }
}

data class AddStaffUiState(
    val isLoading: Boolean = false,
    val isSuccess: Boolean = false,
    val message: String? = null,
    val error: String? = null
)

@HiltViewModel
class StaffDetailViewModel @Inject constructor(
    private val userRepository: UserRepository
) : ViewModel() {

    private val _uiState = MutableStateFlow(StaffDetailUiState())
    val uiState: StateFlow<StaffDetailUiState> = _uiState.asStateFlow()

    private var authToken: String? = null

    fun setAuthToken(token: String) {
        authToken = token
    }

    fun loadStaffDetail(staffId: String) {
        viewModelScope.launch {
            _uiState.value = StaffDetailUiState(isLoading = true, error = null)

            try {
                val token = authToken ?: run {
                    _uiState.value = StaffDetailUiState(
                        isLoading = false,
                        error = "Token tidak tersedia"
                    )
                    return@launch
                }

                val result = userRepository.getStaffDetail(token, staffId)

                if (result.isSuccess) {
                    val response = result.getOrNull()
                    if (response != null) {
                        val staff = Staff(
                            userId = response.userId,
                            name = response.name,
                            phoneNumber = response.phoneNumber,
                            walletBalance = response.walletBalance,
                            dailyTxnCount = response.dailyTxnCount,
                            dailyTxnAmount = response.dailyTxnAmount,
                            marginScheme = response.marginScheme,
                            marginValue = response.marginValue,
                            isActive = response.isActive
                        )
                        _uiState.value = StaffDetailUiState(
                            staff = staff,
                            isLoading = false,
                            error = null
                        )
                    }
                } else {
                    _uiState.value = StaffDetailUiState(
                        isLoading = false,
                        error = result.exceptionOrNull()?.message ?: "Gagal memuat detail staff"
                    )
                }
            } catch (e: Exception) {
                _uiState.value = StaffDetailUiState(
                    isLoading = false,
                    error = e.message ?: "Gagal memuat detail staff"
                )
            }
        }
    }

    fun updateStaff(
        staffId: String,
        name: String? = null,
        marginScheme: String? = null,
        marginValue: Double? = null,
        dailyTxnLimit: Int? = null,
        dailyAmountLimit: Double? = null,
        isActive: Boolean? = null
    ) {
        viewModelScope.launch {
            _uiState.value = StaffDetailUiState(isLoading = true, error = null)

            try {
                val token = authToken ?: return@launch

                val request = com.yonotech.ppob.data.remote.model.UpdateStaffRequest(
                    name = name,
                    marginScheme = marginScheme,
                    marginValue = marginValue,
                    dailyTxnLimit = dailyTxnLimit,
                    dailyAmountLimit = dailyAmountLimit,
                    isActive = isActive
                )

                val result = userRepository.updateStaff(token, staffId, request)

                if (result.isSuccess) {
                    _uiState.value = StaffDetailUiState(
                        isLoading = false,
                        message = "Staff berhasil diperbarui"
                    )
                    loadStaffDetail(staffId) // Reload data
                } else {
                    _uiState.value = StaffDetailUiState(
                        isLoading = false,
                        error = result.exceptionOrNull()?.message ?: "Gagal memperbarui staff"
                    )
                }
            } catch (e: Exception) {
                _uiState.value = StaffDetailUiState(
                    isLoading = false,
                    error = e.message ?: "Gagal memperbarui staff"
                )
            }
        }
    }
}