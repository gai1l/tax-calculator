# Tax Calculator API 

REST API สำหรับคำนวณภาษีเงินได้บุคคลธรรมดา

## เทคโนโลยีที่ใช้
- Go 1.21+
- Gin Web Framework

## วิธีรัน
```bash
go run main.go
API รันที่ http://localhost:8080
Endpoint
POST /tax/calculations
Request ตัวอย่าง (Bonus Test Case)
JSON{
  "totalIncome": 850000,
  "wht": 0,
  "allowances": [
    { "allowanceType": "donation", "amount": 150000 }
  ]
}
Response
JSON{ "tax": 56000 }
Request ตัวอย่างพื้นฐาน
JSON{
  "totalIncome": 1200000,
  "wht": 0,
  "allowances": []
}
Response
JSON{ "tax": 138000 }

## คำนวณภาษีพื้นฐาน
Bashcurl -X POST http://localhost:8080/tax/calculations \
  -H "Content-Type: application/json" \
  -d '{"totalIncome":1200000,"wht":0,"allowances":[]}'
## คำนวณภาษีพร้อม donation (Bonus)
Bashcurl -X POST http://localhost:8080/tax/calculations \
  -H "Content-Type: application/json" \
  -d '{"totalIncome":850000,"wht":0,"allowances":[{"allowanceType":"donation","amount":150000}]}'
## คำนวณภาษีพร้อม WHT
Bashcurl -X POST http://localhost:8080/tax/calculations \
  -H "Content-Type: application/json" \
  -d '{"totalIncome":450000,"wht":8000,"allowances":[]}'
