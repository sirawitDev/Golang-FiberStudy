package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/joho/godotenv"
)

// Book struct to hold book data
type Book struct {
	ID	int `json:"id"`
	Title string `json:"title"`
	Author	string 	`json:"author"`
}

var books []Book // Slice to store books []Slice [1]Array

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Load .env error")
	}

	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	}) // app ตัวแทน Fiber application || app = express() ใน Node.js

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
	app.Get("/test-html" , testHtml)

	app.Get("/config" , getEnv)

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

func testHtml(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Title": "Hello Golang",
		"Name": "FewDev",
	})
}

func getEnv(c *fiber.Ctx) error {
	// if value, exists := os.LookupEnv("SECRET") ; exists {
	// 	return c.JSON(fiber.Map{
	// 		"SECRET": value,
	// 	})
	// }

	// return c.JSON(fiber.Map{
	// 	"SECRET": "defaultsecret",
	// })

	secret := os.Getenv("SECRET")

	if secret == "" {
		secret = "defaultsecret"
	}

	return c.JSON(fiber.Map{
		"SECRET": secret,
	})
}


//nodemon --watch . --ext go --exec go run . --signal SIGTERM
