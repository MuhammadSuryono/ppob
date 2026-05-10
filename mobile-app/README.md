# PPOB Android Mobile Application

## Overview

This is the native Android client application for the **PPOB (Payment Point Online Bank)** multi-tenant platform. Built with modern Android development best practices using Kotlin, Jetpack Compose, and clean architecture.

## Architecture

### Tech Stack

| Layer | Technology |
|-------|-----------|
| Language | Kotlin |
| UI Framework | Jetpack Compose |
| Architecture | MVVM + Clean Architecture |
| Dependency Injection | Hilt (Dagger) |
| Networking | Retrofit + OkHttp |
| Local Storage | Room Database + DataStore |
| Image Loading | Coil |
| Navigation | Jetpack Navigation Compose |
| Authentication | JWT with Refresh Token |
| Push Notifications | FCM (Firebase Cloud Messaging) |
| Biometric Auth | AndroidX Biometric Library |

### Module Structure

```
app/src/main/java/com/yonotech/ppob/
├── di/                              # Dependency Injection (Hilt)
│   └── AppModule.kt
├── data/
│   ├── local/                       # Local data sources
│   │   ├── dao/                     # Room DAOs
│   │   ├── entity/                  # Room entities
│   │   ├── database/               # Room database
│   │   └── datastore/              # DataStore preferences
│   ├── remote/                      # Remote data sources
│   │   ├── api/                     # Retrofit API interfaces
│   │   ├── model/                   # DTOs
│   │   └── interceptor/            # Network interceptors
│   └── repository/                 # Repository implementations
├── domain/
│   ├── model/                       # Domain models
│   ├── repository/                 # Repository interfaces
│   └── usecase/                     # Use cases
├── presentation/
│   ├── auth/                        # Authentication screens
│   ├── home/                        # Home/Dashboard
│   ├── transaction/                # Transaction flow
│   ├── transactiondetail/          # Transaction detail
│   ├── wallet/                      # Wallet screens
│   ├── staff/                       # Staff management
│   ├── profile/                     # Profile/Settings
│   ├── components/                 # Shared UI components
│   └── navigation/                 # Navigation graphs
├── theme/                           # Design system
├── widgets/                         # Reusable components
├── services/                        # Background services
└── ui/                              # Activity & base UI classes
```

## API Reference

All API calls go through the base URL:
```
https://fedora.sinauplatform.id/api/v1/{service}/{endpoint}
```

### Available Services

| Service | Base Path | Description |
|---------|-----------|-------------|
| Auth | `/auth` | Registration, login, OTP, token refresh |
| User | `/users` | Profile, staff management |
| Product | `/products` | Product catalog, categories |
| Wallet | `/wallets` | Balance, holds, transfers, top-up |
| Transaction | `/transactions` | Transaction lifecycle |
| Integration | `/integration` | Digiflazz gateway |

## Installation & Setup

### Prerequisites

- Android Studio Hedgehog (2023.1.1+) or newer
- JDK 17
- Android SDK 34 (API Level 34)
- Google Firebase project configured (see firebase setup below)

### Setup Steps

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd mobile-app
   ```

2. **Configure Firebase**
   - Place `google-services.json` in `app/` directory
   - Configure Firebase project at [Firebase Console](https://console.firebase.google.com/)

3. **Configure API Base URL** (optional)
   - Edit `build.gradle.kts` in app module to change `BuildConfig.API_BASE_URL`

4. **Open in Android Studio**
   - File → Open → Select `mobile-app/`
   - Wait for Gradle sync to complete

5. **Run the app**
   - Select a device/emulator (API 26+)
   - Click Run ▶ or `./gradlew installDebug`

### Build Commands

```bash
# Build debug APK
./gradlew assembleDebug

# Run unit tests
./gradlew testDebugUnitTest

# Run lint
./gradlew lintDebug

# Build release APK
./gradlew assembleRelease
```

## Development Guidelines

### Code Style

- Follow [Android Kotlin Style Guide](https://developer.android.com/kotlin/style-guide)
- Use Kotlin coroutines for async operations
- Use Flow/StateFlow for reactive state management
- All UI must use Jetpack Compose (no XML layouts)
- Apply dependency injection via Hilt for all Android components

### Naming Conventions

- **Use cases**: `{Action}{Subject}UseCase.kt` (e.g., `GetProductsUseCase`)
- **ViewModels**: `{Screen}ViewModel.kt` (e.g., `HomeViewModel`)
- **Composables**: `{Component}Screen.kt` or `{Component}Card.kt`
- **DI modules**: `{Module}Module.kt` (e.g., `NetworkModule`)
- **Extensions**: `{Type}Extensions.kt`
- **Test files**: Match the source file name with `Test` suffix

### Branch Strategy

- `main` — Production-ready code
- `develop` — Integration branch
- `feature/xxx` — Feature branches
- `hotfix/xxx` — Bug fix branches

### Commit Message Convention

```
feat(module): add new feature description
fix(module): fix bug description
chore: maintenance task
docs: documentation changes
test: test additions/changes
refactor: code restructuring
```

## Testing Strategy

### Unit Tests (`src/test/`)

- ViewModel state emission tests using Turbine
- UseCase tests with mocked repositories
- Domain logic validation
- Target coverage: ≥80% for core modules

### Instrumentation Tests (`src/androidTest/`)

- Compose UI tests with `createComposeRule()`
- Navigation graph verification
- User flow tests (login → browse → purchase)
- Offline/online behavior testing

Run tests:
```bash
./gradlew testDebugUnitTest
./gradlew connectedDebugAndroidTest
```

## Security Considerations

- JWT tokens stored in encrypted DataStore
- TLS certificate pinning (Phase 2)
- PIN hashed with bcrypt/argon2id
- Screen capture prevention on sensitive screens
- Input validation on client and server
- SQLCipher for encrypted Room database
- Root/jailbreak detection (optional)

## Performance Targets

| Metric | Target |
|--------|--------|
| APK Size | < 30MB (split per ABI) |
| Cold Start | < 2 seconds |
| Frame Render | < 16ms (60fps) |
| Memory Usage | < 120MB typical |
| API Response | < 1s (non-Digiflazz) |

## Phase 5 Deliverables

This module implements **Phase 5** of the AI Execution Timeline:

- ✅ Project setup with Hilt DI and Retrofit networking
- ✅ Auth flow: Phone Input → OTP → Set Credentials → PIN Login
- ✅ Home screen with balance card, quick actions, recent transactions
- ✅ Product catalog with categories and search
- ✅ Full transaction flow: Category → Product → Confirm → PIN → Result
- ✅ Transaction history with filters
- ✅ Wallet screen with real-time balance
- ✅ Staff management (Mitra role)
- ✅ Profile screen with settings
- ✅ Offline sync queue with WorkManager
- ✅ FCM push notifications
- ✅ Biometric authentication (optional)
- ✅ Reusable widget library (PPOBButton, PPOBInput, etc.)
- ✅ Design system with PPOB Theme

## Dependencies

Key libraries:
- Jetpack Compose BOM 2024.09.02
- Hilt 2.51.1
- Retrofit 2.9.0 + Moshi
- Room 2.6.1
- Firebase BoM 34.13.0
- Coil 3.0.3
- Accompanist 0.36.1
- Paging 3.3.1

Full dependency list: `gradle/libs.versions.toml`

## License

Proprietary — PPOB Multi-Tenant Mobile Application

---

**Owner:** Mobile Team Lead
**Created:** 2026-05-09
**Last Updated:** 2026-05-09