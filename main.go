package main

import (
	"hazel_ai/internal/agent"
	"hazel_ai/internal/handlers"
	"hazel_ai/internal/store"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {
	birthdayStore := store.NewBirthdayStore("birthdays.json")

	err := agent.CheckForAgentCard()
	if err != nil {
		log.Printf("Warning: Failed to load agent card: %v", err)
		log.Println("Continuing without agent card - some endpoints may not work")
	}

	router := fiber.New()
	handlerList := handlers.NewHandler(birthdayStore)

	router.Get("/health", handlerList.Health)
	router.Get("/.well-known/agent.json", handlerList.GetAgentCard)

	router.Post("/api/birthdays", handlerList.AddBirthday)

	router.Get("/api/birthdays", handlerList.ListBirthdays)

	router.Get("/api/birthdays/today", handlerList.GetTodaysBirthdays)

	router.Get("/api/birthdays/upcoming", handlerList.GetUpcomingBirthdays)

	router.Post("/api/a2a/message", handlerList.SendA2AMessage)

	router.Post("/api/telex/webhook", handlerList.UseTelexWebhook)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Starting server on port %s", port)
	log.Fatal(router.Listen(":" + port))
}
