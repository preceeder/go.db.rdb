package rdb

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"time"
)

// 普通指令
type builder func(ctx context.Context, cmd RdCmd, cmdName Command, args map[string]any, includeArgs ...any) *redis.Cmd

// lua脚本
type lua func(ctx context.Context, lua LuaScript, keyInfo map[string]string, valueInfo map[string]any) *redis.Cmd
type Config struct {
	Host        string `json:"host" yaml:"host"`
	Port        string `json:"port" yaml:"port"`
	Password    string `json:"password" yaml:"password"`
	UserName    string `json:"username" yaml:"username"`
	Db          int    `json:"db" yaml:"db"`
	MaxIdle     int    `json:"maxIdle" yaml:"maxIdle"`
	MinIdle     int    `json:"minIdle" yaml:"minIdle"`
	IdleTimeout int    `json:"idleTimeout" yaml:"idleTimeout"`
	PoolSize    int    `json:"poolSize" yaml:"poolSize"`
}

type RedisClient struct {
	lua
	builder
	Config Config
	Client *redis.Client
}

func NewRedisClient(config Config) *RedisClient {
	client := RedisClient{Client: initRedis(config), Config: config}
	client.builder = client.Handler
	client.lua = client.ExecScript
	return &client
}

func initRedis(c Config) *redis.Client {
	//Rd := new(redis.Client)
	slog.Info("redisDb connect", "info", c)
	addr := c.Host + ":" + c.Port
	redisOpt := &redis.Options{
		Addr:         addr,
		Password:     c.Password,
		Username:     c.UserName,
		DB:           c.Db,
		PoolSize:     c.PoolSize,
		MaxIdleConns: c.MaxIdle,
		MinIdleConns: c.MinIdle,
	}
	rdb := redis.NewClient(redisOpt)
	rdb.AddHook(RKParesHook{})
	cmd := rdb.Ping(context.Background())
	if cmd.Err() != nil {
		panic("redis connect fail, " + cmd.Err().Error())
	}

	return rdb
}

func (rdm RedisClient) RedisClose() {
	err := rdm.Client.Close()
	if err != nil {
		slog.Error("close redisDb", "index", rdm.Config.Db, "error", err.Error())
	} else {
		slog.Info("close redisDb", "index", rdm.Config.Db)
	}
}

func (rdm RedisClient) Handler(ctx context.Context, cmd RdCmd, cmdName Command, args map[string]any, includeArgs ...any) *redis.Cmd {
	cmdList, key, subCmd := Build(ctx, cmd, cmdName, args, includeArgs...)
	resultCmd := rdm.Client.Do(ctx, cmdList...)
	// 如果是 nil 且不返回的话就处理一下
	if !subCmd.ReturnNilError {
		if errors.Is(resultCmd.Err(), redis.Nil) {
			resultCmd.SetErr(nil)
		}
	}
	rdm.setExp(ctx, key, subCmd)
	return resultCmd
}

func (rdm RedisClient) setExp(ctx context.Context, key string, subCmd RdSubCmd) {
	if subCmd.Exp != nil {
		exp := subCmd.Exp()
		expireCmd := rdm.Client.Expire(ctx, key, exp)
		if expireCmd.Err() != nil {
			slog.ErrorContext(ctx, "set expire", "key", key, "error", expireCmd.Err())
		}
	}
}

func (rdm RedisClient) PipeLine() *RedisPipeline {
	return newPipeline(rdm)
}

// Exec 执行 Redis 命令
func Exec(ctx context.Context, rdb *redis.Client, cmd redis.Cmder, key string, exp time.Duration) error {
	// 使用 go-redis 执行命令
	err := rdb.Process(ctx, cmd) // 这一步可以单独执行
	if err == nil && exp > 0 {
		if er := rdb.Expire(ctx, key, exp).Err(); er != nil {
			slog.ErrorContext(ctx, "set exp", "key", key, "exp", exp, "error", er.Error())
		}
	}
	return err
}
