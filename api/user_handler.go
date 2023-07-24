package api

import (
	"context"
	"github.com/codebyyogesh/hotel-booking-service/db"
	"github.com/codebyyogesh/hotel-booking-service/types"
	"github.com/gofiber/fiber/v2"
)
type UserHandler struct{
    userStore db.UserStore // is an interface
}

func NewUserHandler(userstore db.UserStore ) *UserHandler{
    return &UserHandler{
        userStore: userstore,
    }
} 
// get a user by id 
func (h *UserHandler)HandleGetUser(c *fiber.Ctx) error{
    var (
        id = c.Params("id")  // user ID = in json it is "id"
        ctx = context.Background()
    )
    user, err := h.userStore.GetUserByID(ctx, id)
    if err != nil{
        return err
    }
    return c.JSON(user)
}

// get users list or all users
func (h *UserHandler)HandleGetUsers(c *fiber.Ctx) error{
    u := types.User{
        FirstName: "Air",
        LastName: "Condition",
    }
    return c.JSON(u)
}
