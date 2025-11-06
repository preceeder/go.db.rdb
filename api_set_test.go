package rdb

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// 测试 Set 操作的 RdCmd 定义
var SetCmd = RdCmd{
	Key: "set:{{keyName}}",
	CMD: map[Command]RdSubCmd{
		SADD: {
			Params: "{{member}}",
			Exp: func() time.Duration {
				return time.Second * 30
			},
		},
		SREM: {
			Params: "{{member}}",
		},
		SMEMBERS: {
			Params: "",
		},
		SCARD: {
			Params: "",
		},
		SISMEMBER: {
			Params: "{{member}}",
		},
		SMOVE: {
			Params: "{{destination}} {{member}}",
		},
		SINTER: {
			Params:   "",
			NoUseKey: true,
		},
		SUNION: {
			Params:   "",
			NoUseKey: true,
		},
		SDIFF: {
			Params:   "",
			NoUseKey: true,
		},
		SINTERSTORE: {
			Params:   "",
			NoUseKey: true,
		},
		SUNIONSTORE: {
			Params:   "",
			NoUseKey: true,
		},
		SDIFFSTORE: {
			Params:   "",
			NoUseKey: true,
		},
	},
}

// TestRedisClient_SAdd 测试 SADD 命令
func TestRedisClient_SAdd(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// SADD 添加多个成员
	cmd := client.SAdd(context.Background(), SetCmd, map[string]any{
		"keyName": "test1",
		"member":  "member1",
	}, "member2", "member3", "member1") // member1 重复，会被忽略

	if cmd.Err() != nil {
		t.Errorf("SAdd failed: %v", cmd.Err())
		return
	}

	fmt.Printf("SADD result: %d members added\n", cmd.Val())
}

// TestRedisClient_SRem 测试 SREM 命令
func TestRedisClient_SRem(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// 先添加成员
	client.SAdd(context.Background(), SetCmd, map[string]any{
		"keyName": "test2",
		"member":  "member1",
	}, "member2", "member3", "member4")

	// 删除成员
	cmd := client.SRem(context.Background(), SetCmd, map[string]any{
		"keyName": "test2",
		"member":  "member1",
	}, "member2") // 可以删除多个

	if cmd.Err() != nil {
		t.Errorf("SRem failed: %v", cmd.Err())
		return
	}

	fmt.Printf("SREM result: %d members removed\n", cmd.Val())
}

// TestRedisClient_SMembers 测试 SMEMBERS 命令
func TestRedisClient_SMembers(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// 先添加成员
	client.SAdd(context.Background(), SetCmd, map[string]any{
		"keyName": "test3",
		"member":  "apple",
	}, "banana", "cherry", "date")

	// 获取所有成员
	cmd := client.SMembers(context.Background(), SetCmd, map[string]any{
		"keyName": "test3",
	})

	if cmd.Err() != nil {
		t.Errorf("SMembers failed: %v", cmd.Err())
		return
	}

	members := cmd.Val()
	fmt.Printf("SMEMBERS result: %v\n", members)
}

// TestRedisClient_SCard 测试 SCARD 命令
func TestRedisClient_SCard(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// 先添加成员
	client.SAdd(context.Background(), SetCmd, map[string]any{
		"keyName": "test4",
		"member":  "a",
	}, "b", "c", "d", "e")

	// 获取集合大小
	cmd := client.SCard(context.Background(), SetCmd, map[string]any{
		"keyName": "test4",
	})

	if cmd.Err() != nil {
		t.Errorf("SCard failed: %v", cmd.Err())
		return
	}

	fmt.Printf("SCARD result: %d\n", cmd.Val())
}

// TestRedisClient_SIsMember 测试 SISMEMBER 命令
func TestRedisClient_SIsMember(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// 先添加成员
	client.SAdd(context.Background(), SetCmd, map[string]any{
		"keyName": "test5",
		"member":  "apple",
	}, "banana", "cherry")

	// 检查成员是否存在
	cmd1 := client.SIsMember(context.Background(), SetCmd, map[string]any{
		"keyName": "test5",
		"member":  "apple",
	})

	if cmd1.Err() != nil {
		t.Errorf("SIsMember failed: %v", cmd1.Err())
		return
	}

	fmt.Printf("SISMEMBER apple: %d (should be 1)\n", cmd1.Val())

	// 检查不存在的成员
	cmd2 := client.SIsMember(context.Background(), SetCmd, map[string]any{
		"keyName": "test5",
		"member":  "orange",
	})

	if cmd2.Err() != nil {
		t.Errorf("SIsMember failed: %v", cmd2.Err())
		return
	}

	fmt.Printf("SISMEMBER orange: %d (should be 0)\n", cmd2.Val())
}

