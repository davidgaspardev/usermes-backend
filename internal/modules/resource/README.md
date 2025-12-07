# Resource Module

API endpoints for resource management including CRUD operations and filtering.

## Base URL
All endpoints are prefixed with `/api/resources`

The Resource module manages resources with the following properties:
- **ID**: Unique identifier (UUID)
- **Code**: Resource code (2-50 characters, unique)
- **ShiftID**: Optional shift assignment
- **Type**: Resource type (2-50 characters)
- **StopFactor**: Non-negative integer (int16)
- **Tags**: Optional array of string tags for classification

## 📝 Create Resource
Create a new resource.

**POST** `/api/resources/`

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
  "stop_factor": 5,
  "tags": ["production", "critical"]  // optional
}
```

**Response (201 Created):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "code": "RES001",
  "shift_id": "SHIFT123",
  "type": "MACHINE",
  "stop_factor": 5,
  "tags": ["production", "critical"],
  "created_at": "2024-01-15T10:00:00Z",
  "updated_at": "2024-01-15T10:00:00Z"
}
```

---

**Response:** `200 OK`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "code": "RES001",
  "shift_id": "SHIFT123",
  "type": "MACHINE",
  "stop_factor": 5,
  "tags": ["production", "critical"],
  "created_at": "2024-01-15T10:00:00Z",
  "updated_at": "2024-01-15T10:00:00Z"
}
```

### Get All Resources
Get all resources with pagination.

**GET** `/api/resources/`

```bash
# Default pagination (limit=10, offset=0)
curl -X GET http://localhost:3001/api/resources/

# With custom pagination
curl -X GET "http://localhost:3001/api/resources/?limit=5&offset=10"
```

**Response (200 OK):**
```json
{
  "resources": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "code": "RES001",
      "shift_id": "SHIFT123",
      "type": "MACHINE",
      "stop_factor": 5,
      "tags": ["production", "critical"],
      "created_at": "2024-01-15T10:00:00Z",
      "updated_at": "2024-01-15T10:00:00Z"
    }
  ],
  "total": 1,
  "limit": 10,
  "offset": 0
}
```

### Get Resource by ID
Get specific resource by UUID.

**GET** `/api/resources/:id`

```bash
curl -X GET http://localhost:3001/api/resources/550e8400-e29b-41d4-a716-446655440000
```

### Get Resource by Code
Get specific resource by its code.

**GET** `/api/resources/code/:code`

```bash
curl -X GET http://localhost:3001/api/resources/code/RES001
```

### Get Resources by Type
Get all resources of a specific type.

**GET** `/api/resources/type/:type`

```bash
# Get all MACHINE type resources
curl -X GET http://localhost:3001/api/resources/type/MACHINE

# With pagination
curl -X GET "http://localhost:3001/api/resources/type/MACHINE?limit=5&offset=0"
```

### Get Resources by Shift ID
Get all resources assigned to a specific shift.

**GET** `/api/resources/shift/:shiftId`

```bash
# Get all resources for SHIFT123
curl -X GET http://localhost:3001/api/resources/shift/SHIFT123

# With pagination
curl -X GET "http://localhost:3001/api/resources/shift/SHIFT123?limit=5&offset=0"
```

---

## ✏️ Update Resource
Update an existing resource.

**PUT** `/api/resources/:id`

```bash
curl -X PUT http://localhost:3001/api/resources/550e8400-e29b-41d4-a716-446655440000 \
  -H "Content-Type: application/json" \
  -d '{
    "code": "RES001-UPDATED",
    "shift_id": "SHIFT456",
    "type": "OPERATOR",
    "stop_factor": 10
  }'
```

**Response (200 OK):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "code": "RES001-UPDATED",
  "shift_id": "SHIFT456",
  "type": "OPERATOR",
  "stop_factor": 10,
  "tags": ["maintenance", "high-priority"]  // optional
}
```

---

## 🗑️ Delete Resource
Delete a resource permanently.

**DELETE** `/api/resources/:id`

```bash
curl -X DELETE http://localhost:3001/api/resources/550e8400-e29b-41d4-a716-446655440000
```

**Response (200 OK):**
```json
{
  "message": "Resource deleted successfully"
}
```

---

## ❌ Error Responses

All endpoints may return error responses in the following format:

```json
{
  "error": "resource_not_found",
  "message": "Resource not found"
}
```

### Common Error Codes:
- `invalid_request` - Invalid request format or parameters
- `invalid_id` - Invalid UUID format
- `invalid_code` - Code validation failed
- `invalid_type` - Type validation failed
- `invalid_stop_factor` - Stop factor must be non-negative
- `resource_not_found` - Resource not found
- `resource_already_exists` - Code already exists
- `internal_error` - Server error

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
    tags := []string{"production", "critical"}
    resource, err := service.Create(ctx, "RES001", &shiftID, "MACHINE", 5, tags)
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
    newTags := []string{"maintenance", "high-priority"}
    updated, err := service.Update(ctx, resource.ID(), "RES001-UPDATED", &newShiftID, "OPERATOR", 10, newTags)
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

## 📋 Request Validation Rules

### Create/Update Resource:
- **code**: Required, 2-50 characters, unique
- **shift_id**: Optional, string
- **type**: Required, 2-50 characters
- **stop_factor**: Required, non-negative integer (>= 0)

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

### Tags
- **Required**: No
- **Type**: Array of strings
- **Format**: Any string array when provided
- **Usage**: For classification and filtering resources

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
- [x] Add resource categories/tags *(Completed)*
- [ ] Implement resource availability tracking
- [ ] Add audit logging for changes
- [ ] Implement search and filtering capabilities
- [ ] Add filtering by tags
- [ ] Add bulk operations (create/update/delete multiple)
- [ ] Implement resource scheduling system
- [ ] Add metrics and monitoring

## License

This module is part of the UserMes Backend project.
