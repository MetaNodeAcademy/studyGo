package main

import (
	"github.com/gin-gonic/gin"
	"work4/authUtil"
	"work4/model"
)

func main() {
	err := model.InitGlobalDB()
	if err != nil {
		panic(err)
	}
	createErr := model.GetDBManager().CreateTable()
	if createErr != nil {
		panic(createErr)
	}
	router := gin.Default()
	router.Use(authUtil.ErrorHandler())
	router.POST("/register", model.Register)
	router.POST("/login", model.Login)
	post := router.Group("/post")
	post.Use(authUtil.JWTAuth())
	post.POST("/create", model.CreatePost)
	post.POST("/list", model.GetAllPosts)
	post.POST("/detail", model.GetPostByID)
	post.POST("/update", model.UpdatePost)
	post.POST("/delete", model.DeletePost)
	comment := router.Group("/comment")
	comment.Use(authUtil.JWTAuth())
	comment.POST("/create", model.CreateComment)
	comment.POST("/list", model.GetCommentList)
	serverletErr := router.Run(":8080")
	if serverletErr != nil {
		return
	}
}
