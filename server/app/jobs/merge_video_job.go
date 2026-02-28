package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/hibiken/asynq"

	"spiritFruit/app/models/async_tasks"
	// "spiritFruit/app/models/video_merges" // 根据您的实际表名引入
	// "spiritFruit/app/models/scripts"

	myAsynq "spiritFruit/pkg/asynq"
	"spiritFruit/pkg/config"
	"spiritFruit/pkg/console"
	"spiritFruit/pkg/database"
	"spiritFruit/pkg/ffmpeg" // 引入您给我的 ffmpeg 扩展包
)

// HandleMergeVideoTask 处理视频合并异步任务
func HandleMergeVideoTask(ctx context.Context, t *asynq.Task) error {
	var p myAsynq.MergeVideoPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json unmarshal failed: %v", err)
	}

	// 1. 获取并标记任务开始
	var taskModel async_tasks.AsyncTask
	if err := database.DB.First(&taskModel, p.AsyncTaskID).Error; err != nil {
		return nil
	}
	taskModel.MarkAsProcessing()
	taskModel.UpdateProgress(10)

	// 同步更新 Merge 表状态
	// database.DB.Model(&video_merges.VideoMerge{}).Where("id = ?", p.MergeID).Update("status", "processing")

	console.Success(fmt.Sprintf("任务[%d] - 开始合并视频, 共有 %d 个片段", p.AsyncTaskID, len(p.Clips)))

	// 2. 准备 FFmpeg 参数
	ffmpegClient := ffmpeg.New()
	defer ffmpegClient.CleanupTempDir() // 任务结束自动清理临时目录

	// 构造输出路径 (public/uploads/videos/merged/xxx.mp4)
	uploadPath := config.GetString("app.upload_path", "public/uploads")
	fileName := fmt.Sprintf("merged_%d_%d.mp4", p.ProjectID, time.Now().Unix())
	outputDir := filepath.Join(uploadPath, "videos", "merged")
	absoluteOutputPath := filepath.Join(outputDir, fileName)

	// 转换 Clip 结构
	var ffmpegClips []ffmpeg.VideoClip
	for _, clip := range p.Clips {
		// 补全相对路径为 URL 或本地绝对路径（您的 ffmpeg 包能自动处理）
		url := clip.URL
		if !strings.HasPrefix(url, "http") && !filepath.IsAbs(url) {
			// 如果是相对路径，例如 /uploads/videos/xxx.mp4，我们转为本地绝对路径或加上APP_URL
			url = filepath.Join(config.GetString("app.public_path", "public"), strings.TrimPrefix(url, "/"))
		}

		ffmpegClips = append(ffmpegClips, ffmpeg.VideoClip{
			URL:        url,
			Duration:   clip.Duration,
			StartTime:  clip.StartTime,
			EndTime:    clip.EndTime,
			Transition: clip.Transition,
		})
	}

	taskModel.UpdateProgress(30)

	// 3. 执行视频合并 (这里会非常耗时)
	console.Success("调用 FFmpeg 开始硬核合并...")
	mergeOpts := &ffmpeg.MergeOptions{
		OutputPath: absoluteOutputPath,
		Clips:      ffmpegClips,
	}

	finalPath, err := ffmpegClient.MergeVideos(mergeOpts)
	if err != nil {
		console.Error(fmt.Sprintf("FFmpeg 合并失败: %v", err))
		taskModel.MarkAsFailed(err)
		// database.DB.Model(&video_merges.VideoMerge{}).Where("id = ?", p.MergeID).Update("status", "failed")
		return err
	}

	taskModel.UpdateProgress(90)

	// 4. 获取合成后的最终时长
	finalDuration, _ := ffmpegClient.GetVideoDuration(finalPath)

	// 5. 组装返回给前端的相对 URL
	// 去掉 public，变为 /uploads/videos/merged/xxx.mp4
	relativePath := "/" + strings.TrimPrefix(finalPath, config.GetString("app.public_path", "public/"))
	relativePath = filepath.ToSlash(relativePath) // 兼容 windows 反斜杠

	// 6. 更新数据库状态
	// database.DB.Model(&video_merges.VideoMerge{}).Where("id = ?", p.MergeID).Updates(map[string]interface{}{
	//	 "status":     "completed",
	//	 "merged_url": relativePath,
	//	 "duration":   finalDuration,
	// })

	// 可选：更新到剧本表/集数表
	// database.DB.Model(&scripts.Scripts{}).Where("id = ?", p.EpisodeID).Update("video_url", relativePath)

	// 7. 任务彻底成功
	taskModel.UpdateProgress(100)
	resData := map[string]interface{}{
		"url":      relativePath,
		"duration": finalDuration,
	}
	resBytes, _ := json.Marshal(resData)
	taskModel.MarkAsSuccess(string(resBytes))

	console.Success(fmt.Sprintf("任务[%d] - 视频合并完成! 输出路径: %s", p.AsyncTaskID, relativePath))
	return nil
}
