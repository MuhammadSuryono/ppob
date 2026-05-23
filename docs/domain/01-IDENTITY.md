# 🔑 Domain: Identity & Access

## 1. Core Mission
The Identity domain is responsible for managing user accounts, verifying authenticity, and controlling access to platform resources.

## 2. Key Concepts

### 2.1 User Identity
A unique entity identified primarily by a validated phone number.
- **Attributes:** Phone Number, Name, Active Role.
- **Rules:** Phone numbers must follow the `+62` format and be unique.

### 2.2 Device Trust
A security signal used to reduce friction for recurring users.
- **Fingerprint:** A combination of device hardware and software identifiers.
- **Trust Score:** A calculated value (0-100) determining the level of authentication required (Adaptive Auth).

### 2.3 Authentication Factors
- **Password:** Used for high-security login and untrusted devices.
- **OTP (One-Time Password):** A 6-digit numeric code sent via SMS/WhatsApp for verification.
- **PIN:** A 6-digit code specifically for authorizing financial transactions.

## 3. Business Rules

### 3.1 Account Security
- **Failed Attempts:** 5 consecutive wrong PIN entries will lock the account for 1 hour.
- **Adaptive Auth:** If trust score < 70, OTP is mandatory.

### 3.2 Session Management
- **Stateless Tokens:** Access is granted via JWT (JSON Web Tokens).
- **Invalidation:** Sessions are immediately invalidated if credentials (Password/PIN) are changed.

## 4. Domain Logic
- **Hashing Policy:** Passwords use `bcrypt`, while transaction PINs use `Argon2id` for maximum resistance against specialized hardware attacks.
- **Adaptive Challenge:** The system dynamically decides whether to present a Password or OTP screen based on real-time trust scoring.
