# PPOB Project Instructions

This file contains the foundational mandates, architecture rules, and workflows for the PPOB project. You **must** adhere to these rules at all times.

## 1. Documentation Precedence & Reference

The project documentation has been newly structured and centralized. When working on any task, you **MUST** consult the relevant documentation in the `docs/` folder before making architectural or logic decisions.

- **Architecture:** `docs/architecture/` (Microservices, Data, Security, Infrastructure, etc.)
- **Business Domains:** `docs/domain/` (Identity, Wallet, Transaction, Catalog, etc.)
- **Operational Flows:** `docs/flow/` (Auth, Transactions, Sync, Reconciliation, etc.)
- **Technical Reference:** `docs/reference/` (API Contracts, Digiflazz RC mapping, Error codes, Data models).

## 2. Documentation Maintenance

Documentation is a living part of the codebase.
- **Mandatory Updates:** If you make any code changes that alter business logic, API contracts, database schemas, or system flows, you **MUST** update the corresponding documentation files in the `docs/` folder to reflect these changes.
- Never leave the documentation out of sync with the implementation.

## 3. Subagent Verification

To ensure high-quality and consistent changes:
- Before finalizing a complex task (e.g., refactoring, adding a new feature, or complex bug fixes), you **MUST** spawn a subagent (e.g., `codebase_investigator` or `generalist`) to review the changes you have made.
- Ask the subagent to verify that the changes adhere to the project's architectural principles, security guidelines, and that they do not break existing flows defined in the `docs/`.

## 4. Build, Validate, and Commit Workflow

For every task that involves modifying code, follow this strict lifecycle:
1.  **Implement:** Write the code according to the rules and documentation.
2.  **Build & Test:** You **MUST** always run the build process (e.g., `make build`) and any relevant tests (e.g., `make test` or `go test ./...`) to ensure the code compiles and tests pass.
3.  **Commit:** If the build and tests are successful, and the subagent review (if applicable) is positive, you **MUST** stage the changes and commit them using `git add` and `git commit`. Use clear, descriptive commit messages.

## 6. Plan Tracking & Monitoring

All complex architectural changes or features must have an implementation plan.
- **Location:** All plans must be stored in the `docs/plan/` directory.
- **Completion:** Once a plan is fully implemented and verified, the filename **MUST** be appended with `_COMPLETED` (e.g., `feature-x.md` becomes `feature-x_COMPLETED.md`).
- This allows for easy auditing of implemented vs. pending architectural changes.
