import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:local_auth/local_auth.dart';
import '../../providers/auth_provider.dart';
import '../../utils/constants.dart';
import '../../services/biometric_service.dart';

class PinEntryScreen extends ConsumerStatefulWidget {
  const PinEntryScreen({super.key});

  @override
  ConsumerState<PinEntryScreen> createState() => _PinEntryScreenState();
}

class _PinEntryScreenState extends ConsumerState<PinEntryScreen> {
  final _pinController = TextEditingController();
  final LocalAuthentication _localAuth = LocalAuthentication();
  final BiometricService _biometricService = BiometricService();
  bool _isLoading = false;
  String? _errorMessage;
  bool _biometricAvailable = false;

  @override
  void initState() {
    super.initState();
    _checkBiometricAvailability();
  }

  Future<void> _checkBiometricAvailability() async {
    final available = await _biometricService.isAvailable();
    if (mounted) {
      setState(() {
        _biometricAvailable = available;
      });
    }
  }

  @override
  void dispose() {
    _pinController.dispose();
    super.dispose();
  }

  Future<void> _handleBiometricAuth() async {
    try {
      setState(() => _isLoading = true);

      final didAuthenticate = await _biometricService.authenticate();

      if (didAuthenticate && mounted) {
        Navigator.of(context).pop(true);
      }
    } catch (e) {
      if (mounted) {
        setState(() {
          _errorMessage = 'Biometric gagal: ${e.toString()}';
        });
      }
    } finally {
      if (mounted) {
        setState(() => _isLoading = false);
      }
    }
  }

  Future<void> _verifyPIN() async {
    final pin = _pinController.text.trim();

    if (pin.length != 6) {
      setState(() {
        _errorMessage = 'PIN harus 6 digit';
      });
      return;
    }

    setState(() {
      _isLoading = true;
      _errorMessage = null;
    });

    // Simulate PIN verification
    await Future.delayed(const Duration(milliseconds: 500));

    if (mounted) {
      Navigator.of(context).pop(true);
    }
  }

  void _showForgotPinDialog() {
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('Lupa PIN?'),
        content: const Text('Hubungi admin untuk reset PIN transaksi Anda.'),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(ctx).pop(),
            child: const Text('OK'),
          ),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final args = GoRouterState.of(context).extra as Map<String, dynamic>;
    final amount = args['amount'] as double;
    final product = args['product'] as dynamic;
    final customerNo = args['customerNo'] as String;

    return Scaffold(
      appBar: AppBar(
        title: const Text('Konfirmasi PIN'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => Navigator.of(context).pop(false),
        ),
      ),
      body: SafeArea(
        child: Padding(
          padding: const EdgeInsets.all(24),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.stretch,
            children: [
              // Transaction summary
              Container(
                padding: const EdgeInsets.all(16),
                decoration: BoxDecoration(
                  color: Colors.green.shade50,
                  borderRadius: BorderRadius.circular(12),
                ),
                child: Column(
                  children: [
                    const Text(
                      'Total Pembayaran',
                      style: TextStyle(color: Colors.grey),
                    ),
                    const SizedBox(height: 8),
                    Text(
                      AppConstants.formatCurrency(amount),
                      style: const TextStyle(
                        fontSize: 32,
                        fontWeight: FontWeight.bold,
                        color: Colors.green,
                      ),
                    ),
                    const SizedBox(height: 16),
                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        Expanded(
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Text(
                                product.name,
                                style: const TextStyle(fontWeight: FontWeight.w600),
                              ),
                              Text(
                                'Ke: $customerNo',
                                style: TextStyle(
                                  fontSize: 12,
                                  color: Colors.grey[600],
                                ),
                              ),
                            ],
                          ),
                        ),
                      ],
                    ),
                  ],
                ),
              ),

              const SizedBox(height: 32),

              // PIN entry
              const Text(
                'Masukkan PIN Transaksi',
                style: TextStyle(fontSize: 16, fontWeight: FontWeight.w600),
                textAlign: TextAlign.center,
              ),
              const SizedBox(height: 8),
              Text(
                '6 digit PIN yang terdaftar',
                style: TextStyle(fontSize: 12, color: Colors.grey[600]),
                textAlign: TextAlign.center,
              ),
              const SizedBox(height: 24),