// TestRedisClient_SMove 测试 SMOVE 命令
func TestRedisClient_SMove(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	sourceKey := "test6_source"
	destKey := "test6_dest"

	// 创建源集合
	client.SAdd(context.Background(), SetCmd, map[string]any{
		"keyName": sourceKey,
		"member":  "member1",
	}, "member2", "member3")

	// 创建目标集合
	client.SAdd(context.Background(), SetCmd, map[string]any{
		"keyName": destKey,
		"member":  "existing",
	})

	// 移动成员
	cmd := client.SMove(context.Background(), SetCmd, map[string]any{
		"keyName":    sourceKey,
		"destination": destKey,
		"member":     "member1",
	})

	if cmd.Err() != nil {
		t.Errorf("SMove failed: %v", cmd.Err())
		return
	}

	fmt.Printf("SMOVE result: %d (1=success, 0=failed)\n", cmd.Val())

	// 查看两个集合
	sourceMembers := client.SMembers(context.Background(), SetCmd, map[string]any{
		"keyName": sourceKey,
	})
	destMembers := client.SMembers(context.Background(), SetCmd, map[string]any{
		"keyName": destKey,
	})
	fmt.Printf("Source set: %v\n", sourceMembers.Val())
	fmt.Printf("Destination set: %v\n", destMembers.Val())
}

// TestRedisClient_SInter 测试 SINTER 命令（交集）
func TestRedisClient_SInter(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// 创建多个集合
	client.SAdd(context.Background(), SetCmd, map[string]any{
		"keyName": "set1",
		"member":  "a",
	}, "b", "c", "d")

	client.SAdd(context.Background(), SetCmd, map[string]any{
		"keyName": "set2",
		"member":  "c",
	}, "d", "e", "f")

	client.SAdd(context.Background(), SetCmd, map[string]any{
		"keyName": "set3",
		"member":  "c",
	}, "d", "g")

	// 计算交集（使用 NoUseKey，通过 includeArgs 传递多个 key）
	cmd := client.SInter(context.Background(), SetCmd, nil,
		"set:set1",
		"set:set2",
		"set:set3",
	)

	if cmd.Err() != nil {
		t.Errorf("SInter failed: %v", cmd.Err())
		return
	}

	intersection := cmd.Val()
	fmt.Printf("SINTER result: %v\n", intersection)
}

// TestRedisClient_SUnion 测试 SUNION 命令（并集）
func TestRedisClient_SUnion(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// 创建多个集合
	client.SAdd(context.Background(), SetCmd, map[string]any{
		"keyName": "union1",
		"member":  "a",
	}, "b", "c")

	client.SAdd(context.Background(), SetCmd, map[string]any{
		"keyName": "union2",
		"member":  "c",
	}, "d", "e")

	// 计算并集
	cmd := client.SUnion(context.Background(), SetCmd, nil,
		"set:union1",
		"set:union2",
	)

	if cmd.Err() != nil {
		t.Errorf("SUnion failed: %v", cmd.Err())
		return
	}

	union := cmd.Val()
	fmt.Printf("SUNION result: %v\n", union)
}

// TestRedisClient_SDiff 测试 SDIFF 命令（差集）
func TestRedisClient_SDiff(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// 创建多个集合
	client.SAdd(context.Background(), SetCmd, map[string]any{
		"keyName": "diff1",
		"member":  "a",
	}, "b", "c", "d")

	client.SAdd(context.Background(), SetCmd, map[string]any{
		"keyName": "diff2",
		"member":  "c",
	}, "d")

	// 计算差集（第一个集合独有的元素）
	cmd := client.SDiff(context.Background(), SetCmd, nil,
		"set:diff1",
		"set:diff2",
	)

	if cmd.Err() != nil {
		t.Errorf("SDiff failed: %v", cmd.Err())
		return
	}

	diff := cmd.Val()
	fmt.Printf("SDIFF result: %v\n", diff)
}

