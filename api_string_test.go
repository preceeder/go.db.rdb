package rdb

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// 测试 String 操作的 RdCmd 定义
var StringCmd = RdCmd{
	Key: "string:{{keyName}}",
	CMD: map[Command]RdSubCmd{
		SET: {
			Params: "{{value}}",
			Exp: func() time.Duration {
				return time.Second * 30
			},
		},
		GET: {
			Params: "",
		},
		MSET: {
			Params:   "",
			NoUseKey: true,
		},
		MGET: {
			Params:   "",
			NoUseKey: true,
		},
		SETEX: {
			Params: "{{seconds}} {{value}}",
		},
		SETNX: {
			Params: "{{value}}",
		},
		GETSET: {
			Params: "{{value}}",
		},
		GETRANGE: {
			Params: "{{start}} {{end}}",
		},
		SETRANGE: {
			Params: "{{offset}} {{value}}",
		},
		INCR: {
			Params: "",
		},
		INCRBY: {
			Params: "{{increment}}",
		},
		INCRBYFLOAT: {
			Params: "{{increment}}",
		},
		DECR: {
			Params: "",
		},
		DECRBY: {
			Params: "{{decrement}}",
		},
		APPEND: {
			Params: "{{value}}",
		},
	},
}

// TestRedisClient_Set 测试 SET 命令
func TestRedisClient_Set(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	cmd := client.Set(context.Background(), StringCmd, map[string]any{
		"keyName": "test1",
		"value":   "Hello World",
	})

	if cmd.Err() != nil {
		t.Errorf("Set failed: %v", cmd.Err())
		return
	}

	fmt.Printf("SET result: %v\n", cmd.Val())
}

// TestRedisClient_Get 测试 GET 命令
func TestRedisClient_Get(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// 先设置一个值
	client.Set(context.Background(), StringCmd, map[string]any{
		"keyName": "test2",
		"value":   "Get Test Value",
	})

	// 获取值
	cmd := client.Get(context.Background(), StringCmd, map[string]any{
		"keyName": "test2",
	})

	if cmd.Err() != nil {
		t.Errorf("Get failed: %v", cmd.Err())
		return
	}

	value := cmd.Val()
	fmt.Printf("GET result: %s\n", value)
}

// TestRedisClient_MSet 测试 MSET 命令
func TestRedisClient_MSet(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// MSET 需要多个 key-value 对
	cmd := client.MSet(context.Background(), StringCmd, nil,
		"key1", "value1",
		"key2", "value2",
		"key3", "value3",
	)

	if cmd.Err() != nil {
		t.Errorf("MSet failed: %v", cmd.Err())
		return
	}

	fmt.Printf("MSET result: %v\n", cmd.Val())
}

// TestRedisClient_MGet 测试 MGET 命令
func TestRedisClient_MGet(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// 先设置多个值
	client.MSet(context.Background(), StringCmd, nil,
		"mget1", "value1",
		"mget2", "value2",
		"mget3", "value3",
	)

	// 批量获取
	cmd := client.MGet(context.Background(), StringCmd, nil,
		"mget1",
		"mget2",
		"mget3",
		"mget4", // 不存在的 key
	)

	if cmd.Err() != nil {
		t.Errorf("MGet failed: %v", cmd.Err())
		return
	}

	values := cmd.Val()
	fmt.Printf("MGET result: %v\n", values)
}

// TestRedisClient_SetEx 测试 SETEX 命令（设置过期时间）
func TestRedisClient_SetEx(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	cmd := client.SetEx(context.Background(), StringCmd, map[string]any{
		"keyName": "setex1",
		"seconds": 60,
		"value":   "This value will expire in 60 seconds",
	})

	if cmd.Err() != nil {
		t.Errorf("SetEx failed: %v", cmd.Err())
		return
	}

	fmt.Printf("SETEX result: %v\n", cmd.Val())
}

