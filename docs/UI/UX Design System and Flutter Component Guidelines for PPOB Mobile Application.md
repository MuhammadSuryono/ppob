# UI/UX Design System and Flutter Component Guidelines for PPOB Mobile Application (Updated)

This document outlines the UI/UX design system and provides guidelines for Flutter component implementation, ensuring a modern, clean, and bright aesthetic consistent with leading Indonesian fintech applications (Mitra Bukalapak, OVO Merchant, GoBiz). Focus: usability for Mitra and Staff with high daily transaction volumes.

**Latest Updates:**
- Added component library inventory
- Defined navigation pattern (bottom nav + drawer)
- Added empty & error state designs
- Added loading states (skeletons)
- Added accessibility requirements
- Added dashboard metrics layout
- Added success/modal patterns

---

## 1. Design Principles

- **Clarity:** Information easily understandable at a glance
- **Efficiency:** Minimal steps to complete tasks
- **Consistency:** Predictable interactions and visual elements
- **Accessibility:** Meets WCAG 2.1 AA (contrast 4.5:1 minimum)
- **Modern & Bright:** Fresh palette, avoids generic AI template look

---

## 2. Color Palette

| Name | Hex Code | Usage |
|---|---|---|
| Primary Green | `#4CAF50` | Primary action buttons, highlights |
| Secondary Blue | `#2196F3` | Secondary actions, info badges |
| Accent Orange | `#FF9800` | Warnings, alerts, special promos |
| Background Light | `#F5F5F5` | Main background (light gray) |
| Surface White | `#FFFFFF` | Cards, inputs, modals |
| Text Dark | `#212121` | Primary text, headings |
| Text Medium | `#757575` | Secondary text, hints |
| Text Light | `#BDBDBD` | Disabled, placeholders |
| Success | `#4CAF50` | Success states |
| Warning | `#FF9800` | Warnings |
| Error | `#F44336` | Error states |

---

## 3. Typography

**Font:** Roboto (Google Fonts — included in Flutter)

| Element | Font Family | Weight | Size (sp) | Color | Line Height |
|---|---|---|---|---|---|
| Display Large | Roboto | Bold | 32 | Text Dark | 40 |
| Headline Medium | Roboto | Medium | 24 | Text Dark | 32 |
| Title Large | Roboto | Medium | 20 | Text Dark | 28 |
| Body Large | Roboto | Regular | 16 | Text Dark | 24 |
| Body Medium | Roboto | Regular | 14 | Text Medium | 20 |
| Label Small | Roboto | Medium | 12 | Text Light | 16 |
| Button Text | Roboto | SemiBold | 14 | White / Primary | 20 |

**Use `TextTheme` from `Theme.of(context).textTheme` everywhere; override as needed.**

---

## 4. Spacing & Layout

- **Grid:** 8dp grid system
- **Padding:** Small (8dp), Medium (16dp), Large (24dp)
- **Margin:** Same as padding
- **Corner Radius:** 8dp (cards, buttons), 4dp (small chips), 16dp (modals, large cards)

**Screen Margins:** 16dp on all sides; content within safe area.

---

## 5. Iconography

- **Style:** Material Icons (filled for active, outlined for inactive)
- **Size:** 24dp standard, 48dp for category icons
- **Color:** Follow text color or use primary green for active

**Custom icons:** Use SVG via `flutter_svg` package; store in `assets/icons/`.

---

## 6. Navigation Pattern

### 6.1 Bottom Navigation Bar

Fixed at bottom, 5 items:

| Icon | Label | Route |
|---|---|---|
| `home` | Beranda | `/` |
| `receipt_long` | Transaksi | `/transactions` |
| `account_balance_wallet` | Wallet | `/wallet` |
| `people` (Mitra only) | Staff | `/staff` |
| `person` | Profil | `/profile` |

**Behavior:**
- Tapping icon navigates to corresponding screen
- Active state: icon filled, primary green color
- Inactive: outlined or gray
- Show badge on `Wallet` if `balance_held > 0` (pending transactions)
- Show badge on `Staff` (Mitra only) if `pending_staff_requests > 0` (new staff awaiting approval)

