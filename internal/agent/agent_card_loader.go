package agent

import (
	"errors"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2/log"
)

var AgentCardData []byte

// LoadAgentCard opens the agent card file and returns its content as bytes
func loadAgentCard(filePath string) ([]byte, error) {
	// Read the entire file content
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Error("Failed to read agent card file:", err)
		return nil, fmt.Errorf("failed to read agent card file: %w", err)
	}

	// Store in global variable for caching
	AgentCardData = data

	log.Info("Successfully loaded agent card file:", filePath)
	return data, nil
}

// LoadDefaultAgentCard attempts to load agent card from common file locations
func LoadDefaultAgentCard() ([]byte, error) {
	// Common agent card file paths - check in order of preference
	possiblePaths := []string{
		"./internal/agent/agent_card.json",
		"./internal/agent/agent.json",
		"./agent_card.json",
		"./.well-known/agent.json",
	}

	for _, path := range possiblePaths {
		data, err := loadAgentCard(path)
		if err == nil {
			log.Info("Found agent card at:", path)
			return data, nil
		}
	}

	return nil, errors.New("no agent card file found in any of the expected locations")
}

func CheckForAgentCard() error {
	// Try to load using the default loader first
	data, err := LoadDefaultAgentCard()
	if err != nil {
		log.Error(err)
		return errors.New("failed to load agent card file")
	}

	// Store the loaded data
	AgentCardData = data

	if len(AgentCardData) == 0 {
		return fmt.Errorf("agent card is empty")
	}
	return nil
}
