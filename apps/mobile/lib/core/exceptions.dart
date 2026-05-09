/// Custom exceptions for API errors
class ApiException implements Exception {
  final String message;
  final int? statusCode;
  final String? code;
  final String? traceId;
  final DateTime timestamp;

  ApiException({
    required this.message,
    this.statusCode,
    this.code,
    this.traceId,
    DateTime? timestamp,
  }) : timestamp = timestamp ?? DateTime.now();

  @override
  String toString() {
    return 'ApiException($statusCode: $message) [code: $code, traceId: $traceId]';
  }
}

class UnauthorizedException extends ApiException {
  UnauthorizedException({String? message, String? traceId})
      : super(
          message: message ?? 'Unauthorized',
          statusCode: 401,
          code: 'UNAUTHORIZED',
          traceId: traceId,
        );
}

class ForbiddenException extends ApiException {
  ForbiddenException({String? message, String? traceId})
      : super(
          message: message ?? 'Forbidden',
          statusCode: 403,
          code: 'FORBIDDEN',
          traceId: traceId,
        );
}

class NotFoundException extends ApiException {
  NotFoundException({String? message, String? traceId})
      : super(
          message: message ?? 'Not Found',
          statusCode: 404,
          code: 'NOT_FOUND',
          traceId: traceId,
        );
}

class ValidationException extends ApiException {
  final Map<String, dynamic>? errors;

  ValidationException({String? message, this.errors, String? traceId})
      : super(
          message: message ?? 'Validation failed',
          statusCode: 400,
          code: 'VALIDATION_ERROR',
          traceId: traceId,
        );
}

class InsufficientBalanceException extends ApiException {
  final double required;
  final double current;

  InsufficientBalanceException({required this.required, required this.current, String? traceId})
      : super(
          message: 'Saldo tidak mencukupi',
          statusCode: 400,
          code: 'TRANSACTION_INSUFFICIENT_BALANCE',
          traceId: traceId,
        );
}

class RateLimitException extends ApiException {
  RateLimitException({String? message, String? traceId})
      : super(
          message: message ?? 'Rate limit exceeded',
          statusCode: 429,
          code: 'RATE_LIMIT_EXCEEDED',
          traceId: traceId,
        );
}

class ServiceUnavailableException extends ApiException {
  ServiceUnavailableException({String? message, String? traceId})
      : super(
          message: message ?? 'Service unavailable',
          statusCode: 503,
          code: 'SERVICE_UNAVAILABLE',
          traceId: traceId,
        );
}

class DigiflazzException extends ApiException {
  DigiflazzException({String? message, String? traceId})
      : super(
          message: message ?? 'Digiflazz API error',
          statusCode: 502,
          code: 'DIGIFLAZZ_ERROR',
          traceId: traceId,
        );
}
