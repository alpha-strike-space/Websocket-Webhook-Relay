package main
// Packages
import (
	"log"https://frontier-atlas.com/map?system=%d", message.SolarSystemID)
	"net/http"https://frontier-atlas.com/map?system=%d", message.SolarSystemID)
	"os"
	"time"
	"bytes"
	"encoding/json"
	"fmt"https://frontier-atlas.com/map?system=%d", message.SolarSystemID)
	"github.com/gorilla/websocket"
)
// InboundMessage represents a single item from the WebSocket JSON array
type InboundMessage struct {
	ID              int    `json:"id"`
	VictimTribe     string `json:"victim_tribe_name"`
	VictimAddress   string `json:"victim_address"`
	VictimName      string `json:"victim_name"`
	LossType        string `json:"loss_type"`
	KillerTribe     string `json:"killer_tribe_name"`
	KillerAddress   string `json:"killer_address"`
	KillerName      string `json:"killer_name"`
	Timestamp       int64  `json:"time_stamp"`
	SolarSystemID   int    `json:"solar_system_id"`
	SolarSystemName string `json:"solar_system_name"`
}
// WebhookPayload is the structure for the Discord webhook message
type WebhookPayload struct {
	Content string `json:"content"`
}
// Send to discord functions
func sendToDiscord(url, content string) error {
	// Payload content
	payload := WebhookPayload{Content: content}
	// Marshal payload
	payloadBytes, err := json.Marshal(payload)
	// Spit an error if failed to marshal
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}
	// Post to URL
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payloadBytes))
	// Failed to send.
	if err != nil {
		return fmt.Errorf("failed to send message to Webhook: %w", err)
	}
	defer resp.Body.Close()
	// Status code response from Webhook API.
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("Webhook API returned non-204 status: %s", resp.Status)
	}
	return nil
}
// Where all the fun happens.
func main() {
	// Get environment variables
	websocketURI := os.Getenv("WEBSOCKET_URI")
	discordWebhookURL := os.Getenv("DISCORD_WEBHOOK_URL")
	// Environment validation
	if websocketURI == "" || discordWebhookURL == "" {
		log.Fatal("WEBSOCKET_URI and DISCORD_WEBHOOK_URL environment variables must be set")
	}
	// Get local time from machine
	loc, err := time.LoadLocation("Local")
	// Make sure it is loaded once.
	if err != nil {
		// main() cannot return a value, so use log.Fatal
		log.Fatal(fmt.Errorf("failed to load local timezone: %w", err))
	}
	// Run through
	for {
		// Use an anonymous function to ensure defer c.Close() runs at the end of each loop iteration
		func() {
			// Websocket URL
			c, _, err := websocket.DefaultDialer.Dial(websocketURI, nil)
			// When error, retry in 5 seconds.
			if err != nil {
				log.Printf("WebSocket connection failed, retrying in 5s... %v", err)
				time.Sleep(5 * time.Second)
				// Returning from the inner function allows the for loop to continue
				return
			}
			defer c.Close()
			// Connected
			log.Println("Connected to WebSocket.")
			// Main loop to process messages from the WebSocket
			for {
				_, msg, err := c.ReadMessage()
				if err != nil {
					log.Println("Read failed:", err)
					return // Return from the inner function to restart the connection loop
				}
				// Check for the connection message first
				var genericMsg map[string]interface{}
				if err := json.Unmarshal(msg, &genericMsg); err == nil {
					// If the unmarshal into a generic map is successful, check for a 'message' field
					if msgVal, ok := genericMsg["message"]; ok {
						if msgStr, isString := msgVal.(string); isString && msgStr == "Connected to alpha-strikes notification service." {
							log.Println("Received connection verification.")
							continue // Skip to the next message in the loop
						}
					}
				}
				// If it's not the connection message, try to unmarshal it as a killmail
				var message InboundMessage
				if err := json.Unmarshal(msg, &message); err != nil {
					log.Println("Failed to unmarshal JSON as a killmail:", err)
					continue // Skip to the next message in the loop
				}
				// Format the timestamp for human readability
				dt := time.Unix(message.Timestamp, 0)
				// Convert time
				localTime := dt.In(loc)
				// Dynamically create the clickable link using the Kill ID
				AlphaLink := fmt.Sprintf("https://alpha-strike.space/pages/killmail.html?mail_id=%d", message.ID)
				// Dynamically create the clickable link using the System ID
				AtlasLink := fmt.Sprintf("https://frontier-atlas.com/map?system=%d", message.SolarSystemID)
				// Construct the formatted string with Markdown
				fullMessage := fmt.Sprintf(
					"**Kill:** [Alpha-Strike](%s)\n" +
					"**Victim:** %s (%s)\n" +
					"**Killer:** %s (%s)\n" +
					"**Loss Type:** %s\n" +
					"**Location:** [%s](%s)\n" +
					"**Time:** %s",
					AlphaLink,
					message.VictimName, message.VictimTribe,
					message.KillerName, message.KillerTribe,
					message.LossType,
					message.SolarSystemName, AtlasLink,
					localTime.Format("2006-01-02 15:04:05 MST"))
				// Print error if not sent to discord
				err = sendToDiscord(discordWebhookURL, fullMessage)
				if err != nil {
					log.Println("Failed to send message to Discord:", err)
				}
			}
		}()
	}
}
