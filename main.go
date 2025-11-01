package main

import (
	"hazel_ai/internal/agent"
	"hazel_ai/internal/handlers"
	"hazel_ai/internal/store"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	birthdayStore := store.NewBirthdayStore("birthdays.json")

	err := agent.CheckForAgentCard()
	if err != nil {
		log.Fatal("Failed to fetch agent card / none found")
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

	log.Println("Starting server on port 3000")
	log.Fatal(router.Listen(":3000"))
}
