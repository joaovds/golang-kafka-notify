package kafka

import (
  "github.com/IBM/sarama"

  "kafka-notify/"
)

func SetupProducer() (sarama.SyncProducer, error) {
  config := sarama.NewConfig()
  config.Producer.Return.Successes = true

  producer, err := sarama.NewSyncProducer([]string{
    main.KafkaBroker,
  })
}

