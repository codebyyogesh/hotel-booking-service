package api

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
    if apiError, ok := err.(Error); ok{
        return c.Status(apiError.Code).JSON(apiError)
    }
        apiErr := NewError(http.StatusInternalServerError, err.Error())
        return c.Status(apiErr.Code).JSON(apiErr)
}



type Error struct {
    Code int    `json:"code"`
    Msg  string `json:"msg"`
}

// Error implements the error interface
func (e Error)Error() string {
    return e.Msg
}


func NewError(code int, msg string) Error {
    return Error{
        Code: code,
        Msg:  msg,
    }
}
func  ErrorUnAuthorized() Error  {
    return Error{
        Code: http.StatusUnauthorized,
        Msg:  "unauthorized request",
    }
}

func ErrorBadRequest() Error {
    return Error{
        Code: http.StatusBadRequest,
        Msg:  "invalid JSON request",
    }
}

func ErrorResourceNotFound(res string) Error {
    return Error{
        Code: http.StatusNotFound,
        Msg:  res + " resource not found",
    }
}
func  ErrorInvalidID() Error  {
    return Error{
        Code: http.StatusBadRequest,
        Msg:  "invalid id",
    }
}