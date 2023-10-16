package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/codebyyogesh/hotel-booking-service/api/middleware"
	"github.com/codebyyogesh/hotel-booking-service/db/fixtures"
	"github.com/codebyyogesh/hotel-booking-service/types"
	"github.com/gofiber/fiber/v2"
)
func TestGetUserBooking(t *testing.T) {
    tdb := setup(t)
    defer tdb.teardown(t)
    var (
        user = fixtures.CreateUser(tdb.Store, "some", "foo", false)
        nonAuthUser = fixtures.CreateUser(tdb.Store, "damm", "foo", false)
        hotel = fixtures.CreateHotel(tdb.Store, "hotel", "location", 5, nil)
        room = fixtures.CreateRoom(tdb.Store, "small", false, 300, hotel.ID)
        from = time.Now()
        till = time.Now().AddDate(0, 0, 5)
        booking = fixtures.CreateBooking(tdb.Store, user.ID, room.ID, 2, from, till)
        app  = fiber.New()
        bookingHandler = NewBookingHandler(tdb.Store)
        route = app.Group("/", middleware.JWTAuthentication(tdb.User))
    )
    _ = booking
    route.Get("/:id", bookingHandler.HandleGetBooking)
    req := httptest.NewRequest("GET", fmt.Sprintf("/%s",booking.ID.Hex()), nil)
    req.Header.Add("X-Api-Token", CreateTokenFromUser(user)) // regular user
    resp, err := app.Test(req)
    if err != nil{
        t.Fatal(err)
    }
    if resp.StatusCode != http.StatusOK {
        t.Fatalf("non 200 code got %d", resp.StatusCode)
    }
    var oneBooking types.Booking // a single booking
    if err := json.NewDecoder(resp.Body).Decode(&oneBooking); err != nil{
        t.Fatal(err)
    }
    if oneBooking.ID != booking.ID {
        t.Fatalf("expected %s got %s", booking.ID, oneBooking.ID)
    }
    if oneBooking.UserID != booking.UserID {
        t.Fatalf("expected %s got %s", booking.UserID, oneBooking.UserID)
    }

    // test a non authorized user, because you created a booking from a user and testing it for another user

    req = httptest.NewRequest("GET", fmt.Sprintf("/%s",booking.ID.Hex()), nil)
    req.Header.Add("X-Api-Token", CreateTokenFromUser(nonAuthUser)) // regular user
    resp, err = app.Test(req)
    if err != nil{
        t.Fatal(err)
    }
    if resp.StatusCode == http.StatusOK {
        t.Fatalf("non 200 code got %d", resp.StatusCode)
    }

}
 

func TestAdminGetBookings(t *testing.T) {
    tdb := setup(t)
    defer tdb.teardown(t)
    var (
        user = fixtures.CreateUser(tdb.Store, "some", "foo", false)
        admin = fixtures.CreateUser(tdb.Store, "admin", "foo", true)
        hotel = fixtures.CreateHotel(tdb.Store, "hotel", "location", 5, nil)
        room = fixtures.CreateRoom(tdb.Store, "small", false, 300, hotel.ID)
        from = time.Now()
        till = time.Now().AddDate(0, 0, 5)
        booking = fixtures.CreateBooking(tdb.Store, user.ID, room.ID, 2, from, till)
        app  = fiber.New()
        bookingHandler = NewBookingHandler(tdb.Store)
        // Group uses variadic params as handlers, here there are two handlers passed
        // first JWTAuthentication gets first called followed by AdminAuth
        adminGrp = app.Group("/", middleware.JWTAuthentication(tdb.User), middleware.AdminAuth )
    )
    _ = booking
    adminGrp.Get("/", bookingHandler.HandleGetBookings)
    req := httptest.NewRequest("GET", "/", nil)
    req.Header.Add("X-Api-Token", CreateTokenFromUser(admin)) // admin user
    resp, err := app.Test(req)
    if err != nil{
        t.Fatal(err)
    }
     if resp.StatusCode != http.StatusOK {
        t.Fatalf("non 200 code got %d", resp.StatusCode)
    }
    var bookings []types.Booking // must be an array because we have multiple bookings
    if err := json.NewDecoder(resp.Body).Decode(&bookings); err != nil{
        t.Fatal(err)
    }
    if(len(bookings) != 1){
        t.Fatalf("non 1 booking got %d", len(bookings))
    }
    have := bookings[0]
    if have.ID != booking.ID {
        t.Fatalf("expected %s got %s", booking.ID, have.ID)
    }
    if have.UserID != booking.UserID {
        t.Fatalf("expected %s got %s", booking.UserID, have.UserID)
    }

    // test to ensure non admins cannot access bookings
    req = httptest.NewRequest("GET", "/", nil)
    req.Header.Add("X-Api-Token", CreateTokenFromUser(user)) // regular user
    resp, err = app.Test(req)
    if err != nil {
        t.Fatal(err)
    }
    if resp.StatusCode == http.StatusOK {
        t.Fatalf("expected a non 200 status code got %d", resp.StatusCode)
    }
}