# PPOB Development Backlog

This document tracks upcoming features, enhancements, and technical debt to be addressed in future development cycles.

## ✅ Completed / Done

1. **Async Commission Calculation**
   - **Status:** COMPLETED
   - **Implementation:** Moved commission processing to an asynchronous Redis-based worker. When a transaction status becomes `success`, a task is pushed to `q:commission_processing`. A background worker in `transaction-service` consumes these tasks and invokes `CommissionService` to calculate and credit earnings.

2. **Integration of Commission Logic into Transaction Flow (Wiring)**
   - **Status:** COMPLETED
   - **Implementation:** Updated `transaction-service`'s `UpdateTransactionStatus` and `UpdateTransactionStatusByProviderRef` methods to trigger the asynchronous commission task publishing. This ensures that every successful transaction (whether updated via API or Webhook) now automatically kicks off the profit distribution process.

## 📌 High Priority

*(Add new high priority items here)*

## 📅 Medium Priority

*(Add medium priority items here)*

## ⏳ Low Priority / Technical Debt

*(Add low priority items here)*
