# 🎂 Hazel the Birthday Bot

> An AI-powered birthday management agent that integrates seamlessly with Telex.im using the A2A protocol

Hazel is an intelligent birthday bot built for HNG Internship Stage 3 Backend Task. She remembers birthdays so you don't have to, generates personalized AI-powered birthday wishes using Google Gemini, and lives right in your chat through Telex's A2A (Agent-to-Agent) protocol.

## 🌟 Why Hazel?

**The Problem**: We all have that friend who never forgets a birthday, and another who couldn't remember their own anniversary. Birthday forgetfulness is universally relatable – and painful.

**The Solution**: An AI agent that combines natural language processing, persistent storage, and Google Gemini AI to create a birthday assistant that actually feels human.

**The Magic**: Just tell Hazel "remember my birthday 2005-01-01" in plain English, and she'll handle everything from storage to generating personalized wishes when the time comes.

## ✨ Features

### 🧠 **AI-Powered Intelligence**
- **Google Gemini Integration**: Generates personalized, thoughtful birthday wishes
- **Natural Language Processing**: Understands various ways people express birthday requests
- **Smart Text Analysis**: Extracts dates and context from conversational input
- **Graceful Fallback**: Beautiful default messages when AI is unavailable

### 💬 **Seamless Chat Integration**
- **Telex A2A Protocol**: Full JSON-RPC 2.0 compliant integration
- **Real-time Messaging**: Instant responses through chat interface
- **Agent Discovery**: Proper `.well-known/agent.json` implementation
- **Workflow Support**: Compatible with Telex workflow automation

### 🎯 **Smart Birthday Management**
- **Natural Commands**: "remember my birthday", "list birthdays", "upcoming birthdays"
- **Flexible Date Formats**: Handles YYYY-MM-DD and various input styles
- **Upcoming Notifications**: Shows birthdays in the next 30 days
- **Persistent Storage**: Survives container restarts and deployments

### 🔧 **Production Ready**
- **Comprehensive Logging**: Debug-friendly with detailed request/response tracking
- **Error Handling**: Graceful degradation with user-friendly error messages
- **Timeout Management**: Prevents hanging on external API calls
- **Container Aware**: Designed for ephemeral container environments like Render

## 🚀 Live Demo

