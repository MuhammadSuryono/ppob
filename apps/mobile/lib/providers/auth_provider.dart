import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import '../models/user.dart';
import '../models/auth_response.dart';
import '../repositories/auth_repository.dart';

final authProvider = StateNotifierProvider<AuthNotifier, AuthState>((ref) {
  return AuthNotifier(ref.watch(authRepositoryProvider), ref.read(secureStorageProvider));
});

final secureStorageProvider = Provider<FlutterSecureStorage>((ref) {
  return const FlutterSecureStorage();
});

class AuthState {
  final bool isAuthenticated;
  final String? userId;
  final String? username;
  final String? email;
  final String? phone;
  final String? role; // 'mitra' or 'staff'
  final String? token;
  final bool isLoading;
  final String? errorMessage;
  final AuthResponse? authResponse;

  const AuthState({
    required this.isAuthenticated,
    this.userId,
    this.username,
    this.email,
    this.phone,
    this.role,
    this.token,
    required this.isLoading,
    this.errorMessage,
    this.authResponse,
  });

  factory AuthState.initial() => AuthState(
        isAuthenticated: false,
        isLoading: false,
      );

  AuthState copyWith({
    bool? isAuthenticated,
    String? userId,
    String? username,
    String? email,
    String? phone,
    String? role,
    String? token,
    bool? isLoading,
    String? errorMessage,
    AuthResponse? authResponse,
  }) {
    return AuthState(
      isAuthenticated: isAuthenticated ?? this.isAuthenticated,
      userId: userId ?? this.userId,
      username: username ?? this.username,
      email: email ?? this.email,
      phone: phone ?? this.phone,
      role: role ?? this.role,
      token: token ?? this.token,
      isLoading: isLoading ?? this.isLoading,
      errorMessage: errorMessage ?? this.errorMessage,
      authResponse: authResponse ?? this.authResponse,
    );
  }
}

class AuthNotifier extends StateNotifier<AuthState> {
  final AuthRepository _authRepository;
  final FlutterSecureStorage _storage;
  static const _tokenKey = 'auth_token';

  AuthNotifier(this._authRepository, this._storage) : super(AuthState.initial()) {
    _loadAuthData();
  }

  Future<void> _loadAuthData() async {
    try {
      final user = await _authRepository.getCurrentUser();
      final token = await _authRepository.getAuthToken();

      if (user != null && token != null) {
        state = state.copyWith(
          isAuthenticated: true,
          userId: user.id,
          username: user.username,
          phone: user.phoneNumber,
          role: user.role,
          token: token,
        );
      }
    } catch (e) {
      // Clear invalid data
      await _authRepository.logout();
      state = AuthState.initial();
    }
  }

  Future<void> register({
    required String email,
    required String phone,
    required String name,
    required String password,
    required String pin,
    String? referralCode,
  }) async {
    state = state.copyWith(isLoading: true, errorMessage: null);
    try {
      final user = await _authRepository.register(
        email: email,
        phone: phone,
        name: name,
        password: password,
        pin: pin,
        referralCode: referralCode,
      );
      final token = await _authRepository.getAuthToken();

      state = state.copyWith(
        isAuthenticated: true,
        userId: user.id,
        username: user.username,
        phone: user.phoneNumber,
        role: user.role,
        token: token,
        isLoading: false,
      );
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        errorMessage: 'Registrasi gagal: ${e.toString()}',
      );
    }
  }

  Future<void> login(String usernameOrPhone, String password) async {
    state = state.copyWith(isLoading: true, errorMessage: null);
    try {
      final user = await _authRepository.login(usernameOrPhone, password);
      final token = await _authRepository.getAuthToken();

      state = state.copyWith(
        isAuthenticated: true,
        userId: user.id,
        username: user.username,
        phone: user.phoneNumber,
        role: user.role,
        token: token,
        isLoading: false,
      );
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        errorMessage: 'Login gagal: ${e.toString()}',
      );
    }
  }

  Future<void> verifyOtp(String phoneNumber, String otp) async {
    state = state.copyWith(isLoading: true, errorMessage: null);
    try {
      final response = await _authRepository.verifyOtp(phoneNumber, otp);
      final user = response.user;

      state = state.copyWith(
        isAuthenticated: true,
        userId: user.id,
        username: user.name,
        email: user.email,
        phone: user.phone,
        role: user.role,
        token: response.accessToken,
        isLoading: false,
        authResponse: response,
      );
    } catch (e) {
      state = state.copyWith(
        isLoading: false,
        errorMessage: 'Verifikasi OTP gagal: ${e.toString()}',
      );
    }
  }

  Future<void> logout() async {
    await _authRepository.logout();
    state = AuthState.initial();
  }

  Future<bool> verifyToken() async {
    final token = await _authRepository.getAuthToken();
    return token != null && token.isNotEmpty;
  }

   Future<void> refreshToken() async {
     try {
       final response = await _authRepository.refreshToken();
       state = state.copyWith(
         token: response.accessToken,
         authResponse: response,
       );
     } catch (e) {
       // Token refresh failed, logout
       await logout();
     }
   }

   void clearError() {
     state = state.copyWith(errorMessage: null);
   }
 }