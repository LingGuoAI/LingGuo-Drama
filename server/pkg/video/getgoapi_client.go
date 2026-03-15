package video

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"spiritFruit/pkg/upload"
	"strings"
)

// GetGoAPIClient GetGoAPI 视频生成客户端
type GetGoAPIClient struct {
	BaseURL       string
	APIKey        string
	Model         string
	Endpoint      string
	QueryEndpoint string
	HTTPClient    *http.Client
}

type GetGoAPIRequest struct {
	Model    string `json:"model"`
	Prompt   string `json:"prompt"`
	ImageURL string `json:"image_url,omitempty"`
	Duration int    `json:"duration,omitempty"`
	Size     string `json:"size,omitempty"`
}

// GetGoAPISoraRequest Sora 模型请求格式
type GetGoAPISoraRequest struct {
	Model          string `json:"model"`
	Prompt         string `json:"prompt"`
	Seconds        string `json:"seconds,omitempty"`
	Size           string `json:"size,omitempty"`
	InputReference string `json:"input_reference,omitempty"`
}

// GetGoAPIDoubaoRequest 豆包/火山模型请求格式
type GetGoAPIDoubaoRequest struct {
	Model   string `json:"model"`
	Content []struct {
		Type     string                 `json:"type"`
		Text     string                 `json:"text,omitempty"`
		ImageURL map[string]interface{} `json:"image_url,omitempty"`
		Role     string                 `json:"role,omitempty"`
	} `json:"content"`
}

type GetGoAPIResponse struct {
	ID     string          `json:"id"`
	TaskID string          `json:"task_id,omitempty"`
	Status string          `json:"status,omitempty"`
	Error  json.RawMessage `json:"error,omitempty"`
	Data   struct {
		ID       string `json:"id,omitempty"`
		Status   string `json:"status,omitempty"`
		VideoURL string `json:"video_url,omitempty"`
	} `json:"data,omitempty"`
}

type GetGoAPITaskResponse struct {
	ID       string          `json:"id,omitempty"`
	TaskID   string          `json:"task_id,omitempty"`
	Status   string          `json:"status,omitempty"`
	VideoURL string          `json:"video_url,omitempty"`
	Error    json.RawMessage `json:"error,omitempty"`
	Data     struct {
		ID       string `json:"id,omitempty"`
		Status   string `json:"status,omitempty"`
		VideoURL string `json:"video_url,omitempty"`
	} `json:"data,omitempty"`
	Content struct {
		VideoURL string `json:"video_url,omitempty"`
	} `json:"content,omitempty"`
}

// getErrorMessage 从 error 字段提取错误信息（支持字符串或对象）
func getErrorMessage(errorData json.RawMessage) string {
	if len(errorData) == 0 {
		return ""
	}

	// 尝试解析为字符串
	var errStr string
	if err := json.Unmarshal(errorData, &errStr); err == nil {
		return errStr
	}

	// 尝试解析为对象
	var errObj struct {
		Message string `json:"message"`
		Code    string `json:"code"`
	}
	if err := json.Unmarshal(errorData, &errObj); err == nil {
		if errObj.Message != "" {
			return errObj.Message
		}
	}

	// 返回原始 JSON 字符串
	return string(errorData)
}

func NewGetGoAPIClient(baseURL, apiKey, model, endpoint, queryEndpoint string) *GetGoAPIClient {
	if endpoint == "" {
		endpoint = "/videos"
	}
	if queryEndpoint == "" {
		queryEndpoint = "/videos/{taskId}/content"
	}
	return &GetGoAPIClient{
		BaseURL:       baseURL,
		APIKey:        apiKey,
		Model:         model,
		Endpoint:      endpoint,
		QueryEndpoint: queryEndpoint,
		HTTPClient:    defaultHTTPClient(),
	}
}

