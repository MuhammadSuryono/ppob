# 🔄 Product Synchronization Flow

## 1. Overview
The Product Service maintains an up-to-date catalog of digital products by synchronizing with the Digiflazz provider API.

## 2. Sync Schedule
| Product Type | Frequency | Trigger |
|---|---|---|
| **Prepaid** | Hourly | Cron Job (at :05) |
| **Postpaid** | Every 5 min | Cron Job |

## 3. Synchronization Flow
1.  **Trigger:** Kubernetes CronJob triggers the `/sync` endpoint in the Product Service.
2.  **Distributed Lock:** System attempts to acquire a Redis Redlock (`product-sync-lock`). If another instance is syncing, the request exits immediately.
3.  **Fetch:** Product Service calls Digiflazz `price-list` API with HMAC signature.
4.  **Upsert Logic:**
    - New products are inserted.
    - Existing products are updated if the price or status has changed.
    - Products not present in the provider's list are marked `is_active = false` (soft deleted).
5.  **Price Calculation:**
    - `Platform Price = Base Price * (1 + platform_markup_percent)`.
    - `platform_markup_percent` is retrieved from the `system_settings` table.
6.  **Cache Invalidation:** After a successful sync, the service invalidates local Redis caches for product lists.

## 4. Error Handling
- **Limit Reached (RC 83):** Digiflazz limits pricelist queries to 1 every 5 minutes. If hit, the system waits 4 minutes before retrying.
- **Drift Detection:** If the price changes by >10% compared to the previous version, a warning is logged for manual review.
- **Circuit Breaker:** If Digiflazz is unreachable, the sync is skipped, and the system continues to serve the last known data.