### 6.2 Drawer Menu (Profile Sidebar)

Access via profile screen top-left avatar or swipe.

Items:
- Settings (`settings`)
- Ganti PIN (`lock`)
- Perangkat Terpercaya (`devices`)
- Bantuan & Pendorongan (`help`)
- Keluar (`logout`)

---

## 7. Component Library

### 7.1 Buttons

| Variant | Usage | Properties |
|---|---|---|
| `PPOBButton.primary` | Main CTA (e.g., "Bayar", "Kirim") | `onPressed`, `label`, `loading?` |
| `PPOBButton.secondary` | Secondary actions (e.g., "Batal") | `onPressed`, `label` |
| `PPOBButton.text` | Tertiary (e.g., "Batal", "Tidak") | `onPressed`, `label` |
| `PPOBButton.danger` | Destructive (e.g., "Hapus akun") | `onPressed`, `label` |

**States:** `enabled`, `disabled` (opacity 0.5), `loading` (spinner inside)

**Code Example:**
```dart
class PPOBButton extends StatelessWidget {
  final String label;
  final VoidCallback? onPressed;
  final bool isLoading;
  final PPOBButtonVariant variant;
  
  const PPOBButton.primary({required this.label, this.onPressed, this.isLoading = false});
  
  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    Color bgColor;
    switch (variant) {
      case primary: bgColor = theme.primaryColor; break;
      case secondary: bgColor = theme.colorScheme.secondary; break;
      case danger: bgColor = theme.colorScheme.error; break;
      default: bgColor = Colors.white;
    }
    
    return ElevatedButton(
      onPressed: isLoading ? null : onPressed,
      style: ElevatedButton.styleFrom(
        backgroundColor: bgColor,
        foregroundColor: Colors.white,
        padding: const EdgeInsets.symmetric(vertical: 16),
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(8)),
        minimumSize: const Size.fromHeight(52), // touch target
      ),
      child: isLoading
        ? const SizedBox(width: 20, height:20, child: CircularProgressIndicator(strokeWidth:2, color:Colors.white))
        : Text(label, style: theme.textTheme.labelLarge),
    );
  }
}
```

---

### 7.2 Input Fields

**Variants:** Text, Password, PIN, Phone, Numeric.

**States:** Default, Focused, Error, Disabled

**Example PIN Input:**
```dart
class PinInput extends StatelessWidget {
  final List<TextEditingController> controllers;
  final FocusNode focusNode;
  final int length;
  
  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceEvenly,
      children: List.generate(6, (i) {
        return SizedBox(
          width: 48, height: 56,
          child: TextField(
            controller: controllers[i],
            focusNode: i==0 ? focusNode : null,
            textAlign: TextAlign.center,
            keyboardType: TextInputType.number,
            obscureText: true,
            maxLength: 1,
            decoration: InputDecoration(
              counterText: '', // hide counter
              border: OutlineInputBorder(borderRadius: BorderRadius.circular(8)),
            ),
          ),
        );
      }),
    );
  }
}
```

---

### 7.3 Cards

| Card Type | Usage |
|---|---|
| `TransactionCard` | Shows single transaction in list (product icon, name, amount, status badge) |
| `ProductCard` | Grid item: product image/icon, name, price |
| `CategoryCard` | Large icon + label for category selection |
| `WalletCard` | Balance display (available + held) |
| `StaffCard` | Staff info, balance, daily stats, actions |

**Example WalletCard:**
```dart
class WalletCard extends StatelessWidget {
  final Decimal available;
  final Decimal held;
  
  @override
  Widget build(BuildContext context) {
    final formatter = NumberFormat.currency(locale: 'id_ID', symbol: 'Rp ');
    return Card(
      elevation: 2,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text('Saldo Tersedia', style: Theme.of(context).textTheme.bodyMedium),
            Text(formatter.format(available), 
                 style: Theme.of(context).textTheme.headlineMedium?.copyWith(color: Colors.green)),
            const SizedBox(height:12),
            Text('Dikunci (Pending)', style: Theme.of(context).textTheme.bodySmall),
            Text(formatter.format(held), style: Theme.of(context).textTheme.titleLarge),
          ],
        ),
      ),
    );
  }
}
```

