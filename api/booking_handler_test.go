package api

import (
	"fmt"
	"testing"
	"time"

	"github.com/codebyyogesh/hotel-booking-service/db/fixtures"
)


func TestAdminGetBookings(t *testing.T) {
    tdb := setup(t)
    defer tdb.teardown(t)

    user := fixtures.CreateUser(tdb.Store, "some", "foo", false)
    hotel := fixtures.CreateHotel(tdb.Store, "hotel", "location", 5, nil)
    room := fixtures.CreateRoom(tdb.Store, "small", false, 300, hotel.ID)
    from := time.Now()
    till := time.Now().AddDate(0, 0, 5)
    booking := fixtures.CreateBooking(tdb.Store, user.ID, room.ID, 2, from, till)
    fmt.Println(booking)
}