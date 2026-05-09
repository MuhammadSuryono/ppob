class AppConstants {
  static const String appName = 'PPOB';
  static const String apiBaseUrl = 'https://192.168.100.23';
  static const String apiVersion = 'v1';

  // Service base URLs (different ports)
  static const String baseAuth = 'https://192.168.100.23:8081/api/v1';
  static const String baseUsers = 'https://192.168.100.23:8082/api/v1';
  static const String baseProducts = 'https://192.168.100.23:8085/api/v1';
  static const String baseWallets = 'https://192.168.100.23:8083/api/v1/wallets';
  static const String baseTransactions = 'https://192.168.100.23:8084/api/v1/transactions';
  static const String baseIntegrations = 'https://192.168.100.23:8086/api/v1/integrations';

  // API endpoints
  static const String endpointAuth = baseAuth + '/auth';
  static const String endpointUsers = baseUsers;
  static const String endpointWallets = baseWallets;
  static const String endpointProducts = baseProducts;
  static const String endpointTransactions = baseTransactions;
  static const String endpointStaff = baseUsers + '/users';
  static const String endpointIntegration = baseIntegrations + '/integrations';

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

  // Microservice ports for development
  static const Map<String, int> servicePorts = {
    'auth': 8081,
    'users': 8082,
    'wallets': 8083,
    'transactions': 8084,
    'products': 8085,
    'integrations': 8086,
  };
}
