package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/extra/rediscmd/v9"
	"github.com/redis/go-redis/v9"
	"github.com/usiot/gbserver/internal/logger"
	"github.com/usiot/gbserver/internal/util"
)

var (
	rds *redis.Client
)

func Init(redisURL string) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		logger.Fatal("redis url解析失败: %s", err)
	}
	rds = redis.NewClient(opt)

	err = rds.Ping(context.Background()).Err()
	if err != nil {
		logger.Fatal("redis连接失败: %s", err)
	}

	rds.AddHook(hook{})
}

type hook struct{}

func (h hook) DialHook(next redis.DialHook) redis.DialHook { return next }

func (h hook) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		start := time.Now()
		err := next(context.Background(), cmd)
		if err != nil {
			logger.ErrOp().
				String("op=rdsErr||", util.CtxLog(ctx)).
				Int64("||cost=", time.Since(start).Milliseconds()).
				String("||cmd=", rediscmd.CmdString(cmd)).
				Error("||errMsg=", err).
				Done()
		} else {
			logger.DbgOp().
				String("op=rdsDtl||", util.CtxLog(ctx)).
				Int64("||cost=", time.Since(start).Milliseconds()).
				String("||cmd=", rediscmd.CmdString(cmd)).
				Done()
		}

		return err
	}
}

func (h hook) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return next
}
