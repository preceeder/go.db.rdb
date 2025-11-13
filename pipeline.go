package rdb

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type RedisPipeline struct {
	lua
	builder
	Client redis.Pipeliner
}

func newPipeline(client RedisClient) *RedisPipeline {
	pip := RedisPipeline{
		Client: client.Client.Pipeline(),
	}
	pip.builder = pip.Handler
	pip.lua = pip.ExecScript
	return &pip
}

func (pip RedisPipeline) Handler(ctx context.Context, cmd RdCmd, cmdName Command, args map[string]any, includeArgs ...any) *CommandBuilder {
	// 返回 CommandBuilder，支持链式调用
	// Pipeline 中的命令会在 Exec() 时执行
	return NewPipelineCommandBuilder(pip.Client, ctx, cmd, cmdName, args, includeArgs...)
}

// 这一步才是真正的执行命令， 之前的所有步骤都是在往数组中添加命令， 实际没有发送到redis中
func (pip RedisPipeline) Exec(ctx context.Context) ([]redis.Cmder, error) {
	return pip.Client.Exec(ctx)
}
