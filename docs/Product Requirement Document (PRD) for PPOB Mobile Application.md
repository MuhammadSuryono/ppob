# Product Requirement Document (PRD) for PPOB Mobile Application

## 1. Introduction

This document outlines the requirements for a modern multi-tenant Payment Point Online Bank (PPOB) mobile application. The application aims to provide a seamless experience for Mitra (partners) to manage their staff and facilitate digital service transactions such as mobile credit, data packages, PLN, and other digital products.

## 2. Goals

- Develop a modern PPOB mobile application with a bright and fresh UI/UX.
- Implement a multi-tenant system where a Mitra can manage multiple staff members.
- Enable staff to perform digital product transactions.
- Provide robust authentication and authorization mechanisms.
- Facilitate margin and revenue sharing between Mitra and staff.
- Offer comprehensive reporting and analytics for both Mitra and staff.

## 3. Features

### 3.1. Authentication & Authorization

- **Login/Register with Phone Number:** Users can register and log in using only their phone number, and will provide their name during registration.
- **OTP Verification:** OTP will be used for new registrations and untrusted device logins.
- **Password Creation:** New users will create a password after OTP verification.
- **PIN for Transaction Authorization:** A 6-digit PIN will be used to authorize all transactions, separate from the login password.
- **Trusted Device Recognition:** For trusted devices, users can directly input their PIN to log in.

### 3.2. Role & Multi-Account

- **Mitra (Partner) Role:** The default user role with privileges to:
    - Add new staff members.
    - Configure staff roles and permissions.
    - Top-up staff wallets.
- **Staff Role:** Staff members can perform transactions using funds from their linked Mitra's wallet.
- **Multi-Role Capability:** A single user can have multiple roles (e.g., Mitra and Staff for another Mitra).
- **Role Switching:** Users can switch between their active roles after logging in (multi-tenant experience).

### 3.3. Wallet & Balance Management

- **Mitra Main Wallet:** Each Mitra will have a primary wallet for their funds.
- **Staff Sub-Wallets:** Staff members will have sub-wallets, topped up by their respective Mitra.
- **Transaction Deduction:** All transactions will deduct funds from the active role's wallet.

### 3.4. Margin & Revenue Sharing

- **Tiered Pricing:**
    - **Base Price (Digiflazz):** The initial price provided by the Digiflazz API.
    - **Platform Markup:** An additional markup applied by the platform to the base price, forming the **Platform Price** for Mitra. This markup represents the platform's revenue.
    - **Mitra Selling Price:** Mitra can set their own selling price to end-users. This price is crucial for calculating Mitra's profit and staff's margin share.
- **Revenue Sharing Schemes (Configurable by Mitra per Staff):**
    - When setting up staff, the default scheme is **Fixed Allowance**.
    - **Fixed Allowance:** Staff receives a fixed fee per transaction.
    - **Margin Share (%):** If selected, the default split is 60% for Staff and 40% for Mitra (60/40), applicable across all products. This percentage is configurable by the Mitra. Staff receives a percentage of the margin (Mitra Selling Price - Platform Price).

### 3.5. Products & Services

- **Categories:** Support for various digital product categories, including:
    - Pulsa (Mobile Credit)
    - Paket Data (Data Packages)
    - PLN (Prepaid & Postpaid Electricity)
    - PDAM (Water Bills)
    - E-wallet Top-ups
    - etc.
- **Platform Management:** The platform will manage product categories and master product lists.
- **Digiflazz Synchronization:**
    - Prepaid products: Synchronized hourly.
    - Postpaid products: Synchronized every 5 minutes.

### 3.6. Reporting & Analytics

- **Mitra Reports:**
    - Total sales.
    - Profit analysis.
    - Staff performance.
- **Staff Reports:**
    - Transaction history.
    - Earned commission.
- **Filters:** Reports can be filtered by daily, weekly, or monthly periods.

### 3.7. Transaction Flow

- **Step-by-step Process:**
    1. Select product category.
    2. Input target number/ID.
    3. Validation (e.g., check postpaid bill details).
    4. Confirmation screen.
    5. Input PIN for authorization.
    6. Transaction status display (success/pending/failed).

## 4. UI/UX Guidelines

- **Style:** Clean, modern, bright UI with vibrant colors (green/blue/soft gradients).
- **Inspiration:** Mitra Bukalapak, OVO Merchant, GoBiz.
- **Character:** Friendly yet professional, focusing on usability for semi-technical users.
- **Homepage Elements:**
    - Grid of PPOB categories with large, attractive icons.
    - Highlighted balance display at the top.
    - Quick access to recent transactions.

## 5. System Design (High-Level)

### 5.1. Backend Architecture

- **Microservices:**
    - **Auth Service:** Handles user authentication (login, registration, OTP).
    - **User & Role Service:** Manages user profiles, roles, and staff assignments.
    - **Wallet Service:** Manages Mitra and staff wallets, balance deductions, and top-ups.
    - **Transaction Service:** Processes all digital product transactions.
    - **Product Service:** Manages product categories, master products, and pricing.
    - **Integration Service (Digiflazz):** Handles communication with the Digiflazz API (product sync, transaction hits, webhook callbacks).
- **Technologies:**
    - **Backend Language:** Golang
    - **Caching:** Redis
    - **Database:** PostgreSQL

### 5.2. Key Flows

- OTP flow.
- PIN validation.
- Wallet deduction.
- Margin calculation.

## 6. Integrations

- **Digiflazz API:**
    - Product synchronization (scheduled).
    - Transaction initiation.
    - Webhook handling for transaction status updates.

## 7. Extra Features (Real Product Enhancements)

- **Simple Fraud Detection:** Transaction limits for staff members.
- **Notifications:** Push notifications and in-app notifications.
- **Audit Log:** Comprehensive logging of all transactions.

## 8. Expected Outputs

- Complete PRD (this document).
- User flow diagrams.
- Wireframes (low & high fidelity).
- Design system (color, typography, spacing).
- Database schema.
- API contract (REST).
- System architecture diagram.
- Edge cases & security concerns documentation.

## 9. References

[1] Digiflazz API Documentation: [https://developer.digiflazz.com/api/](https://developer.digiflazz.com/api/)
[2] UI/UX Study case — Optimizing/redesigning Gobiz Dashboard: [https://miqbalmahulana.medium.com/ui-ux-study-case-optimizing-redesigning-gobiz-dashboard-by-ui-audit-method-d3b6c98f4013](https://miqbalmahulana.medium.com/ui-ux-study-case-optimizing-redesigning-gobiz-dashboard-by-ui-audit-method-d3b6c98f4013)
[3] Redesign Bukalapak App (E-commerce) | UI/UX Case Study: [https://www.behance.net/gallery/115689141/Redesign-Bukalapak-App-(E-commerce)-UIUX-Case-Study/modules/660106397](https://www.behance.net/gallery/115689141/Redesign-Bukalapak-App-(E-commerce)-UIUX-Case-Study/modules/660106397)
[4] OVO APP — UI/UX Design Case Study: [https://medium.com/@boyferton07/ovo-app-ui-ux-design-case-study-40ef0e34e9d9](https://medium.com/@boyferton07/ovo-app-ui-ux-design-case-study-40ef0e34e9d9)
[5] Building a Multi-Tenant Architecture in Golang: A Practical Guide: [https://articles.wesionary.team/building-a-multi-tenant-architecture-in-golang-a-practical-guide-8ee066436678](https://articles.wesionary.team/building-a-multi-tenant-architecture-in-golang-a-practical-guide-8ee066436678)
