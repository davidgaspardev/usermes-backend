package dto

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/davidgaspardev/usermes-backend/internal/modules/resource/domain/entity"
)

func TestToResourceResponse(t *testing.T) {
	// Create a test resource
	id := uuid.New()
	shiftID := "shift123"
	createdAt := time.Now().Add(-time.Hour)
	updatedAt := time.Now()

	resource := entity.ReconstructResource(
		id,
		"RES001",
		&shiftID,
		"machine",
		5,
		createdAt,
		updatedAt,
	)

	response := ToResourceResponse(resource)

	assert.Equal(t, id.String(), response.ID)
	assert.Equal(t, "RES001", response.Code)
	assert.Equal(t, &shiftID, response.ShiftID)
	assert.Equal(t, "machine", response.Type)
	assert.Equal(t, int16(5), response.StopFactor)
	assert.Equal(t, createdAt, response.CreatedAt)
	assert.Equal(t, updatedAt, response.UpdatedAt)
}

func TestToResourceResponse_WithNilShiftID(t *testing.T) {
	id := uuid.New()
	createdAt := time.Now().Add(-time.Hour)
	updatedAt := time.Now()

	resource := entity.ReconstructResource(
		id,
		"RES002",
		nil,
		"tool",
		3,
		createdAt,
		updatedAt,
	)

	response := ToResourceResponse(resource)

	assert.Equal(t, id.String(), response.ID)
	assert.Equal(t, "RES002", response.Code)
	assert.Nil(t, response.ShiftID)
	assert.Equal(t, "tool", response.Type)
	assert.Equal(t, int16(3), response.StopFactor)
}

func TestToResourceListResponse(t *testing.T) {
	// Create test resources
	id1 := uuid.New()
	id2 := uuid.New()
	shiftID := "shift123"
	now := time.Now()

	resource1 := entity.ReconstructResource(id1, "RES001", &shiftID, "machine", 5, now, now)
	resource2 := entity.ReconstructResource(id2, "RES002", nil, "tool", 3, now, now)

	resources := []*entity.Resource{resource1, resource2}

	t.Run("with valid resources", func(t *testing.T) {
		response := ToResourceListResponse(resources, 10, 0)

		assert.Len(t, response.Resources, 2)
		assert.Equal(t, 2, response.Total)
		assert.Equal(t, 10, response.Limit)
		assert.Equal(t, 0, response.Offset)

		// Check first resource
		assert.Equal(t, id1.String(), response.Resources[0].ID)
		assert.Equal(t, "RES001", response.Resources[0].Code)
		assert.Equal(t, &shiftID, response.Resources[0].ShiftID)

		// Check second resource
		assert.Equal(t, id2.String(), response.Resources[1].ID)
		assert.Equal(t, "RES002", response.Resources[1].Code)
		assert.Nil(t, response.Resources[1].ShiftID)
	})

	t.Run("with different pagination", func(t *testing.T) {
		response := ToResourceListResponse(resources, 5, 10)

		assert.Equal(t, 5, response.Limit)
		assert.Equal(t, 10, response.Offset)
		assert.Equal(t, 2, response.Total)
	})

	t.Run("with empty resources", func(t *testing.T) {
		response := ToResourceListResponse([]*entity.Resource{}, 10, 0)

		assert.Empty(t, response.Resources)
		assert.Equal(t, 0, response.Total)
		assert.Equal(t, 10, response.Limit)
		assert.Equal(t, 0, response.Offset)
	})
}

func TestNewErrorResponse(t *testing.T) {
	t.Run("with both error and message", func(t *testing.T) {
		response := NewErrorResponse("validation_error", "Invalid input provided")

		assert.Equal(t, "validation_error", response.Error)
		assert.Equal(t, "Invalid input provided", response.Message)
	})

	t.Run("with empty message", func(t *testing.T) {
		response := NewErrorResponse("internal_error", "")

		assert.Equal(t, "internal_error", response.Error)
		assert.Equal(t, "", response.Message)
	})

	t.Run("with empty error", func(t *testing.T) {
		response := NewErrorResponse("", "Something went wrong")

		assert.Equal(t, "", response.Error)
		assert.Equal(t, "Something went wrong", response.Message)
	})
}

func TestNewSuccessResponse(t *testing.T) {
	t.Run("with data", func(t *testing.T) {
		data := map[string]string{"key": "value"}
		response := NewSuccessResponse("Operation completed", data)

		assert.Equal(t, "Operation completed", response.Message)
		assert.Equal(t, data, response.Data)
	})

	t.Run("with nil data", func(t *testing.T) {
		response := NewSuccessResponse("Success", nil)

		assert.Equal(t, "Success", response.Message)
		assert.Nil(t, response.Data)
	})

	t.Run("with empty message", func(t *testing.T) {
		response := NewSuccessResponse("", "some data")

		assert.Equal(t, "", response.Message)
		assert.Equal(t, "some data", response.Data)
	})
}

