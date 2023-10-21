package api

import (
	"context"
	"log"
	"testing"

	"github.com/codebyyogesh/hotel-booking-service/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testdb struct {
    client *mongo.Client
    *db.Store // note: We use the same db that is used for the main program because NewMongoHotelStore() etc in setup
            // uses db.DBNAME = "hotel-booking" which is the same as the main program.
}

func (tdb *testdb) teardown(t *testing.T) {
    dbname := db.DBNAME
    if err := tdb.client.Database(dbname).Drop(context.TODO()); err != nil {
        t.Fatal(err)
    }
}

func setup(t *testing.T) *testdb {
    dburi  := db.DBURI
    client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dburi))
    if err != nil {
        log.Fatal(err)
    }
    hotelStore := db.NewMongoHotelStore(client)
    roomStore := db.NewMongoRoomStore(client, hotelStore)
    userStore := db.NewMongoUserStore(client)
    bookingStore := db.NewMongoBookingStore(client)

    return &testdb{
        client: client,
        Store: &db.Store{
            Hotels:   hotelStore,
            Rooms:   roomStore,
            User: userStore,
            Booking: bookingStore,
        },
    }
}
