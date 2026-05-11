package com.yonotech.ppob.mobile.data.remote

import com.yonotech.ppob.mobile.data.remote.dto.*
import retrofit2.Response
import retrofit2.http.Body
import retrofit2.http.POST

interface AuthService {
    @POST("auth/send-otp")
    suspend fun sendOtp(@Body request: SendOtpRequest): Response<AuthResponse>

    @POST("auth/register")
    suspend fun register(@Body request: RegisterRequest): Response<AuthResponse>

    @POST("auth/login")
    suspend fun login(@Body request: LoginRequest): Response<AuthResponse>

    @POST("auth/verify-otp")
    suspend fun verifyOtp(@Body request: VerifyOtpRequest): Response<AuthResponse>

    @POST("auth/set-password-pin")
    suspend fun setPasswordPin(@Body request: SetPasswordPinRequest): Response<AuthResponse>

    @POST("auth/pin-login")
    suspend fun pinLogin(@Body request: PinLoginRequest): Response<AuthResponse>

    @POST("auth/refresh")
    suspend fun refreshToken(@Body refreshToken: String): Response<AuthResponse>
}