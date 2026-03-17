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
	tempDir := filepath.Join(os.TempDir(), "video-merge")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		logger.Error(fmt.Sprintf("FFmpeg: Failed to create temp dir: %v", err))
	}

	return &Client{
		tempDir: tempDir,
	}
}

func (c *Client) CleanupTempDir() error {
	return os.RemoveAll(c.tempDir)
}

// GetVideoDuration 获取视频时长 (Job 脚本依赖此方法)
func (c *Client) GetVideoDuration(videoPath string) (float64, error) {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		videoPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("ffprobe failed: %w", err)
	}

	var duration float64
	if _, err := fmt.Sscanf(strings.TrimSpace(string(output)), "%f", &duration); err != nil {
		return 0, fmt.Errorf("parse duration failed: %w", err)
	}

	return duration, nil
}

// MergeVideos 下载、裁剪并合并视频
func (c *Client) MergeVideos(opts *MergeOptions) (string, error) {
	if len(opts.Clips) == 0 {
		return "", fmt.Errorf("no video clips to merge")
	}

	logger.Info(fmt.Sprintf("FFmpeg: Starting video merge process, clips_count: %d", len(opts.Clips)))

	var trimmedPaths []string
	var downloadedPaths []string

	for i, clip := range opts.Clips {
		downloadPath := filepath.Join(c.tempDir, fmt.Sprintf("download_%d_%d.mp4", time.Now().Unix(), i))
		localPath, err := c.downloadVideo(clip.URL, downloadPath)
		if err != nil {
			c.cleanupFiles(downloadedPaths)
			c.cleanupFiles(trimmedPaths)
			return "", fmt.Errorf("failed to download clip %d: %w", i, err)
		}
		downloadedPaths = append(downloadedPaths, localPath)

		trimmedPath := filepath.Join(c.tempDir, fmt.Sprintf("trimmed_%d_%d.mp4", time.Now().Unix(), i))
		err = c.trimAndNormalizeVideo(localPath, trimmedPath, clip.StartTime, clip.EndTime)
		if err != nil {
			c.cleanupFiles(downloadedPaths)
			c.cleanupFiles(trimmedPaths)
			return "", fmt.Errorf("failed to trim clip %d: %w", i, err)
		}
		trimmedPaths = append(trimmedPaths, trimmedPath)
	}

	c.cleanupFiles(downloadedPaths)

	outputDir := filepath.Dir(opts.OutputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		c.cleanupFiles(trimmedPaths)
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	err := c.concatenateVideosWithTransitions(trimmedPaths, opts.Clips, opts.OutputPath)
	c.cleanupFiles(trimmedPaths)

	if err != nil {
		return "", fmt.Errorf("failed to concatenate videos: %w", err)
	}

	return opts.OutputPath, nil
}

// trimAndNormalizeVideo 裁剪并标准化（统一音频参数解决无声问题）
func (c *Client) trimAndNormalizeVideo(inputPath, outputPath string, startTime, endTime float64) error {
	args := []string{"-i", inputPath}
	if endTime > startTime && endTime > 0 {
		args = append(args, "-ss", fmt.Sprintf("%.2f", startTime), "-to", fmt.Sprintf("%.2f", endTime))
	}

	args = append(args,
		"-c:v", "libx264",
		"-pix_fmt", "yuv420p",
		"-c:a", "aac",
		"-ar", "44100",
		"-ac", "2",
		"-y",
		outputPath,
	)

	cmd := exec.Command("ffmpeg", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("trim failed: %s", string(output))
	}
	return nil
}

func (c *Client) concatenateVideosWithTransitions(inputPaths []string, clips []VideoClip, outputPath string) error {
	if len(inputPaths) == 1 {
		return c.copyFile(inputPaths[0], outputPath)
	}

	hasTransitions := false
	for _, clip := range clips {
		if tType, ok := clip.Transition["type"].(string); ok && strings.ToLower(tType) != "none" && tType != "" {
			hasTransitions = true
			break
		}
	}

	if !hasTransitions {
		return c.concatenateSimple(inputPaths, outputPath)
	}
	return c.mergeWithXfade(inputPaths, clips, outputPath)
}

func (c *Client) concatenateSimple(inputPaths []string, outputPath string) error {
	listFile := filepath.Join(c.tempDir, fmt.Sprintf("filelist_%d.txt", time.Now().UnixNano()))
	defer os.Remove(listFile)

	var content strings.Builder
	for _, path := range inputPaths {
		content.WriteString(fmt.Sprintf("file '%s'\n", path))
	}
	os.WriteFile(listFile, []byte(content.String()), 0644)

	cmd := exec.Command("ffmpeg",
		"-f", "concat",
		"-safe", "0",
		"-i", listFile,
		"-c:v", "libx264",
		"-c:a", "aac",
		"-ar", "44100",
		"-ac", "2",
		"-y",
		outputPath,
	)

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("concat failed: %s", string(output))
	}
	return nil
}

