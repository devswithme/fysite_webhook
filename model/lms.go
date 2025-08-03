package model

type GraphQLPayload struct {
	Query string `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

type GetUsersResponse struct {
    Data struct {
        Users struct {
            TotalDocs int `json:"totalDocs"`
        } `json:"Users"`
    } `json:"data"`
}

type CreateUserResponse struct {
    Data struct {
        CreateUser struct {
            ID int `json:"id"`
        } `json:"createUser"`
    } `json:"data"`
}

type GetUserIdByEmailResponse struct {
    Data struct {
        Users struct {
            Docs []struct {
                ID int `json:"id"`
            } `json:"docs"`
        } `json:"Users"`
    } `json:"data"`
}


type GetGrantedStudentsResponse struct {
      Data struct {
        Courses struct {
            Docs []struct {
                ID int `json:"id"`
                GrantedTo []struct {
                    ID int `json:"id"`
                } `json:"grantedTo"`
            } `json:"docs"`
        } `json:"Courses"`
    } `json:"data"`
}
