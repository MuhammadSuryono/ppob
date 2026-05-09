import 'package:dio/dio.dart';
import '../models/wallet.dart';
import '../models/api_models.dart';
import '../core/exceptions.dart';

/// Wallet Service API client
/// Handles financial ledger operations: balance, holds, transfers
/// Port: 8083, Base Path: /api/v1/wallets
class WalletService {
  final Dio _dio;
  final String _baseUrl = 'https://192.168.100.23:8083/api/v1/wallets';

  WalletService(this._dio);

  /// Get wallet balance
  /// GET /:id/balance (requires Bearer token)
  Future<BalanceResponse> getBalance(String walletId) async {
    try {
      final response = await _dio.get('$_baseUrl/$walletId/balance');
      return BalanceResponse.fromJson(response.data as Map<String, dynamic>);
    } on ApiException catch (e) {
      throw e;
    } catch (e) {
      throw ApiException(message: 'Failed to get balance: $e');
    }
  }

  /// Get balance reconstructed from events (for reconciliation)
  /// GET /:id/balance-events (requires Bearer token)
  Future<BalanceResponse> getBalanceFromEvents(String walletId) async {
    try {
      final response = await _dio.get('$_baseUrl/$walletId/balance-events');
      return BalanceResponse.fromJson(response.data as Map<String, dynamic>);
    } on ApiException catch (e) {
      throw e;
    } catch (e) {
      throw ApiException(message: 'Failed to get balance from events: $e');
    }
  }

  /// Place a hold on funds for a transaction
  /// POST /:id/hold (requires Bearer token)
  Future<HoldResponse> holdFunds(String walletId, HoldFundsRequest request) async {
    try {
      final response = await _dio.post(
        '$_baseUrl/$walletId/hold',
        data: request.toJson(),
      );
      return HoldResponse.fromJson(response.data as Map<String, dynamic>);
    } on ApiException catch (e) {
      throw e;
    } catch (e) {
      throw ApiException(message: 'Failed to place hold on funds: $e');
    }
  }

  /// Release a previously held amount
  /// POST /:id/release-hold (requires Bearer token)
  Future<void> releaseHold(String walletId, ReleaseHoldRequest request) async {
    try {
      await _dio.post(
        '$_baseUrl/$walletId/release-hold',
        data: request.toJson(),
      );
    } on ApiException catch (e) {
      throw e;
    } catch (e) {
      throw ApiException(message: 'Failed to release hold: $e');
    }
  }

  /// Direct debit from wallet
  /// POST /:id/debit (requires Bearer token)
  Future<TransactionResponse> debitWallet(
    String walletId,
    double amount,
    String transactionId, {
    String? description,
  }) async {
    try {
      final response = await _dio.post(
        '$_baseUrl/$walletId/debit',
        data: {
          'amount': amount,
          'transaction_id': transactionId,
          if (description != null) 'description': description,
        },
      );
      return TransactionResponse.fromJson(response.data as Map<String, dynamic>);
    } on ApiException catch (e) {
      throw e;
    } catch (e) {
      throw ApiException(message: 'Failed to debit wallet: $e');
    }
  }

  /// Direct credit to wallet
  /// POST /:id/credit (requires Bearer token)
  Future<TransactionResponse> creditWallet(
    String walletId,
    double amount,
    String transactionId, {
    String? description,
  }) async {
    try {
      final response = await _dio.post(
        '$_baseUrl/$walletId/credit',
        data: {
          'amount': amount,
          'transaction_id': transactionId,
          if (description != null) 'description': description,
        },
      );
      return TransactionResponse.fromJson(response.data as Map<String, dynamic>);
    } on ApiException catch (e) {
      throw e;
    } catch (e) {
      throw ApiException(message: 'Failed to credit wallet: $e');
    }
  }

  /// Internal transfer between users
  /// POST /transfer (requires Bearer token)
  Future<TransferResponse> transfer(TransferRequest request) async {
    try {
      final response = await _dio.post(
        '$_baseUrl/transfer',
        data: request.toJson(),
      );
      return TransferResponse.fromJson(response.data as Map<String, dynamic>);
    } on ApiException catch (e) {
      throw e;
    } catch (e) {
      throw ApiException(message: 'Failed to transfer: $e');
    }
  }

  /// Staff top-up by Mitra
  /// POST /staff/topup (requires Bearer token)
  Future<TransactionResponse> topUpStaff(String staffId, double amount) async {
    try {
      final response = await _dio.post(
        '$_baseUrl/staff/$staffId/topup',
        data: {'amount': amount},
      );
      return TransactionResponse.fromJson(response.data as Map<String, dynamic>);
    } on ApiException catch (e) {
      throw e;
    } catch (e) {
      throw ApiException(message: 'Failed to top-up staff: $e');
    }
  }

