package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/example/pkg/kafka"
)

var kafkaProducer *kafka.Producer

func main() {
	// Initialize Kafka producer
	kafkaBrokers := []string{"kafka:9092"}
	if brokers := os.Getenv("KAFKA_BROKERS"); brokers != "" {
		kafkaBrokers = []string{brokers}
	}

	var err error
	kafkaProducer, err = kafka.NewProducer(kafka.Config{
		Brokers: kafkaBrokers,
		Topic:   "events",
	})
	if err != nil {
		log.Printf("Warning: Failed to initialize Kafka producer: %v\n", err)
	} else {
		defer kafkaProducer.Close()
		
		// Send test event
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		testEvent := kafka.Event{
			Type:      "gateway.startup",
			Timestamp: time.Now().Unix(),
			Data: map[string]string{
				"message": "API Gateway started successfully",
				"version": "v1",
			},
		}
		
		if err := kafkaProducer.PublishEvent(ctx, testEvent); err != nil {
			log.Printf("Failed to send test event: %v\n", err)
		}
	}

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/orders", func(c *gin.Context) {
		orderSvc := os.Getenv("ORDER_SERVICE_URL")
		if orderSvc == "" {
			orderSvc = "http://order-service:3000"
		}
		c.JSON(http.StatusOK, gin.H{"message": "Would proxy to " + orderSvc})
	})

	// Endpoint to send test events
	r.POST("/test-event", func(c *gin.Context) {
		if kafkaProducer == nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Kafka producer not available"})
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		testEvent := kafka.Event{
			Type:      "gateway.test",
			Timestamp: time.Now().Unix(),
			Data: map[string]string{
				"message": "Test event from API Gateway",
			},
		}
        err := kafkaProducer.PublishEvent(ctx, testEvent)

		if  err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Test event sent successfully"})
	})

	r.Run(":8080")
}

