package dao

import (
	"context"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql" // 或者其他数据库驱动
	"github.com/jmoiron/sqlx"
	"github.com/usiot/gbserver/internal/dao/internal"
)

type (
	SQLFns    = internal.SQLFns
	SQLPtrFns = internal.SQLPtrFns
	SQLPair   = internal.SQLPair
	SQLPairs  = internal.SQLPairs
)

const (
	TableDevice = "device"
)

var (
	db *sqlx.DB
)

func Init(host string, port uint16, user, password, dbname string) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, dbname)
	db = sqlx.MustConnect("mysql", dsn)
	db.SetConnMaxIdleTime(300 * time.Second)
	db.SetMaxOpenConns(120)
	db.SetMaxIdleConns(100)
}

func Insert(ctx context.Context, table string, sf internal.SQLFns) error {
	_, err := internal.Insert(ctx, db, table, sf)
	return err
}

func Update(ctx context.Context, table string, sf internal.SQLPtrFns) error {
	_, err := internal.Update(ctx, db, table, sf)
	return err
}

func InsertOrUpdate(ctx context.Context, table string, sf internal.SQLPtrFns) error {
	err := internal.InsertOrUpdate(ctx, db, table, sf)
	return err
}

func Delete(ctx context.Context, table string, cond map[string]interface{}) error {
	_, err := internal.Delete(ctx, db, table, cond)
	return err
}
