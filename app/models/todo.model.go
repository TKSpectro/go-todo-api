package models

import (
	"tkspectro/vefeast/app/types/pagination"
	"tkspectro/vefeast/config/database"
	"tkspectro/vefeast/utils"

	"gorm.io/gorm"
)

type Todo struct {
	BaseModel
	Title       string `gorm:"not null" json:"title"`
	Description string `gorm:"" json:"description"`

	AccountID int `gorm:"not null" json:"fkAccountId"`
	// Account   Account
}

func FindTodosByAccount(dest interface{}, meta *pagination.Meta, accountID uint) *gorm.DB {
	query := database.DB.Model(&Todo{}).Where("account_id = ?", accountID)

	utils.CountMeta(meta, query)

	order := utils.ParseOrder(meta.Order)

	return query.
		Offset(meta.Offset).
		Limit(meta.Limit).
		Order(order).
		Find(dest)
}

func FindTodoByID(dest interface{}, id string, accountID uint) *gorm.DB {
	return database.DB.Model(&Todo{}).Where("id = ? AND account_id = ?", id, accountID).Take(dest)
}
