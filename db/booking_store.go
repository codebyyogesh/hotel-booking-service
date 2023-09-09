package db

import (
	"context"

	"github.com/codebyyogesh/hotel-booking-service/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const bookingCollection ="bookings"

type BookingStore interface{
    InsertBooking(context.Context, *types.Booking)(*types.Booking, error)
    GetBookings(context.Context, bson.M ) (*[]types.Booking, error) // *[]types.Booking will be a pointer to a slice
}

type MongoBookingStore struct{
    client *mongo.Client
    collection *mongo.Collection
}

func NewMongoBookingStore(client *mongo.Client) *MongoBookingStore{
    return &MongoBookingStore{
        client: client,
        collection: client.Database(DBNAME).Collection(bookingCollection),
    }
}

func (s *MongoBookingStore)InsertBooking(ctx context.Context, booking *types.Booking) (*types.Booking, error){
    result, err := s.collection.InsertOne(ctx, booking)
    if err != nil{
        return nil, err
    }
    booking.ID = result.InsertedID.(primitive.ObjectID)
    return booking, nil
}

func (s *MongoBookingStore)GetBookings(ctx context.Context, filter bson.M)(*[]types.Booking, error){
    cursor, err:= s.collection.Find(ctx, filter); 
    if err != nil{
        return nil, err
    }
    var bookings []types.Booking
    if err = cursor.All(ctx, &bookings); err != nil {
        return nil, err
    }
    if len(bookings) == 0{ // if there are no bookings
        return nil, nil
    }
    return &bookings, nil
}