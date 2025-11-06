package rdb

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// 测试 List 操作的 RdCmd 定义
var ListCmd = RdCmd{
	Key: "list:{{keyName}}",
	CMD: map[Command]RdSubCmd{
		LPUSH: {
			Params: "{{value}}",
			Exp: func() time.Duration {
				return time.Second * 30
			},
		},
		RPUSH: {
			Params: "{{value}}",
		},
		LPOP: {
			Params: "",
		},
		RPOP: {
			Params: "",
		},
		LLEN: {
			Params: "",
		},
		LINDEX: {
			Params: "{{index}}",
		},
		LRANGE: {
			Params: "{{start}} {{stop}}",
		},
		LREM: {
			Params: "{{count}} {{value}}",
		},
		LSET: {
			Params: "{{index}} {{value}}",
		},
		LTRIM: {
			Params: "{{start}} {{stop}}",
		},
		LINSERT: {
			Params: "{{position}} {{pivot}} {{value}}",
		},
		LPUSHX: {
			Params: "{{value}}",
		},
		RPUSHX: {
			Params: "{{value}}",
		},
		RPOPLPUSH: {
			Params: "{{target}}",
		},
	},
}

// TestRedisClient_LPush 测试 LPUSH 命令
func TestRedisClient_LPush(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// LPUSH 从左侧插入多个值
	cmd := client.LPush(context.Background(), ListCmd, map[string]any{
		"keyName": "test1",
		"value":   "value1",
	}, "value2", "value3")

	if cmd.Err() != nil {
		t.Errorf("LPush failed: %v", cmd.Err())
		return
	}

	fmt.Printf("LPUSH result length: %d\n", cmd.Val())
}

// TestRedisClient_RPush 测试 RPUSH 命令
func TestRedisClient_RPush(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// RPUSH 从右侧插入多个值
	cmd := client.RPush(context.Background(), ListCmd, map[string]any{
		"keyName": "test2",
		"value":   "value1",
	}, "value2", "value3")

	if cmd.Err() != nil {
		t.Errorf("RPush failed: %v", cmd.Err())
		return
	}

	fmt.Printf("RPUSH result length: %d\n", cmd.Val())
}

// TestRedisClient_LPop 测试 LPOP 命令
func TestRedisClient_LPop(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// 先创建列表
	client.LPush(context.Background(), ListCmd, map[string]any{
		"keyName": "test3",
		"value":   "first",
	}, "second", "third")

	// 从左侧弹出
	cmd := client.LPop(context.Background(), ListCmd, map[string]any{
		"keyName": "test3",
	})

	if cmd.Err() != nil {
		t.Errorf("LPop failed: %v", cmd.Err())
		return
	}

	fmt.Printf("LPOP result: %s\n", cmd.Val())
}

// TestRedisClient_RPop 测试 RPOP 命令
func TestRedisClient_RPop(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// 先创建列表
	client.RPush(context.Background(), ListCmd, map[string]any{
		"keyName": "test4",
		"value":   "first",
	}, "second", "third")

	// 从右侧弹出
	cmd := client.RPop(context.Background(), ListCmd, map[string]any{
		"keyName": "test4",
	})

	if cmd.Err() != nil {
		t.Errorf("RPop failed: %v", cmd.Err())
		return
	}

	fmt.Printf("RPOP result: %s\n", cmd.Val())
}

// TestRedisClient_LLen 测试 LLEN 命令
func TestRedisClient_LLen(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// 先创建列表
	client.LPush(context.Background(), ListCmd, map[string]any{
		"keyName": "test5",
		"value":   "item1",
	}, "item2", "item3", "item4")

	// 获取列表长度
	cmd := client.LLen(context.Background(), ListCmd, map[string]any{
		"keyName": "test5",
	})

	if cmd.Err() != nil {
		t.Errorf("LLen failed: %v", cmd.Err())
		return
	}

	fmt.Printf("LLEN result: %d\n", cmd.Val())
}

// TestRedisClient_LIndex 测试 LINDEX 命令
func TestRedisClient_LIndex(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// 先创建列表
	client.RPush(context.Background(), ListCmd, map[string]any{
		"keyName": "test6",
		"value":   "zero",
	}, "one", "two", "three")

	// 获取索引 1 的元素
	cmd := client.LIndex(context.Background(), ListCmd, map[string]any{
		"keyName": "test6",
		"index":   1,
	})

	if cmd.Err() != nil {
		t.Errorf("LIndex failed: %v", cmd.Err())
		return
	}

	fmt.Printf("LINDEX(1) result: %s\n", cmd.Val())
}

