/*
@author: panfengguo
@since: 2024/11/9
@desc: desc
*/
package utils

import (
	"encoding/json"
	"math/rand"
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
