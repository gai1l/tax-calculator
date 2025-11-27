package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	r.POST("/tax/calculations", func(c *gin.Context) {
		var input struct {
			TotalIncome float64 `json:"totalIncome"`
			WHT         float64 `json:"wht"`
			Allowances  []struct {
				AllowanceType string  `json:"allowanceType"`
				Amount        float64 `json:"amount"`
			} `json:"allowances"`
		}

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// === Donation deduction (สูงสุด 100,000) ===
		donation := 0.0
		for _, a := range input.Allowances {
			if a.AllowanceType == "donation" {
				if a.Amount > 100000 {
					donation = 100000
				} else {
					donation += a.Amount
				}
			}
		}

		// === รายได้สุทธิ ===
		taxable := input.TotalIncome - 60000 - donation
		if taxable < 0 {
			taxable = 0
		}

		// === คำนวณภาษีตามอัตราที่ถูกต้องของไทย ===
		var tax float64
		remaining := taxable

		if remaining > 2000000 {
			tax += (remaining - 2000000) * 0.35
			remaining = 2000000
		}
		if remaining > 1000000 {
			tax += (remaining - 1000000) * 0.20
			remaining = 1000000
		}
		if remaining > 500000 {
			tax += (remaining - 500000) * 0.15
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
