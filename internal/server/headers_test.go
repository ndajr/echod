package server

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http/httptest"
	"testing"

	. "github.com/franela/goblin"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func createUserIdServer() *fiber.App {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString(getUserId(c))
	})
	return app
}

func generateUser(n int) (string, string) {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	usr := make([]rune, n)
	pwd := make([]rune, n)
	for i := range usr {
		usr[i] = letters[rand.Intn(len(letters))]
		pwd[i] = letters[rand.Intn(len(letters))]
	}
	return string(usr), string(pwd)
}

func createRequestIdServer() *fiber.App {
	app := fiber.New()
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString(getRequestId(c))
	})
	return app
}

func TestHeaders(t *testing.T) {
	g := Goblin(t)

	g.Describe("Headers", func() {
		g.It("verify user id on Authorization header", func() {
			app := createUserIdServer()
			usr, pwd := generateUser(7)
			req := httptest.NewRequest("GET", "/", nil)
			auth := base64.StdEncoding.EncodeToString([]byte(usr + ":" + pwd))
			req.Header.Add("Authorization", "Basic "+auth)
			resp, _ := app.Test(req)
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalln(err)
			}
			g.Assert(string(body)).Equal(usr)
		})

		g.It("verify request id on X-Request-ID header", func() {
			// send a request without headers
			app := createRequestIdServer()
			req := httptest.NewRequest("GET", "/", nil)
			resp, _ := app.Test(req)
			g.Assert(resp.Header.Get("X-Request-ID")).IsNotZero()

			// send a request with X-Request-ID header
			req = httptest.NewRequest("GET", "/", nil)
			rid := uuid.New().String()
			req.Header.Set("X-Request-ID", rid)
			resp, _ = app.Test(req)
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalln(err)
			}
			g.Assert(string(body)).Equal(rid)
		})
	})

}
