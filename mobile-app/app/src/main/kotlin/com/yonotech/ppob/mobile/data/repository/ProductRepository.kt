package com.yonotech.ppob.mobile.data.repository

import com.yonotech.ppob.mobile.data.remote.ProductService
import com.yonotech.ppob.mobile.data.remote.dto.CategoryCollection
import com.yonotech.ppob.mobile.data.remote.dto.CategoryDto
import com.yonotech.ppob.mobile.data.remote.dto.ProductCollection
import com.yonotech.ppob.mobile.data.remote.dto.ProductDto
import retrofit2.Response
import javax.inject.Inject
import javax.inject.Singleton

@Singleton
class ProductRepository @Inject constructor(
    private val productService: ProductService
) {
    suspend fun getCategories(): Response<CategoryCollection> {
        return productService.getCategories()
    }

    suspend fun getProducts(categoryId: String? = null, brand: String? = null): Response<ProductCollection> {
        return productService.getProducts(categoryId, brand)
    }

    suspend fun searchProducts(query: String): Response<ProductCollection> {
        return productService.searchProducts(query)
    }
}
