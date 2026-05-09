# QA Test Scenarios - PPOB Mobile Application

## 1. Authentication & Authorization Tests

### 1.1 Registration Flow

| TC-ID | Test Scenario | Pre-requisites | Steps | Expected Result |
|-------|--------------|----------------|-------|-----------------|
| AUTH-001 | Register with new phone number | Clean phone number | 1. Enter phone number<br>2. Enter name<br>3. Submit | - OTP sent to phone<br>- User created in pending state |
| AUTH-002 | Register with existing phone number | Phone already registered | 1. Enter existing phone number<br>2. Enter name<br>3. Submit | - Error: "Phone number already registered"<br>- No new user created |
| AUTH-003 | Register with invalid phone format | None | 1. Enter invalid phone format (e.g., letters) | - Error: "Invalid phone number format"<br>- Registration blocked |
| AUTH-004 | OTP verification success | OTP sent | 1. Enter correct OTP within 5 minutes | - OTP verified<br>- User state changes to active |
| AUTH-005 | OTP verification expired | OTP sent > 5 minutes | 1. Enter OTP after 5 minutes | - Error: "OTP expired"<br>- User can request new OTP |
| AUTH-006 | OTP verification wrong | OTP sent | 1. Enter incorrect OTP | - Error: "Invalid OTP"<br>- Retry counter decreases |
| AUTH-007 | OTP max retries exceeded | OTP sent | 1. Enter wrong OTP 3 times | - Account locked temporarily<br>- Error: "Too many attempts" |
| AUTH-008 | Password creation after OTP | OTP verified | 1. Create password<br>2. Confirm password<br>3. Submit | - Password created<br>- User can login |

### 1.2 Login Flow

| TC-ID | Test Scenario | Pre-requisites | Steps | Expected Result |
|-------|--------------|----------------|-------|-----------------|
| AUTH-009 | Login with correct credentials | Active user | 1. Enter phone<br>2. Enter password<br>3. Submit | - JWT token returned<br>- Login successful |
| AUTH-010 | Login with wrong password | Active user | 1. Enter phone<br>2. Enter wrong password | - Error: "Invalid credentials"<br>- No token returned |
| AUTH-011 | Login with non-existent phone | None | 1. Enter unregistered phone | - Error: "User not found"<br>- No token returned |
| AUTH-012 | Trusted device login | Previously trusted device | 1. Login from trusted device<br>2. Enter PIN directly | - Skip password<br>- Direct access to dashboard |
| AUTH-013 | Untrusted device login | New device | 1. Login from new device | - OTP required<br>- Security alert sent |
| AUTH-014 | Token refresh | Valid refresh token | 1. Use refresh token | - New access token returned |
| AUTH-015 | Logout | Authenticated user | 1. Call logout endpoint | - Token blacklisted<br>- Session ended |

## 2. Role & Multi-Account Tests

### 2.1 Mitra Role

| TC-ID | Test Scenario | Pre-requisites | Steps | Expected Result |
|-------|--------------|----------------|-------|-----------------|
| ROLE-001 | Mitra default role assignment | New user | 1. Register user | - User assigned Mitra role by default |
| ROLE-002 | Add new staff member | Mitra authenticated | 1. Enter staff phone<br>2. Assign role<br>3. Submit | - Staff invitation sent<br>- Staff listed in Mitra's dashboard |
| ROLE-003 | Staff role assignment | Staff invited | 1. Staff accepts invitation | - Staff role activated<br>- Linked to Mitra |
| ROLE-004 | Configure staff permissions | Staff exists | 1. Mitra sets staff permissions | - Permissions saved<br>- Staff restricted accordingly |
| ROLE-005 | Top-up staff wallet | Staff exists | 1. Mitra adds funds to staff wallet | - Staff wallet balance increased |
| ROLE-006 | Staff wallet balance check | Staff wallet funded | 1. View staff wallet | - Correct balance displayed |

### 2.2 Multi-Role Capability

