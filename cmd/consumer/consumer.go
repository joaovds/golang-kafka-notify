package main

import (
	"context"
	"errors"
	"fmt"
	"kafka-notify/internals/configuration"
	"kafka-notify/internals/infra/kafka"
	"kafka-notify/pkg/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
  store := models.NewNotificationStore()

  ctx, cancel := context.WithCancel(context.Background())
  go kafka.SetupConsumer(ctx, store)
  defer cancel()

  gin.SetMode(gin.ReleaseMode)
  router := gin.Default()
  router.GET("/notifications/:userID", func(c *gin.Context) {
    handleNotifications(c, store)
  })

  fmt.Printf("Kafka CONSUMER (Group: %s) 👥📥 "+
    "started at http://localhost%s\n", configuration.ConsumerGroup, configuration.ConsumerPort)

  if err := router.Run(configuration.ConsumerPort); err != nil {
    log.Printf("failed to run the server: %v", err)
  }
}

var ErrNoMessagesFound = errors.New("no messages found")

func getUserIDFromRequest(ctx *gin.Context) (string, error) {
  userID := ctx.Param("userID")
  if userID == "" {
    return "", ErrNoMessagesFound
  }
  return userID, nil
}

func handleNotifications(ctx *gin.Context, store *models.NotificationStore) {
  userID, err := getUserIDFromRequest(ctx)
  if err != nil {
    ctx.JSON(http.StatusNotFound, gin.H{ "message": err.Error() })
    return
  }

  notes := store.Get(userID)
  if len(notes) == 0 {
    ctx.JSON(http.StatusOK, gin.H{
      "message": "no notifications found for user",
      "notifications": []models.Notification{},
    })
    return
  }

  ctx.JSON(http.StatusOK, gin.H{ "notifications": notes })
}

