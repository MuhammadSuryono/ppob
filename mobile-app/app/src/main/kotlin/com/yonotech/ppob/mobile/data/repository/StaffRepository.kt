package com.yonotech.ppob.mobile.data.repository

import com.yonotech.ppob.mobile.data.remote.StaffService
import com.yonotech.ppob.mobile.data.remote.dto.CreateStaffRequest
import com.yonotech.ppob.mobile.data.remote.dto.StaffDto
import com.yonotech.ppob.mobile.data.remote.dto.TopUpRequest
import com.yonotech.ppob.mobile.data.remote.dto.TransactionHistoryResponse
import com.yonotech.ppob.mobile.utils.Resource
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class StaffRepository @Inject constructor(
    private val staffService: StaffService
) {
    suspend fun getStaffList(): Resource<List<StaffDto>> {
        return try {
            val response = staffService.getStaffList()
            if (response.isSuccessful) {
                Resource.Success(response.body()!!)
            } else {
                Resource.Error("Failed to get staff list: ${response.message()}")
            }
        } catch (e: Exception) {
            Resource.Error(e.message ?: "Unknown error")
        }
    }

    suspend fun createStaff(request: CreateStaffRequest): Resource<StaffDto> {
        return try {
            val response = staffService.createStaff(request)
            if (response.isSuccessful) {
                Resource.Success(response.body()!!)
            } else {
                Resource.Error("Failed to create staff: ${response.message()}")
            }
        } catch (e: Exception) {
            Resource.Error(e.message ?: "Unknown error")
        }
    }

    suspend fun topUpStaff(staffId: String, amount: Double, pin: String): Resource<TransactionHistoryResponse> {
        return try {
            val response = staffService.topUpStaff(TopUpRequest(staffId, amount, pin))
            if (response.isSuccessful) {
                Resource.Success(response.body()!!)
            } else {
                Resource.Error("Top-up failed: ${response.message()}")
            }
        } catch (e: Exception) {
            Resource.Error(e.message ?: "Unknown error")
        }
    }

    suspend fun getTransactionHistory(): Resource<List<TransactionHistoryResponse>> {
        return try {
            val response = staffService.getTransactionHistory()
            if (response.isSuccessful) {
                Resource.Success(response.body()!!)
            } else {
                Resource.Error("Failed to get transaction history: ${response.message()}")
            }
        } catch (e: Exception) {
            Resource.Error(e.message ?: "Unknown error")
        }
    }
}