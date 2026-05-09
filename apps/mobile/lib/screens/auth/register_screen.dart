import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../providers/auth_provider.dart';
import '../../utils/constants.dart';

class RegisterScreen extends ConsumerStatefulWidget {
  const RegisterScreen({super.key});

  @override
  ConsumerState<RegisterScreen> createState() => _RegisterScreenState();
}

class _RegisterScreenState extends ConsumerState<RegisterScreen> {
  final _formKey = GlobalKey<FormState>();
  final _usernameController = TextEditingController();
  final _emailController = TextEditingController();
  final _phoneController = TextEditingController();
  final _passwordController = TextEditingController();
  final _confirmPasswordController = TextEditingController();
  final _pinController = TextEditingController();
  final _confirmPinController = TextEditingController();
  String? _selectedRole;
  bool _obscurePassword = true;
  bool _obscureConfirm = true;
  bool _obscurePin = true;
  bool _obscureConfirmPin = true;
  bool _acceptTerms = false;

  @override
  void dispose() {
    _usernameController.dispose();
    _emailController.dispose();
    _phoneController.dispose();
    _passwordController.dispose();
    _confirmPasswordController.dispose();
    _pinController.dispose();
    _confirmPinController.dispose();
    super.dispose();
  }

  Future<void> _handleRegister() async {
    if (!_acceptTerms) {
      _showErrorDialog('Anda harus menyetujui syarat dan ketentuan');
      return;
    }

    if (_formKey.currentState!.validate()) {
      // Validate PIN match
      if (_pinController.text != _confirmPinController.text) {
        _showErrorDialog('PIN tidak cocok');
        return;
      }

      final authState = ref.read(authProvider);
      if (authState.isLoading) return;

      await ref.read(authProvider.notifier).register(
            email: _emailController.text.trim(),
            phone: _phoneController.text.trim(),
            name: _usernameController.text.trim(),
            password: _passwordController.text,
            pin: _pinController.text,
            referralCode: null,
          );

      final updatedState = ref.read(authProvider);
      if (updatedState.errorMessage != null) {
        _showErrorDialog(updatedState.errorMessage!);
      } else if (mounted) {
        Navigator.of(context).pushReplacementNamed('/otp');
      }
    }
  }

