package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/codebyyogesh/hotel-booking-service/api"
	"github.com/codebyyogesh/hotel-booking-service/db"
	"github.com/codebyyogesh/hotel-booking-service/db/fixtures"
	"github.com/codebyyogesh/hotel-booking-service/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/exp/rand"
)

var (
    client *mongo.Client
    ctx = context.Background()
    hotelStore db.HotelStore
    roomStore db.RoomStore
    userStore db.UserStore
    bookingStore db.BookingStore
)

func seedUser(isAdmin bool, fname, lname, email, password string) (*types.User){
    user, err := types.NewUserFromParams(types.CreateUserParams{
        Email:     email,
        FirstName: fname,
        LastName:  lname,
        Password: password ,
    })
    if err != nil{
        log.Fatal(err)
    }
    user.IsAdmin = isAdmin
    insertedUser, err := userStore.InsertUser(context.Background(), user)
    if err != nil{
        log.Fatal(err)
    }
    fmt.Printf("%s -> %s\n", user.Email, api.CreateTokenFromUser(user))
    return insertedUser
}
func seedHotel(name string, location string, rating int) (*types.Hotel){

    hotel := types.Hotel{
        Name: name,
        Location: location,
        Rooms: []primitive.ObjectID{},
        Rating: rating,
    }
    insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
    if err != nil{
        log.Fatal(err)
    }
    return insertedHotel
}

func seedRoom(size string, suite bool, price float64, hotelID primitive.ObjectID) (*types.Room){
    room := types.Room{
        Size: size,
        Suite: suite,
        Price: price,
        HotelID: hotelID,
    }
    insertedRoom, err := roomStore.InsertRoom(ctx, &room)
    if err != nil{
        log.Fatal(err)
    }
    return insertedRoom
}
func seedBooking(userID, roomID primitive.ObjectID, numPersons int, from, till time.Time  ) (*types.Booking){
    booking := types.Booking{
        UserID: userID,
        RoomID: roomID,
        NumPersons: numPersons,
        FromDate: from,
        TillDate: till,
    }
    insertedBooking, err := bookingStore.InsertBooking(ctx, &booking)
    if err != nil{
        log.Fatal(err)
    }
    return insertedBooking
}
func main(){
    ctx := context.Background()
    client, err := mongo.Connect(ctx, options.Client().ApplyURI(db.DBURI))
    if err != nil {
        log.Fatal(err)
    }
    if client.Database(db.DBNAME).Drop(ctx); err != nil{
        log.Fatal(err)
    }

    hotelStore := db.NewMongoHotelStore(client)
    roomStore := db.NewMongoRoomStore(client, hotelStore)
    userStore := db.NewMongoUserStore(client)
    bookingStore := db.NewMongoBookingStore(client)
    store := &db.Store {
        Hotels:   hotelStore,
        Rooms :   roomStore,
        User :    userStore,
        Booking : bookingStore,
    }
    // seed user
    user := fixtures.CreateUser(store, "raju", "foo", false)
    fmt.Printf("%s -> %s\n", user.Email, api.CreateTokenFromUser(user))
    fmt.Printf("-----------------------------------------------------------\n")
    admin := fixtures.CreateUser(store, "admin", "admin", true) // admin
    fmt.Printf("%s -> %s\n", admin.Email, api.CreateTokenFromUser(admin))
    // seed hotel
    hotel := fixtures.CreateHotel(store, "The Leela Palace", "Bengaluru", 3, nil)
    // seed room
    fixtures.CreateRoom(store, "small", false, 9999.9, hotel.ID)
    fixtures.CreateRoom(store, "normal", true, 14999.9, hotel.ID)
    room := fixtures.CreateRoom(store, "kingsize", true, 14999.9, hotel.ID)
    // seed booking
    booking := fixtures.CreateBooking(store, user.ID, room.ID, 2, time.Now(), time.Now().AddDate(0, 0, 2))
    fmt.Printf("-----------------------------------------------------------\n")
    fmt.Println("booking ->", booking.ID)
    for i := 0; i < 100; i++ {
        name := fmt.Sprintf("hotel-%d", i)
        location := fmt.Sprintf("location-%d", i)
        fixtures.CreateHotel(store, name, location, rand.Intn(5) + 1, nil)
    }
} 
