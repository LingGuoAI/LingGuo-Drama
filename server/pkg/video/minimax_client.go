package video

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type MinimaxClient struct {
	BaseURL    string
	APIKey     string
	Model      string
	HTTPClient *http.Client
}

type MinimaxSubjectReference struct {
	Type  string   `json:"type"`
	Image []string `json:"image"`
}

type MinimaxRequest struct {
	Prompt           string                    `json:"prompt"`
	FirstFrameImage  string                    `json:"first_frame_image,omitempty"`
	LastFrameImage   string                    `json:"last_frame_image,omitempty"`
	SubjectReference []MinimaxSubjectReference `json:"subject_reference,omitempty"`
	Model            string                    `json:"model"`
	Duration         int                       `json:"duration,omitempty"`
	Resolution       string                    `json:"resolution,omitempty"`
}

type MinimaxCreateResponse struct {
	TaskID   string `json:"task_id"`
	BaseResp struct {
		StatusCode int    `json:"status_code"`
		StatusMsg  string `json:"status_msg"`
	} `json:"base_resp"`
}

type MinimaxQueryResponse struct {
	TaskID      string `json:"task_id"`
	Status      string `json:"status"` // Processing, Success, Failed
	FileID      string `json:"file_id"`
	VideoWidth  int    `json:"video_width"`
	VideoHeight int    `json:"video_height"`
	BaseResp    struct {
		StatusCode int    `json:"status_code"`
		StatusMsg  string `json:"status_msg"`
	} `json:"base_resp"`
}

type MinimaxFileResponse struct {
	File struct {
		FileID      interface{} `json:"file_id"`
		Bytes       int         `json:"bytes"`
		CreatedAt   int64       `json:"created_at"`
		Filename    string      `json:"filename"`
		Purpose     string      `json:"purpose"`
		DownloadURL string      `json:"download_url"`
	} `json:"file"`
	BaseResp struct {
		StatusCode int    `json:"status_code"`
		StatusMsg  string `json:"status_msg"`
	} `json:"base_resp"`
}

func NewMinimaxClient(baseURL, apiKey, model string) *MinimaxClient {
	if baseURL == "" {
		baseURL = "https://api.minimax.chat/v1"
	}
	if model == "" {
		model = "MiniMax-Hailuo-02"
	}
	return &MinimaxClient{
		BaseURL:    baseURL,
		APIKey:     apiKey,
		Model:      model,
		HTTPClient: defaultHTTPClient(),
	}
}

func (c *MinimaxClient) GenerateVideo(prompt string, opts ...VideoOption) (*VideoResult, error) {
	options := &VideoOptions{
		Duration:   6,
		Resolution: "1080P",
	}

	for _, opt := range opts {
		opt(options)
	}

	model := c.Model
	if options.Model != "" {
		model = options.Model
	}

	reqBody := MinimaxRequest{
		Prompt:   prompt,
		Model:    model,
		Duration: options.Duration,
	}

	if options.Resolution != "" {
		reqBody.Resolution = options.Resolution
	}

	// 支持单图 (当首帧处理) 或 首尾帧
	if options.FirstFrameURL != "" {
		reqBody.FirstFrameImage = options.FirstFrameURL
	} else if options.ImageURL != "" {
		reqBody.FirstFrameImage = options.ImageURL
	}

	if options.LastFrameURL != "" {
		reqBody.LastFrameImage = options.LastFrameURL
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	endpoint := c.BaseURL + "/video_generation"
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	fmt.Printf("Minimax: Sending generation request to: %s\n", endpoint)
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

	var result MinimaxCreateResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	if result.BaseResp.StatusCode != 0 {
		return nil, fmt.Errorf("minimax error: %s", result.BaseResp.StatusMsg)
	}

	videoResult := &VideoResult{
		TaskID:    result.TaskID,
		Status:    "Processing",
		Completed: false,
	}

	return videoResult, nil
}

func (c *MinimaxClient) GetTaskStatus(taskID string) (*VideoResult, error) {
	endpoint := fmt.Sprintf("%s/query/video_generation?task_id=%s", c.BaseURL, taskID)
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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var queryResult MinimaxQueryResponse
	if err := json.Unmarshal(body, &queryResult); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	if queryResult.BaseResp.StatusCode != 0 {
		return nil, fmt.Errorf("minimax error: %s", queryResult.BaseResp.StatusMsg)
	}

	videoResult := &VideoResult{
		TaskID:    queryResult.TaskID,
		Status:    queryResult.Status,
		Width:     queryResult.VideoWidth,
		Height:    queryResult.VideoHeight,
		Completed: false,
	}

	if queryResult.Status == "Success" && queryResult.FileID != "" {
		downloadURL, err := c.getFileDownloadURL(queryResult.FileID)
		if err != nil {
			return nil, fmt.Errorf("failed to get download URL: %w", err)
		}
		videoResult.VideoURL = downloadURL
		videoResult.Completed = true
	} else if queryResult.Status == "Failed" {
		videoResult.Error = "Video generation failed"
		videoResult.Completed = true
	}

	return videoResult, nil
}

func (c *MinimaxClient) getFileDownloadURL(fileID string) (string, error) {
	endpoint := fmt.Sprintf("%s/files/retrieve?file_id=%s", c.BaseURL, fileID)
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var fileResult MinimaxFileResponse
	if err := json.Unmarshal(body, &fileResult); err != nil {
		return "", fmt.Errorf("parse response: %w", err)
	}

	if fileResult.BaseResp.StatusCode != 0 {
		return "", fmt.Errorf("minimax error: %s", fileResult.BaseResp.StatusMsg)
	}

	return fileResult.File.DownloadURL, nil
}
