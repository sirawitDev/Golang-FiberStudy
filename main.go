package main

import (
	"github.com/gofiber/fiber/v2"
)

// Book struct to hold book data
type Book struct {
	ID	int `json:"id"`
	Title string `json:"title"`
	Author	string 	`json:"author"`
}

var books []Book // Slice to store books []Slice [1]Array

func main() {
	app := fiber.New() // app ตัวแทน Fiber application || app = express() ใน Node.js

	//initialize in-memory data
	books = append(books, Book{ ID:1, Title:"1984", Author: "George Orwell"})
	books = append(books, Book{ ID:2, Title:"The Great Gatsby", Author:"F. Scott Fitzgerald"})

	// Apply CORS middleware
	// app.Use(cors.New(cors.Config{
	// 	AllowOrigins: "*", // Adjust this to be more restrictive if needed
	// 	AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
	// 	AllowHeaders: "Origin, Content-Type, Accept",
	// })

	app.Get("/books" , getBooks)
	app.Get("/books/:id" , getBook)
	app.Post("/books" , createBook)
	app.Put("/books/:id" , updateBook)
	app.Delete("/books/:id" , deleteBook)

	app.Post("/upload" , uploadFile)

	app.Listen(":8080")
}

func uploadFile(c *fiber.Ctx) error {
	file, err := c.FormFile("image")

	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	err = c.SaveFile(file, "./uploads/" + file.Filename)

	if err != nil {
		return  c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendString("File uploaded successfully: " + file.Filename)
}


//nodemon --watch . --ext go --exec go run . --signal SIGTERM
