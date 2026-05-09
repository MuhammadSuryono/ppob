import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../providers/auth_provider.dart';
import '../../repositories/wallet_repository.dart';
import '../../repositories/staff_repository.dart';
import '../../models/staff.dart';
import '../../models/wallet.dart';
import '../../utils/constants.dart';

class StaffTopUpScreen extends ConsumerStatefulWidget {
  const StaffTopUpScreen({super.key});

  @override
  ConsumerState<StaffTopUpScreen> createState() => _StaffTopUpScreenState();
}

class _StaffTopUpScreenState extends ConsumerState<StaffTopUpScreen> {
  final _amountController = TextEditingController();
  final _formKey = GlobalKey<FormState>();
  bool _isLoading = false;

  // Preset amounts
  final List<double> _presetAmounts = [
    100000,
    250000,
    500000,
    1000000,
    2500000,
    5000000,
  ];

  @override
  void dispose() {
    _amountController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final args = ModalRoute.of(context)!.settings.arguments as String;
    final staffId = args;

    final staffRepo = ref.watch(staffRepositoryProvider);
    final walletRepo = ref.watch(walletRepositoryProvider);
    final authState = ref.watch(authProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Top Up Staff'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => Navigator.of(context).pop(),
        ),
      ),
      body: SafeArea(
        child: FutureBuilder(
          future: Future.wait<dynamic>([
            staffRepo.getStaffDetail(staffId),
            walletRepo.getActiveWallet(authState.userId!),
          ]),
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
                    const SizedBox(height: 16),
                    ElevatedButton(
                      onPressed: () => Navigator.of(context).pop(),
                      child: const Text('Kembali'),
                    ),
                  ],
                ),
              );
            }

            final results = snapshot.data!;
            final staff = results[0] as Staff;
            final mitraWallet = results[1] as Wallet;

            return SingleChildScrollView(
              padding: const EdgeInsets.all(16),
              child: Form(
                key: _formKey,
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    // Staff info
                    Card(
                      child: Padding(
                        padding: const EdgeInsets.all(16),
                        child: Row(
                          children: [
                            CircleAvatar(
                              backgroundColor: Theme.of(context).colorScheme.primaryContainer,
                              child: Text(
                                staff.name.isNotEmpty ? staff.name[0].toUpperCase() : '?',
                                style: TextStyle(
                                  fontWeight: FontWeight.bold,
                                  color: Theme.of(context).colorScheme.primary,
                                ),
                              ),
                            ),
                            const SizedBox(width: 16),
                            Expanded(
                              child: Column(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                children: [
                                  Text(
                                    staff.name,
                                    style: const TextStyle(
                                      fontWeight: FontWeight.w600,
                                      fontSize: 16,
                                    ),
                                  ),
                                  Text(
                                    staff.phoneNumber,
                                    style: TextStyle(
                                      color: Colors.grey[600],
                                      fontSize: 12,
                                    ),
                                  ),
                                ],
                              ),
                            ),
                          ],
                        ),
                      ),
                    ),

                    const SizedBox(height: 24),

                    // Mitra balance info
                    Container(
                      padding: const EdgeInsets.all(16),
                      decoration: BoxDecoration(
                        color: Colors.green.shade50,
                        borderRadius: BorderRadius.circular(12),
                        border: Border.all(color: Colors.green.shade200),
                      ),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Row(
                            children: [
                              Icon(Icons.account_balance_wallet, color: Colors.green[700]),
                              const SizedBox(width: 8),
                              const Text(
                                'Saldo Mitra Tersedia',
                                style: TextStyle(fontWeight: FontWeight.w600),
                              ),
                            ],
                          ),
                          const SizedBox(height: 4),
                          Text(
                            AppConstants.formatCurrency(mitraWallet.availableBalance),
                            style: TextStyle(
                              fontSize: 24,
                              fontWeight: FontWeight.bold,
                              color: Colors.green[700],
                            ),
                          ),
                        ],
                      ),
                    ),

                    const SizedBox(height: 24),

                    // Amount selection
                    const Text(
                      'Pilih Nominal',
                      style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold),
                    ),
                    const SizedBox(height: 12),

                    // Preset amounts
                    Wrap(
                      spacing: 8,
                      runSpacing: 8,
                      children: _presetAmounts.map((amount) {
                        return FilterChip(
                          label: Text(AppConstants.formatCurrency(amount)),
                          selected: _amountController.text == amount.toString(),
                          onSelected: (selected) {
                            setState(() {
                              _amountController.text = amount.toString();
                            });
                          },
                        );
                      }).toList(),
                    ),

                    const SizedBox(height: 16),

                    // Custom amount input
                    TextFormField(
                      controller: _amountController,
                      decoration: const InputDecoration(
                        labelText: 'Jumlah Lain (Rp)',
                        prefixIcon: Icon(Icons.edit),
                      ),
                      keyboardType: TextInputType.number,
                      validator: (value) {
                        if (value == null || value.isEmpty) {
                          return 'Jumlah wajib diisi';
                        }
                        final amount = double.tryParse(value);
                        if (amount == null || amount <= 0) {
                          return 'Jumlah tidak valid';
                        }
                        if (amount > mitraWallet.availableBalance) {
                          return 'Saldo tidak mencukupi';
                        }
                        return null;
                      },
                      onChanged: (value) {
                        setState(() {});
                      },
                    ),

                    const SizedBox(height: 24),

                    // Proceed button
                    SizedBox(
                      width: double.infinity,
                      child: ElevatedButton(
                        onPressed: _isLoading
                            ? null
                            : _amountController.text.isNotEmpty
                                ? () => _processTopUp(staff, mitraWallet)
                                : null,
                        child: _isLoading
                            ? const SizedBox(
                                height: 20,
                                width: 20,
                                child: CircularProgressIndicator(
                                  strokeWidth: 2,
                                  valueColor: AlwaysStoppedAnimation<Color>(Colors.white)),
                              )
                            : const Text('PROSES TOP UP'),
                      ),
                    ),

                    // Warning message
                    const SizedBox(height: 16),
                    Container(
                      padding: const EdgeInsets.all(12),
                      decoration: BoxDecoration(
                        color: Colors.orange.shade50,
                        borderRadius: BorderRadius.circular(8),
                        border: Border.all(color: Colors.orange.shade200),
                      ),
                      child: Row(
                        children: [
                          Icon(Icons.warning, color: Colors.orange[700], size: 20),
                          const SizedBox(width: 8),
                          Expanded(
                            child: Text(
                              'Top up tidak dapat dibatalkan. Pastikan amount sudah benar.',
                              style: TextStyle(
                                fontSize: 12,
                                color: Colors.orange[700],
                              ),
                            ),
                          ),
                        ],
                      ),
                    ),
                  ],
                ),
              ),
            );
          },
        ),
      ),
    );
  }

  Future<void> _processTopUp(Staff staff, Wallet mitraWallet) async {
    if (!_formKey.currentState!.validate()) return;

    final amount = double.parse(_amountController.text);

    setState(() => _isLoading = true);

    try {
      final walletRepo = ref.read(walletRepositoryProvider);

      await walletRepo.topUpStaff(staff.id, amount);

      if (mounted) {
        await showDialog(
          context: context,
          barrierDismissible: false,
          builder: (ctx) => AlertDialog(
            icon: const Icon(Icons.check_circle, color: Colors.green, size: 64),
            title: const Text('Top Up Berhasil'),
            content: Text(
              'Top up Rp ${amount.toStringAsFixed(0)} ke ${staff.name} berhasil.',
            ),
            actions: [
              TextButton(
                onPressed: () {
                  Navigator.of(ctx).pop();
                  Navigator.of(context).pop(true); // Return success
                },
                child: const Text('OK'),
              ),
            ],
          ),
        );
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Gagal: ${e.toString()}'),
            backgroundColor: Colors.red,
          ),
        );
      }
    } finally {
      if (mounted) {
        setState(() => _isLoading = false);
      }
    }
  }
}
