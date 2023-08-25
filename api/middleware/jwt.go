package middleware

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)


func JWTAuthentication(c *fiber.Ctx) error {
    fmt.Println("...JWT...")
    token, ok := c.GetReqHeaders()["X-Api-Token"] // ok=true, if map key exists, else false
    if !ok{
        return fmt.Errorf("unauthorized")
    }
    if err := parseToken(token); err != nil{
        return err
    }
    fmt.Println("token:", token)
    return nil
}

func parseToken(tokenString string) error {

    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
    // Don't forget to validate the alg is what you expect:
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            fmt.Printf("Invalid Signing Method: %v", token.Header["alg"])
            return nil, fmt.Errorf("Unauthorized")
        }
        secret := os.Getenv("JWT_SECRET")
        fmt.Println("Never print the secret:",secret)
        // hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
        return []byte(secret), nil
    })
    if err != nil{
        fmt.Println("Failed to parse JWT token:", err)
        return fmt.Errorf("Unauthorized")
    }

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        fmt.Println(claims)
    }
    return fmt.Errorf("Unauthorized")
}