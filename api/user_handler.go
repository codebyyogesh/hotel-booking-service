package api

import (
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

func (h *UserHandler)HandlePutUser(c *fiber.Ctx) error{
    var (
        UserID = c.Params("id")  // user ID = in json it is "id"
        ctx = c.Context()
    )
    var params types.UpdateUserParams

    if err := c.BodyParser(&params); err != nil {
        return err
    }
    h.userStore.UpdateUser(ctx, UserID,  &params)
    return c.JSON(map[string]string{"Updated": UserID})
}

func (h *UserHandler)HandleDeleteUser(c *fiber.Ctx) error{
    var (
        UserID = c.Params("id")  // user ID = in json it is "id"
        ctx = c.Context()
    )
    if err := h.userStore.DeleteUser(ctx, UserID); err != nil{
        return err
    }
    return c.JSON(map[string]string{"Deleted": UserID})
}

// post handler
func (h *UserHandler)HandlePostUser(c *fiber.Ctx) error{
    var params types.CreateUserParams
    if err:= c.BodyParser(&params); err != nil{  // parse into ValidateParams
        return err
    }
    // Validate User before using it
    if errors := params.ValidateUserParams(); len(errors) > 0{
        return c.JSON(errors)
    }
    // Now we have the correct user after validation
    user, err := types.NewUserFromParams(params)
    if err != nil{
        return err
    }
    // insert into db
    insertedUser, err := h.userStore.InsertUser(c.Context(), user)
    if err != nil{
        return err
    }
    return c.JSON(insertedUser)
}

// get a user by id 
func (h *UserHandler)HandleGetUser(c *fiber.Ctx) error{
    var (
        id = c.Params("id")  // user ID = in json it is "id"
        ctx = c.Context()
    )
    user, err := h.userStore.GetUserByID(ctx, id)
    if err != nil{
        return err
    }
    return c.JSON(user)
}

// get users list or all users
func (h *UserHandler)HandleGetUsers(c *fiber.Ctx) error{
    users, err := h.userStore.GetUsers(c.Context())
    if err != nil{
        return err
    }
    return c.JSON(users)
}
