package middleware

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

type (
	// ClerkAuthConfig defines the config for ClerkAuth middleware
	ClerkAuthConfig struct {
		// Skipper defines a function to skip middleware
		Skipper echomiddleware.Skipper

		// JWKSUrl is the URL to fetch Clerk JWKS
		JWKSUrl string

		// ContextKey is the key used to store user ID in context
		// Default: "userID"
		ContextKey string

		// ErrorHandler defines a function which is executed when an error occurs
		ErrorHandler func(c echo.Context, err error) error
	}

	ClerkJWKS struct {
		Keys []ClerkJWK `json:"keys"`
	}

	ClerkJWK struct {
		Kid string `json:"kid"`
		Kty string `json:"kty"`
		Alg string `json:"alg"`
		Use string `json:"use"`
		N   string `json:"n"`
		E   string `json:"e"`
	}

	ClerkClaims struct {
		Sub string `json:"sub"`
		jwt.RegisteredClaims
	}

	jwksCache struct {
		jwks      *ClerkJWKS
		expiresAt time.Time
		mu        sync.RWMutex
	}
)

var (
	// DefaultClerkAuthConfig is the default ClerkAuth middleware config
	DefaultClerkAuthConfig = ClerkAuthConfig{
		Skipper:    echomiddleware.DefaultSkipper,
		ContextKey: "userID",
	}

	cache = &jwksCache{}
)

// ClerkAuth returns a ClerkAuth middleware with config
func ClerkAuth(config ClerkAuthConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = DefaultClerkAuthConfig.Skipper
	}
	if config.ContextKey == "" {
		config.ContextKey = DefaultClerkAuthConfig.ContextKey
	}
	if config.JWKSUrl == "" {
		panic("clerk auth middleware requires JWKSUrl")
	}
	if config.ErrorHandler == nil {
		config.ErrorHandler = func(c echo.Context, err error) error {
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		}
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			// Extract token from Authorization header
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return config.ErrorHandler(c, errors.New("missing authorization header"))
			}

			// Check if the header has the Bearer prefix
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				return config.ErrorHandler(c, errors.New("invalid authorization header format"))
			}

			tokenString := parts[1]

			// Parse the token without validation to get the kid from header
			token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &ClerkClaims{})
			if err != nil {
				return config.ErrorHandler(c, errors.New("invalid token format"))
			}

			// Get the kid from token header
			kid, ok := token.Header["kid"].(string)
			if !ok {
				return config.ErrorHandler(c, errors.New("missing kid in token header"))
			}

			// Fetch JWKS from Clerk (with caching)
			jwks, err := fetchJWKSWithCache(config.JWKSUrl)
			if err != nil {
				return config.ErrorHandler(c, fmt.Errorf("failed to fetch JWKS: %w", err))
			}

			// Find the key with matching kid
			var publicKey *rsa.PublicKey
			for _, key := range jwks.Keys {
				if key.Kid == kid {
					publicKey, err = convertJWKToPublicKey(key)
					if err != nil {
						return config.ErrorHandler(c, fmt.Errorf("failed to parse public key: %w", err))
					}
					break
				}
			}

			if publicKey == nil {
				return config.ErrorHandler(c, errors.New("no matching key found"))
			}

			// Parse and validate the token with the public key
			claims := &ClerkClaims{}
			parsedToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return publicKey, nil
			})

			if err != nil || !parsedToken.Valid {
				return config.ErrorHandler(c, errors.New("invalid token"))
			}

			// Extract user UUID from sub claim
			userID := claims.Sub
			if userID == "" {
				return config.ErrorHandler(c, errors.New("missing sub claim in token"))
			}

			// Store user ID in context for handlers to use
			c.Set(config.ContextKey, userID)

			return next(c)
		}
	}
}

// ClerkAuthWithConfig returns a ClerkAuth middleware with config
func ClerkAuthWithConfig(jwksURL string) echo.MiddlewareFunc {
	config := DefaultClerkAuthConfig
	config.JWKSUrl = jwksURL
	return ClerkAuth(config)
}

func fetchJWKSWithCache(jwksURL string) (*ClerkJWKS, error) {
	cache.mu.RLock()
	if cache.jwks != nil && time.Now().Before(cache.expiresAt) {
		jwks := cache.jwks
		cache.mu.RUnlock()
		return jwks, nil
	}
	cache.mu.RUnlock()

	cache.mu.Lock()
	defer cache.mu.Unlock()

	// Double-check after acquiring write lock
	if cache.jwks != nil && time.Now().Before(cache.expiresAt) {
		return cache.jwks, nil
	}

	jwks, err := fetchJWKS(jwksURL)
	if err != nil {
		return nil, err
	}

	cache.jwks = jwks
	cache.expiresAt = time.Now().Add(1 * time.Hour)

	return jwks, nil
}

func fetchJWKS(jwksURL string) (*ClerkJWKS, error) {
	resp, err := http.Get(jwksURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch JWKS")
	}

	var jwks ClerkJWKS
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, err
	}

	return &jwks, nil
}

func convertJWKToPublicKey(jwk ClerkJWK) (*rsa.PublicKey, error) {
	// Decode the modulus (n)
	nBytes, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		return nil, err
	}

	// Decode the exponent (e)
	eBytes, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		return nil, err
	}

	// Convert exponent bytes to int
	var eInt int
	for _, b := range eBytes {
		eInt = eInt<<8 + int(b)
	}

	return &rsa.PublicKey{
		N: new(big.Int).SetBytes(nBytes),
		E: eInt,
	}, nil
}
