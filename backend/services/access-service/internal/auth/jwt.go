package auth

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"math/big"
	"net/http"
	"strings"
	"sync"
)

type Jwks struct {
	Keys []struct {
		Kty string `json:"kty"`
		Kid string `json:"kid"`
		N   string `json:"n"`
		E   string `json:"e"`
	} `json:"keys"`
}

type ctxKey string

var userInfoKey ctxKey = "userInfo"

type UserInfo struct {
	Auth0ID string
	Email   string
}

var (
	auth0Domain string
	cachedKeys  = map[string]*rsa.PublicKey{}
	mu          sync.Mutex
)

func SetAuth0Domain(domain string) {
	auth0Domain = domain
	log.Println("Auth0 Domain set in auth package:", auth0Domain)
}

func getPublicKey(jwksURL string, token *jwt.Token) (*rsa.PublicKey, error) {
	mu.Lock()
	defer mu.Unlock()

	kid, ok := token.Header["kid"].(string)
	if !ok {
		return nil, errors.New("kid header not found in token")
	}

	if key, found := cachedKeys[kid]; found {
		return key, nil
	}

	resp, err := http.Get(jwksURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var jwks Jwks
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, err
	}

	for _, key := range jwks.Keys {
		if key.Kid == kid {
			pubKey, err := parseRSAPublicKey(key.N, key.E)
			if err != nil {
				return nil, err
			}
			cachedKeys[kid] = pubKey
			return pubKey, nil
		}
	}

	return nil, errors.New("matching key not found")
}

func parseRSAPublicKey(nStr, eStr string) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(nStr)
	if err != nil {
		return nil, err
	}
	eBytes, err := base64.RawURLEncoding.DecodeString(eStr)
	if err != nil {
		return nil, err
	}

	n := new(big.Int).SetBytes(nBytes)
	e := int(new(big.Int).SetBytes(eBytes).Int64())

	return &rsa.PublicKey{N: n, E: e}, nil
}

func JwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		claims := struct {
			Sub   string `json:"sub"`
			Email string `json:"email"`
			jwt.RegisteredClaims
		}{}

		token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, errors.New("invalid token signature")
			}
			return getPublicKey(auth0Domain+"/.well-known/jwks.json", token)
		})

		if err != nil || !token.Valid {
			log.Println("JWT verification failed:", err)
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		if claims.Sub == "" || claims.Email == "" {
			http.Error(w, "Invalid token: missing claims", http.StatusUnauthorized)
			return
		}

		user := &UserInfo{
			Auth0ID: claims.Sub,
			Email:   claims.Email,
		}

		log.Println("JWT verified successfully for user:", user.Auth0ID)

		ctx := context.WithValue(r.Context(), userInfoKey, user)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func GetUserInfo(ctx context.Context) (*UserInfo, error) {
	user, ok := ctx.Value(userInfoKey).(*UserInfo)
	if !ok || user == nil {
		return nil, errors.New("user info not found in context")
	}
	return user, nil
}
