# Resource Module

API endpoints for resource management including CRUD operations and filtering.

## Base URL
All endpoints are prefixed with `/api/resources`

---

## 📝 Create Resource
Create a new resource.

**POST** `/api/resources/`

```bash
curl -X POST http://localhost:8080/api/resources/ \
  -H "Content-Type: application/json" \
  -d '{
    "code": "RES001",
    "shift_id": "SHIFT123",
    "type": "MACHINE",
    "stop_factor": 5
  }'
```

**Response (201 Created):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "code": "RES001",
  "shift_id": "SHIFT123",
  "type": "MACHINE",
  "stop_factor": 5,
  "created_at": "2023-12-01T10:00:00Z",
  "updated_at": "2023-12-01T10:00:00Z"
}
```

---

## 🔍 Get Resources

### Get All Resources
Get all resources with pagination.

**GET** `/api/resources/`

```bash
# Default pagination (limit=10, offset=0)
curl -X GET http://localhost:8080/api/resources/

# With custom pagination
curl -X GET "http://localhost:8080/api/resources/?limit=5&offset=10"
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
      "created_at": "2023-12-01T10:00:00Z",
      "updated_at": "2023-12-01T10:00:00Z"
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
curl -X GET http://localhost:8080/api/resources/550e8400-e29b-41d4-a716-446655440000
```

### Get Resource by Code
Get specific resource by its code.

**GET** `/api/resources/code/:code`

```bash
curl -X GET http://localhost:8080/api/resources/code/RES001
```

### Get Resources by Type
Get all resources of a specific type.

**GET** `/api/resources/type/:type`

```bash
# Get all MACHINE type resources
curl -X GET http://localhost:8080/api/resources/type/MACHINE

# With pagination
curl -X GET "http://localhost:8080/api/resources/type/MACHINE?limit=5&offset=0"
```

### Get Resources by Shift ID
Get all resources assigned to a specific shift.

**GET** `/api/resources/shift/:shiftId`

```bash
# Get all resources for SHIFT123
curl -X GET http://localhost:8080/api/resources/shift/SHIFT123

# With pagination
curl -X GET "http://localhost:8080/api/resources/shift/SHIFT123?limit=5&offset=0"
```

---

## ✏️ Update Resource
Update an existing resource.

**PUT** `/api/resources/:id`

```bash
curl -X PUT http://localhost:8080/api/resources/550e8400-e29b-41d4-a716-446655440000 \
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
  "created_at": "2023-12-01T10:00:00Z",
  "updated_at": "2023-12-01T11:00:00Z"
}
```

---

## 🗑️ Delete Resource
Delete a resource permanently.

**DELETE** `/api/resources/:id`

```bash
curl -X DELETE http://localhost:8080/api/resources/550e8400-e29b-41d4-a716-446655440000
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

---

## 📋 Request Validation Rules

### Create/Update Resource:
- **code**: Required, 2-50 characters, unique
- **shift_id**: Optional, string
- **type**: Required, 2-50 characters
- **stop_factor**: Required, non-negative integer (>= 0)

### Examples of Valid Data:
```json
{
  "code": "MACHINE_001",
  "shift_id": "MORNING_SHIFT",
  "type": "PRODUCTION_MACHINE",
  "stop_factor": 0
}
```

```json
{
  "code": "OP_042",
  "shift_id": null,
  "type": "OPERATOR",
  "stop_factor": 15
}
```
