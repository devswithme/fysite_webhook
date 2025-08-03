package redis

import (
	"encoding/json"
	"log"
	"tally_webhook/model"
	"time"
)

func SaveTransaction(user model.UserCacheData, payment model.PaymentCacheData){
	data, err := json.Marshal(user)

	if err != nil {
		log.Fatal(err)
	}

	err = Client.Set(Ctx, payment.Reference, data, time.Duration(payment.ExpiredTime - int(time.Now().Unix())) * time.Second).Err()

	if err != nil {
		log.Fatal(err)
	}
}

func GetTransaction(response model.WebhhokResponse) model.UserCacheData {
	val, err := Client.Get(Ctx, response.Reference).Result()

	if err != nil {
		log.Fatal(err)
	}

	var payload model.UserCacheData
	json.Unmarshal([]byte(val), &payload)

	return payload
}

func DeleteTransaction(response model.WebhhokResponse){
	err := Client.Del(Ctx, response.Reference).Err()

	if err != nil {
		log.Fatal(err)
	}
}