// TestRedisClient_LRange 测试 LRANGE 命令
func TestRedisClient_LRange(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// 先创建列表
	client.RPush(context.Background(), ListCmd, map[string]any{
		"keyName": "test7",
		"value":   "a",
	}, "b", "c", "d", "e")

	// 获取范围 [0, 2]
	cmd := client.LRange(context.Background(), ListCmd, map[string]any{
		"keyName": "test7",
		"start":   0,
		"stop":    2,
	})

	if cmd.Err() != nil {
		t.Errorf("LRange failed: %v", cmd.Err())
		return
	}

	values := cmd.Val()
	fmt.Printf("LRANGE(0, 2) result: %v\n", values)

	// 获取所有元素
	cmd2 := client.LRange(context.Background(), ListCmd, map[string]any{
		"keyName": "test7",
		"start":   0,
		"stop":    -1, // -1 表示最后一个元素
	})
	fmt.Printf("LRANGE(0, -1) result: %v\n", cmd2.Val())
}

// TestRedisClient_LRem 测试 LREM 命令
func TestRedisClient_LRem(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// 先创建包含重复元素的列表
	client.RPush(context.Background(), ListCmd, map[string]any{
		"keyName": "test8",
		"value":   "apple",
	}, "banana", "apple", "cherry", "apple")

	// 从左侧删除 2 个 "apple" (count > 0)
	cmd := client.LRem(context.Background(), ListCmd, map[string]any{
		"keyName": "test8",
		"count":   2,
		"value":   "apple",
	})

	if cmd.Err() != nil {
		t.Errorf("LRem failed: %v", cmd.Err())
		return
	}

	fmt.Printf("LREM count=2 result: %d elements removed\n", cmd.Val())

	// 查看剩余元素
	rangeCmd := client.LRange(context.Background(), ListCmd, map[string]any{
		"keyName": "test8",
		"start":   0,
		"stop":    -1,
	})
	fmt.Printf("Remaining elements: %v\n", rangeCmd.Val())
}

// TestRedisClient_LSet 测试 LSET 命令
func TestRedisClient_LSet(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// 先创建列表
	client.RPush(context.Background(), ListCmd, map[string]any{
		"keyName": "test9",
		"value":   "old1",
	}, "old2", "old3")

	// 设置索引 1 的值为新值
	cmd := client.LSet(context.Background(), ListCmd, map[string]any{
		"keyName": "test9",
		"index":   1,
		"value":   "new2",
	})

	if cmd.Err() != nil {
		t.Errorf("LSet failed: %v", cmd.Err())
		return
	}

	fmt.Printf("LSET result: %v\n", cmd.Val())

	// 验证修改
	indexCmd := client.LIndex(context.Background(), ListCmd, map[string]any{
		"keyName": "test9",
		"index":   1,
	})
	fmt.Printf("After LSET, LINDEX(1): %s\n", indexCmd.Val())
}

// TestRedisClient_LTrim 测试 LTRIM 命令
func TestRedisClient_LTrim(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// 先创建列表
	client.RPush(context.Background(), ListCmd, map[string]any{
		"keyName": "test10",
		"value":   "a",
	}, "b", "c", "d", "e", "f")

	// 修剪列表，只保留索引 1 到 3 的元素
	cmd := client.LTrim(context.Background(), ListCmd, map[string]any{
		"keyName": "test10",
		"start":   1,
		"stop":    3,
	})

	if cmd.Err() != nil {
		t.Errorf("LTrim failed: %v", cmd.Err())
		return
	}

	fmt.Printf("LTRIM result: %v\n", cmd.Val())

	// 查看修剪后的列表
	rangeCmd := client.LRange(context.Background(), ListCmd, map[string]any{
		"keyName": "test10",
		"start":   0,
		"stop":    -1,
	})
	fmt.Printf("After LTRIM: %v\n", rangeCmd.Val())
}

// TestRedisClient_LInsert 测试 LINSERT 命令
func TestRedisClient_LInsert(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// 先创建列表
	client.RPush(context.Background(), ListCmd, map[string]any{
		"keyName": "test11",
		"value":   "world",
	}, "hello", "redis")

	// 在 "world" 之前插入 "new"
	cmd := client.LInsert(context.Background(), ListCmd, map[string]any{
		"keyName": "test11",
		"position": "BEFORE",
		"pivot":   "world",
		"value":   "new",
	})

	if cmd.Err() != nil {
		t.Errorf("LInsert failed: %v", cmd.Err())
		return
	}

	fmt.Printf("LINSERT result: new length %d\n", cmd.Val())

	// 查看插入后的列表
	rangeCmd := client.LRange(context.Background(), ListCmd, map[string]any{
		"keyName": "test11",
		"start":   0,
		"stop":    -1,
	})
	fmt.Printf("After LINSERT: %v\n", rangeCmd.Val())
}

