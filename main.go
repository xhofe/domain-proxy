package main

import (
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/proxy"
)

func main() {
	dst := os.Getenv("PROXY_DESTINATION")
	if dst == "" {
		panic("PROXY_DESTINATION is not set")
	}
	app := fiber.New()
	app.Use(proxy.Balancer(proxy.Config{
		Servers: []string{dst},
	}))
	panic(app.Listen(":3000"))
}
