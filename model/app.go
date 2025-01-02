package model

type AppInfo struct {
	AppID          string `json:"app_id,omitempty"`
	AppVersionCode int64  `json:"app_version_code,omitempty"`
	AppVersionName string `json:"app_version_name,omitempty"`
}
