import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../models/product.dart';
import '../models/api_models.dart' hide Provider;
import '../services/product_service.dart';
import '../core/api_client.dart';
import '../utils/constants.dart';

final productRepositoryProvider = Provider<ProductRepository>((ref) {
  final dio = ref.read(dioProvider);
  return ProductRepositoryImpl(ProductService(dio));
});

abstract class ProductRepository {
  Future<List<Product>> getProducts({String? category, String? search, int page = 1, int limit = 20});
  Future<Product> getProduct(String productId);
  Future<List<Product>> searchProducts(String query);
  Future<Product> getProductByCode(String code);
  Future<bool> validatePrice(String productId, double sellingPrice);
  Future<List<Product>> syncProductsFromServer();
  Future<void> updateProductPrice(String productId, double sellingPrice);
  Future<List<Product>> getCategories();
}

class ProductRepositoryImpl implements ProductRepository {
  final ProductService _productService;
  final List<Product> _localProducts = [];

  ProductRepositoryImpl(this._productService) {
    // Seed with some products for offline support
    _seedLocalProducts();
  }

  void _seedLocalProducts() {
    _localProducts.addAll([
      Product(
        id: 'pulsa_001',
        name: 'Pulsa XL 10 Ribu',
        description: 'Pulsa XL 10.000',
        category: 'pulsa',
        productCode: 'XL10',
        platformPrice: 9700.0,
        sellingPrice: 10000.0,
        isActive: true,
        lastSyncedAt: DateTime.now(),
      ),
      Product(
        id: 'pulsa_002',
        name: 'Pulsa XL 25 Ribu',
        description: 'Pulsa XL 25.000',
        category: 'pulsa',
        productCode: 'XL25',
        platformPrice: 24200.0,
        sellingPrice: 25000.0,
        isActive: true,
        lastSyncedAt: DateTime.now(),
      ),
      Product(
        id: 'pulsa_003',
        name: 'Pulsa XL 50 Ribu',
        description: 'Pulsa XL 50.000',
        category: 'pulsa',
        productCode: 'XL50',
        platformPrice: 48500.0,
        sellingPrice: 50000.0,
        isActive: true,
        lastSyncedAt: DateTime.now(),
      ),
      Product(
        id: 'token_001',
        name: 'Token Listrik 20 Ribu',
        description: 'Token Listrik PLN 20.000',
        category: 'token-listrik',
        productCode: 'TOK20',
        platformPrice: 19500.0,
        sellingPrice: 20000.0,
        isActive: true,
        lastSyncedAt: DateTime.now(),
      ),
      Product(
        id: 'token_002',
        name: 'Token Listrik 50 Ribu',
        description: 'Token Listrik PLN 50.000',
        category: 'token-listrik',
        productCode: 'TOK50',
        platformPrice: 48500.0,
        sellingPrice: 50000.0,
        isActive: true,
        lastSyncedAt: DateTime.now(),
      ),
      Product(
        id: 'data_001',
        name: 'Paket Data XL 10GB',
        description: 'Paket data XL 10GB 30 hari',
        category: 'paket-data',
        productCode: 'XLD10',
        platformPrice: 29000.0,
        sellingPrice: 30000.0,
        isActive: true,
        lastSyncedAt: DateTime.now(),
      ),
      Product(
        id: 'data_002',
        name: 'Paket Data XL 20GB',
        description: 'Paket data XL 20GB 30 hari',
        category: 'paket-data',
        productCode: 'XLD20',
        platformPrice: 58000.0,
        sellingPrice: 60000.0,
        isActive: true,
        lastSyncedAt: DateTime.now(),
      ),
    ]);
  }

  @override
  Future<List<Product>> getProducts({
    String? category,
    String? search,
    int page = 1,
    int limit = 20,
  }) async {
    try {
      // Try to fetch from server
      final products = await _productService.getProducts(
        category: category,
        search: search,
        page: page,
        limit: limit,
      );
      // Update local cache
      _updateLocalCache(products);
      return products;
    } catch (e) {
      // Fall back to local data if network fails
      List<Product> result = _localProducts;
      if (category != null) {
        result = result.where((p) => p.category == category).toList();
      }
      if (search != null && search.isNotEmpty) {
        result = result.where((p) => p.name.toLowerCase().contains(search.toLowerCase())).toList();
      }
      return result.skip((page - 1) * limit).take(limit).toList();
    }
  }

  @override
  Future<Product> getProduct(String productId) async {
    try {
      return await _productService.getProduct(productId);
    } catch (e) {
      // Fall back to local
      final product = _localProducts.firstWhere((p) => p.id == productId);
      return product;
    }
  }

  @override
  Future<List<Product>> searchProducts(String query) async {
    try {
      return await _productService.searchProducts(query);
    } catch (e) {
      final result = _localProducts
          .where((p) => p.name.toLowerCase().contains(query.toLowerCase()))
          .toList();
      return result;
    }
  }

  @override
  Future<Product> getProductByCode(String code) async {
    try {
      return await _productService.getProductByCode(code);
    } catch (e) {
      final product = _localProducts.firstWhere((p) => p.productCode == code);
      return product;
    }
  }

  @override
  Future<bool> validatePrice(String productId, double sellingPrice) async {
    try {
      return await _productService.validatePrice(productId, sellingPrice);
    } catch (e) {
      // For offline, check against local product's platform price
      final product = _localProducts.firstWhere((p) => p.id == productId);
      return sellingPrice >= product.platformPrice;
    }
  }

  @override
  Future<List<Product>> syncProductsFromServer() async {
    try {
      // Sync from Digiflazz
      await _productService.syncPrepaidProducts();
      await _productService.syncPostpaidProducts();
      // Then fetch updated list
      final products = await _productService.getProducts(limit: 100);
      _updateLocalCache(products);
      return products;
    } catch (e) {
      throw Exception('Failed to sync products: $e');
    }
  }

  @override
  Future<void> updateProductPrice(String productId, double sellingPrice) async {
    // Update local cache
    final index = _localProducts.indexWhere((p) => p.id == productId);
    if (index >= 0) {
      _localProducts[index] = _localProducts[index].copyWith(sellingPrice: sellingPrice);
    }

    // TODO: Optionally sync to server if admin
    // await _productService.updateProduct(productId, updatedProduct);
  }

  @override
  Future<List<Product>> getCategories() async {
    try {
      final categories = await _productService.getCategories();
      // Convert to Product-like objects for UI
      return categories
          .map((cat) => Product(
                id: cat.id,
                name: cat.name,
                description: cat.description ?? '',
                category: cat.name.toLowerCase(),
                productCode: cat.id,
                platformPrice: 0,
                isActive: true,
                lastSyncedAt: DateTime.now(),
              ))
          .toList();
    } catch (e) {
      // Fallback: extract unique categories from local products
      final categories = _localProducts
          .map((p) => p.category)
          .toSet()
          .map((cat) => Product(
                id: cat,
                name: cat,
                description: '',
                category: cat,
                productCode: cat,
                platformPrice: 0,
                isActive: true,
                lastSyncedAt: DateTime.now(),
              ))
          .toList();
      return categories;
    }
  }

  void _updateLocalCache(List<Product> products) {
    for (final product in products) {
      final index = _localProducts.indexWhere((p) => p.id == product.id);
      if (index >= 0) {
        _localProducts[index] = product;
      } else {
        _localProducts.add(product);
      }
    }
  }
}
