package com.yonotech.ppob.domain.usecase

import com.yonotech.ppob.data.remote.model.ProductsResponse
import com.yonotech.ppob.domain.model.Category
import com.yonotech.ppob.domain.model.Product
import com.yonotech.ppob.domain.repository.ProductRepository

class GetCategoriesUseCase(private val repository: ProductRepository) {
    suspend operator fun invoke(): Result<List<Category>> {
        return repository.getCategories()
    }
}

class GetProductsUseCase(private val repository: ProductRepository) {
    suspend operator fun invoke(categoryId: String, limit: Int = 20, cursor: String? = null): Result<ProductsResponse> {
        return repository.getProducts(categoryId, limit, cursor)
    }

    suspend fun getCachedProducts(categoryId: String): List<Product> {
        return repository.getCachedProducts(categoryId)
    }
}

class SyncProductsUseCase(private val repository: ProductRepository) {
    suspend operator fun invoke(token: String): Result<Unit> {
        return repository.syncProducts(token)
    }
}

class SearchProductsUseCase(private val repository: ProductRepository) {
    suspend operator fun invoke(query: String): Result<List<Product>> {
        // This would search cached products
        return Result.success(emptyList())
    }
}
