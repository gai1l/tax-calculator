package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Allowance struct {
	AllowanceType string  `json:"allowanceType"`
	Amount        float64 `json:"amount"`
}

type TaxRequest struct {
	TotalIncome float64     `json:"totalIncome"`
	WHT         float64     `json:"wht"`
	Allowances  []Allowance `json:"allowances"`
}

type TaxLevel struct {
	Level string  `json:"level"`
	Tax   float64 `json:"tax"`
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	r.POST("/tax/calculations", func(c *gin.Context) {
		var req TaxRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
			return
		}

		// === Validation ทั้งหมดตาม Test Case 4-7 ===
		if req.TotalIncome < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "totalIncome cannot be negative"})
			return
		}
		if req.WHT < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "wht cannot be negative"})
			return
		}
		if req.WHT > req.TotalIncome {
			c.JSON(http.StatusBadRequest, gin.H{"error": "wht cannot be greater than totalIncome"})
			return
		}
		for _, a := range req.Allowances {
			if a.AllowanceType != "donation" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "allowanceType must be 'donation'"})
				return
			}
			if a.Amount < 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "amount cannot be negative"})
				return
			}
		}

		// === Donation deduction (สูงสุด 100,000) ===
		donation := 0.0
		for _, a := range req.Allowances {
			if a.AllowanceType == "donation" {
				if a.Amount > 100000 {
					donation += 100000
				} else {
					donation += a.Amount
				}
			}
		}

		// === คำนวณรายได้สุทธิ ===
		taxable := req.TotalIncome - 60000 - donation
		if taxable < 0 {
			taxable = 0
		}

		// === คำนวณภาษีทีละขั้น + เก็บ taxLevel ===
		var tax float64
		levels := []TaxLevel{
			{Level: "0-150,000", Tax: 0},
			{Level: "150,001-500,000", Tax: 0},
			{Level: "500,001-1,000,000", Tax: 0},
			{Level: "1,000,001-2,000,000", Tax: 0},
			{Level: "2,000,001 ขึ้นไป", Tax: 0},
		}

		remaining := taxable

		// ขั้น 2,000,001 ขึ้นไป
		if remaining > 2000000 {
			excess := remaining - 2000000
			levels[4].Tax = excess * 0.35
			tax += levels[4].Tax
			remaining = 2000000
		}
		// ขั้น 1,000,001-2,000,000
		if remaining > 1000000 {
			excess := remaining - 1000000
			levels[3].Tax = excess * 0.30
			tax += levels[3].Tax
			remaining = 1000000
		}
		// ขั้น 500,001-1,000,000
		if remaining > 500000 {
			excess := remaining - 500000
			levels[2].Tax = excess * 0.20
			tax += levels[2].Tax
			remaining = 500000
		}
		// ขั้น 150,001-500,000
		if remaining > 150000 {
			excess := remaining - 150000
			levels[1].Tax = excess * 0.10
			tax += levels[1].Tax
			remaining = 150000
		}
		// 0-150,000 = 0%

		finalTax := tax - req.WHT

		// Response พร้อม taxLevel
		c.JSON(http.StatusOK, gin.H{
			"tax":      finalTax,
			"taxLevel": levels,
		})
	})

	r.Run(":8080")
}
