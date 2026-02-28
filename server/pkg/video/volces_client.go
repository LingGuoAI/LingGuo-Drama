package video

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type VolcesArkClient struct {
	BaseURL       string
	APIKey        string
	Model         string
	Endpoint      string
	QueryEndpoint string
	HTTPClient    *http.Client
}

type VolcesArkContent struct {
	Type     string                 `json:"type"`
	Text     string                 `json:"text,omitempty"`
	ImageURL map[string]interface{} `json:"image_url,omitempty"`
	Role     string                 `json:"role,omitempty"`
}

type VolcesArkRequest struct {
	Model         string             `json:"model"`
	Content       []VolcesArkContent `json:"content"`
	GenerateAudio bool               `json:"generate_audio,omitempty"`
}

type VolcesArkResponse struct {
	ID      string `json:"id"`
	Model   string `json:"model"`
	Status  string `json:"status"`
	Content struct {
		VideoURL string `json:"video_url"`
	} `json:"content"`
	Duration int         `json:"duration"`
	Error    interface{} `json:"error,omitempty"`
}

func NewVolcesArkClient(baseURL, apiKey, model, endpoint, queryEndpoint string) *VolcesArkClient {
	if baseURL == "" {
		baseURL = "https://open.volcengineapi.com"
	}
	if endpoint == "" {
		endpoint = "/api/v3/contents/generations/tasks"
	}
	if queryEndpoint == "" {
		queryEndpoint = endpoint
	}
	return &VolcesArkClient{
		BaseURL:       baseURL,
		APIKey:        apiKey,
		Model:         model,
		Endpoint:      endpoint,
		QueryEndpoint: queryEndpoint,
		HTTPClient:    defaultHTTPClient(),
	}
}

func (c *VolcesArkClient) GenerateVideo(prompt string, opts ...VideoOption) (*VideoResult, error) {
	options := &VideoOptions{
		Duration:    5,
		AspectRatio: "adaptive",
	}

	for _, opt := range opts {
		opt(options)
	}

	model := c.Model
	if options.Model != "" {
		model = options.Model
	}

	promptText := prompt
	if options.AspectRatio != "" && options.AspectRatio != "adaptive" {
		promptText += fmt.Sprintf("  --ratio %s", options.AspectRatio)
	}
	if options.Duration > 0 {
		promptText += fmt.Sprintf("  --dur %d", options.Duration)
	}

	content := []VolcesArkContent{
		{
			Type: "text",
			Text: promptText,
		},
	}

	// 处理图片：多图优先 -> 其次首尾帧 -> 其次首帧 -> 最后单图
	if len(options.ReferenceImageURLs) > 0 {
		for _, refURL := range options.ReferenceImageURLs {
			content = append(content, VolcesArkContent{
				Type:     "image_url",
				ImageURL: map[string]interface{}{"url": refURL},
				Role:     "reference_image",
			})
		}
	} else if options.FirstFrameURL != "" && options.LastFrameURL != "" {
		content = append(content, VolcesArkContent{
			Type:     "image_url",
			ImageURL: map[string]interface{}{"url": options.FirstFrameURL},
			Role:     "first_frame",
		})
		content = append(content, VolcesArkContent{
			Type:     "image_url",
			ImageURL: map[string]interface{}{"url": options.LastFrameURL},
			Role:     "last_frame",
		})
	} else if options.FirstFrameURL != "" {
		content = append(content, VolcesArkContent{
			Type:     "image_url",
			ImageURL: map[string]interface{}{"url": options.FirstFrameURL},
			Role:     "first_frame",
		})
	} else if options.ImageURL != "" {
		content = append(content, VolcesArkContent{
			Type:     "image_url",
			ImageURL: map[string]interface{}{"url": options.ImageURL},
		})
	}

	generateAudio := strings.Contains(strings.ToLower(model), "seedance-1-5-pro")

	reqBody := VolcesArkRequest{
		Model:         model,
		Content:       content,
		GenerateAudio: generateAudio,
	}

	jsonData, err := json.Marshal(reqBody)
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

	fmt.Printf("Volces: Sending generation request to: %s\n", endpoint)
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

	var result VolcesArkResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	if result.Error != nil {
		return nil, fmt.Errorf("volces error: %v", result.Error)
	}

	videoResult := &VideoResult{
		TaskID:    result.ID,
		Status:    result.Status,
		Completed: result.Status == "completed" || result.Status == "succeeded",
		Duration:  result.Duration,
	}

	if result.Content.VideoURL != "" {
		videoResult.VideoURL = result.Content.VideoURL
		videoResult.Completed = true
	}

	return videoResult, nil
}

func (c *VolcesArkClient) GetTaskStatus(taskID string) (*VideoResult, error) {
	queryPath := c.QueryEndpoint
	if strings.Contains(queryPath, "{taskId}") {
		queryPath = strings.ReplaceAll(queryPath, "{taskId}", taskID)
	} else if strings.Contains(queryPath, "{task_id}") {
		queryPath = strings.ReplaceAll(queryPath, "{task_id}", taskID)
	} else {
		queryPath = queryPath + "/" + taskID
	}

	endpoint := c.BaseURL + queryPath
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

	var result VolcesArkResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	videoResult := &VideoResult{
		TaskID:    result.ID,
		Status:    result.Status,
		Completed: result.Status == "completed" || result.Status == "succeeded",
		Duration:  result.Duration,
	}

	if result.Error != nil {
		videoResult.Error = fmt.Sprintf("%v", result.Error)
	}

	if result.Content.VideoURL != "" {
		videoResult.VideoURL = result.Content.VideoURL
		videoResult.Completed = true
	}

	return videoResult, nil
}
