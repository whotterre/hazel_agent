package handlers

import (
	"hazel_ai/internal/agent"
	"hazel_ai/internal/store"
	a2alogic "hazel_ai/internal/a2a"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	birthdayStore *store.BirthdayStore
}

func NewHandler(birthdayStore *store.BirthdayStore) *Handler {
	return &Handler{
		birthdayStore: birthdayStore,
	}
}

func (h *Handler) Health(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "healthy",
		"agent":  "hazel",
	})

}

func (h *Handler) GetAgentCard(c *fiber.Ctx) error {
	agentCard, err := agent.LoadDefaultAgentCard()
	if err != nil {
		return err
	}
	return c.Status(http.StatusOK).Send(agentCard)
}

func (h *Handler) AddBirthday(c *fiber.Ctx) error {
	type AddBirthdayRequest struct {
		Name string `json:"name"`
		Date string `json:"date"`
	}

	var req AddBirthdayRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	id, err := h.birthdayStore.AddBirthday(req.Name, req.Date)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to add birthday: " + err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{
		"message": "Birthday added successfully",
		"name":    req.Name,
		"date":    req.Date,
		"id":      id,
	})
}

func (h *Handler) ListBirthdays(c *fiber.Ctx) error {
	birthdays := h.birthdayStore.List()
	return c.Status(200).JSON(fiber.Map{
		"birthdays": birthdays,
		"total":     len(birthdays),
	})
}

func (h *Handler) GetTodaysBirthdays(c *fiber.Ctx) error {
	today := time.Now()
	birthdays := h.birthdayStore.List()

	var todaysBirthdays []store.Birthday
	for _, b := range birthdays {
		if b.Month == int(today.Month()) && b.Day == today.Day() {
			todaysBirthdays = append(todaysBirthdays, b)
		}
	}

	return c.Status(200).JSON(fiber.Map{
		"birthdays": todaysBirthdays,
		"count":     len(todaysBirthdays),
	})
}

func (h *Handler) GetUpcomingBirthdays(c *fiber.Ctx) error {
	now := time.Now()
		birthdays := h.birthdayStore.List()

		var upcoming []store.Birthday
		for _, b := range birthdays {
			thisYear := time.Date(now.Year(), time.Month(b.Month), b.Day, 0, 0, 0, 0, now.Location())
			if thisYear.Before(now) {
				thisYear = time.Date(now.Year()+1, time.Month(b.Month), b.Day, 0, 0, 0, 0, now.Location())
			}

			daysUntil := int(thisYear.Sub(now).Hours() / 24)
			if daysUntil <= 30 && daysUntil > 0 {
				upcoming = append(upcoming, b)
			}
		}

		return c.Status(200).JSON(fiber.Map{
			"birthdays": upcoming,
			"count":     len(upcoming),
		})
}

func (h *Handler) SendA2AMessage(c *fiber.Ctx) error {
		type A2AMessage struct {
			From    string      `json:"from"`
			To      string      `json:"to"`
			Content interface{} `json:"content"`
			Type    string      `json:"type"`
		}

		var msg A2AMessage
		if err := c.BodyParser(&msg); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid A2A message"})
		}

		log.Printf("Received A2A message from %s: %+v", msg.From, msg.Content)

		return c.Status(200).JSON(fiber.Map{
			"status":   "processed",
			"response": "Message received and processed",
		})
}

func (h *Handler) UseTelexWebhook (c *fiber.Ctx) error {
	type TelexWebhook struct {
			Event string      `json:"event"`
			Data  interface{} `json:"data"`
		}

		var webhook TelexWebhook
		if err := c.BodyParser(&webhook); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid webhook payload"})
		}

		log.Printf("Received Telex webhook: %s", webhook.Event)

		switch webhook.Event {
		case "daily_check":
			log.Println("Triggering daily birthday check...")
			a2alogic.Remember()
		default:
			log.Printf("Unknown webhook event: %s", webhook.Event)
		}

		return c.Status(200).JSON(fiber.Map{"status": "ok"})
}