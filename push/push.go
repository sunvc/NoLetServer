package push

import (
	"NoLetServer/config"
	"NoLetServer/model"
	"fmt"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"github.com/uuneo/apns2"
	"github.com/uuneo/apns2/payload"
)

// Push message to APNs server
func Push(params *model.ParamsMap, pushType apns2.EPushType, token string) error {
	pl := payload.NewPayload().MutableContent()

	if pushType == apns2.PushTypeBackground {
		pl = pl.ContentAvailable()
	} else {
		pl = pl.AlertTitle(model.PMGet(params, model.Title)).
			AlertSubtitle(model.PMGet(params, model.Subtitle)).
			AlertBody(model.PMGet(params, model.Body)).
			Sound(model.PMGet(params, model.Sound)).
			TargetContentID(model.PMGet(params, model.ID)).
			ThreadID(model.PMGet(params, model.Group)).
			Category(model.PMGet(params, model.Category))
	}

	// 添加自定义参数
	skipKeys := map[string]struct{}{
		model.DeviceKey:   {},
		model.DeviceToken: {},
		model.Title:       {},
		model.Body:        {},
		model.Sound:       {},
		model.Category:    {},
	}

	for pair := params.Oldest(); pair != nil; pair = pair.Next() {
		if _, skip := skipKeys[pair.Key]; skip {
			continue
		}
		pl.Custom(pair.Key, pair.Value)
	}

	CLI := <-CLIENTS // 从池中获取一个客户端
	CLIENTS <- CLI   // 将客户端放回池中

	// 创建并发送通知
	resp, err := CLI.Push(&apns2.Notification{
		DeviceToken: token,
		CollapseID:  fmt.Sprint(params.Value(model.ID)),
		Topic:       config.LocalConfig.Apple.Topic,
		Payload:     pl,
		Expiration:  model.DateNow().Add(24 * time.Hour),
		PushType:    pushType,
	})

	// 错误处理
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("APNs push failed: %s", resp.Reason)
	}
	return nil

}

func BatchPush(params *model.ParamsResult, pushType apns2.EPushType) error {
	var (
		errors []error
		mu     sync.Mutex
		wg     sync.WaitGroup
	)
	if len(params.Results) > 0 {

		for _, param := range params.Results {

			for _, key := range params.DeviceTokens {
				wg.Add(1)
				go func(p *model.ParamsMap) {
					defer wg.Done()
					if err := Push(p, pushType, key); err != nil {
						log.Error(err.Error())
						mu.Lock()
						errors = append(errors, err)
						mu.Unlock()
					}
				}(param)
			}

		}

	} else {

		// 如果 title, subtitle 和 body 都为空，设置无推送模式
		if params.IsNan {
			pushType = apns2.PushTypeBackground
		}
		for _, key := range params.DeviceTokens {
			wg.Add(1)
			go func() {
				if err := Push(params.Params, pushType, key); err != nil {
					errors = append(errors, err)
				}
			}()

		}

	}

	wg.Wait()
	if len(errors) > 0 {
		return fmt.Errorf("APNs push failed: %v", errors)
	}

	return nil
}

func PTTPush(params map[string]interface{}, token string) error {
	pl := payload.NewPayload()

	for key, value := range params {
		if key == model.DeviceKey || key == model.DeviceToken {
			continue // 跳过设备相关的键
		}
		pl.Custom(key, value)
	}
	CLI := <-CLIENTS // 从池中获取一个客户端
	CLIENTS <- CLI   // 将客户端放回池中
	// 创建并发送通知
	resp, err := CLI.Push(&apns2.Notification{
		DeviceToken: token,
		Topic:       config.LocalConfig.Apple.Topic + ".voip-ptt",
		Payload:     pl,
		Expiration:  model.DateNow().Add(24 * time.Hour),
		PushType:    apns2.PushTypePushToTalk,
	})
	// 错误处理
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("APNs push failed: %s", resp.Reason)
	}
	return nil
}
