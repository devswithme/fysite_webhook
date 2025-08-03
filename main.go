package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	_ "embed"
	"tally_webhook/lms"
	"tally_webhook/mail"
	"tally_webhook/model"
	"tally_webhook/redis"
	"tally_webhook/tally"
	"tally_webhook/tripay"
)

func main() {
	godotenv.Load()

	http.HandleFunc("/webhook/tally", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
			return
		}

		header := r.Header.Get("X-Webhook-Secret")

		if header != os.Getenv("TALLY_SECRET") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		bodyBytes, err := io.ReadAll(r.Body)

		if err != nil {
			http.Error(w, "Failed to read body", http.StatusInternalServerError)
			return
		}

		defer r.Body.Close()

		data := tally.TallyResponse(bodyBytes)
		
		res := tripay.CreatePayment(data)
		mail.SendPaymentEmail(res)
	})

	http.HandleFunc("/webhook/tripay", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
			return
		}

		bodyBytes, err := io.ReadAll(r.Body)

		if err != nil {
			http.Error(w, "Failed to read body", http.StatusInternalServerError)
			return
		}

		defer r.Body.Close()

		header := r.Header.Get("X-Callback-Signature")
		signature := tripay.CreateWebhookSignature(bodyBytes)


		if header != signature {
			http.Error(w, "Invalid signature", http.StatusUnauthorized)
			return
		}

		var response model.WebhhokResponse
		json.Unmarshal(bodyBytes, &response)

		if response.Status == "PAID" {
			cache := redis.GetTransaction(response)
	
			lms.AddAccount(cache)
			for _, code := range cache.Codes {
				lms.AddGrantedStudent(code, cache.Email)
			}

			mail.SendNotifEmail(cache)
			redis.DeleteTransaction(response)
		}	

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]bool{
			"success": true,
		})
	
	})

	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
