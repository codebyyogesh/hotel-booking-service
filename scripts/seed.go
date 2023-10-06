package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/codebyyogesh/hotel-booking-service/api"
	"github.com/codebyyogesh/hotel-booking-service/db"
	"github.com/codebyyogesh/hotel-booking-service/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
    fmt.Println("Seeding the user to db...")
    user := seedUser(false, "raju", "gentleman", "raju@me.com", "mybestsecurepassword") //regular user
    seedUser(true, "admin", "admin", "admin@me.com", "admin123") // admin
    fmt.Println("Seeding the hotel to db...")
    seedHotel("The Leela Palace", "Bengaluru", 3)
    seedHotel("Kaldan Samudra", "Mahabalipuram", 5)
    hotel := seedHotel("The Taj", "Mumbai", 4)
    fmt.Println("Seeding the room to db...")
    seedRoom("small", false, 9999.9, hotel.ID)
    seedRoom("normal", true, 14999.9, hotel.ID)
    room :=seedRoom("kingsize", true, 19999.9, hotel.ID)
    fmt.Println("Seeding the booking to db...")
    booking := seedBooking(user.ID, room.ID, 2, time.Now(), time.Now().AddDate(0, 0, 2))
    fmt.Println("booking id:", booking.ID)
} 

// special function gets automatically called before main()
func init(){
    var err error
    client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
    if err != nil {
        log.Fatal(err)
    }
    if client.Database(db.DBNAME).Drop(ctx); err != nil{
        log.Fatal(err)
    }
    hotelStore = db.NewMongoHotelStore(client)
    roomStore = db.NewMongoRoomStore(client, hotelStore)
    userStore = db.NewMongoUserStore(client)
    bookingStore = db.NewMongoBookingStore(client)
}