package db

import (
	"os"

	env "github.com/codebyyogesh/hotel-booking-service/config"
)

// General db functions
var DBNAME string
var DBURI string

type Store struct{
    User  UserStore
    Hotels HotelStore
    Rooms RoomStore
    Booking BookingStore
}

type Pagination struct{
    Limit int
    Page int
}

func init() {
    env.LoadEnv()
    DBNAME = os.Getenv("MONGO_DB_NAME")
    DBURI  = os.Getenv("MONGO_DB_URL")
}