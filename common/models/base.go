package models

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	Id        int            `gorm:"primary_key;column:id;serial;" json:"id" form:"id"`
	CreatedAt time.Time      `gorm:"column:create_time" json:"create_time" form:"create_time"`
	UpdatedAt time.Time      `gorm:"column:update_time" json:"update_time" form:"update_time"`
	DeletedAt gorm.DeletedAt `gorm:"column:delete_time" sql:"index" json:"-"`
}
