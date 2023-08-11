package db

import (
	"context"

	"github.com/codebyyogesh/hotel-booking-service/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const roomCollection ="rooms"

type RoomStore interface{
    InsertRoom(context.Context, *types.Room)(*types.Room, error)
    GetRooms(context.Context, string)(*[]types.Room, error)
}

type MongoRoomStore struct{
    client *mongo.Client
    collection *mongo.Collection
    HotelStore
}

func NewMongoRoomStore(client *mongo.Client, hotelStore HotelStore) *MongoRoomStore{
    return &MongoRoomStore{
        client: client,
        collection: client.Database(DBNAME).Collection(roomCollection),
        HotelStore: hotelStore,
    }
}

func (s *MongoRoomStore)GetRooms(ctx context.Context, id string) (*[]types.Room, error){
    // Mongodb does not accept direct ids, we need some kind of conversion to object id
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil{
        return nil, err
    }
    filter := bson.M{"hotelID": oid}
    cursor, err := s.collection.Find(ctx, filter)
    if err != nil{
        return nil, err
    }
    var rooms []types.Room
    if err := cursor.All(ctx, &rooms); err != nil{
        return nil, err
    }
    return &rooms, nil
}

func (s *MongoRoomStore)InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error){
    result, err := s.collection.InsertOne(ctx, room)
    if err != nil{
        return nil, err
    }
    room.ID = result.InsertedID.(primitive.ObjectID)

    // Update the hotel with this room ID
    filter := bson.M{"_id": room.HotelID}
    update := bson.M{"$push": bson.M{"rooms": room.ID}}
    if s.UpdateHotel(ctx, filter, update); err != nil{
        return nil, err
    }
    return room, nil
}