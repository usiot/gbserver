package cache

import "context"

const (
	CeqKey = "GB:MEDIA:CEQ"
)

func Incr(ctx context.Context, key string) (int64, error) { return rds.Incr(ctx, key).Result() }

func GetCeq(ctx context.Context) (int64, error) { return Incr(ctx, CeqKey) }
