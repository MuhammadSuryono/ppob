package com.yonotech.ppob.data.repository

import android.util.Log
import com.yonotech.ppob.data.local.dao.CategoryDao
import com.yonotech.ppob.data.local.dao.ProductDao
import com.yonotech.ppob.data.local.entity.CategoryEntity
import com.yonotech.ppob.data.local.entity.ProductEntity
import com.yonotech.ppob.data.remote.api.ProductApiService
import com.yonotech.ppob.data.remote.model.ProductsResponse
import com.yonotech.ppob.domain.model.Category as DomainCategory
import com.yonotech.ppob.domain.model.Product as DomainProduct
import com.yonotech.ppob.domain.repository.ProductRepository
import kotlinx.coroutines.flow.first

class ProductRepositoryImpl(
    private val apiService: ProductApiService,
    private val productDao: ProductDao,
    private val categoryDao: CategoryDao
) : ProductRepository {

    override suspend fun getCategories(): Result<List<DomainCategory>> {
        return try {
            val response = apiService.getCategories()
            if (response.isSuccessful) {
                val body = response.body()
                if (body?.success == true && body.data != null) {
                    val categories = body.data.categories.map {
                        DomainCategory(it.categoryId, it.categoryName, it.iconUrl, it.displayOrder)
                    }
                    // Cache locally
                    cacheCategories(categories)
                    Result.success(categories)
                } else {
                    // Return cached
                    val cached = getCachedCategoriesInternal()
                    if (cached.isNotEmpty()) {
                        Result.success(cached)
                    } else {
                        Result.failure(Exception(body?.message ?: "Failed to get categories"))
                    }
                }
            } else {
                val cached = getCachedCategoriesInternal()
                if (cached.isNotEmpty()) Result.success(cached)
                else Result.failure(Exception("HTTP ${response.code()}"))
            }
        } catch (e: Exception) {
            Log.e("ProductRepository", "Get categories error", e)
            val cached = getCachedCategoriesInternal()
            if (cached.isNotEmpty()) Result.success(cached)
            else Result.failure(e)
        }
    }

    override suspend fun getProducts(categoryId: String, limit: Int, cursor: String?): Result<ProductsResponse> {
        return try {
            val response = apiService.getProducts(categoryId, limit, cursor)
            if (response.isSuccessful) {
                val body = response.body()
                if (body?.success == true && body.data != null) {
                    val products = body.data.products.map {
                        DomainProduct(
                            productId = it.productId,
                            buyerSkuCode = it.buyerSkuCode,
                            productName = it.productName,
                            basePrice = it.basePrice,
                            platformPrice = it.platformPrice,
                            isPrepaid = it.isPrepaid,
                            mitraSellingPrice = it.mitraSellingPrice,
                            categoryId = categoryId
                        )
                    }
                    // Cache locally
                    cacheProducts(products, categoryId)
                    Result.success(body.data)
                } else {
                    Result.failure(Exception(body?.message ?: "Failed to get products"))
                }
            } else {
                Result.failure(Exception("HTTP ${response.code()}"))
            }
        } catch (e: Exception) {
            Log.e("ProductRepository", "Get products error", e)
            Result.failure(e)
        }
    }

    override suspend fun getCachedProducts(categoryId: String): List<DomainProduct> {
        return try {
            val entities = productDao.getProductsByCategory(categoryId).first()
            entities.map {
                DomainProduct(
                    productId = it.productId,
                    buyerSkuCode = it.buyerSkuCode,
                    productName = it.productName,
                    basePrice = it.basePrice,
                    platformPrice = it.platformPrice,
                    isPrepaid = it.isPrepaid,
                    mitraSellingPrice = it.mitraSellingPrice,
                    categoryId = it.categoryId
                )
            }
        } catch (e: Exception) {
            Log.e("ProductRepository", "Get cached products error", e)
            emptyList()
        }
    }

    override suspend fun syncProducts(token: String): Result<Unit> {
        return try {
            // Sync prepaid
            val prepaidResponse = apiService.syncPrepaid("Bearer $token")
            // Sync postpaid
            val postpaidResponse = apiService.syncPostpaid("Bearer $token")

            if (prepaidResponse.isSuccessful && postpaidResponse.isSuccessful) {
                Log.d("ProductRepository", "Products synced successfully")
            }
            Result.success(Unit)
        } catch (e: Exception) {
            Log.e("ProductRepository", "Sync products error", e)
            Result.failure(e)
        }
    }

    private suspend fun getCachedCategoriesInternal(): List<DomainCategory> {
        return try {
            val entities = categoryDao.getAllCategories().first()
            entities.map {
                DomainCategory(it.categoryId, it.categoryName, it.iconUrl, it.displayOrder)
            }
        } catch (e: Exception) {
            emptyList()
        }
    }

    private suspend fun cacheCategories(categories: List<DomainCategory>) {
        try {
            val entities = categories.map {
                CategoryEntity(it.categoryId, it.categoryName, it.iconUrl, it.displayOrder)
            }
            categoryDao.insertAll(entities)
        } catch (e: Exception) {
            Log.e("ProductRepository", "Cache categories error", e)
        }
    }

    private suspend fun cacheProducts(products: List<DomainProduct>, categoryId: String) {
        try {
            val entities = products.map {
                ProductEntity(
                    productId = it.productId,
                    buyerSkuCode = it.buyerSkuCode,
                    productName = it.productName,
                    basePrice = it.basePrice,
                    platformPrice = it.platformPrice,
                    isPrepaid = it.isPrepaid,
                    mitraSellingPrice = it.mitraSellingPrice,
                    categoryId = categoryId
                )
            }
            productDao.insertAll(entities)
        } catch (e: Exception) {
            Log.e("ProductRepository", "Cache products error", e)
        }
    }
}
