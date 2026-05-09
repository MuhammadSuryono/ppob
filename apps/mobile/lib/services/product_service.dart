import 'package:dio/dio.dart';
import '../models/product.dart';
import '../models/api_models.dart';
import '../core/exceptions.dart';

/// Product Service API client
/// Handles product catalog and synchronization with Digiflazz
/// Port: 8085, Base Path: /api/v1
class ProductService {
  final Dio _dio;
  final String _baseUrl = 'https://fedora.sinauplatform.id/api/v1/product';

  ProductService(this._dio);

  /// List products with optional filters
  /// GET /products (no auth required)
  Future<List<Product>> getProducts({
    String? category,
    String? status,
    String? search,
    int page = 1,
    int limit = 20,
  }) async {
    try {
      final queryParams = <String, dynamic>{
        'page': page,
        'limit': limit,
      };
      if (category != null) queryParams['category'] = category;
      if (status != null) queryParams['status'] = status;
      if (search != null) queryParams['search'] = search;

      final response = await _dio.get(
        '$_baseUrl/products',
        queryParameters: queryParams,
      );

      final List<dynamic> data = response.data as List<dynamic>;
      return data.map((json) => Product.fromJson(json as Map<String, dynamic>)).toList();
    } on ApiException catch (e) {
      throw e;
    } catch (e) {
      throw ApiException(message: 'Failed to fetch products: $e');
    }
  }

  /// Get product by ID
  /// GET /products/:id (no auth required)
  Future<Product> getProduct(String productId) async {
    try {
      final response = await _dio.get('$_baseUrl/products/$productId');
      return Product.fromJson(response.data as Map<String, dynamic>);
    } on ApiException catch (e) {
      throw e;
    } catch (e) {
      throw ApiException(message: 'Failed to get product: $e');
    }
  }

  /// Search products by name/code
  /// GET /products/search (no auth required)
  Future<List<Product>> searchProducts(String query) async {
    try {
      final response = await _dio.get(
        '$_baseUrl/products/search',
        queryParameters: {'q': query},
      );
      final List<dynamic> data = response.data as List<dynamic>;
      return data.map((json) => Product.fromJson(json as Map<String, dynamic>)).toList();
    } on ApiException catch (e) {
      throw e;
    } catch (e) {
      throw ApiException(message: 'Failed to search products: $e');
    }
  }

  /// Get product by SKU code
  /// GET /products/by-code/:code (no auth required)
  Future<Product> getProductByCode(String code) async {
    try {
      final response = await _dio.get('$_baseUrl/products/by-code/$code');
      return Product.fromJson(response.data as Map<String, dynamic>);
    } on ApiException catch (e) {
      throw e;
    } catch (e) {
      throw ApiException(message: 'Failed to get product by code: $e');
    }
  }

  /// Check if a price meets margin requirements
  /// GET /products/validate-price (no auth required)
  Future<bool> validatePrice(String productId, double sellingPrice) async {
    try {
      final response = await _dio.get(
        '$_baseUrl/products/validate-price',
        queryParameters: {
          'product_id': productId,
          'selling_price': sellingPrice,
        },
      );
      return response.data as bool;
    } on ApiException catch (e) {
      throw e;
    } catch (e) {
      throw ApiException(message: 'Failed to validate price: $e');
    }
  }

  /// Create product manually (admin only)
  /// POST /products (requires Bearer token, Admin)
  Future<Product> createProduct(Product product) async {
    try {
      final response = await _dio.post(
        '$_baseUrl/products',
        data: product.toJson(),
      );
      return Product.fromJson(response.data as Map<String, dynamic>);
    } on ApiException catch (e) {
      throw e;
    } catch (e) {
      throw ApiException(message: 'Failed to create product: $e');
    }
  }

  /// Update product information
  /// PUT /products/:id (requires Bearer token, Admin)
  Future<Product> updateProduct(String productId, Product product) async {
    try {
      final response = await _dio.put(
        '$_baseUrl/products/$productId',
        data: product.toJson(),
      );
      return Product.fromJson(response.data as Map<String, dynamic>);
    } on ApiException catch (e) {
      throw e;
    } catch (e) {
      throw ApiException(message: 'Failed to update product: $e');
    }
  }

