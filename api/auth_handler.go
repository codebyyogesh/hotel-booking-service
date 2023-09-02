package api

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/codebyyogesh/hotel-booking-service/db"
	"github.com/codebyyogesh/hotel-booking-service/types"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthHandler struct{
    userStore db.UserStore // is an interface
}

type AuthParams struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type AuthResponse struct{
    User *types.User    `json:"user"`
    Token string        `json:"token"`
}

type genericResp struct {
    Type  string `json:"type"`
    Msg   string `json:"msg"`
}

func NewAuthHandler(userstore db.UserStore ) *AuthHandler{
    return &AuthHandler{
        userStore: userstore,
    }
}
func invalidCredentialsError(c *fiber.Ctx) error{
    return c.Status(http.StatusBadRequest).JSON(genericResp{
        Type: "error",
        Msg:  "invalid credentials",
    })
}

func (h *AuthHandler)HandleAuthenticate(c *fiber.Ctx) error{
    var params AuthParams
    if err := c.BodyParser(&params); err != nil{
        return err
    }
    user, err := h.userStore.GetUserByEmail(c.Context(), params.Email)
    if err != nil{
        if errors.Is(err, mongo.ErrNoDocuments){
            return invalidCredentialsError(c)
        }
        return err
    }
    if !types.IsPasswordValid(user.EncryptedPassword, params.Password){
        return invalidCredentialsError(c)
    }
    resp := AuthResponse{
        User: user,
        Token: createTokenFromUser(user),
    }
    return c.JSON(resp)
}

func createTokenFromUser(user *types.User) string{
    now := time.Now()
    expires := now.Add(time.Hour * 2).Unix() // valid for two hours
    // Create a new token object, specifying signing method and the claims
    // you would like it to contain.
    claims := jwt.MapClaims{
        "id": user.ID,
        "email": user.Email,
        "exp": expires, // exp or nbf is the keyword to be used if we want to use jwt.MapClaims Api
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    // Sign and get the complete encoded token as a string using the secret
    secret := os.Getenv("JWT_SECRET")
    tokenString, err := token.SignedString([]byte(secret))
    if err != nil{
        fmt.Println("failed to sign token with secret")
    }
    return tokenString
}
