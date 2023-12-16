package utils

import "github.com/google/uuid"

func UUID() string {
	// 生成一个随机的UUID作为请求ID
	requestID := uuid.New().String()
	return requestID
}