| TC-ID | Test Scenario | Pre-requisites | Steps | Expected Result |
|-------|--------------|----------------|-------|-----------------|
| ROLE-007 | User as Mitra and Staff | User with dual roles | 1. Login<br>2. Switch role | - Role switched successfully<br>- Dashboard updated |
| ROLE-008 | Transaction from Mitra wallet | Dual role user, Mitra role active | 1. Perform transaction | - Deducted from Mitra main wallet |
| ROLE-009 | Transaction from Staff wallet | Dual role user, Staff role active | 1. Perform transaction | - Deducted from Staff sub-wallet |
| ROLE-010 | Role switching persistence | Role selected | 1. Close app<br>2. Reopen | - Last active role preserved |

## 3. Wallet & Balance Management Tests

| TC-ID | Test Scenario | Pre-requisites | Steps | Expected Result |
|-------|--------------|----------------|-------|-----------------|
| WALLET-001 | Mitra wallet creation | New Mitra | 1. Register as Mitra | - Main wallet created with 0 balance |
| WALLET-002 | Staff wallet creation | Staff added | 1. Add staff member | - Sub-wallet created with 0 balance |
| WALLET-003 | Wallet top-up | Authenticated Mitra | 1. Add funds to wallet<br>2. Confirm | - Wallet balance increased<br>- Transaction recorded |
| WALLET-004 | Insufficient balance | Low wallet balance | 1. Attempt transaction exceeding balance | - Transaction rejected<br>- Error: "Insufficient balance" |
| WALLET-005 | Balance hold | Active wallet | 1. Place hold for transaction | - Hold amount deducted from available balance<br>- Status: pending |
| WALLET-006 | Hold release | Active hold | 1. Cancel/rollback transaction | - Hold amount released<br>- Available balance restored |
| WALLET-007 | Transaction deduction | Successful transaction | 1. Complete transaction | - Balance deducted<br>- Transaction recorded |
| WALLET-008 | Wallet audit trail | Multiple transactions | 1. View wallet history | - All transactions listed<br>- Correct balances |

## 4. Margin & Revenue Sharing Tests

| TC-ID | Test Scenario | Pre-requisites | Steps | Expected Result |
|-------|--------------|----------------|-------|-----------------|
| MARGIN-001 | Fixed allowance calculation | Staff with fixed allowance scheme | 1. Staff completes transaction | - Staff receives fixed fee<br>- Mitra receives remaining margin |
| MARGIN-002 | Margin share calculation | Staff with % margin scheme | 1. Staff completes transaction | - Staff receives configured % of margin<br>- Mitra receives remainder |
| MARGIN-003 | Default 60/40 split | Staff with margin share, no custom % | 1. Staff completes transaction | - Staff receives 60%<br>- Mitra receives 40% |
| MARGIN-004 | Custom margin percentage | Staff with custom % (e.g., 70/30) | 1. Staff completes transaction | - Staff receives 70%<br>- Mitra receives 30% |
| MARGIN-005 | Margin calculation accuracy | Transaction with known values | 1. Base: 10000<br>Platform: 11000<br>Mitra sell: 12000 | - Margin: 1000<br>- Staff/Mitra split correctly |
| MARGIN-006 | Zero margin scenario | Base = Mitra selling price | 1. Transaction with no margin | - Staff receives allowance only (if fixed)<br>- No revenue sharing |
| MARGIN-007 | Commission accumulation | Multiple transactions | 1. Staff completes 5 transactions | - Total commission accumulated<br>- Balance updated |

## 5. Transaction Flow Tests

### 5.1 Prepaid Transactions

