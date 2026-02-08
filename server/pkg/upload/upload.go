package upload

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"spiritFruit/pkg/str"
	"strings"
	"time"
)

// UploadData 上传数据结构
type UploadData struct {
	FileName   string `json:"file_name"`
	FileSize   int64  `json:"file_size"`
	FileURL    string `json:"file_url"`
	UploadTime string `json:"upload_time"`
}

func UploadsTool(c *gin.Context) (errStr string, fileName, filePath string, headSize int64) {
	// 设置最大内存限制为 10MB
	err := c.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		return errStr, "", "", 0
	}
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		return err.Error(), "", "", 0
	}
	defer file.Close()
	// 验证文件类型
	if !isValidImageType(header.Filename) {
		return
	}

	// 验证文件大小（限制为5MB）
	const maxSize = 5 << 20 // 5MB
	if header.Size > maxSize {
		return "文件大小超过限制，最大支持5MB", "", "", 0
	}

	// 创建上传目录
	uploadDir := "uploads/images/" + time.Now().Format("2006/01/02")
	if mkdirErr := os.MkdirAll(uploadDir, 0755); mkdirErr != nil {
		return "创建上传目录失败" + mkdirErr.Error(), "", "", 0
	}
	// 生成唯一文件名
	ext := strings.ToLower(filepath.Ext(header.Filename))
	fileNameValue := uuid.New().String() + ext
	filePathValue := path.Join(uploadDir, fileNameValue)
	// 保存文件
	if saveErr := c.SaveUploadedFile(header, filePathValue); saveErr != nil {
		return "保存文件失败: " + saveErr.Error(), "", "", 0
	}
	// 生成文件访问URL
	fileURL := fmt.Sprintf("%s", filePathValue)
	return "", fileName, fileURL, header.Size
}

func BatchUpload(c *gin.Context) (errStr string, uploadData []UploadData, failedFiles []string) {
	form, err := c.MultipartForm()
	if err != nil {
		return "解析表单失败: " + err.Error(), uploadData, failedFiles
	}

	files := form.File["files"]
	if len(files) == 0 {
		return "没有上传文件", uploadData, failedFiles
	}

	// 限制批量上传数量
	const maxFiles = 10
	if len(files) > maxFiles {
		return fmt.Sprintf("批量上传最多支持%d个文件", maxFiles), uploadData, failedFiles
	}

	var uploadResults []UploadData
	var failedFilesMap []string

	uploadDir := "uploads/images/" + time.Now().Format("2006/01/02")
	if mkdirErr := os.MkdirAll(uploadDir, 0755); mkdirErr != nil {
		return "创建上传目录失败: " + mkdirErr.Error(), uploadData, failedFiles
	}

	for _, file := range files {
		// 验证文件
		if !isValidImageType(file.Filename) || file.Size > (5<<20) {
			failedFilesMap = append(failedFilesMap, file.Filename)
			continue
		}

		// 生成文件名并保存
		ext := strings.ToLower(filepath.Ext(file.Filename))
		fileName := uuid.New().String() + ext
		filePath := filepath.Join(uploadDir, fileName)

		if saveErr := c.SaveUploadedFile(file, filePath); saveErr != nil {
			failedFilesMap = append(failedFilesMap, file.Filename)
			continue
		}

		uploadResults = append(uploadResults, UploadData{
			FileName:   fileName,
			FileSize:   file.Size,
			FileURL:    fmt.Sprintf("/static/%s", filePath),
			UploadTime: time.Now().Format("2006-01-02 15:04:05"),
		})
	}
	return "", uploadResults, failedFilesMap
}

// isValidImageType 验证是否为有效的图片类型
func isValidImageType(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	validExits := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}

	for _, validExt := range validExits {
		if ext == validExt {
			return true
		}
	}
	return false
}

func DecodeBase64Image(base64Str, fileNameValue string) (errStr string, filePath string) {
	// 解析 Base64 图片数据
	imageData, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return "无效的 Base64 图片数据", ""
	}
	// 生成一个唯一的文件名，以防止文件名冲突
	fileName := str.GenerateUniqueFileName(fileNameValue)
	// 指定文件的保存路径
	filePath = filepath.Join("uploads", fileName)
	// 将 Base64 图片数据保存到指定路径
	err = ioutil.WriteFile(filePath, imageData, 0644)
	if err != nil {
		return "保存图片失败", ""
	}
	return "", filePath
}

// SaveImageByte 通用文件保存逻辑 (供 Job 使用)
// data: 图片的二进制数据
// ext: 文件后缀 (e.g. ".png", ".jpg")
// return: 相对路径 (e.g. "uploads/images/2024/02/08/xxx.png"), error
func SaveImageByte(data []byte, ext string) (string, error) {
	// 1. 创建上传目录 (保持与 UploadsTool 一致的结构)
	uploadDir := "uploads/images/" + time.Now().Format("2006/01/02")
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", fmt.Errorf("create dir failed: %v", err)
	}

	// 2. 生成唯一文件名
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	fileName := uuid.New().String() + ext
	filePath := path.Join(uploadDir, fileName)

	// 3. 写入文件
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return "", fmt.Errorf("write file failed: %v", err)
	}

	// 4. 返回相对路径 (注意：Windows下如果是反斜杠可能需要转义，这里假设是 Linux/Mac 风格或用于 Web URL)
	// 为了前端访问方便，通常统一为 "/"
	return strings.ReplaceAll(filePath, "\\", "/"), nil
}

// DownloadAndSave 从 URL 下载并保存 (针对 OpenAI 返回 URL 的情况)
func DownloadAndSave(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 简单推断后缀，默认 png
	return SaveImageByte(data, ".png")
}

// SaveBase64Image 保存 Base64 图片 (针对 Gemini 返回 DataURI 的情况)
func SaveBase64Image(base64Str string) (string, error) {
	// 去掉 data:image/png;base64, 前缀
	parts := strings.Split(base64Str, ",")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid base64 format")
	}

	// 获取后缀 (简单解析)
	mimeType := strings.Split(strings.Split(parts[0], ";")[0], ":")[1] // image/png
	ext := "." + strings.Split(mimeType, "/")[1]

	data, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", err
	}

	return SaveImageByte(data, ext)
}
