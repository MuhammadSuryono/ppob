# Android Native Implementation Requirements for PPOB Mobile Application

## 1. Project Overview

This document outlines the requirements for implementing the PPOB mobile application using Android native development (Kotlin + Jetpack Compose) based on the UI/UX Design System.

## 2. Technology Stack

- **Language:** Kotlin
- **UI Framework:** Jetpack Compose
- **Architecture:** MVVM + Clean Architecture
- **Dependency Injection:** Hilt
- **Networking:** Retrofit + OkHttp
- **Local Storage:** Room Database + DataStore
- **Image Loading:** Coil
- **Navigation:** Jetpack Navigation Compose
- **Authentication:** JWT with Refresh Token
- **Push Notifications:** FCM (Firebase Cloud Messaging)

## 3. Minimum Requirements

- Android API Level 26 (Android 8.0) minimum
- Target API Level 34 (Android 14)
- APK size under 50MB
- Cold start under 2 seconds
- 60fps smooth scrolling

## 4. Module Structure

```
app/
├── src/main/
│   ├── java/com/ppob/app/
│   │   ├── PPOBApp.kt                    # Application class
│   │   ├── di/                            # Dependency injection (Hilt)
│   │   │   ├── AppModule.kt
│   │   │   ├── NetworkModule.kt
│   │   │   ├── DatabaseModule.kt
│   │   │   └── RepositoryModule.kt
│   │   ├── data/
│   │   │   ├── local/                     # Local data sources
│   │   │   │   ├── dao/                   # Room DAOs
│   │   │   │   ├── entity/                # Room entities
│   │   │   │   ├── database/             # Room database
│   │   │   │   └── datastore/             # DataStore preferences
│   │   │   ├── remote/                    # Remote data sources
│   │   │   │   ├── api/                   # Retrofit API interfaces
│   │   │   │   ├── model/                 # Remote response models (DTOs)
│   │   │   │   └── interceptor/           # Network interceptors
│   │   │   └── repository/               # Repository implementations
│   │   ├── domain/
│   │   │   ├── model/                     # Domain models
│   │   │   ├── repository/               # Repository interfaces
│   │   │   └── usecase/                   # Use cases
│   │   ├── presentation/
│   │   │   ├── auth/                      # Auth screens
│   │   │   │   ├── LoginScreen.kt
│   │   │   │   ├── RegisterScreen.kt
│   │   │   │   ├── OtpScreen.kt
│   │   │   │   ├── PinScreen.kt
│   │   │   │   └── AuthViewModel.kt
│   │   │   ├── home/                      # Home/Dashboard screens
│   │   │   │   ├── HomeScreen.kt
│   │   │   │   ├── DashboardViewModel.kt
│   │   │   │   └── components/
│   │   │   ├── wallet/                    # Wallet screens
│   │   │   │   ├── WalletScreen.kt
│   │   │   │   ├── WalletViewModel.kt
│   │   │   │   └── components/
│   │   │   ├── transaction/              # Transaction screens
│   │   │   │   ├── TransactionListScreen.kt
│   │   │   │   ├── TransactionDetailScreen.kt
│   │   │   │   ├── ProductListScreen.kt
│   │   │   │   ├── ConfirmTransactionScreen.kt
│   │   │   │   └── components/
│   │   │   ├── staff/                     # Staff management screens
│   │   │   │   ├── StaffListScreen.kt
│   │   │   │   ├── AddStaffScreen.kt
│   │   │   │   ├── StaffDetailScreen.kt
│   │   │   │   └── components/
│   │   │   └── profile/                   # Profile/Settings screens
│   │   │       ├── ProfileScreen.kt
│   │   │       ├── SettingsScreen.kt
│   │   │       ├── DeviceManagementScreen.kt
│   │   │       └── components/
│   │   ├── navigation/                   # Navigation
│   │   │   ├── PPOBNavHost.kt
│   │   │   ├── Screen.kt                  # Route definitions
│   │   │   └── NavGraph.kt
│   │   ├── theme/                        # Design system
│   │   │   ├── Color.kt
│   │   │   ├── Theme.kt
│   │   │   ├── Shape.kt
│   │   │   └── Typography.kt
│   │   └── widgets/                      # Reusable components
│   │       ├── PPOBButton.kt
│   │       ├── PPOBInput.kt
│   │       ├── WalletCard.kt
│   │       ├── CategoryCard.kt
│   │       ├── ProductCard.kt
│   │       ├── TransactionItem.kt
│   │       ├── StatusBadge.kt
│   │       ├── EmptyState.kt
│   │       └── LoadingSkeleton.kt
│   ├── res/
│   │   ├── values/
│   │   │   ├── strings.xml               # Indonesian string resources
│   │   │   ├── colors.xml
│   │   │   ├── themes.xml
│   │   │   └── dimens.xml
│   │   ├── values-id/                     # Indonesian localization
│   │   │   └── strings.xml
│   │   ├── drawable/                      # Icons and images
│   │   ├── font/                          # Custom fonts
│   │   └── navigation/                    # Navigation graphs
│   └── assets/
└── src/test/                              # Unit tests
└── src/androidTest/                       # Instrumentation tests

```

