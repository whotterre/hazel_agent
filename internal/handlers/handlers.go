package handlers

import (
	"fmt"
	a2alogic "hazel_ai/internal/a2a"
	"hazel_ai/internal/agent"
	"hazel_ai/internal/clients"
	"hazel_ai/internal/store"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	birthdayStore *store.BirthdayStore
	geminiClient  *clients.GeminiClient
}

func NewHandler(birthdayStore *store.BirthdayStore) *Handler {
	geminiClient, err := clients.NewGeminiClient()
	if err != nil {
		log.Printf("Warning: Failed to initialize Gemini client: %v", err)
		geminiClient = nil
	}

	return &Handler{
		birthdayStore: birthdayStore,
		geminiClient:  geminiClient,
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
	// Handle both simple A2A messages and Telex JSONRPC format
	var telexRequest map[string]interface{}
	if err := c.BodyParser(&telexRequest); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid A2A message"})
	}

	log.Printf("Received A2A request: %+v", telexRequest)

	// Extract text content from Telex JSONRPC format
	var textContent string

	// Check if this is a Telex JSONRPC request
	if params, ok := telexRequest["params"].(map[string]interface{}); ok {
		if message, ok := params["message"].(map[string]interface{}); ok {
			if parts, ok := message["parts"].([]interface{}); ok && len(parts) > 0 {
				if part, ok := parts[0].(map[string]interface{}); ok {
					if text, ok := part["text"].(string); ok {
						textContent = text
					}
				}
			}
		}
	}

	// If no text content found, try simple format
	if textContent == "" {
		if content, ok := telexRequest["content"].(string); ok {
			textContent = content
		}
	}

	log.Printf("Extracted text content: %s", textContent)

	// Process the text content
	return h.processTextContent(c, textContent, telexRequest)
}

// processTextContent analyzes text and determines what action to take
func (h *Handler) processTextContent(c *fiber.Ctx, text string, originalRequest map[string]interface{}) error {
	text = strings.ToLower(strings.TrimSpace(text))

	// Check for different birthday-related commands
	if strings.Contains(text, "birthday wish") || strings.Contains(text, "wish") {
		return h.handleWishRequest(c, text, originalRequest)
	} else if strings.Contains(text, "remember") || strings.Contains(text, "store") {
		return h.handleRememberRequest(c, text, originalRequest)
	} else if strings.Contains(text, "list") || strings.Contains(text, "show") {
		return h.handleListRequest(c, originalRequest)
	} else {
		// Generic response
		response := "Hello! I'm Hazel, your birthday bot. I can help you with:\n• Generate birthday wishes\n• Remember birthdays\n• List stored birthdays\n\nTry asking me to 'generate a birthday wish' or 'remember my birthday'!"
		return h.sendTelexResponse(c, response, originalRequest)
	}
}

// handleWishRequest processes birthday wish requests
func (h *Handler) handleWishRequest(c *fiber.Ctx, text string, originalRequest map[string]interface{}) error {
	// Try to extract a name from the text
	name := ""
	words := strings.Fields(text)

	// Look for patterns like "wish for John" or "birthday wish for Alice"
	for i, word := range words {
		if (word == "for" || word == "to") && i+1 < len(words) {
			name = strings.Title(words[i+1])
			break
		}
	}

	// If no specific name, generate a generic wish
	if name == "" {
		name = "you"
	}

	// Generate the birthday wish
	var wish string
	var source string

	if h.geminiClient == nil {
		wish = fmt.Sprintf("🎉 Happy Birthday%s! 🎂 Wishing you all the joy, happiness, and wonderful surprises on your special day! May this year bring you endless blessings and amazing adventures! 🌟",
			func() string {
				if name == "you" {
					return ""
				} else {
					return ", " + name
				}
			}())
		source = "fallback"
	} else {
		var err error
		if name == "you" {
			wish, err = h.geminiClient.GenerateGenericBirthdayWish("friend")
		} else {
			wish, err = h.geminiClient.GenerateGenericBirthdayWish(name)
		}

		if err != nil {
			wish = fmt.Sprintf("🎉 Happy Birthday%s! 🎂 Wishing you all the joy, happiness, and wonderful surprises on your special day! 🌟",
				func() string {
					if name == "you" {
						return ""
					} else {
						return ", " + name
					}
				}())
			source = "fallback"
		} else {
			source = "gemini"
		}
	}

	log.Printf("Generated %s birthday wish for %s", source, name)
	return h.sendTelexResponse(c, wish, originalRequest)
}

