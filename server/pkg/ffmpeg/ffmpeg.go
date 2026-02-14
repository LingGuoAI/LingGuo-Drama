package ffmpeg

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"spiritFruit/pkg/logger"
)

// Client FFmpeg 处理器
type Client struct {
	tempDir string
}

// VideoClip 定义视频片段结构
type VideoClip struct {
	URL        string
	Duration   float64
	StartTime  float64
	EndTime    float64
	Transition map[string]interface{}
}

// MergeOptions 合并选项
type MergeOptions struct {
	OutputPath string
	Clips      []VideoClip
}

// New 创建一个新的 FFmpeg 客户端
func New() *Client {
	// 使用系统临时目录下的子目录
	tempDir := filepath.Join(os.TempDir(), "video-merge")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		logger.Error(fmt.Sprintf("FFmpeg: Failed to create temp dir: %v", err))
	}

	return &Client{
		tempDir: tempDir,
	}
}

// CleanupTempDir 清理临时目录
func (c *Client) CleanupTempDir() error {
	return os.RemoveAll(c.tempDir)
}

// ====================================================================================
// Public API: 视频合并相关
// ====================================================================================

// MergeVideos 下载、裁剪并合并视频（支持转场）
func (c *Client) MergeVideos(opts *MergeOptions) (string, error) {
	if len(opts.Clips) == 0 {
		return "", fmt.Errorf("no video clips to merge")
	}

	logger.Info(fmt.Sprintf("FFmpeg: Starting video merge process, clips_count: %d", len(opts.Clips)))

	var trimmedPaths []string
	var downloadedPaths []string

	// 1. 下载并裁剪所有片段
	for i, clip := range opts.Clips {
		// A. 下载
		downloadPath := filepath.Join(c.tempDir, fmt.Sprintf("download_%d_%d.mp4", time.Now().Unix(), i))
		localPath, err := c.downloadVideo(clip.URL, downloadPath)
		if err != nil {
			c.cleanupFiles(downloadedPaths)
			c.cleanupFiles(trimmedPaths)
			return "", fmt.Errorf("failed to download clip %d: %w", i, err)
		}
		downloadedPaths = append(downloadedPaths, localPath)

		// B. 裁剪
		trimmedPath := filepath.Join(c.tempDir, fmt.Sprintf("trimmed_%d_%d.mp4", time.Now().Unix(), i))
		err = c.trimVideo(localPath, trimmedPath, clip.StartTime, clip.EndTime)
		if err != nil {
			c.cleanupFiles(downloadedPaths)
			c.cleanupFiles(trimmedPaths)
			return "", fmt.Errorf("failed to trim clip %d: %w", i, err)
		}
		trimmedPaths = append(trimmedPaths, trimmedPath)

		logger.Info(fmt.Sprintf("FFmpeg: Clip prepared, index: %d, duration: %.2f", i, clip.EndTime-clip.StartTime))
	}

	// 清理原始下载文件
	c.cleanupFiles(downloadedPaths)

	// 2. 准备输出目录
	outputDir := filepath.Dir(opts.OutputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		c.cleanupFiles(trimmedPaths)
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	// 3. 执行合并（根据是否有转场选择不同策略）
	err := c.concatenateVideosWithTransitions(trimmedPaths, opts.Clips, opts.OutputPath)

	// 清理裁剪后的中间文件
	c.cleanupFiles(trimmedPaths)

	if err != nil {
		return "", fmt.Errorf("failed to concatenate videos: %w", err)
	}

	logger.Info(fmt.Sprintf("FFmpeg: Video merge completed successfully, output: %s", opts.OutputPath))
	return opts.OutputPath, nil
}

// ====================================================================================
// Public API: 音频处理相关
// ====================================================================================

// ExtractAudio 从视频中提取音频
func (c *Client) ExtractAudio(videoURL, outputPath string) (string, error) {
	logger.Info(fmt.Sprintf("FFmpeg: Extracting audio, url: %s, output: %s", videoURL, outputPath))

	// 1. 下载视频
	downloadPath := filepath.Join(c.tempDir, fmt.Sprintf("video_for_audio_%d.mp4", time.Now().Unix()))
	localVideoPath, err := c.downloadVideo(videoURL, downloadPath)
	if err != nil {
		return "", fmt.Errorf("failed to download video: %w", err)
	}
	defer func() {
		if err := os.Remove(localVideoPath); err != nil {
			logger.Error(fmt.Sprintf("Failed to remove temp file: %v", err))
		}
	}()

	// 2. 检查音频流
	if !c.hasAudioStream(localVideoPath) {
		logger.Warn(fmt.Sprintf("FFmpeg: No audio stream found in video %s, generating silence", videoURL))
		duration, err := c.GetVideoDuration(localVideoPath)
		if err != nil {
			return "", fmt.Errorf("failed to get video duration: %w", err)
		}
		return c.generateSilence(outputPath, duration)
	}

	// 3. 准备输出目录
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	// 4. 提取音频
	// -vn: 禁用视频, -ac: 声道数, -ab: 比特率
	cmd := exec.Command("ffmpeg",
		"-i", localVideoPath,
		"-vn",
		"-acodec", "aac",
		"-ar", "44100",
		"-ac", "2",
		"-ab", "128k",
		"-y",
		outputPath,
	)

	if output, err := cmd.CombinedOutput(); err != nil {
		logger.Error(fmt.Sprintf("FFmpeg audio extraction failed: %v, output: %s", err, string(output)))
		return "", fmt.Errorf("ffmpeg failed: %w", err)
	}

	return outputPath, nil
}

// GetVideoDuration 获取视频时长
func (c *Client) GetVideoDuration(videoPath string) (float64, error) {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		videoPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error(fmt.Sprintf("FFmpeg: Failed to get video duration: %v", err))
		return 0, fmt.Errorf("ffprobe failed: %w", err)
	}

	var duration float64
	if _, err := fmt.Sscanf(strings.TrimSpace(string(output)), "%f", &duration); err != nil {
		return 0, fmt.Errorf("parse duration failed: %w", err)
	}

	return duration, nil
}

