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

var config = fiber.Config{
    ErrorHandler: func(c *fiber.Ctx, err error) error {
        return c.JSON(map[string]string{"error": err.Error()})
    },
}


func main(){
    client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
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
    var(
        mongoUserStore = db.NewMongoUserStore(client, db.DBNAME)
        userHandler    = api.NewUserHandler(mongoUserStore)
    )

    var(
        app   = fiber.New(config)
        apiv1 = app.Group("/api/v1")      // /api/v1
    )

    // user handlers
    apiv1.Post("/user", userHandler.HandlePostUser)
    apiv1.Get("/user", userHandler.HandleGetUsers)
    apiv1.Get("/user/:id", userHandler.HandleGetUser)
    apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
    apiv1.Put("/user/:id", userHandler.HandlePutUser)

    var(
        hotelStore   = db.NewMongoHotelStore(client)
        roomStore    = db.NewMongoRoomStore(client, hotelStore)
        hotelHandler = api.NewHotelHandler(hotelStore, roomStore)
    )
    // hotel handlers
    apiv1.Get("/hotel", hotelHandler.HandleGetHotels)


    app.Listen(*listenAddr)
}
