package api

import (
	"fmt"

	"github.com/codebyyogesh/hotel-booking-service/db"
	"github.com/gofiber/fiber/v2"
)

type HotelHandler struct{
    hotelStore db.HotelStore
    roomStore db.RoomStore
}

func NewHotelHandler(hs db.HotelStore, rs db.RoomStore) *HotelHandler{
    return &HotelHandler{
        hotelStore: hs,
        roomStore : rs,
    }
}

type HotelQueryParam struct{
    Room bool
    Rating int
}
// get hotels list
func (h *HotelHandler)HandleGetHotels(c *fiber.Ctx) error{
    var qparam HotelQueryParam
    err := c.QueryParser(&qparam)
    if err != nil{
        return err
    }
    fmt.Println(qparam)
    hotels, err := h.hotelStore.GetHotels(c.Context(), nil)
    if err != nil{
        return err
    }
    return c.JSON(hotels)
}