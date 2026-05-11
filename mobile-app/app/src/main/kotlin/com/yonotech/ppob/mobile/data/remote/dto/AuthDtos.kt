package com.yonotech.ppob.mobile.data.remote.dto

import com.squareup.moshi.Json
import com.squareup.moshi.JsonClass

@JsonClass(generateAdapter = true)
data class LoginRequest(
    @Json(name = "identifier") val identifier: String, // email or phone
    @Json(name = "password") val password: String,
    @Json(name = "device_id") val deviceId: String
)

@JsonClass(generateAdapter = true)
data class RegisterRequest(
    @Json(name = "email") val email: String,
    @Json(name = "phone") val phone: String,
    @Json(name = "name") val name: String,
    @Json(name = "password") val password: String,
    @Json(name = "pin") val pin: String
)

@JsonClass(generateAdapter = true)
data class VerifyOtpRequest(
    @Json(name = "identifier") val identifier: String,
    @Json(name = "otp_code") val otpCode: String,
    @Json(name = "type") val type: String // "registration" or "login"
)

@JsonClass(generateAdapter = true)
data class SendOtpRequest(
    @Json(name = "phone") val phone: String,
    @Json(name = "device_id") val deviceId: String
)

@JsonClass(generateAdapter = true)
data class SetPasswordPinRequest(
    @Json(name = "phone") val phone: String,
    @Json(name = "password") val password: String,
    @Json(name = "pin") val pin: String
)

@JsonClass(generateAdapter = true)
data class PinLoginRequest(
    @Json(name = "phone") val phone: String,
    @Json(name = "pin") val pin: String,
    @Json(name = "device_id") val deviceId: String
)

@JsonClass(generateAdapter = true)
data class AuthResponse(
    @Json(name = "access_token") val accessToken: String?,
    @Json(name = "refresh_token") val refreshToken: String?,
    @Json(name = "requires_otp") val requiresOtp: Boolean = false,
    @Json(name = "is_new_user") val isNewUser: Boolean = false,
    @Json(name = "user") val user: UserDto? = null
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