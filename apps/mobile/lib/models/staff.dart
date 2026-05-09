import 'package:hive/hive.dart';

part 'staff.g.dart';

@HiveType(typeId: 3)
class Staff {
  @HiveField(0)
  final String id;

  @HiveField(1)
  final String mitraId;

  @HiveField(2)
  final String name;

  @HiveField(3)
  final String phoneNumber;

  @HiveField(4)
  final String? pinHash;

  @HiveField(5)
  final String marginScheme; // 'FixedAllowance' or 'MarginShare'

  @HiveField(6)
  final double? fixedAllowanceAmount; // Used if scheme is FixedAllowance

  @HiveField(7)
  final double? marginSharePercentage; // Used if scheme is MarginShare (e.g., 0.7 for 70%)

  @HiveField(8)
  final double dailyLimitAmount; // Max transaction amount per day

  @HiveField(9)
  final int dailyLimitCount; // Max transaction count per day

  @HiveField(10)
  final bool isActive;

  @HiveField(11)
  final DateTime createdAt;

  @HiveField(12)
  final DateTime updatedAt;

  Staff({
    required this.id,
    required this.mitraId,
    required this.name,
    required this.phoneNumber,
    this.pinHash,
    required this.marginScheme,
    this.fixedAllowanceAmount,
    this.marginSharePercentage,
    required this.dailyLimitAmount,
    required this.dailyLimitCount,
    this.isActive = true,
    required this.createdAt,
    required this.updatedAt,
  });

  Staff copyWith({
    String? id,
    String? mitraId,
    String? name,
    String? phoneNumber,
    String? pinHash,
    String? marginScheme,
    double? fixedAllowanceAmount,
    double? marginSharePercentage,
    double? dailyLimitAmount,
    int? dailyLimitCount,
    bool? isActive,
    DateTime? createdAt,
    DateTime? updatedAt,
  }) {
    return Staff(
      id: id ?? this.id,
      mitraId: mitraId ?? this.mitraId,
      name: name ?? this.name,
      phoneNumber: phoneNumber ?? this.phoneNumber,
      pinHash: pinHash ?? this.pinHash,
      marginScheme: marginScheme ?? this.marginScheme,
      fixedAllowanceAmount: fixedAllowanceAmount ?? this.fixedAllowanceAmount,
      marginSharePercentage: marginSharePercentage ?? this.marginSharePercentage,
      dailyLimitAmount: dailyLimitAmount ?? this.dailyLimitAmount,
      dailyLimitCount: dailyLimitCount ?? this.dailyLimitCount,
      isActive: isActive ?? this.isActive,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
    );
  }

  factory Staff.fromJson(Map<String, dynamic> json) {
    return Staff(
      id: json['id'] as String,
      mitraId: json['mitra_id'] as String,
      name: json['name'] as String,
      phoneNumber: json['phone_number'] as String,
      pinHash: json['pin_hash'] as String?,
      marginScheme: json['margin_scheme'] as String,
      fixedAllowanceAmount: (json['fixed_allowance_amount'] as num?)?.toDouble(),
      marginSharePercentage: (json['margin_share_percentage'] as num?)?.toDouble(),
      dailyLimitAmount: (json['daily_limit_amount'] as num).toDouble(),
      dailyLimitCount: json['daily_limit_count'] as int,
      isActive: json['is_active'] as bool? ?? true,
      createdAt: DateTime.parse(json['created_at'] as String),
      updatedAt: DateTime.parse(json['updated_at'] as String),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'mitra_id': mitraId,
      'name': name,
      'phone_number': phoneNumber,
      'pin_hash': pinHash,
      'margin_scheme': marginScheme,
      'fixed_allowance_amount': fixedAllowanceAmount,
      'margin_share_percentage': marginSharePercentage,
      'daily_limit_amount': dailyLimitAmount,
      'daily_limit_count': dailyLimitCount,
      'is_active': isActive,
      'created_at': createdAt.toIso8601String(),
      'updated_at': updatedAt.toIso8601String(),
    };
  }
}
