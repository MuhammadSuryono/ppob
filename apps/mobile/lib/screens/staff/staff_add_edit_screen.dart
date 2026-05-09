import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../models/staff.dart';
import '../../repositories/staff_repository.dart';
import '../../providers/auth_provider.dart';
import '../../utils/constants.dart';

class StaffAddEditScreen extends ConsumerStatefulWidget {
  final Staff? staff; // If null, it's add mode

  const StaffAddEditScreen({super.key, this.staff});

  @override
  ConsumerState<StaffAddEditScreen> createState() => _StaffAddEditScreenState();
}

class _StaffAddEditScreenState extends ConsumerState<StaffAddEditScreen> {
  final _formKey = GlobalKey<FormState>();
  final _nameController = TextEditingController();
  final _phoneController = TextEditingController();
  final _pinController = TextEditingController();

  String? _marginScheme;
  double? _fixedAllowance;
  double? _marginShare;
  double _dailyLimitAmount = 1000000;
  int _dailyLimitCount = 50;

  bool get _isEdit => widget.staff != null;

  @override
  void initState() {
    super.initState();
    if (_isEdit) {
      _nameController.text = widget.staff!.name;
      _phoneController.text = widget.staff!.phoneNumber;
      _marginScheme = widget.staff!.marginScheme;
      _fixedAllowance = widget.staff!.fixedAllowanceAmount;
      _marginShare = widget.staff!.marginSharePercentage;
      _dailyLimitAmount = widget.staff!.dailyLimitAmount;
      _dailyLimitCount = widget.staff!.dailyLimitCount;
    }
  }

  @override
  void dispose() {
    _nameController.dispose();
    _phoneController.dispose();
    _pinController.dispose();
    super.dispose();
  }

