package v1

import (
	"fmt"
	"spiritFruit/pkg/response"
	"os"
	"path/filepath"
	"strings"
	"time"
	
    "spiritFruit/pkg/upload"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)


type UploadsController struct {
	BaseADMINController
}

// UploadResponse 上传响应结构
type UploadResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    *UploadData `json:"data,omitempty"`
}

// UploadData 上传数据结构
type UploadData struct {
	FileName   string `json:"file_name"`
	FileSize   int64  `json:"file_size"`
	FileURL    string `json:"file_url"`
	UploadTime string `json:"upload_time"`
}

// Upload 处理图片上传
func (u *UploadsController) Upload(c *gin.Context) {


	err, fileName, fileURL, headSize := upload.UploadsTool(c)
	if err != "" {
		response.JSON(c, gin.H{
			"code":    400,
			"message": err,
			"data":    nil,
		})
		return
	}
	// 返回成功响应
	response.JSON(c, UploadResponse{
		Code:    0,
		Message: "上传成功",
		Data: &UploadData{
			FileName:   fileName,
			FileSize:   headSize,
			FileURL:    fileURL,
			UploadTime: time.Now().Format("2006-01-02 15:04:05"),
		},
	})

}

// isValidImageType 验证是否为有效的图片类型
func isValidImageType(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	validExts := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}

	for _, validExt := range validExts {
		if ext == validExt {
			return true
		}
	}
	return false
}

// BatchUpload 批量上传图片（可选功能）
func (u *UploadsController) BatchUpload(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		response.JSON(c, gin.H{
			"code":    400,
			"message": "解析表单失败: " + err.Error(),
		})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		response.JSON(c, gin.H{
			"code":    400,
			"message": "没有找到上传文件",
		})
		return
	}

	// 限制批量上传数量
	const maxFiles = 10
	if len(files) > maxFiles {
		response.JSON(c, gin.H{
			"code":    400,
			"message": fmt.Sprintf("批量上传最多支持%d个文件", maxFiles),
		})
		return
	}

	var uploadResults []UploadData
	var failedFiles []string

	uploadDir := "uploads/images/" + time.Now().Format("2006/01/02")
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
        response.JSON(c, UploadResponse{
            Code:    500,
            Message: "创建上传目录失败: " + err.Error(),
        })
        return
    }

	for _, file := range files {
		// 验证文件
		if !isValidImageType(file.Filename) || file.Size > (5<<20) {
			failedFiles = append(failedFiles, file.Filename)
			continue
		}

		// 生成文件名并保存
		ext := strings.ToLower(filepath.Ext(file.Filename))
		fileName := uuid.New().String() + ext
		filePath := filepath.Join(uploadDir, fileName)

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			failedFiles = append(failedFiles, file.Filename)
			continue
		}

		uploadResults = append(uploadResults, UploadData{
			FileName:   fileName,
			FileSize:   file.Size,
			FileURL:    fmt.Sprintf("/static/%s", filePath),
			UploadTime: time.Now().Format("2006-01-02 15:04:05"),
		})
	}

	responseData := gin.H{
		"code":    200,
		"message": "批量上传完成",
		"data": gin.H{
			"success_count": len(uploadResults),
			"failed_count":  len(failedFiles),
			"files":         uploadResults,
		},
	}

	if len(failedFiles) > 0 {
		responseData["failed_files"] = failedFiles
	}

	response.JSON(c, responseData)
}
