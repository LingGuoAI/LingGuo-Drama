package v1

import (
	"fmt"
	"time"
	"spiritFruit/pkg/response"
	"spiritFruit/pkg/database"

	"github.com/gin-gonic/gin"
)

// StatisticsController 统计控制器
type StatisticsController struct {
	BaseADMINController
}

// StatisticsCard 统计卡片数据结构
type StatisticsCard struct {
	Title     string  `json:"title"`     // 标题
	Number    string  `json:"number"`    // 数量
	UpTrend   string  `json:"upTrend,omitempty"`   // 上升趋势
	DownTrend string  `json:"downTrend,omitempty"` // 下降趋势
	LeftType  string  `json:"leftType"`  // 左侧图标类型
}

// StatisticsData 获取统计数据
// @Summary 获取统计数据
// @Description 获取各个模块的统计数据，包括总量和增长趋势
// @Tags 统计管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]StatisticsCard} "统计数据"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /admin/v1/statistics [get]
func (s *StatisticsController) StatisticsData(c *gin.Context) {
	var statisticsCards []StatisticsCard

	// 获取当前时间和一周前的时间
	now := time.Now()
	oneWeekAgo := now.AddDate(0, 0, -7)

	
    // 角色统计
    charactersCard := s.getTableStatistics("characters", "角色总量(个)", "echarts-line", now, oneWeekAgo)
    statisticsCards = append(statisticsCards, charactersCard)
    
    
    // 短剧项目统计
    projectsCard := s.getTableStatistics("projects", "短剧项目总量(个)", "echarts-bar", now, oneWeekAgo)
    statisticsCards = append(statisticsCards, projectsCard)
    
    
    // 剧本统计
    scriptsCard := s.getTableStatistics("scripts", "剧本总量(个)", "icon-usergroup", now, oneWeekAgo)
    statisticsCards = append(statisticsCards, scriptsCard)
    
    
    // 镜头表统计
    shotsCard := s.getTableStatistics("shots", "镜头表总量(个)", "icon-file-paste", now, oneWeekAgo)
    statisticsCards = append(statisticsCards, shotsCard)
    
    

	response.JSON(c, gin.H{
        "code":    0,
        "message": "success",
        "data":    statisticsCards,
    })
}

// getTableStatistics 获取单个表的统计数据
func (s *StatisticsController) getTableStatistics(tableName, title, iconType string, now, oneWeekAgo time.Time) StatisticsCard {
	// 获取总数量
	var totalCount int64
	database.DB.Table(tableName).Count(&totalCount)

	// 获取一周前的数量
	var weekAgoCount int64
	database.DB.Table(tableName).Where("created_at <= ?", oneWeekAgo).Count(&weekAgoCount)

	// 计算增长数量和增长率
	growthCount := totalCount - weekAgoCount
	var trendStr string
	var upTrend, downTrend string

	if weekAgoCount > 0 {
		growthRate := float64(growthCount) / float64(weekAgoCount) * 100
		trendStr = fmt.Sprintf("%.1f%%", growthRate)

		if growthCount >= 0 {
			upTrend = trendStr
		} else {
			downTrend = trendStr
		}
	} else if growthCount > 0 {
		upTrend = "100.0%"
	}

	return StatisticsCard{
		Title:     title,
		Number:    fmt.Sprintf("%d", totalCount),
		UpTrend:   upTrend,
		DownTrend: downTrend,
		LeftType:  iconType,
	}
}

// GetTableTrend 获取指定表的趋势数据（可选方法，用于图表展示）
// @Summary 获取表趋势数据
// @Description 获取指定表最近7天的数据趋势
// @Tags 统计管理
// @Accept json
// @Produce json
// @Param table_name query string true "表名"
// @Success 200 {object} response.Response{data=[]map[string]interface{}} "趋势数据"
// @Failure 400 {object} response.Response "参数错误"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /admin/v1/statistics/trend [get]
func (s *StatisticsController) GetTableTrend(c *gin.Context) {
	tableName := c.Query("table_name")
	if tableName == "" {
		return
	}

	// 获取最近7天的数据
	var trendData []map[string]interface{}
	now := time.Now()

	for i := 6; i >= 0; i-- {
		date := now.AddDate(0, 0, -i)
		startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
		endOfDay := startOfDay.Add(24 * time.Hour)

		var count int64
		database.DB.Table(tableName).
			Where("created_at >= ? AND created_at < ?", startOfDay, endOfDay).
			Count(&count)

		trendData = append(trendData, map[string]interface{}{
			"date":  startOfDay.Format("2006-01-02"),
			"count": count,
		})
	}

	response.JSON(c, gin.H{
        "code":    0,
        "message": "success",
        "data":    trendData,
    })
}

