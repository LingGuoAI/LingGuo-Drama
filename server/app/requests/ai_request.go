package requests

import (
	"github.com/gin-gonic/gin"
	"github.com/thedevsaddam/govalidator"
)

type GenerateCharactersRequest struct {
	DramaID uint   `json:"dramaId" valid:"dramaId"`
	Count   int    `json:"count" valid:"count"` // 生成数量
	Outline string `json:"outline"`             // 可选：指定大纲内容，不填则使用剧本简介
}

type ExtractScenesRequest struct {
	EpisodeID uint `json:"episodeId" valid:"episodeId"`
}

func ValidateGenerateCharacters(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"dramaId": []string{"required", "numeric"},
		"count":   []string{"numeric_between:1,20"},
	}
	messages := govalidator.MapData{
		"dramaId": []string{"required:剧本ID不能为空"},
		"count":   []string{"numeric_between:生成数量需在1-20之间"},
	}
	return validate(data, rules, messages)
}

func ValidateExtractScenes(data interface{}, c *gin.Context) map[string][]string {
	rules := govalidator.MapData{
		"episodeId": []string{"required", "numeric"},
	}
	messages := govalidator.MapData{
		"episodeId": []string{"required:章节ID不能为空"},
	}
	return validate(data, rules, messages)
}
