package main

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql" // MySQL 驱动
	"github.com/jmoiron/sqlx"
)

// Book 结构体与 books 表字段一一对应，确保类型安全
// db 标签指定数据库字段名，实现结构体与表字段的映射
type Book struct {
	ID     int     `db:"id"`     // 对应表中 id 字段（整数类型）
	Title  string  `db:"title"`  // 对应表中 title 字段（字符串类型）
	Author string  `db:"author"` // 对应表中 author 字段（字符串类型）
	Price  float64 `db:"price"`  // 对应表中 price 字段（浮点类型）
}

func main() {

	var (
		dbUser     = "root"      // 数据库用户名
		dbPassword = "root"      // 数据库密码
		dbHost     = "localhost" // 数据库地址
		dbPort     = "3306"      // 数据库端口
		dbName     = "study_go"  // 数据库名
	)

	// 1. 连接数据库（使用 sqlx.Connect 建立连接）
	dataSourceName := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbUser, dbPassword, dbHost, dbPort, dbName,
	)
	// 连接字符串格式：用户名:密码@tcp(地址:端口)/数据库名?参数
	db, err := sqlx.Connect("mysql", dataSourceName)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	defer db.Close() // 程序退出前关闭数据库连接

	// 验证连接有效性
	if err := db.Ping(); err != nil {
		log.Fatalf("数据库连接测试失败: %v", err)
	}

	// 2. 查询价格大于 50 元的书籍
	priceThreshold := 50.0
	books, err := getBooksPriceGreaterThan(db, priceThreshold)
	if err != nil {
		log.Fatalf("查询书籍失败: %v", err)
	}

	// 3. 输出查询结果
	fmt.Printf("----- 价格大于 %.2f 元的书籍列表 -----\n", priceThreshold)
	if len(books) == 0 {
		fmt.Println("没有符合条件的书籍")
		return
	}
	for _, book := range books {
		fmt.Printf(
			"ID: %d, 书名: %s, 作者: %s, 价格: %.2f 元\n",
			book.ID, book.Title, book.Author, book.Price,
		)
	}
}

// getBooksPriceGreaterThan 查询价格大于指定值的书籍，返回 Book 结构体切片
func getBooksPriceGreaterThan(db *sqlx.DB, minPrice float64) ([]Book, error) {
	var books []Book // 用于接收查询结果的结构体切片

	// 使用 sqlx.Select 执行查询，自动映射到结构体切片
	// SQL 语句中使用占位符 ? 避免 SQL 注入，确保安全性
	query := "SELECT id, title, author, price FROM books WHERE price > ? ORDER BY price DESC"
	err := db.Select(&books, query, minPrice)
	if err != nil {
		return nil, fmt.Errorf("执行查询失败: %w", err) // 使用 %w 包装错误，便于上层处理
	}

	return books, nil
}