## 5. API Integration

The Android application will consume backend microservices documented in [MICROSERVICES_API_DOC.md](../MICROSERVICES_API_DOC.md). All API calls use HTTPS with JWT Bearer token authentication.

### 5.1 Base URLs & Services

All services share the same domain with path-based versioning:

**Base URL:** `https://fedora.sinauplatform.id/api/v1`

| Service | Base Path | Primary Use |
|---|---|---|
| Auth Service | `auth` | Registration, login, OTP, token refresh |
| User Service | `user` | Profile, staff management, notifications |
| Product Service | `product` | Product catalog, category listing, sync status |
| Wallet Service | `wallet` | Balance queries, holds, transfers, top-ups |
| Transaction Service | `transaction` | Transaction lifecycle, history, reports |
| Integration Service | `integration` | Digiflazz gateway, provider config |

Full endpoint examples:
- `POST https://fedora.sinauplatform.id/api/v1/auth/register`
- `GET https://fedora.sinauplatform.id/api/v1/product?category=pulsa`
- `GET https://fedora.sinauplatform.id/api/v1/wallet/:id/balance`
- `POST https://fedora.sinauplatform.id/api/v1/transaction/initiate`
- `GET https://fedora.sinauplatform.id/api/v1/user/staff`

### 5.2 Retrofit API Interfaces

Define Retrofit interfaces per service domain:

```kotlin
interface AuthApiService {
    @POST("auth/register")
    suspend fun register(@Body request: RegisterRequest): ApiResponse<AuthResponse>

    @POST("auth/login")
    suspend fun login(@Body request: LoginRequest): ApiResponse<AuthResponse>

    @POST("auth/verify-otp")
    suspend fun verifyOtp(@Body request: VerifyOtpRequest): ApiResponse<VerifyOtpResponse>

    @POST("auth/refresh")
    suspend fun refreshToken(@Body request: RefreshTokenRequest): ApiResponse<RefreshResponse>

    @POST("auth/logout")
    suspend fun logout(@Header("Authorization") token: String): ApiResponse<Unit>

    @POST("auth/change-pin")
    suspend fun changePin(
        @Header("Authorization") token: String,
        @Body request: ChangePinRequest
    ): ApiResponse<Unit>
}
```

### 5.3 Request Headers

All authenticated requests must include:

```
Authorization: Bearer <JWT_ACCESS_TOKEN>
X-Trace-ID: <UUID>  // Optional, for distributed tracing
Idempotency-Key: <UUID>  // Required for /transaction/initiate
```

### 5.4 Error Handling

Standard error response format:

```json
{
  "error": {
    "code": "TRANSACTION_INSUFFICIET_BALANCE",
    "message": "Saldo tidak mencukupi",
    "details": { "required": 50000, "current": 25000 },
    "trace_id": "uuid-string",
    "timestamp": "2026-05-08T..."
  }
}
```

Common HTTP codes:
- `400` — Validation error or business rule violation
- `401` — Token expired/invalid (trigger token refresh flow)
- `403` — Insufficient permissions or locked account
- `429` — Rate limit exceeded
- `502` — Digiflazz API error or timeout
- `503` — Circuit breaker open or maintenance

Implement automatic token refresh in OkHttp interceptor. Retry failed request after obtaining new token.

### 5.5 Network Layer Architecture

```
NetworkModule.kt (Hilt)
├── providesOkHttpClient(): OkHttpClient
│   ├── Interceptors:
│   │   ├── AuthInterceptor (adds Authorization header)
│   │   ├── IdempotencyInterceptor (adds UUID for initiate calls)
│   │   ├── LoggingInterceptor (debug only)
│   │   └── ErrorInterceptor (parses error body, throws typed exceptions)
│   └── Timeout: 30s connect/read/write
├── providesRetrofit(): Retrofit
│   ├── MoshiConverterFactory
│   └── baseUrl from BuildConfig
└── providesApiServices(): AuthApi, UserApi, ProductApi, WalletApi, TransactionApi, IntegrationApi
```

