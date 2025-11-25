package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New() // app ตัวแทน Fiber application || app = express() ใน Node.js

	app.Get("/hello" , func(c * fiber.Ctx) error {
		return c.SendString("Hello Golang")
	})

	app.Listen(":8080")
}
