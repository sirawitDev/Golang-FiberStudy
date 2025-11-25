package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/gofiber/swagger"
	"github.com/gofiber/template/html/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	_ "github.com/sirawitDev/golang-fiberserver/docs"
)

// Book struct to hold book data
type Book struct {
	ID	int `json:"id"`
	Title string `json:"title"`
	Author	string 	`json:"author"`
}

var books []Book // Slice to store books []Slice [1]Array

// @title Book API
// @description This is a sample server for a book API.
// @version 1.0
// @host localhost:8080
// @BasePath /
// @schemes http
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Load .env error")
	}

	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views: engine,
	}) // app ตัวแทน Fiber application || app = express() ใน Node.js

	app.Get("/swagger/*", swagger.HandlerDefault)

	//initialize in-memory data
	books = append(books, Book{ ID:1, Title:"1984", Author: "George Orwell"})
	books = append(books, Book{ ID:2, Title:"The Great Gatsby", Author:"F. Scott Fitzgerald"})

	// Apply CORS middleware
	// app.Use(cors.New(cors.Config{
	// 	AllowOrigins: "*", // Adjust this to be more restrictive if needed
	// 	AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
	// 	AllowHeaders: "Origin, Content-Type, Accept",
	// })\

	app.Post("/login" , login)

	app.Use(jwtware.New(jwtware.Config{
        SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}))

	app.Use(checkMiddleware)

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

	secret := os.Getenv("JWT_SECRET")

	if secret == "" {
		secret = "defaultsecret"
	}

	return c.JSON(fiber.Map{
		"SECRET": secret,
	})
}

func checkMiddleware(c *fiber.Ctx) error {
	// start := time.Now()

	// fmt.Printf(
	// 	"URL = %s, Method = %s, Time = %s\n",
	// 	c.OriginalURL() , c.Method() , start,
	// )

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)

	if claims["role"] != "admin" {
		return fiber.ErrUnauthorized
	}

	return c.Next()
}

type User struct {
	Email	string 	`json:"email"`
	Password	string `json:"password"`
}

var member = User{
	Email: "admin@admin.com",
	Password: "admin",
}

func login(c *fiber.Ctx) error {
	user := new(User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if user.Email != member.Email || user.Password != member.Password {
		return  fiber.ErrUnauthorized
	}

	// Create Token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = user.Email
	claims["role"] = "admin"
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token send it as response
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{
		"message : ": "Login Success",
		"token" : t,
	})
}


//nodemon --watch . --ext go --exec go run . --signal SIGTERM
