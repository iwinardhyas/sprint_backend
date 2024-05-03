package models

// SignUp struct to describe register a new user.
type SignUp struct {
	Email    string `json:"email" validate:"required,email,lte=255"`
	Name     string `json:"name" validate:"required,min=5,max=50"`
	Password string `json:"password" validate:"required,min=5,max=15,lte=255"`
}

// SignIn struct to describe login user.
type SignIn struct {
	Email    string `json:"email" validate:"required,email,lte=255"`
	Password string `json:"password" validate:"required,min=5,max=15,lte=255"`
}