// ====================================================================================
// Private Helpers: 核心逻辑实现
// ====================================================================================

// downloadVideo 下载视频或复制本地文件
func (c *Client) downloadVideo(url, destPath string) (string, error) {
	// 本地文件处理
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		if _, err := os.Stat(url); err == nil {
			// 修复点：copyFile 返回 error，不能直接作为 string 返回
			if err := c.copyFile(url, destPath); err != nil {
				return "", err
			}
			return destPath, nil
		}
		return "", fmt.Errorf("local file not found: %s", url)
	}

	// 网络文件下载
	logger.Info(fmt.Sprintf("FFmpeg: Downloading remote video: %s", url))
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	out, err := os.Create(destPath)
	if err != nil {
		return "", err
	}
	defer func() {
		// 修复点：处理 Close 错误
		if closeErr := out.Close(); closeErr != nil {
			logger.Error(fmt.Sprintf("Failed to close file %s: %v", destPath, closeErr))
		}
	}()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return "", err
	}

	return destPath, nil
}

// trimVideo 裁剪视频
func (c *Client) trimVideo(inputPath, outputPath string, startTime, endTime float64) error {
	// 如果没有有效的裁剪范围，则重新编码整个视频以统一格式
	isFullCopy := (startTime == 0 && endTime == 0) || endTime <= startTime

	args := []string{"-i", inputPath}

	if !isFullCopy {
		args = append(args, "-ss", fmt.Sprintf("%.2f", startTime))
		if endTime > 0 {
			args = append(args, "-to", fmt.Sprintf("%.2f", endTime))
		}
	}

	// 统一编码参数，确保合并时不冲突
	args = append(args,
		"-c:v", "libx264",
		"-preset", "fast",
		"-crf", "23",
		"-c:a", "aac",
		"-b:a", "128k",
		"-movflags", "+faststart",
		"-y",
		outputPath,
	)

	cmd := exec.Command("ffmpeg", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		logger.Error(fmt.Sprintf("FFmpeg trim failed: %v, output: %s", err, string(output)))
		return fmt.Errorf("ffmpeg failed: %w", err)
	}

	return nil
}

