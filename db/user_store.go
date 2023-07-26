package db

import (
	"context"
	"github.com/codebyyogesh/hotel-booking-service/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const userCollection = "users"

type UserStore interface{
    GetUserByID(context.Context, string) (*types.User, error)
    GetUsers(context.Context) (*[]types.User, error) // *[]types.User will be a pointer to a slice
    // PS: []*types.User will be slice of pointers to object of type User
    InsertUser(context.Context, *types.User)(*types.User, error)
}

type MongoUserStore struct{
    client *mongo.Client
    collection *mongo.Collection
}

func NewMongoUserStore(client *mongo.Client) *MongoUserStore{
    return &MongoUserStore{
        client: client,
        collection: client.Database(DBNAME).Collection(userCollection),
    }
}

func (s *MongoUserStore)InsertUser(ctx context.Context, user *types.User ) (*types.User, error) {
    result, err := s.collection.InsertOne(ctx, user)
    if err != nil{
        return nil, err
    }
    user.ID = result.InsertedID.(primitive.ObjectID)
    return user, nil
}

func (s *MongoUserStore)GetUserByID(ctx context.Context, id string) (*types.User,error) {
    // Mongodb does not accept direct ids, we need some kind of conversion to object id
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil{
        return nil, err
    }
    var user types.User
    if err:= s.collection.FindOne(ctx, bson.M{"_id":oid}).Decode(&user); err != nil{
        return nil, err
    }
    return &user, nil
}

func (s *MongoUserStore)GetUsers(ctx context.Context) (*[]types.User, error) {
    cursor, err:= s.collection.Find(ctx, bson.M{}); 
    if err != nil{
        return nil, err
    }
    var users []types.User
    if err = cursor.All(ctx, &users); err != nil {
        return nil, err
    }
    return &users, nil
}