### 5.6 Data Synchronization

- **Product Sync:** Call `POST /product/sync/prepaid` and `POST /product/sync/postpaid` on app start (with exponential backoff retry). Store `GET /product/sync/status` lastSync timestamps locally.
- **Offline Support:** Queue transaction initiation requests locally when offline; sync when network restored. Display pending status in transaction list.

## 6. User Flow Implementation

Screen flows derived from [UI/UX Design System and Flutter Component Guidelines for PPOB Mobile Application.md](UI/UX%20Design%20System%20and%20Flutter%20Component%20Guidelines%20for%20PPOB%20Mobile%20Application.md). Implement using Jetpack Compose with shared ViewModels per flow.

### 6.1 Login / Register Flow

Sequence:
1. **PhoneInputScreen** — Country code picker (+62 pre-selected), phone field, "Lanjutkan" button. Validate phone format.
2. **OtpVerifyScreen** — 6-digit OTP input with auto-focus and auto-submit. Countdown timer with "Kirim ulang" link after expiry. Resend OTP via `POST /auth/verify-otp`.
3. **SetCredentialsScreen** (new users only) — Password + Confirm Password with visibility toggle, 6-digit PIN pad (3×3 grid), Confirm PIN pad.
4. **PinLoginScreen** (subsequent logins on trusted device) — Direct 6-digit PIN pad entry. Fallback link "Login dengan password" for untrusted devices.

Navigation: `AuthNavHost` with `Screen.PhoneInput`, `Screen.OtpVerify`, `Screen.SetCredentials`, `Screen.PinLogin`.

### 6.2 Transaction Flow

Five screens orchestrated by `TransactionViewModel`:

1. **CategorySelectionScreen** — Grid of `CategoryCard` (48dp icon + label). Tapping navigates to ProductListScreen with category filter.
2. **ProductSelectionScreen** — Search bar, sort menu (popularity, price low→high, high→low). Each `ProductCard` displays platform price and Mitra selling price input field (last-input cached per product). Bottom sticky input for customer number.
3. **TransactionConfirmationScreen** — Product summary, masked customer number, price breakdown (platform price, margin, total), "Bayar" button. Disable if balance insufficient (pre-check via `GET /wallet/:id/balance`).
4. **PinAuthorizationScreen** — 6-digit PIN pad. On submit: call `POST /transaction/initiate` with `Idempotency-Key` header. Show loading state; on response navigate to ResultScreen.
5. **TransactionResultScreen** — Success (green check, amount, reference, "Selesai"/"Bagikan" buttons), Pending (hourglass, "Sedang diproses...", "Lihat status"), or Failed (red cross, error message, "Coba lagi").

All screens support back navigation; transaction state persisted in `TransactionRepository` during flow.

### 6.3 Staff Management Flow (Mitra Role)

1. **StaffListScreen** — App bar with "+ Tambah" FAB. Search bar. Each `StaffCard` shows avatar (initials), name, phone, today's count/amount, wallet balance, active/inactive toggle. Badge on app bar if `GET /user/staff/pending-count > 0`.
2. **StaffDetailScreen** — Editable fields: name, phone, PIN reset. Sections:
   - Margin Settings: Radio Fixed Fee (Rp/transaksi) vs Revenue Share (%); slider 10–80% (default 60).
   - Daily Limit: count override, amount override.
3. **AddStaffScreen** — Similar form; password required for new staff.
4. **StaffTopUpScreen** — ModalBottomSheet: dropdown select staff, numeric amount keypad, shows Mitra wallet deduction preview, confirm → `POST /user/staff/topup`.

Data source: `UserService.getStaff()`, `POST /user/staff`, `PUT /user/staff/:id`.

### 6.4 Wallet Flow

- **WalletScreen** displays `WalletCard` with available balance (primary green headline) and held balance (pending transactions, collapsible). "Top Up" button triggers `StaffTopUpScreen` (Mitra-only internal transfer).
- Transaction history sub-tab uses same `TransactionListScreen` but filtered by wallet ID (`GET /transaction?wallet_id=:id`).
- Real-time updates via FCM data messages: when a transaction completes, held balance decreases; refresh wallet card automatically.

### 6.5 Reports Flow (Mitra/Admin)

