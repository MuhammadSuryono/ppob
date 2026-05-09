import 'dart:convert';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../models/user.dart';
import '../models/auth_response.dart';
import '../services/auth_service.dart';
import '../core/api_client.dart';
import '../core/exceptions.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

final authRepositoryProvider = Provider<AuthRepository>((ref) {
  final dio = ref.read(dioProvider);
  final authService = AuthService(dio);
  final storage = ref.read(secureStorageProvider);
  return AuthRepositoryImpl(authService, storage);
});

abstract class AuthRepository {
  Future<User> login(String username, String password);
  Future<User> register({
    required String email,
    required String phone,
    required String name,
    required String password,
    required String pin,
    String? referralCode,
  });
  Future<void> logout();
  Future<User?> getCurrentUser();
  Future<String?> getAuthToken();
  Future<String?> getRefreshToken();
  Future<AuthResponse> verifyOtp(String phoneNumber, String otp);
  Future<void> changePassword(String oldPassword, String newPassword);
  Future<void> changePin(String oldPin, String newPin);
  Future<AuthResponse> refreshToken();
}

class AuthRepositoryImpl implements AuthRepository {
  final AuthService _authService;
  final FlutterSecureStorage _storage;
  static const _userKey = 'user_data';
  static const _tokenKey = 'auth_token';
  static const _refreshTokenKey = 'refresh_token';

  AuthRepositoryImpl(this._authService, this._storage);

  @override
  Future<String?> getRefreshToken() async {
    return await _storage.read(key: _refreshTokenKey);
  }

  @override
  Future<AuthResponse> refreshToken() async {
    try {
      final refreshToken = await _storage.read(key: _refreshTokenKey);
      if (refreshToken == null) {
        throw Exception('No refresh token available');
      }
      final response = await _authService.refreshToken(refreshToken);
      await _saveAuthData(response);
      return response;
    } on ApiException catch (e) {
      throw Exception(e.message);
    } catch (e) {
      throw Exception('Token refresh failed: $e');
    }
  }

  @override
  Future<User> login(String username, String password) async {
    try {
      final response = await _authService.login(username, password);
      await _saveAuthData(response);
      return _mapUserResponse(response.user);
    } on ApiException catch (e) {
      throw Exception(e.message);
    } catch (e) {
      throw Exception('Login failed: $e');
    }
  }

  @override
  Future<User> register({
    required String email,
    required String phone,
    required String name,
    required String password,
    required String pin,
    String? referralCode,
  }) async {
    try {
      final request = RegisterRequest(
        email: email,
        phone: phone,
        name: name,
        password: password,
        pin: pin,
        referralCode: referralCode,
      );
      final response = await _authService.register(request);
      await _saveAuthData(response);
      return _mapUserResponse(response.user);
    } on ApiException catch (e) {
      throw Exception(e.message);
    } catch (e) {
      throw Exception('Registration failed: $e');
    }
  }

  @override
  Future<void> logout() async {
    try {
      await _authService.logout();
    } catch (e) {
      // Proceed with local logout even if server fails
    } finally {
      await _storage.delete(key: _tokenKey);
      await _storage.delete(key: _refreshTokenKey);
      await _storage.delete(key: _userKey);
    }
  }

  @override
  Future<User?> getCurrentUser() async {
    final userJson = await _storage.read(key: _userKey);
    if (userJson != null) {
      try {
        final Map<String, dynamic> json = jsonDecode(userJson);
        return User.fromJson(json);
      } catch (e) {
        return null;
      }
    }
    return null;
  }

  @override
  Future<String?> getAuthToken() async {
    return await _storage.read(key: _tokenKey);
  }

  @override
  Future<AuthResponse> verifyOtp(String phoneNumber, String otp) async {
    try {
      final response = await _authService.verifyOtp(phoneNumber, otp);
      await _saveAuthData(response);
      return response;
    } on ApiException catch (e) {
      throw Exception(e.message);
    } catch (e) {
      throw Exception('OTP verification failed: $e');
    }
  }

  @override
  Future<void> changePassword(String oldPassword, String newPassword) async {
    try {
      await _authService.changePassword(oldPassword, newPassword);
    } on ApiException catch (e) {
      throw Exception(e.message);
    } catch (e) {
      throw Exception('Password change failed: $e');
    }
  }

  @override
  Future<void> changePin(String oldPin, String newPin) async {
    try {
      await _authService.changePin(oldPin, newPin);
    } on ApiException catch (e) {
      throw Exception(e.message);
    } catch (e) {
      throw Exception('PIN change failed: $e');
    }
  }

  Future<void> _saveAuthData(AuthResponse response) async {
    final user = _mapUserResponse(response.user);
    await _storage.write(key: _tokenKey, value: response.accessToken);
    await _storage.write(key: _refreshTokenKey, value: response.refreshToken);
    await _storage.write(key: _userKey, value: jsonEncode(user.toJson()));
  }

  User _mapUserResponse(UserResponse userResponse) {
    return User(
      id: userResponse.id,
      username: userResponse.name,
      phoneNumber: userResponse.phone,
      role: userResponse.role,
      trustScore: userResponse.trustScore,
      createdAt: userResponse.createdAt,
      updatedAt: userResponse.updatedAt,
    );
  }
}
