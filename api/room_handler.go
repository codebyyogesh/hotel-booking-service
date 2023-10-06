package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/codebyyogesh/hotel-booking-service/db"
	"github.com/codebyyogesh/hotel-booking-service/types"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BookRoomParams struct {
    FromDate   time.Time `json:"fromDate"`
    TillDate   time.Time `json:"tillDate"`
    NumPersons int       `json:"numPersons"`
}

type RoomHandler struct {
    store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
    return &RoomHandler{
        store: store,
    }
}

func (p BookRoomParams) Validate() error  {
    now := time.Now()
    if now.After(p.FromDate) || now.After(p.TillDate) {
        return fmt.Errorf("INVALID DATES")
    }
    return nil
}

// check if the room is already booked in this period. We will use the OR operator
// to check if either of the dates are valid
func (r *RoomHandler)isRoomAvailableForBooking(ctx context.Context, roomID primitive.ObjectID, params BookRoomParams) (bool, error) {
    filter :=   bson.M{
        "roomID": roomID,
        "fromDate": bson.M{
            "$gte": params.FromDate,
        },
        "tillDate": bson.M{
            "$lte": params.TillDate,
        },
    }

    bookings, err := r.store.Booking.GetBookings(ctx, filter)
    if err != nil {
        return false, err
    }
    // There are already some bookings for this room, so return
     if bookings != nil{ 
        return false, nil
    }
    return true, nil
}

func (r *RoomHandler) HandleBookRoom(c *fiber.Ctx) error {
    var params BookRoomParams
    if err := c.BodyParser(&params); err != nil {
        return err
    }
    if err := params.Validate(); err != nil {
        return err
    }
    roomID, err := primitive.ObjectIDFromHex(c.Params("id"))
    if err != nil {
        return err
    }
    user, ok := c.Context().Value("user").(*types.User)
    if !ok {
        return c.Status(http.StatusInternalServerError).JSON(genericResp{
            Type: "error",
            Msg:  "internal server error",
        })
    }

    ok, err = r.isRoomAvailableForBooking(c.Context(), roomID, params)

    if err != nil {
        return err
    }
    if !ok {
        return c.Status(http.StatusBadRequest).JSON(genericResp{
            Type: "error",
            Msg:  fmt.Sprintf("room %s already booked", c.Params("id")),
        })
    }

    // new booking
    booking := types.Booking{
        UserID:     user.ID,
        RoomID:     roomID,
        FromDate:   params.FromDate,
        TillDate:   params.TillDate,
        NumPersons: params.NumPersons,
    }
    inserted, err := r.store.Booking.InsertBooking(c.Context(), &booking)
    if err != nil {
        return err
    }
    return c.JSON(inserted)
}

func (r *RoomHandler) HandleGetRooms(c *fiber.Ctx) error {

    rooms, err := r.store.Rooms.GetRooms(c.Context(), bson.M{})
    if err != nil {
        return err
    }
    return c.JSON(rooms)
}

