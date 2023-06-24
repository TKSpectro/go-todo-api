package types

type LoginDTO struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=100"`
}

type RegisterDTO struct {
	LoginDTO
	Firstname string `json:"firstname" validate:"omitempty,min=2"`
	Lastname  string `json:"lastname" validate:"omitempty,min=2"`
}

type AuthResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refreshToken"`
}
