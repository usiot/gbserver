package internal

import (
	"bytes"
	"context"
	"strconv"
	"strings"
	"sync"

	"github.com/usiot/gbserver/internal/logger"
	"github.com/usiot/gbserver/internal/util"
)

type (
	sqlselect struct {
		sel    string
		from   string
		join   string
		on     string
		where  string
		group  string
		having string
		order  string
		limit  []int
		offset int
		args   []interface{}
	}
)

var (
	plsqlsel = sync.Pool{New: func() interface{} { return &sqlselect{} }}
	plstring = sync.Pool{New: func() interface{} { return &bytes.Buffer{} }}
)

func BeginSelect(sf SQLFns) *sqlselect {
	ss := plsqlsel.Get().(*sqlselect)
	ss.Clear()

	keys, _ := sf.SQLValues().Split()

	ss.sel = strings.Join(keys, ", ")
	return ss
}

func BeginSelectStr(sel string) *sqlselect {
	ss := plsqlsel.Get().(*sqlselect)
	ss.Clear()

	ss.sel = sel
	return ss
}

func (ss *sqlselect) From(tablename string) *sqlselect {
	ss.from = tablename
	return ss
}

func (ss *sqlselect) Join(tablename string) *sqlselect {
	ss.join = tablename
	return ss
}

func (ss *sqlselect) On(args map[string]interface{}) *sqlselect {
	field := []string{}
	for k, v := range args {
		field = append(field, k+" = ? ")
		ss.args = append(ss.args, v)
	}
	ss.on += strings.Join(field, " AND ")

	return ss
}

func (ss *sqlselect) OnStr(on string, args ...interface{}) *sqlselect {
	ss.on += on
	ss.args = append(ss.args, args...)
	return ss
}

func (ss *sqlselect) Where(args map[string]interface{}) *sqlselect {
	field := []string{}
	for k, v := range args {
		field = append(field, k+" = ? ")
		ss.args = append(ss.args, v)
	}
	ss.where += strings.Join(field, " AND ")

	return ss
}

func (ss *sqlselect) WhereStr(where string, args []interface{}) *sqlselect {
	if len(where) > 0 {
		ss.where += where
		ss.args = append(ss.args, args...)
	}
	return ss
}

// l 一个或两个值，否者忽略
func (ss *sqlselect) Limit(l ...int) *sqlselect {
	switch len(l) {
	case 1:
		ss.limit = []int{l[0]}
	case 2:
		ss.limit = []int{l[0], l[1]}
	}
	// ss.limit = l
	return ss
}

func (ss *sqlselect) Offset(o int) *sqlselect {
	ss.offset = o
	return ss
}

func (ss *sqlselect) Group(group string) *sqlselect {
	ss.group = group
	return ss
}

// having 也会有 args
func (ss *sqlselect) Having(having string) *sqlselect {
	ss.having = having
	return ss
}

func (ss *sqlselect) Order(order string) *sqlselect {
	ss.order = order
	return ss
}

func (ss *sqlselect) String() string {
	bs := plstring.Get().(*bytes.Buffer)
	bs.Reset()
	defer func() {
		if bs.Cap() < 1 {
			plstring.Put(bs)
		}
	}()
	bs.WriteString(" SELECT ")
	bs.WriteString(ss.sel)
	bs.WriteString(" FROM ")
	bs.WriteString(ss.from)

	if ss.join != "" {
		bs.WriteString(" JOIN ")
		bs.WriteString(ss.join)
		bs.WriteByte(' ')
		if ss.on != "" {
			bs.WriteString(" ON ")
			bs.WriteString(ss.on)
			bs.WriteByte(' ')
		}
	}

	if ss.where != "" {
		bs.WriteString(" WHERE ")
		bs.WriteString(ss.where)
		bs.WriteByte(' ')
	}

	if ss.group != "" {
		bs.WriteString(" GROUP BY ")
		bs.WriteString(ss.group)
		if ss.having != "" {
			bs.WriteString(" HAVING ")
			bs.WriteString(ss.having)
		}
	}

	if ss.order != "" {
		bs.WriteString(" ORDER BY ")
		bs.WriteString(ss.order)
	}

	if len(ss.limit) != 0 {
		bs.WriteString(" LIMIT ")
		bs.WriteString(strconv.Itoa(ss.limit[0]))
		if len(ss.limit) == 2 {
			bs.WriteByte(',')
			bs.WriteString(strconv.Itoa(ss.limit[1]))
		}
	}

	if ss.offset != 0 {
		bs.WriteString(" OFFSET ")
		bs.WriteString(strconv.Itoa(ss.offset))
	}
	return bs.String()
}

func (ss *sqlselect) GetOne(ctx context.Context, db DBX, data interface{}) error {
	defer plsqlsel.Put(ss)
	sql := ss.String()
	if len(ss.args) == 0 {
		logger.DbgOp().
			String("op=sql||", util.CtxLog(ctx)).
			String("||sql=", sql).
			Done()
	} else {
		printSQL(ctx, sql, ss.args)
	}
	return db.Get(data, sql, ss.args...)
}

func (ss *sqlselect) GetMulti(ctx context.Context, DB DBX, data interface{}) error {
	defer plsqlsel.Put(ss)
	sql := ss.String()
	if len(ss.args) == 0 {
		logger.DbgOp().
			String("op=sql||", util.CtxLog(ctx)).
			String("||sql=", sql).
			Done()
	} else {
		printSQL(ctx, sql, ss.args)
	}
	return DB.Select(data, sql, ss.args...)
}

func (ss *sqlselect) Clear() {
	ss.sel = ""
	ss.from = ""
	ss.join = ""
	ss.on = ""
	ss.where = ""
	ss.group = ""
	ss.having = ""
	ss.order = ""
	ss.limit = ss.limit[:0]
	ss.offset = 0
	ss.args = ss.args[:0]
}

// >=> `{elem}{sep}{elem}`
func RepeatJoin(elem, sep string, count int) string {
	switch count {
	case 0:
		return ""
	case 1:
		return elem
	}

	buf := strings.Builder{}
	buf.Grow((len(elem)+len(sep))*(count-1) + len(elem))

	buf.WriteString(elem)
	for i := 0; i < count-1; i++ {
		buf.WriteString(sep)
		buf.WriteString(elem)
	}

	return buf.String()
}
