package main

import "github.com/gin-gonic/gin"

func main() {
    r := gin.Default()
    r.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "Tax Calculator API Ready!"})
    })
    r.Run(":808 บ") // รันที่ port 8080
}