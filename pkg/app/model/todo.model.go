package model

import (
	"gopkg.in/guregu/null.v4"
	"gopkg.in/guregu/null.v4/zero"
)

type Todo struct {
	BaseModel
	Title       zero.String `gorm:"not null" json:"title" x-search:"true" swaggertype:"string" validate:"required,min=1"`
	Description zero.String `gorm:"" json:"description" x-search:"true" swaggertype:"string"`
	Completed   bool        `gorm:"default:false" json:"completed"`
	CompletedAt null.Time   `gorm:"" json:"completedAt" swaggertype:"string" format:"date-time"`

	AccountID uint `gorm:"not null" json:"fkAccountId"`
	// Account   Account
}

func (todo *Todo) New(remote Todo) {
	todo.Title = remote.Title
	todo.Description = remote.Description
	todo.Completed = remote.Completed
	todo.CompletedAt = remote.CompletedAt
}
