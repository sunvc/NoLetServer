package controller

import (
	"NoLetServer/model"
	"NoLetServer/push"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/sunvc/apns2"
)

// MARK: - 推送任务

var NotPushedDataList sync.Map

func CirclePush() {
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			log.Info("开始检查未推送数据")
			NotPushedDataList.Range(func(key, value any) bool {
				data1, ok := value.(*model.NotPushedData)
				if !ok {
					NotPushedDataList.Delete(key) // 类型异常也清除
					return true
				}

				now := model.DateNow()

				// 超过 24 小时未成功推送，直接清除
				if now.Sub(data1.LastPushDate) > 24*time.Hour {
					NotPushedDataList.Delete(key)
					return true
				}

				// 推送节流策略：每次失败后等待 Count × 10 分钟
				nextTry := data1.LastPushDate.Add(time.Duration(data1.Count) * 10 * time.Minute)
				if nextTry.After(now) {
					return true // 还没到下一次推送时间，跳过
				}

				// 执行推送
				if err := push.BatchPush(data1.Params, data1.PushType); err != nil {
					NotPushedDataList.Delete(key) // 推送失败直接删
				}

				return true
			})
		}
	}()

}

// UpdateNotPushedData 更新已有记录，若不存在则添加
func UpdateNotPushedData(id string, params *model.ParamsResult, pushType apns2.EPushType) {
	if val, ok := NotPushedDataList.Load(id); ok {
		res := val.(*model.NotPushedData)
		res.LastPushDate = model.DateNow()
		res.Count++
		res.Params = params
		res.PushType = pushType
		NotPushedDataList.Store(id, model.Success) // 可省略，但保持一致性
	} else {
		NotPushedDataList.Store(id, &model.NotPushedData{
			ID:           id,
			CreateDate:   model.DateNow(),
			LastPushDate: model.DateNow(),
			Count:        1,
			Params:       params,
			PushType:     pushType,
		})
	}
}

// RemoveNotPushedData 删除数据
func RemoveNotPushedData(id string) {
	NotPushedDataList.Delete(id)
}
