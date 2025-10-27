package model

import (
	"database/sql"
	"errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"sync"
	"time"
)

type DBManager struct {
	DB    *gorm.DB
	SqlDB *sql.DB
}

// 单例模式
var Ins *DBManager
var once sync.Once

func InitGlobalDB() error {
	var err error
	once.Do(func() {
		Ins = &DBManager{}
		err = Ins.Init()
	})
	return err
}

func GetDBManager() *DBManager {
	if Ins == nil {
		log.Fatal("实例暂未初始化")
	}
	return Ins
}

func (dm *DBManager) Init() error {
	username := "root"
	password := "root"
	host := "127.0.0.1"
	port := "3306"
	dbName := "study_go"
	dsn := username + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: dsn,
	}))
	if err != nil {
		return err
	}
	dm.DB = db
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	dm.SqlDB = sqlDB
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(time.Minute * 20)
	return nil
}

func (dm *DBManager) CreateTable() error {
	err := dm.DB.AutoMigrate(&User{}, &Post{}, &Comment{})
	if err != nil {
		return errors.New(err.Error())
	}

	return nil
}
