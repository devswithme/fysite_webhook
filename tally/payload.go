package tally

import (
	"bytes"
	"encoding/json"
	"tally_webhook/model"
)

func TallyResponse(reqBytes []byte) model.UserData {
	var payload model.TallyPayload

	json.NewDecoder(bytes.NewBuffer(reqBytes)).Decode(&payload)

	var data model.UserData

	for _, field := range payload.Data.Fields {
		switch field.Label {
		case "Course code":
			ids, ok := field.Value.([]interface{})
		if ok {
			for _, idRaw := range ids {
				idStr := idRaw.(string)
				for _, option := range field.Options {
					if option.ID == idStr {
						data.Codes = append(data.Codes, option.Text)
					}
				}
			}
		}
		case "Channel":
			if channels, ok := field.Value.([]interface{}); ok {
				channel, _ := channels[0].(string)

				for _, option := range field.Options {
					if option.ID == channel {
						data.Channel = option.Text
						break
					}
				}
			}
		case "Amount":
			if f, ok := field.Value.(float64); ok {
				data.Amount = int(f)
			}
		case "Email":
			if email, ok := field.Value.(string); ok {
				data.Email = email
			}
		case "Name":
			if name, ok := field.Value.(string); ok {
				data.Name = name
			}
		case "Phone":
			if phone, ok := field.Value.(string); ok {
				data.Phone = phone
			}
		}
	}

	return data
}