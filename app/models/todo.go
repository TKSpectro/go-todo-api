package models

type Todo struct {
	BaseModel
	Title       string `gorm:"not null" json:"title"`
	Description string `gorm:"" json:"description"`

	AccountID int
	Account   Account
}
