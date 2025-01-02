package model

import (
	"easywrite-service/mytime"
)

type AppVersion struct {
	ID          uint                   `gorm:"primary_key;auto_increment" json:"-"`
	Appid       string                 `gorm:"type:varchar(255); not null" json:"appid"`
	AppName     string                 `gorm:"type:varchar(255); not null" json:"app_name"`
	VersionCode int                    `gorm:"type:int; not null" json:"version_code"`
	VersionName string                 `gorm:"type:varchar(255); not null" json:"version_name"`
	IsForce     *bool                  `gorm:"type:tinyint(1); default:false; not null" json:"is_force"`
	CanHide     *bool                  `gorm:"type:tinyint(1); default:true; not null" json:"can_hide"`
	ChangeLog   string                 `gorm:"type:varchar(255); not null" json:"changelog"`
	DownloadUrl string                 `gorm:"type:varchar(1024)" json:"download_url"`
	CreatedAt   *mytime.CustomTime     `gorm:"type:datetime; not null" json:"created_at"`
	UpdatedAt   *mytime.CustomTime     `gorm:"type:datetime; not null" json:"updated_at"`
	DeletedAt   *mytime.CustomNullTime `gorm:"type:datetime; default:null" json:"-"`
}
