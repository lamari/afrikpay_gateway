package middleware

import (
    "net/http"
    "os"

    "github.com/gorilla/mux"
)

// AuthMiddleware validates JWT by delegating to Auth service /auth/verify endpoint.
// If AUTH_URL env var is empty, middleware is disabled (passes through).
func AuthMiddleware() mux.MiddlewareFunc {
    authURL := os.Getenv("AUTH_URL")
    if authURL == "" {
        // disabled
        return func(next http.Handler) http.Handler {
            return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                next.ServeHTTP(w, r)
            })
        }
    }
    verifyURL := authURL + "/auth/verify"

    client := &http.Client{}
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                http.Error(w, "missing authorization header", http.StatusUnauthorized)
                return
            }
            req, _ := http.NewRequest(http.MethodGet, verifyURL, nil)
            req.Header.Set("Authorization", authHeader)
            resp, err := client.Do(req)
            if err != nil || resp.StatusCode != http.StatusOK {
                http.Error(w, "unauthorized", http.StatusUnauthorized)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
