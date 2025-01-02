package model

import "easywrite-service/mytime"

type User struct {
	ID int64 `gorm:"id" json:"id"`
	UserInfoReadOnly
	Password  string                `gorm:"password"`
	Salt      string                `gorm:"salt"`
	CreatedAt mytime.CustomTime     `json:"-"`
	UpdatedAt mytime.CustomTime     `json:"-"`
	DeletedAt mytime.CustomNullTime `json:"-" gorm:"index"`
}

type UserInfo struct {
	Nickname      string `json:"nickname" gorm:"nickname"`
	Avatar        string `json:"avatar" gorm:"avatar"`
	SchoolID      int    `json:"schoolID" gorm:"school_id"`
	CampusID      string `json:"campusID" gorm:"campus_id"`
	School        string `json:"school"`
	Campus        string `json:"campus"`
	CampusAddress string `json:"campusAddress"`
}

type UserInfoReadOnly struct {
	Email      string `gorm:"email" json:"email"`
	Phone      string `gorm:"phone" json:"phone"`
	Username   string `gorm:"column:username;unique" json:"username"`
	IsHorseman *bool  `gorm:"is_horseman" json:"isHorseman"`
	UserInfo
}

type UseCodeLoginParameters struct {
	Username         string `json:"username" binding:"required"`
	VerificationCode string `json:"code" binding:"required"`
}

type LoginParameters struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegistrationParameters struct {
	Username         string `json:"username" binding:"required"`
	Password         string `json:"password" binding:"required"`
	VerificationCode string `json:"code" binding:"required"`
}

type VerificationParameters struct {
	PhoneOrEmail string `binding:"required" json:"phoneOrEmail"`
	CodeType     string `binding:"required" json:"codeType"`
}

type ResetPasswordParameters struct {
	Username         string `json:"username"  binding:"required"`
	NewPassword      string `json:"newPassword"  binding:"required"`
	VerificationCode string `json:"code"  binding:"required"`
}

type ChangePasswordParameters struct {
	Username    string `json:"username,omitempty"`
	OldPassword string `json:"oldPassword,omitempty"`
	NewPassword string `json:"newPassword,omitempty"`
}