// concatenateVideosWithTransitions 决策合并策略
func (c *Client) concatenateVideosWithTransitions(inputPaths []string, clips []VideoClip, outputPath string) error {
	if len(inputPaths) == 0 {
		return fmt.Errorf("no input paths")
	}

	if len(inputPaths) == 1 {
		logger.Info("FFmpeg: Only one clip, copying directly")
		return c.copyFile(inputPaths[0], outputPath)
	}

	// 检查是否有转场
	hasTransitions := false
	for _, clip := range clips {
		if clip.Transition != nil && len(clip.Transition) > 0 {
			// 排除 none 类型
			if tType, ok := clip.Transition["type"].(string); ok && strings.ToLower(tType) != "none" && tType != "" {
				hasTransitions = true
				break
			}
		}
	}

	if !hasTransitions {
		logger.Info("FFmpeg: No transitions detected, using simple concatenation")
		return c.concatenateSimple(inputPaths, outputPath)
	}

	logger.Info("FFmpeg: Using complex filter (xfade) for transitions")
	return c.mergeWithXfade(inputPaths, clips, outputPath)
}

// concatenateSimple 简单拼接 (concat demuxer)
func (c *Client) concatenateSimple(inputPaths []string, outputPath string) error {
	listFile := filepath.Join(c.tempDir, fmt.Sprintf("filelist_%d.txt", time.Now().Unix()))
	defer func() {
		// 修复点：处理 Remove 错误
		if err := os.Remove(listFile); err != nil {
			logger.Error(fmt.Sprintf("Failed to remove list file: %v", err))
		}
	}()

	var content strings.Builder
	for _, path := range inputPaths {
		content.WriteString(fmt.Sprintf("file '%s'\n", path))
	}

	if err := os.WriteFile(listFile, []byte(content.String()), 0644); err != nil {
		return err
	}

	cmd := exec.Command("ffmpeg",
		"-f", "concat",
		"-safe", "0",
		"-i", listFile,
		"-c", "copy",
		"-y",
		outputPath,
	)

	if output, err := cmd.CombinedOutput(); err != nil {
		logger.Error(fmt.Sprintf("FFmpeg concat failed: %v, output: %s", err, string(output)))
		return fmt.Errorf("ffmpeg failed: %w", err)
	}

	return nil
}

