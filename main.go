package main

import (
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/proxy"
)

func main() {
	host := os.Getenv("PROXY_HOST")
	if host == "" {
		panic("PROXY_HOST is not set")
	}
	schema := os.Getenv("PROXY_SCHEMA")
	if schema == "" {
		schema = "https"
	}
	app := fiber.New()
	app.Use(proxy.Balancer(proxy.Config{
		Servers: []string{schema + "://" + host},
		ModifyRequest: func(c fiber.Ctx) error {
			c.Request().Header.Set(fiber.HeaderHost, host)
			return nil
		},
	}))
	panic(app.Listen(":3000"))
}
