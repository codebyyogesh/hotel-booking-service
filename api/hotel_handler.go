package api

import (
	"github.com/codebyyogesh/hotel-booking-service/db"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HotelHandler struct{
    store *db.Store 
}

func NewHotelHandler(store *db.Store) *HotelHandler{
    return &HotelHandler{
        store : store,
    }
}

func (h *HotelHandler)HandleGetRooms(c *fiber.Ctx) error{
    id := c.Params("id")
    oid, err := primitive.ObjectIDFromHex(id)
    if err != nil {
        return ErrorInvalidID()
    }
    filter := bson.M{"hotelID": oid}
    rooms, err := h.store.Rooms.GetRooms(c.Context(), filter)
    if err != nil{
        return ErrorResourceNotFound("rooms")
    }
    return c.JSON(rooms)
}

func (h *HotelHandler)HandleGetHotel(c *fiber.Ctx) error{
    id := c.Params("id")
    hotel, err := h.store.Hotels.GetHotelByID(c.Context(), id)
    if err != nil{
        return ErrorResourceNotFound("hotel")
    }
    return c.JSON(hotel)
}
// get hotels list
func (h *HotelHandler)HandleGetHotels(c *fiber.Ctx) error{
    hotels, err := h.store.Hotels.GetHotels(c.Context(), nil)
    if err != nil{
        return ErrorResourceNotFound("hotels")
    }
    return c.JSON(hotels)
}