| TC-ID | Test Scenario | Pre-requisites | Steps | Expected Result |
|-------|--------------|----------------|-------|-----------------|
| TRANS-001 | Successful pulsa transaction | Sufficient balance | 1. Select pulsa<br>2. Enter phone<br>3. Confirm<br>4. Enter PIN | - Transaction successful<br>- Balance deducted<br>- Receipt generated |
| TRANS-002 | Successful data package | Sufficient balance | 1. Select data package<br>2. Enter phone<br>3. Confirm<br>4. Enter PIN | - Transaction successful<br>- Balance deducted |
| TRANS-003 | Failed transaction - invalid number | Valid balance | 1. Enter invalid phone number | - Error: "Invalid phone number"<br>- No deduction |
| TRANS-004 | Failed transaction - insufficient balance | Low balance | 1. Transaction exceeds balance | - Error: "Insufficient balance"<br>- No deduction |
| TRANS-005 | Failed transaction - wrong PIN | Active transaction | 1. Enter wrong PIN 3 times | - Error: "Invalid PIN"<br>- Transaction cancelled |
| TRANS-006 | Transaction pending status | Provider delay | 1. Transaction submitted | - Status: pending<br>- Callback received when complete |

### 5.2 Postpaid Transactions

| TC-ID | Test Scenario | Pre-requisites | Steps | Expected Result |
|-------|--------------|----------------|-------|-----------------|
| TRANS-007 | PLN postpaid inquiry | Valid PLN number | 1. Enter customer ID<br>2. Check bill | - Bill details displayed<br>- Amount shown |
| TRANS-008 | PLN postpaid payment | Bill retrieved | 1. Confirm payment<br>2. Enter PIN | - Payment successful<br>- Bill marked paid |
| TRANS-009 | PDAM payment | Valid customer ID | 1. Enter customer ID<br>2. Check bill<br>3. Pay | - Payment processed<br>- Receipt generated |

### 5.3 E-wallet Transactions

| TC-ID | Test Scenario | Pre-requisites | Steps | Expected Result |
|-------|--------------|----------------|-------|-----------------|
| TRANS-010 | OVO top-up | Sufficient balance | 1. Select OVO<br>2. Enter OVO number<br>3. Enter amount<br>4. Confirm | - Top-up successful<br>- Balance deducted |
| TRANS-011 | GoPay top-up | Sufficient balance | 1. Select GoPay<br>2. Enter GoPay number<br>3. Enter amount<br>4. Confirm | - Top-up successful |
| TRANS-012 | Dana top-up | Sufficient balance | 1. Select Dana<br>2. Enter phone<br>3. Enter amount<br>4. Confirm | - Top-up successful |

## 6. Reporting & Analytics Tests

| TC-ID | Test Scenario | Pre-requisites | Steps | Expected Result |
|-------|--------------|----------------|-------|-----------------|
| REPORT-001 | Mitra sales report | Transactions completed | 1. Filter by day/week/month | - Correct sales data displayed |
| REPORT-002 | Mitra profit analysis | Transactions with margin | 1. View profit report | - Profit calculated correctly |
| REPORT-003 | Staff performance report | Staff transactions | 1. View staff performance | - Transaction count displayed<br>- Commission earned shown |
| REPORT-004 | Staff transaction history | Staff completed transactions | 1. View history | - All transactions listed<br>- Correct details |
| REPORT-005 | Staff commission report | Staff with earnings | 1. View commission | - Total commission displayed<br>- Per-transaction breakdown |
| REPORT-006 | Date filter functionality | Historical data | 1. Apply date filters | - Data filtered correctly |

## 7. Product & Service Tests

| TC-ID | Test Scenario | Pre-requisites | Steps | Expected Result |
|-------|--------------|----------------|-------|-----------------|
| PROD-001 | Product list display | Active products | 1. View product list | - All products displayed<br>- Categories grouped |
| PROD-002 | Product search | Products exist | 1. Search by keyword | - Matching products shown |
| PROD-003 | Product detail view | Product selected | 1. Tap product | - Details displayed<br>- Price shown |
| PROD-004 | Product price accuracy | Known product | 1. Check product price | - Correct margin applied<br>- Accurate Mitra price |
| PROD-005 | Hourly prepaid sync | Scheduled job | 1. Check product update time | - Prepaid products updated |
| PROD-006 | 5-min postpaid sync | Scheduled job | 1. Check postpaid update | - Postpaid products refreshed every 5 min |

