package mail

import (
	"bytes"
	"html/template"
	"log"
	"net/smtp"
	"os"
	"tally_webhook/model"

	_ "embed"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

//go:embed payment.html
var paymentHTML string

//go:embed notification.html
var notificationHTML string

func sendEmail(body bytes.Buffer, email string){
	smtpHost := os.Getenv("SMTP_HOST")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")
	smtpPort := os.Getenv("SMTP_PORT")

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	message := []byte("Subject: Your Order\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
		"\r\n" +
		body.String(),
	)

	smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpUser, []string{email}, message)
}

func SendPaymentEmail(res model.TripayResponse){
	tmpl, err := template.New("payment").Parse(paymentHTML)

	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	var body bytes.Buffer

	res.Data.AmountFormatted = message.NewPrinter(language.Indonesian).Sprintf("Rp %d", res.Data.Amount)

	tmpl.Execute(&body, res.Data)

	sendEmail(body, res.Data.CustomerEmail)
}

func SendNotifEmail(res model.UserCacheData){
	tmpl, err := template.New("notification").Parse(notificationHTML)

	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	var body bytes.Buffer
	tmpl.Execute(&body, res)

	sendEmail(body, res.Email)
}