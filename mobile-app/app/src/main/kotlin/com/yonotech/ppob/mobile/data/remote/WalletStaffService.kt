package com.yonotech.ppob.mobile.data.remote

import com.yonotech.ppob.mobile.data.remote.dto.CreateStaffRequest
import com.yonotech.ppob.mobile.data.remote.dto.StaffDto
import com.yonotech.ppob.mobile.data.remote.dto.TopUpRequest
import com.yonotech.ppob.mobile.data.remote.dto.TransactionHistoryResponse
import com.yonotech.ppob.mobile.data.remote.dto.WalletResponse
import retrofit2.Response
import retrofit2.http.Body
import retrofit2.http.GET
import retrofit2.http.POST
import retrofit2.http.Path
import retrofit2.http.Query

interface WalletService {
    @GET("wallets/v1/balance")
    suspend fun getBalance(): Response<WalletResponse>
}

interface StaffService {
    @GET("user/v1/staff")
    suspend fun getStaffList(): Response<List<StaffDto>>

    @POST("user/v1/staff")
    suspend fun createStaff(@Body request: CreateStaffRequest): Response<StaffDto>

    @POST("wallets/v1/staff/{id}/topup")
    suspend fun topUpStaff(
        @Path("id") staffId: String,
        @Body request: TopUpRequest
    ): Response<TransactionHistoryResponse>

    @GET("transactions/v1/history")
    suspend fun getTransactionHistory(
        @Query("status") status: String? = null,
        @Query("limit") limit: Int = 50,
        @Query("offset") offset: Int = 0
    ): Response<List<TransactionHistoryResponse>>
}