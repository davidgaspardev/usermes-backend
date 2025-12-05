# Resource Module

This module implements a complete CRUD (Create, Read, Update, Delete) system for managing resources following Clean Architecture and Domain-Driven Design principles.

## Overview

The Resource module manages resources with the following properties:
- **ID**: Unique identifier (UUID)
- **Code**: Resource code (2-50 characters, unique)
- **ShiftID**: Optional shift assignment
- **Type**: Resource type (2-50 characters)
- **StopFactor**: Non-negative integer (int16)

## Architecture

The module follows Hexagonal Architecture (Ports & Adapters):

```
resource/
├── domain/                           # Core business logic
│   ├── entity/                       # Domain entities
│   │   ├── resource.go              # Resource entity
│   │   └── resource_test.go         # Entity tests
│   └── errors/                       # Domain errors
│       └── errors.go                # Business rule errors
├── application/                      # Use cases
│   ├── port/
│   │   ├── input/                   # Input ports (interfaces)
│   │   │   └── resource_service.go # Service interface
│   │   └── output/                  # Output ports (interfaces)
│   │       └── resource_repository.go # Repository interface
│   └── usecase/                     # Business logic implementation
│       └── resource_service_impl.go # Service implementation
└── infrastructure/                   # External adapters
    ├── adapter/
    │   ├── input/
    │   │   └── http/                # HTTP handlers
    │   │       ├── resource_handler.go
    │   │       └── routes.go
    │   └── output/
    │       └── persistence/         # Data persistence
    │           └── memory_resource_repository.go
    └── dto/                         # Data Transfer Objects
        └── resource_dto.go
```

## API Endpoints

### Create Resource
```http
POST /api/v1/resources
Content-Type: application/json

{
  "code": "RES001",
  "shift_id": "SHIFT123",  // optional
  "type": "MACHINE",
  "stop_factor": 5
}
```

**Response:** `201 Created`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "code": "RES001",
  "shift_id": "SHIFT123",
  "type": "MACHINE",
  "stop_factor": 5,
  "created_at": "2024-01-15T10:00:00Z",
  "updated_at": "2024-01-15T10:00:00Z"
}
```

### Get Resource by ID
```http
GET /api/v1/resources/:id
```

**Response:** `200 OK`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "code": "RES001",
  "shift_id": "SHIFT123",
  "type": "MACHINE",
  "stop_factor": 5,
  "created_at": "2024-01-15T10:00:00Z",
  "updated_at": "2024-01-15T10:00:00Z"
}
```

### Get Resource by Code
```http
GET /api/v1/resources/code/:code
```

### Get All Resources (with pagination)
```http
GET /api/v1/resources?limit=10&offset=0
```

**Response:** `200 OK`
```json
{
  "resources": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "code": "RES001",
      "shift_id": "SHIFT123",
      "type": "MACHINE",
      "stop_factor": 5,
      "created_at": "2024-01-15T10:00:00Z",
      "updated_at": "2024-01-15T10:00:00Z"
    }
  ],
  "total": 1,
  "limit": 10,
  "offset": 0
}
```

### Get Resources by Type
```http
GET /api/v1/resources/type/:type?limit=10&offset=0
```

### Get Resources by Shift ID
```http
GET /api/v1/resources/shift/:shiftId?limit=10&offset=0
```

### Update Resource
```http
PUT /api/v1/resources/:id
Content-Type: application/json

{
  "code": "RES001-UPDATED",
  "shift_id": "SHIFT456",
  "type": "OPERATOR",
  "stop_factor": 10
}
```

**Response:** `200 OK`

### Delete Resource
```http
DELETE /api/v1/resources/:id
```

**Response:** `200 OK`
```json
{
  "message": "Resource deleted successfully",
  "data": null
}
```

## Error Responses

All error responses follow this format:

```json
{
  "error": "error_code",
  "message": "Human-readable error message"
}
```

### Common Error Codes