  /// Get wallet transaction history
  /// GET /:id/events (requires Bearer token)
  Future<List<WalletEvent>> getWalletEvents(
    String walletId, {
    int limit = 20,
    int offset = 0,
    String? type,
    DateTime? from,
    DateTime? to,
  }) async {
    try {
      final queryParams = <String, dynamic>{
        'limit': limit,
        'offset': offset,
      };
      if (type != null) queryParams['type'] = type;
      if (from != null) queryParams['from'] = from.toIso8601String();
      if (to != null) queryParams['to'] = to.toIso8601String();

      final response = await _dio.get(
        '$_baseUrl/$walletId/events',
        queryParameters: queryParams,
      );
      final List<dynamic> data = response.data as List<dynamic>;
      return data.map((json) => WalletEvent.fromJson(json as Map<String, dynamic>)).toList();
    } on ApiException catch (e) {
      throw e;
    } catch (e) {
      throw ApiException(message: 'Failed to get wallet events: $e');
    }
  }

  /// Check balance vs events drift (reconciliation)
  /// GET /:id/reconcile (requires Bearer token)
  Future<ReconciliationResponse> reconcile(String walletId) async {
    try {
      final response = await _dio.get('$_baseUrl/$walletId/reconcile');
      return ReconciliationResponse.fromJson(response.data as Map<String, dynamic>);
    } on ApiException catch (e) {
      throw e;
    } catch (e) {
      throw ApiException(message: 'Failed to reconcile wallet: $e');
    }
  }
}

/// Hold response
class HoldResponse {
  final String holdId;
  final double amount;
  final String? description;
  final DateTime heldAt;

  HoldResponse({
    required this.holdId,
    required this.amount,
    this.description,
    required this.heldAt,
  });

  factory HoldResponse.fromJson(Map<String, dynamic> json) {
    return HoldResponse(
      holdId: json['hold_id'] as String,
      amount: (json['amount'] as num).toDouble(),
      description: json['description'] as String?,
      heldAt: DateTime.parse(json['held_at'] as String),
    );
  }
}

/// Transaction response for wallet operations
class TransactionResponse {
  final String transactionId;
  final String walletId;
  final double amount;
  final double balanceAfter;
  final String type; // 'debit' or 'credit'
  final DateTime timestamp;

  TransactionResponse({
    required this.transactionId,
    required this.walletId,
    required this.amount,
    required this.balanceAfter,
    required this.type,
    required this.timestamp,
  });

  factory TransactionResponse.fromJson(Map<String, dynamic> json) {
    return TransactionResponse(
      transactionId: json['transaction_id'] as String,
      walletId: json['wallet_id'] as String,
      amount: (json['amount'] as num).toDouble(),
      balanceAfter: (json['balance_after'] as num).toDouble(),
      type: json['type'] as String,
      timestamp: DateTime.parse(json['timestamp'] as String),
    );
  }
}

/// Transfer response
class TransferResponse {
  final String transferId;
  final String fromWalletId;
  final String toWalletId;
  final double amount;
  final DateTime timestamp;

  TransferResponse({
    required this.transferId,
    required this.fromWalletId,
    required this.toWalletId,
    required this.amount,
    required this.timestamp,
  });

  factory TransferResponse.fromJson(Map<String, dynamic> json) {
    return TransferResponse(
      transferId: json['transfer_id'] as String,
      fromWalletId: json['from_wallet_id'] as String,
      toWalletId: json['to_wallet_id'] as String,
      amount: (json['amount'] as num).toDouble(),
      timestamp: DateTime.parse(json['timestamp'] as String),
    );
  }
}

/// Wallet event (transaction history)
class WalletEvent {
  final String id;
  final String walletId;
  final String type; // 'debit', 'credit', 'hold', 'release'
  final double amount;
  final double balanceAfter;
  final String? referenceId; // transaction_id or hold_id
  final String? description;
  final DateTime createdAt;

  WalletEvent({
    required this.id,
    required this.walletId,
    required this.type,
    required this.amount,
    required this.balanceAfter,
    this.referenceId,
    this.description,
    required this.createdAt,
  });

  factory WalletEvent.fromJson(Map<String, dynamic> json) {
    return WalletEvent(
      id: json['id'] as String,
      walletId: json['wallet_id'] as String,
      type: json['type'] as String,
      amount: (json['amount'] as num).toDouble(),
      balanceAfter: (json['balance_after'] as num).toDouble(),
      referenceId: json['reference_id'] as String?,
      description: json['description'] as String?,
      createdAt: DateTime.parse(json['created_at'] as String),
    );
  }
}

/// Reconciliation response
class ReconciliationResponse {
  final bool isConsistent;
  final double? drift;
  final List< ReconciliationIssue>? issues;

  ReconciliationResponse({
    required this.isConsistent,
    this.drift,
    this.issues,
  });

  factory ReconciliationResponse.fromJson(Map<String, dynamic> json) {
    return ReconciliationResponse(
      isConsistent: json['is_consistent'] as bool? ?? true,
      drift: (json['drift'] as num?)?.toDouble(),
      issues: json['issues'] != null
          ? (json['issues'] as List<dynamic>)
              .map((e) => ReconciliationIssue.fromJson(e as Map<String, dynamic>))
              .toList()
          : null,
    );
  }
}

class ReconciliationIssue {
  final String type;
  final String description;
  final double? amount;

  ReconciliationIssue({
    required this.type,
    required this.description,
    this.amount,
  });

  factory ReconciliationIssue.fromJson(Map<String, dynamic> json) {
    return ReconciliationIssue(
      type: json['type'] as String,
      description: json['description'] as String,
      amount: (json['amount'] as num?)?.toDouble(),
    );
  }
}
