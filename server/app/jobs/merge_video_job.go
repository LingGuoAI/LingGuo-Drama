package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"

	"spiritFruit/app/models/async_tasks"
	"spiritFruit/app/models/shot_video_merge"

	myAsynq "spiritFruit/pkg/asynq"
	"spiritFruit/pkg/console"
	"spiritFruit/pkg/database"
	"spiritFruit/pkg/ffmpeg"
)

// HandleMergeVideoTask 处理视频合并异步任务
func HandleMergeVideoTask(ctx context.Context, t *asynq.Task) error {
	var p myAsynq.MergeVideoPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("json unmarshal failed: %v", err)
	}

	// 1. 获取并标记异步任务开始
	var taskModel async_tasks.AsyncTask
	if err := database.DB.First(&taskModel, p.AsyncTaskID).Error; err != nil {
		return nil
	}
	taskModel.MarkAsProcessing()
	taskModel.UpdateProgress(10)

	console.Success(fmt.Sprintf("任务[%d] - 开始合并视频, 共有 %d 个片段", p.AsyncTaskID, len(p.Clips)))

	//  2. 提前在数据库中创建合并记录，状态标记为 processing
	statusProcessing := "processing"
	mergeRecord := shot_video_merge.ShotVideoMerge{
		ProjectId: &p.ProjectID,
		ScriptId:  &p.EpisodeID,
		Title:     &p.Title,
		TaskId:    &p.AsyncTaskID,
		Status:    &statusProcessing, // 初始状态为处理中
	}
	if err := database.DB.Create(&mergeRecord).Error; err != nil {
		taskModel.MarkAsFailed(fmt.Errorf("create merge record failed: %v", err))
		return err
	}

	// 3. 准备 FFmpeg 参数与客户端
	taskModel.UpdateProgress(20)
	ffmpegClient := ffmpeg.New()
	defer ffmpegClient.CleanupTempDir() // 任务结束自动清理临时目录

	uploadDir := "uploads/videos/" + time.Now().Format("2006/01/02")
	fileName := uuid.New().String() + ".mp4"

	localPath := path.Join(uploadDir, fileName)
	absoluteOutputPath, _ := filepath.Abs(filepath.FromSlash(localPath))

	// 确保存放的目录真正生成
	if err := os.MkdirAll(filepath.Dir(absoluteOutputPath), 0755); err != nil {
		taskModel.MarkAsFailed(fmt.Errorf("create dir failed: %v", err))
		return err
	}

	// 4. 转换并校验 Clip 结构
	taskModel.UpdateProgress(30)
	var ffmpegClips []ffmpeg.VideoClip
	for _, clip := range p.Clips {
		url := clip.URL

		// 补全相对路径为本地绝对路径
		if !strings.HasPrefix(url, "http") && !filepath.IsAbs(url) {
			url = strings.TrimPrefix(url, "/")
			absUrl, err := filepath.Abs(filepath.FromSlash(url))
			if err == nil {
				url = absUrl
			}
		}

		ffmpegClips = append(ffmpegClips, ffmpeg.VideoClip{
			URL:        url,
			Duration:   clip.Duration,
			StartTime:  clip.StartTime,
			EndTime:    clip.EndTime,
			Transition: clip.Transition,
		})
	}

	// 5. 执行视频合并
	taskModel.UpdateProgress(40)
	console.Success("调用 FFmpeg 开始硬核合并...")
	mergeOpts := &ffmpeg.MergeOptions{
		OutputPath: absoluteOutputPath, // 传递物理绝对路径
		Clips:      ffmpegClips,
	}

	_, err := ffmpegClient.MergeVideos(mergeOpts)
	if err != nil {
		console.Error(fmt.Sprintf("FFmpeg 合并失败: %v", err))
		taskModel.MarkAsFailed(err)

		// 任务失败：更新合并记录的状态为 failed，并写入 ErrorMsg
		statusFailed := "failed"
		errMsg := err.Error()
		database.DB.Model(&shot_video_merge.ShotVideoMerge{}).
			Where("id = ?", mergeRecord.ID).
			Updates(map[string]interface{}{
				"status":    &statusFailed,
				"error_msg": &errMsg,
			})

		return err
	}

	// 6. 获取合成后的最终时长
	taskModel.UpdateProgress(90)
	finalDuration, _ := ffmpegClient.GetVideoDuration(absoluteOutputPath)
	durationInt := int(finalDuration)

	//  7. 任务成功：更新合并记录的状态为 completed，并填入路径和时长
	statusCompleted := "completed"
	err = database.DB.Model(&shot_video_merge.ShotVideoMerge{}).
		Where("id = ?", mergeRecord.ID).
		Updates(map[string]interface{}{
			"status":     &statusCompleted,
			"merged_url": &localPath,
			"duration":   &durationInt,
			"error_msg":  nil, // 清空错误
		}).Error

	if err != nil {
		taskModel.MarkAsFailed(fmt.Errorf("db update failed: %v", err))
		return err
	}

	// 8. 标记异步任务彻底成功
	taskModel.UpdateProgress(100)
	resData := map[string]interface{}{
		"url":      localPath,
		"duration": finalDuration,
	}
	resBytes, _ := json.Marshal(resData)
	taskModel.MarkAsSuccess(string(resBytes))

	console.Success(fmt.Sprintf("任务[%d] - 视频合并彻底完成! 输出路径: %s", p.AsyncTaskID, localPath))
	return nil
}
