package model

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/wzru/gitran-server/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() error {
	var err error
	if config.DB.Type == "mysql" {
		dsn := fmt.Sprintf("%s:%s@tcp(%s)/?charset=utf8mb4&parseTime=True&loc=Local", config.DB.User, config.DB.Password, config.DB.Host)
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("Database connect ERROR : %v", err)
			return err
		}
		DB.Exec("CREATE DATABASE IF NOT EXISTS " + config.DB.Name)
		dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.DB.User, config.DB.Password, config.DB.Host, config.DB.Name)
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	} else if config.DB.Type == "postgresql" {
		//TODO
	}
	if err != nil {
		log.Fatalf("Database connect ERROR : %v", err)
		return err
	}
	DB.AutoMigrate(&User{}, &Project{}, &Phrase{}, &Translation{})
	return nil
}