import 'package:hive/hive.dart';

part 'user.g.dart';

@HiveType(typeId: 0)
class User {
  @HiveField(0)
  final String id;

  @HiveField(1)
  final String username;

  @HiveField(2)
  final String phoneNumber;

  @HiveField(3)
  final String role; // 'mitra' or 'staff'

  @HiveField(4)
  final String? pinHash;

  @HiveField(5)
  final int trustScore; // 0-100

  @HiveField(6)
  final DateTime createdAt;

  @HiveField(7)
  final DateTime updatedAt;

  User({
    required this.id,
    required this.username,
    required this.phoneNumber,
    required this.role,
    this.pinHash,
    this.trustScore = 50,
    required this.createdAt,
    required this.updatedAt,
  });

  User copyWith({
    String? id,
    String? username,
    String? phoneNumber,
    String? role,
    String? pinHash,
    int? trustScore,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) {
    return User(
      id: id ?? this.id,
      username: username ?? this.username,
      phoneNumber: phoneNumber ?? this.phoneNumber,
      role: role ?? this.role,
      pinHash: pinHash ?? this.pinHash,
      trustScore: trustScore ?? this.trustScore,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
    );
  }

  factory User.fromJson(Map<String, dynamic> json) {
    return User(
      id: json['id'] as String,
      username: json['username'] as String,
      phoneNumber: json['phoneNumber'] as String,
      role: json['role'] as String,
      pinHash: json['pinHash'] as String?,
      trustScore: json['trustScore'] as int? ?? 50,
      createdAt: DateTime.parse(json['createdAt'] as String),
      updatedAt: DateTime.parse(json['updatedAt'] as String),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'username': username,
      'phoneNumber': phoneNumber,
      'role': role,
      'pinHash': pinHash,
      'trustScore': trustScore,
      'createdAt': createdAt.toIso8601String(),
      'updatedAt': updatedAt.toIso8601String(),
    };
  }
}
