package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// respondWithError 返回一个标准化的 JSON 格式错误响应
//
// 参数:
//   - w     : HTTP 响应写入器
//   - code  : HTTP 状态码（如 400 Bad Request, 500 Internal Server Error）
//   - msg   : 人类可读的错误描述信息
//
// 行为:
//   - 自动将错误信息包装为 { "error": "message" } 的 JSON 格式
//   - 如果状态码为 5XX 错误 (code > 499)，会记录错误日志
//   - 调用 respondWithJSON 最终写入响应
//
// 示例:
//
//	respondWithError(w, http.StatusNotFound, "User not found")
//	respondWithError(w, http.StatusInternalServerError, "Database connection failed")
func respondWithError(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5XX errors: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})
}

// respondWithJSON 将数据序列化为 JSON 并发送 HTTP 响应。
//
// 功能：
//   - 设置响应头的 Content-Type 为 "application/json"
//   - 将 payload 参数序列化为 JSON 格式
//   - 若序列化失败，返回 500 状态码并记录错误日志
//   - 若序列化成功，写入指定的 HTTP 状态码和 JSON 数据
//
// 参数：
//   - w       : HTTP 响应写入器
//   - code    : HTTP 状态码（如 200, 400, 500）
//   - payload : 需要序列化的数据（可以是结构体、map、切片等）
//
// 注意：
//   - 若 payload 为 nil，会序列化为 JSON 的 "null"
//   - 若需返回空 JSON 对象，建议传递 struct{}{}
//   - 调用此函数后，不应再对 w 进行写入操作
//
// 示例：
//
//	respondWithJSON(w, http.StatusOK, map[string]string{"status": "ok"})
//	respondWithJSON(w, http.StatusBadRequest, errors.New("invalid request"))
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError) // statusCode : 500
		return
	}
	w.WriteHeader(code)
	_, err = w.Write(data)
	if err != nil {
		return // 可能导致错误???
	}
}
