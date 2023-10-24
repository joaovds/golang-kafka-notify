package main

import (
  "fmt"
  "kafka-notify/internals/configuration"
  "kafka-notify/internals/infra/kafka"
  "kafka-notify/pkg/models"
  "log"
  "net/http"

  "github.com/gin-gonic/gin"
)

func main () {
  users := []models.User{
    { ID: 1, Name: "John" },
    { ID: 2, Name: "Nicoly" },
    { ID: 3, Name: "Peter" },
    { ID: 4, Name: "Jane" },
  }

  producer, err := kafka.SetupProducer()
  if err != nil {
    log.Fatalf("failed to setup kafka producer: %v", err)
  }
  defer producer.Close()

  gin.SetMode(gin.ReleaseMode)
  router := gin.Default()

  router.GET("/users", func(c *gin.Context) {
    c.JSON(http.StatusOK, users)
  })

  fmt.Printf("Kafka PRODUCER running at http://localhost%s\n", configuration.ProducerPort)

  if err := router.Run(configuration.ProducerPort); err != nil {
    log.Printf("Failed to start server: %v\n", err)
  }
}