  void _showErrorDialog(String message) {
    showDialog(
      context: context,
      builder: (ctx) => AlertDialog(
        title: const Text('Error'),
        content: Text(message),
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
    final authState = ref.watch(authProvider);

    if (authState.isAuthenticated) {
      WidgetsBinding.instance.addPostFrameCallback((_) {
        Navigator.of(context).pushReplacementNamed('/home');
      });
    }

    return Scaffold(
      appBar: AppBar(
        title: const Text('Daftar'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => Navigator.of(context).pop(),
        ),
      ),
      body: SafeArea(
        child: SingleChildScrollView(
          padding: const EdgeInsets.all(24.0),
          child: Form(
            key: _formKey,
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                Icon(
                  Icons.person_add,
                  size: 80,
                  color: Theme.of(context).colorScheme.primary,
                ),
                const SizedBox(height: 16),
                Text(
                  'Buat Akun Baru',
                  textAlign: TextAlign.center,
                  style: Theme.of(context).textTheme.headlineSmall?.copyWith(
                        fontWeight: FontWeight.bold,
                      ),
                ),
                const SizedBox(height: 32),

                // Role selection
                DropdownButtonFormField<String>(
                  value: _selectedRole,
                  decoration: const InputDecoration(
                    labelText: 'Daftar sebagai',
                    prefixIcon: Icon(Icons.badge),
                  ),
                  items: const [
                   DropdownMenuItem(
                      value: 'mitra',
                      child: Text('Mitra (Pemilik Usaha)'),
                    ),
                    DropdownMenuItem(
                      value: 'staff',
                      child: Text('Staff (Pegawai)'),
                    ),
                  ],
                  onChanged: (value) {
                    setState(() {
                      _selectedRole = value;
                    });
                  },
                  validator: (value) {
                    if (value == null || value.isEmpty) {
                      return 'Pilih peran Anda';
                    }
                    return null;
                  },
                ),
                const SizedBox(height: 16),

                // Username field
                TextFormField(
                  controller: _usernameController,
                  decoration: const InputDecoration(
                    labelText: 'Username',
                    prefixIcon: Icon(Icons.person_outline),
                    hintText: 'Pilih username unik',
                  ),
                  validator: (value) {
                    if (value == null || value.trim().isEmpty) {
                      return 'Username wajib diisi';
                    }
                    if (value.length < 4) {
                      return 'Username minimal 4 karakter';
                    }
                    return null;
                    },
                  textInputAction: TextInputAction.next,
                ),
                const SizedBox(height: 16),

                // Phone field
                TextFormField(
                  controller: _phoneController,
                  decoration: const InputDecoration(
                    labelText: 'Nomor HP',
                    prefixIcon: Icon(Icons.phone),
                    hintText: '+62xxx',
                  ),
                  keyboardType: TextInputType.phone,
                  validator: (value) {
                    if (value == null || value.trim().isEmpty) {
                      return 'Nomor HP wajib diisi';
                    }
                    if (!RegExp(r'^\+?[0-9]{10,15}$').hasMatch(value.trim())) {
                      return 'Format nomor HP tidak valid';
                    }
                    return null;
                  },
                 ),
                 // Email field
                 TextFormField(
                   controller: _emailController,
                   decoration: const InputDecoration(
                     labelText: 'Email',
                     prefixIcon: Icon(Icons.email),
                     hintText: 'example@domain.com',
                   ),
                   keyboardType: TextInputType.emailAddress,
                   validator: (value) {
                     if (value == null || value.trim().isEmpty) {
                       return 'Email wajib diisi';
                     }
                     if (!RegExp(r'^[^@]+@[^@]+\.[^@]+').hasMatch(value.trim())) {
                       return 'Format email tidak valid';
                     }
                     return null;
                   },
                 ),
                 const SizedBox(height: 16),

                 // Password field
                TextFormField(
                  controller: _passwordController,
                  decoration: InputDecoration(
                    labelText: 'Kata Sandi',
                    prefixIcon: const Icon(Icons.lock_outline),
                    suffixIcon: IconButton(
                      icon: Icon(_obscurePassword ? Icons.visibility : Icons.visibility_off),
                      onPressed: () {
                        setState(() {
                          _obscurePassword = !_obscurePassword;
                        });
                      },
                    ),
                  ),
                  obscureText: _obscurePassword,
                  validator: (value) {
                    if (value == null || value.isEmpty) {
                      return 'Kata sandi wajib diisi';
                    }
                    if (value.length < 6) {
                      return 'Kata sandi minimal 6 karakter';
                    }
                    return null;
                  },
                ),
                const SizedBox(height: 16),

                // Confirm password field
                TextFormField(
                  controller: _confirmPasswordController,
                  decoration: InputDecoration(
                    labelText: 'Konfirmasi Kata Sandi',
                    prefixIcon: const Icon(Icons.lock_reset),
                    suffixIcon: IconButton(
                      icon: Icon(_obscureConfirm ? Icons.visibility : Icons.visibility_off),
                      onPressed: () {
                        setState(() {
                          _obscureConfirm = !_obscureConfirm;
                        });
                      },
                    ),
                  ),
                  obscureText: _obscureConfirm,
                  validator: (value) {
                    if (value != _passwordController.text) {
                      return 'Kata sandi tidak cocok';
                    }
                    return null;
                  },
                 ),
                 const SizedBox(height: 16),

                 // PIN field
                 TextFormField(
                   controller: _pinController,
                   decoration: InputDecoration(
                     labelText: 'PIN Transaksi',
                     prefixIcon: const Icon(Icons.pin),
                     suffixIcon: IconButton(
                       icon: Icon(_obscurePin ? Icons.visibility : Icons.visibility_off),
                       onPressed: () {
                         setState(() {
                           _obscurePin = !_obscurePin;
                         });
                       },
                     ),
                   ),
                   obscureText: _obscurePin,
                   keyboardType: TextInputType.number,
                   maxLength: 6,
                   validator: (value) {
                     if (value == null || value.isEmpty) {
                       return 'PIN wajib diisi';
                     }
                     if (value.length < 6) {
                       return 'PIN minimal 6 digit';
                     }
                     return null;
                   },
                 ),
                 const SizedBox(height: 16),

                 // Confirm PIN field
                 TextFormField(
                   controller: _confirmPinController,
                   decoration: InputDecoration(
                     labelText: 'Konfirmasi PIN',
                     prefixIcon: const Icon(Icons.pin),
                     suffixIcon: IconButton(
                       icon: Icon(_obscureConfirmPin ? Icons.visibility : Icons.visibility_off),
                       onPressed: () {
                         setState(() {
                           _obscureConfirmPin = !_obscureConfirmPin;
                         });
                       },
                     ),
                   ),
                   obscureText: _obscureConfirmPin,
                   keyboardType: TextInputType.number,
                   maxLength: 6,
                   validator: (value) {
                     if (value != _pinController.text) {
                       return 'PIN tidak cocok';
                     }
                     return null;
                   },
                 ),
                 const SizedBox(height: 16),

                 // Terms checkbox
                Row(
                  children: [
                    Checkbox(
                      value: _acceptTerms,
                      onChanged: (value) {
                        setState(() {
                          _acceptTerms = value ?? false;
                        });
                      },
                    ),
                    Expanded(
                      child: RichText(
                        text: TextSpan(
                          style: DefaultTextStyle.of(context).style,
                          children: [
                            const TextSpan(text: 'Saya menyetujui '),
                            TextSpan(
                              text: 'Syarat & Ketentuan',
                              style: TextStyle(
                                color: Theme.of(context).colorScheme.primary,
                                fontWeight: FontWeight.w600,
                              ),
                            ),
                            const TextSpan(text: ' dan '),
                            TextSpan(
                              text: 'Kebijakan Privasi',
                              style: TextStyle(
                                color: Theme.of(context).colorScheme.primary,
                                fontWeight: FontWeight.w600,
                              ),
                            ),
                          ],
                        ),
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 24),

                // Register button
                ElevatedButton(
                  onPressed: authState.isLoading ? null : _handleRegister,
                  child: authState.isLoading
                      ? const SizedBox(
                          height: 20,
                          width: 20,
                          child: CircularProgressIndicator(
                            strokeWidth: 2,
                            valueColor: AlwaysStoppedAnimation<Color>(Colors.white)),
                        )
                      : const Text('DAFTAR'),
                ),
                const SizedBox(height: 16),

                // Already have account
                TextButton(
                  onPressed: () {
                    Navigator.of(context).pop();
                  },
                  child: const Text('Sudah punya akun? Masuk'),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