1. **ReportsScreen** — Date range picker (Today | 7d | 30d | Custom).
2. **KPIsRow** — 2×2 grid cards: Total Sales, Platform Profit, Staff Count, Success Rate.
3. **SalesTrendChart** — `fl_chart` LineChart showing daily sales over selected range (fetch `GET /transaction/reports?metric=sales&from=...&to=...`).
4. **StaffPerformanceChart** — Bar chart ranking top 5 staff by transaction count.
5. **RecentTransactionsTable** — Paginated list; tap row → `TransactionDetailScreen`.

## 7. Screen Component Mapping (Jetpack Compose)

| UI Component | Flutter Reference | Compose Implementation |
|---|---|---|
| Primary Button | `PPOBButton.primary` | `PPOBButton(variant = Primary, ...)` |
| Secondary Button | `PPOBButton.secondary` | `PPOBButton(variant = Secondary, ...)` |
| Input Field | `PPOBInput` | `OutlinedTextField` with `PPOBTheme.InputStyle` |
| PIN Pad | PIN input grid (3×3) | `PinInputScreen` Composable with `LazyVerticalGrid` |
| Transaction Card | `TransactionCard` | `TransactionItem` Composable row |
| Wallet Card | `WalletCard` | `WalletCard` Composable with balance typography |
| Category Card | `CategoryCard` | `CategoryCard` Composable with `Icon` + `Text` |
| Status Badge | `statusBadge()` | `StatusBadge(status: String)` Composable |
| Empty State | Empty state designs | `EmptyStateScreen(icon, message, cta?)` |
| Skeleton Loader | `TransactionListSkeleton` | `TransactionListShimmer()` using `accompanist-placeholder` |

### 7.1 Theming

Create `PPOBTheme.kt` mirroring Flutter `ThemeData`:

```kotlin
object PPOBTheme {
    val colorScheme = lightColorScheme(
        primary = Color(0xFF4CAF50),      // Primary Green
        secondary = Color(0xFF2196F3),    // Secondary Blue
        tertiary = Color(0xFFFF9800),     // Accent Orange
        background = Color(0xFFF5F5F5),   // Background Light
        surface = Color(0xFFFFFFFF),      // Surface White
        error = Color(0xFFF44336)         // Error Red
    )
    val typography = Typography(
        displayLarge = TextStyle(weight = FontWeight.Bold, size = 32, lineHeight = 40),
        headlineMedium = TextStyle(weight = FontWeight.Medium, size = 24, lineHeight = 32),
        titleLarge = TextStyle(weight = FontWeight.Medium, size = 20, lineHeight = 28),
        bodyLarge = TextStyle(weight = FontWeight.Normal, size = 16, lineHeight = 24),
        bodyMedium = TextStyle(weight = FontWeight.Normal, size = 14, lineHeight = 20),
        labelSmall = TextStyle(weight = FontWeight.Medium, size = 12, lineHeight = 16)
    )
    val shapes = Shapes(
        small = RoundedCornerShape(4.dp),
        medium = RoundedCornerShape(8.dp),
        large = RoundedCornerShape(16.dp)
    )
}
```

Apply via `MaterialTheme(colorScheme = PPOBTheme.colorScheme, typography = PPOBTheme.typography, shapes = PPOBTheme.shapes)` in `PPOBApp.kt`.

## 8. Navigation Pattern

Bottom navigation bar (5 items) + Drawer menu. Implement via `NavGraph.kt` and `PPOBScaffold.kt`:

- **Routes:** `/home`, `/transactions`, `/wallet`, `/staff` (Mitra only), `/profile`
- **Active state:** Primary green icon (filled); inactive: outlined gray
- **Badges:** Wallet badge if `wallet.balance.held > 0`; Staff badge (Mitra) if `pending_staff_requests > 0`
- **Drawer items:** Settings, Ganti PIN, Perangkat Terpercaya, Bantuan & Pendorongan, Keluar

Use `NavigationSuiteScaffold` or custom `BottomNavigationBar` Composable.

## 9. Offline & Caching Strategy

### 9.1 Local Storage
- **Room entities:** `UserEntity`, `WalletEntity`, `TransactionEntity`, `ProductEntity`, `CategoryEntity`
- **DataStore:** `AuthPreferences` (access_token, refresh_token, user_id), `AppPreferences` (theme, language)
- **Caching:** Repository layer returns cached data immediately; network call refreshes asynchronously. `TransactionRepository.getTransactions()` reads from Room; `refreshTransactions()` fetches from API and updates Room.

