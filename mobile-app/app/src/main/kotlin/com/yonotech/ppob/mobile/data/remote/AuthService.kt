package com.yonotech.ppob.mobile.data.remote

import com.yonotech.ppob.mobile.data.remote.dto.*
import retrofit2.Response
import retrofit2.http.Body
import retrofit2.http.POST

interface AuthService {
    @POST("auth/initiate")
    suspend fun initiateAuth(@Body request: InitiateAuthRequest): Response<InitiateAuthResponse>

    @POST("auth/send-otp")
    suspend fun sendOtp(@Body request: SendOtpRequest): Response<SendOtpResponse>

    @POST("auth/verify-otp")
    suspend fun verifyOtp(@Body request: VerifyOtpRequest): Response<VerifyOtpResponse>

    @POST("auth/register")
    suspend fun register(@Body request: RegisterRequest): Response<AuthResponse>

    @POST("auth/login")
    suspend fun login(@Body request: LoginRequest): Response<AuthResponse>

    @POST("auth/verify-password")
    suspend fun verifyPassword(@Body request: VerifyPasswordRequest): Response<AuthResponse>

    @POST("auth/verify-pin")
    suspend fun verifyPin(@Body request: LoginRequest): Response<AuthResponse>

    @POST("auth/refresh")
    suspend fun refreshToken(@Body refreshToken: String): Response<AuthResponse>
}