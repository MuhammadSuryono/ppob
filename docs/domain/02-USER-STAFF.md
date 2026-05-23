# 👥 Domain: User & Staff Management

## 1. Core Mission
This domain manages the hierarchy and operational relationships between different types of business partners (Mitra) and their employees (Staff).

## 2. Key Concepts

### 2.1 Roles
- **Mitra (Partner):** The owner of a PPOB business. Has the authority to manage staff, set margins, and fund staff wallets.
- **Staff:** An employee assigned to a Mitra. Can perform transactions on behalf of the Mitra but has limited administrative capabilities.
- **Super-Admin:** Platform operator with global access (limited to system management).

### 2.2 Multi-Tenancy
A user can possess multiple roles across different Mitra organizations.
- **Active Role:** The role currently used by the user, which determines the operational context (wallet, staff list, reports).

### 2.3 Staff Lifecycle
- **Assignment:** Staff must be assigned by a Mitra.
- **Permissions:** Mitra can toggle staff status (Active/Inactive) and override default transaction limits.

## 3. Business Rules

### 3.1 Hierarchy
- Staff are directly linked to a Mitra via the `assigned_by` relationship.
- A Mitra can only manage staff that they have created/assigned.

### 3.2 Operating Constraints
- **Role Inactivity:** Users cannot switch to a role that has been marked as inactive.
- **Ownership:** Staff cannot manage other staff or access the Mitra's main reporting dashboard.

## 4. Domain Logic
- **Context Switching:** When a user switches roles, the system must transparently update the session context (including the active wallet ID) to prevent data leakage between tenants.
- **Automatic Provisioning:** Assigning a staff role automatically triggers the creation of a sub-wallet and initialization of default margin settings.