---

### 7.4 Bottom Sheets

**ModalBottomSheet** for:
- Staff creation/edit (margin settings)
- Top-up confirmation
- Transaction cancellation reason

Use `showModalBottomSheet` with `isScrollControlled: true` for full-height on small screens.

---

### 7.5 Loading States

**Skeleton Loader:** Shimmer effect while data fetching.

```dart
class TransactionListSkeleton extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return ListView.builder(
      itemCount: 5,
      itemBuilder: (_, i) => const ListTile(
        leading: CircleAvatar(backgroundColor: Colors.grey),
        title: SizedBox(width: 100, child: LinearProgressIndicator()),
        subtitle: SizedBox(width: 150, child: LinearProgressIndicator()),
        trailing: SizedBox(width: 80, child: LinearProgressIndicator()),
      ),
    );
  }
}
```

**Full-screen loading:** Center spinner with optional text "Memuat..."

---

## 8. Empty State Designs

| Screen | Illustration (icon) | Message | CTA |
|---|---|---|---|
| Transactions | `receipt_long` (gray) | "Belum ada transaksi" | "Lakukan transaksi pertama" (navigates to home) |
| Staff | `person_add_disabled` | "Belum ada staff" | "Tambah Staff" (FAB) |
| Wallet (zero) | `account_balance_wallet` (gray) | "Saldo Rp 0" | "Top Up" (if enabled) |
| No Network | `wifi_off` | "Tidak ada koneksi internet" | "Coba lagi" (retry) |
| Search empty | `search_off` | "Tidak ditemukan" | — |

**Design:** Centered icon (48dp), message (body medium), optional button below.

---

## 9. Error State Screens

### 9.1 Network Error (Full-screen)
- Illustration: Wi-Fi with slash
- Title: "Tidak dapat terhubung"
- Description: "Periksa koneksi internet dan coba lagi"
- Button: "Coba lagi" (primary, reloads)

### 9.2 Server Error (500)
- Illustration: Server with error
- Title: "Server sedang sibuk"
- Description: "Tim kami sudah diberitahu. Coba dalam beberapa menit."
- Button: "Kembali ke Beranda"

### 9.3 Validation Errors
- Inline below input field
- Small red text (`caption` style)
- Example: "Nomor HP harus dimulai dengan +62"

---

## 10. Success/Confirmation Modals

**Transaction Success Modal:**

```dart
showDialog(
  context: context,
  builder: (_) => AlertDialog(
    icon: Icon(Icons.check_circle, color: Colors.green, size: 64),
    title: Text('Transaksi Berhasil'),
    content: Column(
      mainAxisSize: MainAxisSize.min,
      children: [
        Text('Rp 27.000', style: Theme.of(context).textTheme.headlineSmall),
        Text('XL Data 25GB', style: Theme.of(context).textTheme.bodyMedium),
        Text('081234567890', style: Theme.of(context).textTheme.bodySmall),
        const SizedBox(height: 8),
        Text('Ref: ABC123XYZ', style: Theme.of(context).textTheme.labelSmall),
      ],
    ),
    actions: [
      TextButton(onPressed: () => Navigator.pop(context), child: Text('Selesai')),
      TextButton(onPressed: () => shareReceipt(tx), child: Text('Bagikan')),
    ],
  ),
);
```

---

## 11. Dashboard Metrics Layout

**Home Screen (Mitra):**

```
[Header: Greeting + Notification Bell]
[Wallet Card: Available Rp X, Held Rp Y]  [Top Up Button]
[Quick Actions Grid: 6 categories (Pulsa, Data, PLN, PDAM, E-Wallet, ...)]
[Recent Transactions: horizontal scroll list (3 items)]
[Staff Quick Stats: if Mitra — "5 staff, Rp 1.2M sales today" ]
```

**Report Dashboard (Mitra):**
```
[Date Range Picker: Today | 7d | 30d | Custom]
[KPIs Row: 2×2 grid]
  Card 1: Total Sales (Rp X)
  Card 2: Platform Profit (Rp Y)
  Card 3: Staff Count (N)
  Card 4: Success Rate (%)
[Chart: Line chart — last 30 days sales trend]
[Chart: Bar chart — top 5 staff by transactions]
[Table: Recent transactions with pagination]
```