              // PIN input
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: List.generate(6, (index) {
                  return Container(
                    margin: const EdgeInsets.symmetric(horizontal: 4),
                    width: 50,
                    height: 60,
                    decoration: BoxDecoration(
                      border: Border.all(color: Colors.grey),
                      borderRadius: BorderRadius.circular(8),
                    ),
                    child: Center(
                      child: Text(
                        index < _pinController.text.length
                            ? '•'
                            : '',
                        style: const TextStyle(fontSize: 24),
                      ),
                    ),
                  );
                }),
              ),

              // Hidden text field for input
              SizedBox(
                height: 0,
                child: TextField(
                  controller: _pinController,
                  keyboardType: TextInputType.number,
                  maxLength: 6,
                  obscureText: true,
                  onChanged: (value) {
                    if (value.length == 6) {
                      _verifyPIN();
                    }
                  },
                ),
              ),

              const SizedBox(height: 16),

              // Error message
              if (_errorMessage != null)
                Text(
                  _errorMessage!,
                  style: const TextStyle(color: Colors.red),
                  textAlign: TextAlign.center,
                ),

              const SizedBox(height: 24),

              // Number pad
              _buildNumberPad(),

               const Spacer(),

               // Biometric button
              if (_biometricAvailable)
                ElevatedButton.icon(
                  onPressed: _isLoading ? null : _handleBiometricAuth,
                  icon: const Icon(Icons.fingerprint),
                  label: const Text('Gunakan Fingerprint/Face ID'),
                  style: ElevatedButton.styleFrom(
                    backgroundColor: Colors.transparent,
                    foregroundColor: Theme.of(context).colorScheme.primary,
                    elevation: 0,
                  ),
                ),

              TextButton(
                onPressed: _showForgotPinDialog,
                child: const Text('Lupa PIN?'),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildNumberPad() {
    return Column(
      children: [
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceEvenly,
          children: List.generate(3, (rowIndex) {
            final number = rowIndex * 3 + 1;
            return _buildNumberButton(number.toString());
          }),
        ),
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceEvenly,
          children: List.generate(3, (rowIndex) {
            final number = rowIndex * 3 + 4;
            return _buildNumberButton(number.toString());
          }),
        ),
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceEvenly,
          children: List.generate(3, (rowIndex) {
            final number = rowIndex * 3 + 7;
            return _buildNumberButton(number.toString());
          }),
        ),
        Row(
          mainAxisAlignment: MainAxisAlignment.spaceEvenly,
          children: [
            const SizedBox(width: 64), // Empty for alignment
            _buildNumberButton('0'),
            IconButton(
              icon: const Icon(Icons.backspace),
              onPressed: () {
                if (_pinController.text.isNotEmpty) {
                  _pinController.text = _pinController.text.substring(0, _pinController.text.length - 1);
                  _pinController.selection = TextSelection.fromPosition(
                    TextPosition(offset: _pinController.text.length),
                  );
                }
              },
              iconSize: 24,
            ),
          ],
        ),
      ],
    );
  }

  Widget _buildNumberButton(String number) {
    return SizedBox(
      width: 80,
      height: 80,
      child: ElevatedButton(
        onPressed: () {
          if (_pinController.text.length < 6) {
            _pinController.text += number;
            _pinController.selection = TextSelection.fromPosition(
              TextPosition(offset: _pinController.text.length),
            );

            if (_pinController.text.length == 6) {
              _verifyPIN();
            }
          }
        },
        style: ElevatedButton.styleFrom(
          shape: const CircleBorder(),
          backgroundColor: Colors.grey[200],
          foregroundColor: Colors.black87,
          elevation: 0,
        ),
        child: Text(
          number,
          style: const TextStyle(fontSize: 24, fontWeight: FontWeight.w500),
        ),
      ),
    );
  }
}
