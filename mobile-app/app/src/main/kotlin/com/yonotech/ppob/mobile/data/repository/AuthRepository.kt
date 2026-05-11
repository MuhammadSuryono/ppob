package com.yonotech.ppob.mobile.data.repository

import com.yonotech.ppob.mobile.data.remote.AuthService
import com.yonotech.ppob.mobile.data.remote.dto.*
import retrofit2.Response
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class AuthRepository @Inject constructor(
    private val authService: AuthService
) {
    suspend fun sendOtp(request: SendOtpRequest): Response<AuthResponse> {
        return authService.sendOtp(request)
    }

    suspend fun login(request: LoginRequest): Response<AuthResponse> {
        return authService.login(request)
    }

    suspend fun register(request: RegisterRequest): Response<AuthResponse> {
        return authService.register(request)
    }

    suspend fun verifyOtp(request: VerifyOtpRequest): Response<AuthResponse> {
        return authService.verifyOtp(request)
    }

    suspend fun setPasswordPin(request: SetPasswordPinRequest): Response<AuthResponse> {
        return authService.setPasswordPin(request)
    }

    suspend fun pinLogin(request: PinLoginRequest): Response<AuthResponse> {
        return authService.pinLogin(request)
    }
}