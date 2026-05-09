import 'package:dio/dio.dart';
import '../models/auth_response.dart';
import '../models/error_response.dart';
import '../core/exceptions.dart';

/// Auth Service API client
/// Handles all authentication-related endpoints
/// Port: 8081, Base Path: /api/v1/auth
class AuthService {
  final Dio _dio;
  final String _baseUrl = 'https://fedora.sinauplatform.id/api/v1/auth';

  AuthService(this._dio);

  /// Register a new user
  /// POST /register
  Future<AuthResponse> register(RegisterRequest request) async {
    try {
      final response = await _dio.post(
        '$_baseUrl/register',
        data: request.toJson(),
      );
      return AuthResponse.fromJson(response.data);
    } on ApiException catch (e) {
      throw e;
    } catch (e) {
      throw ApiException(message: 'Registration failed: $e');
    }
  }

  /// Login with email/phone and password
  /// POST /login
  Future<AuthResponse> login(String username, String password) async {
    try {
      final response = await _dio.post(
        '$_baseUrl/login',
        data: {
          'username': username,
          'password': password,
        },
      );
      return AuthResponse.fromJson(response.data);
    } on ApiException catch (e) {
      throw e;
    } catch (e) {
      throw ApiException(message: 'Login failed: $e');
    }
  }

  /// Verify OTP for login or registration
  /// POST /verify-otp
  Future<AuthResponse> verifyOtp(String phoneNumber, String otp) async {
    try {
      final response = await _dio.post(
        '$_baseUrl/verify-otp',
        data: {
          'phone_number': phoneNumber,
          'otp': otp,
        },
      );
      return AuthResponse.fromJson(response.data);
    } on ApiException catch (e) {
      throw e;
    } catch (e) {
      throw ApiException(message: 'OTP verification failed: $e');
    }
  }

  /// Refresh access token using refresh token
  /// POST /refresh
  Future<AuthResponse> refreshToken(String refreshToken) async {
    try {
      final response = await _dio.post(
        '$_baseUrl/refresh',
        data: {
          'refresh_token': refreshToken,
        },
      );
      return AuthResponse.fromJson(response.data);
    } on ApiException catch (e) {
      throw e;
    } catch (e) {
      throw ApiException(message: 'Token refresh failed: $e');
    }
  }

  /// Logout and invalidate session
  /// POST /logout (requires Bearer token)
  Future<void> logout() async {
    try {
      await _dio.post('$_baseUrl/logout');
    } on ApiException catch (e) {
      if (e.statusCode != 401) rethrow;
      // Ignore 401 on logout (already invalid)
    } catch (e) {
      // Non-critical error, proceed with local logout
    }
  }

  /// Change password
  /// POST /change-password (requires Bearer token)
  Future<void> changePassword(String oldPassword, String newPassword) async {
    try {
      await _dio.post(
        '$_baseUrl/change-password',
        data: {
          'old_password': oldPassword,
          'new_password': newPassword,
        },
      );
    } on ApiException catch (e) {
      throw e;
    } catch (e) {
      throw ApiException(message: 'Password change failed: $e');
    }
  }

  /// Change transaction PIN
  /// POST /change-pin (requires Bearer token)
  Future<void> changePin(String oldPin, String newPin) async {
    try {
      await _dio.post(
        '$_baseUrl/change-pin',
        data: {
          'old_pin': oldPin,
          'new_pin': newPin,
        },
      );
    } on ApiException catch (e) {
      throw e;
    } catch (e) {
      throw ApiException(message: 'PIN change failed: $e');
    }
  }

  /// Request OTP for login/register
  /// POST /request-otp (not in doc but commonly needed)
  Future<void> requestOtp(String phoneNumber) async {
    try {
      await _dio.post(
        '$_baseUrl/request-otp',
        data: {
          'phone_number': phoneNumber,
        },
      );
    } catch (e) {
      throw ApiException(message: 'Failed to request OTP: $e');
    }
  }
}