### 9.2 Idempotency
For `POST /transaction/initiate`, generate UUID client-side and store in `PendingTransaction` Room entity with status PENDING. Retry same UUID on retry to avoid duplicate charges.

### 9.3 Optimistic UI
Update wallet balance and transaction list optimistically on transaction submission; rollback on error.

## 10. Performance Targets

- Cold start < 2s (measured from process start to first frame)
- 60fps smooth scrolling: use `LazyColumn`/`LazyVerticalGrid`, avoid layouts with deep nesting in list items
- APK < 50MB: enable R8 shrinker, split per ABI, remove unused resources
- Memory < 150MB: use `remember`/`DisposableEffect` to clean large bitmaps; Coil with `memoryCache` enabled
- Network: Enable HTTP/2, connection pooling, GZIP compression

Load product icons with `AsyncImage` (Coil) and `crossfade(true)`. Use `rememberImagePainter`.

## 11. Security Requirements

- Store tokens in `EncryptedSharedPreferences` or `DataStore` with `Encryption`.
- Do not log sensitive data (passwords, PINs, tokens); use `HttpLoggingInterceptor.Level.BODY` only in debug.
- Implement certificate pinning for API domains (optional phase 2).
- Prevent screen capture on sensitive screens (PIN entry) via `FLAG_SECURE` in activity.
- Validate all input client-side before API submission; never rely solely on server validation.

## 12. Accessibility

- Content descriptions for icon-only buttons (`contentDescription = "Notifikasi"`).
- Touch targets minimum 48×48 dp.
- Font scaling support: use `sp` for text sizes; test at 200% system font.
- Color contrast ratio ≥ 4.5:1 (verified via contrast checker).
- Logical focus order: set `android:focusable` and `android:focusableInTouchMode` appropriately.

## 13. Internationalization

Locale: Indonesian (`id_ID`) only for MVP.

String resources in `res/values/strings.xml`:
```xml
<string name="login_title">Masuk / Daftar</string>
<string name="error_insufficient_balance">Saldo tidak mencukupi</string>
```

Number formatting: `NumberFormat.getCurrencyInstance(Locale("id", "ID")).apply { currency = Currency.getInstance("IDR") }`.

## 14. Testing Strategy

### 14.1 Unit Tests
- Repository implementations with mocked `ApiService` and `Dao`
- UseCases: `InitiateTransactionUseCase`, `GetProductsUseCase`
- ViewModels: test state emissions for loading/success/error

### 14.2 Integration Tests
- Room database CRUD operations
- Retrofit service calls with `MockWebServer`

### 14.3 UI Tests (Compose)
- Login flow: phone → OTP → PIN success path
- Transaction flow: category → product → confirm → PIN → success modal
- Staff list: add staff, set margin, verify card updates

Run with `./gradlew test` and `./gradlew connectedAndroidTest`.

## 15. DevOps & Monitoring

- Crash reporting: Firebase Crashlytics
- Analytics: Firebase Analytics for events (`transaction_initiated`, `transaction_success`, `staff_added`)
- Performance monitoring: Firebase Performance Monitoring (cold start, network latency)
- Feature flags: Remote Config to toggle product sync, top-up availability

## Appendix

### A. API Endpoint Quick Reference

See [MICROSERVICES_API_DOC.md](../MICROSERVICES_API_DOC.md) for complete endpoint catalog.

**Key endpoints by feature:**
- Auth: `POST /api/v1/auth/register`, `POST /api/v1/auth/login`, `POST /api/v1/auth/verify-otp`
- Products: `GET /api/v1/product`, `POST /api/v1/product/sync/prepaid`, `GET /api/v1/category`
- Wallet: `GET /api/v1/wallet/:id/balance`, `POST /api/v1/wallet/:id/hold`, `POST /api/v1/wallet/:id/debit`
- Transaction: `POST /api/v1/transaction/initiate` (idempotent), `GET /api/v1/transaction/:id`, `GET /api/v1/transaction/history`
- Staff: `GET /api/v1/user/staff`, `POST /api/v1/user/staff`, `PUT /api/v1/user/staff/:id`

### B. Flutter Component Mapping

All Compose widgets should adhere to visual and behavioral specifications defined in [UI/UX Design System and Flutter Component Guidelines for PPOB Mobile Application.md](UI/UX%20Design%20System%20and%20Flutter%20Component%20Guidelines%20for%20PPOB%20Mobile%20Application.md) — including colors, typography, spacing (8dp grid), corner radii (4/8/16dp), and component states (enabled, disabled, loading, error).