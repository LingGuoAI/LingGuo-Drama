// Package v1 处理业务逻辑,  控制器 v1
package v1

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

// BaseADMINController 基础控制器
type BaseADMINController struct {
}

// GetAdminID 获取当前管理员ID
func (ctrl *BaseADMINController) GetAdminID(c *gin.Context) uint64 {
	adminIDStr := c.GetString("current_admin_id")
	if adminIDStr == "" {
		return 0
	}
	adminID, _ := strconv.ParseUint(adminIDStr, 10, 64)
	return adminID
}