// TestRedisClient_LPushx 测试 LPUSHX 命令（仅当列表存在时插入）
func TestRedisClient_LPushx(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	keyName := "test12"

	// 第一次尝试（列表不存在，应该失败）
	cmd1 := client.LPushx(context.Background(), ListCmd, map[string]any{
		"keyName": keyName,
		"value":   "value1",
	})

	if cmd1.Err() != nil {
		t.Errorf("First LPushx failed: %v", cmd1.Err())
		return
	}

	fmt.Printf("First LPUSHX (list not exists) result: %d (should be 0)\n", cmd1.Val())

	// 先创建列表
	client.LPush(context.Background(), ListCmd, map[string]any{
		"keyName": keyName,
		"value":   "existing",
	})

	// 第二次尝试（列表存在，应该成功）
	cmd2 := client.LPushx(context.Background(), ListCmd, map[string]any{
		"keyName": keyName,
		"value":   "value2",
	})

	if cmd2.Err() != nil {
		t.Errorf("Second LPushx failed: %v", cmd2.Err())
		return
	}

	fmt.Printf("Second LPUSHX (list exists) result: %d\n", cmd2.Val())
}

// TestRedisClient_RPushx 测试 RPUSHX 命令（仅当列表存在时插入）
func TestRedisClient_RPushx(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	keyName := "test13"

	// 先创建列表
	client.RPush(context.Background(), ListCmd, map[string]any{
		"keyName": keyName,
		"value":   "existing",
	})

	// RPUSHX 插入
	cmd := client.RPushx(context.Background(), ListCmd, map[string]any{
		"keyName": keyName,
		"value":   "new_value",
	})

	if cmd.Err() != nil {
		t.Errorf("RPushx failed: %v", cmd.Err())
		return
	}

	fmt.Printf("RPUSHX result: %d\n", cmd.Val())
}

// TestRedisClient_RPopLPush 测试 RPOPLPUSH 命令
func TestRedisClient_RPopLPush(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	sourceKey := "test14_source"
	targetKey := "test14_target"

	// 创建源列表
	client.RPush(context.Background(), ListCmd, map[string]any{
		"keyName": sourceKey,
		"value":   "item1",
	}, "item2", "item3")

	// 创建目标列表（可以为空）
	client.RPush(context.Background(), ListCmd, map[string]any{
		"keyName": targetKey,
		"value":   "target1",
	})

	// 从源列表右侧弹出并推入目标列表左侧
	cmd := client.RPopLPush(context.Background(), ListCmd, map[string]any{
		"keyName": sourceKey,
		"target":  targetKey,
	})

	if cmd.Err() != nil {
		t.Errorf("RPopLPush failed: %v", cmd.Err())
		return
	}

	fmt.Printf("RPOPLPUSH result: %s\n", cmd.Val())

	// 查看两个列表
	sourceRange := client.LRange(context.Background(), ListCmd, map[string]any{
		"keyName": sourceKey,
		"start":   0,
		"stop":    -1,
	})
	targetRange := client.LRange(context.Background(), ListCmd, map[string]any{
		"keyName": targetKey,
		"start":   0,
		"stop":    -1,
	})
	fmt.Printf("Source list: %v\n", sourceRange.Val())
	fmt.Printf("Target list: %v\n", targetRange.Val())
}

// TestRedisClient_List_Integration 集成测试：List 操作的完整流程
func TestRedisClient_List_Integration(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	keyName := "integration_list"

	// 1. LPUSH 从左侧插入
	client.LPush(context.Background(), ListCmd, map[string]any{
		"keyName": keyName,
		"value":   "left1",
	}, "left2")

	// 2. RPUSH 从右侧插入
	client.RPush(context.Background(), ListCmd, map[string]any{
		"keyName": keyName,
		"value":   "right1",
	}, "right2")

	// 3. LLEN
	lenCmd := client.LLen(context.Background(), ListCmd, map[string]any{
		"keyName": keyName,
	})
	fmt.Printf("1. LLEN: %d\n", lenCmd.Val())

	// 4. LRANGE
	rangeCmd := client.LRange(context.Background(), ListCmd, map[string]any{
		"keyName": keyName,
		"start":   0,
		"stop":    -1,
	})
	fmt.Printf("2. LRANGE: %v\n", rangeCmd.Val())

	// 5. LINDEX
	indexCmd := client.LIndex(context.Background(), ListCmd, map[string]any{
		"keyName": keyName,
		"index":   0,
	})
	fmt.Printf("3. LINDEX(0): %s\n", indexCmd.Val())

	// 6. LPOP
	popCmd := client.LPop(context.Background(), ListCmd, map[string]any{
		"keyName": keyName,
	})
	fmt.Printf("4. LPOP: %s\n", popCmd.Val())

	// 7. 最终列表
	finalRange := client.LRange(context.Background(), ListCmd, map[string]any{
		"keyName": keyName,
		"start":   0,
		"stop":    -1,
	})
	fmt.Printf("5. Final list: %v\n", finalRange.Val())
}

