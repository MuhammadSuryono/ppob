import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../models/product.dart';
import '../../models/wallet.dart';
import '../../providers/auth_provider.dart';
import '../../repositories/transaction_repository.dart';
import '../../repositories/wallet_repository.dart';
import '../../utils/constants.dart';

class TransactionInitiateScreen extends ConsumerStatefulWidget {
  final Product product;

  const TransactionInitiateScreen({
    super.key,
    required this.product,
  });

  @override
  ConsumerState<TransactionInitiateScreen> createState() => _TransactionInitiateScreenState();
}

class _TransactionInitiateScreenState extends ConsumerState<TransactionInitiateScreen> {
  final _formKey = GlobalKey<FormState>();
  final _customerNoController = TextEditingController();
  final _noteController = TextEditingController();

  String? _selectedProvider; // For phone: Telkomsel, XL, etc.
  bool _isLoading = false;
  double? _finalPrice;
  double? _selectedNominal;

  final Map<String, List<double>> _providerNominals = {
    'Telkomsel': [10000, 25000, 50000, 100000],
    'XL': [10000, 25000, 50000, 100000],
    'Indosat': [10000, 25000, 50000],
    'Tri': [10000, 25000, 50000],
    'Smartfren': [10000, 25000, 50000],
  };

  @override
  void initState() {
    super.initState();
    // For token listrik, no provider selection needed
    if (widget.product.category != 'token-listrik') {
      _selectedProvider = _providerNominals.keys.first;
    }
  }

  @override
  void dispose() {
    _customerNoController.dispose();
    _noteController.dispose();
    super.dispose();
  }

