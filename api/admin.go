package api

import (
	"github.com/codebyyogesh/hotel-booking-service/types"
	"github.com/gofiber/fiber/v2"
)


func  AdminAuth(c *fiber.Ctx) error{
    
   user,ok := c.Context().UserValue("user").(*types.User)
   if !ok{
       return ErrorUnAuthorized()
   }
   if !user.IsAdmin{
       return ErrorUnAuthorized()
   }
   return c.Next()
}