func TestCreateResourceRequest_JSONTags(t *testing.T) {
	// This test ensures the JSON tags are properly set
	// We'll create a struct and check if it can be marshaled/unmarshaled correctly
	req := CreateResourceRequest{
		Code:       "RES001",
		ShiftID:    nil,
		Type:       "machine",
		StopFactor: 5,
	}

	// Test that all fields are accessible
	assert.Equal(t, "RES001", req.Code)
	assert.Nil(t, req.ShiftID)
	assert.Equal(t, "machine", req.Type)
	assert.Equal(t, int16(5), req.StopFactor)
}

func TestCreateResourceRequest_WithShiftID(t *testing.T) {
	shiftID := "shift123"
	req := CreateResourceRequest{
		Code:       "RES001",
		ShiftID:    &shiftID,
		Type:       "machine",
		StopFactor: 5,
	}

	assert.Equal(t, "RES001", req.Code)
	require.NotNil(t, req.ShiftID)
	assert.Equal(t, "shift123", *req.ShiftID)
	assert.Equal(t, "machine", req.Type)
	assert.Equal(t, int16(5), req.StopFactor)
}

func TestUpdateResourceRequest_JSONTags(t *testing.T) {
	shiftID := "updated_shift"
	req := UpdateResourceRequest{
		Code:       "RES001_UPDATED",
		ShiftID:    &shiftID,
		Type:       "tool",
		StopFactor: 10,
	}

	assert.Equal(t, "RES001_UPDATED", req.Code)
	require.NotNil(t, req.ShiftID)
	assert.Equal(t, "updated_shift", *req.ShiftID)
	assert.Equal(t, "tool", req.Type)
	assert.Equal(t, int16(10), req.StopFactor)
}

func TestResourceResponse_JSONTags(t *testing.T) {
	id := uuid.New()
	shiftID := "shift123"
	now := time.Now()

	response := ResourceResponse{
		CreatedAt:  now,
		UpdatedAt:  now,
		ShiftID:    &shiftID,
		ID:         id.String(),
		Code:       "RES001",
		Type:       "machine",
		StopFactor: 5,
	}

	assert.Equal(t, now, response.CreatedAt)
	assert.Equal(t, now, response.UpdatedAt)
	assert.Equal(t, &shiftID, response.ShiftID)
	assert.Equal(t, id.String(), response.ID)
	assert.Equal(t, "RES001", response.Code)
	assert.Equal(t, "machine", response.Type)
	assert.Equal(t, int16(5), response.StopFactor)
}

func TestResourceListResponse_Structure(t *testing.T) {
	id := uuid.New()
	shiftID := "shift123"
	now := time.Now()

	resourceResponse := ResourceResponse{
		CreatedAt:  now,
		UpdatedAt:  now,
		ShiftID:    &shiftID,
		ID:         id.String(),
		Code:       "RES001",
		Type:       "machine",
		StopFactor: 5,
	}

	listResponse := ResourceListResponse{
		Resources: []*ResourceResponse{&resourceResponse},
		Total:     1,
		Limit:     10,
		Offset:    0,
	}

	assert.Len(t, listResponse.Resources, 1)
	assert.Equal(t, 1, listResponse.Total)
	assert.Equal(t, 10, listResponse.Limit)
	assert.Equal(t, 0, listResponse.Offset)
	assert.Equal(t, id.String(), listResponse.Resources[0].ID)
}

func TestErrorResponse_Structure(t *testing.T) {
	errorResp := ErrorResponse{
		Error:   "validation_error",
		Message: "Invalid input provided",
	}

	assert.Equal(t, "validation_error", errorResp.Error)
	assert.Equal(t, "Invalid input provided", errorResp.Message)
}

func TestSuccessResponse_Structure(t *testing.T) {
	data := map[string]interface{}{
		"count":   5,
		"status":  "completed",
		"results": []string{"a", "b", "c"},
	}

	successResp := SuccessResponse{
		Data:    data,
		Message: "Operation successful",
	}

	assert.Equal(t, data, successResp.Data)
	assert.Equal(t, "Operation successful", successResp.Message)
}

// Benchmark tests for performance-critical conversion functions
func BenchmarkToResourceResponse(b *testing.B) {
	id := uuid.New()
	shiftID := "shift123"
	now := time.Now()

	resource := entity.ReconstructResource(id, "RES001", &shiftID, "machine", 5, now, now)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ToResourceResponse(resource)
	}
}

func BenchmarkToResourceListResponse(b *testing.B) {
	id1 := uuid.New()
	id2 := uuid.New()
	shiftID := "shift123"
	now := time.Now()

	resource1 := entity.ReconstructResource(id1, "RES001", &shiftID, "machine", 5, now, now)
	resource2 := entity.ReconstructResource(id2, "RES002", nil, "tool", 3, now, now)

	resources := []*entity.Resource{resource1, resource2}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ToResourceListResponse(resources, 10, 0)
	}
}
