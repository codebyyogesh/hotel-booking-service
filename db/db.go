package db

// General db functions
const (
    DBNAME     = "hotel-booking"
    TestDBNAME = "hotel-booking-test"
    DBURI      =  "mongodb://localhost:27017"
)

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