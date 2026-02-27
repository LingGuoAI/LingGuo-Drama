package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"spiritFruit/app/models/characters"
	"spiritFruit/app/models/props"
	"spiritFruit/app/models/scenes"
	"spiritFruit/app/models/scripts"
	"spiritFruit/app/models/shots"
	"spiritFruit/app/requests"
	"spiritFruit/app/services"
	"spiritFruit/pkg/database"
	"spiritFruit/pkg/response"
)

type AiController struct {
	BaseADMINController
}

// GenerateCharacters 异步生成角色
func (ctrl *AiController) GenerateCharacters(c *gin.Context) {
	// 1. 验证参数
	request := requests.GenerateCharactersRequest{}
	if ok := requests.Validate(c, &request, requests.ValidateGenerateCharacters); !ok {
		return
	}

	// 2. 调用 TaskService
	taskService := new(services.TaskService)
	task, err := taskService.CreateGenerateCharactersTask(uint64(request.ProjectId), request.Count, request.Outline)

	if err != nil {
		response.Abort500(c, "任务启动失败: "+err.Error())
		return
	}

	// 3. 返回结果 (taskId 是数据库主键 ID)
	response.JSON(c, gin.H{
		"status":  0,
		"message": "角色生成任务已在后台运行",
		"data": map[string]interface{}{
			"task_id": task.ID,
			"status":  task.Status,
		},
	})
}

// ExtractScenes 异步提取场景
func (ctrl *AiController) ExtractScenes(c *gin.Context) {
	request := requests.ExtractScenesRequest{}
	if ok := requests.Validate(c, &request, requests.ValidateExtractScenes); !ok {
		return
	}

	// 查询关联的项目ID
	var scriptsInfo scripts.Scripts
	if err := database.DB.First(&scriptsInfo, request.ScriptId).Error; err != nil {
		response.Abort500(c, "未找到对应章节")
		return
	}
	// 注意：episode.ProjectId 是指针
	var projectID uint64
	if scriptsInfo.ProjectId != nil {
		projectID = *scriptsInfo.ProjectId
	}

	// 2. 调用 TaskService
	taskService := new(services.TaskService)
	task, err := taskService.CreateExtractScenesTask(projectID, uint64(request.ScriptId))

	if err != nil {
		response.Abort500(c, "任务启动失败: "+err.Error())
		return
	}

	// 3. 返回结果
	response.JSON(c, gin.H{
		"status":  0,
		"message": "场景提取任务已在后台运行",
		"data": map[string]interface{}{
			"task_id": task.ID,
			"status":  task.Status,
		},
	})
}

// GenerateCharacterImage 生成角色图片
func (ctrl *AiController) GenerateCharacterImage(c *gin.Context) {
	// 1. 定义简易请求参数 (或者在 requests 包中定义验证器)
	type Req struct {
		CharacterID uint64 `json:"characterId" binding:"required"`
	}
	var req Req
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Abort500(c, "参数错误")
		return
	}

	// 2. 查询角色信息以获取 visualPrompt
	var char characters.Characters
	if err := database.DB.First(&char, req.CharacterID).Error; err != nil {
		response.Abort500(c, "角色不存在")
		return
	}

	prompt := ""
	if char.VisualPrompt != nil {
		prompt = *char.VisualPrompt
	}
	if prompt == "" {
		// 如果没有 visualPrompt，使用外观描述回退
		if char.AppearanceDesc != nil {
			prompt = *char.AppearanceDesc
		} else {
			response.Abort500(c, "角色缺少外貌描述，无法生成")
			return
		}
	}

	// 确保 ProjectId 存在
	projectID := uint64(0)
	if char.ProjectId != nil {
		projectID = *char.ProjectId
	}

	// 3. 调用 Service 创建任务
	taskService := new(services.TaskService)
	task, err := taskService.CreateImageGenerationTask(projectID, req.CharacterID, prompt)

	if err != nil {
		response.Abort500(c, "任务启动失败: "+err.Error())
		return
	}

	// 4. 返回结果
	response.JSON(c, gin.H{
		"status":  200,
		"message": "图片生成任务已在后台运行",
		"data": map[string]interface{}{
			"task_id": task.ID,
			"status":  task.Status,
		},
	})
}

