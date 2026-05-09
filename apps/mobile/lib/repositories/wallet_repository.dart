import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../models/wallet.dart';
import '../models/transaction.dart';
import '../models/pending_sync_item.dart';
import '../models/api_models.dart' hide Provider;
import '../services/wallet_service.dart';
import '../services/user_service.dart';
import '../core/api_client.dart';
import '../core/exceptions.dart';
import '../utils/constants.dart';

final walletRepositoryProvider = Provider<WalletRepository>((ref) {
  final dio = ref.read(dioProvider);
  return WalletRepositoryImpl(WalletService(dio), UserService(dio));
});

abstract class WalletRepository {
  Future<Wallet> getActiveWallet(String userId);
  Future<void> topUpStaff(String staffId, double amount);
  Future<List<Transaction>> getWalletTransactions(String walletId, {int limit = 20});
  Future<void> holdAmount(String walletId, double amount, String referenceId);
  Future<void> releaseHold(String holdId);
  Future<Wallet?> getStaffWallet(String staffId);
  Future<BalanceResponse> getBalance(String walletId);
}

class WalletRepositoryImpl implements WalletRepository {
  final WalletService _walletService;
  final UserService _userService;
  final Map<String, Wallet> _cachedWallets = {};

  WalletRepositoryImpl(this._walletService, this._userService);

  @override
  Future<Wallet> getActiveWallet(String userId) async {
    try {
      // Fetch user profile for name and role
      final user = await _userService.getUserProfile(userId);
      // Fetch wallet balance
      final balance = await _walletService.getBalance(userId);

      final wallet = Wallet(
        id: balance.walletId,
        userId: userId,
        role: user.role,
        ownerName: user.name,
        availableBalance: balance.availableBalance,
        heldBalance: balance.heldBalance,
        dailySpentAmount: 0, // Could compute from events
        dailyTransactionCount: 0,
        date: DateTime.now(),
        updatedAt: balance.updatedAt,
      );

      _cachedWallets[wallet.id] = wallet;
      return wallet;
    } on ApiException catch (e) {
      throw Exception(e.message);
    } catch (e) {
      throw Exception('Failed to get wallet: $e');
    }
  }

  @override
  Future<void> topUpStaff(String staffId, double amount) async {
    try {
      await _walletService.topUpStaff(staffId, amount);
    } on ApiException catch (e) {
      if (e is InsufficientBalanceException) {
        throw Exception('Saldo mitra tidak mencukupi untuk top up');
      }
      throw Exception(e.message);
    } catch (e) {
      throw Exception('Failed to top up staff: $e');
    }
  }

  @override
  Future<List<Transaction>> getWalletTransactions(String walletId, {int limit = 20}) async {
    try {
      final events = await _walletService.getWalletEvents(walletId, limit: limit);
      // Convert WalletEvent to Transaction format for UI
      return events
          .where((e) => e.referenceId != null)
          .map((e) => Transaction(
                id: e.referenceId!,
                productId: 'WALLET_TRANSFER',
                productName: e.type == 'credit' ? 'Top Up' : 'Pembayaran',
                customerNo: walletId,
                sellingPrice: e.amount.abs(),
                platformPrice: e.amount.abs(),
                status: 'Success',
                createdAt: e.createdAt,
              ))
          .toList();
    } catch (e) {
      throw Exception('Failed to get wallet transactions: $e');
    }
  }

  @override
  Future<void> holdAmount(String walletId, double amount, String referenceId) async {
    try {
      await _walletService.holdFunds(
        walletId,
        HoldFundsRequest(
          transactionId: referenceId,
          amount: amount,
        ),
      );
    } on ApiException catch (e) {
      throw Exception(e.message);
    } catch (e) {
      throw Exception('Failed to hold amount: $e');
    }
  }

  @override
  Future<void> releaseHold(String holdId) async {
    try {
      await _walletService.releaseHold(
        holdId,
        ReleaseHoldRequest(holdId: holdId),
      );
    } on ApiException catch (e) {
      throw Exception(e.message);
    } catch (e) {
      throw Exception('Failed to release hold: $e');
    }
  }

  @override
  Future<Wallet?> getStaffWallet(String staffId) async {
    try {
      return await getActiveWallet(staffId);
    } catch (e) {
      // Staff might not have wallet yet
      return null;
    }
  }

  @override
  Future<BalanceResponse> getBalance(String walletId) async {
    try {
      return await _walletService.getBalance(walletId);
    } catch (e) {
      throw Exception('Failed to get balance: $e');
    }
   }
}
