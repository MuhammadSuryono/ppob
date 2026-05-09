class AppConstants {
  static const String appName = 'PPOB';
  static const String apiBaseUrl = 'https://fedora.sinauplatform.id';
  static const String apiVersion = 'v1';

  // Service base URLs (path-based routing)
  static const String baseAuth = 'https://fedora.sinauplatform.id/api/v1/auth';
  static const String baseUsers = 'https://fedora.sinauplatform.id/api/v1/user';
  static const String baseProducts = 'https://fedora.sinauplatform.id/api/v1/product';
  static const String baseWallets = 'https://fedora.sinauplatform.id/api/v1/wallet';
  static const String baseTransactions = 'https://fedora.sinauplatform.id/api/v1/transaction';
  static const String baseIntegrations = 'https://fedora.sinauplatform.id/api/v1/integration';

  // API endpoints
  static const String endpointAuth = baseAuth;
  static const String endpointUsers = baseUsers;
  static const String endpointWallets = baseWallets;
  static const String endpointProducts = baseProducts;
  static const String endpointTransactions = baseTransactions;
  static const String endpointStaff = baseUsers + 's';
  static const String endpointIntegration = baseIntegrations;

  // Hive box names
  static const String boxUser = 'user_box';
  static const String boxWallet = 'wallet_box';
  static const String boxTransactions = 'transactions_box';
  static const String boxProducts = 'products_box';
  static const String boxStaff = 'staff_box';
  static const String boxPendingSync = 'pending_sync_queue';
  static const String boxSettings = 'settings_box';

  // Auth keys
  static const String keyAuthToken = 'auth_token';
  static const String keyUserData = 'user_data';
  static const String keyPIN = 'user_pin';
  static const String keyBiometricEnabled = 'biometric_enabled';

  // Transaction status
  static const String statusInitiated = 'Initiated';
  static const String statusPending = 'Pending';
  static const String statusSuccess = 'Success';
  static const String statusFailed = 'Failed';
  static const String statusExpired = 'Expired';

  // Roles
  static const String roleMitra = 'mitra';
  static const String roleStaff = 'staff';

  //印尼盾格式
  static String formatCurrency(double amount) {
    return 'Rp ${amount.toStringAsFixed(0).replaceAllMapped(RegExp(r'(\d{1,3})(?=(\d{3})+(?!\d))'), (Match m) => '${m[1]}.')}';
  }

  // Phone number formatting
  static String formatPhoneNumber(String phone) {
    if (phone.startsWith('+62') || phone.startsWith('62')) {
      return phone.replaceAllMapped(RegExp(r'^(\+?62)(\d{3,})(\d{4,})(\d{4,})$'),
          (Match m) => '+62 ${m[2]} ${m[3]} ${m[4]}');
    }
    return phone;
  }

  // Service ports for legacy reference (no longer used)
  static const Map<String, int> servicePorts = {
    'auth': 8081,
    'users': 8082,
    'wallets': 8083,
    'transactions': 8084,
    'products': 8085,
    'integrations': 8086,
  };
}
