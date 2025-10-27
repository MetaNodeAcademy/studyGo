package model

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Username string `gorm:"not null"`
	Password string
	Email    string
}

type Post struct {
	gorm.Model
	Title   string `gorm:"not null"`
	Content string `gorm:"not null"`
	UserID  uint
	User    User
}

type Comment struct {
	gorm.Model
	Content string `gorm:"not null"`
	PostID  uint
	Post    Post
	UserID  uint
	User    User
}

var dm *DBManager

func init() {
	dm = GetDBManager()
}

// 注册方法
func Register(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = string(hashedPassword)

	if err := dm.DB.Create(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(201, gin.H{"message": "User registered successfully"})
}

// 登录方法
func Login(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	var storedUser User
	if err := dm.DB.Where("username = ?", user.Username).First(&storedUser).Error; err != nil {
		c.JSON(401, gin.H{"error": "Invalid username or password"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password)); err != nil {
		c.JSON(401, gin.H{"error": "Invalid username or password"})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       storedUser.ID,
		"username": storedUser.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate token"})
		return
	}
	c.SetCookie("token", tokenString, 3600*24, "/", "", false, true)
}

// 创建文章
func CreatePost(c *gin.Context) {
	var post Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := dm.DB.Create(&post).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create post"})
		return
	}
	token, err := c.Cookie("token")
	if err != nil {
		c.JSON(401, gin.H{"error": "Please login first"})
	}
	claims, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if claimsMap, ok := claims.Claims.(jwt.MapClaims); ok {
		userID := uint(claimsMap["id"].(float64))
		post.UserID = userID
		if err := dm.DB.Save(&post).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to save post"})
			return
		}
	} else {
		c.JSON(401, gin.H{"error": "Invalid token"})
		return
	}
}

// 获取所有文章列表
func GetAllPosts(c *gin.Context) {
	var posts []Post
	if err := dm.DB.Preload("User").Find(&posts).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to get posts"})
	} else {
		c.JSON(200, posts)
	}
	return
}
func GetPostByID(c *gin.Context) {
	id := c.Param("id")
	var post Post
	if err := dm.DB.Preload("User").First(&post, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Post not found"})
		return
	}
	c.JSON(200, post)
}

// 更新修改所有文章列表
func UpdatePost(c *gin.Context) {
	var post Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}
	//应该验证post里面的用户id和token是否一致，放在上一层校验
	if err := dm.DB.Save(&post).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update post"})
	} else {
		c.JSON(200, post)
	}
}
func DeletePost(c *gin.Context) {
	id := c.Param("id")
	var post Post
	if err := dm.DB.First(&post, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Post not found"})
		return
	}
	if err := dm.DB.Delete(&post).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete post"})
		return
	}
}

// 评论文章
func CreateComment(c *gin.Context) {
	var comment Comment
	var post Post
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if err := dm.DB.First(&post, comment.PostID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Post not found"})
		return
	}
	if err := dm.DB.Create(&comment).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create comment"})
	}
	c.JSON(201, gin.H{"message": "Comment created successfully"})
}

// 根据文章获取评论
func GetCommentList(c *gin.Context) {
	var comments []Comment
	var postId int
	if err := c.ShouldBindQuery(&postId); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
	}
	if err := dm.DB.Where("post_id = ?", postId).Find(&comments).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to get comments"})
	} else {
		c.JSON(200, comments)
	}
}
