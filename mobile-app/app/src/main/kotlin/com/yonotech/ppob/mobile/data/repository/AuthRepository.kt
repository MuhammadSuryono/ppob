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
    suspend fun initiateAuth(request: InitiateAuthRequest): Response<InitiateAuthResponse> {
        return authService.initiateAuth(request)
    }

    suspend fun sendOtp(request: SendOtpRequest): Response<SendOtpResponse> {
        return authService.sendOtp(request)
    }

    suspend fun verifyOtp(request: VerifyOtpRequest): Response<VerifyOtpResponse> {
        return authService.verifyOtp(request)
    }

    suspend fun register(request: RegisterRequest): Response<AuthResponse> {
        return authService.register(request)
    }

    suspend fun login(request: LoginRequest): Response<AuthResponse> {
        return authService.login(request)
    }

    suspend fun verifyPassword(request: VerifyPasswordRequest): Response<AuthResponse> {
        return authService.verifyPassword(request)
    }

    suspend fun verifyPin(request: LoginRequest): Response<AuthResponse> {
        return authService.verifyPin(request)
    }
}