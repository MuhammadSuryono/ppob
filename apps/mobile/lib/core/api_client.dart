import 'dart:math';
import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import '../utils/constants.dart';
import '../models/error_response.dart';
import '../core/exceptions.dart';

/// Provider for secure storage
final secureStorageProvider = Provider<FlutterSecureStorage>((ref) {
  return const FlutterSecureStorage();
});

/// Provider for Dio instance with interceptors
final dioProvider = Provider<Dio>((ref) {
  final dio = Dio(
    BaseOptions(
      connectTimeout: const Duration(seconds: 30),
      receiveTimeout: const Duration(seconds: 30),
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
      },
    ),
  );

  // Add auth interceptor
  dio.interceptors.add(AuthInterceptor(ref));

  // Add logging interceptor (only in debug mode)
  if (const bool.fromEnvironment('dart.vm.product') == false) {
    dio.interceptors.add(LoggingInterceptor());
  }

  // Add error handling interceptor
  dio.interceptors.add(ErrorInterceptor());

  return dio;
});

/// Auth interceptor - adds JWT token to requests
class AuthInterceptor extends Interceptor {
  final Ref _ref;
  static const _tokenKey = 'auth_token';
  static final Random _random = Random();

  AuthInterceptor(this._ref);

  @override
  Future<void> onRequest(
    RequestOptions options,
    RequestInterceptorHandler handler,
  ) async {
    final storage = _ref.read(secureStorageProvider);
    final token = await storage.read(key: _tokenKey);

    if (token != null && token.isNotEmpty) {
      options.headers['Authorization'] = 'Bearer $token';
    }

    // Add trace ID if not present
    options.headers.putIfAbsent('X-Trace-ID', () => _generateTraceId());

    handler.next(options);
  }

  @override
  void onError(DioException err, ErrorInterceptorHandler handler) {
    if (err.response?.statusCode == 401) {
      // Token expired - clear storage and redirect to login
      _ref.read(secureStorageProvider).delete(key: _tokenKey);
      // Navigate to login screen (will be handled by auth provider)
    }
    handler.next(err);
  }

  String _generateTraceId() {
    final timestamp = DateTime.now().millisecondsSinceEpoch;
    final random = _random.nextInt(10000);
    return 'trace_${timestamp}_$random';
  }
}

/// Logging interceptor for debug mode
class LoggingInterceptor extends Interceptor {
  @override
  void onRequest(RequestOptions options, RequestInterceptorHandler handler) {
    print('📤 REQUEST: ${options.method} ${options.uri}');
    print('Headers: ${options.headers}');
    if (options.data != null) {
      print('Body: ${options.data}');
    }
    handler.next(options);
  }

  @override
  void onResponse(Response response, ResponseInterceptorHandler handler) {
    print('📥 RESPONSE: ${response.statusCode}');
    print('Data: ${response.data}');
    handler.next(response);
  }

  @override
  void onError(DioException err, ErrorInterceptorHandler handler) {
    print('❌ ERROR: ${err.message}');
    if (err.response != null) {
      print('Status: ${err.response?.statusCode}');
      print('Data: ${err.response?.data}');
    }
    handler.next(err);
  }
}

/// Error interceptor - converts errors to ApiException
class ErrorInterceptor extends Interceptor {
  @override
  Future<void> onError(DioException err, ErrorInterceptorHandler handler) async {
    if (err.response != null) {
      final statusCode = err.response!.statusCode ?? 500;
      final data = err.response!.data;

      ApiException exception;

      try {
        if (data is Map<String, dynamic> && data.containsKey('error')) {
          final errorResponse = ErrorResponse.fromJson(data);
          final apiError = errorResponse.error;

          switch (statusCode) {
            case 400:
              if (apiError.code == 'TRANSACTION_INSUFFICIENT_BALANCE') {
                final details = apiError.details;
                exception = InsufficientBalanceException(
                  required: (details?['required'] as num?)?.toDouble() ?? 0.0,
                  current: (details?['current'] as num?)?.toDouble() ?? 0.0,
                  traceId: apiError.traceId,
                );
              } else {
                exception = ValidationException(
                  message: apiError.message,
                  errors: apiError.details,
                  traceId: apiError.traceId,
                );
              }
              break;
            case 401:
              exception = UnauthorizedException(
                message: apiError.message,
                traceId: apiError.traceId,
              );
              break;
            case 403:
              exception = ForbiddenException(
                message: apiError.message,
                traceId: apiError.traceId,
              );
              break;
            case 404:
              exception = NotFoundException(
                message: apiError.message,
                traceId: apiError.traceId,
              );
              break;
            case 429:
              exception = RateLimitException(
                message: apiError.message,
                traceId: apiError.traceId,
              );
              break;
            case 502:
              exception = DigiflazzException(
                message: apiError.message,
                traceId: apiError.traceId,
              );
              break;
            case 503:
              exception = ServiceUnavailableException(
                message: apiError.message,
                traceId: apiError.traceId,
              );
              break;
            default:
              exception = ApiException(
                message: apiError.message,
                statusCode: statusCode,
                code: apiError.code,
                traceId: apiError.traceId,
              );
          }
        } else {
          exception = ApiException(
            message: data?.toString() ?? err.message ?? 'Unknown error',
            statusCode: statusCode,
          );
        }
      } catch (e) {
        exception = ApiException(
          message: 'Failed to parse error response',
          statusCode: statusCode,
        );
      }

      handler.reject(err.copyWith(error: exception));
    } else {
      // Network error or timeout
      final exception = ApiException(
        message: err.message ?? 'Network error occurred',
        statusCode: 0,
      );
      handler.reject(err.copyWith(error: exception));
    }
  }
}
