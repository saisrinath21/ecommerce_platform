package main

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    r.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": "ok"})
    })

    r.POST("/orders", func(c *gin.Context) {
        var payload map[string]interface{}
        if err := c.BindJSON(&payload); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }
        c.JSON(http.StatusCreated, gin.H{"order": payload, "id": "order_123"})
    })

    r.GET("/orders", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"orders": []interface{}{}})
    })

    r.Run(":3000")
}
