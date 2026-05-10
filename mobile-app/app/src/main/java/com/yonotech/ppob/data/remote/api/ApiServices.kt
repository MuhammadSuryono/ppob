package com.yonotech.ppob.data.remote.api

import com.yonotech.ppob.data.remote.model.ApiResponse
import com.yonotech.ppob.data.remote.model.AuthResponse
import com.yonotech.ppob.data.remote.model.CategoriesResponse
import com.yonotech.ppob.data.remote.model.ChangePinRequest
import com.yonotech.ppob.data.remote.model.LoginRequest
import com.yonotech.ppob.data.remote.model.ProductsResponse
import com.yonotech.ppob.data.remote.model.RegisterRequest
import com.yonotech.ppob.data.remote.model.TransactionHistoryResponse
import com.yonotech.ppob.data.remote.model.TransactionResponse
import com.yonotech.ppob.data.remote.model.VerifyOtpRequest
import com.yonotech.ppob.data.remote.model.WalletBalanceResponse
import okhttp3.ResponseBody
import retrofit2.Response
import retrofit2.http.Body
import retrofit2.http.GET
import retrofit2.http.Header
import retrofit2.http.Headers
import retrofit2.http.POST
import retrofit2.http.PUT
import retrofit2.http.Path
import retrofit2.http.Query

interface AuthApiService {

    @POST("auth/register")
    suspend fun register(@Body request: RegisterRequest): Response<ApiResponse<AuthResponse>>

    @POST("auth/login")
    suspend fun login(@Body request: LoginRequest): Response<ApiResponse<AuthResponse>>

    @POST("auth/verify-otp")
    suspend fun verifyOtp(@Body request: VerifyOtpRequest): Response<ApiResponse<AuthResponse>>

    @POST("auth/refresh")
    suspend fun refreshToken(@Body request: com.yonotech.ppob.data.remote.model.RefreshTokenRequest): Response<ApiResponse<AuthResponse>>

    @POST("auth/logout")
    suspend fun logout(@Header("Authorization") token: String): Response<ApiResponse<Unit>>

    @POST("auth/change-pin")
    suspend fun changePin(
        @Header("Authorization") token: String,
        @Body request: ChangePinRequest
    ): Response<ApiResponse<Unit>>
}

interface UserApiService {

    @GET("users/v1/profile")
    suspend fun getProfile(
        @Header("Authorization") token: String
    ): Response<ApiResponse<com.yonotech.ppob.data.remote.model.UserProfileResponse>>

    @POST("users/v1/switch-role")
    suspend fun switchRole(
        @Header("Authorization") token: String,
        @Body request: com.yonotech.ppob.data.remote.model.SwitchRoleRequest
    ): Response<ApiResponse<com.yonotech.ppob.data.remote.model.SwitchRoleResponse>>

    @GET("users/v1/staff")
    suspend fun getStaff(
        @Header("Authorization") token: String,
        @Query("limit") limit: Int = 20,
        @Query("offset") offset: Int = 0
    ): Response<ApiResponse<com.yonotech.ppob.data.remote.model.StaffListResponse>>

    @POST("users/v1/staff")
    suspend fun addStaff(
        @Header("Authorization") token: String,
        @Body request: com.yonotech.ppob.data.remote.model.AddStaffRequest
    ): Response<ApiResponse<com.yonotech.ppob.data.remote.model.StaffDetailResponse>>

    @PUT("users/v1/staff/{staff_id}")
    suspend fun updateStaff(
        @Header("Authorization") token: String,
        @Path("staff_id") staffId: String,
        @Body request: com.yonotech.ppob.data.remote.model.UpdateStaffRequest
    ): Response<ApiResponse<com.yonotech.ppob.data.remote.model.StaffDetailResponse>>

    @GET("users/v1/staff/{staff_id}")
    suspend fun getStaffDetail(
        @Header("Authorization") token: String,
        @Path("staff_id") staffId: String
    ): Response<ApiResponse<com.yonotech.ppob.data.remote.model.StaffDetailResponse>>

    @GET("users/v1/staff/pending-count")
    suspend fun getStaffPendingCount(
        @Header("Authorization") token: String
    ): Response<ApiResponse<com.yonotech.ppob.data.remote.model.ApiResponse<Int>>>
}

interface ProductApiService {

    @GET("products/v1/categories")
    suspend fun getCategories(): Response<ApiResponse<CategoriesResponse>>

    @GET("products/v1")
    suspend fun getProducts(
        @Query("category_id") categoryId: String,
        @Query("limit") limit: Int = 20,
        @Query("cursor") cursor: String? = null
    ): Response<ApiResponse<ProductsResponse>>

    @POST("products/v1/sync/prepaid")
    suspend fun syncPrepaid(
        @Header("Authorization") token: String
    ): Response<ApiResponse<Unit>>

    @POST("products/v1/sync/postpaid")
    suspend fun syncPostpaid(
        @Header("Authorization") token: String
    ): Response<ApiResponse<Unit>>

    @GET("products/v1/sync/status")
    suspend fun getSyncStatus(): Response<ApiResponse<com.yonotech.ppob.data.remote.model.SyncStatusResponse>>
}

interface WalletApiService {

    @GET("wallets/v1/balance")
    suspend fun getBalance(
        @Header("Authorization") token: String
    ): Response<ApiResponse<WalletBalanceResponse>>

    @POST("wallets/v1/staff/{staff_id}/topup")
    suspend fun topUpStaff(
        @Header("Authorization") token: String,
        @Path("staff_id") staffId: String,
        @Body request: com.yonotech.ppob.data.remote.model.TopUpRequest
    ): Response<ApiResponse<com.yonotech.ppob.data.remote.model.TopUpResponse>>

    @POST("wallets/v1/:id/topup")
    suspend fun topUpSelf(
        @Header("Authorization") token: String,
        @Body request: com.yonotech.ppob.data.remote.model.TopUpRequest
    ): Response<ApiResponse<WalletBalanceResponse>>
}

interface TransactionApiService {

    @Headers("Idempotency-Key: UUID")
    @POST("transactions/v1/initiate")
    suspend fun initiateTransaction(
        @Header("Authorization") token: String,
        @Header("Idempotency-Key") idempotencyKey: String,
        @Body request: com.yonotech.ppob.data.remote.model.InitiateTransactionRequest
    ): Response<ApiResponse<TransactionResponse>>

    @GET("transactions/v1/history")
    suspend fun getTransactionHistory(
        @Header("Authorization") token: String,
        @Query("limit") limit: Int = 20,
        @Query("cursor") cursor: String? = null,
        @Query("status") status: String? = null
    ): Response<ApiResponse<TransactionHistoryResponse>>

    @GET("transactions/v1/{transaction_id}/status")
    suspend fun getTransactionStatus(
        @Header("Authorization") token: String,
        @Path("transaction_id") transactionId: String
    ): Response<ApiResponse<com.yonotech.ppob.data.remote.model.TransactionDetailResponse>>

    @POST("transactions/v1/{transaction_id}/cancel")
    suspend fun cancelTransaction(
        @Header("Authorization") token: String,
        @Path("transaction_id") transactionId: String
    ): Response<ApiResponse<Unit>>
}

interface IntegrationApiService {

    @GET("integration/v1/providers")
    suspend fun getProviders(
        @Header("Authorization") token: String
    ): Response<ResponseBody>

    @GET("integration/v1/errors")
    suspend fun getErrorCodes(
        @Header("Authorization") token: String
    ): Response<ResponseBody>
}