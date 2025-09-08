package model

import (
	"fmt"
	"time"
)

type CommonResp struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp int64       `json:"timestamp"`
}

// for the fast return Success result

func Success(data ...interface{}) CommonResp {
	var result interface{}

	if len(data) > 0 {
		result = data[0]
	}
	return CommonResp{
		Code:      200,
		Message:   "success",
		Data:      result,
		Timestamp: DateNow().Unix(),
	}
}

// for the fast return Failed result

func Failed(code int, message string, args ...interface{}) CommonResp {
	return CommonResp{
		Code:      code,
		Message:   fmt.Sprintf(message, args...),
		Timestamp: DateNow().Unix(),
	}
}

// for the fast return result with custom data

func BaseRes(code int, message string, data ...interface{}) CommonResp {
	var result interface{}

	if len(data) > 0 {
		result = data[0]
	}

	return CommonResp{
		Code:      code,
		Message:   message,
		Data:      result,
		Timestamp: DateNow().Unix(),
	}
}

type DeviceInfo struct {
	Key   string `json:"key"`
	Token string `json:"token"`
	Voice bool   `json:"voice"`
}

func DateNow() time.Time {
	return time.Now().UTC()
}
