package main

import (
	"context"
	"flag"
	"log"
	"os"

	api "github.com/codebyyogesh/hotel-booking-service/api"
	env "github.com/codebyyogesh/hotel-booking-service/config"
	"github.com/codebyyogesh/hotel-booking-service/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var config = fiber.Config{
    ErrorHandler: api.ErrorHandler,
}


func main(){
    mongoDB_EndPoint:= os.Getenv("MONGO_DB_URL")
    client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoDB_EndPoint))
    if err != nil {
        log.Fatal(err)
    }

    flag.Parse()

    // create new user handler - handler initialization after db initialization.
    // PS: mongoUserStore is a pointer, but satisfies the interface and 
    // can be passed as param to NewUserHandler which accepts interface
    // as param. i.e. *mongoUserStore implements the UserStore interface.
    // Interfaces even works for the pointers.
    var(
        userStore      = db.NewMongoUserStore(client)
        hotelStore     = db.NewMongoHotelStore(client)
        roomStore      = db.NewMongoRoomStore(client, hotelStore)
        bookingStore   = db.NewMongoBookingStore(client)
        store          = &db.Store{
                        Hotels: hotelStore,
                        Rooms: roomStore,
                        User: userStore,
                        Booking: bookingStore,
        }
        userHandler    = api.NewUserHandler(userStore)
        hotelHandler   = api.NewHotelHandler(store)
        authHandler    = api.NewAuthHandler(userStore)
        roomHandler    = api.NewRoomHandler(store)
        bookingHandler = api.NewBookingHandler(store)
        app            = fiber.New(config)
        auth           = app.Group("/api")
        apiv1          = app.Group("/api/v1", api.JWTAuthentication(userStore))      // /api/v1
        admin          = apiv1.Group("/admin", api.AdminAuth)
    )
    // auth handlers
    auth.Post("/auth", authHandler.HandleAuthenticate)

    // Versioned API routes
    // user handlers
    apiv1.Post("/user", userHandler.HandlePostUser)
    apiv1.Get("/user", userHandler.HandleGetUsers)
    apiv1.Get("/user/:id", userHandler.HandleGetUser)
    apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
    apiv1.Put("/user/:id", userHandler.HandlePutUser)

    // hotel handlers
    apiv1.Get("/hotel", hotelHandler.HandleGetHotels)
    apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
    apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)

    // room handlers
    apiv1.Get("/room", roomHandler.HandleGetRooms)
    apiv1.Post("/room/:id/book", roomHandler.HandleBookRoom)

    // booking handler
    admin.Get("/booking", bookingHandler.HandleGetBookings) // admin
    apiv1.Get("/booking/:id", bookingHandler.HandleGetBooking) // user

    // cancel booking
    apiv1.Get("/booking/:id/cancel", bookingHandler.HandleCancelBooking)

    listenAddr := os.Getenv("HTTP_LISTEN_ADDRESS")
    app.Listen(listenAddr)
}


func init() {
    env.LoadEnv()
}