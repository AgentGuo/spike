/*
@author: panfengguo
@since: 2024/11/9
@desc: desc
*/
package utils

import (
	"encoding/json"
	"fmt"
	"github.com/sony/sonyflake"
	"math/rand"
	"sync"
	"time"
)

func GetJson(data interface{}) string {
	output, _ := json.Marshal(data)
	return string(output)
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

// SonyFlakeSingleton 定义一个结构体用于封装 SonyFlake 实例
type SonyFlakeSingleton struct {
	flake *sonyflake.Sonyflake
}

var (
	instance *SonyFlakeSingleton // 单例实例
	once     sync.Once           // 确保只初始化一次
)

// GetSonyFlakeInstance 获取 SonyFlake 的单例实例
func GetSonyFlakeInstance() *SonyFlakeSingleton {
	once.Do(func() {
		sf := sonyflake.NewSonyflake(sonyflake.Settings{})
		if sf == nil {
			panic("Failed to initialize Sonyflake!")
		}
		instance = &SonyFlakeSingleton{flake: sf}
	})
	return instance
}

// GenerateID 使用单例生成唯一 ID
func (s *SonyFlakeSingleton) GenerateID() uint64 {
	id, err := s.flake.NextID()
	if err != nil {
		panic(fmt.Sprintf("Failed to generate ID: %v", err))
	}
	return id
}
