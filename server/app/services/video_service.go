package services

import (
	"encoding/json"
	"fmt"
	"strconv"

	"spiritFruit/app/models/async_tasks"
	"spiritFruit/app/models/shots"

	"spiritFruit/pkg/asynq"
	"spiritFruit/pkg/console"
	"spiritFruit/pkg/database"
)

type TimelineClipReq struct {
	AssetID    interface{}            `json:"asset_id"`
	ShotID     string                 `json:"shotId"`
	Order      int                    `json:"order"`
	StartTime  float64                `json:"start_time"`
	EndTime    float64                `json:"end_time"`
	Duration   float64                `json:"duration"`
	Transition map[string]interface{} `json:"transition"`
}

type FinalizeEpisodeReq struct {
	ProjectID     uint64            `json:"projectId" binding:"required"`
	EpisodeNumber uint64            `json:"episodeNumber" binding:"required"`
	Clips         []TimelineClipReq `json:"clips"`
}

type VideoService struct{}

// FinalizeEpisode 参数类型改为使用本包的 FinalizeEpisodeReq
func (s *VideoService) FinalizeEpisode(req FinalizeEpisodeReq) (map[string]interface{}, error) {
	var shotList []shots.Shots
	// 🔴 这里改为 req.EpisodeNumber
	database.DB.Where("project_id = ? AND episode_number = ?", req.ProjectID, req.EpisodeNumber).Order("sequence_no asc").Find(&shotList)

	if len(shotList) == 0 {
		return nil, fmt.Errorf("该集数下没有找到任何分镜")
	}

	sceneMap := make(map[string]shots.Shots)
	for _, shot := range shotList {
		sceneMap[fmt.Sprintf("%d", shot.ID)] = shot
	}

	var mergeClips []asynq.MergeClip
	var skippedScenes []int

	if len(req.Clips) > 0 {
		console.Success(fmt.Sprintf("使用前端时间线数据合成，片段数: %d", len(req.Clips)))

		for _, clip := range req.Clips {
			var videoURL string

			if videoURL == "" && clip.ShotID != "" {
				if scene, exists := sceneMap[clip.ShotID]; exists {
					if scene.VideoUrl != nil && *scene.VideoUrl != "" {
						videoURL = *scene.VideoUrl
					}
				}
			}

			if videoURL == "" {
				shotIDInt, _ := strconv.Atoi(clip.ShotID)
				skippedScenes = append(skippedScenes, shotIDInt)
				continue
			}

			mergeClips = append(mergeClips, asynq.MergeClip{
				URL:        videoURL,
				Duration:   clip.Duration,
				StartTime:  clip.StartTime,
				EndTime:    clip.EndTime,
				Transition: clip.Transition,
			})
		}
	} else {
		fmt.Println("[INFO] 无时间线数据，按分镜默认顺序拼接")

		for _, scene := range shotList {
			var videoURL string

			if scene.VideoUrl != nil && *scene.VideoUrl != "" {
				videoURL = *scene.VideoUrl
			}

			if videoURL == "" {
				seqNo := 0
				if scene.SequenceNo != nil {
					seqNo = int(*scene.SequenceNo)
				}
				skippedScenes = append(skippedScenes, seqNo)
				continue
			}

			duration := 3.0
			if scene.DurationMs != nil {
				duration = float64(*scene.DurationMs) / 1000.0
			}

			mergeClips = append(mergeClips, asynq.MergeClip{
				URL:      videoURL,
				Duration: duration,
			})
		}
	}

	if len(mergeClips) == 0 {
		return nil, fmt.Errorf("没有找到任何可用的视频片段用于合成")
	}

	mergeRecordID := uint64(1) // 您的实际 Merge ID

	payload := asynq.MergeVideoPayload{
		MergeID:   mergeRecordID,
		ProjectID: req.ProjectID,
		EpisodeID: req.EpisodeNumber,
		Title:     fmt.Sprintf("项目%d-第%d集合成", req.ProjectID, req.EpisodeNumber),
		Clips:     mergeClips,
	}

	payloadBytes, _ := json.Marshal(payload)
	task := async_tasks.AsyncTask{
		ProjectID: req.ProjectID,
		RelID:     mergeRecordID,
		Type:      asynq.TypeMergeVideo,
		Status:    async_tasks.StatusPending,
		Payload:   string(payloadBytes),
	}
	database.DB.Create(&task)

	payload.AsyncTaskID = task.ID
	_, err := asynq.EnqueueMergeVideo(payload)
	if err != nil {
		task.MarkAsFailed(err)
		return nil, err
	}

	result := map[string]interface{}{
		"message":      "视频合成任务已创建，正在后台处理",
		"merge_id":     mergeRecordID,
		"task_id":      task.ID,
		"scenes_count": len(mergeClips),
	}

	if len(skippedScenes) > 0 {
		result["skipped_scenes"] = skippedScenes
		result["warning"] = fmt.Sprintf("已跳过 %d 个未生成视频的场景", len(skippedScenes))
	}

	return result, nil
}
