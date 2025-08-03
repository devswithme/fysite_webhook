package lms

import (
	"log"
	"tally_webhook/model"

	"github.com/google/uuid"
)



func isRegistered(email string) bool {
	payload := model.GraphQLPayload {
		Variables: map[string]interface{}{
			"email": email,
		},
		Query: `query GetUsers($email: EmailAddress!){
			Users(where: {email: {equals: $email}}){
				totalDocs
			}
		}`,
	}

	var response model.GetUsersResponse
	if err := FetchGraphQL(payload, &response); err != nil {
		log.Fatalf("FetchGraphQL failed: %v", err)
	}

	return response.Data.Users.TotalDocs > 0	
}

func AddAccount(data model.UserCacheData){
	if is_registred := isRegistered(data.Email); is_registred {
		return
	}

	password := uuid.New().String()

	payload := model.GraphQLPayload {
		Variables: map[string]interface{}{
			"data": map[string]interface{}{
				"email":    data.Email,
				"role":     "student",
				"name":     data.Name,
				"password": password,
			},
		},
		Query: `mutation CreateUser($data: mutationUserInput!){
			createUser(data: $data){
				id
			}
		}`,
	}

	if err := FetchGraphQL(payload, nil); err != nil {
		log.Fatalf("FetchGraphQL failed: %v", err)
	}
}
