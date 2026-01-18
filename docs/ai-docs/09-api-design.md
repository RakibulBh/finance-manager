# API Design & Standards

**Document ID:** 09
**Goal:** Interface standards to ensure the API is easy to use for the frontend.

## 1. The Tech Stack: Chi

The API is built using **chi**, a lightweight, idiomatic router for Go.

```go
import "github.com/go-chi/chi/v5"
```

## 2. Response Format

We follow a consistent JSON structure for all responses.

### Success Response

Wrapped in a `data` object.

```json
{
  "data": {
    "id": "e98e826b-67a5-485a-a388-75c13865612c",
    "email": "user@example.com",
    "family_id": "a1b2c3d4..."
  }
}
```

### Error Response

Contains a single `error` string.

```json
{
  "error": "Invalid request body"
}
```

## 3. Endpoints

### Authentication (Public)

* `POST /api/register`: Register a new user and family.
* `POST /api/login`: Authenticate and receive a JWT.

### Accounts (Protected)

* `GET /api/accounts`: List all family accounts and net worth.
* `POST /api/accounts`: Create a new account (Asset/Liability).

### Transactions (Protected)

* `POST /api/transactions`: Record a standard spend/income.
* `POST /api/transfers`: Move money between two accounts.

### Investments (Protected)

* `POST /api/investments/trade`: Record a security Trade (Buy/Sell).

## 4. Authentication Middleware

All protected routes require a `Authorization: Bearer <token>` header.
The JWT contains the following claims:

* `user_id`: UUID of the authenticated user.
* `family_id`: UUID of the family the user belongs to.
* `exp`: Expiration timestamp (1 week from issuance).

These values are injected into the request context for use by handlers.

## 5. Implementation Roadmap

The Backend API is fully implemented and verified with integration tests.

✅ **API Implementation Complete**
✅ **Integration Testing Complete**
✅ **CI Pipeline Established**

👉 **Proceed to `10-frontend-integration.md`**
