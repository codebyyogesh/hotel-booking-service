package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/codebyyogesh/hotel-booking-service/db/fixtures"
	"github.com/gofiber/fiber/v2"
)

func TestAuthenticateWithWrongPassword(t *testing.T) {
    tdb := setup(t)
    defer tdb.teardown(t)
    fixtures.CreateUser(tdb.Store, "some", "foo", false)

    app := fiber.New()
    authHandler := NewAuthHandler(tdb.User)
    app.Post("/auth", authHandler.HandleAuthenticate)

    params := AuthParams{
        Email:    "some@foo.com",
        Password: "mybestsecurepasswordwrongpassword",
    }
    b, _ := json.Marshal(params)
    req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
    req.Header.Add("Content-Type", "application/json")
    resp, err := app.Test(req)
    if err != nil {
        t.Fatal(err)
    }
    if resp.StatusCode != http.StatusBadRequest {
        t.Fatalf("expected status code 400 but got %d", resp.StatusCode)
    }

    var genResp genericResp
    if err := json.NewDecoder(resp.Body).Decode(&genResp); err != nil {
        t.Error(err)
    }
    if genResp.Type != "error" {
        t.Fatalf("expected genResponse type to be error but got %s", genResp.Type)
    }
    if genResp.Msg != "invalid credentials" {
        t.Fatalf("expected genResponse msg to be invalid credentials but got %s", genResp.Msg)
    }
}
func TestAuthenticateSuccess(t *testing.T) {
    tdb := setup(t)
    defer tdb.teardown(t)
    insertedUser := fixtures.CreateUser(tdb.Store, "some", "foo", false)

    app := fiber.New()
    authHandler := NewAuthHandler(tdb.User)
    app.Post("/auth", authHandler.HandleAuthenticate)

    params := AuthParams{
        Email:    "some@foo.com",
        Password: "some_foo",
    }
    b, _ := json.Marshal(params)
    req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
    req.Header.Add("Content-Type", "application/json")
    resp, err := app.Test(req)
    if err != nil {
        t.Fatal(err)
    }
    if resp.StatusCode != http.StatusOK {
        t.Fatalf("expected status code %d, got %d", http.StatusOK, resp.StatusCode)
    }

    var authResp AuthResponse
    if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
        t.Error(err)
    }
    if authResp.Token == "" {
        t.Fatalf("expected the JWT token to be present in the auth response")
    }
    // Because we do not return encrypted password in any of the JSON response, thus we need to
    // set it to null string when comparing with authResp.User
    insertedUser.EncryptedPassword = ""
    if !reflect.DeepEqual(authResp.User, insertedUser) {
        t.Fatalf("expected the user to be %v, got %v", insertedUser, authResp.User)
    }
}
