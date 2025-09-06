package controller

import (
	"NoLetServer/config"
	"NoLetServer/model"
	"NoLetServer/push"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

var VoiceList = make(chan model.PttVoice, 1000)

// 在线用户列表：用 map 模拟集合（userID -> true）
var onlineUsers = struct {
	sync.RWMutex
	Users map[string]model.PttUser
}{Users: make(map[string]model.PttUser)}

// 群组字典：群名 -> 用户集合
var groups = struct {
	sync.RWMutex
	Data map[string]map[string]model.PttUser
}{Data: make(map[string]map[string]model.PttUser)}

func JoinPTTChannel(user model.PttUser) {
	onlineUsers.Lock()
	defer onlineUsers.Unlock()

	if _, exists := onlineUsers.Users[user.ID]; !exists {
		onlineUsers.Users[user.ID] = user
	}

	groups.Lock()
	defer groups.Unlock()

	if _, exists := groups.Data[user.Channel]; !exists {
		groups.Data[user.Channel] = make(map[string]model.PttUser)
	}
	groups.Data[user.Channel][user.ID] = user
}

func LeavePTTChannel(user model.PttUser) {
	onlineUsers.Lock()
	defer onlineUsers.Unlock()

	if _, exists := onlineUsers.Users[user.ID]; exists {
		delete(onlineUsers.Users, user.ID)
	}

	groups.Lock()
	defer groups.Unlock()

	if _, exists := groups.Data[user.Channel]; exists {
		delete(groups.Data[user.Channel], user.ID)
		if len(groups.Data[user.Channel]) == 0 {
			delete(groups.Data, user.Channel) // 删除空群组
		}
	}
}

func GetChannelUsers(channel string) []model.PttUser {
	groups.RLock()
	defer groups.RUnlock()
	if users, exists := groups.Data[channel]; exists {
		result := make([]model.PttUser, 0, len(users))
		for _, user := range users {
			result = append(result, user)
		}
		return result
	}
	return nil
}

func CheckUserInChannel(user model.PttUser) bool {
	groups.RLock()
	defer groups.RUnlock()

	if users, exists := groups.Data[user.Channel]; exists {
		if _, exists := users[user.ID]; exists {
			return true
		}
	}
	return false
}

func JoinChannel(c *fiber.Ctx) error {
	var user model.PttUser
	if err := c.BodyParser(&user); err != nil {
		return c.JSON(model.BaseRes(-1, "Failed to parse request body: "+err.Error(), 0))
	}

	if len(user.ID) < 10 || len(user.Channel) < 5 || strings.Count(user.Channel, "-") != 2 || len(user.Token) < 10 {
		return c.JSON(model.BaseRes(-1, "Invalid user ID, channel name, or token", 0))
	}

	JoinPTTChannel(user)

	return c.JSON(
		model.BaseRes(0, "Joined channel successfully", len(GetChannelUsers(user.Channel))),
	)
}

func LeaveChannel(c *fiber.Ctx) error {
	var user model.PttUser
	if err := c.BodyParser(&user); err != nil {
		return c.JSON(model.BaseRes(-1, "Failed to parse request body: "+err.Error(), 0))
	}

	if len(user.ID) < 10 || len(user.Channel) < 5 || strings.Count(user.Channel, "-") != 2 {
		return c.JSON(model.BaseRes(-1, "Invalid user ID, channel name, or token", 0))
	}

	LeavePTTChannel(user)
	return c.JSON(model.BaseRes(0, "Left channel successfully", 0))
}
func PingPTT(c *fiber.Ctx) error {
	channel := c.Params("channel")
	token := c.Get("X-Q")
	id := c.Get(fiber.HeaderAuthorization)

	user := model.PttUser{
		ID:      id,
		Token:   token,
		Channel: channel,
	}
	if !CheckUserInChannel(user) {
		JoinPTTChannel(user)
	}

	parts := strings.Split(channel, "-")
	if channel == "" && len(parts) != 2 {
		return c.Status(fiber.StatusNotFound).SendString("Not Found")
	}
	users := GetChannelUsers(channel)
	if users == nil {

		return c.JSON(model.BaseRes(-1, "No users found in this channel", 0))
	}

	return c.JSON(model.BaseRes(0, "Users retrieved successfully", len(users)))
}

func UploadVoice(c *fiber.Ctx) error {

	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(err.Error())
	}

	parts := strings.Split(file.Filename, "-")
	if len(parts) != 5 {
		log.Info("Invalid file name format.")
		return c.JSON("Invalid file name format.")
	}

	channel := parts[0] + "-" + parts[1] + "-" + parts[2]
	id := parts[3]

	if file.Size > 512*1024 { // 512KB limit
		log.Info("Long file size")
		return c.JSON("Long file size")
	}
	fileUrl := config.BaseDir("voices", file.Filename)
	_ = EnsurePath(fileUrl)
	if err := c.SaveFile(file, fileUrl); err != nil {
		log.Error(err.Error())
		return c.JSON(err.Error())
	}

	log.Info("File uploaded successfully:", file.Filename)

	VoiceList <- model.PttVoice{
		ID:       id,
		Channel:  channel,
		FileName: file.Filename,
	}

	return c.JSON("ok")
}

func GetVoice(c *fiber.Ctx) error {

	filename := c.Params("fileName")
	if filename == "" || len(filename) < 10 {
		return c.JSON(model.Failed(-1, "Filename is required"))
	}

	filePath := "./voices/" + filename

	//  检查有没有文件
	if _, err := os.Stat(filePath); os.IsNotExist(err) {

		return c.Status(fiber.StatusNotFound).JSON("No such file")
	}

	return c.SendFile(filePath)
}

// EnsurePath 确保路径存在，如果不存在则创建
func EnsurePath(path string) error {
	// 获取路径的目录部分（如果是文件路径，则获取其所在目录）
	dir := path

	// 检查是否是文件路径（包含扩展名）
	if filepath.Ext(path) != "" {
		dir = filepath.Dir(path)
	}

	// 创建目录（包括所有父目录）
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	return nil
}

func CirclePushPTT() {
	log.Info("PTT Voice Push List Handler initialized")
	go func() {
		for {
			select {
			case voice := <-VoiceList:
				users := GetChannelUsers(voice.Channel)
				go func() {
					for _, user := range users {

						if len(user.ID) < 10 {
							log.Info("Invalid user ID:", user.ID)
							continue
						}

						if user.ID == voice.ID {
							continue
						}

						log.Info(user.ID, voice.ID)
						err := push.PTTPush(map[string]interface{}{
							"fileName": voice.FileName,
							"channel":  voice.Channel,
							"id":       voice.ID,
						}, user.Token)
						if err != nil {
							log.Error(err.Error())
						}
						log.Info("push success")
					}
				}()

			}
		}
	}()
}
