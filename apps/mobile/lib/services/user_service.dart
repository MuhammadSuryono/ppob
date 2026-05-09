import 'package:dio/dio.dart';
import '../models/auth_response.dart';
import '../models/error_response.dart';
import '../models/api_models.dart';
import '../core/exceptions.dart';

/// User Service API client
/// Handles user profiles, roles, and administrative lists
/// Port: 8082, Base Path: /api/v1
class UserService {
  final Dio _dio;
  final String _baseUrl = 'https://192.168.100.23:8082/api/v1';

  UserService(this._dio);

  /// Get user profile by ID
  /// GET /users/:id (requires Bearer token, Owner/Admin only)
  Future<UserResponse> getUserProfile(String userId) async {
    try {
      final response = await _dio.get('$_baseUrl/users/$userId');
      return UserResponse.fromJson(response.data);
    } on ApiException catch (e) {
      throw e;
    } catch (e) {
      throw ApiException(message: 'Failed to get user profile: $e');
    }
  }

  /// Update user profile
  /// PUT /users/:id (requires Bearer token, Owner/Admin only)
  Future<UserResponse> updateUserProfile(
    String userId,
    {
      String? name,
      String? phone,
      String? address,
      DateTime? dateOfBirth,
    }
  ) async {
    try {
      final data = <String, dynamic>{};
      if (name != null) data['name'] = name;
      if (phone != null) data['phone'] = phone;
      if (address != null) data['address'] = address;
      if (dateOfBirth != null) data['date_of_birth'] = dateOfBirth.toIso8601String();

      final response = await _dio.put(
        '$_baseUrl/users/$userId',
        data: data,
      );
      return UserResponse.fromJson(response.data);
    } on ApiException catch (e) {
      throw e;
    } catch (e) {
      throw ApiException(message: 'Failed to update user profile: $e');
    }
  }

  /// List all users (paginated)
  /// GET /users (requires Bearer token, Admin/Staff only)
  Future<PaginatedResponse<UserResponse>> listUsers({
    int page = 1,
    int limit = 20,
    String? role,
    String? search,
  }) async {
    try {
      final queryParams = <String, dynamic>{
        'page': page,
        'limit': limit,
      };
      if (role != null) queryParams['role'] = role;
      if (search != null) queryParams['search'] = search;

      final response = await _dio.get(
        '$_baseUrl/users',
        queryParameters: queryParams,
      );

      final data = response.data as Map<String, dynamic>;
      return PaginatedResponse<UserResponse>.fromJson(
        data,
        (item) => UserResponse.fromJson(item as Map<String, dynamic>),
      );
    } on ApiException catch (e) {
      throw e;
    } catch (e) {
      throw ApiException(message: 'Failed to list users: $e');
    }
  }

  /// Get roles assigned to a user
  /// GET /users/:id/roles (requires Bearer token, Owner/Admin only)
  Future<List<RoleResponse>> getUserRoles(String userId) async {
    try {
      final response = await _dio.get('$_baseUrl/users/$userId/roles');
      final List<dynamic> data = response.data as List<dynamic>;
      return data.map((item) => RoleResponse.fromJson(item as Map<String, dynamic>)).toList();
    } on ApiException catch (e) {
      throw e;
    } catch (e) {
      throw ApiException(message: 'Failed to get user roles: $e');
    }
  }

  /// Assign a role to a user
  /// POST /users/:id/roles (requires Bearer token, Admin only)
  Future<RoleResponse> assignRole(String userId, String roleId) async {
    try {
      final response = await _dio.post(
        '$_baseUrl/users/$userId/roles',
        data: {
          'role_id': roleId,
        },
      );
      return RoleResponse.fromJson(response.data);
    } on ApiException catch (e) {
      throw e;
    } catch (e) {
      throw ApiException(message: 'Failed to assign role: $e');
    }
  }

  /// List all available system roles
  /// GET /roles (requires Bearer token, Admin only)
  Future<List<RoleResponse>> listRoles() async {
    try {
      final response = await _dio.get('$_baseUrl/roles');
      final List<dynamic> data = response.data as List<dynamic>;
      return data.map((item) => RoleResponse.fromJson(item as Map<String, dynamic>)).toList();
    } on ApiException catch (e) {
      throw e;
    } catch (e) {
      throw ApiException(message: 'Failed to list roles: $e');
    }
  }

  /// Create a new role definition
  /// POST /roles (requires Bearer token, Admin only)
  Future<RoleResponse> createRole(String name, String description, List<String> permissions) async {
    try {
      final response = await _dio.post(
        '$_baseUrl/roles',
        data: {
          'name': name,
          'description': description,
          'permissions': permissions,
        },
      );
      return RoleResponse.fromJson(response.data);
    } on ApiException catch (e) {
      throw e;
    } catch (e) {
      throw ApiException(message: 'Failed to create role: $e');
    }
  }
}

/// Role response model
class RoleResponse {
  final String id;
  final String name;
  final String description;
  final List<String> permissions;
  final bool isSystemRole;

  RoleResponse({
    required this.id,
    required this.name,
    required this.description,
    required this.permissions,
    required this.isSystemRole,
  });

  factory RoleResponse.fromJson(Map<String, dynamic> json) {
    return RoleResponse(
      id: json['id'] as String,
      name: json['name'] as String,
      description: json['description'] as String? ?? '',
      permissions: (json['permissions'] as List<dynamic>?)?.cast<String>() ?? [],
      isSystemRole: json['is_system_role'] as bool? ?? false,
    );
  }
}
