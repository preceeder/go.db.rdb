package rdb

import (
	"context"
	"fmt"
	"testing"
)

// TestExecuteIntCmd_Chain 测试链式调用 ExecuteIntCmd - 直接获取 *redis.IntCmd 类型
func TestExecuteIntCmd_Chain(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	var IntCmd = RdCmd{
		Key: "int:{{keyName}}",
		CMD: map[Command]RdSubCmd{
			INCR: {
				Params: "",
			},
			ZCARD: {
				Params: "",
			},
		},
	}

	// 先设置一个值
	client.Set(context.Background(), StringCmd, map[string]any{
		"keyName": "test_int",
		"value":   "10",
	})

	// 使用链式调用 Int() - 返回的是 *redis.IntCmd，可以直接调用类型安全的方法
	cmd := client.Incr(context.Background(), IntCmd, map[string]any{
		"keyName": "test_int",
	}).Int()

	if cmd.Err() != nil {
		t.Errorf("Int() failed: %v", cmd.Err())
		return
	}

	// 直接使用类型安全的方法，无需类型断言
	val, err := cmd.Result()
	if err != nil {
		t.Errorf("Result() failed: %v", err)
		return
	}

	fmt.Printf("INCR with ExecuteIntCmd result: %d (type: %T)\n", val, cmd)
	fmt.Printf("可以直接调用 cmd.Int64(): %d\n", cmd.Val())
}

// TestExecuteStringCmd_Chain 测试链式调用 ExecuteStringCmd - 直接获取 *redis.StringCmd 类型
func TestExecuteStringCmd_Chain(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	var StringCmdDirect = RdCmd{
		Key: "string:{{keyName}}",
		CMD: map[Command]RdSubCmd{
			GET: {
				Params: "",
			},
			SET: {
				Params: "{{value}}",
			},
		},
	}

	// 先设置一个值
	client.Set(context.Background(), StringCmdDirect, map[string]any{
		"keyName": "test_string",
		"value":   "Hello World",
	})

	// 使用链式调用 String() - 返回的是 *redis.StringCmd
	cmd := client.Get(context.Background(), StringCmdDirect, map[string]any{
		"keyName": "test_string",
	}).String()

	if cmd.Err() != nil {
		t.Errorf("String() failed: %v", cmd.Err())
		return
	}

	// 直接使用类型安全的方法
	val, err := cmd.Result()
	if err != nil {
		t.Errorf("Result() failed: %v", err)
		return
	}

	fmt.Printf("GET with ExecuteStringCmd result: %s (type: %T)\n", val, cmd)
	fmt.Printf("可以直接调用 cmd.Text(): %s\n", cmd.String())
}

// TestExecuteSliceCmd_Chain 测试链式调用 ExecuteSliceCmd - 直接获取 *redis.SliceCmd 类型
func TestExecuteSliceCmd_Chain(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	var HashCmdWithSlice = RdCmd{
		Key: "hash:{{keyName}}",
		CMD: map[Command]RdSubCmd{
			HGETALL: {
				Params: "",
			},
			LRANGE: {
				Params: "{{start}} {{stop}}",
			},
		},
	}

	// 先设置一些数据
	client.HMSet(context.Background(), HashCmd, map[string]any{
		"keyName": "test_slice",
	}, "name", "Alice", "age", "30", "city", "Beijing")

	// 使用链式调用 Slice() - 返回的是 *redis.SliceCmd
	cmd := client.HGetAll(context.Background(), HashCmdWithSlice, map[string]any{
		"keyName": "test_slice",
	}).Slice()

	if cmd.Err() != nil {
		t.Errorf("Slice() failed: %v", cmd.Err())
		return
	}

	// 直接使用类型安全的方法
	values, err := cmd.Result()
	if err != nil {
		t.Errorf("Result() failed: %v", err)
		return
	}

	fmt.Printf("HGETALL with ExecuteSliceCmd result: %v (type: %T)\n", values, cmd)
	fmt.Printf("Length: %d\n", len(values))
	fmt.Printf("可以直接调用 cmd.Slice(): %v\n", cmd.Val())
}

