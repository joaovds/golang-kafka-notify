package main

import (
	"fmt"
	"kafka-notify/pkg/models"
  "log"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
  ProducerPort = ":8080"
  KafkaBroker = "localhost:9092"
  KafkaTopic = "notifications"
)

func main () {
  users := []models.User{
    { ID: 1, Name: "John" },
    { ID: 2, Name: "Nicoly" },
    { ID: 3, Name: "Peter" },
    { ID: 4, Name: "Jane" },
  }

  gin.SetMode(gin.ReleaseMode)
  router := gin.Default()

  router.GET("/users", func(c *gin.Context) {
    c.JSON(http.StatusOK, users)
  })

  fmt.Printf("Kafka PRODUCER running at http://localhost%s\n", ProducerPort)

  if err := router.Run(ProducerPort); err != nil {
    log.Printf("Failed to start server: %v\n", err)
  }
}

