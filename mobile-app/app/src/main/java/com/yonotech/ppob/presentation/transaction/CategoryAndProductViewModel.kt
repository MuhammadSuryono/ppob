package com.yonotech.ppob.presentation.transaction

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.yonotech.ppob.data.remote.model.Product
import com.yonotech.ppob.domain.model.Category
import com.yonotech.ppob.domain.model.Product as DomainProduct
import com.yonotech.ppob.domain.usecase.GetCategoriesUseCase
import com.yonotech.ppob.domain.usecase.GetProductsUseCase
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.StateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

data class CategoryUiState(
    val categories: List<Category> = emptyList(),
    val isLoading: Boolean = true,
    val error: String? = null
)

@HiltViewModel
class CategoryViewModel @Inject constructor(
    private val getCategoriesUseCase: GetCategoriesUseCase
) : ViewModel() {

    private val _uiState = MutableStateFlow(CategoryUiState())
    val uiState: StateFlow<CategoryUiState> = _uiState.asStateFlow()

    init {
        loadCategories()
    }

    fun loadCategories() {
        viewModelScope.launch {
            _uiState.value = _uiState.value.copy(isLoading = true, error = null)
            try {
                val result = getCategoriesUseCase()
                if (result.isSuccess) {
                    _uiState.value = _uiState.value.copy(
                        categories = result.getOrDefault(emptyList()),
                        isLoading = false
                    )
                } else {
                    _uiState.value = _uiState.value.copy(
                        isLoading = false,
                        error = result.exceptionOrNull()?.message ?: "Gagal memuat kategori"
                    )
                }
            } catch (e: Exception) {
                _uiState.value = _uiState.value.copy(
                    isLoading = false,
                    error = e.message ?: "Gagal memuat kategori"
                )
            }
        }
    }

    fun clearError() {
        _uiState.value = _uiState.value.copy(error = null)
    }
}

data class ProductUiState(
    val products: List<DomainProduct> = emptyList(),
    val categoryId: String = "",
    val isLoading: Boolean = true,
    val error: String? = null,
    val searchQuery: String = "",
    val hasMore: Boolean = true,
    val cursor: String? = null
)

@HiltViewModel
class ProductCatalogViewModel @Inject constructor(
    private val getProductsUseCase: GetProductsUseCase
) : ViewModel() {

    private val _uiState = MutableStateFlow(ProductUiState())
    val uiState: StateFlow<ProductUiState> = _uiState.asStateFlow()

    fun loadProducts(categoryId: String, refresh: Boolean = false) {
        viewModelScope.launch {
            if (refresh || _uiState.value.categoryId != categoryId) {
                _uiState.value = ProductUiState(
                    categoryId = categoryId,
                    isLoading = true,
                    error = null
                )
            }

            try {
                val result = getProductsUseCase(
                    categoryId,
                    limit = 20,
                    cursor = if (refresh) null else _uiState.value.cursor
                )

                if (result.isSuccess) {
                    val response = result.getOrNull()
                    if (response != null) {
                        val products = response.products.map { product ->
                            DomainProduct(
                                productId = product.productId,
                                buyerSkuCode = product.buyerSkuCode,
                                productName = product.productName,
                                basePrice = product.basePrice,
                                platformPrice = product.platformPrice,
                                isPrepaid = product.isPrepaid,
                                mitraSellingPrice = product.mitraSellingPrice,
                                categoryId = categoryId
                            )
                        }

                        _uiState.value = _uiState.value.copy(
                            products = if (refresh) products else _uiState.value.products + products,
                            hasMore = response.pagination?.hasMore ?: false,
                            cursor = response.pagination?.nextCursor,
                            isLoading = false,
                            error = null
                        )
                    }
                } else {
                    _uiState.value = _uiState.value.copy(
                        isLoading = false,
                        error = result.exceptionOrNull()?.message ?: "Gagal memuat produk"
                    )
                }
            } catch (e: Exception) {
                _uiState.value = _uiState.value.copy(
                    isLoading = false,
                    error = e.message ?: "Gagal memuat produk"
                )
            }
        }
    }

    fun searchProducts(query: String) {
        _uiState.value = _uiState.value.copy(searchQuery = query)
        // Filter locally when search is active
    }

    fun loadMore() {
        val currentState = _uiState.value
        if (!currentState.isLoading && currentState.hasMore) {
            loadProducts(currentState.categoryId)
        }
    }

    fun clearError() {
        _uiState.value = _uiState.value.copy(error = null)
    }
}