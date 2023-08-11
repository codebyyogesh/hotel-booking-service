package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Hotel struct{
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    Name      string `bson:"Name" json:"Name"`
    Location  string `bson:"location" json:"location"`
    Rooms     []primitive.ObjectID `bson:"rooms" json:"rooms"`
    Rating    int `bson:"rating" json:"rating"`
}

type RoomType int
const (
    _ RoomType = iota
    SingleRoomType
    DoubleRoomType
    DeluxeRoomType
    SuiteRoomType
)

type Room struct{
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
    // can be small, normal and kingsize
    Size      string             `bson:"size" json:"size"`
    Suite     bool               `bson:"suite" json:"suite"`
    Price     float64            `bson:"price" json:"price"`
    HotelID   primitive.ObjectID `bson:"hotelID" json:"hotelID"`
}