// mergeWithXfade 使用 xfade 滤镜处理转场
func (c *Client) mergeWithXfade(inputPaths []string, clips []VideoClip, outputPath string) error {
	// 1. 预处理：检测分辨率和音频
	maxWidth, maxHeight := 0, 0
	hasAnyAudio := false
	audioStreams := make([]bool, len(inputPaths))

	for i, path := range inputPaths {
		w, h := c.getVideoResolution(path)
		if w > maxWidth {
			maxWidth = w
		}
		if h > maxHeight {
			maxHeight = h
		}
		audioStreams[i] = c.hasAudioStream(path)
		if audioStreams[i] {
			hasAnyAudio = true
		}
	}
	logger.Info(fmt.Sprintf("FFmpeg: Merge config - width: %d, height: %d, has_audio: %v", maxWidth, maxHeight, hasAnyAudio))

	// 2. 构建 filter_complex
	var videoFilters []string
	var inputArgs []string

	for _, path := range inputPaths {
		inputArgs = append(inputArgs, "-i", path)
	}

	// 2.1 视频流缩放与 TPad (用于转场时延长尾帧)
	for i := 0; i < len(inputPaths); i++ {
		tpadDuration := 0.0
		// 如果不是最后一个视频，且有转场，需要延长
		if i < len(clips)-1 && clips[i].Transition != nil {
			if dur, ok := clips[i].Transition["duration"].(float64); ok && dur > 0 {
				tType, _ := clips[i].Transition["type"].(string)
				if strings.ToLower(tType) != "none" {
					tpadDuration = dur
				}
			} else {
				tpadDuration = 1.0 // 默认
			}
		}

		filter := fmt.Sprintf("[%d:v]scale=%d:%d:force_original_aspect_ratio=decrease,pad=%d:%d:(ow-iw)/2:(oh-ih)/2",
			i, maxWidth, maxHeight, maxWidth, maxHeight)

		if tpadDuration > 0 {
			filter += fmt.Sprintf(",tpad=stop_mode=clone:stop_duration=%.2f", tpadDuration)
		}
		filter += fmt.Sprintf("[v%d]", i)
		videoFilters = append(videoFilters, filter)
	}

	// 2.2 XFade 转场连接
	var offset float64 = 0
	for i := 0; i < len(inputPaths)-1; i++ {
		clipDuration := clips[i].Duration
		if clips[i].EndTime > 0 {
			clipDuration = clips[i].EndTime - clips[i].StartTime
		}
		offset += clipDuration

		// 获取转场参数
		tType := "fade"
		tDuration := 1.0
		if clips[i].Transition != nil {
			if t, ok := clips[i].Transition["type"].(string); ok && t != "" {
				tType = c.mapTransitionType(t)
			}
			if d, ok := clips[i].Transition["duration"].(float64); ok && d > 0 {
				tDuration = d
			}
			if tType == "none" { // 特殊处理无转场
				tDuration = 0
			}
		}

		input1 := fmt.Sprintf("[v%d]", 0)
		if i > 0 {
			input1 = fmt.Sprintf("[vx%02d]", i-1)
		}
		input2 := fmt.Sprintf("[v%d]", i+1)
		output := fmt.Sprintf("[vx%02d]", i)
		if i == len(inputPaths)-2 {
			output = "[outv]"
		}

		xfadeCmd := fmt.Sprintf("%s%sxfade=transition=%s:duration=%.1f:offset=%.1f%s",
			input1, input2, tType, tDuration, offset, output)
		videoFilters = append(videoFilters, xfadeCmd)
	}

	// 2.3 音频处理 (如果存在音频流)
	var audioFilters []string
	if hasAnyAudio {
		// 标准化音频流：无音频的补静音，有音频的延长
		for i := 0; i < len(inputPaths); i++ {
			clipDuration := clips[i].Duration
			if clips[i].EndTime > 0 {
				clipDuration = clips[i].EndTime - clips[i].StartTime
			}

			// 计算延长时长 (pad)
			padDuration := 0.0
			if i < len(clips)-1 && clips[i].Transition != nil {
				if dur, ok := clips[i].Transition["duration"].(float64); ok && dur > 0 {
					padDuration = dur
				} else {
					padDuration = 1.0
				}
			}

			if !audioStreams[i] {
				// 生成静音流
				audioFilters = append(audioFilters,
					fmt.Sprintf("anullsrc=channel_layout=stereo:sample_rate=44100:duration=%.2f[a%d]", clipDuration+padDuration, i))
			} else {
				if padDuration > 0 {
					audioFilters = append(audioFilters, fmt.Sprintf("[%d:a]apad=pad_dur=%.2f[a%d]", i, padDuration, i))
				} else {
					audioFilters = append(audioFilters, fmt.Sprintf("[%d:a]acopy[a%d]", i, i))
				}
			}
		}

		// 音频交叉淡入淡出 (acrossfade)
		for i := 0; i < len(inputPaths)-1; i++ {
			input1 := fmt.Sprintf("[a%d]", 0)
			if i > 0 {
				input1 = fmt.Sprintf("[ax%02d]", i-1)
			}
			input2 := fmt.Sprintf("[a%d]", i+1)
			output := fmt.Sprintf("[ax%02d]", i)
			if i == len(inputPaths)-2 {
				output = "[outa]"
			}

			// 获取转场时长
			tDuration := 1.0
			if clips[i].Transition != nil {
				if d, ok := clips[i].Transition["duration"].(float64); ok && d > 0 {
					tDuration = d
				}
				if t, ok := clips[i].Transition["type"].(string); ok && t == "none" {
					tDuration = 0
				}
			}

			crossfadeCmd := fmt.Sprintf("%s%sacrossfade=d=%.2f:c1=tri:c2=tri%s",
				input1, input2, tDuration, output)
			audioFilters = append(audioFilters, crossfadeCmd)
		}
	}

	// 3. 组装命令
	fullFilter := strings.Join(videoFilters, ";")
	if hasAnyAudio {
		fullFilter += ";" + strings.Join(audioFilters, ";")
	}

	args := append(inputArgs, "-filter_complex", fullFilter, "-map", "[outv]")
	if hasAnyAudio {
		args = append(args, "-map", "[outa]")
		args = append(args, "-c:a", "aac", "-b:a", "128k")
	}
	args = append(args, "-c:v", "libx264", "-preset", "medium", "-crf", "23", "-y", outputPath)

	logger.Info(fmt.Sprintf("FFmpeg: Executing Xfade Merge, filter length: %d", len(fullFilter)))

	cmd := exec.Command("ffmpeg", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		logger.Error(fmt.Sprintf("FFmpeg xfade failed: %v, output: %s", err, string(output)))
		return fmt.Errorf("ffmpeg xfade failed: %w", err)
	}

	return nil
}