---

## 12. Screen Wireframe Descriptions (Key Screens)

### 12.1 Login / Register Flow

**Adaptive Auth Flow Overview:**

| Screen | Condition | Navigation |
|---|---|---|
| Screen 1: Phone Input | Always first | → Screen 2 or Screen 4 |
| Screen 2: OTP | `requires_otp == true` | → Screen 3a (new) or Screen 3b (existing) |
| Screen 3a: Create PIN & Password | `is_new_user == true` (after OTP) | → Home |
| Screen 3b: Input PIN/Password | `is_new_user == false` (after OTP) | → Home |
| Screen 4: PIN Login | `is_trusted == true` (skip OTP) | → Home |

---

**Screen 1 — Phone Input:**
- Large heading "Masuk / Daftar"
- Phone field with country code picker (pre-set +62)
- "Lanjutkan" button (primary)
- Link: "Sudah punya akun? Masuk"
- **API Call:** `POST /api/v1/auth/initiate`
  - Request: `{ phone, device_id, fingerprint }`
  - Response: `{ is_registered, is_trusted, requires_otp, user_id }`
  - If `is_trusted=true` → jump to **Screen 4 (PIN Login)**
  - If `requires_otp=true` → go to **Screen 2 (OTP)**

**Screen 2 — OTP:**
- "Masukkan kode OTP"
- 6-digit OTP input (auto-focus, auto-submit on 6th digit)
- Timer: "Kirim ulang dalam MM:SS" (calculated from `expires_at` in SendOTP response)
- "Kirim ulang" link (after timer)
- **API Call (request OTP):** `POST /api/v1/auth/send-otp`
  - Request: `{ phone, type: "login" or "register" }`
  - Response: `{ request_id, expires_at }`
- **API Call (verify OTP):** `POST /api/v1/auth/verify-otp`
  - Request: `{ request_id, phone, code, type }`
  - Response: `{ request_id, is_verified, is_new_user }`
  - If `is_new_user=true` → go to **Screen 3a (Create PIN & Password)**
  - If `is_new_user=false` → go to **Screen 3b (Input PIN or Password)**

**Screen 3a — Set PIN & Password (new user — after OTP verified):**
- Password field (show/hide toggle)
- Confirm Password
- 6-digit PIN pad (numeric grid 3×3)
- Confirm PIN
- "Buat Akun" button (CTA)
- **API Call:** `POST /api/v1/auth/register`
  - Request: `{ email, phone, full_name, password, pin, device_id, request_id }`
  - Response: `{ user_id, email, phone, full_name, token, refresh_token, expires_at, refresh_expires_at }`
  - Device automatically marked as trusted
  - Verified flag consumed (single-use)

**Screen 3b — Input PIN or Password (existing user — after OTP verified):**
- Default view: 6-digit PIN pad
- Switch toggle: "Gunakan password" / "Gunakan PIN"
- PIN input: 6-digit PIN pad (numeric grid 3×3)
- Password input: password field with show/hide toggle
- "Masuk" button (CTA)
- **API Call:** `POST /api/v1/auth/verify-credential`
  - Request: `{ phone, device_id, request_id, auth_method: "pin" or "password", value }`
  - Response: `{ user_id, email, phone, full_name, token, refresh_token, expires_at }`
  - Device automatically marked as trusted after success
  - Verified flag consumed (single-use)

**Screen 4 — PIN Login (trusted device — skip OTP):**
- Large 6-digit PIN pad
- Option: "Login dengan password" below (redirects to full password+OTP flow)
- **API Call:** `POST /api/v1/auth/verify-pin`
  - Request: `{ phone, pin, device_id }`
  - Response: `{ user_id, email, phone, full_name, token, refresh_token, expires_at }`

---

### 12.2 Transaction Flow

**Screen 1 — Category Selection:**
- Grid of category cards (icon + label)
- Search bar optional

**Screen 2 — Product Selection:**
- Search bar
- Sort: popularity, price low-high, price high-low
- List of ProductCard (name, platform price, Mitra selling price input? Not in list — set in settings)
- Input field for customer number at bottom (sticky)

