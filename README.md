# ZeroTrust Proxy

A lightweight, high-performance Zero-Trust Identity Aware Proxy (IAP) written in Go. This proxy sits in front of internal microservices and enforces strict authentication, authorization, and network boundaries without requiring changes to the underlying services.

## Overview

In a Zero Trust Architecture, no user or service is inherently trusted, even if they are inside the corporate network. `ZeroTrust Proxy` intercepts all traffic, cryptographically verifies Identity (via JWT tokens), enforces Role-Based Access Control (RBAC), and proxies the validated request with injected identity headers to the backend.

## Key Features

- **Strict Identity Verification:** Validates JWT access tokens signed by an Identity Provider (IdP).
- **Role-Based Access Control (RBAC):** Blocks unauthorized roles before traffic reaches the backend application.
- **Header Injection:** Securely passes the validated `X-Authenticated-User` and `X-User-Role` downstream.
- **High Performance:** Built on Go's robust `net/http/httputil` standard library reverse proxy.
- **Logging & Auditing:** Comprehensive audit logs for allowed and denied requests.

## Architecture

```text
[Client] --> HTTPS/mTLS --> [ZeroTrust Proxy] -- HTTP --> [Backend Microservice]
                              |
                        (JWT Validation & RBAC)
```

## Setup

1. Make sure you have Go 1.21+ installed.
2. Clone the repository and fetch dependencies:
   ```bash
   go mod download
   ```
3. Run the proxy:
   ```bash
   go run main.go
   ```
4. The proxy listens on `:8443` and forwards to `:8080`.

## Testing the Proxy

**1. Unauthorized Request (Blocked):**
```bash
curl -i http://localhost:8443/api/data
# HTTP/1.1 401 Unauthorized
```

**2. Authorized Request (Allowed):**
You will need to generate a valid JWT signed with the `JWTSecret` defined in `main.go`.
```bash
curl -i -H "Authorization: Bearer <YOUR_JWT_TOKEN>" http://localhost:8443/api/data
# Proxies to backend and returns response.
```
