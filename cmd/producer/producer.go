package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"kafka-notify/internals/configuration"
	"kafka-notify/internals/infra/kafka"
	"kafka-notify/pkg/models"
	"log"
	"net/http"
	"strconv"

	"github.com/IBM/sarama"
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
  router.POST("/send", sendMessageHandler(producer, users))

  fmt.Printf("Kafka PRODUCER running at http://localhost%s\n", configuration.ProducerPort)

  if err := router.Run(configuration.ProducerPort); err != nil {
    log.Printf("Failed to start server: %v\n", err)
  }
}

func sendMessageHandler(producer sarama.SyncProducer, users []models.User) gin.HandlerFunc {
  return func(ctx *gin.Context) {
    fromID, err := getIDFromRequest("fromID", ctx)
    if err != nil {
      ctx.JSON(http.StatusBadRequest, gin.H{ "message": err.Error()})
      return
    }

    toID, err := getIDFromRequest("toID", ctx)
    if err != nil {
      ctx.JSON(http.StatusBadRequest, gin.H{ "message": err.Error()})
      return
    }

    err = sendKafkaMessage(producer, users, ctx, fromID, toID)
    if errors.Is(err, ErrUserNotFoundInProducer) {
      ctx.JSON(http.StatusNotFound, gin.H{ "message": "user not found" })
      return
    }
    if err != nil {
      ctx.JSON(
        http.StatusInternalServerError,
        gin.H{ "message": fmt.Sprintf("failed to send message: %v", err) },
        )
      return
    }

    ctx.JSON(http.StatusOK, gin.H{ "message": "Notification sent successfully" })
  }
}

var ErrUserNotFoundInProducer = errors.New("user not found in producer")

func findUserById(userId int, users []models.User) (models.User, error) {
  for _, user := range users {
    if user.ID == userId {
      return user, nil
    }
  }
  return models.User{}, ErrUserNotFoundInProducer
}

func getIDFromRequest(formKey string, ctx *gin.Context) (int, error) {
  fmt.Println(ctx.PostForm(formKey))
  id, err := strconv.Atoi(ctx.PostForm(formKey))
  if err != nil {
    return 0, fmt.Errorf("failed to convert %s to int: %v", formKey, err)
  }

  return id, nil
}

func sendKafkaMessage(producer sarama.SyncProducer, users []models.User, ctx *gin.Context, fromID, toID int) error {
  message := ctx.PostForm("message")

  fromUser, err := findUserById(fromID, users)
  if err != nil {
    return err
  }

  toUser, err := findUserById(toID, users)
  if err != nil {
    return err
  }

  notification := models.Notification{
    From: fromUser,
    To: toUser,
    Message: message,
  }

  notificationJSON, err := json.Marshal(notification)
  if err != nil {
    return fmt.Errorf("failed to marshal notification to JSON: %v", err)
  }

  msg := &sarama.ProducerMessage{
    Topic: configuration.KafkaTopic,
    Key: sarama.StringEncoder(strconv.Itoa(toUser.ID)),
    Value: sarama.StringEncoder(notificationJSON),
  }

  _, _, err = producer.SendMessage(msg)
  return err
}

