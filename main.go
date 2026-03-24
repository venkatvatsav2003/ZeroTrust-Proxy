package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	JWTSecret       = "super-secret-key-change-in-production"
	TargetBackend   = "http://localhost:8080"
	ProxyPort       = ":8443"
)

func main() {
	target, err := url.Parse(TargetBackend)
	if err != nil {
		log.Fatalf("Invalid target URL: %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	
	// Custom Director to modify the request before sending to backend
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		// Internal services shouldn't trust the client IP directly
		req.Header.Set("X-Forwarded-Host", req.Header.Get("Host"))
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", ZeroTrustMiddleware(proxy))

	log.Printf("Starting Zero-Trust Proxy on port %s forwarding to %s\n", ProxyPort, TargetBackend)
	// In a real scenario, use ListenAndServeTLS for mTLS
	log.Fatal(http.ListenAndServe(ProxyPort, mux))
}

// ZeroTrustMiddleware enforces JWT authentication and RBAC
func ZeroTrustMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized: Missing or invalid Authorization header", http.StatusUnauthorized)
			log.Printf("Blocked request from %s: Missing token", r.RemoteAddr)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(JWTSecret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			log.Printf("Blocked request from %s: Invalid token (%v)", r.RemoteAddr, err)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			role := claims["role"]
			// Basic RBAC Example
			if role != "admin" && role != "service" {
				http.Error(w, "Forbidden: Insufficient permissions", http.StatusForbidden)
				log.Printf("Blocked request from %s: Insufficient permissions for role %v", r.RemoteAddr, role)
				return
			}
			// Inject validated identity context into headers for the backend service
			r.Header.Set("X-Authenticated-User", fmt.Sprintf("%v", claims["sub"]))
			r.Header.Set("X-User-Role", fmt.Sprintf("%v", role))
		} else {
			http.Error(w, "Unauthorized: Invalid claims", http.StatusUnauthorized)
			return
		}

		log.Printf("Allowed request from %s to %s (Time: %s)", r.RemoteAddr, r.URL.Path, time.Since(start))
		next.ServeHTTP(w, r)
	}
}