// handleRememberRequest processes remember birthday requests
func (h *Handler) handleRememberRequest(c *fiber.Ctx, text string, originalRequest map[string]interface{}) error {
	response := "I'd love to remember your birthday! Please tell me the date in YYYY-MM-DD format (like 2003-09-09) and I'll store it for you."
	return h.sendTelexResponse(c, response, originalRequest)
}

// handleListRequest processes list birthdays requests
func (h *Handler) handleListRequest(c *fiber.Ctx, originalRequest map[string]interface{}) error {
	birthdays := h.birthdayStore.List()

	if len(birthdays) == 0 {
		response := "📝 No birthdays stored yet! Ask me to 'remember your birthday' to get started."
		return h.sendTelexResponse(c, response, originalRequest)
	}

	response := fmt.Sprintf("🎂 Stored Birthdays (%d total):\n\n", len(birthdays))
	for _, b := range birthdays {
		response += fmt.Sprintf("• %s - %s %d\n", b.Name, time.Month(b.Month), b.Day)
	}

	return h.sendTelexResponse(c, response, originalRequest)
}

// sendTelexResponse sends a properly formatted response back to Telex
func (h *Handler) sendTelexResponse(c *fiber.Ctx, message string, originalRequest map[string]interface{}) error {
	// Check if this is a Telex JSONRPC request and respond accordingly
	if jsonrpc, ok := originalRequest["jsonrpc"]; ok {
		if id, ok := originalRequest["id"]; ok {
			return c.Status(200).JSON(fiber.Map{
				"jsonrpc": jsonrpc,
				"id":      id,
				"result": fiber.Map{
					"message": fiber.Map{
						"kind": "message",
						"role": "assistant",
						"parts": []fiber.Map{
							{
								"kind": "text",
								"text": message,
							},
						},
					},
				},
			})
		}
	}

	// Fallback to simple response
	return c.Status(200).JSON(fiber.Map{
		"status":   "success",
		"response": message,
	})
}

// Old A2A handlers replaced with new Telex-compatible ones above

func (h *Handler) UseTelexWebhook(c *fiber.Ctx) error {
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

// GenerateBirthdayWish generates a personalized birthday wish using Gemini AI
func (h *Handler) GenerateBirthdayWish(c *fiber.Ctx) error {
	type WishRequest struct {
		Name string `json:"name"`
		Age  int    `json:"age,omitempty"`
	}

	var req WishRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if req.Name == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Name is required"})
	}

	if h.geminiClient == nil {
		// Fallback to simple message if Gemini is not available
		fallbackWish := "🎉 Happy Birthday, " + req.Name + "! 🎂 Wishing you all the joy and happiness on your special day! 🌟"
		return c.Status(200).JSON(fiber.Map{
			"name":   req.Name,
			"wish":   fallbackWish,
			"source": "fallback"})
	}

	var wish string
	var err error

	if req.Age > 0 {
		wish, err = h.geminiClient.GenerateBirthdayWish(req.Name, req.Age)
	} else {
		wish, err = h.geminiClient.GenerateGenericBirthdayWish(req.Name)
	}

	if err != nil {
		log.Printf("Error generating birthday wish: %v", err)
		fallbackWish := "🎉 Happy Birthday, " + req.Name + "! 🎂 Wishing you all the joy and happiness on your special day! 🌟"
		return c.Status(200).JSON(fiber.Map{
			"name":   req.Name,
			"wish":   fallbackWish,
			"source": "fallback"})
	}

	return c.Status(200).JSON(fiber.Map{
		"name":   req.Name,
		"wish":   wish,
		"age":    req.Age,
		"source": "gemini"})
}