// TestRedisClient_SetNx 测试 SETNX 命令（仅当 key 不存在时设置）
func TestRedisClient_SetNx(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	keyName := "setnx1"

	// 第一次设置，应该成功
	cmd1 := client.SetNx(context.Background(), StringCmd, map[string]any{
		"keyName": keyName,
		"value":   "First value",
	})

	if cmd1.Err() != nil {
		t.Errorf("First SetNx failed: %v", cmd1.Err())
		return
	}

	fmt.Printf("First SETNX result: %v (should be 1)\n", cmd1.Val())

	// 第二次设置，应该失败（key 已存在）
	cmd2 := client.SetNx(context.Background(), StringCmd, map[string]any{
		"keyName": keyName,
		"value":   "Second value",
	})

	if cmd2.Err() != nil {
		t.Errorf("Second SetNx failed: %v", cmd2.Err())
		return
	}

	fmt.Printf("Second SETNX result: %v (should be 0)\n", cmd2.Val())
}

// TestRedisClient_GetSet 测试 GETSET 命令
func TestRedisClient_GetSet(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	keyName := "getset1"

	// 先设置一个值
	client.Set(context.Background(), StringCmd, map[string]any{
		"keyName": keyName,
		"value":   "Old value",
	})

	// 使用 GETSET 获取旧值并设置新值
	cmd := client.GetSet(context.Background(), StringCmd, map[string]any{
		"keyName": keyName,
		"value":   "New value",
	})

	if cmd.Err() != nil {
		t.Errorf("GetSet failed: %v", cmd.Err())
		return
	}

	oldValue := cmd.Val()
	fmt.Printf("GETSET old value: %s\n", oldValue)

	// 验证新值
	getCmd := client.Get(context.Background(), StringCmd, map[string]any{
		"keyName": keyName,
	})
	fmt.Printf("GETSET new value: %s\n", getCmd.Val())
}

// TestRedisClient_GetRange 测试 GETRANGE 命令
func TestRedisClient_GetRange(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	keyName := "getrange1"
	value := "Hello Redis"

	// 设置值
	client.Set(context.Background(), StringCmd, map[string]any{
		"keyName": keyName,
		"value":   value,
	})

	// 获取子字符串
	cmd := client.GetRange(context.Background(), StringCmd, map[string]any{
		"keyName": keyName,
		"start":   0,
		"end":     4,
	})

	if cmd.Err() != nil {
		t.Errorf("GetRange failed: %v", cmd.Err())
		return
	}

	fmt.Printf("GETRANGE(0, 4) of '%s': %s\n", value, cmd.Val())
}

// TestRedisClient_SetRange 测试 SETRANGE 命令
func TestRedisClient_SetRange(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	keyName := "setrange1"

	// 先设置一个值
	client.Set(context.Background(), StringCmd, map[string]any{
		"keyName": keyName,
		"value":   "Hello World",
	})

	// 从偏移量 6 开始覆写
	cmd := client.SetRange(context.Background(), StringCmd, map[string]any{
		"keyName": keyName,
		"offset":  6,
		"value":   "Redis",
	})

	if cmd.Err() != nil {
		t.Errorf("SetRange failed: %v", cmd.Err())
		return
	}

	fmt.Printf("SETRANGE result length: %d\n", cmd.Val())

	// 获取更新后的值
	getCmd := client.Get(context.Background(), StringCmd, map[string]any{
		"keyName": keyName,
	})
	fmt.Printf("After SETRANGE: %s\n", getCmd.Val())
}

// TestRedisClient_Incr 测试 INCR 命令
func TestRedisClient_Incr(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	keyName := "incr1"

	// 先设置一个数字
	client.Set(context.Background(), StringCmd, map[string]any{
		"keyName": keyName,
		"value":   "10",
	})

	// 自增
	cmd := client.Incr(context.Background(), StringCmd, map[string]any{
		"keyName": keyName,
	})

	if cmd.Err() != nil {
		t.Errorf("Incr failed: %v", cmd.Err())
		return
	}

	fmt.Printf("INCR result: %d\n", cmd.Val())
}

// TestRedisClient_IncrBy 测试 INCRBY 命令
func TestRedisClient_IncrBy(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	keyName := "incrby1"

	// 先设置一个数字
	client.Set(context.Background(), StringCmd, map[string]any{
		"keyName": keyName,
		"value":   "10",
	})

	// 增加 5
	cmd := client.IncrBy(context.Background(), StringCmd, map[string]any{
		"keyName":  keyName,
		"increment": 5,
	})

	if cmd.Err() != nil {
		t.Errorf("IncrBy failed: %v", cmd.Err())
		return
	}

	fmt.Printf("INCRBY 5 result: %d\n", cmd.Val())
}

