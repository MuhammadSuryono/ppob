package com.yonotech.ppob.mobile.data.repository

import com.yonotech.ppob.mobile.data.remote.AuthService
import com.yonotech.ppob.mobile.data.remote.dto.AuthResponse
import com.yonotech.ppob.mobile.data.remote.dto.LoginRequest
import com.yonotech.ppob.mobile.data.remote.dto.RegisterRequest
import com.yonotech.ppob.mobile.data.remote.dto.VerifyOtpRequest
import retrofit2.Response
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class AuthRepository @Inject constructor(
    private val authService: AuthService
) {
    suspend fun login(request: LoginRequest): Response<AuthResponse> {
        return authService.login(request)
    }

    suspend fun register(request: RegisterRequest): Response<AuthResponse> {
        return authService.register(request)
    }

    suspend fun verifyOtp(request: VerifyOtpRequest): Response<AuthResponse> {
        return authService.verifyOtp(request)
    }
}