// GenerateBirthdayWishForPerson generates a birthday wish for a specific person by ID
func (h *Handler) GenerateBirthdayWishForPerson(c *fiber.Ctx) error {
	personID := c.Params("id")
	if personID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Person ID is required"})
	}

	// Find the person in the birthday store
	birthdays := h.birthdayStore.List()
	var targetPerson *store.Birthday

	for _, b := range birthdays {
		if b.ID == personID {
			targetPerson = &b
			break
		}
	}

	if targetPerson == nil {
		return c.Status(404).JSON(fiber.Map{"error": "Person not found"})
	}

	// For now, we'll generate generic wishes since we only store month/day
	age := 0

	if h.geminiClient == nil {
		// Fallback message
		fallbackWish := "🎉 Happy Birthday, " + targetPerson.Name + "! 🎂 Wishing you all the joy and happiness on your special day! 🌟"
		return c.Status(200).JSON(fiber.Map{
			"id":     targetPerson.ID,
			"name":   targetPerson.Name,
			"wish":   fallbackWish,
			"source": "fallback"})
	}

	var wish string
	var err error

	if age > 0 {
		wish, err = h.geminiClient.GenerateBirthdayWish(targetPerson.Name, age)
	} else {
		wish, err = h.geminiClient.GenerateGenericBirthdayWish(targetPerson.Name)
	}

	if err != nil {
		log.Printf("Error generating birthday wish for %s: %v", targetPerson.Name, err)
		fallbackWish := "🎉 Happy Birthday, " + targetPerson.Name + "! 🎂 Wishing you all the joy and happiness on your special day! 🌟"
		return c.Status(200).JSON(fiber.Map{
			"id":     targetPerson.ID,
			"name":   targetPerson.Name,
			"wish":   fallbackWish,
			"source": "fallback"})
	}

	return c.Status(200).JSON(fiber.Map{
		"id":     targetPerson.ID,
		"name":   targetPerson.Name,
		"wish":   wish,
		"age":    age,
		"source": "gemini",
	})
}

// GenerateSimpleBirthdayWish generates a birthday wish with minimal input - just name required
func (h *Handler) GenerateSimpleBirthdayWish(c *fiber.Ctx) error {
	name := c.Query("name")
	if name == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Name parameter is required"})
	}

	// Optional age parameter
	ageStr := c.Query("age")
	age := 0
	if ageStr != "" {
		if parsedAge, err := strconv.Atoi(ageStr); err == nil && parsedAge > 0 {
			age = parsedAge
		}
	}

	// If Gemini is not available, return a nice fallback wish
	if h.geminiClient == nil {
		fallbackWish := fmt.Sprintf("🎉 Happy Birthday, %s! 🎂 Wishing you all the joy, happiness, and wonderful surprises on your special day! May this year bring you endless blessings and amazing adventures! 🌟", name)
		return c.Status(200).JSON(fiber.Map{
			"name":   name,
			"wish":   fallbackWish,
			"source": "fallback",
		})
	}

	var wish string
	var err error

	if age > 0 {
		wish, err = h.geminiClient.GenerateBirthdayWish(name, age)
	} else {
		wish, err = h.geminiClient.GenerateGenericBirthdayWish(name)
		log.Println(err)
	}

	if err != nil {
		log.Printf("Error generating birthday wish: %v", err)
		fallbackWish := fmt.Sprintf("🎉 Happy Birthday, %s! 🎂 Wishing you all the joy, happiness, and wonderful surprises on your special day! May this year bring you endless blessings and amazing adventures! 🌟", name)
		return c.Status(200).JSON(fiber.Map{
			"name":   name,
			"wish":   fallbackWish,
			"source": "fallback",
		})
	}

	response := fiber.Map{
		"name":   name,
		"wish":   wish,
		"source": "gemini",
	}

	if age > 0 {
		response["age"] = age
	}

	return c.Status(200).JSON(response)
}
