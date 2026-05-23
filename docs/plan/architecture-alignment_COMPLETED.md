# Architecture Alignment Plan

## 1. Background & Motivation
The current implementation diverges from the `docs/architecture` blueprint in several key areas. To ensure the code adheres to the defined architecture, a comprehensive refactoring is required. 

## 2. Scope
This alignment will span across the shared libraries and multiple microservices (`auth`, `wallet`, `transaction`, `product`). 

## 3. Implementation Phases

### Phase 0: Plan Tracking & Rules Update
- Update `GEMINI.md` to add the new plan monitoring rule: All implementation plans must reside in `docs/plan/`, and upon completion, their filenames must be appended with `_COMPLETED`.
- Move this plan file into the newly created `docs/plan/` directory.

### Phase 1: Security Upgrade
- **Argon2id for PINs:** Update `auth-service` (`Register`, `ChangePIN`, `VerifyPINLogin`, `AuthorizeTransaction`) to use Argon2id for hashing transaction PINs instead of bcrypt.
- **JWT RS256:** 
  - Introduce RSA key pair generation/loading in the environment or configuration.
  - Update `auth-service` to sign JWTs using RS256 instead of HS256.
  - Update `ValidateTokenStatic` to verify JWTs using the public RSA key.

### Phase 2: gRPC Migration
- **Protobuf Contracts:** Create a `shared/proto` directory. Define `.proto` files for internal service communication (e.g., `wallet.proto`, `product.proto`).
- **Code Generation:** Generate Go gRPC client and server code within the `shared/` directory.
- **Server Implementation:** Expose gRPC servers in `wallet-service` (for hold/release/debit operations) and `product-service` (for price/availability validation).
- **Client Migration:** Refactor `transaction-service` internal clients (currently using HTTP REST) to use the new gRPC clients.

### Phase 3: Event-Driven Architecture (Redis Streams)
- **Shared Event Utility:** Build a Redis Streams publisher and consumer wrapper in the `shared/` directory.
- **Event Publishers:** 
  - `auth-service` to publish `user.registered`.
  - `transaction-service` to publish `transaction.success`.
- **Event Consumers:**
  - `wallet-service` to listen to `user.registered` for automatic sub-wallet creation.
  - `wallet-service` (or a dedicated worker) to listen to `transaction.success` for commission distribution.

## 4. Verification & Testing
- Run `make build` and ensure all services compile successfully.
- Ensure existing integration tests pass and write new tests covering gRPC endpoints.
- Spawn `codebase_investigator` subagent after each phase to verify alignment with `docs/architecture`.
- Rename this plan file to `architecture-alignment_COMPLETED.md` once all phases are finished.