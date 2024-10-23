package repository

import (
	"github.com/glebarez/sqlite"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"stash-go/model/entity"
)

type Repository struct {
	Db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		Db: db,
	}
}

func NewDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(""+
		"data.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("数据库连接失败： %s", err.Error())
		return nil
	}

	db.Set("gorm:table_options", "ENGINE=InnoDB")

	// 自动迁移
	err = db.AutoMigrate(
		&entity.Activity{},
		&entity.Package{},
	)
	if err != nil {
		log.Fatalf("数据库自动迁移失败： %s", err.Error())
		return nil
	}

	return db
}
