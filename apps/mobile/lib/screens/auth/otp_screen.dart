import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../providers/auth_provider.dart';
import '../../utils/constants.dart';
import 'login_screen.dart';

class OTPScreen extends ConsumerStatefulWidget {
  const OTPScreen({super.key});

  @override
  ConsumerState<OTPScreen> createState() => _OTPScreenState();
}

class _OTPScreenState extends ConsumerState<OTPScreen> {
  final _formKey = GlobalKey<FormState>();
  final _otpController = TextEditingController();
  final _focusNode = FocusNode();
  int _resendCooldown = 0;

  @override
  void initState() {
    super.initState();
    _startResendCooldown();
    // Focus on OTP field
    WidgetsBinding.instance.addPostFrameCallback((_) {
      FocusScope.of(context).requestFocus(_focusNode);
    });
  }

  @override
  void dispose() {
    _otpController.dispose();
    _focusNode.dispose();
    super.dispose();
  }

  void _startResendCooldown() {
    setState(() {
      _resendCooldown = 60;
    });
    Future.doWhile(() async {
      await Future.delayed(const Duration(seconds: 1));
      if (mounted) {
        setState(() {
          if (_resendCooldown > 0) _resendCooldown--;
        });
        return _resendCooldown > 0;
      }
      return false;
    });
  }

  Future<void> _handleVerifyOTP() async {
    if (_formKey.currentState!.validate()) {
      final authState = ref.read(authProvider);
      if (authState.isLoading) return;

      // Note: In real app, this would call verifyOtp API
      // For now, just simulate success
      await Future.delayed(const Duration(seconds: 1));

      if (mounted) {
        // Navigate to home on success
        context.pushReplacement('/home');
      }
    }
  }

  Future<void> _handleResendOTP() async {
    if (_resendCooldown > 0) return;

    // Show confirmation dialog
    final confirm = await showDialog<bool>(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('Kirim Ulang OTP'),
        content: const Text('Kirim kode OTP baru ke nomor HP Anda?'),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(ctx).pop(false),
            child: const Text('Batal'),
          ),
          TextButton(
            onPressed: () => Navigator.of(ctx).pop(true),
            child: const Text('Kirim'),
          ),
        ],
      ),
    );

    if (confirm == true) {
      // Simulate resend
      await Future.delayed(const Duration(milliseconds: 500));
      _startResendCooldown();
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Kode OTP telah dikirim ulang'),
            backgroundColor: Colors.green,
          ),
        );
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final authState = ref.watch(authProvider);

    if (authState.isAuthenticated) {
      WidgetsBinding.instance.addPostFrameCallback((_) {
        context.pushReplacement('/home');
      });
    }

    // Get phone number from auth state (mock)
    final phoneNumber = authState.username ?? '+62xxx-xxxx-xxxx';

    return Scaffold(
      appBar: AppBar(
        title: const Text('Verifikasi OTP'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => Navigator.of(context).pop(),
        ),
      ),
      body: SafeArea(
        child: Center(
          child: SingleChildScrollView(
            padding: const EdgeInsets.all(24.0),
            child: Form(
              key: _formKey,
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  // Icon
                  Icon(
                    Icons.sms,
                    size: 80,
                    color: Theme.of(context).colorScheme.primary,
                  ),
                  const SizedBox(height: 24),

                  // Title
                  Text(
                    'Verifikasi Nomor HP',
                    style: Theme.of(context).textTheme.titleLarge?.copyWith(
                          fontWeight: FontWeight.bold,
                        ),
                  ),
                  const SizedBox(height: 8),

                  // Subtitle
                  Text(
                    'Kami telah mengirim kode OTP ke',
                    style: Theme.of(context).textTheme.bodyMedium,
                  ),
                  const SizedBox(height: 4),
                  Text(
                    phoneNumber,
                    style: Theme.of(context).textTheme.bodyMedium?.copyWith(
                          fontWeight: FontWeight.bold,
                          color: Theme.of(context).colorScheme.primary,
                        ),
                  ),
                  const SizedBox(height: 32),

                  // OTP input
                  TextFormField(
                    controller: _otpController,
                    focusNode: _focusNode,
                    decoration: const InputDecoration(
                      labelText: 'Kode OTP',
                      hintText: 'Masukkan 6 digit kode',
                      prefixIcon: Icon(Icons.password),
                    ),
                    keyboardType: TextInputType.number,
                    maxLength: 6,
                    textAlign: TextAlign.center,
                    style: const TextStyle(
                      fontSize: 24,
                      letterSpacing: 8,
                      fontWeight: FontWeight.bold,
                    ),
                    validator: (value) {
                      if (value == null || value.isEmpty) {
                        return 'Kode OTP wajib diisi';
                      }
                      if (value.length != 6) {
                        return 'Kode OTP harus 6 digit';
                      }
                      return null;
                    },
                  ),
                  const SizedBox(height: 24),

                  // Verify button
                  SizedBox(
                    width: double.infinity,
                    child: ElevatedButton(
                      onPressed: authState.isLoading ? null : _handleVerifyOTP,
                      child: authState.isLoading
                          ? const SizedBox(
                              height: 20,
                              width: 20,
                              child: CircularProgressIndicator(
                                strokeWidth: 2,
                                valueColor: AlwaysStoppedAnimation<Color>(Colors.white)),
                            )
                          : const Text('VERIFIKASI'),
                    ),
                  ),
                  const SizedBox(height: 16),

                  // Resend OTP
                  TextButton(
                    onPressed: _resendCooldown > 0 ? null : _handleResendOTP,
                    child: Text(
                      _resendCooldown > 0
                          ? 'Kirim ulang dalam $_resendCooldown detik'
                          : 'Kirim ulang kode OTP',
                    ),
                  ),

                  // Change phone number
                  TextButton(
                    onPressed: () {
                      Navigator.of(context).pop();
                      context.pushReplacement('/login');
                    },
                    child: const Text('Ubah nomor HP'),
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}
