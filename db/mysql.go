package db

import (
	"easywrite-service/model"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	Mysql *gorm.DB
)

type MysqlConfig struct {
	Username string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
}

func (c *MysqlConfig) toUrl() string {
	return fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.Username, c.Password, c.Host, c.Port, c.Database)
}

func InitMySql(mysqlConfig MysqlConfig) {
	var err error
	Mysql, err = gorm.Open(mysql.Open(mysqlConfig.toUrl()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}

	err = Mysql.AutoMigrate(&model.User{})
	if err != nil {
		panic(err)
	}
	err = Mysql.AutoMigrate(&model.Bill{})
	err = Mysql.AutoMigrate(&model.Feedback{})
	err = Mysql.AutoMigrate(&model.AppVersion{})
	if err != nil {
		panic(err)
	}
}
