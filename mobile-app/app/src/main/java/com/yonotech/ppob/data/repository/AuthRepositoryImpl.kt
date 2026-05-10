package com.yonotech.ppob.data.repository

import android.util.Log
import com.squareup.moshi.Moshi
import com.yonotech.ppob.data.local.dao.UserDao
import com.yonotech.ppob.data.local.datastore.AuthPreferencesDataStore
import com.yonotech.ppob.data.local.entity.UserEntity
import com.yonotech.ppob.data.remote.api.AuthApiService
import com.yonotech.ppob.data.remote.model.ApiError
import com.yonotech.ppob.data.remote.model.AuthResponse
import com.yonotech.ppob.data.remote.model.ChangePinRequest
import com.yonotech.ppob.data.remote.model.LoginRequest
import com.yonotech.ppob.data.remote.model.RegisterRequest
import com.yonotech.ppob.data.remote.model.VerifyOtpRequest
import com.yonotech.ppob.data.remote.model.RefreshTokenRequest
import com.yonotech.ppob.domain.repository.AuthRepository
import kotlinx.coroutines.flow.first

class AuthRepositoryImpl(
    private val apiService: AuthApiService,
    private val userDao: UserDao,
    private val authPreferences: AuthPreferencesDataStore
) : AuthRepository {

    override suspend fun register(request: RegisterRequest): Result<AuthResponse> {
        return try {
            val response = apiService.register(request)
            if (response.isSuccessful) {
                val body = response.body()
                if (body?.data != null) {
                    saveUserCredentials(body.data)
                    Result.success(body.data)
                } else {
                    Result.failure(Exception(body?.message ?: "Registration failed"))
                }
            } else {
                Result.failure(handleError(response.errorBody()?.string()))
            }
        } catch (e: Exception) {
            Log.e("AuthRepository", "Register error", e)
            Result.failure(e)
        }
    }

    override suspend fun login(request: LoginRequest): Result<AuthResponse> {
        return try {
            val response = apiService.login(request)
            if (response.isSuccessful) {
                val body = response.body()
                if (body?.data != null) {
                    saveUserCredentials(body.data)
                    Result.success(body.data)
                } else {
                    Result.failure(Exception(body?.message ?: "Login failed"))
                }
            } else {
                Result.failure(handleError(response.errorBody()?.string()))
            }
        } catch (e: Exception) {
            Log.e("AuthRepository", "Login error", e)
            Result.failure(e)
        }
    }

    override suspend fun verifyOtp(request: VerifyOtpRequest): Result<AuthResponse> {
        return try {
            val response = apiService.verifyOtp(request)
            if (response.isSuccessful) {
                val body = response.body()
                if (body?.data != null) {
                    saveUserCredentials(body.data)
                    Result.success(body.data)
                } else {
                    Result.failure(Exception(body?.message ?: "OTP verification failed"))
                }
            } else {
                Result.failure(handleError(response.errorBody()?.string()))
            }
        } catch (e: Exception) {
            Log.e("AuthRepository", "Verify OTP error", e)
            Result.failure(e)
        }
    }

    override suspend fun refreshToken(refreshToken: String): Result<AuthResponse> {
        return try {
            val request = RefreshTokenRequest(refreshToken)
            val response = apiService.refreshToken(request)
            if (response.isSuccessful) {
                val body = response.body()
                if (body?.data != null) {
                    body.data.let {
                        authPreferences.saveAccessToken(it.accessToken)
                        authPreferences.saveRefreshToken(it.refreshToken)
                    }
                    Result.success(body.data)
                } else {
                    Result.failure(Exception(body?.message ?: "Token refresh failed"))
                }
            } else {
                Result.failure(handleError(response.errorBody()?.string()))
            }
        } catch (e: Exception) {
            Log.e("AuthRepository", "Refresh token error", e)
            Result.failure(e)
        }
    }

    override suspend fun logout(token: String): Result<Unit> {
        return try {
            val response = apiService.logout("Bearer $token")
            if (response.isSuccessful) {
                authPreferences.clearAuthData()
                Result.success(Unit)
            } else {
                Result.failure(Exception("Logout failed"))
            }
        } catch (e: Exception) {
            Log.e("AuthRepository", "Logout error", e)
            Result.failure(e)
        }
    }

    override suspend fun changePin(token: String, request: ChangePinRequest): Result<Unit> {
        return try {
            val response = apiService.changePin("Bearer $token", request)
            if (response.isSuccessful) {
                val body = response.body()
                if (body?.success == true) {
                    Result.success(Unit)
                } else {
                    Result.failure(Exception(body?.message ?: "Change PIN failed"))
                }
            } else {
                Result.failure(handleError(response.errorBody()?.string()))
            }
        } catch (e: Exception) {
            Log.e("AuthRepository", "Change PIN error", e)
            Result.failure(e)
        }
    }

    override suspend fun getStoredToken(): String? {
        return try {
            authPreferences.authPreferences.first().accessToken
        } catch (e: Exception) {
            null
        }
    }

    override suspend fun clearTokens() {
        authPreferences.clearAuthData()
    }

    private suspend fun saveUserCredentials(data: AuthResponse) {
        try {
            authPreferences.saveUserCredentials(
                accessToken = data.accessToken,
                refreshToken = data.refreshToken,
                userId = data.user.userId,
                phoneNumber = data.user.phoneNumber,
                userName = data.user.name,
                activeRole = data.user.activeRole.roleName,
                walletId = data.user.walletId
            )

            val userEntity = UserEntity(
                userId = data.user.userId,
                phoneNumber = data.user.phoneNumber,
                name = data.user.name,
                activeRole = data.user.activeRole.roleName,
                walletId = data.user.walletId,
                accessToken = data.accessToken,
                refreshToken = data.refreshToken,
                tokenExpiresAt = System.currentTimeMillis() + (data.expiresIn * 1000L)
            )
            userDao.insert(userEntity)
        } catch (e: Exception) {
            Log.e("AuthRepository", "Error saving credentials", e)
        }
    }

    private fun handleError(errorBody: String?): Exception {
        return try {
            val moshi = Moshi.Builder().build()
            val adapter = moshi.adapter(ApiError::class.java)
            if (errorBody != null) {
                adapter.fromJson(errorBody)?.let { apiError ->
                    return Exception(apiError.message)
                }
            }
            Exception("API Error")
        } catch (e: Exception) {
            Exception(errorBody ?: "Unknown error")
        }
    }
}
