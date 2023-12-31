package infra

import (
	"github.com/Daka-0424/my-go-server/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type MySQLConnector struct {
	DB *gorm.DB
}

func NewMySQLConnector(cfg *config.Config) *MySQLConnector {
	conn := cfg.MySQL.DBConn

	log := logger.Default
	if cfg.IsDevelopment() {
		log = logger.Default.LogMode(logger.Info)
	}
	db, err := gorm.Open(mysql.Open(conn), &gorm.Config{
		Logger: log,
	})
	if err != nil {
		panic(err)
	}

	return &MySQLConnector{
		DB: db,
	}
}
