package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Job struct {
	gorm.Model
	Name        string
	JobURL      string
	Description string
	Company     Company
}

type Company struct {
	gorm.Model
	Name        string
	Industry    string
	CompanyURL  string
	Description string
	JobID       uint
}

func connect() *gorm.DB {
	db, err := gorm.Open(postgres.Open("host=db user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Taipei"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	// 檢查 Job 和 Company 表是否存在
	if db.Migrator().HasTable(&Job{}) {
		// 如果存在，則刪除
		db.Migrator().DropTable(&Job{})
	}
	if db.Migrator().HasTable(&Company{}) {
		// 如果存在，則刪除
		db.Migrator().DropTable(&Company{})
	}
	// Migrate the schema
	db.AutoMigrate(&Job{}, &Company{})
	return db
}
