package util

import (
	"context"
	"strconv"
	"strings"
	"time"
)

type CtxField string

const (
	CtxTraceId CtxField = "traceId"
	CtxSpanId  CtxField = "spanId"
	CtxParent  CtxField = "parentId"

	CtxDid  CtxField = "did"
	CtxCid  CtxField = "cid"
	CtxCost CtxField = "cost"
	CtxTime CtxField = "time"
)

func CtxString(ctx context.Context, key CtxField) string {
	if ctx == nil {
		return ""
	}

	v, _ := ctx.Value(key).(string)
	return v
}

func CtxLog(ctx context.Context) string {
	buf := strings.Builder{}
	for _, v := range []CtxField{
		CtxDid, CtxCid, CtxCost,
		CtxTraceId, CtxSpanId, CtxParent,
	} {
		if val := CtxString(ctx, v); val != "" {
			buf.WriteString(string(v))
			buf.Write([]byte("="))
			buf.WriteString(val)
			buf.WriteString("||")
		}
	}
	buf.WriteString("time=")
	buf.WriteString(strconv.FormatInt(time.Now().UnixMilli(), 10))

	return buf.String()
}
