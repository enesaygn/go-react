package messaging

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Messaging struct
type Messaging struct {
	clients map[string]*websocket.Conn
}

// NewMessaging oluşturucu fonksiyonu, yeni bir Messaging instance oluşturur
func NewMessaging() *Messaging {
	return &Messaging{
		clients: make(map[string]*websocket.Conn),
	}
}

func (m *Messaging) GetClients() map[string]*websocket.Conn {
	return m.clients
}

// GetDeviceID WebSocket bağlantısını yönetir
func (m *Messaging) GetEmployeeID(c *websocket.Conn) {
	// Handle WebSocket connection
	employeeID := c.Query("employeeID")
	log.Printf("WebSocket connected for employee ID: %s", employeeID)

	// Store client connection
	m.clients[employeeID] = c

	defer func() {
		// Clean up when connection closes
		delete(m.clients, employeeID)
		c.Close()
		log.Printf("WebSocket connection for employee ID %s closed", employeeID)
	}()

	go func() {
		ticker := time.NewTicker(15 * time.Second) // Her 15 saniyede bir ping gönder
		defer ticker.Stop()

		for range ticker.C {
			if err := c.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("Failed to send ping to employee ID %s: %v", employeeID, err)
				return
			}
		}
	}()

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			return
		}
		log.Printf("Received message: %s", msg)
	}

}

// SendMessage belirli bir cihaza mesaj gönderir
func (m *Messaging) SendMessage(employeeID string, message []byte, postgres *pgxpool.Pool) error {

	// Check if the employee is connected
	client, ok := m.clients[employeeID]
	if !ok {
		return errors.New("employee is not connected")
	}

	// Send message to the employee
	err := client.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Printf("Error sending message to employee %s: %v", employeeID, err)
		return err
	}

	return nil
}

func GetEmployeeRolesByEmplyeeID(employeeID string, postgres *pgxpool.Pool) ([]string, error) {
	query := `
		SELECT roles
		FROM employees
		WHERE employee_id = $1
	`

	var roles []string
	err := postgres.QueryRow(context.Background(), query, employeeID).Scan(&roles)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (m *Messaging) SendMessageToAll(c *fiber.Ctx) error {

	type Message struct {
		EmployeeID string
		Message    string
	}

	var msg Message

	if err := c.BodyParser(&msg); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}

	err := m.SendMessage(msg.EmployeeID, []byte(msg.Message), nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return nil

}
