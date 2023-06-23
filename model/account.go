package model

type Account struct {
	BaseModel
	Email       string `gorm:"uniqueIndex;not null" json:"email"`
	Password    string `gorm:"not null" json:"password"`
	Firstname   string `gorm:"" json:"firstname"`
	Lastname    string `gorm:"" json:"lastname"`
	SecretToken string `gorm:"type:varchar(8)" json:"secretToken"`

	Todos []Todo
}