// generateSilence 生成静音文件
func (c *Client) generateSilence(outputPath string, duration float64) (string, error) {
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return "", err
	}

	cmd := exec.Command("ffmpeg",
		"-f", "lavfi",
		"-i", "anullsrc=channel_layout=stereo:sample_rate=44100",
		"-t", fmt.Sprintf("%.2f", duration),
		"-acodec", "aac",
		"-ab", "128k",
		"-y",
		outputPath,
	)

	if output, err := cmd.CombinedOutput(); err != nil {
		logger.Error(fmt.Sprintf("FFmpeg silence generation failed: %v, output: %s", err, string(output)))
		return "", fmt.Errorf("failed to generate silence: %w", err)
	}
	return outputPath, nil
}

// mapTransitionType 映射转场名称
func (c *Client) mapTransitionType(t string) string {
	t = strings.ToLower(t)
	// 这里可以添加更多映射，目前默认直接返回，如果需要兼容前端命名
	switch t {
	case "fade", "fadein", "fadeout":
		return "fade"
	case "none":
		return "none"
	default:
		// FFmpeg 支持很多类型，默认直接透传，如果不支持会报错回退到 fade
		return t
	}
}

// 辅助工具方法
func (c *Client) copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	// 修复点：处理 Close 错误
	defer func() {
		if err := destFile.Close(); err != nil {
			logger.Error(fmt.Sprintf("Failed to close dest file %s: %v", dst, err))
		}
	}()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return err
	}
	return nil
}

func (c *Client) cleanupFiles(paths []string) {
	for _, path := range paths {
		// 修复点：处理 Remove 错误
		if err := os.Remove(path); err != nil {
			// 可以选择忽略错误或打印日志，这里选择打印日志
			logger.Warn(fmt.Sprintf("Failed to cleanup file %s: %v", path, err))
		}
	}
}

func (c *Client) hasAudioStream(path string) bool {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-select_streams", "a:0",
		"-show_entries", "stream=codec_type",
		"-of", "default=noprint_wrappers=1:nokey=1",
		path,
	)
	output, _ := cmd.CombinedOutput()
	return strings.TrimSpace(string(output)) == "stream_codec_type=audio" // ffprobe output format varies
}

func (c *Client) getVideoResolution(path string) (int, int) {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-select_streams", "v:0",
		"-show_entries", "stream=width,height",
		"-of", "csv=p=0",
		path,
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 1920, 1080
	}
	var w, h int
	fmt.Sscanf(strings.TrimSpace(string(output)), "%d,%d", &w, &h)
	if w == 0 || h == 0 {
		return 1920, 1080
	}
	return w, h
}
