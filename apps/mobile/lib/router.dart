import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../providers/auth_provider.dart';
import '../models/product.dart';
import '../models/staff.dart';
import 'screens/auth/login_screen.dart';
import 'screens/auth/register_screen.dart';
import 'screens/auth/otp_screen.dart';
import 'screens/home/home_screen.dart';
import 'screens/transaction/transaction_history_screen.dart';
import 'screens/transaction/transaction_initiate_screen.dart';
import 'screens/transaction/confirmation_screen.dart';
import 'screens/transaction/pin_entry_screen.dart';
import 'screens/transaction/transaction_detail_screen.dart';
import 'screens/wallet/wallet_screen.dart';
import 'screens/staff/staff_list_screen.dart';
import 'screens/staff/staff_add_edit_screen.dart';
import 'screens/staff/staff_topup_screen.dart';
import 'screens/settings/settings_screen.dart';

final routerProvider = Provider<GoRouter>((ref) {
  final authState = ref.watch(authProvider);

  return GoRouter(
    initialLocation: '/login',
    redirect: (context, state) {
      final isAuthenticated = authState.isAuthenticated;
      final isLoginRoute = state.matchedLocation == '/login';
      final isRegisterRoute = state.matchedLocation == '/register';
      final isOtpRoute = state.matchedLocation == '/otp';

      if (!isAuthenticated) {
        if (isLoginRoute || isRegisterRoute || isOtpRoute) {
          return null;
        }
        return '/login';
      }

      if (isLoginRoute || isRegisterRoute || isOtpRoute) {
        return '/home';
      }

      return null;
    },
    routes: [
      GoRoute(
        path: '/login',
        builder: (context, state) => const LoginScreen(),
      ),
      GoRoute(
        path: '/register',
        builder: (context, state) => const RegisterScreen(),
      ),
      GoRoute(
        path: '/otp',
        builder: (context, state) => const OTPScreen(),
      ),
      GoRoute(
        path: '/home',
        builder: (context, state) => const HomeScreen(),
      ),
      GoRoute(
        path: '/transactions',
        builder: (context, state) => const TransactionHistoryScreen(),
      ),
      GoRoute(
        path: '/transaction/initiate',
        builder: (context, state) => TransactionInitiateScreen(
          product: state.extra as Product,
        ),
      ),
      GoRoute(
        path: '/transaction/confirm',
        builder: (context, state) => const ConfirmationScreen(),
      ),
      GoRoute(
        path: '/transaction/pin',
        builder: (context, state) => const PinEntryScreen(),
      ),
      GoRoute(
        path: '/transactions/detail',
        builder: (context, state) => const TransactionDetailScreen(),
      ),
      GoRoute(
        path: '/wallet',
        builder: (context, state) => const WalletScreen(),
      ),
      GoRoute(
        path: '/staff',
        builder: (context, state) => const StaffListScreen(),
      ),
      GoRoute(
        path: '/staff/add',
        builder: (context, state) => const StaffAddEditScreen(),
      ),
      GoRoute(
        path: '/staff/edit',
        builder: (context, state) => StaffAddEditScreen(
          staff: state.extra as Staff,
        ),
      ),
      GoRoute(
        path: '/staff/topup',
        builder: (context, state) => const StaffTopUpScreen(),
      ),
      GoRoute(
        path: '/settings',
        builder: (context, state) => const SettingsScreen(),
      ),
    ],
  );
});

class AuthFeel {
  final bool isAuthenticated;
  final String? role;

  AuthFeel({required this.isAuthenticated, this.role});
}

final authFeelProvider = Provider<AuthFeel>((ref) {
  final authState = ref.watch(authProvider);
  return AuthFeel(
    isAuthenticated: authState.isAuthenticated,
    role: authState.role,
  );
});