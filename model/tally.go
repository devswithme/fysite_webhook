package model

type UserData struct {
	Price  string
	Email  string
	Name string
	Phone string
	Codes []string
	Channel string
	Amount int
}

type UserCacheData struct {
	Name string
	Email string 
	Codes []string
}

type TallyPayload struct {
	Data struct {
		Fields []struct {
			Label   string        `json:"label"`
			Value   interface{}   `json:"value"`
			Options []struct {
				ID   string `json:"id"`
				Text string `json:"text"`
			} `json:"options"`
		} `json:"fields"`
	} `json:"data"`
}
