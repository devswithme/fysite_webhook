package lms

import (
	"log"
	"tally_webhook/model"
)

func contains(slice []int, val int) bool {
    for _, item := range slice {
        if item == val {
            return true
        }
    }
    return false
}


func getGrantedStudents(code string) (int, []int){
	payload := model.GraphQLPayload {
		Variables: map[string]interface{}{
			 "code": code,
		},
		Query: `query GetCourses($code: String!){
			Courses(where: {code: {equals: $code}}){
				docs {
					id
					grantedTo {
						id
					}
				}
			}
		}`,		
	}

	var response model.GetGrantedStudentsResponse
	
	if err := FetchGraphQL(payload, &response); err != nil {
		log.Fatalf("FetchGraphQL failed: %v", err)
	}

	ids := make([]int, 0, len(response.Data.Courses.Docs[0].GrantedTo))

	for _, granted := range response.Data.Courses.Docs[0].GrantedTo {
		ids = append(ids, granted.ID)
	}


	return response.Data.Courses.Docs[0].ID, ids
}

func getUserIdByEmail(email string) int {
	payload := model.GraphQLPayload {
		Variables: map[string]interface{}{
			"email": email,
		},
		Query: `query GetUserIdByEmail($email: EmailAddress!){
			Users(where: {email: {equals: $email}}){
				docs {
					id
				}
			}
		}`,
	}

	var response model.GetUserIdByEmailResponse

	if err := FetchGraphQL(payload, &response); err != nil {
		log.Fatalf("FetchGraphQL failed: %v", err)
	}

	return response.Data.Users.Docs[0].ID
}

func AddGrantedStudent(code string, email string) {
	userId := getUserIdByEmail(email)

	
	courseId, grantedIds := getGrantedStudents(code)
	
	if contains(grantedIds, userId){
		return
	}

	payload := model.GraphQLPayload {
		Variables: map[string]interface{}{
			"id": courseId,
			"data": map[string]interface{}{
				"grantedTo": append(grantedIds, userId),
			},
		},
		Query: `mutation updateCourse($id: Int!, $data: mutationCourseUpdateInput!){
			updateCourse(id: $id, data: $data){
				id
			}
		}`,
	}

	if err := FetchGraphQL(payload, nil); err != nil {
		log.Fatalf("FetchGraphQL failed: %v", err)
	}
}
