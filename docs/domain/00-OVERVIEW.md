# 🗺️ PPOB Business Domains Overview

## 1. Introduction
The PPOB system is organized into several distinct business domains, each representing a specific area of functionality and expertise. This document outlines the boundaries and core responsibilities of these domains.

## 2. Domain Map
The business is logically divided into six core areas:

### 2.1 Identity & Access
Focuses on securing the platform and identifying users.
- **Concepts:** Users, Devices, Trust Scores, MFA (OTP), JWT Sessions.
- **Goal:** Provide seamless access while protecting against unauthorized use.

### 2.2 User & Staff Management
Focuses on the relationships and hierarchy between business entities.
- **Concepts:** Roles (Mitra, Staff), Hierarchy (Assigned By), Permissions.
- **Goal:** Enable multi-tenant business operations.

### 2.3 Wallet & Financial Engine
The "General Ledger" of the system.
- **Concepts:** Wallets, Balance (Available/Held), Events, Commissions, Profit.
- **Goal:** Ensure 100% financial accuracy and auditability.

### 2.4 Product & Pricing Catalog
Manages what is being sold and for how much.
- **Concepts:** Categories, SKUs, Platform Price, Selling Price, Markup.
- **Goal:** Provide a competitive and profitable digital product catalog.

### 2.5 Transaction & Operations
The core revenue-generating domain.
- **Concepts:** Orders, Statuses (State Machine), Idempotency, Receipts.
- **Goal:** Facilitate smooth and reliable fulfillment of digital purchases.

### 2.6 Provider Integration
The bridge to external supply chains.
- **Concepts:** Provider API, Webhooks, Error Mapping, Circuit Breakers.
- **Goal:** Abstract external complexities from the core business logic.

## 3. Documentation Index
Detailed domain specifications are available in the following documents:

1.  **[Identity & Access](01-IDENTITY.md):** Security and authentication domain.
2.  **[User & Staff Management](02-USER-STAFF.md):** The Partner-Staff ecosystem.
3.  **[Wallet & Financials](03-WALLET-FINANCIAL.md):** Accounting and ledger logic.
4.  **[Product & Catalog](04-PRODUCT-CATALOG.md):** Digital products and pricing tiers.
5.  **[Transaction Operations](05-TRANSACTION.md):** Sales and lifecycle management.
6.  **[Integration Boundaries](06-INTEGRATION.md):** External provider interactions.