  Future<void> _saveStaff() async {
    if (!_formKey.currentState!.validate()) return;

    final authState = ref.read(authProvider);
    final staffRepo = ref.read(staffRepositoryProvider);

    try {
      if (_isEdit) {
        final updatedStaff = widget.staff!.copyWith(
          name: _nameController.text.trim(),
          phoneNumber: _phoneController.text.trim(),
          marginScheme: _marginScheme!,
          fixedAllowanceAmount: _fixedAllowance,
          marginSharePercentage: _marginShare,
          dailyLimitAmount: _dailyLimitAmount,
          dailyLimitCount: _dailyLimitCount,
          updatedAt: DateTime.now(),
        );
        await staffRepo.updateStaff(updatedStaff);
      } else {
        final newStaff = Staff(
          id: 'staff_${DateTime.now().millisecondsSinceEpoch}',
          mitraId: authState.userId!,
          name: _nameController.text.trim(),
          phoneNumber: _phoneController.text.trim(),
          pinHash: _pinController.text.isNotEmpty ? _pinController.text : null,
          marginScheme: _marginScheme!,
          fixedAllowanceAmount: _fixedAllowance,
          marginSharePercentage: _marginShare,
          dailyLimitAmount: _dailyLimitAmount,
          dailyLimitCount: _dailyLimitCount,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        );
        await staffRepo.createStaff(newStaff);
      }

      if (mounted) {
        Navigator.of(context).pop(true);
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
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(_isEdit ? 'Edit Staff' : 'Tambah Staff'),
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
                // Basic info
                const Text(
                  'Informasi Staff',
                  style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold),
                ),
                const SizedBox(height: 16),

                TextFormField(
                  controller: _nameController,
                  decoration: const InputDecoration(
                    labelText: 'Nama Lengkap',
                    prefixIcon: Icon(Icons.person),
                  ),
                  validator: (value) {
                    if (value == null || value.trim().isEmpty) {
                      return 'Nama wajib diisi';
                    }
                    return null;
                  },
                ),
                const SizedBox(height: 16),

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
                      return 'Format nomor tidak valid';
                    }
                    return null;
                  },
                ),
                if (!_isEdit) ...[
                  const SizedBox(height: 16),
                  TextFormField(
                    controller: _pinController,
                    decoration: const InputDecoration(
                      labelText: 'PIN Transaksi (Opsional)',
                      prefixIcon: Icon(Icons.lock),
                      hintText: '6 digit',
                    ),
                    keyboardType: TextInputType.number,
                    maxLength: 6,
                    obscureText: true,
                  ),
                ],

                const SizedBox(height: 24),

                // Margin scheme
                const Text(
                  'Skema Margin',
                  style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold),
                ),
                const SizedBox(height: 12),

                RadioListTile<String>(
                  title: const Text('Gaji Tetap (Fixed Allowance)'),
                  subtitle: const Text('Staff menerima Rp X per hari bulanan'),
                  value: 'FixedAllowance',
                  groupValue: _marginScheme,
                  onChanged: (value) {
                    setState(() {
                      _marginScheme = value;
                    });
                  },
                ),
                RadioListTile<String>(
                  title: const Text('Bagi Hasil (Margin Share)'),
                  subtitle: const Text('Staff mendapat % dari margin transaksi'),
                  value: 'MarginShare',
                  groupValue: _marginScheme,
                  onChanged: (value) {
                    setState(() {
                      _marginScheme = value;
                    });
                  },
                ),

                if (_marginScheme == 'FixedAllowance') ...[
                  const SizedBox(height: 16),
                  TextFormField(
                    decoration: const InputDecoration(
                      labelText: 'Jumlah Gaji (Rp)',
                      prefixIcon: Icon(Icons.attach_money),
                    ),
                    keyboardType: TextInputType.number,
                    initialValue: _fixedAllowance?.toString() ?? '',
                    validator: (value) {
                      if (_marginScheme == 'FixedAllowance') {
                        if (value == null || value.isEmpty) {
                          return 'Jumlah gaji wajib diisi';
                        }
                        final amount = double.tryParse(value);
                        if (amount == null || amount <= 0) {
                          return 'Jumlah tidak valid';
                        }
                      }
                      return null;
                    },
                    onChanged: (value) {
                      _fixedAllowance = double.tryParse(value);
                    },
                  ),
                ],

                if (_marginScheme == 'MarginShare') ...[
                  const SizedBox(height: 16),
                  TextFormField(
                    decoration: const InputDecoration(
                      labelText: 'Persentase Bagi Hasil (%)',
                      prefixIcon: Icon(Icons.percent),
                      hintText: '70',
                    ),
                    keyboardType: TextInputType.number,
                    initialValue: _marginShare != null ? '${(_marginShare! * 100).toInt()}' : '',
                    validator: (value) {
                      if (_marginScheme == 'MarginShare') {
                        if (value == null || value.isEmpty) {
                          return 'Persentase wajib diisi';
                        }
                        final percent = double.tryParse(value);
                        if (percent == null || percent <= 0 || percent > 100) {
                          return 'Persentase 0-100';
                        }
                      }
                      return null;
                    },
                    onChanged: (value) {
                      _marginShare = double.tryParse(value) != null ? double.tryParse(value)! / 100 : null;
                    },
                  ),
                ],

                const SizedBox(height: 24),

                // Daily limits
                const Text(
                  'Limit Harian',
                  style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold),
                ),
                const SizedBox(height: 12),
                TextFormField(
                  decoration: const InputDecoration(
                    labelText: 'Limit Amount (Rp)',
                    prefixIcon: Icon(Icons.money),
                  ),
                  keyboardType: TextInputType.number,
                  initialValue: _dailyLimitAmount.toStringAsFixed(0),
                  validator: (value) {
                    if (value == null || value.isEmpty) {
                      return 'Limit amount wajib diisi';
                    }
                    final amount = double.tryParse(value);
                    if (amount == null || amount <= 0) {
                      return 'Jumlah tidak valid';
                    }
                    return null;
                  },
                  onChanged: (value) {
                    _dailyLimitAmount = double.tryParse(value) ?? 0;
                  },
                ),
                const SizedBox(height: 16),
                TextFormField(
                  decoration: const InputDecoration(
                    labelText: 'Limit Jumlah Transaksi',
                    prefixIcon: Icon(Icons.countertops),
                  ),
                  keyboardType: TextInputType.number,
                  initialValue: _dailyLimitCount.toString(),
                  validator: (value) {
                    if (value == null || value.isEmpty) {
                      return 'Limit count wajib diisi';
                    }
                    final count = int.tryParse(value);
                    if (count == null || count <= 0) {
                      return 'Jumlah tidak valid';
                    }
                    return null;
                  },
                  onChanged: (value) {
                    _dailyLimitCount = int.tryParse(value) ?? 0;
                  },
                ),

                const SizedBox(height: 32),

                // Save button
                SizedBox(
                  width: double.infinity,
                  child: ElevatedButton(
                    onPressed: _marginScheme == null ? null : _saveStaff,
                    child: Text(_isEdit ? 'UPDATE STAFF' : 'TAMBAH STAFF'),
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
