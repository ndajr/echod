package server

import (
	"encoding/base64"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// getRequestId gets the request id from the X-Request-ID header, in case it doesn't exist it creates a new uuid and adds it to the headers.
func getRequestId(c *fiber.Ctx) string {
	rid := string(c.Request().Header.Peek("X-Request-ID"))
	if rid == "" {
		rid = uuid.New().String()
		c.Append("X-Request-ID", rid)
	}
	return rid
}

// getUserId gets the user id from the Authorization, base64 decoding the second part of the header. E.g: "Basic YQ==" -> decode("YQ==") -> "a"
func getUserId(c *fiber.Ctx) string {
	authorization := string(c.Request().Header.Peek("Authorization"))
	if authorization == "" {
		return ""
	}
	s := strings.Split(authorization, " ")
	if len(s) > 0 {
		decoded, err := base64.StdEncoding.DecodeString(s[1])
		if err != nil {
			return ""
		}
		user := strings.Split(string(decoded), ":")
		return string(user[0])
	}
	return ""
}
