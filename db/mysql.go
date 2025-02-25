package db

import (
	"easywrite-service/model"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
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
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Username, c.Password, c.Host, c.Port, c.Database)
}

// 自动重连函数
func connectWithRetry(mysqlConfig MysqlConfig, maxRetries int, retryInterval time.Duration) *gorm.DB {
	var db *gorm.DB
	var err error

	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(mysql.Open(mysqlConfig.toUrl()), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
		if err == nil {
			log.Println("MySQL 连接成功！")
			return db
		}

		log.Printf("连接 MySQL 失败: %v，%d 秒后重试 (%d/%d)...\n", err, retryInterval/time.Second, i+1, maxRetries)
		time.Sleep(retryInterval)
	}

	panic(fmt.Sprintf("无法连接到 MySQL，已重试 %d 次: %v", maxRetries, err))
}

func InitMySql(mysqlConfig MysqlConfig) {
	// 允许最多 5 次重试，每次间隔 5 秒
	Mysql = connectWithRetry(mysqlConfig, 5, 5*time.Second)

	// 自动迁移数据库表结构
	if err := Mysql.AutoMigrate(&model.User{}, &model.Bill{}, &model.Feedback{}, &model.AppVersion{}); err != nil {
		panic(fmt.Sprintf("数据库迁移失败: %v", err))
	}

	// 定时检查数据库连接
	go keepAlive(30*time.Second, mysqlConfig)
}

// 定时检查数据库连接并重连
func keepAlive(interval time.Duration, mysqlConfig MysqlConfig) {
	for {
		time.Sleep(interval)

		sqlDB, err := Mysql.DB()
		if err != nil {
			log.Println("获取数据库连接失败:", err)
			continue
		}

		if err = sqlDB.Ping(); err != nil {
			log.Println("数据库连接丢失，尝试重连...")
			Mysql = connectWithRetry(mysqlConfig, 5, 5*time.Second)
		}
	}
}
