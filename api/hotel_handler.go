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

type ResourceResp struct{
    NumberOfItems int `json:"numberofitems"`
    Data any `json:"data"`
    Page int `json:"page"`
}
type HotelQueryParams struct{
    db.Pagination
    Rating int
    //Add more filter query params here (say roomsize or number of rooms etc)
}
// get hotels list
// ToDo : Remove bson.M filter from here and use a generic struct suitable for all databases
// currently only mongodb is supported
func (h *HotelHandler)HandleGetHotels(c *fiber.Ctx) error{
    //var pagination db.Pagination
    var params HotelQueryParams


    if err := c.QueryParser(&params); err != nil{
        return ErrorBadRequest()
    }
    if params.Rating <= 0 || params.Rating > 5{
        return ErrorBadRequest()
    }
    filter := bson.M{"rating": bson.M{"$eq": params.Rating}}

    hotels, err := h.store.Hotels.GetHotels(c.Context(), filter, &params.Pagination)
    if err != nil{
        return ErrorResourceNotFound("hotels")
    }
    resp := ResourceResp{
        NumberOfItems: len(*hotels),
        Data: hotels,
        Page: params.Page,
    }
    return c.JSON(resp)
}