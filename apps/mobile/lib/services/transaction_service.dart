import 'package:dio/dio.dart';
import '../models/transaction.dart';
import '../models/product.dart';
import '../models/api_models.dart';
import '../core/exceptions.dart';

/// Transaction Service API client
/// Orchestrates the lifecycle of a PPOB transaction
/// Port: 8084, Base Path: /api/v1/transactions
class TransactionService {
  final Dio _dio;
  final String _baseUrl = 'https://fedora.sinauplatform.id/api/v1/transaction';

  TransactionService(this._dio);

  /// Start a new transaction
  /// POST /initiate (requires Bearer token, Idempotency-Key header required)
  Future<Transaction> initiateTransaction(
    TransactionInitiateRequest request, {
    required String idempotencyKey,
  }) async {
    try {
      final response = await _dio.post(
        '$_baseUrl/initiate',
        data: request.toJson(),
        options: Options(
          headers: {'Idempotency-Key': idempotencyKey},
        ),
      );
      return Transaction.fromJson(response.data as Map<String, dynamic>);
    } on ApiException catch (e) {
      throw e;
    } catch (e) {
      throw ApiException(message: 'Failed to initiate transaction: $e');
    }
  }

  /// Get transaction status by internal ID
  /// GET /:id (requires Bearer token)
  Future<Transaction> getTransaction(String transactionId) async {
    try {
      final response = await _dio.get('$_baseUrl/$transactionId');
      return Transaction.fromJson(response.data as Map<String, dynamic>);
    } on ApiException catch (e) {
      throw e;
    } catch (e) {
      throw ApiException(message: 'Failed to get transaction: $e');
    }
  }

  /// Get transaction by Transaction UUID
  /// GET /by-id/:id (requires Bearer token)
  Future<Transaction> getTransactionByUuid(String uuid) async {
    try {
      final response = await _dio.get('$_baseUrl/by-id/$uuid');
      return Transaction.fromJson(response.data as Map<String, dynamic>);
    } on ApiException catch (e) {
      throw e;
    } catch (e) {
      throw ApiException(message: 'Failed to get transaction: $e');
    }
  }

  /// List transactions with filters
  /// GET / (requires Bearer token)
  Future<PaginatedResponse<Transaction>> listTransactions({
    int page = 1,
    int limit = 20,
    String? status,
    String? productId,
    String? userId,
    DateTime? from,
    DateTime? to,
  }) async {
    try {
      final queryParams = <String, dynamic>{
        'page': page,
        'limit': limit,
      };
      if (status != null) queryParams['status'] = status;
      if (productId != null) queryParams['product_id'] = productId;
      if (userId != null) queryParams['user_id'] = userId;
      if (from != null) queryParams['from'] = from.toIso8601String();
      if (to != null) queryParams['to'] = to.toIso8601String();

      final response = await _dio.get(
        _baseUrl,
        queryParameters: queryParams,
      );

      final data = response.data as Map<String, dynamic>;
      return PaginatedResponse<Transaction>.fromJson(
        data,
        (item) => Transaction.fromJson(item as Map<String, dynamic>),
      );
    } on ApiException catch (e) {
      throw e;
    } catch (e) {
      throw ApiException(message: 'Failed to list transactions: $e');
    }
  }

  /// Get paginated transaction history
  /// GET /history (requires Bearer token)
  Future<PaginatedResponse<Transaction>> getTransactionHistory({
    required String userId,
    int limit = 20,
    String? cursor,
    String? status,
  }) async {
    try {
      final queryParams = <String, dynamic>{
        'user_id': userId,
        'limit': limit,
      };
      if (cursor != null) queryParams['cursor'] = cursor;
      if (status != null) queryParams['status'] = status;

      final response = await _dio.get(
        '$_baseUrl/history',
        queryParameters: queryParams,
      );

      final data = response.data as Map<String, dynamic>;
      return PaginatedResponse<Transaction>.fromJson(
        data,
        (item) => Transaction.fromJson(item as Map<String, dynamic>),
      );
    } on ApiException catch (e) {
      throw e;
    } catch (e) {
      throw ApiException(message: 'Failed to get transaction history: $e');
    }
  }

  /// Manually update transaction status (admin only)
  /// POST /:id/status (requires Bearer token, Admin/Staff)
  Future<Transaction> updateTransactionStatus(
    String transactionId,
    String status, {
    String? notes,
  }) async {
    try {
      final response = await _dio.post(
        '$_baseUrl/$transactionId/status',
        data: {
          'status': status,
          if (notes != null) 'notes': notes,
        },
      );
      return Transaction.fromJson(response.data as Map<String, dynamic>);
    } on ApiException catch (e) {
      throw e;
    } catch (e) {
      throw ApiException(message: 'Failed to update transaction status: $e');
    }
  }

  /// Cancel a pending transaction
  /// POST /:id/cancel (requires Bearer token)
  Future<Transaction> cancelTransaction(String transactionId, {String? reason}) async {
    try {
      final response = await _dio.post(
        '$_baseUrl/$transactionId/cancel',
        data: {
          if (reason != null) 'reason': reason,
        },
      );
      return Transaction.fromJson(response.data as Map<String, dynamic>);
    } on ApiException catch (e) {
      throw e;
    } catch (e) {
      throw ApiException(message: 'Failed to cancel transaction: $e');
    }
  }

  /// Webhook endpoint for Digiflazz updates (called by Digiflazz)
  /// POST /webhook/digiflazz (no auth required)
  /// Note: This is for server-to-server calls, not used by mobile app directly
  Future<void> verifyWebhookSignature(String signature, Map<String, dynamic> payload) async {
    // This would be used by backend, not mobile app
    throw UnimplementedError('Webhook verification is handled by backend');
  }
}
