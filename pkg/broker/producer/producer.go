package producer

import (
	"github.com/Shopify/sarama"
	log "github.com/acikkaynak/musahit-harita-backend/pkg/logger"
	"os"
	"strings"
)

func NewProducer() (sarama.SyncProducer, error) {
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")

	if kafkaBrokers == "" {
		log.Logger().Panic("KAFKA_BROKERS is not set")
	}

	brokers := strings.Split(kafkaBrokers, ",")

	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true

	return sarama.NewSyncProducer(brokers, cfg)

}
