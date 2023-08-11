package api

import (

	"github.com/codebyyogesh/hotel-booking-service/db"
	"github.com/gofiber/fiber/v2"
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
    rooms, err := h.store.Rooms.GetRooms(c.Context(), id)
    if err != nil{
        return err
    }
    return c.JSON(rooms)
}

func (h *HotelHandler)HandleGetHotel(c *fiber.Ctx) error{
    id := c.Params("id")
    hotel, err := h.store.Hotels.GetHotelByID(c.Context(), id)
    if err != nil{
        return err
    }
    return c.JSON(hotel)
}
// get hotels list
func (h *HotelHandler)HandleGetHotels(c *fiber.Ctx) error{
    hotels, err := h.store.Hotels.GetHotels(c.Context(), nil)
    if err != nil{
        return err
    }
    return c.JSON(hotels)
}