| Error Code | Status | Description |
|-----------|--------|-------------|
| `invalid_request` | 400 | Invalid request body |
| `invalid_id` | 400 | Invalid UUID format |
| `invalid_code` | 400 | Code must be 2-50 characters |
| `invalid_type` | 400 | Type must be 2-50 characters |
| `invalid_stop_factor` | 400 | Stop factor must be non-negative |
| `resource_not_found` | 404 | Resource not found |
| `code_already_exists` | 409 | Code already in use |
| `resource_already_exists` | 409 | Resource ID already exists |
| `internal_error` | 500 | Internal server error |

## Usage Example

### Integration in main.go

```go
package main

import (
    "log"
    
    resourceHttp "github.com/davidgaspardev/usermes-backend/internal/modules/resource/infrastructure/adapter/input/http"
    resourcePersistence "github.com/davidgaspardev/usermes-backend/internal/modules/resource/infrastructure/adapter/output/persistence"
    resourceUseCase "github.com/davidgaspardev/usermes-backend/internal/modules/resource/application/usecase"
)

func main() {
    // Initialize repository
    resourceRepo := resourcePersistence.NewMemoryResourceRepository()
    
    // Initialize service
    resourceService := resourceUseCase.NewResourceService(resourceRepo)
    
    // Initialize HTTP server (Fiber example)
    app := fiber.New()
    
    // Register routes
    api := app.Group("/api/v1")
    resourceHttp.RegisterRoutes(api, resourceService)
    
    // Start server
    log.Fatal(app.Listen(":3000"))
}
```

### Programmatic Usage

```go
package example

import (
    "context"
    
    "github.com/davidgaspardev/usermes-backend/internal/modules/resource/application/usecase"
    "github.com/davidgaspardev/usermes-backend/internal/modules/resource/infrastructure/adapter/output/persistence"
)

func Example() {
    ctx := context.Background()
    
    // Setup
    repo := persistence.NewMemoryResourceRepository()
    service := usecase.NewResourceService(repo)
    
    // Create resource
    shiftID := "SHIFT123"
    resource, err := service.Create(ctx, "RES001", &shiftID, "MACHINE", 5)
    if err != nil {
        panic(err)
    }
    
    // Get resource by ID
    found, err := service.GetByID(ctx, resource.ID())
    if err != nil {
        panic(err)
    }
    
    // Update resource
    newShiftID := "SHIFT456"
    updated, err := service.Update(ctx, resource.ID(), "RES001-UPDATED", &newShiftID, "OPERATOR", 10)
    if err != nil {
        panic(err)
    }
    
    // Delete resource
    err = service.Delete(ctx, updated.ID())
    if err != nil {
        panic(err)
    }
}
```

## Validation Rules

### Code
- **Required**: Yes
- **Type**: String
- **Length**: 2-50 characters
- **Unique**: Yes
- **Format**: Any non-empty string

### Type
- **Required**: Yes
- **Type**: String
- **Length**: 2-50 characters
- **Format**: Any non-empty string

### Stop Factor
- **Required**: Yes
- **Type**: int16
- **Range**: >= 0 (non-negative)

### Shift ID
- **Required**: No
- **Type**: String (pointer, can be nil)
- **Format**: Any string when provided

## Testing

Run tests for the Resource module:

```bash
# Run all tests
go test ./internal/modules/resource/... -v

# Run with coverage
go test ./internal/modules/resource/... -coverprofile=coverage.out
go tool cover -html=coverage.out

# Run specific tests
go test ./internal/modules/resource/domain/entity -v
```

## Domain Rules

1. **Uniqueness**: Each resource must have a unique code
2. **Stop Factor**: Must be non-negative (>= 0)
3. **Code Validation**: Cannot be empty and must be within length limits
4. **Type Validation**: Cannot be empty and must be within length limits
5. **Shift Assignment**: Optional - resources can exist without a shift assignment

## Future Enhancements

- [ ] Add database persistence (PostgreSQL, MySQL)
- [ ] Implement soft delete functionality
- [ ] Add resource status field (active/inactive)
- [ ] Add resource categories/tags
- [ ] Implement resource availability tracking
- [ ] Add audit logging for changes
- [ ] Implement search and filtering capabilities
- [ ] Add bulk operations (create/update/delete multiple)
- [ ] Implement resource scheduling system
- [ ] Add metrics and monitoring

## License

This module is part of the UserMes Backend project.