import 'package:dio/dio.dart';
import '../models/api_models.dart';
import '../core/exceptions.dart';

/// Integration Service API client
/// Gateway for external provider (Digiflazz) communication
/// Port: 8086, Base Path: /api/v1/integrations
class IntegrationService {
  final Dio _dio;
  final String _baseUrl = 'https://192.168.100.23:8086/api/v1/integrations';

  IntegrationService(this._dio);

  /// Forward transaction to Digiflazz API
  /// POST /digiflazz/transaction (requires Bearer token)
  Future<DigiflazzResponse> forwardToDigiflazz(Map<String, dynamic> payload) async {
    try {
      final response = await _dio.post(
        '$_baseUrl/digiflazz/transaction',
        data: payload,
      );
      return DigiflazzResponse.fromJson(response.data as Map<String, dynamic>);
    } on ApiException catch (e) {
      throw e;
    } catch (e) {
      throw ApiException(message: 'Failed to forward to Digiflazz: $e');
    }
  }

  /// List provider configurations
  /// GET /providers (requires Bearer token)
  Future<List<Provider>> listProviders() async {
    try {
      final response = await _dio.get('$_baseUrl/providers');
      final List<dynamic> data = response.data as List<dynamic>;
      return data.map((json) => Provider.fromJson(json as Map<String, dynamic>)).toList();
    } on ApiException catch (e) {
      throw e;
    } catch (e) {
      throw ApiException(message: 'Failed to list providers: $e');
    }
  }

  /// Retrieve mapped error code catalog
  /// GET /errors (requires Bearer token)
  Future<List<ErrorCode>> getErrorCodes({String? provider}) async {
    try {
      final queryParams = <String, dynamic>{};
      if (provider != null) queryParams['provider'] = provider;

      final response = await _dio.get(
        '$_baseUrl/errors',
        queryParameters: queryParams,
      );
      final List<dynamic> data = response.data as List<dynamic>;
      return data.map((json) => ErrorCode.fromJson(json as Map<String, dynamic>)).toList();
    } on ApiException catch (e) {
      throw e;
    } catch (e) {
      throw ApiException(message: 'Failed to get error codes: $e');
    }
  }

  /// List background compensation/retry jobs
  /// GET /compensation/jobs (requires Bearer token, Admin)
  Future<List<CompensationJob>> getCompensationJobs({
    String? status,
    int limit = 20,
    int offset = 0,
  }) async {
    try {
      final queryParams = <String, dynamic>{
        'limit': limit,
        'offset': offset,
      };
      if (status != null) queryParams['status'] = status;

      final response = await _dio.get(
        '$_baseUrl/compensation/jobs',
        queryParameters: queryParams,
      );
      final List<dynamic> data = response.data as List<dynamic>;
      return data.map((json) => CompensationJob.fromJson(json as Map<String, dynamic>)).toList();
    } on ApiException catch (e) {
      throw e;
    } catch (e) {
      throw ApiException(message: 'Failed to get compensation jobs: $e');
    }
  }

  /// List failed jobs in dead letter queue
  /// GET /compensation/dead-letter (requires Bearer token, Admin)
  Future<List<DeadLetterJob>> getDeadLetterJobs({
    int limit = 20,
    int offset = 0,
  }) async {
    try {
      final response = await _dio.get(
        '$_baseUrl/compensation/dead-letter',
        queryParameters: {
          'limit': limit,
          'offset': offset,
        },
      );
      final List<dynamic> data = response.data as List<dynamic>;
      return data.map((json) => DeadLetterJob.fromJson(json as Map<String, dynamic>)).toList();
    } on ApiException catch (e) {
      throw e;
    } catch (e) {
      throw ApiException(message: 'Failed to get dead letter jobs: $e');
    }
  }

  /// Webhook verification and processing (server-to-server)
  /// POST /webhook/digiflazz (no auth required)
  /// Note: This is for backend webhook handling, not mobile
  Future<void> verifyWebhook(String signature, Map<String, dynamic> payload) async {
    throw UnimplementedError('Webhook verification handled by backend service');
  }
}

/// Digiflazz API response
class DigiflazzResponse {
  final String requestId;
  final String? trxId; // Digiflazz transaction ID
  final String status; // 'pending', 'success', 'failed'
  final String? message;
  final Map<String, dynamic>? data; // additional response data

  DigiflazzResponse({
    required this.requestId,
    this.trxId,
    required this.status,
    this.message,
    this.data,
  });

  factory DigiflazzResponse.fromJson(Map<String, dynamic> json) {
    return DigiflazzResponse(
      requestId: json['request_id'] as String,
      trxId: json['trx_id'] as String?,
      status: json['status'] as String,
      message: json['message'] as String?,
      data: json['data'] as Map<String, dynamic>?,
    );
  }
}

/// Compensation job
class CompensationJob {
  final String id;
  final String type;
  final String status; // 'pending', 'processing', 'completed', 'failed'
  final String? transactionId;
  final int retryCount;
  final DateTime scheduledAt;
  final DateTime? executedAt;
  final String? errorMessage;

  CompensationJob({
    required this.id,
    required this.type,
    required this.status,
    this.transactionId,
    required this.retryCount,
    required this.scheduledAt,
    this.executedAt,
    this.errorMessage,
  });

  factory CompensationJob.fromJson(Map<String, dynamic> json) {
    return CompensationJob(
      id: json['id'] as String,
      type: json['type'] as String,
      status: json['status'] as String,
      transactionId: json['transaction_id'] as String?,
      retryCount: json['retry_count'] as int? ?? 0,
      scheduledAt: DateTime.parse(json['scheduled_at'] as String),
      executedAt: json['executed_at'] != null
          ? DateTime.parse(json['executed_at'] as String)
          : null,
      errorMessage: json['error_message'] as String?,
    );
  }
}

/// Dead letter job
class DeadLetterJob {
  final String id;
  final String originalJobId;
  final String reason;
  final Map<String, dynamic>? payload;
  final DateTime failedAt;
  final int retryCount;

  DeadLetterJob({
    required this.id,
    required this.originalJobId,
    required this.reason,
    this.payload,
    required this.failedAt,
    required this.retryCount,
  });

  factory DeadLetterJob.fromJson(Map<String, dynamic> json) {
    return DeadLetterJob(
      id: json['id'] as String,
      originalJobId: json['original_job_id'] as String,
      reason: json['reason'] as String,
      payload: json['payload'] as Map<String, dynamic>?,
      failedAt: DateTime.parse(json['failed_at'] as String),
      retryCount: json['retry_count'] as int? ?? 0,
    );
  }
}
