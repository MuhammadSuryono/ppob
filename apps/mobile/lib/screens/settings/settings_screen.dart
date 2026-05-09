import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import '../../providers/auth_provider.dart';
import '../../services/biometric_service.dart';
import '../../utils/constants.dart';

class SettingsScreen extends ConsumerStatefulWidget {
  const SettingsScreen({super.key});

  @override
  ConsumerState<SettingsScreen> createState() => _SettingsScreenState();
}

class _SettingsScreenState extends ConsumerState<SettingsScreen> {
  final _storage = const FlutterSecureStorage();
  final BiometricService _biometricService = BiometricService();

  bool _biometricEnabled = false;
  bool _isLoading = false;

  @override
  void initState() {
    super.initState();
    _loadBiometricSetting();
  }

  Future<void> _loadBiometricSetting() async {
    final value = await _storage.read(key: AppConstants.keyBiometricEnabled);
    setState(() {
      _biometricEnabled = value == 'true';
    });
  }

  @override
  Widget build(BuildContext context) {
    final authState = ref.watch(authProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Pengaturan'),
        leading: IconButton(
          icon: const Icon(Icons.arrow_back),
          onPressed: () => Navigator.of(context).pop(),
        ),
      ),
      body: ListView(
        children: [
          // Account section
          _buildSectionHeader('Akun'),
          ListTile(
            leading: const Icon(Icons.person),
            title: Text(authState.username ?? 'User'),
            subtitle: Text(_formatRole(authState.role)),
          ),
          ListTile(
            leading: const Icon(Icons.phone),
            title: const Text('Nomor HP'),
            subtitle: const Text('+62 812-3456-7890'), // Mock
          ),

          const Divider(),

          // Security section
          _buildSectionHeader('Keamanan'),
          SwitchListTile(
            secondary: const Icon(Icons.fingerprint),
            title: const Text('Autentikasi Biometrik'),
            subtitle: const Text('Gunakan fingerprint/Face ID untuk login'),
            value: _biometricEnabled,
            onChanged: _isLoading ? null : (value) async {
              setState(() => _isLoading = true);
              try {
                final isAvailable = await _biometricService.isAvailable();
                if (!isAvailable) {
                  if (mounted) {
                    ScaffoldMessenger.of(context).showSnackBar(
                      const SnackBar(
                        content: Text('Biometric tidak tersedia di perangkat ini'),
                      ),
                    );
                  }
                  return;
                }

                // Authenticate first to confirm
                final didAuth = await _biometricService.authenticate();
                if (didAuth) {
                  await _storage.write(
                    key: AppConstants.keyBiometricEnabled,
                    value: value.toString(),
                  );
                  setState(() {
                    _biometricEnabled = value;
                  });
                  if (mounted) {
                    ScaffoldMessenger.of(context).showSnackBar(
                      SnackBar(
                        content: Text(value ? 'Biometric diaktifkan' : 'Biometric dimatikan'),
                        backgroundColor: Colors.green,
                      ),
                    );
                  }
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
            },
          ),

          const Divider(),

          // Other settings
          _buildSectionHeader('Umum'),
          ListTile(
            leading: const Icon(Icons.notifications),
            title: const Text('Notifikasi'),
            trailing: Switch(
              value: true,
              onChanged: (value) {
                // TODO: Implement notification settings
              },
            ),
          ),
          ListTile(
            leading: const Icon(Icons.language),
            title: const Text('Bahasa'),
            trailing: const Text('Indonesia', style: TextStyle(color: Colors.grey)),
            onTap: () {
              // TODO: Language selection
            },
          ),

          const Divider(),

          // About
          _buildSectionHeader('Tentang'),
          const ListTile(
            leading: Icon(Icons.info),
            title: Text('Versi'),
            subtitle: Text('1.0.0'),
          ),
          ListTile(
            leading: const Icon(Icons.privacy_tip),
            title: const Text('Kebijakan Privasi'),
            onTap: () {
              // TODO: Show privacy policy
            },
          ),
          ListTile(
            leading: const Icon(Icons.description),
            title: const Text('Syarat & Ketentuan'),
            onTap: () {
              // TODO: Show terms
            },
          ),

          const SizedBox(height: 32),

          // Logout button
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 16),
            child: SizedBox(
              width: double.infinity,
              child: OutlinedButton.icon(
                onPressed: _isLoading
                    ? null
                    : () async {
                        final confirm = await showDialog<bool>(
                          context: context,
                          builder: (ctx) => AlertDialog(
                            title: const Text('Logout'),
                            content: const Text('Yakin ingin logout?'),
                            actions: [
                              TextButton(
                                onPressed: () => Navigator.of(ctx).pop(false),
                                child: const Text('Batal'),
                              ),
                              TextButton(
                                onPressed: () => Navigator.of(ctx).pop(true),
                                child: const Text('Logout', style: TextStyle(color: Colors.red)),
                              ),
                            ],
                          ),
                        );

                        if (confirm == true) {
                          await ref.read(authProvider.notifier).logout();
                          if (mounted) {
                            Navigator.of(context).pushReplacementNamed('/login');
                          }
                        }
                      },
                icon: const Icon(Icons.logout, color: Colors.red),
                label: const Text('Logout', style: TextStyle(color: Colors.red)),
                style: OutlinedButton.styleFrom(
                  side: const BorderSide(color: Colors.red),
                ),
              ),
            ),
          ),

          const SizedBox(height: 32),
        ],
      ),
    );
  }

  Widget _buildSectionHeader(String title) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(16, 16, 16, 8),
      child: Text(
        title,
        style: TextStyle(
          fontSize: 12,
          fontWeight: FontWeight.bold,
          color: Colors.grey[600],
          letterSpacing: 1.0,
        ),
      ),
    );
  }

  String _formatRole(String? role) {
    switch (role) {
      case 'mitra':
        return 'Mitra';
      case 'staff':
        return 'Staff';
      default:
        return 'User';
    }
  }
}
