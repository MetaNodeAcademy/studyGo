package test

import (
	"fmt"
	"log"
	_ "math"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func saveTable() {
	db := getDB()
	user := User{
		UserName: "zhangsan",
		Email:    "907560@qq.com",
	}
	saveUser(db, &user)
	// 同时新增文章和评论
	post := Post{
		Title:   "测试文章",
		Content: "测试内容",
		UserID:  1,
		Comments: []Comment{
			{Content: "牛逼"},
			{Content: "123"},
			{Content: "456"},
		},
	}
	savePost(db, &post)
}

func getDB() *gorm.DB {
	// 连接数据库
	dsn := "root:root@tcp(localhost:3306)/study_go?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	return db
}
func searchTable(userId uint) {
	db := getDB()
	// 1. 查询指定用户（ID=1）的所有文章及对应的评论
	userPosts, err := getUserArticlesWithComments(db, userId)
	if err != nil {
		log.Printf("查询用户文章及评论失败: %v", err)
	} else {
		printUserArticles(userPosts)
	}

	// 2. 查询评论数量最多的文章
	topPost, err := getPostWithMostComments(db)
	if err != nil {
		log.Printf("查询评论最多的文章失败: %v", err)
	} else {
		printTopCommentedPost(topPost)
	}
}

// getUserArticlesWithComments 查询指定用户的所有文章及评论
func getUserArticlesWithComments(db *gorm.DB, userID uint) (User, error) {
	var user User
	// 1. 先查询用户，通过 Preload 预加载关联的 Posts
	// 2. 对 Posts 再通过 Preload 预加载关联的 Comments
	result := db.Preload("Posts.Comments").First(&user, userID)
	if result.Error != nil {
		return User{}, fmt.Errorf("查询用户失败: %w", result.Error)
	}
	return user, nil
}

// getPostWithMostComments 查询评论数量最多的文章
func getPostWithMostComments(db *gorm.DB) (Post, error) {
	var post Post
	// 子查询：统计每篇文章的评论数，按评论数降序取第一条
	subQuery := db.Model(&Comment{}).Select("post_id, COUNT(*) as comment_count").Group("post_id").Order("comment_count DESC").Limit(1)

	// 主查询：关联子查询结果，获取评论最多的文章详情
	result := db.Joins("JOIN (?) as sub ON posts.id = sub.post_id", subQuery).Preload("Comments").First(&post)
	if result.Error != nil {
		return Post{}, fmt.Errorf("查询文章失败: %w", result.Error)
	}
	return post, nil
}

// 打印用户的文章及评论（辅助函数）
func printUserArticles(user User) {
	fmt.Printf("\n----- 用户 %s（ID: %d）的文章及评论 -----\n", user.UserName, user.ID)
	if len(user.Posts) == 0 {
		fmt.Println("该用户暂无文章")
		return
	}
	for _, post := range user.Posts {
		fmt.Printf("\n文章 ID: %d, 标题: %s\n", post.ID, post.Title)
		fmt.Printf("文章内容: %s\n", post.Content)
		fmt.Printf("评论数: %d\n", len(post.Comments))
		if len(post.Comments) > 0 {
			fmt.Println("评论列表:")
			for _, comment := range post.Comments {
				fmt.Printf("  - 评论 ID: %d, 内容: %s (发布时间: %s)\n", comment.ID, comment.Content, comment.CreatedAt.Format("2006-01-02 15:04:05"))
			}
		}
	}
}

// 打印评论最多的文章（辅助函数）
func printTopCommentedPost(post Post) {
	fmt.Printf("\n----- 评论最多的文章 -----\n")
	fmt.Printf("文章 ID: %d\n", post.ID)
	fmt.Printf("标题: %s\n", post.Title)
	fmt.Printf("内容: %s\n", post.Content)
	fmt.Printf("评论总数: %d\n", len(post.Comments))
	if len(post.Comments) > 0 {
		fmt.Println("最新评论:")
		// 取最后一条评论（假设按创建时间排序）
		lastComment := post.Comments[len(post.Comments)-1]
		fmt.Printf("  - %s (发布时间: %s)\n", lastComment.Content, lastComment.CreatedAt.Format("2006-01-02 15:04:05"))
	}
}

func savePost(db *gorm.DB, post *Post) {
	db.Create(post)
}

func saveUser(db *gorm.DB, user *User) {
	db.Create(user)
}
