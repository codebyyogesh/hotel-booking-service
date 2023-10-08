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
    *db.Store
}

func (tdb *testdb) teardown(t *testing.T) {
    if err := tdb.client.Database(db.DBNAME).Drop(context.TODO()); err != nil {
        t.Fatal(err)
    }
}

func setup(t *testing.T) *testdb {
    client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.DBURI))
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
