package com.yonotech.ppob

import app.cash.turbine.test
import com.yonotech.ppob.domain.model.Product
import com.yonotech.ppob.domain.usecase.GetCategoriesUseCase
import com.yonotech.ppob.domain.usecase.GetProductsUseCase
import com.yonotech.ppob.domain.usecase.InitiateTransactionUseCase
import com.yonotech.ppob.presentation.transaction.CategoryViewModel
import com.yonotech.ppob.presentation.transaction.ProductCatalogViewModel
import io.mockk.coEvery
import io.mockk.mockk
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.ExperimentalCoroutinesApi
import kotlinx.coroutines.test.*
import org.junit.After
import org.junit.Assert.*
import org.junit.Before
import org.junit.Test
import com.yonotech.ppob.data.remote.model.Product as RemoteProduct
import com.yonotech.ppob.data.remote.model.CategoriesResponse
import com.yonotech.ppob.data.remote.model.ProductsResponse
import com.yonotech.ppob.data.remote.model.Pagination

@OptIn(ExperimentalCoroutinesApi::class)
class TransactionUseCasesTest {

    private lateinit var getCategoriesUseCase: GetCategoriesUseCase
    private lateinit var getProductsUseCase: GetProductsUseCase
    private lateinit var initiateTransactionUseCase: InitiateTransactionUseCase
    private lateinit var mockProductRepository: com.yonotech.ppob.domain.repository.ProductRepository
    private lateinit var mockTransactionRepository: com.yonotech.ppob.domain.repository.TransactionRepository

    private lateinit var categoryViewModel: CategoryViewModel
    private lateinit var productCatalogViewModel: ProductCatalogViewModel

    private val testCategories = listOf(
        com.yonotech.ppob.domain.model.Category("cat1", "Pulsa", "https://cdn.ppob.co.id/icons/pulsa.png", 1),
        com.yonotech.ppob.domain.model.Category("cat2", "PLN", "https://cdn.ppob.co.id/icons/pln.png", 2)
    )

    private val testProducts = listOf(
        Product("prod1", "XL10", "XL Data 10GB", 10000.0, 11000.0, true, 12000.0, "cat1"),
        Product("prod2", "XL25", "XL Data 25GB", 25000.0, 26250.0, true, 27000.0, "cat1")
    )

    @Before
    fun setUp() {
        Dispatchers.setMain(StandardTestDispatcher())

        mockProductRepository = mockk(relaxed = true)
        mockTransactionRepository = mockk(relaxed = true)

        getCategoriesUseCase = GetCategoriesUseCase(mockProductRepository)
        getProductsUseCase = GetProductsUseCase(mockProductRepository)
        initiateTransactionUseCase = InitiateTransactionUseCase(mockTransactionRepository)

        categoryViewModel = CategoryViewModel(mockProductRepository)
        productCatalogViewModel = ProductCatalogViewModel(mockProductRepository)
    }

    @After
    fun tearDown() {
        Dispatchers.resetMain()
    }

    @Test
    fun `getCategoriesUseCase should return categories on success`() = runTest {
        coEvery { mockProductRepository.getCategories() } returns Result.success(testCategories)

        val result = getCategoriesUseCase()

        assertTrue(result.isSuccess)
        assertEquals(2, result.getOrNull()?.size)
        assertEquals("Pulsa", result.getOrNull()?.first()?.categoryName)
    }

    @Test
    fun `getCategoriesUseCase should return failure on error`() = runTest {
        coEvery { mockProductRepository.getCategories() } returns Result.failure(Exception("Network error"))

        val result = getCategoriesUseCase()

        assertFalse(result.isSuccess)
        assertEquals("Network error", result.exceptionOrNull()?.message)
    }

    @Test
    fun `getProductsUseCase should return products on success`() = runTest {
        val response = ProductsResponse(
            products = listOf(
                RemoteProduct("prod1", "XL10", "XL Data 10GB", 10000.0, 11000.0, true, 12000.0),
                RemoteProduct("prod2", "XL25", "XL Data 25GB", 25000.0, 26250.0, true, 27000.0)
            ),
            pagination = Pagination(2, 20, 0, false)
        )
        coEvery { mockProductRepository.getProducts("cat1", 20, null) } returns Result.success(response)

        val result = getProductsUseCase("cat1")

        assertTrue(result.isSuccess)
        assertEquals(2, result.getOrNull()?.products?.size)
    }

    @Test
    fun `getProductsUseCase should return failure on error`() = runTest {
        coEvery { mockProductRepository.getProducts(any(), any(), any()) } returns
            Result.failure(Exception("Failed to load"))

        val result = getProductsUseCase("cat1")

        assertFalse(result.isSuccess)
    }

    @Test
    fun `CategoryViewModel should load categories on init`() = runTest {
        coEvery { mockProductRepository.getCategories() } returns Result.success(testCategories)

        categoryViewModel.loadCategories()
        advanceUntilIdle()

        val state = categoryViewModel.uiState.value
        assertFalse(state.isLoading)
        assertEquals(2, state.categories.size)
        assertNull(state.error)
    }