// TestRedisClient_IncrByFloat 测试 INCRBYFLOAT 命令
func TestRedisClient_IncrByFloat(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	keyName := "incrbyfloat1"

	// 先设置一个浮点数
	client.Set(context.Background(), StringCmd, map[string]any{
		"keyName": keyName,
		"value":   "10.5",
	})

	// 增加 2.3
	cmd := client.IncrByFloat(context.Background(), StringCmd, map[string]any{
		"keyName":  keyName,
		"increment": 2.3,
	})

	if cmd.Err() != nil {
		t.Errorf("IncrByFloat failed: %v", cmd.Err())
		return
	}

	fmt.Printf("INCRBYFLOAT 2.3 result: %s\n", cmd.Val())
}

// TestRedisClient_Decr 测试 DECR 命令
func TestRedisClient_Decr(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	keyName := "decr1"

	// 先设置一个数字
	client.Set(context.Background(), StringCmd, map[string]any{
		"keyName": keyName,
		"value":   "10",
	})

	// 自减
	cmd := client.Decr(context.Background(), StringCmd, map[string]any{
		"keyName": keyName,
	})

	if cmd.Err() != nil {
		t.Errorf("Decr failed: %v", cmd.Err())
		return
	}

	fmt.Printf("DECR result: %d\n", cmd.Val())
}

// TestRedisClient_DecrBy 测试 DECRBY 命令
func TestRedisClient_DecrBy(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	keyName := "decrby1"

	// 先设置一个数字
	client.Set(context.Background(), StringCmd, map[string]any{
		"keyName": keyName,
		"value":   "10",
	})

	// 减少 3
	cmd := client.DecrBy(context.Background(), StringCmd, map[string]any{
		"keyName":  keyName,
		"decrement": 3,
	})

	if cmd.Err() != nil {
		t.Errorf("DecrBy failed: %v", cmd.Err())
		return
	}

	fmt.Printf("DECRBY 3 result: %d\n", cmd.Val())
}

// TestRedisClient_Append 测试 APPEND 命令
func TestRedisClient_Append(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	keyName := "append1"

	// 先设置一个值
	client.Set(context.Background(), StringCmd, map[string]any{
		"keyName": keyName,
		"value":   "Hello",
	})

	// 追加字符串
	cmd := client.StringAppend(context.Background(), StringCmd, map[string]any{
		"keyName": keyName,
		"value":   " World",
	})

	if cmd.Err() != nil {
		t.Errorf("Append failed: %v", cmd.Err())
		return
	}

	fmt.Printf("APPEND result length: %d\n", cmd.Val())

	// 获取最终值
	getCmd := client.Get(context.Background(), StringCmd, map[string]any{
		"keyName": keyName,
	})
	fmt.Printf("After APPEND: %s\n", getCmd.Val())
}

// TestRedisClient_String_Integration 集成测试：String 操作的完整流程
func TestRedisClient_String_Integration(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	keyName := "integration_string"

	// 1. SET
	client.Set(context.Background(), StringCmd, map[string]any{
		"keyName": keyName,
		"value":   "Initial",
	})

	// 2. GET
	getCmd := client.Get(context.Background(), StringCmd, map[string]any{
		"keyName": keyName,
	})
	fmt.Printf("1. GET: %s\n", getCmd.Val())

	// 3. APPEND
	client.StringAppend(context.Background(), StringCmd, map[string]any{
		"keyName": keyName,
		"value":   " Value",
	})

	// 4. GETRANGE
	rangeCmd := client.GetRange(context.Background(), StringCmd, map[string]any{
		"keyName": keyName,
		"start":   0,
		"end":     6,
	})
	fmt.Printf("2. GETRANGE(0,6): %s\n", rangeCmd.Val())

	// 5. INCR (需要先设置为数字)
	client.Set(context.Background(), StringCmd, map[string]any{
		"keyName": keyName + "_num",
		"value":   "100",
	})
	incrCmd := client.Incr(context.Background(), StringCmd, map[string]any{
		"keyName": keyName + "_num",
	})
	fmt.Printf("3. INCR: %d\n", incrCmd.Val())
}

