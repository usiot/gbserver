package dao

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
)

func init() {
	var (
		user     = "root"
		password = "1234"
		host     = "127.0.0.1"
		port     = "3306"
		dbname   = "usiot"
	)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, dbname)
	db = sqlx.MustConnect("mysql", dsn)
}

func TestNamed(t *testing.T) {
	data := map[string]interface{}{
		"id":   time.Now().UnixMilli(),
		"name": "jmoiron",
	}
	insertSQL := `
		INSERT INTO test (id, name)
		VALUES (:id, :name)
	`

	// 使用db.Rebind重新绑定SQL语句中的参数
	query, args, err := sqlx.Named(insertSQL, data)
	if err != nil {
		log.Fatal(err)
	}

	// 打印完整的SQL语句和参数值
	fmt.Println("Executing SQL:", db.Rebind(query))
	fmt.Println("With Args:", args)

	// 执行INSERT语句
	_, err = db.NamedExec(query, data)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Insert completed successfully!")
}