    @Test
    fun `CategoryViewModel should set error on failure`() = runTest {
        coEvery { mockProductRepository.getCategories() } returns
            Result.failure(Exception("Network error"))

        categoryViewModel.loadCategories()
        advanceUntilIdle()

        val state = categoryViewModel.uiState.value
        assertFalse(state.isLoading)
        assertTrue(state.categories.isEmpty())
        assertEquals("Network error", state.error)
    }

    @Test
    fun `CategoryViewModel should clear error`() {
        categoryViewModel.uiState.value = com.yonotech.ppob.presentation.transaction.CategoryUiState(
            error = "Some error"
        )
        categoryViewModel.clearError()
        assertNull(categoryViewModel.uiState.value.error)
    }

    @Test
    fun `ProductCatalogViewModel should load products`() = runTest {
        val response = ProductsResponse(
            products = listOf(
                RemoteProduct("prod1", "XL10", "XL Data 10GB", 10000.0, 11000.0, true, 12000.0)
            ),
            pagination = Pagination(1, 20, 0, false)
        )
        coEvery { mockProductRepository.getProducts("cat1", 20, null) } returns Result.success(response)

        productCatalogViewModel.loadProducts("cat1")
        advanceUntilIdle()

        val state = productCatalogViewModel.uiState.value
        assertFalse(state.isLoading)
        assertEquals(1, state.products.size)
        assertNull(state.error)
    }

    @Test
    fun `ProductCatalogViewModel should handle pagination`() = runTest {
        val page1 = ProductsResponse(
            products = listOf(
                RemoteProduct("prod1", "XL10", "XL Data 10GB", 10000.0, 11000.0, true, 12000.0)
            ),
            pagination = Pagination(2, 1, 0, true)
        )
        val page2 = ProductsResponse(
            products = listOf(
                RemoteProduct("prod2", "XL25", "XL Data 25GB", 25000.0, 26250.0, true, 27000.0)
            ),
            pagination = Pagination(2, 1, 1, false)
        )

        coEvery { mockProductRepository.getProducts("cat1", 20, null) } returns Result.success(page1)
        coEvery { mockProductRepository.getProducts("cat1", 20, "cursor1") } returns Result.success(page2)

        productCatalogViewModel.loadProducts("cat1")
        advanceUntilIdle()

        assertEquals(1, productCatalogViewModel.uiState.value.products.size)

        productCatalogViewModel.loadMore()
        advanceUntilIdle()

        assertEquals(2, productCatalogViewModel.uiState.value.products.size)
    }

    @Test
    fun `ProductCatalogViewModel should search products`() = runTest {
        val response = ProductsResponse(
            products = listOf(
                RemoteProduct("prod1", "XL10", "XL Data 10GB", 10000.0, 11000.0, true, 12000.0),
                RemoteProduct("prod2", "XL25", "XL Data 25GB", 25000.0, 26250.0, true, 27000.0)
            ),
            pagination = null
        )
        coEvery { mockProductRepository.getProducts("cat1", 20, null) } returns Result.success(response)

        productCatalogViewModel.loadProducts("cat1")
        advanceUntilIdle()

        productCatalogViewModel.searchProducts("25GB")
        val filtered = productCatalogViewModel.uiState.value.products.filter {
            it.productName.contains("25GB", ignoreCase = true)
        }
        assertEquals(1, filtered.size)
        assertEquals("XL Data 25GB", filtered.first().productName)
    }

    @Test
    fun `initiateTransactionUseCase should succeed with valid input`() = runTest {
        val response = com.yonotech.ppob.data.remote.model.TransactionResponse(
            transactionId = "txn_123",
            refId = "ref_123",
            status = "Success",
            message = "Transaksi berhasil",
            sellingPrice = 27000.0,
            platformPrice = 26250.0
        )
        coEvery {
            mockTransactionRepository.initiateTransaction(any(), any(), any(), any(), any(), any())
        } returns Result.success(response)

        val result = initiateTransactionUseCase("token", "productId", "081234567890", 27000.0, "123456")

        assertTrue(result.isSuccess)
        assertEquals("txn_123", result.getOrNull()?.transactionId)
    }

    @Test
    fun `initiateTransactionUseCase should fail on insufficient balance`() = runTest {
        coEvery {
            mockTransactionRepository.initiateTransaction(any(), any(), any(), any(), any(), any())
        } returns Result.failure(Exception("INSUFFICIENT_BALANCE"))

        val result = initiateTransactionUseCase("token", "productId", "081234567890", 27000.0, "123456")

        assertFalse(result.isSuccess)
        assertTrue(result.exceptionOrNull()?.message?.contains("INSUFFICIENT_BALANCE") == true)
    }

    @Test
    fun `initiateTransactionUseCase should fail on validation error`() = runTest {
        coEvery {
            mockTransactionRepository.initiateTransaction(any(), any(), any(), any(), any(), any())
        } returns Result.failure(Exception("VALIDATION_CUSTOMER_NUMBER_INVALID"))

        val result = initiateTransactionUseCase("token", "productId", "invalid", 27000.0, "123456")

        assertFalse(result.isSuccess)
    }
}

fun formatRupiah(amount: Double): String {
    return "Rp %,d".format(amount.toInt())
}