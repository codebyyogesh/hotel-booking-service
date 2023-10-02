package middleware

import (
	"fmt"

	"github.com/codebyyogesh/hotel-booking-service/types"
	"github.com/gofiber/fiber/v2"
)


func  AdminAuth(c *fiber.Ctx) error{
    
   user,ok := c.Context().UserValue("user").(*types.User)
   if !ok{
       return fmt.Errorf("not authorized")
   }
   if !user.IsAdmin{
       return fmt.Errorf("not authorized")
   }
   return c.Next()
}