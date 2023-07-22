package api

import (
	"github.com/codebyyogesh/hotel-booking-service/types"
	"github.com/gofiber/fiber/v2"
)


func HandleGetUsers(c *fiber.Ctx) error{
    u := types.User{
        FirstName: "Air",
        LastName: "Condition",
    }
    return c.JSON(u)
}

func HandleGetUser(c *fiber.Ctx) error{
    return c.JSON(map[string]string{"Byomkesh": "Bakshi"})
}