func (c *Client) mergeWithXfade(inputPaths []string, clips []VideoClip, outputPath string) error {
	maxWidth, maxHeight := 1280, 720
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

	var videoFilters []string
	var audioFilters []string
	var inputArgs []string

	for i, path := range inputPaths {
		inputArgs = append(inputArgs, "-i", path)
		tpadDur := 0.0
		if i < len(clips)-1 {
			if d, ok := clips[i].Transition["duration"].(float64); ok {
				tpadDur = d
			} else {
				tpadDur = 1.0
			}
		}

		vFilter := fmt.Sprintf("[%d:v]scale=%d:%d:force_original_aspect_ratio=decrease,pad=%d:%d:(ow-iw)/2:(oh-ih)/2", i, maxWidth, maxHeight, maxWidth, maxHeight)
		if tpadDur > 0 {
			vFilter += fmt.Sprintf(",tpad=stop_mode=clone:stop_duration=%.2f", tpadDur)
		}
		videoFilters = append(videoFilters, vFilter+fmt.Sprintf("[v%d]", i))

		clipDur := clips[i].Duration
		if clips[i].EndTime > 0 {
			clipDur = clips[i].EndTime - clips[i].StartTime
		}

		if !audioStreams[i] {
			audioFilters = append(audioFilters, fmt.Sprintf("anullsrc=channel_layout=stereo:sample_rate=44100:duration=%.2f[a%d]", clipDur+tpadDur, i))
		} else {
			aFilter := fmt.Sprintf("[%d:a]aresample=44100,aformat=sample_fmts=fltp:sample_rates=44100:channel_layouts=stereo", i)
			if tpadDur > 0 {
				aFilter += fmt.Sprintf(",apad=pad_dur=%.2f", tpadDur)
			}
			audioFilters = append(audioFilters, aFilter+fmt.Sprintf("[a%d]", i))
		}
	}

	lastV := "[v0]"
	var offset float64 = 0
	for i := 0; i < len(inputPaths)-1; i++ {
		dur := clips[i].Duration
		if clips[i].EndTime > 0 {
			dur = clips[i].EndTime - clips[i].StartTime
		}
		offset += dur
		tType := c.mapTransitionType(fmt.Sprintf("%v", clips[i].Transition["type"]))
		tDur := 1.0
		if d, ok := clips[i].Transition["duration"].(float64); ok && d > 0 {
			tDur = d
		}
		outV := fmt.Sprintf("[vx%d]", i)
		if i == len(inputPaths)-2 {
			outV = "[outv]"
		}
		videoFilters = append(videoFilters, fmt.Sprintf("%s[v%d]xfade=transition=%s:duration=%.2f:offset=%.2f%s", lastV, i+1, tType, tDur, offset, outV))
		lastV = fmt.Sprintf("[vx%d]", i)
	}

	lastA := "[a0]"
	for i := 0; i < len(inputPaths)-1; i++ {
		tDur := 1.0
		if d, ok := clips[i].Transition["duration"].(float64); ok && d > 0 {
			tDur = d
		}
		outA := fmt.Sprintf("[ax%d]", i)
		if i == len(inputPaths)-2 {
			outA = "[outa]"
		}
		audioFilters = append(audioFilters, fmt.Sprintf("%s[a%d]acrossfade=d=%.2f:c1=tri:c2=tri%s", lastA, i+1, tDur, outA))
		lastA = fmt.Sprintf("[ax%d]", i)
	}

	fullFilter := strings.Join(videoFilters, ";") + ";" + strings.Join(audioFilters, ";")
	args := append(inputArgs, "-filter_complex", fullFilter, "-map", "[outv]")
	if hasAnyAudio {
		args = append(args, "-map", "[outa]", "-c:a", "aac", "-ar", "44100", "-ac", "2")
	}
	args = append(args, "-c:v", "libx264", "-pix_fmt", "yuv420p", "-y", outputPath)

	cmd := exec.Command("ffmpeg", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("xfade merge failed: %s", string(output))
	}
	return nil
}

func (c *Client) hasAudioStream(path string) bool {
	cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "a", "-show_entries", "stream=index", "-of", "csv=p=0", path)
	output, _ := cmd.CombinedOutput()
	return len(strings.TrimSpace(string(output))) > 0
}

func (c *Client) getVideoResolution(path string) (int, int) {
	cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries", "stream=width,height", "-of", "csv=p=0", path)
	output, _ := cmd.CombinedOutput()
	var w, h int
	fmt.Sscanf(strings.TrimSpace(string(output)), "%d,%d", &w, &h)
	if w == 0 {
		return 1280, 720
	}
	return w, h
}

func (c *Client) downloadVideo(url, destPath string) (string, error) {
	if !strings.HasPrefix(url, "http") {
		return url, nil
	}
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	out, _ := os.Create(destPath)
	defer out.Close()
	io.Copy(out, resp.Body)
	return destPath, nil
}

func (c *Client) mapTransitionType(t string) string {
	switch strings.ToLower(t) {
	case "fade":
		return "fadeblack"
	case "flash":
		return "fadewhite"
	case "crossfade":
		return "fade"
	default:
		return "fade"
	}
}

func (c *Client) copyFile(src, dst string) error {
	s, _ := os.Open(src)
	defer s.Close()
	d, _ := os.Create(dst)
	defer d.Close()
	io.Copy(d, s)
	return nil
}

func (c *Client) cleanupFiles(paths []string) {
	for _, p := range paths {
		os.Remove(p)
	}
}
