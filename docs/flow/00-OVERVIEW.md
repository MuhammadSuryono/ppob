# 🔄 PPOB System Flows Overview

## 1. Introduction
This document provides a summary of the primary operational flows within the PPOB system. These flows define how users interact with the platform, how transactions are processed, and how data is synchronized across microservices.

## 2. Core User Flows
- **Authentication & Onboarding:** Adaptive flow covering registration, OTP verification, and device-trust based login.
- **Transaction Flow:** End-to-end lifecycle of digital product purchases, including secure PIN authorization and provider integration.
- **Staff Management:** Relationship management between Partners (Mitra) and their Staff, including role assignment and wallet funding.

## 3. Core System Flows
- **Financial Ledgering:** Event-sourced wallet operations including holds, debits, and commission distribution.
- **Product Synchronization:** Automated catalog updates from external providers (Digiflazz).
- **Reconciliation & Recovery:** Background processes to ensure data consistency and handle pending state timeouts.

## 4. Documentation Index
Refer to the following documents for detailed step-by-step specifications and sequence diagrams:

1.  **[Authentication & Registration](01-AUTH-REGISTRATION.md):** The adaptive login and new user onboarding flow.
2.  **[Transaction Lifecycle](02-TRANSACTION-LIFECYCLE.md):** Detailed flow for Prepaid and Postpaid transactions.
3.  **[Staff Management](03-STAFF-MANAGEMENT.md):** Managing the Mitra-Staff hierarchy and permissions.
4.  **[Wallet & Financial Flows](04-WALLET-OPERATIONS.md):** Top-ups, holds, and commission distribution logic.
5.  **[Product Sync Flow](05-PRODUCT-SYNC.md):** How the system keeps pricing and availability up to date.
6.  **[Reconciliation & Recovery](06-RECONCILIATION.md):** Handling system drifts and transaction timeouts.
