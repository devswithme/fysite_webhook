package tripay

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"tally_webhook/model"
	"tally_webhook/redis"

	_ "embed"

	"github.com/google/uuid"
)

//go:embed course.json
var courseJSON []byte

func CreatePayment(data model.UserData) model.TripayResponse {
	merchant_ref := uuid.New().String()
	channel := data.Channel

	var courses []model.Course

	if err := json.Unmarshal(courseJSON, &courses); err != nil {
		log.Fatalf("Failed to parse course JSON: %v", err)
	}

	var items []model.OrderItem

	for _, code := range data.Codes {
		for _, course := range courses {
			if course.Code == code {
				items = append(items, model.OrderItem {
					Name:  course.Code,
					Quantity:   1,
					Price: course.Price,
				})
				break;
			}
		}
	}

	payload := model.TripayPayload{
		Method: channel,
		MerchantRef: merchant_ref,
		Amount: data.Amount,
		CustomerName: data.Name,
		CustomerPhone: data.Phone,
		CustomerEmail: data.Email,
		OrderItems: items,
		Signature: CreatePaymentSignature(data.Amount, merchant_ref),
	}

	request, err := json.Marshal(payload)

	if err != nil {
		log.Fatalf("Failed to stringify payload: %v", err)
	}

	
	req, err := http.NewRequest(http.MethodPost, "https://tripay.co.id/api-sandbox/transaction/create", bytes.NewBuffer(request))
	
	if err != nil {
		log.Fatalf("Failed to create HTTP request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("TRIPAY_API_KEY"))
	
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalf("HTTP request failed: %v", err)
	}

	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatalf("Failed to read response body: %v", err)
	}
	
	var tripayResponse model.TripayResponse
	
	if err := json.Unmarshal(body, &tripayResponse); err != nil {
		log.Fatalf("Failed to parse JSON response: %v", err)
	}

 	user := model.UserCacheData {
		Name: data.Name,
		Email: data.Email,
		Codes: data.Codes,
	}

	payment := model.PaymentCacheData {
		Reference: tripayResponse.Data.Reference,
		ExpiredTime: tripayResponse.Data.ExpiredTime,
	}

	redis.SaveTransaction(user, payment)

	return tripayResponse
}


