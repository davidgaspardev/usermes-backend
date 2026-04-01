# UserMes Backend

[![CI](https://github.com/davidgaspardev/usermes-backend/actions/workflows/ci.yml/badge.svg)](https://github.com/davidgaspardev/usermes-backend/actions/workflows/ci.yml)

Go backend for a Manufacturing Execution System (MES). Manages users, organizational location hierarchies, shift patterns, and production resources.

## Architecture

Hexagonal architecture (Ports & Adapters) organized as a modular monolith. Each module has its own domain, application, and infrastructure layers. Modules reference each other only by ID — no direct object coupling.

```
internal/
├── modules/
│   ├── iam/                   # Identity & access management
│   ├── organization/          # Location hierarchy & shift patterns
│   └── production/            # Production resources
└── shared/
    └── infrastructure/
        ├── http/server/       # Fiber server setup
        └── security/          # JWT implementation
```

Each module follows the same structure:

```
module/
├── domain/
│   ├── entity/
│   └── errors/
├── application/
│   ├── port/
│   │   ├── input/     # Use case interfaces + commands
│   │   └── output/    # Repository interfaces
│   └── usecase/
└── infrastructure/
    ├── adapter/
    │   ├── input/http/          # Fiber handlers + routes + middleware
    │   └── output/persistence/  # In-memory repositories
    └── dto/
```

## Modules

### IAM

User registration, login, profile management, and JWT token issuance.

- **Entities:** User
- **Value Objects:** Email, Password

### Organization

Location hierarchy (Plant → Area → Line → Section) and shift pattern definitions. A location is a node in a tree — the root is always a Plant.

- **Entities:** Location, ShiftPattern
- **Notes:** Locations track `created_by` (user UUID). Child locations automatically inherit a default shift pattern if none is specified.

### Production

Production resources scoped to a location code.

- **Entities:** Resource
- **Notes:** All resource routes require a valid `:location_code` in the path (enforced by middleware).

## API Routes

All routes are prefixed with `/v1/api`.

### IAM — `/v1/api/iam`

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/users/register` | — | Register user |
| POST | `/users/login` | — | Login, returns JWT |
| GET | `/users/me` | ✅ | Get own profile |
| GET | `/users/:id` | ✅ | Get user by ID |
| PUT | `/users/:id` | ✅ | Update user |
| POST | `/users/:id/change-password` | ✅ | Change password |
| POST | `/users/:id/deactivate` | ✅ | Deactivate user |
| POST | `/users/:id/activate` | ✅ | Activate user |

### Organization — `/v1/api/organization`

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/locations/` | ✅ | Create root location (Plant) |
| POST | `/locations/add` | ✅ | Add child location |
| GET | `/locations/` | — | List all location trees |
| GET | `/locations/:location_code` | — | Get location tree by code |
| POST | `/locations/:location_code/shift-patterns` | ✅ | Assign shift pattern to location |
| GET | `/shift-patterns/:id` | ✅ | Get shift pattern |
| PUT | `/shift-patterns/:id` | ✅ | Update shift pattern |
| DELETE | `/shift-patterns/:id` | ✅ | Delete shift pattern |

### Production — `/v1/api/production`

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/locations/:location_code/resources/` | ✅ | Create resource |
| GET | `/locations/:location_code/resources/` | ✅ | List resources |
| GET | `/locations/:location_code/resources/:id` | ✅ | Get resource by ID |
| GET | `/locations/:location_code/resources/code/:res_code` | ✅ | Get resource by code |
| GET | `/locations/:location_code/resources/type/:type` | ✅ | Get resources by type |
| PUT | `/locations/:location_code/resources/:id` | ✅ | Update resource |
| DELETE | `/locations/:location_code/resources/:id` | ✅ | Delete resource |

> ✅ = requires `Authorization: Bearer <token>`. Production routes also validate that `:location_code` exists before reaching the handler.

## Running

**Requirements:** Go 1.22+

```bash
make run        # start server on :3001
make dev        # start with hot reload (requires air)
```

**Environment variables:**

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `3001` | HTTP port |

## Testing

```bash
make test           # run all tests
make test-coverage  # run tests with coverage threshold check
make coverage-html  # generate HTML coverage report
```
