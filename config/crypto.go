package config

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"

	"github.com/bytedance/sonic"
)

// RequestMethod 请求方法类型
type RequestMethod string

const (
	MethodGet  RequestMethod = "GET"
	MethodPost RequestMethod = "POST"
)

// flattenParams 扁平化嵌套参数（与Swift版本保持一致）
func flattenParams(params map[string]interface{}, prefix string) map[string]string {
	result := make(map[string]string)

	for key, value := range params {
		newKey := key
		if prefix != "" {
			newKey = prefix + "." + key
		}

		switch v := value.(type) {
		case map[string]interface{}:
			// 递归处理子字典
			subResult := flattenParams(v, newKey)
			for k, val := range subResult {
				result[k] = val
			}
		case []interface{}:
			// 处理数组
			for index, item := range v {
				arrayKey := fmt.Sprintf("%s[%d]", newKey, index)
				if subDict, ok := item.(map[string]interface{}); ok {
					subResult := flattenParams(subDict, arrayKey)
					for k, val := range subResult {
						result[k] = val
					}
				} else {
					result[arrayKey] = fmt.Sprintf("%v", item)
				}
			}
		default:
			result[newKey] = fmt.Sprintf("%v", value)
		}
	}

	return result
}

// VerifySignature 校验签名
func VerifySignature(url string, method RequestMethod, params map[string]interface{},
	signature, secretKey string) (bool, error) {

	// 构建与Swift相同的参数结构
	baseParams := map[string]interface{}{
		"url":    url,
		"method": string(method),
	}

	// 合并参数
	allParams := make(map[string]interface{})
	for k, v := range baseParams {
		allParams[k] = v
	}
	for k, v := range params {
		allParams[k] = v
	}

	// 扁平化参数
	flatParams := flattenParams(allParams, "")

	// 按key升序排序并构建参数字符串
	keys := make([]string, 0, len(flatParams))
	for k := range flatParams {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var paramStrings []string
	for _, key := range keys {
		paramStrings = append(paramStrings, fmt.Sprintf("%s:%s", key, flatParams[key]))
	}
	paramsStr := strings.Join(paramStrings, ",")

	// 计算HMAC-SHA256
	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(paramsStr))
	expectedSignature := hex.EncodeToString(mac.Sum(nil))

	// 比较签名（不区分大小写）
	return strings.EqualFold(expectedSignature, signature), nil
}

// VerifySignatureFromJSON 从JSON字符串校验签名
func VerifySignatureFromJSON(jsonData, signature, secretKey string) (bool, error) {
	var request struct {
		URL    string                 `json:"url"`
		Method string                 `json:"method"`
		Params map[string]interface{} `json:"params"`
	}

	if err := sonic.Unmarshal([]byte(jsonData), &request); err != nil {
		return false, fmt.Errorf("解析JSON失败: %v", err)
	}

	return VerifySignature(
		request.URL,
		RequestMethod(request.Method),
		request.Params,
		signature,
		secretKey,
	)
}