  void _calculatePrice() {
    double basePrice = widget.product.platformPrice;
    if (widget.product.sellingPrice != null) {
      basePrice = widget.product.sellingPrice!;
    }

    if (_selectedNominal != null) {
      setState(() {
        _finalPrice = basePrice * (_selectedNominal! / 10000);
      });
    } else {
      setState(() {
        _finalPrice = basePrice;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    final authState = ref.watch(authProvider);
    final walletRepository = ref.watch(walletRepositoryProvider);

    _calculatePrice();

    return Scaffold(
      appBar: AppBar(
        title: const Text('Transaksi'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => Navigator.of(context).pop(),
        ),
      ),
      body: SafeArea(
        child: SingleChildScrollView(
          padding: const EdgeInsets.all(16),
          child: Form(
            key: _formKey,
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                // Product summary
                Card(
                  child: Padding(
                    padding: const EdgeInsets.all(16),
                    child: Row(
                      children: [
                        Container(
                          width: 60,
                          height: 60,
                          decoration: BoxDecoration(
                            color: Theme.of(context).colorScheme.primaryContainer,
                            borderRadius: BorderRadius.circular(8),
                          ),
                          child: Icon(
                            _getCategoryIcon(),
                            size: 30,
                            color: Theme.of(context).colorScheme.primary,
                          ),
                        ),
                        const SizedBox(width: 16),
                        Expanded(
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Text(
                                widget.product.name,
                                style: Theme.of(context).textTheme.titleMedium?.copyWith(
                                      fontWeight: FontWeight.bold,
                                    ),
                              ),
                              const SizedBox(height: 4),
                              Text(
                                widget.product.description,
                                style: Theme.of(context).textTheme.bodySmall,
                              ),
                            ],
                          ),
                        ),
                        Column(
                          crossAxisAlignment: CrossAxisAlignment.end,
                          children: [
                            Text(
                              _finalPrice != null ? AppConstants.formatCurrency(_finalPrice!) : 'Rp -',
                              style: Theme.of(context).textTheme.titleMedium?.copyWith(
                                    color: Theme.of(context).colorScheme.primary,
                                    fontWeight: FontWeight.bold,
                                  ),
                            ),
                            if (widget.product.sellingPrice != null)
                              Text(
                                AppConstants.formatCurrency(widget.product.platformPrice),
                                style: const TextStyle(
                                  fontSize: 12,
                                  decoration: TextDecoration.lineThrough,
                                  color: Colors.grey,
                                ),
                              ),
                          ],
                        ),
                      ],
                    ),
                  ),
                ),

                const SizedBox(height: 24),

                // Customer input section
                Text(
                  'Informasi Pelanggan',
                  style: Theme.of(context).textTheme.titleMedium?.copyWith(
                        fontWeight: FontWeight.bold,
                      ),
                ),
                const SizedBox(height: 16),

                // Provider selection (for non-token products)
                if (widget.product.category != 'token-listrik') ...[
                  DropdownButtonFormField<String>(
                    decoration: const InputDecoration(
                      labelText: 'Provider',
                      prefixIcon: Icon(Icons.network_cell),
                    ),
                    value: _selectedProvider,
                    items: _providerNominals.keys
                        .map((provider) => DropdownMenuItem(
                              value: provider,
                              child: Text(provider),
                            ))
                        .toList(),
                    onChanged: (value) {
                      setState(() {
                        _selectedProvider = value;
                        _selectedNominal = null;
                      });
                    },
                    validator: (value) => value == null ? 'Pilih provider' : null,
                  ),
                  const SizedBox(height: 16),
                ],

                // Customer number input
                TextFormField(
                  controller: _customerNoController,
                  decoration: InputDecoration(
                    labelText: widget.product.category == 'token-listrik'
                        ? 'Nomor Meter / ID Pelanggan'
                        : 'Nomor HP',
                    prefixIcon: Icon(widget.product.category == 'token-listrik'
                        ? Icons.electrical_services
                        : Icons.phone),
                    hintText: widget.product.category == 'token-listrik'
                        ? 'Masukkan nomor meter listrik'
                        : '08xxxxxxxxxx',
                  ),
                  keyboardType: TextInputType.phone,
                  validator: (value) {
                    if (value == null || value.trim().isEmpty) {
                      return 'Nomor wajib diisi';
                    }
                    if (value.length < 8) {
                      return 'Nomor terlalu pendek';
                    }
                    return null;
                  },
                ),

                // Nominal selection (for non-token)
                if (widget.product.category != 'token-listrik' && _selectedProvider != null) ...[
                  const SizedBox(height: 16),
                  const Text('Pilih Nominal:'),
                  const SizedBox(height: 8),
                  Wrap(
                    spacing: 8,
                    runSpacing: 8,
                    children: _providerNominals[_selectedProvider]!
                        .map((nominal) => ChoiceChip(
                              label: Text(AppConstants.formatCurrency(nominal.toDouble())),
                              selected: _selectedNominal == nominal,
                              onSelected: (selected) {
                                setState(() {
                                  _selectedNominal = selected ? nominal.toDouble() : null;
                                });
                              },
                            ))
                        .toList(),
                  ),
                ],

                // Additional notes (for token listrik - SN)
                if (widget.product.category == 'token-listrik') ...[
                  const SizedBox(height: 16),
                  TextFormField(
                    controller: _noteController,
                    decoration: const InputDecoration(
                      labelText: 'Serial Number (SN) - Optional',
                      prefixIcon: Icon(Icons.numbers),
                      hintText: 'Masukkan SN token jika ada',
                    ),
                    maxLines: 2,
                  ),
                ],

                const SizedBox(height: 32),

                // Proceed button
                SizedBox(
                  width: double.infinity,
                  child: ElevatedButton(
                    onPressed: _isLoading
                        ? null
                        : _finalPrice != null && _finalPrice! > 0
                            ? () => _proceedToConfirmation(authState, walletRepository)
                            : null,
                    child: _isLoading
                        ? const SizedBox(
                            height: 20,
                            width: 20,
                            child: CircularProgressIndicator(
                              strokeWidth: 2,
                              valueColor: AlwaysStoppedAnimation<Color>(Colors.white)),
                          )
                        : const Text('LANJUTKAN'),
                  ),
                ),

                // Wallet balance info
                const SizedBox(height: 16),
                FutureBuilder<Wallet>(
                  future: walletRepository.getActiveWallet(authState.userId!),
                  builder: (context, snapshot) {
                    if (snapshot.connectionState == ConnectionState.waiting) {
                      return const SizedBox.shrink();
                    }
                    final wallet = snapshot.data;
                    if (wallet == null) return const SizedBox.shrink();

                    final canAfford = wallet.availableBalance >= (_finalPrice ?? 0);

                    return Container(
                      padding: const EdgeInsets.all(12),
                      decoration: BoxDecoration(
                        color: canAfford ? Colors.green.shade50 : Colors.red.shade50,
                        borderRadius: BorderRadius.circular(8),
                        border: Border.all(
                          color: canAfford ? Colors.green : Colors.red,
                        ),
                      ),
                      child: Row(
                        children: [
                          Icon(
                            Icons.account_balance_wallet,
                            color: canAfford ? Colors.green : Colors.red,
                          ),
                          const SizedBox(width: 8),
                          Expanded(
                            child: Text(
                              'Saldo Tersedia: ${AppConstants.formatCurrency(wallet.availableBalance)}',
                              style: TextStyle(
                                color: canAfford ? Colors.green : Colors.red,
                                fontWeight: FontWeight.w600,
                              ),
                            ),
                          ),
                        ],
                      ),
                    );
                  },
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  IconData _getCategoryIcon() {
    switch (widget.product.category) {
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

  Future<void> _proceedToConfirmation(AuthState authState, WalletRepository walletRepo) async {
    if (!_formKey.currentState!.validate()) return;

    setState(() => _isLoading = true);

    try {
      // Check wallet balance
      final wallet = await walletRepo.getActiveWallet(authState.userId!);
      if (wallet.availableBalance < (_finalPrice ?? 0)) {
        throw Exception('Saldo tidak mencukupi');
      }

      if (!mounted) return;

      // Navigate to confirmation screen
      Navigator.of(context).pushNamed(
        '/transaction/confirm',
        arguments: {
          'product': widget.product,
          'customerNo': _customerNoController.text.trim(),
          'finalPrice': _finalPrice,
          'nominal': _selectedNominal,
          'provider': _selectedProvider,
          'notes': _noteController.text.isNotEmpty ? [_noteController.text.trim()] : null,
        },
      );
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Error: ${e.toString()}'),
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
