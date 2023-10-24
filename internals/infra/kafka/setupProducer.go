package kafka

import (
  "fmt"
  "kafka-notify/internals/configuration"

  "github.com/IBM/sarama"
)

func SetupProducer() (sarama.SyncProducer, error) {
  config := sarama.NewConfig()
  config.Producer.Return.Successes = true

  producer, err := sarama.NewSyncProducer(
    []string{configuration.KafkaBroker},
    config,
    )
  if err != nil {
    return nil, fmt.Errorf("failed to setup kafka producer: %w", err)
  }

  return producer, nil
}

