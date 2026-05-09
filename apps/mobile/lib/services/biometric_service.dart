import 'package:local_auth/local_auth.dart';
import 'package:flutter/services.dart';

class BiometricService {
  final LocalAuthentication _localAuth = LocalAuthentication();

  Future<bool> isAvailable() async {
    try {
      return await _localAuth.canCheckBiometrics;
    } catch (e) {
      return false;
    }
  }

   Future<bool> authenticate() async {
     try {
       return await _localAuth.authenticate(
         localizedReason: 'Autentikasi untuk transaksi',
         biometricOnly: false,
         sensitiveTransaction: true,
         persistAcrossBackgrounding: true,
       );
     } on PlatformException catch (e) {
       throw Exception('Biometric authentication failed: ${e.message}');
     } catch (e) {
       throw Exception('Biometric not available');
     }
   }

  Future<List<BiometricType>> getAvailableBiometrics() async {
    try {
      return await _localAuth.getAvailableBiometrics();
    } catch (e) {
      return [];
    }
  }
}
