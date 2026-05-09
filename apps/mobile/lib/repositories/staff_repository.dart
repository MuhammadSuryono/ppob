import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../models/staff.dart';
import '../models/wallet.dart';
import '../models/auth_response.dart';
import '../services/user_service.dart';
import '../core/api_client.dart';
import '../utils/constants.dart';

final staffRepositoryProvider = Provider<StaffRepository>((ref) {
  final dio = ref.read(dioProvider);
  return StaffRepositoryImpl(UserService(dio));
});

abstract class StaffRepository {
  Future<List<Staff>> getStaffList(String mitraId);
  Future<Staff> getStaffDetail(String staffId);
  Future<Staff> createStaff(Staff staff);
  Future<Staff> updateStaff(Staff staff);
  Future<void> deleteStaff(String staffId);
  Future<Wallet> getStaffWallet(String staffId);
  Future<List<UserResponse>> searchUsers(String query, {String? role});
}

class StaffRepositoryImpl implements StaffRepository {
  final UserService _userService;
  final List<Staff> _localStaff = [];

  StaffRepositoryImpl(this._userService) {
    _seedLocalStaff();
  }

  void _seedLocalStaff() {
    _localStaff.addAll([
      Staff(
        id: 'staff_001',
        mitraId: 'mitra_1',
        name: 'Ahmad Wijaya',
        phoneNumber: '+628123456001',
        marginScheme: 'FixedAllowance',
        fixedAllowanceAmount: 500000.0,
        dailyLimitAmount: 1000000.0,
        dailyLimitCount: 50,
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      ),
      Staff(
        id: 'staff_002',
        mitraId: 'mitra_1',
        name: 'Siti Nurhaliza',
        phoneNumber: '+628123456002',
        marginScheme: 'MarginShare',
        marginSharePercentage: 0.7,
        dailyLimitAmount: 2000000.0,
        dailyLimitCount: 100,
        createdAt: DateTime.now(),
        updatedAt: DateTime.now(),
      ),
    ]);
  }

  @override
  Future<List<Staff>> getStaffList(String mitraId) async {
    try {
      // In real implementation, you'd fetch users with role='staff' assigned to mitra
      // For now, return local filtered
      await Future.delayed(const Duration(milliseconds: 300));
      return _localStaff.where((s) => s.mitraId == mitraId).toList();
    } catch (e) {
      throw Exception('Failed to get staff list: $e');
    }
  }

  @override
  Future<Staff> getStaffDetail(String staffId) async {
    try {
      await Future.delayed(const Duration(milliseconds: 200));
      final staff = _localStaff.firstWhere((s) => s.id == staffId);
      return staff;
    } catch (e) {
      throw Exception('Staff not found: $e');
    }
  }

  @override
  Future<Staff> createStaff(Staff staff) async {
    try {
      // In real implementation, create user with role='staff' and associate with mitra
      await Future.delayed(const Duration(milliseconds: 500));
      _localStaff.add(staff);
      return staff;
    } catch (e) {
      throw Exception('Failed to create staff: $e');
    }
  }

  @override
  Future<Staff> updateStaff(Staff staff) async {
    try {
      await Future.delayed(const Duration(milliseconds: 500));
      final index = _localStaff.indexWhere((s) => s.id == staff.id);
      if (index >= 0) {
        _localStaff[index] = staff;
        return staff;
      }
      throw Exception('Staff not found');
    } catch (e) {
      throw Exception('Failed to update staff: $e');
    }
  }

  @override
  Future<void> deleteStaff(String staffId) async {
    try {
      await Future.delayed(const Duration(milliseconds: 300));
      _localStaff.removeWhere((s) => s.id == staffId);
    } catch (e) {
      throw Exception('Failed to delete staff: $e');
    }
  }

  @override
  Future<Wallet> getStaffWallet(String staffId) async {
    try {
      // Fetch wallet from wallet service
      // This would be handled by walletRepository in practice
      await Future.delayed(const Duration(milliseconds: 200));
      return Wallet(
        id: 'wallet_$staffId',
        userId: staffId,
        role: 'staff',
        ownerName: 'Staff',
        availableBalance: 100000.0,
        heldBalance: 0.0,
        date: DateTime.now(),
        updatedAt: DateTime.now(),
      );
    } catch (e) {
      throw Exception('Failed to get staff wallet: $e');
    }
  }

  @override
  Future<List<UserResponse>> searchUsers(String query, {String? role}) async {
    try {
      // Search users from backend
      final users = await _userService.listUsers(search: query, limit: 20);
      if (role != null) {
        return users.data.where((u) => u.role == role).toList();
      }
      return users.data;
    } catch (e) {
      throw Exception('Failed to search users: $e');
    }
  }
}