// TestExecuteFloatCmd_Chain 测试链式调用 ExecuteFloatCmd - 直接获取 *redis.FloatCmd 类型
func TestExecuteFloatCmd_Chain(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	var FloatCmd = RdCmd{
		Key: "float:{{keyName}}",
		CMD: map[Command]RdSubCmd{
			INCRBYFLOAT: {
				Params: "{{increment}}",
			},
			ZSCORE: {
				Params: "{{member}}",
			},
		},
	}

	// 先设置一个值
	client.Set(context.Background(), StringCmd, map[string]any{
		"keyName": "test_float",
		"value":   "10.5",
	})

	// 使用链式调用 Float()
	cmd := client.IncrByFloat(context.Background(), FloatCmd, map[string]any{
		"keyName":   "test_float",
		"increment": 2.5,
	}).Float()

	if cmd.Err() != nil {
		t.Errorf("Float() failed: %v", cmd.Err())
		return
	}

	// 直接使用类型安全的方法
	val, err := cmd.Result()
	if err != nil {
		t.Errorf("Result() failed: %v", err)
		return
	}

	fmt.Printf("INCRBYFLOAT with ExecuteFloatCmd result: %f (type: %T)\n", val, cmd)
	fmt.Printf("可以直接调用 cmd.Val(): %f\n", cmd.Val())
}

// TestExecuteBoolCmd_Chain 测试链式调用 ExecuteBoolCmd - 直接获取 *redis.BoolCmd 类型
func TestExecuteBoolCmd_Chain(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	var BoolCmd = RdCmd{
		Key: "bool:{{keyName}}",
		CMD: map[Command]RdSubCmd{
			SETNX: {
				Params: "{{value}}",
			},
			HEXISTS: {
				Params: "{{field}}",
			},
		},
	}

	// 使用链式调用 Bool()
	cmd := client.SetNx(context.Background(), BoolCmd, map[string]any{
		"keyName": "test_bool",
		"value":   "test_value",
	}).Bool()

	if cmd.Err() != nil {
		t.Errorf("Bool() failed: %v", cmd.Err())
		return
	}

	// 直接使用类型安全的方法
	val, err := cmd.Result()
	if err != nil {
		t.Errorf("Result() failed: %v", err)
		return
	}

	fmt.Printf("SETNX with ExecuteBoolCmd result: %v (type: %T)\n", val, cmd)
	fmt.Printf("可以直接调用 cmd.Val(): %v\n", cmd.Val())
}

// TestExecute_DefaultBehavior 测试默认的 Execute 方法（返回 *redis.Cmd）
func TestExecute_DefaultBehavior(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	var HashCmdDefault = RdCmd{
		Key: "hash:{{keyName}}",
		CMD: map[Command]RdSubCmd{
			HGETALL: {
				Params: "",
			},
			GET: {
				Params: "",
			},
		},
	}

	// 先设置数据
	client.HMSet(context.Background(), HashCmd, map[string]any{
		"keyName": "test_default",
	}, "field1", "value1", "field2", "value2")

	// 使用链式调用 Slice() 获取具体类型
	cmd := client.HGetAll(context.Background(), HashCmdDefault, map[string]any{
		"keyName": "test_default",
	}).Slice()

	if cmd.Err() != nil {
		t.Errorf("Slice() failed: %v", cmd.Err())
		return
	}

	result, _ := cmd.Result()
	fmt.Printf("HGETALL with Slice() result type: %T, value: %v\n", result, result)
}

// TestExecute_Comparison 对比不同 Execute 方法的使用方式
func TestExecute_Comparison(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// 准备数据
	client.HMSet(context.Background(), HashCmd, map[string]any{
		"keyName": "test_comparison",
	}, "name", "Bob", "age", "25")

	// 方式1: 使用 Slice() 获取原始数组
	cmd1 := client.HGetAll(context.Background(), HashCmd, map[string]any{
		"keyName": "test_comparison",
	}).Slice()
	result1, _ := cmd1.Result()
	fmt.Printf("1. Slice() result type: %T, value: %v\n", result1, result1)

	// 方式2: 使用 Slice() 获取原始数组
	var HashCmdSlice = RdCmd{
		Key: "hash:{{keyName}}",
		CMD: map[Command]RdSubCmd{
			HGETALL: {
				Params: "",
			},
		},
	}
	cmd2 := client.HGetAll(context.Background(), HashCmdSlice, map[string]any{
		"keyName": "test_comparison",
	}).Slice()
	result2, _ := cmd2.Result()
	fmt.Printf("2. Slice() result type: %T, value: %v\n", result2, result2)

	// 方式3: 使用 Int() 获取整数
	var IntCmd = RdCmd{
		Key: "int:{{keyName}}",
		CMD: map[Command]RdSubCmd{
			INCR: {
				Params: "",
			},
		},
	}
	client.Set(context.Background(), StringCmd, map[string]any{
		"keyName": "test_comparison_int",
		"value":   "100",
	})
	cmd3 := client.Incr(context.Background(), IntCmd, map[string]any{
		"keyName": "test_comparison_int",
	}).Int()
	result3, _ := cmd3.Result()
	fmt.Printf("3. Int() result type: %T, value: %d\n", result3, result3)
}

