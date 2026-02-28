package video

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// --- Runway ---
type RunwayClient struct {
	BaseURL    string
	APIKey     string
	Model      string
	HTTPClient *http.Client
}

type RunwayRequest struct {
	Model       string `json:"model"`
	PromptImage string `json:"prompt_image,omitempty"`
	PromptText  string `json:"prompt_text"`
	Duration    int    `json:"duration,omitempty"`
	AspectRatio string `json:"aspect_ratio,omitempty"`
	Seed        int64  `json:"seed,omitempty"`
}

type RunwayResponse struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	Output struct {
		URL string `json:"url"`
	} `json:"output"`
	Error string `json:"error,omitempty"`
}

func NewRunwayClient(baseURL, apiKey, model string) *RunwayClient {
	if baseURL == "" {
		baseURL = "https://api.runwayml.com"
	}
	if model == "" {
		model = "gen3a_turbo" // 占位
	}
	return &RunwayClient{
		BaseURL:    baseURL,
		APIKey:     apiKey,
		Model:      model,
		HTTPClient: defaultHTTPClient(),
	}
}

func (c *RunwayClient) GenerateVideo(prompt string, opts ...VideoOption) (*VideoResult, error) {
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

	reqBody := RunwayRequest{
		Model:       model,
		PromptImage: options.ImageURL,
		PromptText:  prompt,
		Duration:    options.Duration,
		AspectRatio: options.AspectRatio,
		Seed:        options.Seed,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	endpoint := c.BaseURL + "/v1/image_to_video" // 根据最新文档可能不同
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("X-Runway-Version", "2024-09-13")

	fmt.Printf("Runway: Sending generation request to: %s\n", endpoint)
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var result RunwayResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if result.Error != "" {
		return nil, fmt.Errorf("runway error: %s", result.Error)
	}

	videoResult := &VideoResult{
		TaskID:    result.ID,
		Status:    result.Status,
		Completed: result.Status == "SUCCEEDED",
	}
	if result.Output.URL != "" {
		videoResult.VideoURL = result.Output.URL
	}
	return videoResult, nil
}

func (c *RunwayClient) GetTaskStatus(taskID string) (*VideoResult, error) {
	endpoint := c.BaseURL + "/v1/tasks/" + taskID
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("X-Runway-Version", "2024-09-13")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result RunwayResponse
	json.Unmarshal(body, &result)

	videoResult := &VideoResult{
		TaskID:    result.ID,
		Status:    result.Status,
		Completed: result.Status == "SUCCEEDED",
		Error:     result.Error,
	}
	if result.Output.URL != "" {
		videoResult.VideoURL = result.Output.URL
	}
	return videoResult, nil
}

// --- Pika --- (实现类似, 保持结构一致)
type PikaClient struct {
	BaseURL    string
	APIKey     string
	Model      string
	HTTPClient *http.Client
}

func NewPikaClient(baseURL, apiKey, model string) *PikaClient {
	return &PikaClient{BaseURL: baseURL, APIKey: apiKey, Model: model, HTTPClient: defaultHTTPClient()}
}

// 示例实现，防止接口报错
func (c *PikaClient) GenerateVideo(prompt string, opts ...VideoOption) (*VideoResult, error) {
	return nil, fmt.Errorf("pika implementation omitted for brevity")
}
func (c *PikaClient) GetTaskStatus(taskID string) (*VideoResult, error) {
	return nil, fmt.Errorf("pika implementation omitted for brevity")
}
