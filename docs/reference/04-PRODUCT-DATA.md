# 📦 Product Data Reference

## 1. Schema Definitions

### 1.1 Prepaid Product (SKU)
```json
{
  "product_name": "Axis 10.000",
  "category": "Pulsa",
  "brand": "AXIS",
  "type": "Umum",
  "price": 10893,
  "buyer_sku_code": "ax10",
  "buyer_product_status": true,
  "seller_product_status": true,
  "unlimited_stock": true,
  "stock": 0,
  "multi": true,
  "start_cut_off": "23:30",
  "end_cut_off": "0:30",
  "desc": "Reguler"
}
```

### 1.2 Category Metadata
Categories define UI input requirements.
- `PULSA`: Requires `Phone Number`.
- `PLN`: Requires `Customer ID` (Meter ID).
- `E-MONEY`: Requires `Phone Number`.

## 2. Sample Data: Core Brands

### 2.1 Telco Operators
- **TELKOMSEL:** SKU Prefix `s` (e.g., `s10`, `s50`).
- **XL:** SKU Prefix `x` (e.g., `x10`, `x100`).
- **AXIS:** SKU Prefix `ax`.
- **INDOSAT:** SKU Prefix `i`.

### 2.2 Billers (Postpaid)
- **PLN:** SKU `pln`.
- **PDAM:** Variable by region (e.g., `pdam_jakarta`).
- **BPJS:** SKU `bpjs`.

## 3. Synchronizaton Rules
- **Prepaid:** Full catalog refresh every 1 hour.
- **Postpaid:** Full catalog refresh every 5 minutes.
- **Price Calculation:** `Internal Platform Price = Digiflazz Price + 5% markup`.
