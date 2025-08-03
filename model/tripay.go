package model

type OrderItem struct {
	Name  string `json:"name"`
	Quantity   int    `json:"quantity"`
	Price int    `json:"price"`
}

type Course struct {
	Code string `json:"code"`
	Price int `json:"price"`
}

type TripayPayload struct {
	Method string `json:"method"`
	MerchantRef string `json:"merchant_ref"`
	Amount int `json:"amount"`
	CustomerName string `json:"customer_name"`
	CustomerEmail string `json:"customer_email"`
	CustomerPhone string `json:"customer_phone"`
	OrderItems []OrderItem `json:"order_items"`
	Signature string `json:"signature"`
}

type TripayResponse struct {
	Data struct {
		CustomerName string `json:"customer_name"`
		CustomerEmail string `json:"customer_email"`
		Reference string `json:"reference"`
		PaymentName string `json:"payment_name"`
		Amount int `json:"amount"`
		AmountFormatted string
		ExpiredTime int `json:"expired_time"`
		CheckoutUrl string `json:"checkout_url"`
	} `json:"data"`
}

type PaymentCacheData struct {
		Reference string
		ExpiredTime int
}

type WebhhokResponse struct {
	Reference string `json:"reference"`
	Status string `json:"status"`
}