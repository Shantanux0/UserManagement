package models

type CreateUserRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
	DOB  string `json:"dob" validate:"required,datetime=2006-01-02,past_date"`
}

type UpdateUserRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
	DOB  string `json:"dob" validate:"required,datetime=2006-01-02,past_date"`
}

type UserResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	DOB  string `json:"dob"`
}

type UserWithAgeResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	DOB  string `json:"dob"`
	Age  int    `json:"age"`
}
