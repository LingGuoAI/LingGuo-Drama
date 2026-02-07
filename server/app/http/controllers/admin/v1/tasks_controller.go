package v1

import (
	"github.com/gin-gonic/gin"
	"spiritFruit/app/models/async_tasks"
	"spiritFruit/pkg/database"
	"spiritFruit/pkg/response"
)

type TasksController struct {
	BaseADMINController
}

// Show 获取任务详情 (用于轮询)
func (ctrl *TasksController) Show(c *gin.Context) {
	id := c.Param("id")
	var task async_tasks.AsyncTask

	// 查询任务表
	if err := database.DB.First(&task, id).Error; err != nil {
		response.Abort404(c, "任务不存在")
		return
	}

	// 返回包含 Process(进度) 和 Status(状态) 的数据
	response.Data(c, task)
}
