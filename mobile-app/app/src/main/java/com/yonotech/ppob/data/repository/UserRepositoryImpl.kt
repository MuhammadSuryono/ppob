package com.yonotech.ppob.data.repository

import android.util.Log
import com.yonotech.ppob.data.local.dao.UserDao
import com.yonotech.ppob.data.remote.api.UserApiService
import com.yonotech.ppob.data.remote.model.AddStaffRequest
import com.yonotech.ppob.data.remote.model.StaffDetailResponse
import com.yonotech.ppob.data.remote.model.StaffListResponse
import com.yonotech.ppob.data.remote.model.UpdateStaffRequest
import com.yonotech.ppob.data.remote.model.UserProfileResponse
import com.yonotech.ppob.domain.repository.UserRepository

class UserRepositoryImpl(
    private val apiService: UserApiService,
    private val userDao: UserDao
) : UserRepository {

    override suspend fun getProfile(token: String): Result<UserProfileResponse> {
        return try {
            val response = apiService.getProfile("Bearer $token")
            if (response.isSuccessful) {
                val body = response.body()
                if (body?.success == true && body.data != null) {
                    Result.success(body.data)
                } else {
                    Result.failure(Exception(body?.message ?: "Failed to get profile"))
                }
            } else {
                Result.failure(Exception("HTTP ${response.code()}"))
            }
        } catch (e: Exception) {
            Log.e("UserRepository", "Get profile error", e)
            Result.failure(e)
        }
    }

    override suspend fun switchRole(token: String, roleId: String): Result<com.yonotech.ppob.data.remote.model.SwitchRoleResponse> {
        return try {
            val response = apiService.switchRole(
                "Bearer $token",
                com.yonotech.ppob.data.remote.model.SwitchRoleRequest(roleId)
            )
            if (response.isSuccessful) {
                val body = response.body()
                if (body?.success == true && body.data != null) {
                    Result.success(body.data)
                } else {
                    Result.failure(Exception(body?.message ?: "Failed to switch role"))
                }
            } else {
                Result.failure(Exception("HTTP ${response.code()}"))
            }
        } catch (e: Exception) {
            Log.e("UserRepository", "Switch role error", e)
            Result.failure(e)
        }
    }

    override suspend fun getStaff(token: String, limit: Int, offset: Int): Result<StaffListResponse> {
        return try {
            val response = apiService.getStaff("Bearer $token", limit, offset)
            if (response.isSuccessful) {
                val body = response.body()
                if (body?.success == true && body.data != null) {
                    Result.success(body.data)
                } else {
                    Result.failure(Exception(body?.message ?: "Failed to get staff"))
                }
            } else {
                Result.failure(Exception("HTTP ${response.code()}"))
            }
        } catch (e: Exception) {
            Log.e("UserRepository", "Get staff error", e)
            Result.failure(e)
        }
    }

    override suspend fun addStaff(token: String, request: AddStaffRequest): Result<StaffDetailResponse> {
        return try {
            val response = apiService.addStaff("Bearer $token", request)
            if (response.isSuccessful) {
                val body = response.body()
                if (body?.success == true && body.data != null) {
                    Result.success(body.data)
                } else {
                    Result.failure(Exception(body?.message ?: "Failed to add staff"))
                }
            } else {
                Result.failure(Exception("HTTP ${response.code()}"))
            }
        } catch (e: Exception) {
            Log.e("UserRepository", "Add staff error", e)
            Result.failure(e)
        }
    }

    override suspend fun updateStaff(token: String, staffId: String, request: UpdateStaffRequest): Result<StaffDetailResponse> {
        return try {
            val response = apiService.updateStaff("Bearer $token", staffId, request)
            if (response.isSuccessful) {
                val body = response.body()
                if (body?.success == true && body.data != null) {
                    Result.success(body.data)
                } else {
                    Result.failure(Exception(body?.message ?: "Failed to update staff"))
                }
            } else {
                Result.failure(Exception("HTTP ${response.code()}"))
            }
        } catch (e: Exception) {
            Log.e("UserRepository", "Update staff error", e)
            Result.failure(e)
        }
    }

    override suspend fun getStaffDetail(token: String, staffId: String): Result<StaffDetailResponse> {
        return try {
            val response = apiService.getStaffDetail("Bearer $token", staffId)
            if (response.isSuccessful) {
                val body = response.body()
                if (body?.success == true && body.data != null) {
                    Result.success(body.data)
                } else {
                    Result.failure(Exception(body?.message ?: "Failed to get staff detail"))
                }
            } else {
                Result.failure(Exception("HTTP ${response.code()}"))
            }
        } catch (e: Exception) {
            Log.e("UserRepository", "Get staff detail error", e)
            Result.failure(e)
        }
    }
}