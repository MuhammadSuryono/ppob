package com.yonotech.ppob.mobile.data.remote.dto

import com.squareup.moshi.Json
import com.squareup.moshi.JsonClass

@JsonClass(generateAdapter = true)
data class InitiateAuthRequest(
    @Json(name = "phone") val phone: String,
    @Json(name = "device_id") val deviceId: String,
    @Json(name = "fingerprint") val fingerprint: String? = null
)

@JsonClass(generateAdapter = true)
data class InitiateAuthResponse(
    @Json(name = "user_id") val userId: Int? = null,
    @Json(name = "is_registered") val isRegistered: Boolean,
    @Json(name = "is_trusted") val isTrusted: Boolean,
    @Json(name = "requires_otp") val requiresOtp: Boolean
)

@JsonClass(generateAdapter = true)
data class SendOtpRequest(
    @Json(name = "phone") val phone: String,
    @Json(name = "type") val type: String // "login" or "register"
)

@JsonClass(generateAdapter = true)
data class SendOtpResponse(
    @Json(name = "request_id") val requestId: String,
    @Json(name = "expires_at") val expiresAt: Long
)

@JsonClass(generateAdapter = true)
data class VerifyOtpRequest(
    @Json(name = "request_id") val requestId: String,
    @Json(name = "phone") val phone: String,
    @Json(name = "code") val code: String,
    @Json(name = "type") val type: String // "login" or "register"
)

@JsonClass(generateAdapter = true)
data class VerifyOtpResponse(
    @Json(name = "request_id") val requestId: String,
    @Json(name = "is_verified") val isVerified: Boolean
)

@JsonClass(generateAdapter = true)
data class RegisterRequest(
    @Json(name = "email") val email: String,
    @Json(name = "phone") val phone: String,
    @Json(name = "full_name") val fullName: String,
    @Json(name = "password") val password: String,
    @Json(name = "pin") val pin: String,
    @Json(name = "device_id") val deviceId: String? = null,
    @Json(name = "request_id") val requestId: String
)

@JsonClass(generateAdapter = true)
data class VerifyPasswordRequest(
    @Json(name = "phone") val phone: String,
    @Json(name = "password") val password: String,
    @Json(name = "device_id") val deviceId: String,
    @Json(name = "request_id") val requestId: String
)

@JsonClass(generateAdapter = true)
data class AuthResponse(
    @Json(name = "access_token") val accessToken: String? = null,
    @Json(name = "refresh_token") val refreshToken: String? = null,
    @Json(name = "expires_at") val expiresAt: Long? = null,
    @Json(name = "user_id") val userId: Int? = null,
    @Json(name = "email") val email: String? = null,
    @Json(name = "phone") val phone: String? = null,
    @Json(name = "full_name") val fullName: String? = null
)

@JsonClass(generateAdapter = true)
data class UserDto(
    @Json(name = "id") val id: String,
    @Json(name = "email") val email: String,
    @Json(name = "phone") val phone: String,
    @Json(name = "name") val name: String
)

@JsonClass(generateAdapter = true)
data class ErrorResponse(
    @Json(name = "error") val error: ErrorDetail
)

@JsonClass(generateAdapter = true)
data class ErrorDetail(
    @Json(name = "code") val code: String,
    @Json(name = "message") val message: String
)