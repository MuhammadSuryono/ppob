import 'dart:math';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:hive/hive.dart';
import '../models/transaction.dart';
import '../models/pending_sync_item.dart';
import '../models/api_models.dart' hide Provider;
import '../services/transaction_service.dart';
import '../services/wallet_service.dart';
import '../services/product_service.dart';
import '../core/api_client.dart';
import '../core/exceptions.dart';
import '../utils/constants.dart';

final transactionRepositoryProvider = Provider<TransactionRepository>((ref) {
  final dio = ref.read(dioProvider);
  return TransactionRepositoryImpl(
    TransactionService(dio),
    WalletService(dio),
    ProductService(dio),
  );
});

final transactionBoxProvider = Provider<Box<Transaction>>((ref) {
  throw UnimplementedError('Hive box must be opened in main.dart');
});

abstract class TransactionRepository {
  Future<Transaction> initiateTransaction({
    required String productId,
    required String productName,
    required String customerNo,
    required double sellingPrice,
    required double platformPrice,
    String? staffId,
    String? mitraId,
  });

  Future<List<Transaction>> getTransactionHistory(String userId, {String? status, int limit = 20, int offset = 0});

  Future<Transaction> getTransactionDetail(String transactionId);

  Future<Transaction> updateTransactionStatus(
    String transactionId,
    String status, {
    String? rcCode,
    String? rcMessage,
    List<String>? notes,
  });

  Future<void> cancelTransaction(String transactionId, {String? reason});

  Future<void> addToPendingSync(PendingSyncItem item);
  Future<List<PendingSyncItem>> getPendingSyncItems();
  Future<void> removePendingSyncItem(String itemId);
}

class TransactionRepositoryImpl implements TransactionRepository {
  final TransactionService _transactionService;
  final WalletService _walletService;
  final ProductService _productService;
  final List<Transaction> _localTransactions = [];

  TransactionRepositoryImpl(this._transactionService, this._walletService, this._productService);

  @override
  Future<Transaction> initiateTransaction({
    required String productId,
    required String productName,
    required String customerNo,
    required double sellingPrice,
    required double platformPrice,
    String? staffId,
    String? mitraId,
  }) async {
    try {
      // Validate price
      final isValid = await _productService.validatePrice(productId, sellingPrice);
      if (!isValid) {
        throw Exception('Harga jual tidak memenuhi syarat minimum');
      }

      // Get user ID
      final userId = staffId ?? mitraId;
      if (userId == null) {
        throw Exception('User ID required');
      }

      // Check wallet balance (assuming walletId = userId)
      final balance = await _walletService.getBalance(userId);
      if (balance.availableBalance < sellingPrice) {
        throw InsufficientBalanceException(
          required: sellingPrice,
          current: balance.availableBalance,
        );
      }

      // Generate idempotency key
      final random = Random();
      final idempotencyKey = 'tx_${DateTime.now().millisecondsSinceEpoch}_${random.nextInt(10000)}';

      final request = TransactionInitiateRequest(
        productId: productId,
        customerNo: customerNo,
        sellingPrice: sellingPrice,
      );

      final transaction = await _transactionService.initiateTransaction(
        request,
        idempotencyKey: idempotencyKey,
      );

      // Cache locally
      _addToLocalCache(transaction);
      return transaction;
    } on ApiException catch (e) {
      throw Exception(e.message);
    } catch (e) {
      throw Exception('Failed to initiate transaction: $e');
    }
  }

  @override
  Future<List<Transaction>> getTransactionHistory(
    String userId, {
    String? status,
    int limit = 20,
    int offset = 0,
  }) async {
    try {
      final transactions = await _transactionService.getTransactionHistory(
        userId: userId,
        limit: limit,
        cursor: offset.toString(),
        status: status,
      );
      _updateLocalCache(transactions.data);
      return transactions.data;
    } catch (e) {
      // Return local transactions sorted by date
      List<Transaction> result = _localTransactions;
      if (status != null) {
        result = result.where((t) => t.status == status).toList();
      }
      return result.skip(offset).take(limit).toList();
    }
  }

  @override
  Future<Transaction> getTransactionDetail(String transactionId) async {
    try {
      return await _transactionService.getTransaction(transactionId);
    } catch (e) {
      final transaction = _localTransactions.firstWhere(
        (t) => t.id == transactionId,
        orElse: () => throw Exception('Transaction not found'),
      );
      return transaction;
    }
  }

  @override
  Future<Transaction> updateTransactionStatus(
    String transactionId,
    String status, {
    String? rcCode,
    String? rcMessage,
    List<String>? notes,
  }) async {
    try {
      final updated = await _transactionService.updateTransactionStatus(
        transactionId,
        status,
        notes: notes?.join(', '),
      );
      _updateLocalCache([updated]);
      return updated;
    } catch (e) {
      throw Exception('Failed to update transaction status: $e');
    }
  }

  @override
  Future<void> cancelTransaction(String transactionId, {String? reason}) async {
    try {
      await _transactionService.cancelTransaction(transactionId, reason: reason);
      // Update local cache
      final index = _localTransactions.indexWhere((t) => t.id == transactionId);
      if (index >= 0) {
        _localTransactions[index] = _localTransactions[index].copyWith(
          status: AppConstants.statusFailed,
          updatedAt: DateTime.now(),
        );
      }
    } catch (e) {
      throw Exception('Failed to cancel transaction: $e');
    }
  }

  @override
  Future<void> addToPendingSync(PendingSyncItem item) async {
    // Store in Hive offline queue
    // Implementation depends on Hive box setup
  }

  @override
  Future<List<PendingSyncItem>> getPendingSyncItems() async {
    return [];
  }

  @override
  Future<void> removePendingSyncItem(String itemId) async {
    // Remove from Hive box
  }

  void _addToLocalCache(Transaction transaction) {
    final index = _localTransactions.indexWhere((t) => t.id == transaction.id);
    if (index >= 0) {
      _localTransactions[index] = transaction;
    } else {
      _localTransactions.insert(0, transaction);
    }
  }

  void _updateLocalCache(List<Transaction> transactions) {
    for (final tx in transactions) {
      _addToLocalCache(tx);
    }
  }
}