// BatchGenerateCharacterImages 批量生成角色图片
func (ctrl *AiController) BatchGenerateCharacterImages(c *gin.Context) {
	// 1. 定义请求参数：接收 ID 数组
	type BatchReq struct {
		CharacterIDs []uint64 `json:"characterIds" binding:"required,min=1,max=10"` // 限制一次最多10个，防止超时
	}
	var req BatchReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Abort500(c, "参数错误: "+err.Error())
		return
	}

	// 2. 准备返回结果结构
	type TaskResult struct {
		CharacterID uint64 `json:"character_id"`
		TaskID      uint64 `json:"task_id"`
		Status      int    `json:"status"`
	}
	var results []TaskResult
	taskService := new(services.TaskService)

	// 3. 遍历 ID，逐个创建任务
	// 注意：这里是在 Controller 层循环调用 Service。
	// 虽然不是最高效（比如可以做批量 Insert），但对于 Asynq 任务投递来说，逐个投递更稳健，
	// 且能确保每个角色都有独立的 TaskID 供前端追踪进度。
	for _, charID := range req.CharacterIDs {
		// A. 查询角色信息获取 Prompt
		var char characters.Characters
		if err := database.DB.First(&char, charID).Error; err != nil {
			// 如果某个角色没找到，记录错误或跳过，这里选择跳过
			continue
		}

		prompt := ""
		if char.VisualPrompt != nil {
			prompt = *char.VisualPrompt
		}
		if prompt == "" && char.AppearanceDesc != nil {
			prompt = *char.AppearanceDesc
		}
		if prompt == "" {
			continue // 无描述无法生成
		}

		projectID := uint64(0)
		if char.ProjectId != nil {
			projectID = *char.ProjectId
		}

		// B. 创建任务
		task, err := taskService.CreateImageGenerationTask(projectID, charID, prompt)
		if err == nil {
			results = append(results, TaskResult{
				CharacterID: charID,
				TaskID:      task.ID,
				Status:      task.Status,
			})
		}
	}

	// 4. 返回结果列表
	response.JSON(c, gin.H{
		"status":  200,
		"message": "批量任务已提交",
		"data":    results,
	})
}

// GenerateSceneImage 单个场景生图
func (ctrl *AiController) GenerateSceneImage(c *gin.Context) {
	type Req struct {
		SceneID uint64 `json:"sceneId" binding:"required"`
	}
	var req Req
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Abort500(c, "参数错误")
		return
	}

	// 1. 查询场景信息
	var scene scenes.Scenes
	if err := database.DB.First(&scene, req.SceneID).Error; err != nil {
		response.Abort500(c, "场景不存在")
		return
	}

	// 2. 准备 Prompt (优先用 VisualPrompt，如果没有则用 Atmosphere)
	// 如果 VisualPrompt 已经是 URL 了 (http开头)，需要处理这种情况
	// 通常这里应该有一个原始 Prompt 字段，或者重新拼接
	prompt := ""
	if scene.VisualPrompt != nil && *scene.VisualPrompt != "" {
		// 简单的判断，如果不是 URL，则认为是 Prompt
		if len(*scene.VisualPrompt) < 4 || (*scene.VisualPrompt)[:4] != "http" && (*scene.VisualPrompt)[0] != '/' {
			prompt = *scene.VisualPrompt
		}
	}
	// 如果 VisualPrompt 是空的或者是 URL，尝试使用 Atmosphere + Location + Time
	if prompt == "" {
		loc := ""
		if scene.Location != nil {
			loc = *scene.Location
		}
		tm := ""
		if scene.Time != nil {
			tm = *scene.Time
		}
		atm := ""
		if scene.Atmosphere != nil {
			atm = *scene.Atmosphere
		}

		prompt = fmt.Sprintf("%s, %s, %s", loc, tm, atm)
	}

	projectID := uint64(0)
	if scene.ProjectId != nil {
		projectID = *scene.ProjectId
	}

	// 3. 创建任务
	taskService := new(services.TaskService)
	task, err := taskService.CreateSceneImageGenerationTask(projectID, req.SceneID, prompt)
	if err != nil {
		response.Abort500(c, "任务启动失败: "+err.Error())
		return
	}

	response.JSON(c, gin.H{
		"code":    0,
		"message": "任务已提交",
		"data": map[string]interface{}{
			"task_id": task.ID,
		},
	})
}

