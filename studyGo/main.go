package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

type student struct {
	gorm.Model
	Name  string `gorm:"column:name"`
	Age   int    `gorm:"column:age"`
	Grade string `gorm:"column:grade"`
}

func main() {
	dsn := "root:root@tcp(127.0.0.1:3306)/study_go?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err) // 或者返回错误给调用者
	}
	//db.AutoMigrate(&student{})
	//db.Create(&student{Name: "李斯", Age: 14, Grade: "三年级"})
	//var stu student
	//db.Where("age > ?", 18).Find(&stu)
	//fmt.Println(stu)
	//
	//db.Where("name = ?", "张三").Update("grade", "四年级")
	//
	//db.Where("age < ?", 15).Delete(&student{})

	result := map[string]interface{}{}
	db.Table("students").Find(&result)
	fmt.Println(result)

}
