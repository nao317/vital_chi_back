package main

import (
    "fmt"
    "net/http"
    "os"
    "vital_chi_back/auth"

    "github.com/golang-jwt/jwt/v5"
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
)

func main() {
    e := echo.New()

    // ミドルウェアの設定
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())
    e.Use(middleware.CORS())

    // パブリックなルート
    e.GET("/", func(c echo.Context) error {
        return c.JSON(http.StatusOK, map[string]string{
            "message": "Hello, Vital Chi API!",
        })
    })

    // ヘルスチェックエンドポイント
    e.GET("/health", func(c echo.Context) error {
        return c.JSON(http.StatusOK, map[string]string{
            "status": "OK",
        })
    })

    // テスト用認証状態確認エンドポイント（開発時のみ）
    e.GET("/auth-test", func(c echo.Context) error {
        authHeader := c.Request().Header.Get("Authorization")
        if authHeader == "" {
            return c.JSON(http.StatusOK, map[string]interface{}{
                "message": "No Authorization header provided",
                "has_auth": false,
            })
        }
        
        headerPreview := authHeader
        if len(authHeader) > 50 {
            headerPreview = authHeader[:50] + "..."
        }
        
        return c.JSON(http.StatusOK, map[string]interface{}{
            "message": "Authorization header detected",
            "has_auth": true,
            "header": headerPreview,
        })
    })

    // Supabase認証情報確認エンドポイント（開発時のみ）
    e.GET("/supabase-info", func(c echo.Context) error {
        projectRef := os.Getenv("SUPABASE_PROJECT_REF")
        jwksURL := fmt.Sprintf("https://%s.supabase.co/auth/v1/.well-known/jwks.json", projectRef)
        
        return c.JSON(http.StatusOK, map[string]interface{}{
            "project_ref": projectRef,
            "jwks_url": jwksURL,
            "auth_url": fmt.Sprintf("https://%s.supabase.co/auth/v1", projectRef),
            "message": "Supabaseの認証画面でユーザーを作成し、JWTトークンを取得してください",
        })
    })

    // 認証が必要なルート
    e.GET("/protected", func(c echo.Context) error {
        user := c.Get("user").(*jwt.Token)
        claims, ok := user.Claims.(jwt.MapClaims)
        if !ok {
            return echo.NewHTTPError(http.StatusInternalServerError, "failed to parse token claims")
        }

        sub, ok := claims["sub"].(string)
        if !ok {
            return echo.NewHTTPError(http.StatusInternalServerError, "invalid user ID in token")
        }

        return c.JSON(http.StatusOK, map[string]interface{}{
            "message": "success",
            "user_id": sub,
        })
    }, auth.SupabaseJWTMiddleware())

    e.Logger.Fatal(e.Start(":8080"))
}