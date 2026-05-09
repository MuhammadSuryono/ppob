/// API response wrapper for paginated results
class PaginatedResponse<T> {
  final List<T> data;
  final int page;
  final int limit;
  final int total;
  final int totalPages;

  PaginatedResponse({
    required this.data,
    required this.page,
    required this.limit,
    required this.total,
    required this.totalPages,
  });

  factory PaginatedResponse.fromJson(
    Map<String, dynamic> json,
    T Function(Map<String, dynamic>) fromJson,
  ) {
    return PaginatedResponse(
      data: (json['data'] as List<dynamic>)
          .map((item) => fromJson(item as Map<String, dynamic>))
          .toList(),
      page: json['page'] as int? ?? 1,
      limit: json['limit'] as int? ?? 20,
      total: json['total'] as int? ?? 0,
      totalPages: json['total_pages'] as int? ?? 1,
    );
  }
}

/// Standard API response for single resources
class ApiResponse<T> {
  final T data;
  final String? message;

  ApiResponse({required this.data, this.message});

  factory ApiResponse.fromJson(
    Map<String, dynamic> json,
    T Function(Map<String, dynamic>) fromJson,
  ) {
    return ApiResponse(
      data: fromJson(json['data'] as Map<String, dynamic>),
      message: json['message'] as String?,
    );
  }
}

/// Balance response from wallet service
class BalanceResponse {
  final String walletId;
  final double availableBalance;
  final double heldBalance;
  final double totalBalance;
  final DateTime updatedAt;

  BalanceResponse({
    required this.walletId,
    required this.availableBalance,
    required this.heldBalance,
    required this.totalBalance,
    required this.updatedAt,
  });

  factory BalanceResponse.fromJson(Map<String, dynamic> json) {
    return BalanceResponse(
      walletId: json['wallet_id'] as String,
      availableBalance: (json['available_balance'] as num).toDouble(),
      heldBalance: (json['held_balance'] as num).toDouble(),
      totalBalance: (json['total_balance'] as num).toDouble(),
      updatedAt: DateTime.parse(json['updated_at'] as String),
    );
  }
}

/// Hold funds request
class HoldFundsRequest {
  final String transactionId;
  final double amount;
  final String? description;

  HoldFundsRequest({
    required this.transactionId,
    required this.amount,
    this.description,
  });

  Map<String, dynamic> toJson() {
    return {
      'transaction_id': transactionId,
      'amount': amount,
      if (description != null) 'description': description,
    };
  }
}

/// Release hold request
class ReleaseHoldRequest {
  final String holdId;
  final String? reason;

  ReleaseHoldRequest({required this.holdId, this.reason});

  Map<String, dynamic> toJson() {
    return {
      'hold_id': holdId,
      if (reason != null) 'reason': reason,
    };
  }
}

/// Transfer request
class TransferRequest {
  final String fromWalletId;
  final String toWalletId;
  final double amount;
  final String? description;

  TransferRequest({
    required this.fromWalletId,
    required this.toWalletId,
    required this.amount,
    this.description,
  });

  Map<String, dynamic> toJson() {
    return {
      'from_wallet_id': fromWalletId,
      'to_wallet_id': toWalletId,
      'amount': amount,
      if (description != null) 'description': description,
    };
  }
}

/// Transaction initiate request
class TransactionInitiateRequest {
  final String productId;
  final String customerNo; // phone number or meter number
  final double? sellingPrice; // optional, defaults to product's selling price
  final String? notes; // additional notes for products like token listrik

  TransactionInitiateRequest({
    required this.productId,
    required this.customerNo,
    this.sellingPrice,
    this.notes,
  });

  Map<String, dynamic> toJson() {
    return {
      'product_id': productId,
      'customer_no': customerNo,
      if (sellingPrice != null) 'selling_price': sellingPrice,
      if (notes != null) 'notes': notes,
    };
  }
}

/// Product sync request
class SyncProductsRequest {
  final String provider; // 'digiflazz'
  final String type; // 'prepaid' or 'postpaid'
  final List<String>? productCodes; // optional specific codes

  SyncProductsRequest({
    required this.provider,
    required this.type,
    this.productCodes,
  });

  Map<String, dynamic> toJson() {
    return {
      'provider': provider,
      'type': type,
      if (productCodes != null) 'product_codes': productCodes,
    };
  }
}

/// Provider response
class Provider {
  final String id;
  final String name;
  final String code;
  final bool isActive;
  final Map<String, dynamic>? config;

  Provider({
    required this.id,
    required this.name,
    required this.code,
    required this.isActive,
    this.config,
  });

  factory Provider.fromJson(Map<String, dynamic> json) {
    return Provider(
      id: json['id'] as String,
      name: json['name'] as String,
      code: json['code'] as String,
      isActive: json['is_active'] as bool? ?? true,
      config: json['config'] as Map<String, dynamic>?,
    );
  }
}

/// Error code catalog entry
class ErrorCode {
  final String code;
  final String message;
  final String? description;
  final String category;

  ErrorCode({
    required this.code,
    required this.message,
    this.description,
    required this.category,
  });

  factory ErrorCode.fromJson(Map<String, dynamic> json) {
    return ErrorCode(
      code: json['code'] as String,
      message: json['message'] as String,
      description: json['description'] as String?,
      category: json['category'] as String? ?? 'general',
    );
  }
}

/// Sync status response from product service
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
