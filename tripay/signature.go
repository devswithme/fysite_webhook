package tripay

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"strconv"
)


func CreatePaymentSignature(amount int, merchantRef string) string {
	privateKey := os.Getenv("TRIPAY_PRIVATE_KEY")
	merchantCode := os.Getenv("TRIPAY_MERCHANT_CODE")
	
	h := hmac.New(sha256.New, []byte(privateKey))
	h.Write([]byte(merchantCode+merchantRef+strconv.Itoa(amount)))
	
	signature := hex.EncodeToString(h.Sum(nil))
	
	return signature
}

func CreateWebhookSignature(body []byte) string {
	privateKey := os.Getenv("TRIPAY_PRIVATE_KEY")

	h := hmac.New(sha256.New, []byte(privateKey))
	h.Write(body)

	signature := hex.EncodeToString(h.Sum(nil))

	return signature
}