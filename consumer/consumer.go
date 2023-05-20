package consumer

import "github.com/acikkaynak/musahit-harita-backend/repository"

const (
	observerTopicName = "topic.election.observer"
)

type Consumer struct {
	Ready chan bool
	repo  *repository.Repository
}

func NewConsumer() *Consumer {

	return &Consumer{
		Ready: make(chan bool),
		repo:  repository.New(),
	}

}
