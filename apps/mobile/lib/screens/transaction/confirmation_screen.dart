import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../providers/auth_provider.dart';
import '../../repositories/transaction_repository.dart';
import '../../repositories/wallet_repository.dart';
import '../../models/product.dart';
import '../../utils/constants.dart';

class ConfirmationScreen extends ConsumerStatefulWidget {
  const ConfirmationScreen({super.key});

  @override
  ConsumerState<ConfirmationScreen> createState() => _ConfirmationScreenState();
}

class _ConfirmationScreenState extends ConsumerState<ConfirmationScreen> {
  bool _isLoading = false;

  @override
  Widget build(BuildContext context) {
    final args = ModalRoute.of(context)!.settings.arguments as Map<String, dynamic>;
    final product = args['product'] as Product;
    final customerNo = args['customerNo'] as String;
    final finalPrice = args['finalPrice'] as double;
    final nominal = args['nominal'] as double?;
    final provider = args['provider'] as String?;
    final notes = args['notes'] as List<String>?;

    final authState = ref.watch(authProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Konfirmasi'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => Navigator.of(context).pop(),
        ),
      ),
      body: SafeArea(
        child: Column(
          children: [
            Expanded(
              child: SingleChildScrollView(
                padding: const EdgeInsets.all(16),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    // Title
                    const Text(
                      'Konfirmasi Transaksi',
                      style: TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
                    ),
                    const SizedBox(height: 24),

                    // Product details
                    _buildSectionTitle('Detail Produk'),
                    _buildDetailRow('Produk', product.name),
                    if (nominal != null) _buildDetailRow('Nominal', AppConstants.formatCurrency(nominal)),
                    if (provider != null) _buildDetailRow('Provider', provider),
                    _buildDetailRow('Harga', AppConstants.formatCurrency(finalPrice)),
                    const Divider(height: 32),

                    // Customer details
                    _buildSectionTitle('Detail Pelanggan'),
                    _buildDetailRow('Nomor Tujuan', customerNo),
                    if (notes != null && notes.isNotEmpty)
                      _buildDetailRow('SN / Catatan', notes.first),
                    const Divider(height: 32),

                    // Summary
                    _buildSectionTitle('Ringkasan Pembayaran'),
                    Container(
                      padding: const EdgeInsets.all(16),
                      decoration: BoxDecoration(
                        color: Colors.grey[100],
                        borderRadius: BorderRadius.circular(8),
                      ),
                      child: Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          const Text(
                            'Total Dibayar:',
                            style: TextStyle(fontSize: 16),
                          ),
                          Text(
                            AppConstants.formatCurrency(finalPrice),
                            style: const TextStyle(
                              fontSize: 20,
                              fontWeight: FontWeight.bold,
                              color: Colors.green,
                            ),
                          ),
                        ],
                      ),
                    ),
                  ],
                ),
              ),
            ),

            // Bottom action button
            Container(
              padding: const EdgeInsets.all(16),
              decoration: BoxDecoration(
                color: Colors.white,
                boxShadow: [
                  BoxShadow(
                    color: Colors.grey.withOpacity(0.2),
                    spreadRadius: 1,
                    blurRadius: 4,
                    offset: const Offset(0, -2),
                  ),
                ],
              ),
              child: SizedBox(
                width: double.infinity,
                child: ElevatedButton(
                  onPressed: _isLoading ? null : () => _confirmTransaction(args),
                  style: ElevatedButton.styleFrom(
                    backgroundColor: Colors.green,
                    foregroundColor: Colors.white,
                    padding: const EdgeInsets.symmetric(vertical: 16),
                  ),
                  child: _isLoading
                      ? const SizedBox(
                          height: 20,
                          width: 20,
                          child: CircularProgressIndicator(
                            strokeWidth: 2,
                            valueColor: AlwaysStoppedAnimation<Color>(Colors.white)),
                        )
                      : const Text('BAYAR & PROSES', style: TextStyle(fontSize: 16)),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildSectionTitle(String title) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 12),
      child: Text(
        title,
        style: const TextStyle(fontSize: 16, fontWeight: FontWeight.w600),
      ),
    );
  }

  Widget _buildDetailRow(String label, String value) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(
            label,
            style: TextStyle(color: Colors.grey[600]),
          ),
          Text(
            value,
            style: const TextStyle(fontWeight: FontWeight.w500),
          ),
        ],
      ),
    );
  }

  Future<void> _confirmTransaction(Map<String, dynamic> args) async {
    setState(() => _isLoading = true);

    try {
      final product = args['product'] as Product;
      final customerNo = args['customerNo'] as String;
      final finalPrice = args['finalPrice'] as double;
      final nominal = args['nominal'] as double?;
      final notes = args['notes'] as List<String>?;
      final authState = ref.read(authProvider);
      final walletRepo = ref.read(walletRepositoryProvider);
      final txRepo = ref.read(transactionRepositoryProvider);

      // First, hold the amount from wallet
      final wallet = await walletRepo.getActiveWallet(authState.userId!);
      await walletRepo.holdAmount(wallet.id, finalPrice, 'tx_hold_${DateTime.now().millisecondsSinceEpoch}');

      if (!mounted) return;

      // Navigate to PIN screen
      final result = await Navigator.of(context).pushNamed<bool>(
        '/transaction/pin',
        arguments: {
          'amount': finalPrice,
          'product': product,
          'customerNo': customerNo,
        },
      );

      if (result == true) {
        // Transaction confirmed via PIN
        final tx = await txRepo.initiateTransaction(
          productId: product.id,
          productName: product.name,
          customerNo: customerNo,
          sellingPrice: finalPrice,
          platformPrice: product.platformPrice,
          staffId: authState.role == 'staff' ? authState.userId : null,
          mitraId: authState.role == 'mitra' ? authState.userId : null,
        );

        if (!mounted) return;

        // Show success dialog
        await showDialog(
          context: context,
          barrierDismissible: false,
          builder: (ctx) => AlertDialog(
            icon: const Icon(Icons.check_circle, color: Colors.green, size: 64),
            title: const Text('Transaksi Berhasil'),
            content: Text('Transaksi ${product.name} berhasil diproses.\nID: ${tx.id}'),
            actions: [
              TextButton(
                onPressed: () {
                  Navigator.of(ctx).pop();
                  Navigator.of(context).popUntil((route) => route.isFirst);
                },
                child: const Text('OK'),
              ),
            ],
          ),
        );
      } else {
        // Transaction cancelled or PIN failed
        if (!mounted) return;
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Transaksi dibatalkan'),
            backgroundColor: Colors.orange,
          ),
        );
        Navigator.of(context).pop();
      }
    } catch (e) {
      if (!mounted) return;
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(
          content: Text('Gagal: ${e.toString()}'),
          backgroundColor: Colors.red,
        ),
      );
    } finally {
      if (mounted) {
        setState(() => _isLoading = false);
      }
    }
  }
}
