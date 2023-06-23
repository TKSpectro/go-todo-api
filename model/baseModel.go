package model

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint           `gorm:"primarykey"`
	CreatedAt time.Time      `gorm:"autoCreateTime:milli"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime:milli"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
