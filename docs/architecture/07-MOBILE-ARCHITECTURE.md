# 📱 Mobile Architecture

## 1. Overview
The PPOB mobile application is a native Android application built with **Kotlin** and **Jetpack Compose**. It follows an offline-first strategy to ensure a smooth user experience in low-connectivity environments common in Indonesia.

## 2. Technology Stack
- **UI Framework:** Jetpack Compose (Declarative UI).
- **Architecture:** MVVM (Model-View-ViewModel) + Clean Architecture.
- **Dependency Injection:** Hilt.
- **Networking:** Retrofit + OkHttp with Coroutines.
- **Local Storage:** Room Database (SQLCipher encrypted).
- **Asynchronous Flow:** Kotlin Coroutines & StateFlow.
- **Image Loading:** Coil with disk caching.

## 3. Core Architecture Layers

### 3.1 UI Layer (Presentation)
- **Jetpack Compose Screens:** Stateless composables for UI rendering.
- **ViewModels:** Maintain UI state using `StateFlow`; survive configuration changes.
- **Navigation:** Jetpack Navigation Compose with deep link support.

### 3.2 Domain Layer
- **Use Cases:** Encapsulate single business actions (e.g., `InitiateTransactionUseCase`).
- **Repository Interfaces:** Define data contracts for the presentation layer.

### 3.3 Data Layer
- **Repositories:** Single source of truth. Decide whether to return cached data from Room or fetch fresh data from the API.
- **Remote Data Source:** Retrofit interfaces for microservice endpoints.
- **Local Data Source:** Room DAOs for persistent offline storage.

## 4. Offline-First Strategy
- **Local Cache:** Product catalogs and last 30 days of transaction history are stored in Room.
- **Optimistic UI:** Local balance and history are updated immediately upon transaction initiation, then reconciled once the server responds.
- **Write Queue:** Transactions initiated offline are queued in a `PendingSync` table and processed by **WorkManager** as soon as connectivity is restored.

## 5. Security Protocols
- **Token Storage:** JWT tokens are stored in `EncryptedSharedPreferences`.
- **Biometric Integration:** Uses `BiometricPrompt` to authorize session unlocks and high-value transactions.
- **Hardware Security:** PIN hashes and encryption keys are stored in the **Android Keystore** (hardware-backed).
- **Anti-Tamper:** Basic root detection and certificate pinning for microservice domains.

## 6. Real-Time Updates
- **Push Notifications:** Firebase Cloud Messaging (FCM) for transaction status changes and low-balance alerts.
- **Notification Channels:** Segregated channels for `Transactions`, `Security`, and `Staff Management`.
