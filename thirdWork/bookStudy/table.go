package test

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

type User struct {
	gorm.Model
	UserName string `gorm:"column:username"`
	Email    string `gorm:"column:email"`
	PostNum  uint
	Posts    []Post
}
type Post struct {
	gorm.Model
	Title      string `gorm:"column:title"`
	Content    string `gorm:"column:content"`
	CommentNum uint
	UserID     uint
	User       User
	Comments   []Comment
}
type Comment struct {
	gorm.Model
	Content string `gorm:"column:content"`
	PostID  uint   `gorm:"not null"`
	Post    Post
}

func createTable() {
	dsn := "root:root@tcp(127.0.0.1:3306)/study_go?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err) // 或者返回错误给调用者
	}
	err = db.AutoMigrate(&User{}, &Post{}, &Comment{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("数据库创建更新成功")
}
