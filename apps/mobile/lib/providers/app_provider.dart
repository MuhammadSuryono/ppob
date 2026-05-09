import 'package:flutter_riverpod/flutter_riverpod.dart';

final appProvider = StateNotifierProvider<AppNotifier, AppState>((ref) {
  return AppNotifier();
});

class AppState {
  final bool isInitialized;
  final bool isLoading;

  const AppState({
    required this.isInitialized,
    required this.isLoading,
  });

  factory AppState.initial() => AppState(
        isInitialized: false,
        isLoading: true,
      );

  AppState copyWith({
    bool? isInitialized,
    bool? isLoading,
  }) {
    return AppState(
      isInitialized: isInitialized ?? this.isInitialized,
      isLoading: isLoading ?? this.isLoading,
    );
  }
}

class AppNotifier extends StateNotifier<AppState> {
  AppNotifier() : super(AppState.initial()) {
    initializeApp();
  }

  Future<void> initializeApp() async {
    // Simulate initialization delay
    await Future.delayed(const Duration(seconds: 2));
    state = state.copyWith(isInitialized: true, isLoading: false);
  }
}