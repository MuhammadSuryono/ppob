import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../providers/auth_provider.dart';
import '../../repositories/product_repository.dart';
import '../../models/product.dart';
import '../../utils/constants.dart';
import '../../widgets/product/product_card.dart';
import 'product_detail_screen.dart';
import '../transaction/transaction_initiate_screen.dart';

class HomeScreen extends ConsumerStatefulWidget {
  const HomeScreen({super.key});

  @override
  ConsumerState<HomeScreen> createState() => _HomeScreenState();
}

class _HomeScreenState extends ConsumerState<HomeScreen> {
  final _searchController = TextEditingController();
  String? _selectedCategory;
  bool _isSearching = false;

  final List<String> _categories = [
    'pulsa',
    'token-listrik',
    'paket-data',
    'multi',
  ];

  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final authState = ref.watch(authProvider);
    final productRepository = ref.watch(productRepositoryProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text('PPOB'),
        actions: [
          IconButton(
            icon: const Icon(Icons.history),
            onPressed: () {
              context.push('/transactions');
            },
          ),
          IconButton(
            icon: const Icon(Icons.account_balance_wallet),
            onPressed: () {
              context.push('/wallet');
            },
          ),
          IconButton(
            icon: const Icon(Icons.settings),
            onPressed: () {
              context.push('/settings');
            },
          ),
        ],
      ),
      drawer: _buildDrawer(authState),
      body: Column(
        children: [
          // Search bar
          Padding(
            padding: const EdgeInsets.all(16.0),
            child: TextField(
              controller: _searchController,
              decoration: InputDecoration(
                hintText: 'Cari produk Pulsa, Token, Paket Data...',
                prefixIcon: const Icon(Icons.search),
                suffixIcon: _isSearching
                    ? IconButton(
                        icon: const Icon(Icons.clear),
                        onPressed: () {
                          setState(() {
                            _searchController.clear();
                            _isSearching = false;
                          });
                        },
                      )
                    : null,
                filled: true,
                fillColor: Colors.white,
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(30),
                  borderSide: BorderSide.none,
                ),
                contentPadding: const EdgeInsets.symmetric(horizontal: 20, vertical: 12),
              ),
              onChanged: (value) {
                setState(() {
                  _isSearching = value.isNotEmpty;
                });
              },
              onSubmitted: (value) {
                setState(() {});
              },
            ),
          ),

          // Category chips
          SizedBox(
            height: 50,
            child: ListView(
              scrollDirection: Axis.horizontal,
              padding: const EdgeInsets.symmetric(horizontal: 16),
              children: [
                _buildCategoryChip('Semua', null),
                ..._categories.map((cat) => _buildCategoryChip(
                      _formatCategoryName(cat),
                      cat,
                    )),
              ],
            ),
          ),

          const SizedBox(height: 16),

          // Product grid
          Expanded(
            child: FutureBuilder<List<Product>>(
              future: productRepository.getProducts(
                category: _selectedCategory,
                search: _isSearching ? _searchController.text : null,
              ),
              builder: (context, snapshot) {
                if (snapshot.connectionState == ConnectionState.waiting) {
                  return const Center(child: CircularProgressIndicator());
                }

                if (snapshot.hasError) {
                  return _buildErrorView(snapshot.error.toString());
                }

                final products = snapshot.data ?? [];

                if (products.isEmpty) {
                  return _buildEmptyView();
                }

                return GridView.builder(
                  padding: const EdgeInsets.all(16),
                  gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
                    crossAxisCount: 2,
                    crossAxisSpacing: 16,
                    mainAxisSpacing: 16,
                    childAspectRatio: 0.75,
                  ),
                  itemCount: products.length,
                   itemBuilder: (context, index) {
                     final product = products[index];
                     return ProductCard(
                       product: product,
                       onTap: () {
                         context.push('/transaction/initiate', extra: product);
                       },
                     );
                   },
                );
              },
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildCategoryChip(String label, String? value) {
    final isSelected = _selectedCategory == value;
    return Padding(
      padding: const EdgeInsets.only(right: 8.0),
      child: FilterChip(
        label: Text(label),
        selected: isSelected,
        onSelected: (selected) {
          setState(() {
            _selectedCategory = selected ? value : null;
          });
        },
        selectedColor: Theme.of(context).colorScheme.primaryContainer,
        checkmarkColor: Theme.of(context).colorScheme.primary,
      ),
    );
  }

  String _formatCategoryName(String category) {
    switch (category) {
      case 'pulsa':
        return 'Pulsa';
      case 'token-listrik':
        return 'Token Listrik';
      case 'paket-data':
        return 'Paket Data';
      case 'multi':
        return 'Multi';
      default:
        return category;
    }
  }

  Widget _buildErrorView(String error) {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          const Icon(Icons.error_outline, size: 64, color: Colors.red),
          const SizedBox(height: 16),
          Text('Terjadi kesalahan: $error'),
          const SizedBox(height: 16),
          ElevatedButton(
            onPressed: () => setState(() {}),
            child: const Text('Coba Lagi'),
          ),
        ],
      ),
    );
  }

  Widget _buildEmptyView() {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: const [
          Icon(Icons.search_off, size: 64, color: Colors.grey),
          SizedBox(height: 16),
          Text(
            'Tidak ada produk ditemukan',
            style: TextStyle(fontSize: 16, color: Colors.grey),
          ),
        ],
      ),
    );
  }

  Widget _buildDrawer(AuthState authState) {
    return Drawer(
      child: ListView(
        padding: EdgeInsets.zero,
        children: [
          UserAccountsDrawerHeader(
            accountName: Text(authState.username ?? 'User'),
            accountEmail: Text(authState.role == 'mitra' ? 'Mitra' : 'Staff'),
            currentAccountPicture: CircleAvatar(
              backgroundColor: Colors.white,
              child: Icon(
                authState.role == 'mitra' ? Icons.store : Icons.person,
                size: 40,
                color: Theme.of(context).colorScheme.primary,
              ),
            ),
            decoration: BoxDecoration(
              color: Theme.of(context).colorScheme.primary,
            ),
          ),
          ListTile(
            leading: const Icon(Icons.home),
            title: const Text('Home'),
            onTap: () {
              Navigator.of(context).pop();
            },
          ),
           if (authState.role == 'mitra')
             ListTile(
               leading: const Icon(Icons.people),
               title: const Text('Kelola Staff'),
               onTap: () {
                 Navigator.of(context).pop();
                 context.push('/staff');
               },
             ),
           ListTile(
             leading: const Icon(Icons.account_balance_wallet),
             title: const Text('Wallet'),
             onTap: () {
               Navigator.of(context).pop();
               context.push('/wallet');
             },
           ),
           ListTile(
             leading: const Icon(Icons.history),
             title: const Text('Riwayat Transaksi'),
             onTap: () {
               Navigator.of(context).pop();
               context.push('/transactions');
             },
           ),
           ListTile(
             leading: const Icon(Icons.settings),
             title: const Text('Pengaturan'),
             onTap: () {
               Navigator.of(context).pop();
               context.push('/settings');
             },
           ),
          const Divider(),
          ListTile(
            leading: const Icon(Icons.logout, color: Colors.red),
            title: const Text('Logout', style: TextStyle(color: Colors.red)),
            onTap: () async {
              await ref.read(authProvider.notifier).logout();
              if (mounted) {
                Navigator.of(context).pop();
              }
            },
          ),
        ],
      ),
    );
  }
}
