package api

import (
	"errors"
	"fmt"

	"github.com/codebyyogesh/hotel-booking-service/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct{
    userStore db.UserStore // is an interface
}
type AuthParams struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

func NewAuthHandler(userstore db.UserStore ) *AuthHandler{
    return &AuthHandler{
        userStore: userstore,
    }
}

func (h *AuthHandler)HandleAuthenticate(c *fiber.Ctx) error{
    var params AuthParams
    if err := c.BodyParser(&params); err != nil{
        return err
    }
    user, err := h.userStore.GetUserByEmail(c.Context(), params.Email)
    if err != nil{
        if errors.Is(err, mongo.ErrNoDocuments){
            return fmt.Errorf("invalid credentials")
        }
        return err
    }
    err =  bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(params.Password))
    if err != nil{
        return fmt.Errorf("invalid credentials")
    }
    fmt.Println("authenticated ->:",user)
    return nil
} 