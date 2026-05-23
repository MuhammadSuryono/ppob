# 📦 Domain: Product & Catalog

## 1. Core Mission
This domain manages the digital product inventory, categorizing them for ease of use and establishing the multi-tiered pricing structure.

## 2. Key Concepts

### 2.1 Categories
Grouping of products for UI presentation (e.g., Pulsa, PLN, PDAM).
- **Metadata:** Each category defines specific input requirements (e.g., Phone Number vs. Meter ID).

### 2.2 Products (SKUs)
Individual sellable items (e.g., "Indosat 10k", "PLN Prepaid 50k").
- **State:** Active/Inactive.

### 2.3 Pricing Tiers
- **Base Price:** The price provided by the supplier (Digiflazz).
- **Platform Price:** Base Price + Platform Markup (revenue for the platform owner).
- **Selling Price:** The price the end customer pays.

## 3. Business Rules

### 3.1 Markup Policy
- **Platform Markup:** A global percentage (default 5%) applied by the system to all base prices.
- **Selling Price Guard:** Mitra cannot set a Selling Price that is lower than the Platform Price.

### 3.2 Product Lifecycle
- **Supplier Status:** If the supplier marks a product as "inactive" or "out of stock," it must be immediately hidden from the catalog.
- **Sync Priority:** Postpaid products are synchronized more frequently (every 5 mins) than prepaid (hourly) due to higher status volatility.

## 4. Domain Logic

### 4.1 Tiered Pricing Formula
1. `Platform Price = Base Price * (1 + 0.05)`
2. `Margin = Selling Price - Platform Price`
3. `Selling Price` is either the Mitra's custom price or the `Platform Price` by default.

### 4.2 Distributed Catalog
The catalog is synchronized from external providers via the Integration Service, but the business logic for pricing and categorization remains within the Product Catalog domain.
