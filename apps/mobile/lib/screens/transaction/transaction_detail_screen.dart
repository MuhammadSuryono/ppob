import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../providers/auth_provider.dart';
import '../../repositories/transaction_repository.dart';
import '../../models/transaction.dart';
import '../../utils/constants.dart';

class TransactionDetailScreen extends ConsumerStatefulWidget {
  const TransactionDetailScreen({super.key});

  @override
  ConsumerState<TransactionDetailScreen> createState() => _TransactionDetailScreenState();
}

class _TransactionDetailScreenState extends ConsumerState<TransactionDetailScreen> {
  @override
  Widget build(BuildContext context) {
    final args = ModalRoute.of(context)!.settings.arguments as String;
    final txId = args;

    return Scaffold(
      appBar: AppBar(
        title: const Text('Detail Transaksi'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => Navigator.of(context).pop(),
        ),
      ),
      body: FutureBuilder<Transaction>(
        future: ref.read(transactionRepositoryProvider).getTransactionDetail(txId),
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.waiting) {
            return const Center(child: CircularProgressIndicator());
          }

          if (snapshot.hasError) {
            return Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  const Icon(Icons.error_outline, size: 64, color: Colors.red),
                  const SizedBox(height: 16),
                  Text('Error: ${snapshot.error}'),
                ],
              ),
            );
          }

          final transaction = snapshot.data!;
          final statusColor = _getStatusColor(transaction.status);

          return SingleChildScrollView(
            padding: const EdgeInsets.all(16),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                // Status header
                Center(
                  child: Container(
                    padding: const EdgeInsets.all(16),
                    decoration: BoxDecoration(
                      color: statusColor.withOpacity(0.1),
                      shape: BoxShape.circle,
                    ),
                    child: Icon(
                      _getStatusIcon(transaction.status),
                      size: 64,
                      color: statusColor,
                    ),
                  ),
                ),
                const SizedBox(height: 16),

                Center(
                  child: Text(
                    _formatStatus(transaction.status),
                    style: Theme.of(context).textTheme.titleLarge?.copyWith(
                          color: statusColor,
                          fontWeight: FontWeight.bold,
                        ),
                  ),
                ),
                const SizedBox(height: 32),

                // Transaction info
                _buildSection('Informasi Transaksi'),
                _buildInfoRow('ID Transaksi', transaction.id),
                _buildInfoRow('Produk', transaction.productName),
                _buildInfoRow('Nomor Tujuan', transaction.customerNo),
                if (transaction.notes != null && transaction.notes!.isNotEmpty)
                  _buildInfoRow('Serial Number', transaction.notes!.join(', ')),
                _buildInfoRow(
                  'Tanggal',
                  '${transaction.createdAt.day}/${transaction.createdAt.month}/${transaction.createdAt.year} ${transaction.createdAt.hour}:${transaction.createdAt.minute.toString().padLeft(2, '0')}',
                ),
                const SizedBox(height: 24),

                // Payment info
                _buildSection('Pembayaran'),
                _buildInfoRow(
                  'Nominal',
                  AppConstants.formatCurrency(transaction.sellingPrice),
                  isBold: true,
                ),
                _buildInfoRow(
                  'Biaya Platform',
                  AppConstants.formatCurrency(transaction.platformPrice),
                ),
                if (transaction.sellingPrice > transaction.platformPrice)
                  _buildInfoRow(
                    'Margin',
                    AppConstants.formatCurrency(transaction.sellingPrice - transaction.platformPrice),
                    color: Colors.green,
                  ),
                const SizedBox(height: 24),

                // Response info (if any)
                if (transaction.rcCode != null || transaction.rcMessage != null) ...[
                  _buildSection('Respon'),
                  if (transaction.rcCode != null)
                    _buildInfoRow('Kode Respon', transaction.rcCode!),
                  if (transaction.rcMessage != null)
                    _buildInfoRow('Pesan', transaction.rcMessage!),
                  const SizedBox(height: 24),
                ],

                // Action buttons
                if (transaction.status == AppConstants.statusPending) ...[
                  SizedBox(
                    width: double.infinity,
                    child: ElevatedButton.icon(
                      onPressed: () => _checkStatus(transaction.id),
                      icon: const Icon(Icons.refresh),
                      label: const Text('Perbarui Status'),
                    ),
                  ),
                  const SizedBox(height: 12),
                ],
              ],
            ),
          );
        },
      ),
    );
  }

  Widget _buildSection(String title) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 12),
      child: Text(
        title,
        style: const TextStyle(fontSize: 16, fontWeight: FontWeight.w600),
      ),
    );
  }

  Widget _buildInfoRow(String label, String value, {bool isBold = false, Color? color}) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(
            label,
            style: TextStyle(color: Colors.grey[600]),
          ),
          Flexible(
            child: Text(
              value,
              textAlign: TextAlign.end,
              style: TextStyle(
                fontWeight: isBold ? FontWeight.bold : FontWeight.normal,
                color: color,
              ),
            ),
          ),
        ],
      ),
    );
  }

  Color _getStatusColor(String status) {
    switch (status) {
      case AppConstants.statusSuccess:
        return Colors.green;
      case AppConstants.statusPending:
        return Colors.orange;
      case AppConstants.statusFailed:
        return Colors.red;
      case AppConstants.statusInitiated:
        return Colors.blue;
      case AppConstants.statusExpired:
        return Colors.grey;
      default:
        return Colors.grey;
    }
  }

  IconData _getStatusIcon(String status) {
    switch (status) {
      case AppConstants.statusSuccess:
        return Icons.check_circle;
      case AppConstants.statusPending:
        return Icons.pending;
      case AppConstants.statusFailed:
        return Icons.cancel;
      case AppConstants.statusInitiated:
        return Icons.pending_actions;
      case AppConstants.statusExpired:
        return Icons.timer_off;
      default:
        return Icons.help_outline;
    }
  }

  String _formatStatus(String status) {
    switch (status) {
      case AppConstants.statusSuccess:
        return 'Transaksi Berhasil';
      case AppConstants.statusPending:
        return 'Sedang Diproses';
      case AppConstants.statusFailed:
        return 'Transaksi Gagal';
      case AppConstants.statusInitiated:
        return 'Menunggu Konfirmasi';
      case AppConstants.statusExpired:
        return 'Transaksi Kadaluarsa';
      default:
        return status;
    }
  }

  Future<void> _checkStatus(String transactionId) async {
    // Refresh transaction detail
    setState(() {});
    if (mounted) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('Status transaksi diperbarui'),
          backgroundColor: Colors.green,
        ),
      );
    }
  }
}