// TestExecute_TypeSafety 演示类型安全的使用方式
func TestExecute_TypeSafety(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	var IntCmd = RdCmd{
		Key: "int:{{keyName}}",
		CMD: map[Command]RdSubCmd{
			INCR: {
				Params: "",
			},
		},
	}

	client.Set(context.Background(), StringCmd, map[string]any{
		"keyName": "test_type_safety",
		"value":   "10",
	})

	// 使用链式调用，直接获取类型安全的 *redis.IntCmd
	cmd := client.Incr(context.Background(), IntCmd, map[string]any{
		"keyName": "test_type_safety",
	}).Int()

	if cmd.Err() != nil {
		t.Errorf("Int() failed: %v", cmd.Err())
		return
	}

	// 由于返回的是 *redis.IntCmd，我们可以直接使用类型安全的方法
	val, _ := cmd.Result()
	fmt.Printf("类型安全，直接使用: %d\n", val)
	fmt.Printf("或者直接调用: %d\n", cmd.Val())
	fmt.Printf("返回类型: %T\n", cmd)
}

// TestExecute_MultipleCommands 测试多个命令使用不同的 Execute 方法
func TestExecute_MultipleCommands(t *testing.T) {
	client := InitRedis()
	defer client.RedisClose()

	// 定义命令（不需要指定 NewCmder）
	var MultiCmd = RdCmd{
		Key: "multi:{{keyName}}",
		CMD: map[Command]RdSubCmd{
			SET: {
				Params: "{{value}}",
			},
			GET: {
				Params: "",
			},
			INCR: {
				Params: "",
			},
			HGETALL: {
				Params: "",
			},
		},
	}

	// 测试 SET - 使用 String()
	setCmd := client.Set(context.Background(), MultiCmd, map[string]any{
		"keyName": "test_multi",
		"value":   "test_value",
	}).String()
	if setCmd.Err() != nil {
		t.Errorf("SET String() failed: %v", setCmd.Err())
		return
	}
	fmt.Printf("SET returns: %T\n", setCmd)

	// 测试 GET - 使用 String()
	getCmd := client.Get(context.Background(), MultiCmd, map[string]any{
		"keyName": "test_multi",
	}).String()
	if getCmd.Err() != nil {
		t.Errorf("GET String() failed: %v", getCmd.Err())
		return
	}
	val, _ := getCmd.Result()
	fmt.Printf("GET returns: %T, value: %s\n", getCmd, val)

	// 测试 INCR - 使用 Int()
	incrCmd := client.Incr(context.Background(), MultiCmd, map[string]any{
		"keyName": "test_multi",
	}).Int()
	if incrCmd.Err() != nil {
		t.Errorf("INCR Int() failed: %v", incrCmd.Err())
		return
	}
	intVal, _ := incrCmd.Result()
	fmt.Printf("INCR returns: %T, value: %d\n", incrCmd, intVal)

	// 测试 HGETALL - 使用 Slice()
	client.HMSet(context.Background(), HashCmd, map[string]any{
		"keyName": "test_multi",
	}, "field1", "value1")
	hgetallCmd := client.HGetAll(context.Background(), MultiCmd, map[string]any{
		"keyName": "test_multi",
	}).Slice()
	if hgetallCmd.Err() != nil {
		t.Errorf("HGETALL Slice() failed: %v", hgetallCmd.Err())
		return
	}
	sliceVal, _ := hgetallCmd.Result()
	fmt.Printf("HGETALL returns: %T, value: %v\n", hgetallCmd, sliceVal)
}
