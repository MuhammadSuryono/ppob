import 'package:hive/hive.dart';

part 'wallet.g.dart';

@HiveType(typeId: 4)
class Wallet {
  @HiveField(0)
  final String id;

  @HiveField(1)
  final String userId;

  @HiveField(2)
  final String role; // 'mitra' or 'staff'

  @HiveField(3)
  final String? ownerName; // For display

  @HiveField(4)
  final double availableBalance;

  @HiveField(5)
  final double heldBalance;

  @HiveField(6)
  final double dailySpentAmount; // Track daily spending for limits

  @HiveField(7)
  final int dailyTransactionCount;

  @HiveField(8)
  final DateTime date; // Date of the daily counters (YYYY-MM-DD)

  @HiveField(9)
  final DateTime updatedAt;

  Wallet({
    required this.id,
    required this.userId,
    required this.role,
    this.ownerName,
    required this.availableBalance,
    required this.heldBalance,
    this.dailySpentAmount = 0.0,
    this.dailyTransactionCount = 0,
    required this.date,
    required this.updatedAt,
  });

  double get totalBalance => availableBalance + heldBalance;

  Wallet copyWith({
    String? id,
    String? userId,
    String? role,
    String? ownerName,
    double? availableBalance,
    double? heldBalance,
    double? dailySpentAmount,
    int? dailyTransactionCount,
    DateTime? date,
    DateTime? updatedAt,
  }) {
    return Wallet(
      id: id ?? this.id,
      userId: userId ?? this.userId,
      role: role ?? this.role,
      ownerName: ownerName ?? this.ownerName,
      availableBalance: availableBalance ?? this.availableBalance,
      heldBalance: heldBalance ?? this.heldBalance,
      dailySpentAmount: dailySpentAmount ?? this.dailySpentAmount,
      dailyTransactionCount: dailyTransactionCount ?? this.dailyTransactionCount,
      date: date ?? this.date,
      updatedAt: updatedAt ?? this.updatedAt,
    );
  }

  factory Wallet.fromJson(Map<String, dynamic> json) {
    return Wallet(
      id: json['id'] as String,
      userId: json['user_id'] as String,
      role: json['role'] as String,
      ownerName: json['owner_name'] as String?,
      availableBalance: (json['available_balance'] as num).toDouble(),
      heldBalance: (json['held_balance'] as num).toDouble(),
      dailySpentAmount: (json['daily_spent_amount'] as num?)?.toDouble() ?? 0.0,
      dailyTransactionCount: json['daily_transaction_count'] as int? ?? 0,
      date: json['date'] != null ? DateTime.parse(json['date'] as String) : DateTime.now(),
      updatedAt: DateTime.parse(json['updated_at'] as String),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'user_id': userId,
      'role': role,
      'owner_name': ownerName,
      'available_balance': availableBalance,
      'held_balance': heldBalance,
      'daily_spent_amount': dailySpentAmount,
      'daily_transaction_count': dailyTransactionCount,
      'date': date.toIso8601String().split('T')[0], // YYYY-MM-DD
      'updated_at': updatedAt.toIso8601String(),
    };
  }
}
