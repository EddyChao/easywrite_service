package common

var UploadSavePath string
var UploadDir = "/upload"

func InitUploadSavePath(dir string) {
	UploadSavePath = dir
}
