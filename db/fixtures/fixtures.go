package fixtures

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/codebyyogesh/hotel-booking-service/db"
	"github.com/codebyyogesh/hotel-booking-service/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)


func CreateUser( store *db.Store, fname, lname string, isAdmin bool) (*types.User){
     user, err := types.NewUserFromParams(types.CreateUserParams{
        FirstName: fname,
        LastName:  lname,
        Email:     fmt.Sprintf("%s@%s.com", fname, lname),
        Password:  fmt.Sprintf("%s_%s", fname, lname),
    })
    if err != nil{
        log.Fatal(err)
    }
    user.IsAdmin = isAdmin
    insertedUser, err := store.User.InsertUser(context.Background(), user)
    if err != nil{
        log.Fatal(err)
    }
    return insertedUser
}

func CreateHotel( store *db.Store, name string, location string, rating int, rooms []primitive.ObjectID) (*types.Hotel){
    var roomIDs = rooms
    if rooms == nil {
        roomIDs = []primitive.ObjectID{}
    }
    hotel := types.Hotel{
        Name: name,
        Location: location,
        Rooms: roomIDs,
        Rating: rating,
    }
    insertedHotel, err := store.Hotels.InsertHotel(context.Background(), &hotel)
    if err != nil{
        log.Fatal(err)
    }
    return insertedHotel
}

func CreateRoom( store *db.Store, size string, suite bool, price float64, hotelID primitive.ObjectID) (*types.Room){
    room := types.Room{
        Size: size,
        Suite: suite,
        Price: price,
        HotelID: hotelID,
    }
    insertedRoom, err := store.Rooms.InsertRoom(context.Background(), &room)
    if err != nil{
        log.Fatal(err)
    }
    return insertedRoom
}

func CreateBooking(store *db.Store, userID, roomID primitive.ObjectID, numPersons int, from, till time.Time  ) (*types.Booking){
    booking := types.Booking{
        UserID: userID,
        RoomID: roomID,
        NumPersons: numPersons,
        FromDate: from,
        TillDate: till,
    }
    insertedBooking, err := store.Booking.InsertBooking(context.Background(), &booking)
    if err != nil{
        log.Fatal(err)
    }
    return insertedBooking
}