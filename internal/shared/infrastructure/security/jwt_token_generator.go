package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/davidgaspardev/usermes-backend/internal/modules/user/application/port/output"
	"github.com/google/uuid"
)

// JWTTokenGenerator implements the TokenGenerator interface using JWT
type JWTTokenGenerator struct {
	secretKey []byte
	issuer    string
}

// JWTClaims represents the JWT token claims
type JWTClaims struct {
	UserID    string `json:"sub"`
	Email     string `json:"email"`
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
	Issuer    string `json:"iss"`
}

// JWTHeader represents the JWT token header
type JWTHeader struct {
	Algorithm string `json:"alg"`
	Type      string `json:"typ"`
}

// NewJWTTokenGenerator creates a new instance of JWTTokenGenerator
func NewJWTTokenGenerator(secretKey string, issuer string) output.TokenGenerator {
	return &JWTTokenGenerator{
		secretKey: []byte(secretKey),
		issuer:    issuer,
	}
}

// GenerateToken creates a new JWT token for the given user ID
func (g *JWTTokenGenerator) GenerateToken(userID uuid.UUID, email string, expiresIn time.Duration) (string, error) {
	now := time.Now()

	// Create header
	header := JWTHeader{
		Algorithm: "HS256",
		Type:      "JWT",
	}

	// Create claims
	claims := JWTClaims{
		UserID:    userID.String(),
		Email:     email,
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(expiresIn).Unix(),
		Issuer:    g.issuer,
	}

	// Encode header
	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("failed to marshal header: %w", err)
	}
	encodedHeader := base64.RawURLEncoding.EncodeToString(headerJSON)

	// Encode claims
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("failed to marshal claims: %w", err)
	}
	encodedClaims := base64.RawURLEncoding.EncodeToString(claimsJSON)

	// Create signature
	message := encodedHeader + "." + encodedClaims
	signature := g.sign(message)

	// Combine all parts
	token := message + "." + signature

	return token, nil
}

// ValidateToken validates a JWT token and returns the user ID
func (g *JWTTokenGenerator) ValidateToken(token string) (uuid.UUID, error) {
	// Split token into parts
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return uuid.Nil, fmt.Errorf("invalid token format")
	}

	// Verify signature
	message := parts[0] + "." + parts[1]
	expectedSignature := g.sign(message)
	if parts[2] != expectedSignature {
		return uuid.Nil, fmt.Errorf("invalid token signature")
	}

	// Decode claims
	claimsJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to decode claims: %w", err)
	}

	var claims JWTClaims
	if err := json.Unmarshal(claimsJSON, &claims); err != nil {
		return uuid.Nil, fmt.Errorf("failed to unmarshal claims: %w", err)
	}

	// Check expiration
	if time.Now().Unix() > claims.ExpiresAt {
		return uuid.Nil, fmt.Errorf("token expired")
	}

	// Parse user ID
	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID in token: %w", err)
	}

	return userID, nil
}

// RefreshToken generates a new token from an existing valid token
func (g *JWTTokenGenerator) RefreshToken(token string, expiresIn time.Duration) (string, error) {
	// Validate the existing token
	userID, err := g.ValidateToken(token)
	if err != nil {
		return "", fmt.Errorf("cannot refresh invalid token: %w", err)
	}

	// Extract email from the old token
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", fmt.Errorf("invalid token format")
	}

	claimsJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", fmt.Errorf("failed to decode claims: %w", err)
	}

	var claims JWTClaims
	if err := json.Unmarshal(claimsJSON, &claims); err != nil {
		return "", fmt.Errorf("failed to unmarshal claims: %w", err)
	}

	// Generate new token with the same user ID and email
	return g.GenerateToken(userID, claims.Email, expiresIn)
}

// sign creates a signature for the given message using HMAC-SHA256
func (g *JWTTokenGenerator) sign(message string) string {
	h := hmac.New(sha256.New, g.secretKey)
	h.Write([]byte(message))
	signature := h.Sum(nil)
	return base64.RawURLEncoding.EncodeToString(signature)
}
