package main

import (
	"fmt"
	"go-postgre/models"
	"go-postgre/storage"
	"log"
	"os"

	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

type Book struct {
	Author    string `json:"author"`
	Title     string `json:"title"`
	Publisher string `json:"publisher"`
}

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) CreateBook(c *fiber.Ctx) error {
	book := Book{}
	err := c.BodyParser(&book)
	if err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "request failed"})
		return err
	}
	err = r.DB.Create(&book).Error
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not create book"})
		return err
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{"message": "book has been added"})
	return nil
}

func (r *Repository) DeleteBook(c *fiber.Ctx) error {
	BookModel := models.Book{}
	id := c.Params("id")
	if id == "" {
		c.Status(http.StatusNotFound).JSON(&fiber.Map{
			"message": "not found book to delete",
		})
		return nil
	}
	err := r.DB.Delete(BookModel, id).Error
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not delete book",
		})
		return err
	}
	c.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "Delete book successfull",
	})
	return nil
}
func (r *Repository) GetBookById(c *fiber.Ctx) error {
	BookModel := &models.Book{}
	id := c.Params("id")
	if id == "" {
		c.Status(http.StatusNotFound).JSON(&fiber.Map{
			"message": "not found book",
		})
		return nil
	}

	fmt.Println("the id is %s", id)
	err := r.DB.Where("id=?", id).First(BookModel).Error
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"message": "could not get the book",
		})
		return err
	}
	c.Status(http.StatusOK).JSON(&fiber.Map{
		"message": "get book success",
		"data":    &BookModel,
	})
	return nil
}
func (r *Repository) GetBooks(c *fiber.Ctx) error {
	BookModel := &[]models.Book{}
	err := r.DB.Find(BookModel).Error
	if err != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not read the books"})
		return err
	}
	c.Status(http.StatusOK).JSON(
		&fiber.Map{
			"message": "fetch book successfull",
			"data":    BookModel,
		})
	return nil
}

func (r *Repository) SetUpRouter(app *fiber.App) {
	api := app.Group("/api")
	api.Get("/create_book", r.CreateBook)
	api.Get("/delete_book/:id", r.DeleteBook)
	api.Get("/get_book/:id", r.GetBookById)
	api.Get("/books/:id", r.GetBooks)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_Port"),
		Password: os.Getenv("DB_Password"),
		User:     os.Getenv("DB_User"),
		SSLMode:  os.Getenv("DB_SSLMode"),
		DBName:   os.Getenv("DB_Name"),
	}
	db, err := storage.NewConnection(config)
	if err != nil {
		log.Fatal("could not load database")
	}

	err = models.MigrateBook(db)
	if err != nil {
		log.Fatal("could not migrate data")
	}

	r := Repository{
		DB: db,
	}
	app := fiber.New()
	r.SetUpRouter(app)
	app.Listen(":3000")

}
