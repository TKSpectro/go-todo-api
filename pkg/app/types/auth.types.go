package types

type LoginDTOBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=100"`
}

type LoginDTO struct {
	Account LoginDTOBody `json:"account"`
}

type RegisterDTOBody struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6,max=100"`
	Firstname string `json:"firstname" validate:"omitempty,min=2"`
	Lastname  string `json:"lastname" validate:"omitempty,min=2"`
}

type RegisterDTO struct {
	Account RegisterDTOBody `json:"account"`
}

type AuthResponseBody struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}

type AuthResponse struct {
	Auth AuthResponseBody `json:"auth"`
}
