/// API error response model matching backend error format

class ErrorResponse {
  final ApiError error;

  ErrorResponse({required this.error});

  factory ErrorResponse.fromJson(Map<String, dynamic> json) {
    return ErrorResponse(
      error: ApiError.fromJson(json['error'] as Map<String, dynamic>),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'error': error.toJson(),
    };
  }
}

class ApiError {
  final String code;
  final String message;
  final Map<String, dynamic>? details;
  final String? traceId;
  final String timestamp;

  ApiError({
    required this.code,
    required this.message,
    this.details,
    this.traceId,
    required this.timestamp,
  });

  factory ApiError.fromJson(Map<String, dynamic> json) {
    return ApiError(
      code: json['code'] as String,
      message: json['message'] as String,
      details: json['details'] as Map<String, dynamic>?,
      traceId: json['trace_id'] as String?,
      timestamp: json['timestamp'] as String? ?? DateTime.now().toIso8601String(),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'code': code,
      'message': message,
      'details': details,
      'trace_id': traceId,
      'timestamp': timestamp,
    };
  }
}
