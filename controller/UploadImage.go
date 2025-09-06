package controller

import (
	"NoLetServer/config"
	"NoLetServer/model"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

// 允许的图片类型
var allowedImageTypes = map[string]bool{
	"image/jpeg":    true,
	"image/png":     true,
	"image/gif":     true,
	"image/webp":    true,
	"image/bmp":     true,
	"image/svg+xml": true,
}

// UploadController 处理图片上传请求
// 支持GET和POST两种请求方式:
// GET: 返回上传页面
// POST: 上传图片并保存
func UploadController(c *fiber.Ctx) error {
	// 验证管理员权限
	ok := model.Admin(c)

	if c.Method() == fiber.MethodGet {
		return c.Render("upload", c.Queries())
	}

	if !ok {
		log.Info("Unauthorized upload attempt")
		return c.SendStatus(http.StatusUnauthorized)
	}

	// 获取文件名
	file, err := c.FormFile("filename")

	if err != nil || file.Filename == "" {

		return c.JSON(fiber.Map{
			"error": "filename is required",
		})
	}

	// 创建上传目录
	uploadDir := config.BaseDir("images")
	_ = EnsurePath(uploadDir)

	// 验证文件是否在uploads目录下
	if isTrue, _ := isFileInDirectory(file.Filename, uploadDir); isTrue {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "file not found in uploads directory",
		})
	}

	// 验证文件类型
	if !allowedImageTypes[file.Header.Get(fiber.HeaderContentType)] {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "only image files are allowed",
		})
	}

	// 生成安全的文件名
	ext := filepath.Ext(file.Filename)
	if ext == "" {
		ext = filepath.Ext(file.Filename)
	}
	if ext == "" {
		ext = ".jpg" // 默认扩展名
	}

	if file.Size > 1024*1024*3 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Image size is too big",
		})
	}

	// 使用时间戳和随机字符串生成唯一文件名
	safeFileName := fmt.Sprintf("%s%s",
		strings.TrimSuffix(file.Filename, ext),
		strings.ToLower(ext))

	// 完整的文件路径
	filePath := filepath.Join(uploadDir, safeFileName)

	// 保存文件

	if err = c.SaveFile(file, filePath); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to save file: " + err.Error(),
		})
	}

	// 返回成功响应
	return c.JSON(fiber.Map{
		"message":  "file uploaded successfully",
		"filename": safeFileName,
		"path":     filePath,
		"size":     file.Size,
		"type":     file.Header.Get(fiber.HeaderContentType),
	})
}

func GetImage(c *fiber.Ctx) error {
	fileName := c.Params("filename")
	if fileName == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "filename is required",
		})
	}

	// 构建文件路径
	filePath := filepath.Join("./images", fileName)

	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {

		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "file not found",
		})
	}

	// 返回文件
	return c.SendFile(filePath)
}

func isFileInDirectory(dirPath, fileName string) (bool, error) {
	// 对目录路径进行规范化处理
	dirPath = filepath.Clean(dirPath)

	// 检查目录是否存在
	dirInfo, err := os.Stat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil // 目录不存在时，直接返回文件不在目录中
		}
		return false, fmt.Errorf("检查目录状态出错: %w", err)
	}

	// 确认路径指向的是一个目录
	if !dirInfo.IsDir() {
		return false, fmt.Errorf("路径 %q 不是一个目录", dirPath)
	}

	// 构建文件的完整路径
	filePath := filepath.Join(dirPath, fileName)

	// 检查文件是否存在
	_, err = os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil // 文件不存在
		}
		return false, fmt.Errorf("检查文件状态出错: %w", err)
	}

	return true, nil
}
