package models

// CreateUserRequest represents the body payload when creating a user.
type CreateUserRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
	DOB  string `json:"dob" validate:"required,datetime=2006-01-02,past_date"`
}

// UpdateUserRequest represents the body payload when updating a user.
type UpdateUserRequest struct {
	Name string `json:"name" validate:"required,min=1,max=100"`
	DOB  string `json:"dob" validate:"required,datetime=2006-01-02,past_date"`
}

// UserResponse represents the JSON response for create/update operations.
type UserResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	DOB  string `json:"dob"`
}

// UserWithAgeResponse represents the JSON response for retrieve operations, including the calculated age.
type UserWithAgeResponse struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	DOB  string `json:"dob"`
	Age  int    `json:"age"`
}
