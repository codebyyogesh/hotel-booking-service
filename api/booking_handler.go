package api

import (
	"net/http"

	"github.com/codebyyogesh/hotel-booking-service/db"
	"github.com/codebyyogesh/hotel-booking-service/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)


type BookingHandler struct{
    store *db.Store 
}

func NewBookingHandler(store *db.Store) *BookingHandler{
    return &BookingHandler{
        store : store,
    }
}

// should be admin authorized
func (h *BookingHandler)HandleGetBookings(c *fiber.Ctx) error{
    bookings, err := h.store.Booking.GetBookings(c.Context(), bson.M{})
    if err != nil{
        return err
    }
    return c.JSON(bookings)
}

// should be user authorized. Booking done can be seen only by the user who booked
// the room and nobody else.
func (h *BookingHandler)HandleGetBooking(c *fiber.Ctx) error{
    id := c.Params("id")
    booking, err := h.store.Booking.GetBookingByID(c.Context(), id)
    if err != nil{
        return err
    }
    user, ok := c.Context().Value("user").(*types.User)
    if !ok{
        return err
    }
    if booking.UserID != user.ID{
        return c.Status(http.StatusUnauthorized).JSON(genericResp{
            Type : "error",
            Msg: "not authorized to view this booking",
        })
    }
    return c.JSON(booking)
}