package com.yonotech.ppob.mobile.data.remote

import com.yonotech.ppob.mobile.data.remote.dto.CategoryCollection
import com.yonotech.ppob.mobile.data.remote.dto.ProductCollection
import com.yonotech.ppob.mobile.data.remote.dto.ProductDto
import retrofit2.Response
import retrofit2.http.GET
import retrofit2.http.Query

interface ProductService {
    @GET("categories")
    suspend fun getCategories(): Response<CategoryCollection>

    @GET("products")
    suspend fun getProducts(
        @Query("category_id") categoryId: String? = null,
        @Query("brand") brand: String? = null,
        @Query("status") status: String? = "active"
    ): Response<ProductCollection>

    @GET("products/search")
    suspend fun searchProducts(@Query("q") query: String): Response<ProductCollection>
}
