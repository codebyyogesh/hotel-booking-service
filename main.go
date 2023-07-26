package main

import (
	"context"
	"flag"
	"log"

	api "github.com/codebyyogesh/hotel-booking-service/api"
	"github.com/codebyyogesh/hotel-booking-service/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dburi = "mongodb://localhost:27017"

func main(){
    client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
    if err != nil {
        log.Fatal(err)
    }

    listenAddr := flag.String("listenAddr", ":5000", "Listen address of the API server")
    flag.Parse()

    // create new user handler - handler initialization after db initialization.
    // PS: mongoUserStore is a pointer, but satisfies the interface and 
    // can be passed as param to NewUserHandler which accepts interface
    // as param. i.e. *mongoUserStore implements the UserStore interface.
    // Interfaces even works for the pointers.
    mongoUserStore := db.NewMongoUserStore(client)
    userHandler := api.NewUserHandler(mongoUserStore)

    app := fiber.New()
    apiv1 := app.Group("/api/v1")      // /api/v1

    apiv1.Post("/user", userHandler.HandlePostUser)
    apiv1.Get("/user", userHandler.HandleGetUsers)
    apiv1.Get("/user/:id", userHandler.HandleGetUser)
    app.Listen(*listenAddr)
}
