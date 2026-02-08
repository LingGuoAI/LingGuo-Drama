package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hibiken/asynq"
	"gorm.io/gorm"

	"spiritFruit/app/models/async_tasks"
	"spiritFruit/app/models/characters"
	"spiritFruit/app/models/scenes"
	"spiritFruit/app/models/scripts"
	"spiritFruit/app/models/shots" // 假设您的分镜模型在这个包
	myAsynq "spiritFruit/pkg/asynq"
	"spiritFruit/pkg/config"
	"spiritFruit/pkg/console"
	"spiritFruit/pkg/database"
	"spiritFruit/pkg/openai"
	"spiritFruit/pkg/utils"
)

// HandleGenerateShots 处理分镜生成任务
func HandleGenerateShots(ctx context.Context, t *asynq.Task) error {
	var p myAsynq.GenerateShotsPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json unmarshal failed: %v", err)
	}

	// 1. 获取任务记录
	taskModel := async_tasks.AsyncTask{}
	if err := database.DB.First(&taskModel, p.AsyncTaskID).Error; err != nil {
		return nil
	}
	taskModel.MarkAsProcessing()
	console.Success(fmt.Sprintf("任务[%d] - 开始拆分分镜", p.AsyncTaskID))

	// 2. 准备数据
	taskModel.UpdateProgress(10)

	// A. 获取剧本内容
	var script scripts.Scripts
	if err := database.DB.First(&script, p.ScriptID).Error; err != nil {
		err = fmt.Errorf("script not found: %v", err)
		taskModel.MarkAsFailed(err)
		return nil
	}
	if script.Content == nil || *script.Content == "" {
		err := fmt.Errorf("script content is empty")
		taskModel.MarkAsFailed(err)
		return nil
	}
	scriptContent := *script.Content

	// B. 获取角色列表 (用于 Prompt 上下文)
	var charList []characters.Characters
	database.DB.Where("project_id = ?", p.ProjectID).Find(&charList)
	charInfoStr := "无角色"
	if len(charList) > 0 {
		var infos []string
		for _, c := range charList {
			cName := ""
			if c.Name != nil {
				cName = *c.Name
			}
			infos = append(infos, fmt.Sprintf(`{"id": %d, "name": "%s"}`, c.ID, cName))
		}
		charInfoStr = fmt.Sprintf("[%s]", strings.Join(infos, ", "))
	}

	// C. 获取场景列表 (用于 Prompt 上下文)
	var sceneList []scenes.Scenes
	database.DB.Where("project_id = ?", p.ProjectID).Find(&sceneList)
	sceneInfoStr := "无场景"
	if len(sceneList) > 0 {
		var infos []string
		for _, s := range sceneList {
			sName := ""
			if s.Name != nil {
				sName = *s.Name
			}
			infos = append(infos, fmt.Sprintf(`{"id": %d, "name": "%s"}`, s.ID, sName))
		}
		sceneInfoStr = fmt.Sprintf("[%s]", strings.Join(infos, ", "))
	}

	// 3. 构建 Prompt
	taskModel.UpdateProgress(30)
	systemPrompt := getStoryboardSystemPrompt() // 下方定义的辅助函数
	userPrompt := fmt.Sprintf(`
【本剧可用角色列表(JSON)】:
%s

【本剧已提取的场景列表(JSON)】:
%s

【剧本原文】:
%s
`, charInfoStr, sceneInfoStr, scriptContent)

	// 4. 调用 AI
	taskModel.UpdateProgress(40)
	aiConfig := openai.Config{
		Provider:      config.GetString("ai.provider"),
		OpenAIBaseURL: config.GetString("ai.openai.base_url"),
		OpenAIKey:     config.GetString("ai.openai.api_key"),
		OpenAIModel:   p.Model, // 如果 payload 没传，Provide 内部会处理默认值
		// ... Gemini 配置 ...
		GeminiBaseURL: config.GetString("ai.gemini.base_url"),
		GeminiKey:     config.GetString("ai.gemini.api_key"),
		GeminiModel:   config.GetString("ai.gemini.model"),
	}
	// 如果没有指定模型，使用默认长文本能力强的模型
	if aiConfig.OpenAIModel == "" {
		aiConfig.OpenAIModel = "gpt-4-turbo"
	}

	provider := openai.NewProvider(aiConfig)

	// 使用较大的 MaxTokens 确保返回完整 JSON
	// 注意：openai.ScriptRequest 结构体可能需要扩展 MaxTokens 字段，或者在 GenerateScript 内部处理
	// 这里假设 provider.GenerateScript 能够处理基本的对话
	aiReq := openai.ScriptRequest{
		Messages: []openai.ChatMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userPrompt},
		},
		Temperature: 0.7,
	}

	aiResp, err := provider.GenerateScript(aiReq)
	if err != nil {
		taskModel.MarkAsFailed(err)
		return err
	}

	// 5. 解析结果
	taskModel.UpdateProgress(70)

	// 定义接收 AI 结果的结构
	type AIStoryboard struct {
		ShotNumber   int      `json:"shot_number"`
		Title        string   `json:"title"`
		ShotType     string   `json:"shot_type"` // 景别
		Angle        string   `json:"angle"`
		Movement     string   `json:"movement"`
		Time         string   `json:"time"`
		Location     string   `json:"location"`
		SceneID      *uint64  `json:"scene_id"` // 可空
		Action       string   `json:"action"`
		Dialogue     string   `json:"dialogue"`
		Duration     int      `json:"duration"`
		VisualPrompt string   `json:"result"` // 对应 Prompt 里的 "result" 字段，作为画面描述
		CharacterIDs []uint64 `json:"characters"`
	}

	var resultWrapper struct {
		Storyboards []AIStoryboard `json:"storyboards"`
	}

	// 尝试解析对象格式 {"storyboards": [...]}
	if err := utils.SafeParseAIJSON(aiResp, &resultWrapper); err != nil {
		// 尝试直接解析数组 [...]
		var arr []AIStoryboard
		if err2 := utils.SafeParseAIJSON(aiResp, &arr); err2 == nil {
			resultWrapper.Storyboards = arr
		} else {
			taskModel.MarkAsFailed(fmt.Errorf("failed to parse JSON: %v", err))
			return nil
		}
	}

	if len(resultWrapper.Storyboards) == 0 {
		taskModel.MarkAsFailed(fmt.Errorf("no shots parsed"))
		return nil
	}

	// 6. 入库 (使用事务)
	taskModel.UpdateProgress(80)

	err = database.DB.Transaction(func(tx *gorm.DB) error {
		// A. 清理旧分镜 (可选，根据业务需求，这里选择覆盖模式)
		if err := tx.Where("script_id = ?", p.ScriptID).Delete(&shots.Shots{}).Error; err != nil {
			return err
		}

		// B. 插入新分镜
		for i, sb := range resultWrapper.Storyboards {
			// 准备基础变量
			projectID := p.ProjectID
			scriptID := p.ScriptID
			sequenceNo := uint64(i + 1)
			duration := uint64(sb.Duration * 1000)
			status := int8(0)

			// 准备字符串变量 (Go中需要取地址，先定义变量更清晰)
			title := sb.Title
			shotType := sb.ShotType
			cameraMove := sb.Movement
			angle := sb.Angle

			// 独立字段赋值，不合并
			timeDesc := sb.Time             // 时间描述
			locDesc := sb.Location          // 地点描述
			action := sb.Action             // 动作描述
			visualResult := sb.VisualPrompt // 画面结果/详细描述 (对应AI返回的 result 字段)
			dialogue := sb.Dialogue

			imagePrompt := visualResult

			// 构造 Shot 模型
			newShot := shots.Shots{
				ProjectId:  &projectID,
				ScriptId:   &scriptID,
				SceneId:    sb.SceneID, // AI返回的 scene_id
				SequenceNo: &sequenceNo,

				Title:          &title,
				ShotType:       &shotType,
				CameraMovement: &cameraMove,
				Angle:          &angle,

				Time:       &timeDesc,     // 独立存
				Location:   &locDesc,      // 独立存
				Action:     &action,       // 独立存
				VisualDesc: &visualResult, // 独立存 (对应 result)

				Dialogue: &dialogue,

				ImagePrompt: &imagePrompt, // 初始 Prompt
				DurationMs:  &duration,
				Status:      &status,
			}

			if err := tx.Create(&newShot).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		taskModel.MarkAsFailed(fmt.Errorf("db transaction failed: %v", err))
		return err
	}

	// 7. 完成
	resultData := map[string]interface{}{
		"total_shots": len(resultWrapper.Storyboards),
	}
	resBytes, _ := json.Marshal(resultData)
	taskModel.MarkAsSuccess(string(resBytes))

	console.Success(fmt.Sprintf("任务[%d] - 分镜拆分完成，共 %d 个镜头", p.AsyncTaskID, len(resultWrapper.Storyboards)))
	return nil
}

// getStoryboardSystemPrompt 返回复杂的 Prompt 提示词
func getStoryboardSystemPrompt() string {
	return `你是一个专业的影视分镜师。请根据用户提供的【剧本原文】，结合【角色列表】和【场景列表】，将其拆解为详细的分镜头脚本。

【分镜要素要求】每个镜头聚焦单一动作，描述要详尽具体：
1. **镜头标题(title)**：3-5字概括核心内容。
2. **场景(scene_id)**：从提供的【场景列表】中选择最匹配的ID（数字）。如果没有合适的，返回null。
3. **镜头设计**：
   - **景别(shot_type)**：[远景/全景/中景/近景/特写]
   - **角度(angle)**：[平视/仰视/俯视/侧面/背面]
   - **运镜(movement)**：[固定/推/拉/摇/跟/移]
4. **人物(characters)**：该镜头出现的人物ID数组。仅使用【角色列表】中存在的ID。
5. **对话(dialogue)**：提取该镜头中的台词（保留角色名）。无对话留空。
6. **画面结果(result)**：用于AI生图的详细画面描述。包含：人物动作+表情+光影+氛围+环境细节。
   - 必须包含≥20字的详细描述。
   - 像给盲人讲画面一样详细。
7. **时长(duration)**：估算该镜头时长（秒），范围4-12秒。

【输出JSON格式】:
{
  "storyboards": [
    {
      "shot_number": 1,
      "title": "...",
      "shot_type": "...",
      "angle": "...",
      "movement": "...",
      "time": "...",
      "location": "...",
      "scene_id": 1, 
      "action": "...",
      "dialogue": "...",
      "result": "详细的画面描述...",
      "duration": 5,
      "characters": [1, 2]
    }
  ]
}
注意：
1. 必须100%覆盖剧本内容，不得遗漏。
2. 严格输出 JSON 格式，不要包含 Markdown 代码块标记。
`
}
