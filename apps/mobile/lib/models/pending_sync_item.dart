import 'package:hive/hive.dart';

part 'pending_sync_item.g.dart';

@HiveType(typeId: 5)
class PendingSyncItem {
  @HiveField(0)
  final String id;

  @HiveField(1)
  final String type; // 'transaction', 'topup', etc.

  @HiveField(2)
  final Map<String, dynamic> data;

  @HiveField(3)
  final int retryCount;

  @HiveField(4)
  final DateTime createdAt;

  @HiveField(5)
  final DateTime? lastAttemptAt;

  @HiveField(6)
  final String? errorMessage;

  PendingSyncItem({
    required this.id,
    required this.type,
    required this.data,
    this.retryCount = 0,
    required this.createdAt,
    this.lastAttemptAt,
    this.errorMessage,
  });

  PendingSyncItem copyWith({
    String? id,
    String? type,
    Map<String, dynamic>? data,
    int? retryCount,
    DateTime? createdAt,
    DateTime? lastAttemptAt,
    String? errorMessage,
  }) {
    return PendingSyncItem(
      id: id ?? this.id,
      type: type ?? this.type,
      data: data ?? this.data,
      retryCount: retryCount ?? this.retryCount,
      createdAt: createdAt ?? this.createdAt,
      lastAttemptAt: lastAttemptAt ?? this.lastAttemptAt,
      errorMessage: errorMessage ?? this.errorMessage,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'type': type,
      'data': data,
      'retry_count': retryCount,
      'created_at': createdAt.toIso8601String(),
      'last_attempt_at': lastAttemptAt?.toIso8601String(),
      'error_message': errorMessage,
    };
  }

  factory PendingSyncItem.fromJson(Map<String, dynamic> json) {
    return PendingSyncItem(
      id: json['id'] as String,
      type: json['type'] as String,
      data: Map<String, dynamic>.from(json['data'] as Map),
      retryCount: json['retry_count'] as int? ?? 0,
      createdAt: DateTime.parse(json['created_at'] as String),
      lastAttemptAt: json['last_attempt_at'] != null
          ? DateTime.parse(json['last_attempt_at'] as String)
          : null,
      errorMessage: json['error_message'] as String?,
    );
  }
}