**Deployed at**: [hazel-agent.onrender.com](https://hazel-agent.onrender.com)

**Try it in Telex**: Search for "Hazel the Birthday Bot" in your Telex colleagues

### Quick Test Commands:
```
remember my birthday 2005-01-01
list birthdays
generate a birthday wish for Alice
list upcoming birthdays
```

## 📁 Project Architecture

```
hazel_agent/
├── main.go                     # Application entry point & server setup
├── internal/
│   ├── handlers/
│   │   └── handlers.go         # HTTP handlers & A2A message processing
│   ├── clients/
│   │   └── gemini.go          # Google Gemini AI client integration
│   ├── store/
│   │   └── store.go           # Birthday data storage & persistence
│   └── agent/
│       ├── agent_card.json    # Telex agent capability definition
│       └── agent_card_loader.go # Agent card loading utilities
├── birthday_workflow.json      # Telex workflow configuration
├── Dockerfile                  # Container deployment configuration
├── go.mod                      # Go module dependencies
└── .env.example               # Environment variables template
```

## 🛠 Quick Start

### Prerequisites
- **Go 1.21+** installed
- **Google Gemini API Key** (for AI birthday wishes)
- **Telex.im account** (for chat integration)

### Environment Setup

1. **Clone and navigate:**
   ```bash
   git clone https://github.com/whotterre/hazel_agent.git
   cd hazel_agent
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Set up environment variables:**
   ```bash
   cp .env.example .env
   # Edit .env and add your GEMINI_API_KEY
   ```

4. **Build and run:**
   ```bash
   go build -o hazel_bot
   ./hazel_bot
   ```

The server starts on port 3000 (or PORT environment variable).

### 🔑 Environment Variables

Create a `.env` file with:
```env
GEMINI_API_KEY=your_google_gemini_api_key_here
PORT=3000
```

## 🔌 API Reference

### Core Endpoints

#### **Agent Discovery**
```http
GET /.well-known/agent.json
```
Returns the Telex agent card defining Hazel's capabilities.

#### **Health Check**
```http
GET /health
```
Basic health check endpoint.

#### **Telex A2A Integration**
```http
POST /
Content-Type: application/json

{
  "jsonrpc": "2.0",
  "id": "unique-id",
  "method": "message/send",
  "params": {
    "message": {
      "kind": "message", 
      "role": "user",
      "parts": [{"kind": "text", "text": "remember my birthday 2005-01-01"}]
    }
  }
}
```

### REST API Endpoints

#### **Add Birthday**
```bash
  -H "Content-Type: application/json" \
  -d '{"name": "Alice", "date": "2005-01-01"}'
```

#### **List All Birthdays**
```bash
curl http://localhost:3000/api/birthdays
```

#### **Get Today's Birthdays**  
```bash
curl http://localhost:3000/api/birthdays/today
```

#### **Get Upcoming Birthdays**
```bash
curl http://localhost:3000/api/birthdays/upcoming
```

#### **Generate Birthday Wish**
```bash
curl -X POST http://localhost:3000/api/wishes/generate \
  -H "Content-Type: application/json" \
  -d '{"name": "Alice", "age": 25}'
```

## 🤖 How Hazel Works

### Natural Language Processing
Hazel understands various ways people might express birthday-related requests:

```
✅ "remember my birthday 2005-01-01"
✅ "my birthday is January 1st 2005"  
✅ "remember my birthday - 2003-09-09"
✅ "list upcoming birthdays"
✅ "generate a birthday wish for Alice"
```

### AI-Powered Responses
Using Google Gemini AI, Hazel generates personalized birthday wishes:

**Input**: "generate a birthday wish for Alice"
**Output**: "🎉 Happy Birthday, Alice! May your special day sparkle with joy, laughter, and all the wonderful surprises life has to offer. Here's to another year of amazing adventures and cherished memories! 🌟�"

### A2A Protocol Integration
Hazel communicates with Telex using JSON-RPC 2.0:

```json
{
  "jsonrpc": "2.0",
  "id": "request-id", 
  "result": {
    "message": {
      "kind": "message",
      "role": "assistant",
      "parts": [{"kind": "text", "text": "🎂 Perfect! I've remembered your birthday..."}]
    }
  }
}
```

## 🏗 Architecture & Design

### Tech Stack
- **Backend**: Go with Fiber web framework
- **AI Integration**: Google Gemini AI API
- **Storage**: JSON file-based persistence  
- **Protocol**: JSON-RPC 2.0 for A2A communication
- **Deployment**: Docker container on Render
- **Environment**: Supports development and production configurations

### Key Design Decisions

#### **Go + Fiber Framework**
- Excellent concurrency for handling multiple AI calls and chat messages
- Lightweight and fast HTTP handling
- Simple, Express.js-like API design

#### **Google Gemini AI**  
- State-of-the-art language model for natural birthday wishes
- 10-second timeout to prevent hanging
- Graceful fallback to pre-written messages

#### **A2A Protocol Implementation**
- Full JSON-RPC 2.0 compliance for Telex integration
- Single POST endpoint routing (Telex requirement)
- Proper agent card discovery mechanism

#### **Container-Aware Architecture**
- Handles ephemeral container storage limitations
- Comprehensive logging for debugging in production
- Environment variable configuration

### Message Flow
```
User Input → Natural Language Processing → Intent Detection → Action Router → Response Generation → AI Enhancement → Telex Response
```

## 🚀 Deployment

### Docker Deployment
```bash
# Build Docker image
docker build -t hazel-bot .

# Run container
docker run -p 3000:3000 -e GEMINI_API_KEY=your_key hazel-bot
```

### Render Deployment
1. Connect your GitHub repository to Render
2. Set environment variables in Render dashboard
3. Deploy automatically on git push

### Environment Variables
```env
GEMINI_API_KEY=your_google_gemini_api_key
PORT=3000                    # Optional, defaults to 3000
```

## 🧪 Testing

### Manual Testing with Telex
1. Search for "Hazel the Birthday Bot" in Telex colleagues
2. Try these commands:
   - `remember my birthday 2005-01-01`
   - `list birthdays`
   - `generate a birthday wish for Alice`
   - `list upcoming birthdays`

### API Testing
```bash
# Test health endpoint
curl http://localhost:3000/health

# Test agent discovery
curl http://localhost:3000/.well-known/agent.json

# Test birthday storage
curl -X POST http://localhost:3000/api/birthdays \
  -H "Content-Type: application/json" \
  -d '{"name": "Test User", "date": "2000-01-01"}'
```

## 📚 Dependencies

```go
require (
    github.com/gofiber/fiber/v2 v2.52.9
    github.com/google/uuid v1.6.0
    github.com/joho/godotenv v1.5.1
    google.golang.org/genai v0.18.0
)
```

## 🐛 Troubleshooting

### Common Issues

#### **"Error while streaming" in Telex**
- **Cause**: Missing Content-Type headers or AI timeout
- **Solution**: Check logs for timeout errors, verify Gemini API key

#### **Birthdays disappearing**  
- **Cause**: Container restart wiped ephemeral storage
- **Solution**: Expected behavior on Render, data persists during active sessions

#### **Agent not responding in Telex**
- **Cause**: Agent card URL incorrect or server down
- **Solution**: Verify deployment URL matches agent card configuration

### Debug Mode
Enable detailed logging by checking server logs after requests.

## 🎯 HNG Stage 3 Requirements

This project fulfills all HNG Internship Stage 3 Backend Task requirements:

✅ **Working AI Agent**: Google Gemini integration with natural language processing  
✅ **Telex.im Integration**: Full A2A protocol implementation with JSON-RPC 2.0  
✅ **Live Demo**: Deployed at hazel-agent.onrender.com  
✅ **Clean API**: RESTful endpoints with proper error handling  
✅ **Documentation**: Comprehensive README and inline code documentation  
✅ **Creativity**: Solves universal problem of birthday forgetfulness  
✅ **Error Handling**: Graceful degradation and comprehensive logging

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)  
5. Open a Pull Request

## 📝 License

This project is part of the HNG Internship Stage 3 Backend Task.

## 🙏 Acknowledgments

- **HNG Internship** for the challenging and educational task requirements
- **Telex.im** for providing the A2A protocol and platform integration
- **Google Gemini AI** for enabling personalized birthday wish generation
- **Go Community** for excellent libraries and documentation

---

**Built with ❤️ for HNG Internship Stage 3 Backend Task**

*Making birthdays memorable, one AI-generated wish at a time.*
