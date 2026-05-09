import 'package:hive/hive.dart';

part 'transaction.g.dart';

@HiveType(typeId: 2)
class Transaction {
  @HiveField(0)
  final String id;

  @HiveField(1)
  final String productId;

  @HiveField(2)
  final String productName;

  @HiveField(3)
  final String customerNo; // Phone number or meter number

  @HiveField(4)
  final double sellingPrice;

  @HiveField(5)
  final double platformPrice;

  @HiveField(6)
  final String status; // Initiated, Pending, Success, Failed, Expired

  @HiveField(7)
  final String? rcCode; // Response code from Digiflazz

  @HiveField(8)
  final String? rcMessage; // Response message

  @HiveField(9)
  final String? staffId; // Who processed this (if staff transaction)

  @HiveField(10)
  final String? mitraId; // Which mitra processed this

  @HiveField(11)
  final DateTime createdAt;

  @HiveField(12)
  final DateTime? updatedAt;

  @HiveField(13)
  final List<String>? notes; // Additional notes/SN for token listrik

  Transaction({
    required this.id,
    required this.productId,
    required this.productName,
    required this.customerNo,
    required this.sellingPrice,
    required this.platformPrice,
    required this.status,
    this.rcCode,
    this.rcMessage,
    this.staffId,
    this.mitraId,
    required this.createdAt,
    this.updatedAt,
    this.notes,
  });

  Transaction copyWith({
    String? id,
    String? productId,
    String? productName,
    String? customerNo,
    double? sellingPrice,
    double? platformPrice,
    String? status,
    String? rcCode,
    String? rcMessage,
    String? staffId,
    String? mitraId,
    DateTime? createdAt,
    DateTime? updatedAt,
    List<String>? notes,
  }) {
    return Transaction(
      id: id ?? this.id,
      productId: productId ?? this.productId,
      productName: productName ?? this.productName,
      customerNo: customerNo ?? this.customerNo,
      sellingPrice: sellingPrice ?? this.sellingPrice,
      platformPrice: platformPrice ?? this.platformPrice,
      status: status ?? this.status,
      rcCode: rcCode ?? this.rcCode,
      rcMessage: rcMessage ?? this.rcMessage,
      staffId: staffId ?? this.staffId,
      mitraId: mitraId ?? this.mitraId,
      createdAt: createdAt ?? this.createdAt,
      updatedAt: updatedAt ?? this.updatedAt,
      notes: notes ?? this.notes,
    );
  }

  factory Transaction.fromJson(Map<String, dynamic> json) {
    return Transaction(
      id: json['id'] as String,
      productId: json['product_id'] as String,
      productName: json['product_name'] as String,
      customerNo: json['customer_no'] as String,
      sellingPrice: (json['selling_price'] as num).toDouble(),
      platformPrice: (json['platform_price'] as num).toDouble(),
      status: json['status'] as String,
      rcCode: json['rc_code'] as String?,
      rcMessage: json['rc_message'] as String?,
      staffId: json['staff_id'] as String?,
      mitraId: json['mitra_id'] as String?,
      createdAt: DateTime.parse(json['created_at'] as String),
      updatedAt: json['updated_at'] != null ? DateTime.parse(json['updated_at'] as String) : null,
      notes: json['notes'] != null ? List<String>.from(json['notes']) : null,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'product_id': productId,
      'product_name': productName,
      'customer_no': customerNo,
      'selling_price': sellingPrice,
      'platform_price': platformPrice,
      'status': status,
      'rc_code': rcCode,
      'rc_message': rcMessage,
      'staff_id': staffId,
      'mitra_id': mitraId,
      'created_at': createdAt.toIso8601String(),
      'updated_at': updatedAt?.toIso8601String(),
      'notes': notes,
    };
  }

  // Calculate margin (only applicable for staff transactions)
  double get margin {
    if (status != 'Success') return 0.0;
    return sellingPrice - platformPrice;
  }
}
