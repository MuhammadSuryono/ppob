import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../utils/constants.dart';

final dioProvider = Provider<Dio>((ref) {
  final dio = Dio(
    BaseOptions(
      baseUrl: AppConstants.apiBaseUrl,
      connectTimeout: const Duration(seconds: 30),
      receiveTimeout: const Duration(seconds: 30),
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
      },
    ),
  );

  // Add interceptors
  dio.interceptors.add(
    InterceptorsWrapper(
      onRequest: (options, handler) async {
        // Add auth token if available
        // TODO: Get token from secure storage
        // final token = await ref.read(secureStorageProvider).read(key: AppConstants.keyAuthToken);
        // if (token != null) {
        //   options.headers['Authorization'] = 'Bearer $token';
        // }
        return handler.next(options);
      },
      onResponse: (response, handler) {
        return handler.next(response);
      },
      onError: (error, handler) {
        // Handle common errors
        if (error.response?.statusCode == 401) {
          // Token expired, redirect to login
          // TODO: Clear auth state and navigate to login
        }
        return handler.next(error);
      },
    ),
  );

  return dio;
});

final apiServiceProvider = Provider<ApiService>((ref) {
  final dio = ref.watch(dioProvider);
  return ApiService(dio);
});

class ApiService {
  final Dio _dio;

  ApiService(this._dio);

  // Auth endpoints
  Future<Map<String, dynamic>> login(String username, String password) async {
    try {
      final response = await _dio.post(
        '${AppConstants.endpointAuth}/login',
        data: {'username': username, 'password': password},
      );
      return response.data as Map<String, dynamic>;
    } catch (e) {
      rethrow;
    }
  }

  Future<Map<String, dynamic>> register(Map<String, dynamic> data) async {
    try {
      final response = await _dio.post(
        '${AppConstants.endpointAuth}/register',
        data: data,
      );
      return response.data as Map<String, dynamic>;
    } catch (e) {
      rethrow;
    }
  }

  Future<Map<String, dynamic>> verifyOtp(String phoneNumber, String otp) async {
    try {
      final response = await _dio.post(
        '${AppConstants.endpointAuth}/verify-otp',
        data: {'phoneNumber': phoneNumber, 'otp': otp},
      );
      return response.data as Map<String, dynamic>;
    } catch (e) {
      rethrow;
    }
  }

  // User endpoints
  Future<Map<String, dynamic>> getUserProfile(String userId) async {
    try {
      final response = await _dio.get('${AppConstants.endpointUsers}/$userId');
      return response.data as Map<String, dynamic>;
    } catch (e) {
      rethrow;
    }
  }

  // Product endpoints
  Future<List<dynamic>> getProducts({String? category, String? search}) async {
    try {
      final queryParameters = <String, dynamic>{};
      if (category != null) queryParameters['category'] = category;
      if (search != null) queryParameters['search'] = search;

      final response = await _dio.get(
        AppConstants.endpointProducts,
        queryParameters: queryParameters,
      );
      return response.data as List<dynamic>;
    } catch (e) {
      rethrow;
    }
  }

  Future<Map<String, dynamic>> getProduct(String productId) async {
    try {
      final response = await _dio.get('${AppConstants.endpointProducts}/$productId');
      return response.data as Map<String, dynamic>;
    } catch (e) {
      rethrow;
    }
  }

  // Wallet endpoints
  Future<Map<String, dynamic>> getWallet(String walletId) async {
    try {
      final response = await _dio.get('${AppConstants.endpointWallets}/$walletId');
      return response.data as Map<String, dynamic>;
    } catch (e) {
      rethrow;
    }
  }

  Future<Map<String, dynamic>> topUpStaff(String staffId, double amount) async {
    try {
      final response = await _dio.post(
        '${AppConstants.endpointWallets}/staff/$staffId/topup',
        data: {'amount': amount},
      );
      return response.data as Map<String, dynamic>;
    } catch (e) {
      rethrow;
    }
  }

  // Transaction endpoints
  Future<Map<String, dynamic>> initiateTransaction(Map<String, dynamic> data) async {
    try {
      final response = await _dio.post(
        '${AppConstants.endpointTransactions}/initiate',
        data: data,
      );
      return response.data as Map<String, dynamic>;
    } catch (e) {
      rethrow;
    }
  }

  Future<List<dynamic>> getTransactionHistory({
    required String userId,
    String? status,
    int limit = 20,
    int offset = 0,
  }) async {
    try {
      final queryParameters = <String, dynamic>{
        'user_id': userId,
        'limit': limit,
        'offset': offset,
      };
      if (status != null) queryParameters['status'] = status;

      final response = await _dio.get(
        AppConstants.endpointTransactions,
        queryParameters: queryParameters,
      );
      return response.data as List<dynamic>;
    } catch (e) {
      rethrow;
    }
  }

  Future<Map<String, dynamic>> getTransactionDetail(String transactionId) async {
    try {
      final response = await _dio.get('${AppConstants.endpointTransactions}/$transactionId');
      return response.data as Map<String, dynamic>;
    } catch (e) {
      rethrow;
    }
  }

  // Staff endpoints (for Mitra)
  Future<List<dynamic>> getStaffList(String mitraId) async {
    try {
      final response = await _dio.get(
        '${AppConstants.endpointStaff}/mitra/$mitraId',
      );
      return response.data as List<dynamic>;
    } catch (e) {
      rethrow;
    }
  }

  Future<Map<String, dynamic>> createStaff(Map<String, dynamic> data) async {
    try {
      final response = await _dio.post(
        AppConstants.endpointStaff,
        data: data,
      );
      return response.data as Map<String, dynamic>;
    } catch (e) {
      rethrow;
    }
  }

  Future<Map<String, dynamic>> updateStaff(String staffId, Map<String, dynamic> data) async {
    try {
      final response = await _dio.put(
        '${AppConstants.endpointStaff}/$staffId',
        data: data,
      );
      return response.data as Map<String, dynamic>;
    } catch (e) {
      rethrow;
    }
  }

  Future<void> deleteStaff(String staffId) async {
    try {
      await _dio.delete('${AppConstants.endpointStaff}/$staffId');
    } catch (e) {
      rethrow;
    }
  }
}
