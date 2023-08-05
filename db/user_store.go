package db

import (
	"context"
	"fmt"

	"github.com/codebyyogesh/hotel-booking-service/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const userCollection = "users"

type Dropper interface {
    Drop(context.Context) error
}

type UserStore interface{
    Dropper
    GetUserByID(context.Context, string) (*types.User, error)
    GetUsers(context.Context) (*[]types.User, error) // *[]types.User will be a pointer to a slice
    // PS: []*types.User will be slice of pointers to object of type User
    InsertUser(context.Context, *types.User)(*types.User, error)
    DeleteUser(context.Context, string)(error)
    UpdateUser(context.Context, string, *types.UpdateUserParams)(error)
}

type MongoUserStore struct{
    client *mongo.Client
    collection *mongo.Collection
}


func NewMongoUserStore(client *mongo.Client, dbname string) *MongoUserStore{
    return &MongoUserStore{
        client: client,
        collection: client.Database(dbname).Collection(userCollection),
    }
}

func (s *MongoUserStore) Drop(ctx context.Context) error {
    fmt.Println("--- dropping user collection")
    return s.collection.Drop(ctx)
}

// PS: ToBsonD() converts to BsonD format, also it does some basic validation
// before update. You can as well use Bson.M{} format for update.
func (s *MongoUserStore)UpdateUser(ctx context.Context, id string, params *types.UpdateUserParams ) error {
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil{
        return err
    }
    // PS: You can use filter as both, bson.D or bson.M. Enable only one of them
    //filter := bson.D{{Key: "_id", Value: oid}}
    filter := bson.M{"_id": oid}
    update := bson.D{{Key: "$set", Value: params.ToBsonD()}}
    _, err = s.collection.UpdateOne(context.TODO(), filter, update)
    if err != nil {
        return err
    }
    return nil
}

    func (s *MongoUserStore)DeleteUser(ctx context.Context, id string ) error {
    // Mongodb does not accept direct ids, we need some kind of conversion to object id
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil{
        return err
    }
    filter := bson.D{{Key: "_id", Value: oid}}
    _, err = s.collection.DeleteOne(ctx, filter)
    if err != nil{
        return err
    }
    return nil
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