// BatchGenerateSceneImages 批量场景生图
func (ctrl *AiController) BatchGenerateSceneImages(c *gin.Context) {
	type Req struct {
		SceneIDs []uint64 `json:"sceneIds" binding:"required,min=1,max=20"`
	}
	var req Req
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Abort500(c, "参数错误: "+err.Error())
		return
	}

	type TaskResult struct {
		SceneID uint64 `json:"scene_id"`
		TaskID  uint64 `json:"task_id"`
	}
	var results []TaskResult
	taskService := new(services.TaskService)

	for _, sceneID := range req.SceneIDs {
		var scene scenes.Scenes
		if err := database.DB.First(&scene, sceneID).Error; err != nil {
			continue
		}

		// 构造 Prompt (逻辑同上)
		prompt := ""
		if scene.VisualPrompt != nil && *scene.VisualPrompt != "" && (*scene.VisualPrompt)[0] != '/' && (*scene.VisualPrompt)[:4] != "http" {
			prompt = *scene.VisualPrompt
		}
		if prompt == "" {
			loc := ""
			if scene.Location != nil {
				loc = *scene.Location
			}
			tm := ""
			if scene.Time != nil {
				tm = *scene.Time
			}
			atm := ""
			if scene.Atmosphere != nil {
				atm = *scene.Atmosphere
			}
			prompt = fmt.Sprintf("%s, %s, %s", loc, tm, atm)
		}

		projectID := uint64(0)
		if scene.ProjectId != nil {
			projectID = *scene.ProjectId
		}

		task, err := taskService.CreateSceneImageGenerationTask(projectID, sceneID, prompt)
		if err == nil {
			results = append(results, TaskResult{
				SceneID: sceneID,
				TaskID:  task.ID,
			})
		}
	}

	response.JSON(c, gin.H{
		"code":    0,
		"message": "批量任务已提交",
		"data":    results,
	})
}

// GenerateShots 智能拆分分镜
func (ctrl *AiController) GenerateShots(c *gin.Context) {
	type Req struct {
		ScriptID uint64 `json:"scriptId" binding:"required"`
		Model    string `json:"model"` // 可选
	}
	var req Req
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Abort500(c, "参数错误")
		return
	}

	// 1. 校验剧本是否存在
	var script scripts.Scripts
	if err := database.DB.First(&script, req.ScriptID).Error; err != nil {
		response.Abort500(c, "剧本不存在")
		return
	}

	projectID := uint64(0)
	if script.ProjectId != nil {
		projectID = *script.ProjectId
	}

	// 2. 创建任务
	taskService := new(services.TaskService)
	task, err := taskService.CreateGenerateShotsTask(projectID, req.ScriptID, req.Model)
	if err != nil {
		response.Abort500(c, "任务启动失败: "+err.Error())
		return
	}

	response.JSON(c, gin.H{
		"code":    0,
		"message": "分镜拆分任务已提交",
		"data": map[string]interface{}{
			"task_id": task.ID,
		},
	})
}

// ExtractProps 从剧本提取道具
func (ctrl *AiController) ExtractProps(c *gin.Context) {
	type Req struct {
		EpisodeID uint64 `json:"episodeId" binding:"required"`
	}
	var req Req
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Abort500(c, "参数错误: "+err.Error())
		return
	}

	// 查剧本信息
	var script scripts.Scripts
	if err := database.DB.First(&script, req.EpisodeID).Error; err != nil {
		response.Abort500(c, "剧本不存在")
		return
	}

	taskService := new(services.TaskService)
	projectID := uint64(0)
	if script.ProjectId != nil {
		projectID = *script.ProjectId
	}

	// 创建任务
	task, err := taskService.CreateExtractPropsTask(projectID, req.EpisodeID)
	if err != nil {
		response.Abort500(c, "任务创建失败: "+err.Error())
		return
	}

	response.JSON(c, gin.H{
		"code": 0,
		"data": gin.H{
			"task_id": task.ID,
		},
		"message": "道具提取任务已提交",
	})
}

// GeneratePropImage 单个道具生图
func (ctrl *AiController) GeneratePropImage(c *gin.Context) {
	type Req struct {
		PropID uint64 `json:"propId" binding:"required"`
	}
	var req Req
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Abort500(c, "参数错误: "+err.Error())
		return
	}

	// 查找道具
	var prop props.Props
	if err := database.DB.First(&prop, req.PropID).Error; err != nil {
		response.Abort500(c, "道具不存在")
		return
	}

	// 构造 Prompt
	prompt := ""
	if prop.ImagePrompt != nil && *prop.ImagePrompt != "" {
		prompt = *prop.ImagePrompt
	} else {
		// 如果没有专门的 ImagePrompt，使用描述或名字兜底
		desc := ""
		if prop.Description != nil {
			desc = *prop.Description
		}
		name := ""
		if prop.Name != nil {
			name = *prop.Name
		}
		prompt = fmt.Sprintf("%s, %s", name, desc)
	}

	taskService := new(services.TaskService)
	projectID := uint64(0)
	if prop.ProjectId != nil {
		projectID = *prop.ProjectId
	}

	task, err := taskService.CreatePropImageGenerationTask(projectID, req.PropID, prompt)
	if err != nil {
		response.Abort500(c, "任务创建失败: "+err.Error())
		return
	}

	response.JSON(c, gin.H{
		"code":    0,
		"data":    gin.H{"task_id": task.ID},
		"message": "生图任务已提交",
	})
}

