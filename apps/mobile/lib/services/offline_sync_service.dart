import 'dart:convert';
import 'package:hive/hive.dart';
import '../models/pending_sync_item.dart';
import '../utils/constants.dart';

class OfflineSyncService {
  static final OfflineSyncService _instance = OfflineSyncService._internal();
  factory OfflineSyncService() => _instance;
  OfflineSyncService._internal();

  late Box<PendingSyncItem> _pendingBox;
  bool _isSyncing = false;

  Future<void> initialize() async {
    _pendingBox = await Hive.openBox<PendingSyncItem>(AppConstants.boxPendingSync);
  }

  Future<void> queueForSync({
    required String type,
    required Map<String, dynamic> data,
  }) async {
    final item = PendingSyncItem(
      id: '${type}_${DateTime.now().millisecondsSinceEpoch}',
      type: type,
      data: data,
      createdAt: DateTime.now(),
    );

    await _pendingBox.put(item.id, item);
  }

  Future<List<PendingSyncItem>> getPendingItems() async {
    return _pendingBox.values.toList();
  }

  Future<void> removePendingItem(String itemId) async {
    await _pendingBox.delete(itemId);
  }

  Future<void> startBackgroundSync({int maxRetries = 3}) async {
    if (_isSyncing) return;
    _isSyncing = true;

    try {
      final pending = await getPendingItems();

      for (final item in pending) {
        if (item.retryCount >= maxRetries) {
          // Mark as failed, stop retrying
          await removePendingItem(item.id);
          continue;
        }

        try {
          // Attempt to sync with server
          // TODO: Implement actual API call based on item.type
          // await _syncItem(item);

          // If successful, remove from pending
          await removePendingItem(item.id);
        } catch (e) {
          // Increment retry count and update
          await _pendingBox.put(
            item.id,
            item.copyWith(
              retryCount: item.retryCount + 1,
              lastAttemptAt: DateTime.now(),
              errorMessage: e.toString(),
            ),
          );
        }
      }
    } finally {
      _isSyncing = false;
    }
  }

  // Check if device is online
  Future<bool> isOnline() async {
    // TODO: Implement connectivity check using connectivity_plus
    return true; // Mock
  }

  // Schedule periodic sync
  void schedulePeriodicSync(Duration interval) {
    // TODO: Use Workmanager or background_fetch for periodic sync
  }

  // Sync transactions when coming back online
  Future<void> syncWhenOnline() async {
    final online = await isOnline();
    if (online) {
      await startBackgroundSync();
    }
  }
}
