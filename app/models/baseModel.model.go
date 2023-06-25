package models

import (
	"time"

	"github.com/TKSpectro/go-todo-api/app/types/pagination"
	"github.com/TKSpectro/go-todo-api/config/database"
	"github.com/TKSpectro/go-todo-api/utils"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `gorm:"autoCreateTime:milli" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:milli" json:"updatedAt"`
}

func FindWithMeta(dest interface{}, model interface{}, meta *pagination.Meta, where *gorm.DB) *gorm.DB {
	search, searchArgs := utils.SearchWhere(meta.Search, model)

	query := database.DB.Model(model).Where(search, searchArgs...)

	if where != nil {
		query = query.Where(where)
	}

	utils.CountMeta(meta, query)

	return query.
		Offset(meta.Offset).
		Limit(meta.Limit).
		Order(meta.Order).
		Find(dest)
}
