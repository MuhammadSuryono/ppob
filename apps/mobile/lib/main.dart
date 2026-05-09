import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:hive_flutter/hive_flutter.dart';
import 'package:firebase_core/firebase_core.dart';
import 'firebase_options.dart';
import 'models/user.dart';
import 'models/product.dart';
import 'models/transaction.dart';
import 'models/staff.dart';
import 'models/wallet.dart';
import 'models/pending_sync_item.dart';
import 'providers/app_provider.dart';
import 'services/offline_sync_service.dart';
import 'services/notification_service.dart';
import 'utils/constants.dart';
import 'router.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();

  // Initialize Firebase
  await Firebase.initializeApp(
    options: DefaultFirebaseOptions.currentPlatform,
  );

  // Initialize Hive
  await Hive.initFlutter();

  // Register Hive adapters
  Hive.registerAdapter(UserAdapter());
  Hive.registerAdapter(ProductAdapter());
  Hive.registerAdapter(TransactionAdapter());
  Hive.registerAdapter(StaffAdapter());
  Hive.registerAdapter(WalletAdapter());
  Hive.registerAdapter(PendingSyncItemAdapter());

  // Open Hive boxes
  await Hive.openBox<User>(AppConstants.boxUser);
  await Hive.openBox<Product>(AppConstants.boxProducts);
  await Hive.openBox<Transaction>(AppConstants.boxTransactions);
  await Hive.openBox<Staff>(AppConstants.boxStaff);
  await Hive.openBox<Wallet>(AppConstants.boxWallet);
  await Hive.openBox<PendingSyncItem>(AppConstants.boxPendingSync);
  await Hive.openBox<String>(AppConstants.boxSettings);

  // Initialize offline sync service
  await OfflineSyncService().initialize();

  // Initialize notifications
  await NotificationService().initialize();

  runApp(const ProviderScope(child: MyApp()));
}

class MyApp extends ConsumerWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final appState = ref.watch(appProvider);
    final router = ref.watch(routerProvider);

    return MaterialApp.router(
      title: AppConstants.appName,
      theme: ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: const Color(0xFF4CAF50)),
        scaffoldBackgroundColor: const Color(0xFFF5F5F5),
        appBarTheme: const AppBarTheme(
          backgroundColor: Color(0xFF4CAF50),
          foregroundColor: Colors.white,
          elevation: 0,
        ),
        elevatedButtonTheme: ElevatedButtonThemeData(
          style: ElevatedButton.styleFrom(
            padding: const EdgeInsets.symmetric(vertical: 16),
            shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(8)),
            textStyle: const TextStyle(fontSize: 16, fontWeight: FontWeight.w600),
          ),
        ),
        inputDecorationTheme: const InputDecorationTheme(
          filled: true,
          fillColor: Colors.white,
          border: OutlineInputBorder(borderRadius: BorderRadius.all(Radius.circular(8))),
          contentPadding: EdgeInsets.symmetric(horizontal: 16, vertical: 12),
        ),
      ),
      routerConfig: router,
      debugShowCheckedModeBanner: false,
    );
  }
}