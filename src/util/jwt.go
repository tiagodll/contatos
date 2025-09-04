package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// ValidateJWT validates a JWT token and returns the user ID
func ValidateJWT(tokenString string, config JwtConfig) (string, error) {
	if tokenString == "" {
		return "", errors.New("no token provided")
	}

	if config.Secret == "" {
		return "", errors.New("JWT secret not configured")
	}

	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Secret), nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to parse token: %v", err)
	}

	// Extract claims
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token")
	}

	// Validate issuer if configured
	if config.Issuer != "" && claims.Issuer != config.Issuer {
		return "", errors.New("invalid token issuer")
	}

	// Validate audience if configured
	if config.Audience != "" {
		if !slices.Contains(claims.Audience, config.Audience) {
			return "", errors.New("invalid token audience")
		}
	}

	// Check if token is expired
	if claims.ExpiresAt != nil && time.Now().After(claims.ExpiresAt.Time) {
		return "", errors.New("token has expired")
	}

	// Validate user ID
	if claims.Subject == "" {
		return "", errors.New("no user ID in token")
	}

	return claims.Subject, nil
}

// GenerateJWT generates a new JWT token for a user (utility function for testing)
func GenerateJWT(userID string, config JwtConfig) (string, error) {
	if config.Secret == "" {
		return "", errors.New("JWT secret not configured")
	}

	// Set expiration time (24 hours from now)
	expirationTime := time.Now().Add(24 * time.Hour)

	// Create claims
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// Add issuer if configured
	if config.Issuer != "" {
		claims.Issuer = config.Issuer
	}

	// Add audience if configured
	if config.Audience != "" {
		claims.Audience = []string{config.Audience}
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret
	tokenString, err := token.SignedString([]byte(config.Secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return tokenString, nil
}

func DecodeJWT(token, secret string) (map[string]any, error) {
	// Split token into parts
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid JWT format")
	}
	header, payload, signature := parts[0], parts[1], parts[2]

	// Decode payload
	payloadBytes, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to decode payload: %v", err)
	}

	// Parse payload JSON
	var payloadMap map[string]any
	if err := json.Unmarshal(payloadBytes, &payloadMap); err != nil {
		return nil, fmt.Errorf("failed to parse payload JSON: %v", err)
	}

	// Verify signature
	unsignedToken := header + "." + payload
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(unsignedToken))
	expectedSignature := base64.RawURLEncoding.EncodeToString(h.Sum(nil))
	if signature != expectedSignature {
		return nil, fmt.Errorf("invalid signature")
	}

	// Check expiration
	if exp, ok := payloadMap["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return nil, fmt.Errorf("token expired")
		}
	} else {
		return nil, fmt.Errorf("missing or invalid exp claim")
	}

	return payloadMap, nil
}