## 8. Security & Fraud Detection Tests

| TC-ID | Test Scenario | Pre-requisites | Steps | Expected Result |
|-------|--------------|----------------|-------|-----------------|
| SEC-001 | Daily transaction limit | Limit set for staff | 1. Exceed staff daily limit | - Transaction blocked<br>- Error: "Daily limit exceeded" |
| SEC-002 | PIN attempt limit | Wrong PIN attempts | 1. Enter wrong PIN 5 times | - Account locked<br>- OTP required to unlock |
| SEC-003 | Session timeout | Inactive session | 1. App inactive for 30 min | - Session expired<br>- Re-authentication required |
| SEC-004 | Concurrent session limit | Multiple logins | 1. Login from multiple devices | - Previous session invalidated |
| SEC-005 | Audit log creation | Transaction performed | 1. Check audit logs | - All actions logged<br>- Timestamp accurate |

## 9. Error & Edge Cases

| TC-ID | Test Scenario | Pre-requisites | Steps | Expected Result |
|-------|--------------|----------------|-------|-----------------|
| EDGE-001 | Network timeout | Poor connection | 1. Attempt transaction | - Retry mechanism<br>- User notified |
| EDGE-002 | Digiflazz API error | API down | 1. Submit transaction | - Error handled gracefully<br>- Retry queue |
| EDGE-003 | Database connection error | DB down | 1. Query data | - Fallback/Circuit breaker<br>- User notified |
| EDGE-004 | Duplicate transaction | Same request twice | 1. Submit same transaction twice | - Deduplication<br>- Second request rejected |
| EDGE-005 | Invalid input handling | Malformed data | 1. Enter special characters | - Input validation<br>- Error: "Invalid input" |

## 10. API Endpoint Tests

### Auth Service (Port 8081)

| Endpoint | Method | Test Case | Expected |
|----------|--------|-----------|----------|
| /api/v1/auth/register | POST | Valid registration | 201 Created |
| /api/v1/auth/register | POST | Duplicate phone | 409 Conflict |
| /api/v1/auth/login | POST | Valid credentials | 200 OK + token |
| /api/v1/auth/login | POST | Invalid credentials | 401 Unauthorized |
| /api/v1/auth/verify-otp | POST | Valid OTP | 200 OK |
| /api/v1/auth/refresh | POST | Valid refresh token | 200 OK + new token |
| /api/v1/auth/logout | POST | Authenticated | 200 OK |

### Integration Service (Port 8086)

| Endpoint | Method | Test Case | Expected |
|----------|--------|-----------|----------|
| /api/v1/integrations/digiflazz/transaction | POST | Valid request | 200 OK |
| /api/v1/integrations/providers | GET | Authenticated | 200 OK |
| /webhook/digiflazz | POST | Valid callback | 200 OK |
| /health | GET | Service running | 200 OK |
| /health/live | GET | Service running | 200 OK |
| /health/ready | GET | All deps ready | 200 OK |

## Test Execution Matrix

| Module | Priority | Manual | Automated | Status |
|--------|----------|--------|-----------|--------|
| Authentication | High | ✓ | | Pending |
| Role Management | High | ✓ | | Pending |
| Wallet Management | High | ✓ | | Pending |
| Transactions | High | ✓ | | Pending |
| Margin Calculation | High | ✓ | | Pending |
| Reporting | Medium | ✓ | | Pending |
| Error Handling | Medium | ✓ | | Pending |

## Entry/Exit Criteria

### Entry Criteria
- Test environment ready
- Test data prepared
- API services deployed
- Test accounts created

### Exit Criteria
- All High priority tests passed (100%)
- All Medium priority tests passed (>95%)
- All Low priority tests passed (>90%)
- Critical bugs fixed and retested