package internal

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"runtime/debug"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/usiot/gbserver/internal/logger"
	"github.com/usiot/gbserver/internal/util"
)

type (
	SQLFns interface {
		SQLValues() SQLPairs
	}

	SQLPtrFns interface {
		SQLFns
		SQLFixPtr()
		SQLPtrNotNilValues() SQLPairs
		SQLPtrNotPtrValues() SQLPairs
	}

	DBX interface {
		Exec(string, ...interface{}) (sql.Result, error)
		Get(interface{}, string, ...interface{}) error
		Select(interface{}, string, ...interface{}) error
		Preparex(query string) (*sqlx.Stmt, error)
	}

	SQLPair struct {
		K string
		V interface{}
	}

	SQLPairs []SQLPair
)

// 事务处理
func Transaction(ctx context.Context, dbx *sqlx.DB, f func(ctx context.Context, tx *sqlx.Tx) (err error)) (err error) {
	var tx *sqlx.Tx
	tx, err = dbx.Beginx()
	if err != nil {
		return
	}
	defer func() {
		if p := recover(); p != nil {
			logger.ErrOp().
				Any("op=PANIC||err=", p).
				Bytes("||stack=", util.RmBSpace(debug.Stack())).
				Done()
			tx.Rollback()
		} else if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	err = f(ctx, tx)

	return
}

func Insert(ctx context.Context, dbx DBX, table string, sf SQLFns) (sql.Result, error) {
	if v, ok := sf.(interface{ SQLFixPtr() }); ok {
		v.SQLFixPtr()
	}
	now := time.Now()
	if v, ok := sf.(interface{ SetCreateTime(time.Time) }); ok {
		v.SetCreateTime(now)
	}

	vals := sf.SQLValues()
	phs := make([]string, len(vals))
	for i := 0; i < len(vals); i++ {
		phs[i] = "?"
	}
	names, values := vals.Split()

	sql := fmt.Sprintf(
		"INSERT INTO `%s`(%s) VALUES(%s);",
		table,
		strings.Join(names, ", "),
		strings.Join(phs, ", "),
	)

	args := printSQL(ctx, sql, values)
	result, err := dbx.Exec(sql, values...)
	if err != nil {
		logger.Error("op=sqlErr||errMsg=%s||sql=%s||args=%s", err, sql, args)
	}

	return result, err
}

func Delete(ctx context.Context, dbx DBX, table string, condMap map[string]interface{}) (sql.Result, error) {
	if len(condMap) == 0 { // 不允许清表
		return nil, errors.New("not allowed to empty")
	}
	var (
		cond   = make([]string, 0, len(condMap))
		values = make([]interface{}, 0, len(condMap))
	)
	for k, v := range condMap {
		cond = append(cond, k+" = ?")
		values = append(values, v)
	}

	sql := fmt.Sprintf(
		"DELETE FROM `%s` WHERE %s",
		table,
		strings.Join(cond, " AND "),
	)

	args := printSQL(ctx, sql, values)

	result, err := dbx.Exec(sql, values...)
	if err != nil {
		logger.Error("op=sqlErr||errMsg=%s||sql=%s||args=%s", err, sql, args)
	}

	return result, err
}

func Update(ctx context.Context, dbx DBX, table string, sf SQLPtrFns) (sql.Result, error) {
	setstr := []string{}
	wherestr := []string{}
	args := []interface{}{}

	now := time.Now()
	if v, ok := sf.(interface{ SetUpdateTime(time.Time) }); ok {
		v.SetUpdateTime(now)
	}

	keys, values := sf.SQLPtrNotNilValues().Split()
	args = append(args, values...)
	for _, v := range keys {
		setstr = append(setstr, v+" = ?")
	}

	keys, values = sf.SQLPtrNotPtrValues().Split()
	args = append(args, values...)
	for _, v := range keys {
		wherestr = append(wherestr, v+" = ?")
	}

	sql := fmt.Sprintf(
		"UPDATE %s SET %s WHERE %s",
		table,
		strings.Join(setstr, ", "),
		strings.Join(wherestr, ", "),
	)
	printSQL(ctx, sql, args)

	return dbx.Exec(sql, args...)
}

// 建议使用时，各字段都赋值
func InsertOrUpdate(ctx context.Context, dbx DBX, table string, sf SQLPtrFns) (err error) {
	vals := sf.SQLPtrNotPtrValues()
	if len(vals) == 0 { // 没有参数，忽略
		return
	}
	var tot int
	BeginSelectStr("1").From(table).Where(sf.SQLPtrNotPtrValues().ToMap()).GetOne(ctx, dbx, &tot)
	if tot == 1 { // 存在更新
		_, err = Update(ctx, dbx, table, sf)
	} else { // 不存在创建
		_, err = Insert(ctx, dbx, table, sf)
	}
	if err != nil {
		logger.Error(err.Error())
	}
	return
}

func (sp SQLPairs) Split() (ks []string, vs []interface{}) {
	for _, v := range sp {
		ks = append(ks, v.K)
		vs = append(vs, v.V)
	}
	return
}

func (sp SQLPairs) ToMap() map[string]interface{} {
	mp := make(map[string]interface{}, len(sp))
	for _, v := range sp {
		mp[v.K] = v.V
	}
	return mp
}

func (sp SQLPairs) Keys() (keys []string) {
	ks, _ := sp.Split()
	return ks
}

func (sp SQLPairs) Values() (values []interface{}) {
	_, vs := sp.Split()
	return vs
}

func printSQL(ctx context.Context, sql string, args []interface{}) []byte {
	vs, err := json.Marshal(args)
	logger.DbgOp().
		String("op=sql||", util.CtxLog(ctx)).
		String("||sql=", sql).
		Bytes("||args=", vs).
		IgError("||errMsg=", err).
		Done()
	return vs
}
