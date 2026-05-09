import 'package:firebase_messaging/firebase_messaging.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter_local_notifications/flutter_local_notifications.dart';

class NotificationService {
  static final NotificationService _instance = NotificationService._internal();
  factory NotificationService() => _instance;
  NotificationService._internal();

  final FirebaseMessaging _firebaseMessaging = FirebaseMessaging.instance;
  final FlutterLocalNotificationsPlugin _localNotifications = FlutterLocalNotificationsPlugin();

  Future<void> initialize() async {
    // Request permissions
    await _firebaseMessaging.requestPermission(
      alert: true,
      badge: true,
      sound: true,
    );

    // Get FCM token
    String? token = await _firebaseMessaging.getToken();
    if (kDebugMode) {
      print('FCM Token: $token');
    }
    // TODO: Send token to backend

    // Initialize local notifications
    const AndroidInitializationSettings initializationSettingsAndroid =
        AndroidInitializationSettings('@mipmap/ic_launcher');

    const DarwinInitializationSettings initializationSettingsDarwin =
        DarwinInitializationSettings();

    const InitializationSettings initializationSettings = InitializationSettings(
      android: initializationSettingsAndroid,
      iOS: initializationSettingsDarwin,
    );

     final bool? result = await _localNotifications.initialize(
       settings: initializationSettings,
       onDidReceiveNotificationResponse: (NotificationResponse response) {
         // Handle notification tap
       },
     );
     if (kDebugMode) {
       print('Local notifications initialized: $result');
     }

    // Handle foreground messages
    FirebaseMessaging.onMessage.listen((RemoteMessage message) {
      _showNotification(message);
    });

    // Handle when app is in background but opened from notification
    FirebaseMessaging.onMessageOpenedApp.listen((RemoteMessage message) {
      _handleNotificationTap(message);
    });
  }

  Future<void> _showNotification(RemoteMessage message) async {
    final notification = message.notification;
    final android = message.notification?.android;

    if (notification != null) {
      const AndroidNotificationDetails androidPlatformChannelSpecifics =
          AndroidNotificationDetails(
        'ppob_transactions',
        'PPOB Transactions',
        channelDescription: 'Transaction notifications',
        importance: Importance.max,
        priority: Priority.high,
        icon: '@mipmap/ic_launcher',
      );

      const NotificationDetails platformChannelSpecifics = NotificationDetails(
        android: androidPlatformChannelSpecifics,
      );

       await _localNotifications.show(
         id: notification.hashCode,
         title: notification.title,
         body: notification.body,
         notificationDetails: platformChannelSpecifics,
       );
    }
  }

  void _handleNotificationTap(RemoteMessage message) {
    // Navigate based on notification type
    final data = message.data;
    final type = data['type'];

    switch (type) {
      case 'transaction_success':
        // Navigate to transaction detail
        // TODO: Implement navigation
        break;
      case 'transaction_failed':
        // Show error or navigate
        break;
      case 'low_balance':
        // Navigate to wallet
        break;
    }
  }

  Future<void> subscribeToTopic(String topic) async {
    await _firebaseMessaging.subscribeToTopic(topic);
  }

  Future<void> unsubscribeFromTopic(String topic) async {
    await _firebaseMessaging.unsubscribeFromTopic(topic);
  }

  Future<String?> getToken() async {
    return await _firebaseMessaging.getToken();
  }
}
