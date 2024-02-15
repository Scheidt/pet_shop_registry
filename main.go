package main

// go run main.go pet.go petController.go

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	
	"gorm.io/gorm"
)

type Pet struct {
	ID    			int			`json: "id"`
	name    		*string		`json: "name"`
	age  			int			`json: "age"`
	owner   		string		`json: "owner"`
	size    		string		`json: "size"`
	weight			float32		`json: "weight"`
	paidThisMonth 	bool		`json: "paid"`
}

type Repository struct {
	DB *gorm.db
}

func (r *Repository) CreatePet(context *fiber.Ctx) error {
	pet := Pet{}
	err := context.BodyParser(&pet)
	if err != nil {
        context.Status(http.StatusUnprocessableEntity).JSON(&fiber.Map{"message": "request failed"})
        return err
	}

	err = r.DB.Create(&pet).Error
	if err != nil {
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message":"couldn't create pet"})
		return err
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{"message":"pet has been added"})
	return nil
}

func (r *Repository) GetPet(context *fiber.Ctx) error{
	petModels :=&models.Pet{}

	err := r.DB.Find(petModels).Error
	if err != nil{
		context.Status(http.StatusBadRequest).JSON(&fiber.Map{"message":"could not get pets"})
	}
	context.Status(http.StatusOK).JSON(&fiber.Map{	"message": "pets fetched successfully",
													"data": petModels})
	return nil
}

func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/create_pets", r.CreatePet)
	api.Delete("delete_pet/:id", r.DeletePet)
	api.Get("/get_pets/:id", r.GetPetByID)
	api.Get("/pet", r.GetPet)
}

//type pet struct {
//	title    string
//	pubYear  int
//	author   string
//	genre    [2]string
//	borrowed bool
//}

func getStringInput(texto string) string {
	fmt.Print(texto)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()
	input = strings.ReplaceAll(input, "'", "''")
	return input
}

func getIntInput(message string) int {
	var userInput int
	for {
		fmt.Print(message)
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input := strings.TrimSpace(scanner.Text())

		value, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Only integer numbers are allowed, please, try again: ")
			continue
		}
		userInput = value
		break
	}
	return userInput
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	config := &database.Config{
		host: os.Getenv("DB_HOST"),
		port: os.Getenv("DB_PORT"),
		password: os.Getenv("DB_PASS"),
		user: os.Getenv("DB_USER"),
		sslmode: os.Getenv("DB_SSLMODE"),
		dbName: os.Getenv("DB_NAME"),
	}

	db, err := database.establishConnection(config)

	if err != nil {
		log.Fatal("could not load the database")
	}

	r := Repository{
		DB: db,
	}
	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080")
}