**Screen 3 — Confirmation:**
- Product summary (name, icon)
- Customer number (masked partially)
- Price breakdown:
  - Harga platform: Rp 26,250
  - Margin Anda: Rp 750
  - Total penjualan: **Rp 27,000**
- "Bayar" button (primary)

**Screen 4 — PIN Authorization:**
- 6-digit PIN pad
- "Masukkan PIN untuk konfirmasi"

**Screen 5 — Result:**
- Success: Green checkmark, "Berhasil", amount, reference, "Selesai" & "Bagikan" buttons
- Pending: Yellow hourglass, "Sedang diproses...", "Lihat status" link
- Failed: Red cross, reason, "Coba lagi" button

---

### 12.3 Manage Staff (Mitra)

**Screen 1 — Staff List:**
- AppBar: "Kelola Staff" + "+ Tambah" button
- Search bar
- List: Avatar (initials), name, phone, today count/amount, wallet balance, status toggle (active/inactive)
- Tap row → detail/edit screen

**Screen 2 — Add/Edit Staff:**
- Form fields: Phone, Name, Password (if new), PIN
- Margin Settings section:
  - Radio: "Biaya Tetap ( Rp X per transaksi )" vs "Bagi Hasil ( Staff % )"
  - If Fixed: numeric input Rp amount
  - If Share: slider 10–80% (default 60)
- Daily Limit section: inputs for count and amount (optional overrides)
- Save button

**Screen 3 — Staff Top-Up Modal:**
- Select staff from dropdown/recent
- Amount input (numeric keypad)
- "Top Up" confirmation → shows Mitra wallet deduction

---

### 12.4 Wallet Screen

- Balance Card (large)
- Held amount section (collapsible details: list of pending transactions?)
- "Top Up" button (if enabled — currently only from Mitra main wallet through internal transfer? Feature TBD)
- Transaction history sub-tab (same as main Transactions screen but filtered to this wallet)

---

### 12.5 Reports (Mitra)

- Date range picker (presets: Hari ini, 7 hari terakhir, 1 bulan, custom)
- KPI cards row (2 columns × 2 rows)
- Line chart: daily sales last 30 days (using `fl_chart` package)
- Bar chart: staff performance (transaction count)
- Transaction list below charts (same as history but with all staff included)

---

## 13. Accessibility

- **Color Contrast:** All text/background combos pass WCAG AA (4.5:1). Verify with `flutter_contrast_checker`.
- **Touch Target:** Minimum 48×48 dp for buttons and tappable items.
- **Font Scaling:** Use `MediaQuery.textScaleFactor`; UI should accommodate up to 200% without overflow.
- **Semantic Labels:** Add `Semantics(label: '...')` for icon-only buttons (e.g., notification bell reads "Notifikasi").
- **Screen Reader:** Content order logical (use `Column` with `mainAxisSize.min` properly).

---

## 14. Internationalization (i18n)

**Locale:** Indonesian (`id_ID`) only for MVP.

**Number Formatting:**
```dart
final rpFormatter = NumberFormat.currency(locale: 'id_ID', symbol: 'Rp ', decimalDigits: 0);
String formatted = rpFormatter.format(amount); // Rp 1.250.000
```

**String Resources:** `lib/l10n/intl_id.arb` with keys:
- `"login_title": "Masuk / Daftar"`
- `"error_insufficient_balance": "Saldo tidak mencukupi"`
- ...

Generate code with `flutter gen-l10n`.

---

## 15. Animations & Microinteractions

- **Success checkmark:** Scale + fade-in 300ms
- **Error shake:** Horizontal shake 200ms on invalid input
- **Page transitions:** Platform-specific (Cupertino on iOS, Material on Android)
- **Loading:** `CircularProgressIndicator` with primary color

Keep animations subtle; avoid slowing user down.

---

## 16. Dark Mode (Future)

Not in initial scope, but design system should support `ThemeData.dark()` with adjusted colors:
- Background: `#121212`
- Surface: `#1E1E1E`
- Text: White on dark

---

## 17. Flutter Project Structure

