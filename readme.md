# QC System - Engineering Drawing Workflow

*   **Workflow Management**: Strict state transitions (Drafting -> First QC -> Final QC -> Approved).
*   **Role-Based Access Control (RBAC)**: Fine-grained permissions for Admins, Drafters, Shift Leads, and Final QC inspectors using **Casbin** (https://github.com/casbin/casbin).
*   **Real-time Collaboration**: Instant updates via **Server-Sent Events (SSE)** and **Redis Pub/Sub** when drawings are claimed or updated. (Here sockets will be overkill as we do not need two way changes| Also cannot just broadcast event to all users therefore used redis - A simple mimic of socket.io for go)
*   **Concurrency Control**: specialized locking mechanisms to prevent race conditions (see "Concurrency Strategy" below).
*   **Audit Logging**: Immutable logs for every workflow transition for accountability. The logs are sent to the kafka (Not consumed anywhere for now: But should be consumed by s3 or can put in some DB async for later retrieval)

---

## 🧠 Concurrency & Data Integrity Strategy

The core challenge in a collaborative QC system is preventing **Race Conditions**, such as two users claiming the same drawing simultaneously ("Double Claiming").

### 1. Pessimistic Locking (Primary Defense)
We utilize **Database-Level Pessimistic Locking** (`SELECT ... FOR UPDATE`) within a strict transaction boundary for critical actions like **Claiming** a drawing.

*   **How it works**: When a user attempts to "Claim" a drawing, the backend initiates a transaction and immediately locks the specific drawing row.
*   **Result**: If `User A` and `User B` click "Claim" at the exact same millisecond:
    1.  The database grants the lock to one transaction (say, `User A`).
    2.  `User B`'s transaction is forced to **wait** (block) until `User A` commits.
    3.  Once `User A` commits (state changes to `Drafting`), `User B`'s transaction resumes, reads the *new* state, sees it is already assigned, and fails gracefully.
*   **Code Reference**: `repositories/drawing_repo.go` -> `GetForUpdate(id)`

### 2. Optimistic Locking (Secondary Defense)
All drawings have a `version` column. Every update operation checks `WHERE id = ? AND version = ?`.
*   This ensures that if a user is looking at stale data (e.g., they loaded the page 5 minutes ago) and tries to perform an action, the update will fail because the version in the database has incremented.

### 3. Valid State Transitions
We implement a strict **Finite State Machine (FSM)** in `models/workflow.go`.
*   Transitions are validated against the current `Stage`, the requested `Action`, and the user's `Role`.
*   An invalid transition (e.g., a Drafter trying to Approve a drawing) is rejected at the domain level before hitting the database.

---

##  Architecture

### Backend (Go / Gin)
*   **Layered Architecture**: Strictly separated into `Controllers` (HTTP), `Services` (Business Logic), and `Repositories` (Data Access).
*   **Dependency Injection**: All dependencies (Repositories, Audit Service, Realtime Broadcaster) are injected via interfaces, making the system highly testable and loosely coupled.

### Frontend (React / Vite)
*   **Real-time UI**: Listens to SSE streams to update the dashboard instantly without manual refreshes.
*   **Optimistic UI**: Provides immediate feedback to users while handling backend rejections gracefully.

---

## Setup & Installation


### 1. Backend Setup
```bash
cd backend
go mod download

# Set environment variables (or rely on defaults)

# Run the server
go run main.go
```
*The server starts on port `8081`.*

### 2. Frontend Setup
```bash
cd frontend
npm install

# Run the development server
npm run dev
```
*The frontend starts on `http://localhost:5173`.*
