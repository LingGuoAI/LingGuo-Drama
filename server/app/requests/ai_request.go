package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

type GenerateCharactersRequest struct {
	ProjectId uint   `json:"projectId" valid:"projectId"`
	Count     int    `json:"count" valid:"count"` // 生成数量
	Outline   string `json:"outline"`             // 可选：指定大纲内容，不填则使用剧本简介
}

type ExtractScenesRequest struct {
	ScriptId uint `json:"scriptId" valid:"scriptId"`
}

func ValidateGenerateCharacters(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"projectId": []string{"required", "numeric"},
		"count":     []string{"numeric_between:1,20"},
	}
	messages := govalidator.MapData{
		"projectId": []string{"required:剧本ID不能为空"},
		"count":     []string{"numeric_between:生成数量需在1-20之间"},
	}
	return validate(data, rules, messages)
}

func ValidateExtractScenes(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"scriptId": []string{"required", "numeric"},
	}
	messages := govalidator.MapData{
		"scriptId": []string{"required:章节ID不能为空"},
	}
	return validate(data, rules, messages)
}

// ExtractPropsRequest 提取道具请求
type ExtractPropsRequest struct {
	EpisodeId uint64 `json:"episodeId" valid:"required"`
}

// GeneratePropImageRequest 单个道具生图请求
type GeneratePropImageRequest struct {
	PropId uint64 `json:"propId" valid:"required"`
	Model  string `json:"model"` // 可选：指定绘画模型
}

// BatchGeneratePropImagesRequest 批量道具生图请求
type BatchGeneratePropImagesRequest struct {
	PropIds []uint64 `json:"propIds" valid:"required"`
	Model   string   `json:"model"`
}
