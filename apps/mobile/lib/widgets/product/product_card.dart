import 'package:flutter/material.dart';
import '../../models/product.dart';
import '../../utils/constants.dart';

class ProductCard extends StatelessWidget {
  final Product product;
  final VoidCallback? onTap;

  const ProductCard({
    super.key,
    required this.product,
    this.onTap,
  });

  @override
  Widget build(BuildContext context) {
    final displayPrice = product.sellingPrice ?? product.platformPrice;

    return Card(
      elevation: 2,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(12),
      ),
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(12),
        child: Padding(
          padding: const EdgeInsets.all(12.0),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // Product image/icon placeholder
              Expanded(
                flex: 2,
                child: Container(
                  width: double.infinity,
                  decoration: BoxDecoration(
                    color: _getCategoryColor().withOpacity(0.1),
                    borderRadius: BorderRadius.circular(8),
                  ),
                  child: Icon(
                    _getCategoryIcon(),
                    size: 40,
                    color: _getCategoryColor(),
                  ),
                ),
              ),

              const SizedBox(height: 8),

              // Product name
              Expanded(
                flex: 1,
                child: Text(
                  product.name,
                  style: Theme.of(context).textTheme.titleSmall?.copyWith(
                        fontWeight: FontWeight.w600,
                      ),
                  maxLines: 2,
                  overflow: TextOverflow.ellipsis,
                ),
              ),

              const SizedBox(height: 4),

              // Price
              Text(
                AppConstants.formatCurrency(displayPrice),
                style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                      color: Theme.of(context).colorScheme.primary,
                      fontWeight: FontWeight.bold,
                    ),
              ),

              // Platform price (strikethrough) if selling price differs
              if (product.sellingPrice != null && product.sellingPrice! > product.platformPrice)
                Text(
                  AppConstants.formatCurrency(product.platformPrice),
                  style: Theme.of(context).textTheme.bodySmall?.copyWith(
                        color: Colors.grey[500],
                        decoration: TextDecoration.lineThrough,
                      ),
                ),
            ],
          ),
        ),
      ),
    );
  }

  Color _getCategoryColor() {
    switch (product.category) {
      case 'pulsa':
        return Colors.blue;
      case 'token-listrik':
        return Colors.orange;
      case 'paket-data':
        return Colors.purple;
      default:
        return Colors.green;
    }
  }

  IconData _getCategoryIcon() {
    switch (product.category) {
      case 'pulsa':
        return Icons.phone_android;
      case 'token-listrik':
        return Icons.electrical_services;
      case 'paket-data':
        return Icons.wifi;
      default:
        return Icons.card_giftcard;
    }
  }
}
