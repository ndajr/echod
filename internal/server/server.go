// Echod is an HTTP server which implements POST /api/echo, PUT /api/echo, GET /health and GET /metrics endpoints.
package server

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	totalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Number of total requests by status code, method and path",
		},
		[]string{"status", "method", "path"},
	)
	requestDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name: "http_request_duration_milliseconds",
			Help: "Duration of all HTTP requests",
		},
	)
)

func handler(c *fiber.Ctx) error {
	c.Accepts("application/json")
	res, err := Echo(c.Body())
	if err != nil {
		if err.Error() == "invalid json" || err.Error() == "echoed already set" {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		log.Error().Err(err).Dict("context", zerolog.Dict().
			Str("userId", getUserId(c)).
			Str("requestId", getRequestId(c)).
			Str("status", strconv.Itoa(fiber.StatusInternalServerError)),
		).Msg("internal error")
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.JSON(res)
}

// statsMiddleware intercepts any incoming request capturing its details such as method, path and requestId, gives the control back to the called endpoint and intercepts the end of the request logging and computing the metrics with its response.
func statsMiddleware(c *fiber.Ctx) error {
	// Before
	start := time.Now()
	method := string(c.Context().Method())
	path := string(c.Context().Path())

	if err := c.Next(); err != nil {
		return err
	}

	// After
	status := strconv.Itoa(c.Response().StatusCode())
	log.Info().Dict("context", zerolog.Dict().
		Str("userId", getUserId(c)).
		Str("requestId", getRequestId(c)).
		Str("status", status),
	).Msg(fmt.Sprintf("%s %s %s", method, path, c.Protocol()))
	totalRequests.WithLabelValues(status, method, path).Inc()
	elapsed := float64(time.Since(start).Nanoseconds()) / 1000000
	c.Append("X-Response-Time", fmt.Sprintf("%.2fms", elapsed))
	requestDuration.Observe(elapsed)
	return nil
}

// New creates an echod server
func New() *fiber.App {
	app := fiber.New()
	r := prometheus.NewRegistry()
	r.MustRegister(totalRequests)
	r.MustRegister(requestDuration)
	promHandler := promhttp.HandlerFor(r, promhttp.HandlerOpts{})
	app.Get("/health", func(c *fiber.Ctx) error { return c.SendStatus(fiber.StatusOK) })
	app.Get("/metrics", adaptor.HTTPHandler(promHandler))
	app.Use(statsMiddleware)
	app.Use(basicauth.New(basicauth.Config{
		Users: map[string]string{
			"remarkablebob": "2rFtex",
			"dextroussheep": "vSKjGK",
			"luckyslug":     "VRsBMW",
			"pettyrabbit":   "2BpgwH",
		},
	}))
	app.Post("/api/echo", handler)
	app.Put("/api/echo", handler)
	return app
}