// GenerateVideo 发起生成请求，统一使用 VideoOption
func (c *GetGoAPIClient) GenerateVideo(prompt string, opts ...VideoOption) (*VideoResult, error) {
	options := &VideoOptions{
		Duration:    5,
		AspectRatio: "16:9",
	}

	for _, opt := range opts {
		opt(options)
	}

	model := c.Model
	if options.Model != "" {
		model = options.Model
	}

	var jsonData []byte
	var err error

	if strings.Contains(model, "doubao") || strings.Contains(model, "seedance") {
		// 豆包/火山格式
		reqBody := GetGoAPIDoubaoRequest{
			Model: model,
		}

		// 构建prompt文本（包含duration和ratio参数）
		promptText := prompt
		if options.AspectRatio != "" {
			promptText += fmt.Sprintf("  --ratio %s", options.AspectRatio)
		}
		if options.Duration > 0 {
			promptText += fmt.Sprintf("  --dur %d", options.Duration)
		}

		// 添加文本内容
		reqBody.Content = append(reqBody.Content, struct {
			Type     string                 `json:"type"`
			Text     string                 `json:"text,omitempty"`
			ImageURL map[string]interface{} `json:"image_url,omitempty"`
			Role     string                 `json:"role,omitempty"`
		}{Type: "text", Text: promptText})

		// 处理不同的图片模式
		if len(options.ReferenceImageURLs) > 0 {
			// 1. 组图模式
			for _, refURL := range options.ReferenceImageURLs {
				reqBody.Content = append(reqBody.Content, struct {
					Type     string                 `json:"type"`
					Text     string                 `json:"text,omitempty"`
					ImageURL map[string]interface{} `json:"image_url,omitempty"`
					Role     string                 `json:"role,omitempty"`
				}{
					Type: "image_url",
					ImageURL: map[string]interface{}{
						"url": refURL,
					},
					Role: "reference_image",
				})
			}
		} else if options.FirstFrameURL != "" && options.LastFrameURL != "" {
			// 2. 首尾帧模式
			reqBody.Content = append(reqBody.Content, struct {
				Type     string                 `json:"type"`
				Text     string                 `json:"text,omitempty"`
				ImageURL map[string]interface{} `json:"image_url,omitempty"`
				Role     string                 `json:"role,omitempty"`
			}{
				Type:     "image_url",
				ImageURL: map[string]interface{}{"url": options.FirstFrameURL},
				Role:     "first_frame",
			})
			reqBody.Content = append(reqBody.Content, struct {
				Type     string                 `json:"type"`
				Text     string                 `json:"text,omitempty"`
				ImageURL map[string]interface{} `json:"image_url,omitempty"`
				Role     string                 `json:"role,omitempty"`
			}{
				Type:     "image_url",
				ImageURL: map[string]interface{}{"url": options.LastFrameURL},
				Role:     "last_frame",
			})
		} else if options.ImageURL != "" {
			// 3. 单图模式（默认）
			reqBody.Content = append(reqBody.Content, struct {
				Type     string                 `json:"type"`
				Text     string                 `json:"text,omitempty"`
				ImageURL map[string]interface{} `json:"image_url,omitempty"`
				Role     string                 `json:"role,omitempty"`
			}{
				Type:     "image_url",
				ImageURL: map[string]interface{}{"url": options.ImageURL},
			})
		} else if options.FirstFrameURL != "" {
			// 4. 只有首帧
			reqBody.Content = append(reqBody.Content, struct {
				Type     string                 `json:"type"`
				Text     string                 `json:"text,omitempty"`
				ImageURL map[string]interface{} `json:"image_url,omitempty"`
				Role     string                 `json:"role,omitempty"`
			}{
				Type:     "image_url",
				ImageURL: map[string]interface{}{"url": options.FirstFrameURL},
				Role:     "first_frame",
			})
		}

		jsonData, err = json.Marshal(reqBody)

	} else if strings.Contains(model, "sora") {
		// Sora 格式
		seconds := fmt.Sprintf("%d", options.Duration)
		size := options.AspectRatio
		if size == "16:9" {
			size = "1280x720"
		} else if size == "9:16" {
			size = "720x1280"
		}

		// 提取图片（优先 ImageURL，其次 FirstFrame）
		refImg := options.ImageURL
		if refImg == "" {
			refImg = options.FirstFrameURL
		}

		reqBody := GetGoAPISoraRequest{
			Model:          model,
			Prompt:         prompt,
			Seconds:        seconds,
			Size:           size,
			InputReference: refImg,
		}
		jsonData, err = json.Marshal(reqBody)

	} else {
		// 默认格式
		refImg := options.ImageURL
		if refImg == "" {
			refImg = options.FirstFrameURL
		}

		reqBody := GetGoAPIRequest{
			Model:    model,
			Prompt:   prompt,
			ImageURL: refImg,
			Duration: options.Duration,
			Size:     options.AspectRatio,
		}
		jsonData, err = json.Marshal(reqBody)
	}

	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	endpoint := c.BaseURL + c.Endpoint
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	fmt.Printf("[GetGoAPI] Response body: %s\n", string(body))

	var result GetGoAPIResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parse response: %w, body: %s", err, string(body))
	}

	// 优先使用 id 字段，其次使用 task_id，再次使用 data 里的 id
	taskID := result.ID
	if taskID == "" {
		taskID = result.TaskID
	}
	if result.Data.ID != "" {
		taskID = result.Data.ID
	}

	status := result.Status
	if status == "" && result.Data.Status != "" {
		status = result.Data.Status
	}

	fmt.Printf("[GetGoAPI] Parsed result - TaskID: %s, Status: %s\n", taskID, status)

	if errMsg := getErrorMessage(result.Error); errMsg != "" {
		return nil, fmt.Errorf("GetGoAPI error: %s", errMsg)
	}

	videoResult := &VideoResult{
		TaskID:    taskID,
		Status:    status,
		Completed: status == "completed" || status == "succeeded",
		Duration:  options.Duration,
	}

	return videoResult, nil
}

