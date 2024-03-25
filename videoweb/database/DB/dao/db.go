package dao

import (
	"context"
	"strings"
	"videoweb/config"
	"videoweb/database/DB/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	Db *gorm.DB
)

func InitDB() {
	//config.Init()

	dsn := strings.Join([]string{config.SqlUserName, ":", config.SqlPassword,
		"@tcp(" + config.MysqlIP + ")/", config.DataBase, //"")
		"?charset=utf8mb4&parseTime=True&loc=Local"}, "")

	var err error
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	Migrate()
}

func Migrate() {
	Db.AutoMigrate(&model.User{})
	Db.AutoMigrate(&model.Video{})
	Db.AutoMigrate(&model.Relation{})
	Db.AutoMigrate(&model.Comment{})
	Db.AutoMigrate(&model.Like{})
	Db.AutoMigrate(&model.Message{})
}

func NewDBClient(ctx context.Context) *gorm.DB {
	db := Db
	return db.WithContext(ctx)
}
