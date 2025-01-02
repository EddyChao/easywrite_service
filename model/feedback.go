package model

import "easywrite-service/mytime"

type Feedback struct {
	ID        uint                  `json:"id" gorm:"primaryKey"`
	Label     string                `form:"label"`
	Feedback  string                `form:"feedback"`
	Picture   string                `form:"-"`
	Phone     string                `form:"phone"`
	CreatedAt mytime.CustomTime     `json:"created_at"`
	UpdatedAt mytime.CustomTime     `json:"updated_at"`
	DeletedAt mytime.CustomNullTime `json:"deleted_at" gorm:"index"`
}