```
lib/
  main.dart
  app.dart               // MaterialApp + router
  core/
    constants/           // colors, strings, endpoints
    errors/              // exceptions, error handler
    theme/               // theme_data.dart
    utils/               // formatters, validators
    widgets/             // reusable: ppob_button, ppob_input
  data/
    models/              // User, Transaction, Product (with fromJson)
    repositories/        // impl of abstract repo (auth_repo_impl.dart)
    datasources/         // local (hive), remote (dio)
    services/            // dio_client, hive_service, token_service
  presentation/
    pages/               // screens: home_screen, login_screen, ...
    providers/           // Riverpod providers
    widgets/             // screen-specific widgets
    routes/              // router configuration
l10n/
  intl_id.arb
assets/
  icons/
  images/
```

---

## 18. State Management (Riverpod) Example

```dart
final transactionProvider = StateNotifierProvider<TransactionNotifier, TransactionState>((ref) {
  return TransactionNotifier(ref.watch(transactionRepositoryProvider));
});

class TransactionNotifier extends StateNotifier<TransactionState> {
  final TransactionRepository _repo;
  
  TransactionNotifier(this._repo) : super(const TransactionState.initial());
  
  Future<void> initiate(InitiateRequest req) async {
    state = const TransactionState.loading();
    try {
      final tx = await _repo.initiate(req);
      state = TransactionState.success(tx);
    } on AppException catch (e) {
      state = TransactionState.error(message: e.message);
    }
  }
}
```

---

## 19. Testing UI

**Widget Tests:** Use `flutter_test` to verify:
- PIN input fills 6 fields correctly
- Login button disabled until all fields filled
- Error message displays on failed login
- Transaction success modal appears

**Golden Tests:** Capture screenshots for visual regression (different screen sizes, locales).

---

## 20. Performance Considerations

- `const` widgets wherever possible
- `ListView.builder` for long lists (transactions, products)
- `AutomaticKeepAliveClientMixin` for preserving state in tabs
- `RepaintBoundary` for complex animations that don't need to repaint often
- `precacheImage` for category icons on app start

---

## Appendix A — Common Widgets Code (Snippets)

**PPOBInput:**
```dart
class PPOBInput extends StatelessWidget {
  final TextEditingController controller;
  final String label;
  final String? errorText;
  final bool obscure;
  
  @override
  Widget build(BuildContext context) {
    return TextField(
      controller: controller,
      obscureText: obscure,
      decoration: InputDecoration(
        labelText: label,
        errorText: errorText,
        filled: true,
        fillColor: Colors.white,
        border: OutlineInputBorder(borderRadius: BorderRadius.circular(8)),
        contentPadding: EdgeInsets.symmetric(horizontal: 16, vertical: 12),
      ),
    );
  }
}
```

**Status Badge:**
```dart
Widget statusBadge(String status) {
  Color color;
  switch (status) {
    case 'Success': color = Colors.green; break;
    case 'Pending': color = Colors.orange; break;
    case 'Failed': case 'Expired': color = Colors.red; break;
    default: color = Colors.grey;
  }
  return Container(
    padding: EdgeInsets.symmetric(horizontal: 8, vertical: 4),
    decoration: BoxDecoration(color: color.withOpacity(0.1), borderRadius: BorderRadius.circular(12)),
    child: Text(status, style: TextStyle(color: color, fontSize: 12)),
  );
}
```

---

## Appendix B — Artboard Specifications

Use **Figma** or **Adobe XD** with 375×812 (iPhone 13) canvas for mobile design. Scale at 1×, 2×, 3× for different densities.

**Grid:** 8px
**Margins:** 16px left/right
**Safe Area:** Top notch, bottom home indicator

Export icons as SVG; raster images as WebP @ 2x for retina.

---

**Owner:** Design Team + Mobile Team Lead  
**Assets Location:** `assets/design/` (Figma files, exported icons)

---

**Related:**  
- `MOBILE_APP_SPEC.md` (state management, offline strategy)  
- `BUSINESS_LOGIC_SPEC.md` (field meanings, validation rules)  
- `ERROR_HANDLING.md` (error display UI guidelines)
