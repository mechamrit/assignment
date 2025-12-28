package audit

import (
	"context"
	"encoding/json"
	"log"

	"backend/config"
	"backend/models"

	"github.com/segmentio/kafka-go"
)

type Service struct {
	writer *kafka.Writer
}

func New(cfg *config.Config) *Service {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(cfg.KafkaURL),
		Topic:    "qc-audit-logs",
		Balancer: &kafka.LeastBytes{},
	}
	log.Println("Kafka producer initialized")
	return &Service{writer: writer}
}

func (s *Service) ProduceAuditLog(workflowLog models.WorkflowLog) {
	go func() {
		data, err := json.Marshal(workflowLog)
		if err != nil {
			log.Printf("Failed to marshal audit log: %v", err)
			return
		}

		err = s.writer.WriteMessages(context.Background(),
			kafka.Message{
				Value: data,
			},
		)
		if err != nil {
			log.Printf("Failed to write to Kafka: %v", err)
		}
	}()
}

func (s *Service) Shutdown() {
	if s.writer != nil {
		s.writer.Close()
	}
}
