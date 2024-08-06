package services

import (
	"log"
	"sasa-elterminali-service/internal/messaging"

	"github.com/streadway/amqp"
)

type RabbitMQConsumerService struct {
	rabbitMQ *messaging.RabbitMQ
}

func NewRabbitMQConsumerService(rabbitMQ *messaging.RabbitMQ) *RabbitMQConsumerService {
	return &RabbitMQConsumerService{rabbitMQ: rabbitMQ}
}

func (s *RabbitMQConsumerService) StartConsumers() error {
	// Burada tüm consumer'larınızı başlatabilirsiniz
	err := s.rabbitMQ.Consume("employee_queue", s.handleEmployeeMessage)
	if err != nil {
		return err
	}

	// İleride farklı queue'lar eklediğinizde burada yeni consumer'lar başlatabilirsiniz
	// err = s.rabbitMQ.Consume("another_queue", s.handleAnotherMessage)
	// if err != nil {
	//     return err
	// }

	return nil
}

func (s *RabbitMQConsumerService) handleEmployeeMessage(msg amqp.Delivery) {
	log.Printf("Received a message: %s", msg.Body)
	// Mesajı işleyin
}

// İleride farklı queue'lar için handler fonksiyonları ekleyebilirsiniz
// func (s *RabbitMQConsumerService) handleAnotherMessage(msg amqp.Delivery) {
//     log.Printf("Received another message: %s", msg.Body)
//     // Mesajı işleyin
// }
