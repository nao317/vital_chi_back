package auth

import (
    "fmt"
    "log"
    "os"

    "github.com/MicahParks/keyfunc/v2"
    "github.com/labstack/echo/v4"
    echojwt "github.com/labstack/echo-jwt/v4"
)

func SupabaseJWTMiddleware() echo.MiddlewareFunc {
    projectRef := os.Getenv("SUPABASE_PROJECT_REF")
    if projectRef == "" {
        log.Fatal("SUPABASE_PROJECT_REF environment variable is required")
    }

    jwksURL := fmt.Sprintf("https://%s.supabase.co/auth/v1/.well-known/jwks.json", projectRef)

    jwks, err := keyfunc.Get(jwksURL, keyfunc.Options{})
    if err != nil {
        log.Fatalf("failed to get the JWKS from %s: %v", jwksURL, err)
    }

    return echojwt.WithConfig(echojwt.Config{
        KeyFunc: jwks.Keyfunc,
        ContextKey: "user",
        ErrorHandler: func(c echo.Context, err error) error {
            return echo.NewHTTPError(401, "invalid token")
        },
    })
}