  /// Soft-delete a product
  /// DELETE /products/:id (requires Bearer token, Admin)
  Future<void> deleteProduct(String productId) async {
    try {
      await _dio.delete('$_baseUrl/products/$productId');
    } on ApiException catch (e) {
      throw e;
    } catch (e) {
      throw ApiException(message: 'Failed to delete product: $e');
    }
  }

  /// Trigger prepaid product sync from Digiflazz
  /// POST /sync/prepaid (requires Bearer token)
  Future<SyncStatusResponse> syncPrepaidProducts() async {
    try {
      final response = await _dio.post('$_baseUrl/sync/prepaid');
      return SyncStatusResponse.fromJson(response.data as Map<String, dynamic>);
    } on ApiException catch (e) {
      throw e;
    } catch (e) {
      throw ApiException(message: 'Failed to trigger prepaid sync: $e');
    }
  }

  /// Trigger postpaid product sync from Digiflazz
  /// POST /sync/postpaid (requires Bearer token)
  Future<SyncStatusResponse> syncPostpaidProducts() async {
    try {
      final response = await _dio.post('$_baseUrl/sync/postpaid');
      return SyncStatusResponse.fromJson(response.data as Map<String, dynamic>);
    } on ApiException catch (e) {
      throw e;
    } catch (e) {
      throw ApiException(message: 'Failed to trigger postpaid sync: $e');
    }
  }

  /// Get last sync timestamps
  /// GET /sync/status (no auth required)
  Future<SyncStatusResponse> getSyncStatus() async {
    try {
      final response = await _dio.get('$_baseUrl/sync/status');
      return SyncStatusResponse.fromJson(response.data as Map<String, dynamic>);
    } on ApiException catch (e) {
      throw e;
    } catch (e) {
      throw ApiException(message: 'Failed to get sync status: $e');
    }
  }

  /// List all categories
  /// GET /categories (no auth required)
  Future<List<CategoryResponse>> getCategories() async {
    try {
      final response = await _dio.get('$_baseUrl/categories');
      final List<dynamic> data = response.data as List<dynamic>;
      return data.map((json) => CategoryResponse.fromJson(json as Map<String, dynamic>)).toList();
    } on ApiException catch (e) {
      throw e;
    } catch (e) {
      throw ApiException(message: 'Failed to get categories: $e');
    }
  }

  /// Create a category (admin only)
  /// POST /categories (requires Bearer token, Admin)
  Future<CategoryResponse> createCategory(String name, String? description) async {
    try {
      final response = await _dio.post(
        '$_baseUrl/categories',
        data: {
          'name': name,
          if (description != null) 'description': description,
        },
      );
      return CategoryResponse.fromJson(response.data as Map<String, dynamic>);
    } on ApiException catch (e) {
      throw e;
    } catch (e) {
      throw ApiException(message: 'Failed to create category: $e');
    }
  }
}

/// Sync status response
class SyncStatusResponse {
  final DateTime? lastPrepaidSyncAt;
  final DateTime? lastPostpaidSyncAt;
  final bool isSyncing;

  SyncStatusResponse({
    this.lastPrepaidSyncAt,
    this.lastPostpaidSyncAt,
    required this.isSyncing,
  });

  factory SyncStatusResponse.fromJson(Map<String, dynamic> json) {
    return SyncStatusResponse(
      lastPrepaidSyncAt: json['last_prepaid_sync_at'] != null
          ? DateTime.parse(json['last_prepaid_sync_at'] as String)
          : null,
      lastPostpaidSyncAt: json['last_postpaid_sync_at'] != null
          ? DateTime.parse(json['last_postpaid_sync_at'] as String)
          : null,
      isSyncing: json['is_syncing'] as bool? ?? false,
    );
  }
}

/// Category response
class CategoryResponse {
  final String id;
  final String name;
  final String? description;
  final int productCount;

  CategoryResponse({
    required this.id,
    required this.name,
    this.description,
    required this.productCount,
  });

  factory CategoryResponse.fromJson(Map<String, dynamic> json) {
    return CategoryResponse(
      id: json['id'] as String,
      name: json['name'] as String,
      description: json['description'] as String?,
      productCount: json['product_count'] as int? ?? 0,
    );
  }
}
