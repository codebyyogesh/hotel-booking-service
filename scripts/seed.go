package main

import (
	"context"
	"fmt"
	"log"

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
)

func hotelSeed(name string, location string, rating int){

    hotel := types.Hotel{
        Name: name,
        Location: location,
        Rooms: []primitive.ObjectID{},
        Rating: rating,
    }
    rooms := []types.Room{
        {
            Type: types.SingleRoomType,
            Price: 9999.9,
        },
        {
            Type : types.DoubleRoomType,
            Price: 14999.9,
        },
        {
            Type: types.DeluxeRoomType,
            Price: 19999.9,
        },
    }
    fmt.Println("Seeding the db...")
    insertedHotel, err := hotelStore.InsertHotel(ctx, &hotel)
    if err != nil{
        log.Fatal(err)
    }

    for _, room := range rooms {
        room.HotelID = insertedHotel.ID
        _,err := roomStore.InsertRoom(ctx, &room)
        if err != nil{
            log.Fatal(err)
        }
    }

}
func main(){
    hotelSeed("The Taj", "Mumbai", 4)
    hotelSeed("The Leela Palace", "Bengaluru", 3)
    hotelSeed("Kaldan Samudra", "Mahabalipuram", 5)
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
}