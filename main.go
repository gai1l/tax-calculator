package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode) // เพื่อไม่ให้ขึ้น debug เยอะ
	r := gin.New()
	r.Use(gin.Recovery())

	r.POST("/tax/calculations", func(c *gin.Context) {
		var input struct {
			TotalIncome float64 `json:"totalIncome"`
			WHT         float64 `json:"wht"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		// หักลดหย่อนส่วนตัว 60,000
		taxable := input.TotalIncome - 60000

		// คำนวณภาษีแบบขั้นบันได (แบบง่ายสุดก่อน)
		var tax float64
		remaining := taxable

		if remaining > 2000000 {
			tax += (remaining - 2000000) * 0.35
			remaining = 2000000
		}
		if remaining > 1000000 {
			tax += (remaining - 1000000) * 0.30
			remaining = 1000000
		}
		if remaining > 500000 {
			tax += (remaining - 500000) * 0.20
			remaining = 500000
		}
		if remaining > 150000 {
			tax += (remaining - 150000) * 0.10
			remaining = 150000
		}

		finalTax := tax - input.WHT

		c.JSON(http.StatusOK, gin.H{"tax": finalTax})
	})

	r.Run(":8080")
}
