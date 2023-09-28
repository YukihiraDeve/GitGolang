package server

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func Start() error {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Erreur de chargement du fichier .env")
	}

	app := fiber.New()
	app.Use(logger.New())

	// Ajoutez vos routes ici

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not set")
	}

	app.Get("/download", func(c *fiber.Ctx) error {
		return c.SendFile("repositories.zip", true)
	})

	log.Printf("Server is running on http://localhost:%s", port)
	return app.Listen(fmt.Sprintf(":%s", port))

}
