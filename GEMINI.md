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

## 5. Database Migration Workflow

To maintain schema consistency across environments, you **MUST** follow these rules for any database changes:
- **Migration Files:** Every schema change (ADD/ALTER/DROP) must be defined in a new SQL migration file within the `migrations/` directory using the naming convention `0XX_description.sql`.

## 7. Environment Specific Configuration

To ensure smooth development across different operating systems:
- **Mobile App Gradle:** If operating on a Linux system, you **MUST** ensure the Linux-specific Java home path (`org.gradle.java.home=/usr/lib/jvm/java-21-openjdk`) is enabled in `mobile-app/gradle.properties` during build/test tasks. However, **DO NOT** commit these local environment changes to the repository; always rollback to the default Windows path before pushing.
- **Goose Synchronization:** You **MUST** run migrations officially using the `goose` tool to ensure the `goose_db_version` table is updated.
  - **Docker Workflow:** Use `docker run --rm -v $(pwd)/migrations:/migrations --network <internal_network> -e GOOSE_DRIVER=postgres -e GOOSE_DBSTRING="..." gomicro/goose:latest goose -dir /migrations up`.
- **Model Sync:** When a database column is added or modified, you **MUST** immediately update the corresponding Go models in all relevant microservices (e.g., adding `DeletedAt` for soft delete support, or using pointers for nullable fields).
- **No Manual SQL Only:** Never perform "SQL only" updates to the database without a corresponding migration file and official `goose` execution.

## 6. Plan Tracking & Monitoring

All complex architectural changes or features must have an implementation plan.
- **Location:** All plans must be stored in the `docs/plan/` directory.
- **Completion:** Once a plan is fully implemented and verified, the filename **MUST** be appended with `_COMPLETED` (e.g., `feature-x.md` becomes `feature-x_COMPLETED.md`).
- This allows for easy auditing of implemented vs. pending architectural changes.
