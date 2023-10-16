package api

import (
	"fmt"
	"os"

	"github.com/codebyyogesh/hotel-booking-service/db"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)


func JWTAuthentication(userStore db.UserStore) fiber.Handler{
    return func(c *fiber.Ctx) error{
        token, ok := c.GetReqHeaders()["X-Api-Token"] // ok=true, if map key exists, else false
        if !ok{
            fmt.Println("token not present in the header")
            return ErrorUnAuthorized()
        }
        claims, err := validateToken(token)
        if err != nil{
            return err
        }
        err = claims.Valid() // checks for time expiration of the token.
        if err != nil{
            return fmt.Errorf("token expired")
        } 
        userID := claims["id"].(string)
        user, err := userStore.GetUserByID(c.Context(), userID)
        if err != nil{
           return fmt.Errorf("unauthorized")
        }
        c.Context().SetUserValue("user", user)
        return c.Next() // move to the next part
    }
}

func validateToken(tokenString string) (jwt.MapClaims, error) {

    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        // Don't forget to validate the alg is what you expect:
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            fmt.Printf("Invalid Signing Method: %v", token.Header["alg"])
            return nil, fmt.Errorf("Unauthorized")
        }
        secret := os.Getenv("JWT_SECRET")
        // hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
        return []byte(secret), nil
    })

    if err != nil{
        fmt.Println("Failed to parse JWT token:", err)
        return nil, fmt.Errorf("Unauthorized")
    }
    if !token.Valid {
        fmt.Println("invalid token")
        return nil, fmt.Errorf("Unauthorized")
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok{
        return nil, fmt.Errorf("Unauthorized")
    }
    return claims, nil 
}