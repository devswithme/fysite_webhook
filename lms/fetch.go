package lms

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"tally_webhook/model"
)

func FetchGraphQL(payload model.GraphQLPayload, result interface{}) error {
	data, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, "http://localhost:3000/api/graphql", bytes.NewBuffer(data))
	req.Header.Set("Authorization", "Bearer "+os.Getenv("API_SECRET"))
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		return err
	}

	client := &http.Client{}

	res, err := client.Do(req)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)

	if err != nil {
		return err
	}

	if result != nil {
        if err := json.Unmarshal(bodyBytes, result); err != nil {
            return err
        }
    }

	return nil
}