package model

import (
	"easywrite-service/mytime"
)

type Bill struct {
	ID                 int64                 `json:"id" gorm:"column:id;autoIncrement:false;type:bigint"`
	ThirdPartyID       string                `json:"third_party_id" gorm:"column:third_party_id;type:varchar(2048);index;unique"`
	CreatedAt          mytime.CustomTime     `json:"created_at"`
	UpdatedAt          mytime.CustomTime     `json:"updated_at"`
	DeletedAt          mytime.CustomNullTime `json:"deleted_at" gorm:"index"`
	Username           string                `json:"username" gorm:"column:username;type:varchar(255)"`
	Amount             int64                 `json:"amount" gorm:"column:amount;type:bigint"`
	Comment            string                `json:"comment" gorm:"column:comment;varchar(1024)"`
	Datetime           mytime.CustomTime     `json:"datetime" gorm:"column:datetime;type:datetime"`
	Date               mytime.CustomTime     `json:"date" gorm:"column:date;type:date"`
	Category           string                `json:"category" gorm:"column:category;varchar(255)"`
	TransactionPartner string                `json:"transaction_partner" gorm:"column:transaction_partner;varchar(255)"`
	Name               string                `json:"name" gorm:"column:name;varchar(255)"`
	Type               string                `json:"type" gorm:"column:type;varchar(255)"`
	ImagesComment      StringList            `json:"images_comment" gorm:"column:images_comment;text"`
}

func (m *Bill) TableName() string {
	return "bills"
}

type GetBillParameters struct {
	Limit  int `form:"limit"`
	Offset int `form:"offset"`
}

type AddBillResponse struct {
	ID            int64      `json:"id"`
	ThirdPartyId  string     `json:"third_party_id"`
	ImagesComment StringList `json:"images_comment"`
}
