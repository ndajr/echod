package server

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var app *fiber.App = New()

func authenticate(req *http.Request) *http.Request {
	auth := base64.StdEncoding.EncodeToString([]byte("remarkablebob" + ":" + "2rFtex"))
	req.Header.Add("Authorization", "Basic "+auth)
	return req
}

func ExampleNew() {
	req := authenticate(httptest.NewRequest("POST", "/api/echo", strings.NewReader(`{"foo":"bar"}`)))
	resp, _ := app.Test(req)
	fmt.Println(resp.StatusCode)
	// Output: 200
}

func ExampleNew_notFound() {
	req := authenticate(httptest.NewRequest("PUT", "/unknown", strings.NewReader(`{}`)))
	resp, _ := app.Test(req)
	fmt.Println(resp.StatusCode)
	// Output: 404
}

func ExampleNew_badRequest() {
	req := authenticate(httptest.NewRequest("POST", "/api/echo", strings.NewReader(`{"foo":"bar", "echoed": true}`)))
	resp, _ := app.Test(req)
	fmt.Println(resp.StatusCode)
	// Output: 400
}

func ExampleNew_unauthorizedRequest() {
	req := httptest.NewRequest("POST", "/api/echo", strings.NewReader(`{"foo":"bar"}`))
	resp, _ := app.Test(req)
	fmt.Println(resp.StatusCode)
	// Output: 401
}
