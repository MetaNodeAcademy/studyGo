package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // MySQL 驱动（需提前安装：go get github.com/go-sql-driver/mysql）
	"log"
)

/**
假设有两个表： accounts 表（包含字段 id 主键， balance 账户余额）和 transactions
表（包含字段 id 主键， from_account_id 转出账户ID， to_account_id 转入账户ID， amount 转账金额）。
要求 ：
编写一个事务，实现从账户 A 向账户 B 转账 100 元的操作。在事务中，需要先检查账户 A 的余额是否足够，
如果足够则从账户 A 扣除 100 元，向账户 B 增加 100 元，并在 transactions 表中记录该笔转账信息。如果余额不足，则回滚事务。
*/

// transfer 实现从 fromAccount 向 toAccount 转账 amount 金额的事务
func transfer(db *sql.DB, fromAccount, toAccount int, amount float64) error {
	// 1. 开启事务
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("开启事务失败: %v", err)
	}
	// 延迟处理：若后续操作出错，自动回滚事务
	defer func() {
		if r := recover(); r != nil {
			// 发生恐慌时回滚
			isRollback(tx, err)
			log.Printf("事务因恐慌回滚: %v", r)
		}
	}()

	// 2. 查询转出账户余额（加行锁，防止并发修改）
	var fromBalance float64
	// FOR UPDATE 用于锁定当前行，避免其他事务同时修改该账户余额
	err = tx.QueryRow(
		"SELECT balance FROM accounts WHERE id = ? FOR UPDATE",
		fromAccount,
	).Scan(&fromBalance)
	if err != nil {
		isRollback(tx, err) // 查询失败，回滚
		return fmt.Errorf("查询账户 %d 余额失败: %v", fromAccount, err)
	}

	// 3. 检查余额是否充足
	if fromBalance < amount {
		isRollback(tx, err) // 余额不足，回滚
		return fmt.Errorf("账户 %d 余额不足（当前: %.2f, 需转出: %.2f）", fromAccount, fromBalance, amount)
	}

	// 4. 扣减转出账户余额
	_, err = tx.Exec(
		"UPDATE accounts SET balance = balance - ? WHERE id = ?",
		amount, fromAccount,
	)
	if err != nil {
		isRollback(tx, err) // 更新失败，回滚
		return fmt.Errorf("扣减账户 %d 余额失败: %v", fromAccount, err)
	}

	// 5. 增加转入账户余额
	_, err = tx.Exec(
		"UPDATE accounts SET balance = balance + ? WHERE id = ?",
		amount, toAccount,
	)
	if err != nil {
		isRollback(tx, err) // 更新失败，回滚
		return fmt.Errorf("增加账户 %d 余额失败: %v", toAccount, err)
	}

	// 6. 记录转账交易
	_, err = tx.Exec(
		"INSERT INTO transactions (from_account_id, to_account_id, amount) VALUES (?, ?, ?)",
		fromAccount, toAccount, amount,
	)
	if err != nil {
		isRollback(tx, err) // 记录失败，回滚
		return fmt.Errorf("记录交易失败: %v", err)
	}

	// 7. 所有操作成功，提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("提交事务失败: %v", err)
	}
	return nil
}

func isRollback(tx *sql.Tx, err error) {
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			er := errors.New("rollback fail")
			if er != nil {
				fmt.Println(er)
			}
		}
	}
}

func main() {
	// 连接 MySQL 数据库
	dataSourceName := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		"root", "root", "127.0.0.1", "3306", "study_go",
	)
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatalf("无法连接数据库: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(db)

	// 测试连接是否有效
	if err := db.Ping(); err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	// 执行转账：从账户 1 向账户 2 转账 100 元
	err = transfer(db, 1, 2, 100)
	if err != nil {
		log.Fatalf("转账失败: %v", err)
	}
	fmt.Println("转账成功！")
}
