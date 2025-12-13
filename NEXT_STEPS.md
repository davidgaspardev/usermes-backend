# Next Steps - Complete the Migration

## 🚀 Quick Guide to Complete the Work

### Step 1: Fix Resource Tests (5 min)

```bash
# Update test files to include plantCode
find internal/modules/production/resource -name "*_test.go" -exec \
  sed -i '' 's/NewResource(/NewResource("SP01", /g' {} \;

# Run tests
go test ./internal/modules/production/resource/... -v
```

### Step 2: Create Plant HTTP Handler (10 min)

File: `internal/modules/organization/plant/infrastructure/adapter/input/http/plant_handler.go`

```go
package http

import (
    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"
    "github.com/davidgaspardev/usermes-backend/internal/modules/organization/plant/application/port/input"
    "github.com/davidgaspardev/usermes-backend/internal/modules/organization/plant/infrastructure/dto"
)

type PlantHandler struct {
    service input.PlantService
}

func NewPlantHandler(service input.PlantService) *PlantHandler {
    return &PlantHandler{service: service}
}

func (h *PlantHandler) Create(c *fiber.Ctx) error {
    var req dto.CreatePlantRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(dto.NewErrorResponse("bad_request", err.Error()))
    }
    
    // Get user ID from JWT context
    userID := c.Locals("user_id").(uuid.UUID)
    
    plant, err := h.service.Create(c.Context(), req.Code, req.Name, req.Latitude, req.Longitude, userID)
    if err != nil {
        return c.Status(400).JSON(dto.NewErrorResponse("error", err.Error()))
    }
    
    return c.Status(201).JSON(dto.ToPlantResponse(plant))
}

// Add GetByCode, ListAll, Update, Delete methods...
```

### Step 3: Create Plant Routes (5 min)

File: `internal/modules/organization/plant/infrastructure/adapter/input/http/routes.go`

```go
package http

import "github.com/gofiber/fiber/v2"

type PlantRoutes struct {
    handler *PlantHandler
}

func NewPlantRoutes(handler *PlantHandler) *PlantRoutes {
    return &PlantRoutes{handler: handler}
}

func (r *PlantRoutes) Register(app *fiber.App, prefix string) {
    plants := app.Group(prefix + "/plants")
    
    plants.Post("/", r.handler.Create)              // POST /v1/plants
    plants.Get("/", r.handler.ListAll)              // GET /v1/plants
    plants.Get("/:code", r.handler.GetByCode)       // GET /v1/plants/:code
    plants.Put("/:id", r.handler.Update)            // PUT /v1/plants/:id
    plants.Delete("/:id", r.handler.Delete)         // DELETE /v1/plants/:id
}
```

### Step 4: Update Resource Handler (10 min)

Update `internal/modules/production/resource/infrastructure/adapter/input/http/resource_handler.go`:

```go
func (h *ResourceHandler) Create(c *fiber.Ctx) error {
    plantCode := c.Params("plant_code")  // Get from URL
    
    var req dto.CreateResourceRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(dto.NewErrorResponse("bad_request", err.Error()))
    }
    
    resource, err := h.service.Create(
        c.Context(),
        plantCode,  // Add this parameter
        req.Code,
        req.ShiftID,
        req.Type,
        req.StopFactor,
        req.Tags,
    )
    // ...
}
```

### Step 5: Wire Everything in main.go (10 min)

```go
// ... existing imports ...
import (
    plantUseCase "github.com/davidgaspardev/usermes-backend/internal/modules/organization/plant/application/usecase"
    plantHTTP "github.com/davidgaspardev/usermes-backend/internal/modules/organization/plant/infrastructure/adapter/input/http"
    plantPersistence "github.com/davidgaspardev/usermes-backend/internal/modules/organization/plant/infrastructure/adapter/output/persistence"
)

func main() {
    // ... existing code ...
    
    // Plant module
    plantRepo := plantPersistence.NewInMemoryPlantRepository()
    plantService := plantUseCase.NewPlantService(plantRepo)
    plantHandler := plantHTTP.NewPlantHandler(plantService)
    plantRoutes := plantHTTP.NewPlantRoutes(plantHandler)
    
    // Register routes
    userRoutes.Register(app, "/v1")
    plantRoutes.Register(app, "/v1")
    
    // Resource routes with plant scope
    resourceGroup := app.Group("/v1/plants/:plant_code/production")
    resourceHandler := resourceHTTP.NewResourceHandler(resourceService)
    // Register resource routes on the group
    
    // ... rest of code ...
}
```

### Step 6: Test Everything (5 min)

```bash
# Compile
go build -o bin/usermes cmd/main.go

# Run
./bin/usermes

# Test in another terminal
curl -X POST http://localhost:3001/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","password":"Test123!","username":"test","name":"Test User"}'

curl -X POST http://localhost:3001/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"Test123!"}'

# Use the token to create a plant
curl -X POST http://localhost:3001/v1/plants \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"code":"SP01","name":"Test Plant","latitude":-23.5,"longitude":-46.6}'
```

## ⏱️ Time Estimate

Total: ~45 minutes to complete everything

## 🎯 Priority Order

1. **Fix Resource tests** (must do - breaks build)
2. **Update Resource Handler** (must do - breaks build)
3. **Create Plant Handler & Routes** (important - new feature)
4. **Wire in main.go** (important - makes it work)
5. **Manual testing** (validation)

## 📚 Reference

- Module structure example: `internal/modules/iam/user/`
- Handler example: `internal/modules/iam/user/infrastructure/adapter/input/http/`
- Routes example: `internal/modules/iam/user/infrastructure/adapter/input/http/routes.go`

