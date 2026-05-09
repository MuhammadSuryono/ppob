import 'package:hive/hive.dart';

part 'product.g.dart';

@HiveType(typeId: 1)
class Product {
  @HiveField(0)
  final String id;

  @HiveField(1)
  final String name;

  @HiveField(2)
  final String description;

  @HiveField(3)
  final String category; // 'pulsa', 'token-listrik', 'paket-data', etc.

  @HiveField(4)
  final String productCode; // Digiflazz product code

  @HiveField(5)
  final double platformPrice; // Cost from Digiflazz

  @HiveField(6)
  final double? sellingPrice; // Can be null (use platform price or custom)

  @HiveField(7)
  final String? imageUrl;

  @HiveField(8)
  final bool isActive;

  @HiveField(9)
  final DateTime lastSyncedAt;

  Product({
    required this.id,
    required this.name,
    required this.description,
    required this.category,
    required this.productCode,
    required this.platformPrice,
    this.sellingPrice,
    this.imageUrl,
    this.isActive = true,
    required this.lastSyncedAt,
  });

  Product copyWith({
    String? id,
    String? name,
    String? description,
    String? category,
    String? productCode,
    double? platformPrice,
    double? sellingPrice,
    String? imageUrl,
    bool? isActive,
    DateTime? lastSyncedAt,
  }) {
    return Product(
      id: id ?? this.id,
      name: name ?? this.name,
      description: description ?? this.description,
      category: category ?? this.category,
      productCode: productCode ?? this.productCode,
      platformPrice: platformPrice ?? this.platformPrice,
      sellingPrice: sellingPrice ?? this.sellingPrice,
      imageUrl: imageUrl ?? this.imageUrl,
      isActive: isActive ?? this.isActive,
      lastSyncedAt: lastSyncedAt ?? this.lastSyncedAt,
    );
  }

  factory Product.fromJson(Map<String, dynamic> json) {
    return Product(
      id: json['id'] as String? ?? json['product_code'] as String,
      name: json['name'] as String? ?? json['product_name'] as String,
      description: json['description'] as String? ?? '',
      category: json['category'] as String? ?? 'pulsa',
      productCode: json['product_code'] as String,
      platformPrice: (json['platform_price'] as num?)?.toDouble() ?? 0.0,
      sellingPrice: (json['selling_price'] as num?)?.toDouble(),
      imageUrl: json['image_url'] as String?,
      isActive: json['is_active'] as bool? ?? true,
      lastSyncedAt: json['last_synced_at'] != null
          ? DateTime.parse(json['last_synced_at'] as String)
          : DateTime.now(),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'description': description,
      'category': category,
      'product_code': productCode,
      'platform_price': platformPrice,
      'selling_price': sellingPrice,
      'image_url': imageUrl,
      'is_active': isActive,
      'last_synced_at': lastSyncedAt.toIso8601String(),
    };
  }
}
