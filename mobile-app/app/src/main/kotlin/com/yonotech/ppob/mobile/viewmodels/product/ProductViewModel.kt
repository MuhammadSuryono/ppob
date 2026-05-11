package com.yonotech.ppob.mobile.viewmodels.product

import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.yonotech.ppob.mobile.data.remote.dto.CategoryDto
import com.yonotech.ppob.mobile.data.remote.dto.ProductDto
import com.yonotech.ppob.mobile.data.repository.ProductRepository
import com.yonotech.ppob.mobile.utils.Resource
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.MutableStateFlow
import kotlinx.coroutines.flow.asStateFlow
import kotlinx.coroutines.launch
import javax.inject.Inject

@HiltViewModel
class ProductViewModel @Inject constructor(
    private val repository: ProductRepository
) : ViewModel() {

    private val _categories = MutableStateFlow<Resource<List<CategoryDto>>>(Resource.Idle)
    val categories = _categories.asStateFlow()

    private val _products = MutableStateFlow<Resource<List<ProductDto>>>(Resource.Idle)
    val products = _products.asStateFlow()

    fun getCategories() {
        viewModelScope.launch {
            _categories.value = Resource.Loading
            try {
                val response = repository.getCategories()
                if (response.isSuccessful) {
                    _categories.value = Resource.Success(response.body() ?: emptyList())
                } else {
                    _categories.value = Resource.Error("Failed to fetch categories: ${response.message()}")
                }
            } catch (e: Exception) {
                _categories.value = Resource.Error(e.message ?: "An unknown error occurred")
            }
        }
    }

    fun getProducts(categoryId: String? = null, brand: String? = null) {
        viewModelScope.launch {
            _products.value = Resource.Loading
            try {
                val response = repository.getProducts(categoryId, brand)
                if (response.isSuccessful) {
                    _products.value = Resource.Success(response.body() ?: emptyList())
                } else {
                    _products.value = Resource.Error("Failed to fetch products: ${response.message()}")
                }
            } catch (e: Exception) {
                _products.value = Resource.Error(e.message ?: "An unknown error occurred")
            }
        }
    }

    fun searchProducts(query: String) {
        viewModelScope.launch {
            _products.value = Resource.Loading
            try {
                val response = repository.searchProducts(query)
                if (response.isSuccessful) {
                    _products.value = Resource.Success(response.body() ?: emptyList())
                } else {
                    _products.value = Resource.Error("Failed to search products: ${response.message()}")
                }
            } catch (e: Exception) {
                _products.value = Resource.Error(e.message ?: "An unknown error occurred")
            }
        }
    }
}
