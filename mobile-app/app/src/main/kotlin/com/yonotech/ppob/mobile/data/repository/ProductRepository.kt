package com.yonotech.ppob.mobile.data.repository

import com.yonotech.ppob.mobile.data.remote.ProductService
import com.yonotech.ppob.mobile.data.remote.dto.CategoryDto
import com.yonotech.ppob.mobile.data.remote.dto.ProductDto
import retrofit2.Response
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class ProductRepository @Inject constructor(
    private val productService: ProductService
) {
    suspend fun getCategories(): Response<List<CategoryDto>> {
        return productService.getCategories()
    }

    suspend fun getProducts(categoryId: String? = null, brand: String? = null): Response<List<ProductDto>> {
        return productService.getProducts(categoryId, brand)
    }

    suspend fun searchProducts(query: String): Response<List<ProductDto>> {
        return productService.searchProducts(query)
    }
}
