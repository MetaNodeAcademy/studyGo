package main

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql" // MySQL 驱动
	"github.com/jmoiron/sqlx"
)

/*
题目1：使用SQL扩展库进行查询
假设你已经使用Sqlx连接到一个数据库，并且有一个 employees 表，包含字段 id 、 name 、 department 、 salary 。
要求 ：
编写Go代码，使用Sqlx查询 employees 表中所有部门为 "技术部" 的员工信息，并将结果映射到一个自定义的 Employee 结构体切片中。
编写Go代码，使用Sqlx查询 employees 表中工资最高的员工信息，并将结果映射到一个 Employee 结构体中。
*/

// 数据库配置（根据实际环境修改）
const (
	dbUser     = "root"      // 数据库用户名
	dbPassword = "root"      // 数据库密码
	dbHost     = "localhost" // 数据库地址
	dbPort     = "3306"      // 数据库端口
	dbName     = "study_go"  // 数据库名
)

// Employee 结构体映射 employees 表字段
type Employee struct {
	ID         int     `db:"id"`         // 对应表中 id 字段
	Name       string  `db:"name"`       // 对应表中 name 字段
	Department string  `db:"department"` // 对应表中 department 字段
	Salary     float64 `db:"salary"`     // 对应表中 salary 字段
}

func main() {
	// 1. 连接数据库（使用 sqlx.Connect）
	dataSourceName := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		dbUser, dbPassword, dbHost, dbPort, dbName,
	)
	db, err := sqlx.Connect("mysql", dataSourceName)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	defer db.Close()

	// 测试连接
	if err := db.Ping(); err != nil {
		log.Fatalf("数据库 ping 失败: %v", err)
	}

	// 2. 查询技术部所有员工
	techEmployees, err := getTechDepartmentEmployees(db)
	if err != nil {
		log.Printf("查询技术部员工失败: %v", err)
	} else {
		fmt.Println("----- 技术部员工列表 -----")
		for _, emp := range techEmployees {
			fmt.Printf("ID: %d, 姓名: %s, 部门: %s, 工资: %.2f\n",
				emp.ID, emp.Name, emp.Department, emp.Salary)
		}
	}

	// 3. 查询工资最高的员工
	topEmployee, err := getHighestSalaryEmployee(db)
	if err != nil {
		log.Printf("查询工资最高员工失败: %v", err)
	} else {
		fmt.Println("\n----- 工资最高的员工 -----")
		fmt.Printf("ID: %d, 姓名: %s, 部门: %s, 工资: %.2f\n",
			topEmployee.ID, topEmployee.Name, topEmployee.Department, topEmployee.Salary)
	}
}

// getTechDepartmentEmployees 查询部门为"技术部"的所有员工
func getTechDepartmentEmployees(db *sqlx.DB) ([]Employee, error) {
	var employees []Employee
	// 使用 Select 方法查询多条记录，直接映射到结构体切片
	err := db.Select(&employees, "SELECT id, name, department, salary FROM employees WHERE department = ?", "技术部")
	if err != nil {
		return nil, fmt.Errorf("查询失败: %v", err)
	}
	return employees, nil
}

// getHighestSalaryEmployee 查询工资最高的员工
func getHighestSalaryEmployee(db *sqlx.DB) (Employee, error) {
	var employee Employee
	// 使用 Get 方法查询单条记录，映射到结构体
	// 按 salary 降序排序，取第一条即为工资最高的员工
	err := db.Get(&employee, "SELECT id, name, department, salary FROM employees ORDER BY salary DESC LIMIT 1")
	if err != nil {
		return Employee{}, fmt.Errorf("查询失败: %v", err)
	}
	return employee, nil
}