func (c *GetGoAPIClient) GetTaskStatus(taskID string) (*VideoResult, error) {
	// 1. 查询任务基础状态
	queryPath := c.QueryEndpoint
	// 如果配置的路径包含 content，说明是专用 content 接口，先回退到基础路径查询状态
	if strings.Contains(queryPath, "/content") {
		queryPath = "/videos/{taskId}"
	}

	if strings.Contains(queryPath, "{taskId}") {
		queryPath = strings.ReplaceAll(queryPath, "{taskId}", taskID)
	} else {
		queryPath = strings.TrimRight(queryPath, "/") + "/" + taskID
	}

	endpoint := strings.TrimRight(c.BaseURL, "/") + "/" + strings.TrimLeft(queryPath, "/")

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	// 处理 Sora 任务进行中的 400 状态 (GetGoAPI 特色)
	if resp.StatusCode == http.StatusBadRequest {
		bodyStr := string(body)
		if strings.Contains(bodyStr, "IN_PROGRESS") || strings.Contains(bodyStr, "not completed") {
			return &VideoResult{
				TaskID:    taskID,
				Status:    "processing",
				Completed: false,
			}, nil
		}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var result GetGoAPITaskResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parse response: %w, body: %s", err, string(body))
	}

	// 提取状态和 ID
	status := result.Status
	if status == "" {
		status = result.Data.Status
	}

	videoResult := &VideoResult{
		TaskID:    taskID,
		Status:    status,
		Completed: status == "completed" || status == "succeeded",
	}

	// 提取 URL
	videoURL := result.VideoURL
	if videoURL == "" {
		videoURL = result.Data.VideoURL
	}
	if videoURL == "" {
		videoURL = result.Content.VideoURL
	}
	videoResult.VideoURL = videoURL

	// 🔴 核心逻辑：如果任务已完成，但是基础接口没给 URL，就去调用专用的 content 接口提取或保存视频
	if videoResult.Completed && videoResult.VideoURL == "" {
		contentURL, fetchErr := c.fetchVideoContent(taskID)
		if fetchErr != nil {
			// 如果 content 接口报进度错误，说明还没真正完成
			if strings.Contains(fetchErr.Error(), "IN_PROGRESS") {
				videoResult.Completed = false
				videoResult.Status = "processing"
			} else {
				videoResult.Error = fmt.Sprintf("fetch content failed: %v", fetchErr)
			}
		} else {
			videoResult.VideoURL = contentURL
		}
	}

	return videoResult, nil
}

// fetchVideoContent 处理二进制流或 JSON URL
func (c *GetGoAPIClient) fetchVideoContent(taskID string) (string, error) {
	// 构造 content 接口地址
	endpoint := fmt.Sprintf("%s/videos/%s/content", strings.TrimRight(c.BaseURL, "/"), taskID)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 处理进行中报错
	if resp.StatusCode == http.StatusBadRequest && strings.Contains(string(body), "IN_PROGRESS") {
		return "", fmt.Errorf("IN_PROGRESS")
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("content api error, status: %d, body: %s", resp.StatusCode, string(body))
	}

	// 🔴 1. 检测是否直接返回了视频二进制流 (mp4)
	contentType := resp.Header.Get("Content-Type")
	// 检查 Header 或者检查文件头魔数 (ftypmp42 等)
	isMP4 := strings.Contains(contentType, "video/") || (len(body) > 8 && string(body[4:8]) == "ftyp")

	if isMP4 {
		// 💡 直接落盘：将二进制流保存到本地服务器
		savedPath, err := upload.SaveFileDirByte(body, "videos", ".mp4")
		if err != nil {
			return "", fmt.Errorf("failed to save binary video: %w", err)
		}
		fmt.Printf("[GetGoAPI] 二进制视频已保存至本地: %s\n", savedPath)
		return savedPath, nil
	}

	// 2. 尝试解析为 JSON (部分厂商完成时会返回包含 url 的 JSON)
	var contentRes struct {
		URL      string `json:"url"`
		VideoURL string `json:"video_url"`
		Data     struct {
			URL string `json:"url"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &contentRes); err == nil {
		if contentRes.URL != "" {
			return contentRes.URL, nil
		}
		if contentRes.VideoURL != "" {
			return contentRes.VideoURL, nil
		}
		if contentRes.Data.URL != "" {
			return contentRes.Data.URL, nil
		}
	}

	// 3. 尝试纯文本 URL
	bodyStr := strings.TrimSpace(string(body))
	if strings.HasPrefix(bodyStr, "http://") || strings.HasPrefix(bodyStr, "https://") {
		return bodyStr, nil
	}

	return "", fmt.Errorf("failed to parse video source from content api")
}
