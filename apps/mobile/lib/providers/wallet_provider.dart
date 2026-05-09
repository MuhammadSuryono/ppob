import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../models/wallet.dart';
import '../repositories/wallet_repository.dart';

final walletProvider = StateNotifierProvider<WalletNotifier, WalletState>((ref) {
  return WalletNotifier(ref.watch(walletRepositoryProvider));
});

class WalletState {
  final Wallet? wallet;
  final WalletStatus status;
  final String? errorMessage;

  const WalletState({
    this.wallet,
    required this.status,
    this.errorMessage,
  });

  factory WalletState.initial() => WalletState(
        status: WalletStatus.loading,
      );

  WalletState copyWith({
    Wallet? wallet,
    WalletStatus? status,
    String? errorMessage,
  }) {
    return WalletState(
      wallet: wallet ?? this.wallet,
      status: status ?? this.status,
      errorMessage: errorMessage ?? this.errorMessage,
    );
  }
}

enum WalletStatus { loading, loaded, error }

class WalletNotifier extends StateNotifier<WalletState> {
  final WalletRepository _repo;

  WalletNotifier(this._repo) : super(WalletState.initial());

  Future<void> loadWallet(String userId) async {
    state = const WalletState(status: WalletStatus.loading);
    try {
      final wallet = await _repo.getActiveWallet(userId);
      state = WalletState(
        wallet: wallet,
        status: WalletStatus.loaded,
      );
    } catch (e) {
      state = WalletState(
        status: WalletStatus.error,
        errorMessage: e.toString(),
      );
    }
  }

  Future<void> refresh() async {
    final userId = state.wallet?.userId;
    if (userId != null) {
      await loadWallet(userId);
    }
  }
}
