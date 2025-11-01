package main

import (
	"hazel_ai/internal/agent"
	"hazel_ai/internal/handlers"
	"hazel_ai/internal/store"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists (optional - won't fail if missing)
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("No .env file found or failed to load: %v", err)
		log.Println("Using system environment variables")
	} else {
		log.Println("Successfully loaded .env file")
		apiKey := os.Getenv("GEMINI_API_KEY")
		if apiKey != "" && len(apiKey) > 10 {
			log.Printf("GEMINI_API_KEY loaded: %s...", apiKey[:10])
		} else if apiKey != "" {
			log.Printf("GEMINI_API_KEY loaded: %s", apiKey)
		} else {
			log.Println("GEMINI_API_KEY not found in environment")
		}
	}

	birthdayStore := store.NewBirthdayStore("birthdays.json")

	err = agent.CheckForAgentCard()
	if err != nil {
		log.Printf("Warning: Failed to load agent card: %v", err)
		log.Println("Continuing without agent card - some endpoints may not work")
	}

	router := fiber.New()
	handlerList := handlers.NewHandler(birthdayStore)

	// Telex A2A endpoint - ALL A2A communication goes through POST /
	router.Post("/", handlerList.HandleTelexA2A)

	router.Get("/health", handlerList.Health)
	router.Get("/.well-known/agent.json", handlerList.GetAgentCard)

	router.Post("/api/birthdays", handlerList.AddBirthday)

	router.Get("/api/birthdays", handlerList.ListBirthdays)

	router.Get("/api/birthdays/today", handlerList.GetTodaysBirthdays)

	router.Get("/api/birthdays/upcoming", handlerList.GetUpcomingBirthdays)

	// Birthday wish generation endpoints
	router.Post("/api/wishes/generate", handlerList.GenerateBirthdayWish)
	router.Get("/api/wishes/person/:id", handlerList.GenerateBirthdayWishForPerson)
	router.Get("/api/wishes/simple", handlerList.GenerateSimpleBirthdayWish)

	router.Post("/api/a2a/message", handlerList.SendA2AMessage)

	router.Post("/api/telex/webhook", handlerList.UseTelexWebhook)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Starting Hazel Birthday Bot server on port %s", port)
	log.Fatal(router.Listen(":" + port))
}
