package com.yonotech.ppob.mobile.data.remote

import com.yonotech.ppob.mobile.data.remote.dto.CategoryDto
import com.yonotech.ppob.mobile.data.remote.dto.ProductDto
import retrofit2.Response
import retrofit2.http.GET
import retrofit2.http.Query

interface ProductService {
    @GET("product/categories")
    suspend fun getCategories(): Response<List<CategoryDto>>

    @GET("product/products")
    suspend fun getProducts(
        @Query("category_id") categoryId: String? = null,
        @Query("brand") brand: String? = null,
        @Query("status") status: String? = "active"
    ): Response<List<ProductDto>>

    @GET("product/products/search")
    suspend fun searchProducts(@Query("q") query: String): Response<List<ProductDto>>
}
