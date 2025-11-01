# Birthday Bot

A birthday management system with Telex integration for automated reminders and wishes.

## 📁 Project Structure

```
├── main.go                     # Main application entry point
├── internal/
│   ├── a2a/
│   │   └── logic.go            # Birthday reminder logic
│   └── agent/
│       ├── agent_card.json     # Agent capability definition
│       ├── agent.json          # Agent metadata
│       └── agent_card_loader.go # Agent card loading utilities
├── birthday_workflow.json      # Telex workflow definition
├── birthdays.json              # Birthday data storage
├── go.mod                      # Go module dependencies
└── README.md                   # Documentation
```

## ✨ Features

- **Birthday Storage**: JSON-based persistent storage for birthday data
- **Automated Reminders**: 24-hour advance notifications
- **Birthday Wishes**: Automated greetings on the actual day
- **Telex Integration**: A2A protocol support for chat platforms
- **REST API**: HTTP endpoints for birthday management
- **Webhooks**: Telex webhook support for event triggers

## Quick Start

1. **Build the project:**
   ```bash
   go build -o birthday-bot.exe .
   ```

2. **Run the server:**
   ```bash
   ./birthday-bot.exe
   ```

The server starts on port 3000 and provides a REST API for birthday management.

## API Usage

### Adding Birthdays

```bash
curl -X POST http://localhost:3000/api/birthdays \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "date": "2024-12-25"}'
```

### Listing Birthdays

```bash
curl http://localhost:3000/api/birthdays
```

### Manual Birthday Check

```bash
curl -X POST http://localhost:3000/api/trigger/birthday-check
```

### Message Types

The agent sends two types of messages:

#### 1. Reminder Messages (Day Before Birthday)
```json
{
  "kind": "message",
  "role": "agent",
  "parts": [
    {
      "kind": "text",
      "text": "Reminder: John's birthday is tomorrow! 🎉"
    }
  ],
  "messageId": "uuid",
  "metadata": {
    "birthday_id": "uuid",
    "birthday_name": "John",
    "type": "reminder"
  }
}
```

#### 2. Birthday Wishes (On Actual Birthday)
```json
{
  "kind": "message",
  "role": "agent", 
  "parts": [
    {
      "kind": "text",
      "text": "Happy Birthday, John! 🎉🎂 May this special day bring you endless joy, wonderful surprises, and all the happiness your heart can hold. Here's to another amazing year ahead!"
    }
  ],
  "messageId": "uuid",
  "metadata": {
    "birthday_id": "uuid",
    "birthday_name": "John",
    "type": "birthday_wish"
  }
}
```

## Telex Integration

The bot integrates with Telex via the A2A protocol:

### Agent Discovery
```bash
curl http://localhost:3000/.well-known/agent.json
```

### Telex Webhooks
```bash
curl -X POST http://localhost:3000/api/telex/webhook \
  -H "Content-Type: application/json" \
  -d '{"event": "daily_check", "data": {}}'
```

### A2A Messages
```bash
curl -X POST http://localhost:3000/api/a2a/message \
  -H "Content-Type: application/json" \
  -d '{"from": "telex", "to": "birthday_bot", "content": "check birthdays", "type": "command"}'
```

## Configuration

### Agent Card (`internal/agent/agent_card.json`)
Defines the bot's capabilities for Telex integration:
- Birthday management functionality
- A2A protocol endpoints
- Input/output modes
- Skill definitions

### Data Storage
- **Format**: JSON file storage
- **Location**: `birthdays.json` in project root  
- **Structure**: Array of birthday objects with ID, name, month, day, and timestamps

## Development

### Available Endpoints

```
GET  /health                     - Health check
GET  /.well-known/agent.json     - Agent discovery
POST /api/birthdays              - Add birthday
GET  /api/birthdays              - List birthdays
GET  /api/birthdays/today        - Today's birthdays
GET  /api/birthdays/upcoming     - Upcoming birthdays (30 days)
POST /api/a2a/message            - A2A message processing
POST /api/telex/webhook          - Telex webhook handler
POST /api/trigger/birthday-check - Manual birthday check
```

### Dependencies

- `github.com/gofiber/fiber/v2` - HTTP web framework
- `github.com/google/uuid` - UUID generation
- Standard Go libraries for HTTP, JSON, and time handling
