package services

import (
	"fmt"
	"spiritFruit/app/models/ai_config"
)

// AiConfigService AI配置服务结构体
type AiConfigService struct {
}

// GetActiveConfigByType 获取指定类型（优先级最高）的激活配置
// serviceType: text, image, video
func (s *AiConfigService) GetActiveConfigByType(serviceType string) (error, ai_config.AiConfig) {
	// 获取指定服务类型且状态为启用的配置
	// 注：如果在 model 层的查询已按 priority DESC 排序，这里取到的就是优先级最高的可用配置
	configInfo := ai_config.GetByWhereMap(map[string]interface{}{
		"service_type": serviceType,
		"is_active":    1,
	})

	if configInfo.ID == 0 {
		return fmt.Errorf("未找到激活的 %s AI配置，请前往后台配置", serviceType), ai_config.AiConfig{}
	}

	return nil, configInfo
}

// GetAllActiveConfigsByType 获取指定类型的所有激活配置列表（用于业务层做轮询或降级策略）
func (s *AiConfigService) GetAllActiveConfigsByType(serviceType string) (error, []ai_config.AiConfig) {
	configList := ai_config.GetMapDataByWhereMap(map[string]interface{}{
		"service_type": serviceType,
		"is_active":    1,
	})

	if len(configList) == 0 {
		return fmt.Errorf("未找到任何激活的 %s AI配置", serviceType), []ai_config.AiConfig{}
	}

	return nil, configList
}

// GetActiveConfigByProvider 获取指定厂商和类型的激活配置 (例如明确指定要用 openai 的文本模型)
func (s *AiConfigService) GetActiveConfigByProvider(serviceType, provider string) (error, ai_config.AiConfig) {
	configInfo := ai_config.GetByWhereMap(map[string]interface{}{
		"service_type": serviceType,
		"provider":     provider,
		"is_active":    1,
	})

	if configInfo.ID == 0 {
		return fmt.Errorf("未找到厂商为 %s 且类型为 %s 的激活AI配置", provider, serviceType), ai_config.AiConfig{}
	}

	return nil, configInfo
}

// GetSpecificModelConfig 检查指定提供商是否包含所需的具体模型
func (s *AiConfigService) GetSpecificModelConfig(serviceType, provider, specificModel string) (error, ai_config.AiConfig) {
	err, configInfo := s.GetActiveConfigByProvider(serviceType, provider)
	if err != nil {
		return err, ai_config.AiConfig{}
	}

	// 遍历检查该配置支持的模型列表中是否包含需要的具体模型
	hasModel := false
	for _, m := range configInfo.Model {
		if m == specificModel {
			hasModel = true
			break
		}
	}

	if !hasModel {
		return fmt.Errorf("厂商 %s 的配置中未包含 %s 模型支持", provider, specificModel), ai_config.AiConfig{}
	}

	return nil, configInfo
}