// TestRedisClient_SInterStore 测试 SINTERSTORE 命令
func TestRedisClient_SInterStore(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// 创建源集合
	client.SAdd(context.Background(), SetCmd, map[string]any{
		"keyName": "inter1",
		"member":  "a",
	}, "b", "c")

	client.SAdd(context.Background(), SetCmd, map[string]any{
		"keyName": "inter2",
		"member":  "c",
	}, "d", "e")

	// 计算交集并存储到目标集合
	cmd := client.SInterStore(context.Background(), SetCmd, nil,
		"set:inter_result",
		"set:inter1",
		"set:inter2",
	)

	if cmd.Err() != nil {
		t.Errorf("SInterStore failed: %v", cmd.Err())
		return
	}

	fmt.Printf("SINTERSTORE result: %d elements stored\n", cmd.Val())

	// 查看结果集合
	result := client.SMembers(context.Background(), SetCmd, map[string]any{
		"keyName": "inter_result",
	})
	fmt.Printf("Stored intersection: %v\n", result.Val())
}

// TestRedisClient_SUnionStore 测试 SUNIONSTORE 命令
func TestRedisClient_SUnionStore(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// 创建源集合
	client.SAdd(context.Background(), SetCmd, map[string]any{
		"keyName": "union_store1",
		"member":  "a",
	}, "b")

	client.SAdd(context.Background(), SetCmd, map[string]any{
		"keyName": "union_store2",
		"member":  "b",
	}, "c")

	// 计算并集并存储
	cmd := client.SUnionStore(context.Background(), SetCmd, nil,
		"set:union_result",
		"set:union_store1",
		"set:union_store2",
	)

	if cmd.Err() != nil {
		t.Errorf("SUnionStore failed: %v", cmd.Err())
		return
	}

	fmt.Printf("SUNIONSTORE result: %d elements stored\n", cmd.Val())

	// 查看结果集合
	result := client.SMembers(context.Background(), SetCmd, map[string]any{
		"keyName": "union_result",
	})
	fmt.Printf("Stored union: %v\n", result.Val())
}

// TestRedisClient_SDiffStore 测试 SDIFFSTORE 命令
func TestRedisClient_SDiffStore(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// 创建源集合
	client.SAdd(context.Background(), SetCmd, map[string]any{
		"keyName": "diff_store1",
		"member":  "a",
	}, "b", "c")

	client.SAdd(context.Background(), SetCmd, map[string]any{
		"keyName": "diff_store2",
		"member":  "c",
	}, "d")

	// 计算差集并存储
	cmd := client.SDiffStore(context.Background(), SetCmd, nil,
		"set:diff_result",
		"set:diff_store1",
		"set:diff_store2",
	)

	if cmd.Err() != nil {
		t.Errorf("SDiffStore failed: %v", cmd.Err())
		return
	}

	fmt.Printf("SDIFFSTORE result: %d elements stored\n", cmd.Val())

	// 查看结果集合
	result := client.SMembers(context.Background(), SetCmd, map[string]any{
		"keyName": "diff_result",
	})
	fmt.Printf("Stored diff: %v\n", result.Val())
}

// TestRedisClient_Set_Integration 集成测试：Set 操作的完整流程
func TestRedisClient_Set_Integration(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	keyName := "integration_set"

	// 1. SADD 添加成员
	client.SAdd(context.Background(), SetCmd, map[string]any{
		"keyName": keyName,
		"member":  "member1",
	}, "member2", "member3", "member4")

	// 2. SCARD
	cardCmd := client.SCard(context.Background(), SetCmd, map[string]any{
		"keyName": keyName,
	})
	fmt.Printf("1. SCARD: %d\n", cardCmd.Val())

	// 3. SMEMBERS
	membersCmd := client.SMembers(context.Background(), SetCmd, map[string]any{
		"keyName": keyName,
	})
	fmt.Printf("2. SMEMBERS: %v\n", membersCmd.Val())

	// 4. SISMEMBER
	isMemberCmd := client.SIsMember(context.Background(), SetCmd, map[string]any{
		"keyName": keyName,
		"member":  "member1",
	})
	fmt.Printf("3. SISMEMBER member1: %d\n", isMemberCmd.Val())

	// 5. SREM
	remCmd := client.SRem(context.Background(), SetCmd, map[string]any{
		"keyName": keyName,
		"member":  "member4",
	})
	fmt.Printf("4. SREM member4: %d\n", remCmd.Val())

	// 6. 最终集合
	finalMembers := client.SMembers(context.Background(), SetCmd, map[string]any{
		"keyName": keyName,
	})
	fmt.Printf("5. Final set: %v\n", finalMembers.Val())
}

