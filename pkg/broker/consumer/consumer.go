package consumer

import (
	"github.com/Shopify/sarama"
	log "github.com/acikkaynak/musahit-harita-backend/pkg/logger"
	"os"
	"strings"
)

func NewConsumer(groupId string) (sarama.ConsumerGroup, error) {
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")

	if kafkaBrokers == "" {
		log.Logger().Panic("KAFKA_BROKERS is not set")
	}

	brokers := strings.Split(kafkaBrokers, ",")

	cfg := sarama.NewConfig()

	return sarama.NewConsumerGroup(brokers, groupId, cfg)
}
