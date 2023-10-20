// Package redis 工具包
package redis

import (
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
	"os"
	"quick/conf"
	"quick/pkg/log"
	"sync"
	"time"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// RedisClient Redis 服务
type RedisClient struct {
	Client  *redis.Client
	Context context.Context
}

// once 确保全局的 Redis 对象只实例一次
var once sync.Once

// Redis 全局 Redis，使用 db 1
var Redis *RedisClient

// ConnectRedis 连接 redis 数据库，设置全局的 Redis 对象
func ConnectRedis(address string, username string, password string, db int) {
	once.Do(func() {
		Redis = NewClient(conf.RedisConfig)
	})
}

// NewClient 创建一个新的 redis 连接
func NewClient(conf *conf.Redis) *RedisClient {

	// 初始化自定的 RedisClient 实例
	rds := &RedisClient{}
	// 使用默认的 context
	rds.Context = context.Background()

	// 使用 redis 库里的 NewClient 初始化连接
	rds.Client = redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Password: conf.Password,
		DB:       conf.Db,
	})

	// 测试一下连接
	if err := rds.Ping(); err != nil {
		log.Sugar.Errorf("redis连接错误%s", err)
		os.Exit(1)
	}

	return rds
}

// Ping 用以测试 redis 连接是否正常
func (rds RedisClient) Ping() error {
	_, err := rds.Client.Ping(rds.Context).Result()
	return err
}
func (rds RedisClient) GetDeviceKey(iccid string) string {
	return fmt.Sprintf("device_info:%s", iccid)

}
func (rds RedisClient) GetSlaveKey(iccid string) string {
	return fmt.Sprintf("slave_info:%s", iccid)

}
func (rds RedisClient) GetPropertyKey() string {
	return fmt.Sprintf("property_info")

}
func (rds RedisClient) GetAwaitSendKey(iccid, slaveId string, markType string, groupId int) string {
	return fmt.Sprintf("await_send:%v:%v:%s:%s", markType, groupId, iccid, slaveId)

}

// Set 存储 key 对应的 value，且设置 expiration 过期时间
func (rds RedisClient) Set(key string, value interface{}, expiration time.Duration) bool {
	if err := rds.Client.Set(rds.Context, key, value, expiration).Err(); err != nil {
		log.Sugar.Warnf("Redis", "Set", err.Error())
		return false
	}
	return true
}

// Get 获取 key 对应的 value
func (rds RedisClient) Get(key string, model interface{}) error {
	result, err := rds.Client.Get(rds.Context, key).Result()
	if err != nil {
		log.Sugar.Errorf("Redis", "Get", err.Error())
		return err
	}
	err = json.Unmarshal([]byte(result), &model)
	return nil
}

// Has 判断一个 key 是否存在，内部错误和 redis.Nil 都返回 false
func (rds RedisClient) Has(key string) bool {
	_, err := rds.Client.Get(rds.Context, key).Result()
	if err != nil {
		if err != redis.Nil {
			log.Sugar.Errorf("Redis", "Has", err.Error())
		}
		return false
	}
	return true
}

// Del 删除存储在 redis 里的数据，支持多个 key 传参
func (rds RedisClient) Del(keys ...string) bool {
	if err := rds.Client.Del(rds.Context, keys...).Err(); err != nil {
		log.Sugar.Errorf("Redis", "Del", err.Error())
		return false
	}
	return true
}

// FlushDB 清空当前 redis db 里的所有数据
func (rds RedisClient) FlushDB() bool {
	if err := rds.Client.FlushDB(rds.Context).Err(); err != nil {
		log.Sugar.Errorf("Redis", "FlushDB", err.Error())
		return false
	}

	return true
}

// Increment 当参数只有 1 个时，为 key，其值增加 1。
// 当参数有 2 个时，第一个参数为 key ，第二个参数为要增加的值 int64 类型。
func (rds RedisClient) Increment(parameters ...interface{}) bool {
	switch len(parameters) {
	case 1:
		key := parameters[0].(string)
		if err := rds.Client.Incr(rds.Context, key).Err(); err != nil {
			log.Sugar.Errorf("Redis", "Increment", err.Error())
			return false
		}
	case 2:
		key := parameters[0].(string)
		value := parameters[1].(int64)
		if err := rds.Client.IncrBy(rds.Context, key, value).Err(); err != nil {
			log.Sugar.Errorf("Redis", "Increment", err.Error())
			return false
		}
	default:
		log.Sugar.Errorf("Redis", "Increment", "参数过多")
		return false
	}
	return true
}

// Decrement 当参数只有 1 个时，为 key，其值减去 1。
// 当参数有 2 个时，第一个参数为 key ，第二个参数为要减去的值 int64 类型。
func (rds RedisClient) Decrement(parameters ...interface{}) bool {
	switch len(parameters) {
	case 1:
		key := parameters[0].(string)
		if err := rds.Client.Decr(rds.Context, key).Err(); err != nil {
			log.Sugar.Errorf("Redis", "Decrement", err.Error())
			return false
		}
	case 2:
		key := parameters[0].(string)
		value := parameters[1].(int64)
		if err := rds.Client.DecrBy(rds.Context, key, value).Err(); err != nil {
			log.Sugar.Errorf("Redis", "Decrement", err.Error())
			return false
		}
	default:
		log.Sugar.Errorf("Redis", "Decrement", "参数过多")
		return false
	}
	return true
}

// Rpush 写入队列
func (rds RedisClient) Rpush(key string, value interface{}) bool {
	if err := rds.Client.RPush(rds.Context, key, value).Err(); err != nil {
		fmt.Printf("Redis   Rpush error:%s", err.Error())
		return false
	}
	return true
}

// Lpop Loop 出列
func (rds RedisClient) Lpop(key string) string {
	result, err := rds.Client.LPop(rds.Context, key).Result()
	if err != nil {
		if err != redis.Nil {
			fmt.Printf("Redis,Lpop,Error:%s", err.Error())
		}
		return ""
	}
	return result
}
func (rds RedisClient) HGetAll(key string, model interface{}) error {
	if err := rds.Client.HGetAll(rds.Context, key).Scan(&model); err != nil {
		log.Sugar.Warnf("Redis", "hash get", err)
		return err
	}
	return nil
}
func (rds RedisClient) HGet(key, field string, model interface{}) error {
	re, err := rds.Client.HGet(rds.Context, key, field).Result()
	if err != nil {
		log.Sugar.Warnf("Redis", "hash get", err)
		return err
	}
	err = json.Unmarshal([]byte(re), &model)
	if err != nil {
		return err
	}
	return nil
}
func (rds RedisClient) HSet(key, field string, model interface{}) {
	marshal, err := json.Marshal(model)
	if err != nil {
		return
	}
	_, err = rds.Client.HSet(rds.Context, key, field, marshal).Result()
	if err != nil {
		return
	}

}
func (rds RedisClient) HMSet(key string, model interface{}) {

	sd := cast.ToStringMap(&model)
	result, err := rds.Client.HMSet(rds.Context, key, sd).Result()
	if err != nil {
		return
	}

	fmt.Println(result)
}
