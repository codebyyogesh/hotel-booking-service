package main

import (
	"flag"
	api "github.com/codebyyogesh/hotel-booking-service/api"
	"github.com/gofiber/fiber/v2"
)


func main(){
    listenAddr := flag.String("listenAddr", ":5000", "Listen address of the API server")
    flag.Parse()
    app := fiber.New()
    apiv1 := app.Group("/api/v1")      // /api/v1
    apiv1.Get("/user", api.HandleGetUsers)
    apiv1.Get("/user:id", api.HandleGetUser)
    app.Listen(*listenAddr)
}
