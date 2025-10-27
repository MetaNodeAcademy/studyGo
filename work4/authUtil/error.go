package authUtil

import (
	"github.com/gin-gonic/gin"
	"log"
	error2 "work4/error"
)

func ErrorHandler() gin.HandlerFunc {
	return func(context *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				context.JSON(500, gin.H{
					"code":    500,
					"message": "服务器错误",
				})
				context.Abort()
			}
		}()
		context.Next()
		if len(context.Errors) > 0 {
			for _, err := range context.Errors {
				log.Printf("error: %v", err)
			}
			firstError := context.Errors[0].Err
			appErr, ok := firstError.(*error2.ResultError)
			if ok {
				context.JSON(appErr.Code, gin.H{
					"code":    appErr.Code,
					"message": appErr.Message,
				})
			} else {
				context.JSON(500, gin.H{
					"code":    500,
					"message": "服务器错误",
				})
			}
		}
	}

}
