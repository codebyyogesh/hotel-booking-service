package db

import (
	"context"

	"github.com/codebyyogesh/hotel-booking-service/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const hotelCollection = "hotels"

type HotelStore interface{
    InsertHotel(context.Context, *types.Hotel)(*types.Hotel, error)
    UpdateHotel(context.Context, bson.M, bson.M)(error)
    GetHotels(context.Context, bson.M, *Pagination)(*[]types.Hotel, error)
    GetHotelByID(context.Context, string) (*types.Hotel, error)
}

type MongoHotelStore struct{
    client *mongo.Client
    collection *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client) *MongoHotelStore{
    return &MongoHotelStore{
        client: client,
        collection: client.Database(DBNAME).Collection(hotelCollection),
    }
}

func (s *MongoHotelStore)InsertHotel(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error){
    result, err := s.collection.InsertOne(ctx, hotel)
    if err != nil{
        return nil, err
    }
    hotel.ID = result.InsertedID.(primitive.ObjectID)
    return hotel, nil
}

func (s *MongoHotelStore)UpdateHotel(ctx context.Context, 
                                    filter bson.M, 
                                    update bson.M) error{

    _, err := s.collection.UpdateOne(ctx, filter, update)
    if err != nil{
        return err
    }
    return nil
}

func (s *MongoHotelStore)GetHotelByID(ctx context.Context,
                                      id string)(*types.Hotel, error){
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil{
        return nil, err
    }
    var hotel types.Hotel
    if err = s.collection.FindOne(ctx, bson.M{"_id":oid}).Decode(&hotel); err!= nil{
        return nil, err
    }
    return &hotel, nil
}

    func (s *MongoHotelStore)GetHotels(ctx context.Context, filter bson.M, p *Pagination)(*[]types.Hotel, error){
        opts := options.FindOptions{}
        if p != nil{
            opts.SetSkip(int64(p.Page-1) * int64(p.Limit))
            opts.SetLimit(int64(p.Limit))
        }
        cursor, err:= s.collection.Find(ctx, filter, &opts); 
        if err != nil{
            return nil, err
        }
        var hotels []types.Hotel
        if err = cursor.All(ctx, &hotels); err != nil {
            return nil, err
        }
    return &hotels, nil
}

