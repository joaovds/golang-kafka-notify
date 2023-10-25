package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"kafka-notify/internals/configuration"
	"kafka-notify/pkg/models"
	"log"

	"github.com/IBM/sarama"
)

type Consumer struct {
  store *models.NotificationStore
}

func (*Consumer) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (*Consumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }
func (consumer *Consumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
  for msg := range claim.Messages() {
    userID := string(msg.Key)
    var notification models.Notification

    err := json.Unmarshal(msg.Value, &notification)
    if err != nil {
      log.Printf("failed to unmarshal notification: %v", err)
      continue
    }

    consumer.store.Add(userID, notification)
    sess.MarkMessage(msg, "")
  }
  return nil
}

func SetupConsumer(ctx context.Context, store *models.NotificationStore) {
  consumerGroup, err := initializeConsumerGroup()
  if err != nil {
    log.Printf("initializing consumer group error: %v", err)
  }
  defer consumerGroup.Close()

  consumer := &Consumer{
    store: store,
  }

  for {
    err = consumerGroup.Consume(ctx, []string{configuration.KafkaTopic}, consumer)
  }
}

func initializeConsumerGroup() (sarama.ConsumerGroup, error) {
  config := sarama.NewConfig()

  consumerGroup, err := sarama.NewConsumerGroup(
    []string{configuration.KafkaBroker},
    configuration.ConsumerGroup,
    config,
  )
  if err != nil {
    return nil, fmt.Errorf("failed to initialize consumer group: %w", err)
  }

  return consumerGroup, nil
}