// BatchGeneratePropImages 批量道具生图
func (ctrl *AiController) BatchGeneratePropImages(c *gin.Context) {
	type Req struct {
		PropIDs []uint64 `json:"propIds" binding:"required,min=1,max=20"`
	}
	var req Req
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Abort500(c, "参数错误: "+err.Error())
		return
	}

	type TaskResult struct {
		PropID uint64 `json:"prop_id"`
		TaskID uint64 `json:"task_id"`
	}
	var results []TaskResult
	taskService := new(services.TaskService)

	for _, propID := range req.PropIDs {
		var prop props.Props
		if err := database.DB.First(&prop, propID).Error; err != nil {
			continue
		}

		// 构造 Prompt
		prompt := ""
		if prop.ImagePrompt != nil && *prop.ImagePrompt != "" {
			prompt = *prop.ImagePrompt
		} else {
			desc := ""
			if prop.Description != nil {
				desc = *prop.Description
			}
			name := ""
			if prop.Name != nil {
				name = *prop.Name
			}
			prompt = fmt.Sprintf("%s, %s", name, desc)
		}

		projectID := uint64(0)
		if prop.ProjectId != nil {
			projectID = *prop.ProjectId
		}

		task, err := taskService.CreatePropImageGenerationTask(projectID, propID, prompt)
		if err == nil {
			results = append(results, TaskResult{
				PropID: propID,
				TaskID: task.ID,
			})
		}
	}

	response.JSON(c, gin.H{
		"code":    0,
		"message": "批量任务已提交",
		"data":    results,
	})
}

// ExtractPrompt 提取分镜图片提示词
func (ctrl *AiController) ExtractPrompt(c *gin.Context) {
	type Req struct {
		ShotID    uint64 `json:"shotId" binding:"required"`
		FrameType string `json:"frameType" binding:"required,oneof=first last key action panel"`
		Model     string `json:"model"` // 可选
	}
	var req Req
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Abort500(c, "参数错误: "+err.Error())
		return
	}

	// 1. 查找分镜信息，验证是否存在
	var shot shots.Shots
	if err := database.DB.First(&shot, req.ShotID).Error; err != nil {
		response.Abort500(c, "分镜镜头不存在")
		return
	}

	projectID := uint64(0)
	if shot.ProjectId != nil {
		projectID = *shot.ProjectId
	}

	// 2. 创建异步任务
	taskService := new(services.TaskService)
	task, err := taskService.CreateExtractFramePromptTask(projectID, req.ShotID, req.FrameType, req.Model)
	if err != nil {
		response.Abort500(c, "任务创建失败: "+err.Error())
		return
	}

	// 3. 返回任务ID供前端轮询
	response.JSON(c, gin.H{
		"code":    0,
		"message": "提示词提取任务已提交",
		"data": map[string]interface{}{
			"task_id": task.ID,
			"status":  task.Status,
		},
	})
}

// GenerateImageByPrompt 根据帧提示词生成图片
func (ctrl *AiController) GenerateImageByPrompt(c *gin.Context) {
	type Req struct {
		ShotID    uint64 `json:"shotId" binding:"required"`
		FrameType string `json:"frameType" binding:"required,oneof=first last action key"`
		Prompt    string `json:"prompt" binding:"required"` // 前端传递要生成的最终提示词
	}
	var req Req
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Abort500(c, "参数错误: "+err.Error())
		return
	}

	// 1. 获取关联项目ID
	var shot shots.Shots
	if err := database.DB.First(&shot, req.ShotID).Error; err != nil {
		response.Abort500(c, "分镜不存在")
		return
	}

	projectID := uint64(0)
	if shot.ProjectId != nil {
		projectID = *shot.ProjectId
	}

	// 2. 调用 Service 创建任务
	taskService := new(services.TaskService)
	task, err := taskService.CreateGenerateFrameImageTask(projectID, req.ShotID, req.FrameType, req.Prompt)
	if err != nil {
		response.Abort500(c, "任务启动失败: "+err.Error())
		return
	}

	// 3. 返回结果给前端进行轮询
	response.JSON(c, gin.H{
		"code":    0,
		"message": "图片生成任务已在后台运行",
		"data": map[string]interface{}{
			"task_id": task.ID,
		},
	})
}