// GetDetailStatistics 获取详细统计信息（可选方法）
// @Summary 获取详细统计信息
// @Description 获取更详细的统计信息，包括各种维度的数据分析
// @Tags 统计管理
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=map[string]interface{}} "详细统计数据"
// @Failure 500 {object} response.Response "服务器错误"
// @Router /admin/v1/statistics/detail [get]
func (s *StatisticsController) GetDetailStatistics(c *gin.Context) {
	detailStats := make(map[string]interface{})

	// 获取今日、昨日、本周、上周、本月、上月的统计数据
	now := time.Now()

	// 今日统计
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	tomorrowStart := todayStart.Add(24 * time.Hour)

	// 昨日统计
	yesterdayStart := todayStart.Add(-24 * time.Hour)

	// 本周统计（周一开始）
	weekday := int(now.Weekday())
	if weekday == 0 { // 周日
		weekday = 7
	}
	thisWeekStart := todayStart.AddDate(0, 0, -(weekday-1))

	// 上周统计
	lastWeekStart := thisWeekStart.AddDate(0, 0, -7)

	// 本月统计
	thisMonthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	// 上月统计
	lastMonthStart := thisMonthStart.AddDate(0, -1, 0)

	
	// 角色详细统计
	detailStats["c"] = map[string]interface{}{
		"today":      s.getCountBetween("characters", todayStart, tomorrowStart),
		"yesterday":  s.getCountBetween("characters", yesterdayStart, todayStart),
		"this_week":  s.getCountBetween("characters", thisWeekStart, now),
		"last_week":  s.getCountBetween("characters", lastWeekStart, thisWeekStart),
		"this_month": s.getCountBetween("characters", thisMonthStart, now),
		"last_month": s.getCountBetween("characters", lastMonthStart, thisMonthStart),
	}
	
	// 短剧项目详细统计
	detailStats["p"] = map[string]interface{}{
		"today":      s.getCountBetween("projects", todayStart, tomorrowStart),
		"yesterday":  s.getCountBetween("projects", yesterdayStart, todayStart),
		"this_week":  s.getCountBetween("projects", thisWeekStart, now),
		"last_week":  s.getCountBetween("projects", lastWeekStart, thisWeekStart),
		"this_month": s.getCountBetween("projects", thisMonthStart, now),
		"last_month": s.getCountBetween("projects", lastMonthStart, thisMonthStart),
	}
	
	// 剧本详细统计
	detailStats["s"] = map[string]interface{}{
		"today":      s.getCountBetween("scripts", todayStart, tomorrowStart),
		"yesterday":  s.getCountBetween("scripts", yesterdayStart, todayStart),
		"this_week":  s.getCountBetween("scripts", thisWeekStart, now),
		"last_week":  s.getCountBetween("scripts", lastWeekStart, thisWeekStart),
		"this_month": s.getCountBetween("scripts", thisMonthStart, now),
		"last_month": s.getCountBetween("scripts", lastMonthStart, thisMonthStart),
	}
	
	// 镜头表详细统计
	detailStats["s"] = map[string]interface{}{
		"today":      s.getCountBetween("shots", todayStart, tomorrowStart),
		"yesterday":  s.getCountBetween("shots", yesterdayStart, todayStart),
		"this_week":  s.getCountBetween("shots", thisWeekStart, now),
		"last_week":  s.getCountBetween("shots", lastWeekStart, thisWeekStart),
		"this_month": s.getCountBetween("shots", thisMonthStart, now),
		"last_month": s.getCountBetween("shots", lastMonthStart, thisMonthStart),
	}
	

	response.JSON(c, gin.H{
        "code":    0,
        "message": "success",
        "data":    detailStats,
    })
}

// getCountBetween 获取指定时间范围内的数据数量
func (s *StatisticsController) getCountBetween(tableName string, startTime, endTime time.Time) int64 {
	var count int64
	database.DB.Table(tableName).
		Where("created_at >= ? AND created_at < ?", startTime, endTime).
		Count(&count)
	return count
}