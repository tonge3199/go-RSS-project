package auth

import (
	"errors"
	"net/http"
	"strings"
)

var ErrNoAuthHeaderIncluded = errors.New("no authorization header included")

// GetAPIKey 从 HTTP 请求头中解析 API Key
// 预期格式: Authorization: ApiKey {your_api_key_here}  - 此格式是 HTTP 协议的标准设计
// 参数:
//
//	headers - HTTP 请求头对象
//
// 返回值:
//
//	string - 成功时返回 API 密钥
//	error  - 错误类型包括:
//	           - ErrNoAuthHeaderIncluded (无认证头)
//	           - malformed header (格式错误)
//
// 示例:
//
//	案例1: 有效请求头
//	  headers: {"Authorization": ["ApiKey abc123"]}
//	  → 返回 "abc123", nil
//
//	案例2: 缺失认证头
//	  headers: {}
//	  → 返回 "", ErrNoAuthHeaderIncluded
//
//	案例3: 错误格式
//	  headers: {"Authorization": ["Bearer token"]}  // 错误类型前缀
//	  → 返回 "", "malformed authorization header"
//
//	案例4: 非法拆分
//	  headers: {"Authorization": ["ApiKey"]}        // 缺少密钥
//	  → 返回 "", "malformed authorization header"
func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}
	splitAuth := strings.Split(authHeader, " ")
	// if len(splitAuth) < 2 || splitAuth[0] != "ApiKey" { // origin code.
	// improve :if headers: {"Authorization": ["ApiKey abc123 abc123 abc123"]},will not be allowed.
	if len(splitAuth) != 2 || splitAuth[0] != "ApiKey" {
		return "", errors.New("malformed authorization header")
	}
	return splitAuth[1], nil
}
