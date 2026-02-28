package video

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"path/filepath"
	"strings"
)

type OpenAISoraClient struct {
	BaseURL    string
	APIKey     string
	Model      string
	HTTPClient *http.Client
}

type OpenAISoraResponse struct {
	ID          string `json:"id"`
	Object      string `json:"object"`
	Model       string `json:"model"`
	Status      string `json:"status"`
	Progress    int    `json:"progress"`
	CreatedAt   int64  `json:"created_at"`
	CompletedAt int64  `json:"completed_at"`
	Size        string `json:"size"`
	Seconds     string `json:"seconds"`
	Quality     string `json:"quality"`
	VideoURL    string `json:"video_url"`
	Video       struct {
		URL string `json:"url"`
	} `json:"video"`
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error"`
}

func NewOpenAISoraClient(baseURL, apiKey, model string) *OpenAISoraClient {
	if baseURL == "" {
		baseURL = "https://api.openai.com"
	}
	if model == "" {
		model = "sora-1.0-turbo"
	}
	return &OpenAISoraClient{
		BaseURL:    baseURL,
		APIKey:     apiKey,
		Model:      model,
		HTTPClient: defaultHTTPClient(),
	}
}

func (c *OpenAISoraClient) GenerateVideo(prompt string, opts ...VideoOption) (*VideoResult, error) {
	options := &VideoOptions{
		Duration: 4,
	}

	for _, opt := range opts {
		opt(options)
	}

	model := c.Model
	if options.Model != "" {
		model = options.Model
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	_ = writer.WriteField("model", model)
	_ = writer.WriteField("prompt", prompt)

	if options.Duration > 0 {
		_ = writer.WriteField("seconds", fmt.Sprintf("%d", options.Duration))
	}

	if options.Resolution != "" {
		_ = writer.WriteField("size", options.Resolution)
	}

	// Sora API 要求 'input_reference' 为文件上传 (binary)
	if options.ImageURL != "" {
		if err := c.addImageToMultipart(writer, "input_reference", options.ImageURL); err != nil {
			return nil, err
		}
	}

	writer.Close()

	endpoint := c.BaseURL + "/v1/videos"
	req, err := http.NewRequest("POST", endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	fmt.Printf("Sora: Sending generation request to: %s\n", endpoint)
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var result OpenAISoraResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	if result.Error.Message != "" {
		return nil, fmt.Errorf("openai error: %s", result.Error.Message)
	}

	videoResult := &VideoResult{
		TaskID:    result.ID,
		Status:    result.Status,
		Completed: result.Status == "completed",
	}

	if result.VideoURL != "" {
		videoResult.VideoURL = result.VideoURL
	} else if result.Video.URL != "" {
		videoResult.VideoURL = result.Video.URL
	}

	return videoResult, nil
}

func (c *OpenAISoraClient) GetTaskStatus(taskID string) (*VideoResult, error) {
	endpoint := c.BaseURL + "/v1/videos/" + taskID
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

	var result OpenAISoraResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	videoResult := &VideoResult{
		TaskID:    result.ID,
		Status:    result.Status,
		Completed: result.Status == "completed",
	}

	if result.Error.Message != "" {
		videoResult.Error = result.Error.Message
	}

	if result.VideoURL != "" {
		videoResult.VideoURL = result.VideoURL
	} else if result.Video.URL != "" {
		videoResult.VideoURL = result.Video.URL
	}

	return videoResult, nil
}

// 辅助方法：将 URL 或 Base64 图片添加到 multipart
func (c *OpenAISoraClient) addImageToMultipart(writer *multipart.Writer, fieldName, imageURL string) error {
	var imageData []byte
	var mimeType string
	var filename string = "reference_image.png"

	if strings.HasPrefix(imageURL, "data:") {
		parts := strings.Split(imageURL, ",")
		if len(parts) != 2 {
			return fmt.Errorf("invalid data URI format")
		}

		header := parts[0]
		if strings.Contains(header, "image/jpeg") || strings.Contains(header, "image/jpg") {
			mimeType = "image/jpeg"
			filename = "reference.jpg"
		} else if strings.Contains(header, "image/png") {
			mimeType = "image/png"
			filename = "reference.png"
		} else if strings.Contains(header, "image/webp") {
			mimeType = "image/webp"
			filename = "reference.webp"
		} else {
			mimeType = "image/png"
		}

		decoded, err := base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			return fmt.Errorf("failed to decode base64 image: %w", err)
		}
		imageData = decoded

	} else {
		resp, err := http.Get(imageURL)
		if err != nil {
			return fmt.Errorf("failed to download reference image: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("failed to download reference image, status: %d", resp.StatusCode)
		}

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read downloaded image: %w", err)
		}
		imageData = data

		mimeType = resp.Header.Get("Content-Type")

		if mimeType == "" || mimeType == "application/octet-stream" {
			ext := filepath.Ext(imageURL)
			switch strings.ToLower(ext) {
			case ".jpg", ".jpeg":
				mimeType = "image/jpeg"
			case ".png":
				mimeType = "image/png"
			case ".webp":
				mimeType = "image/webp"
			default:
				mimeType = "image/png"
			}
		}

		base := filepath.Base(imageURL)
		if base != "" && base != "." {
			if idx := strings.Index(base, "?"); idx != -1 {
				base = base[:idx]
			}
			filename = base
		}
	}

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fieldName, filename))
	h.Set("Content-Type", mimeType)

	part, err := writer.CreatePart(h)
	if err != nil {
		return fmt.Errorf("create part: %w", err)
	}
	if _, err := part.Write(imageData); err != nil {
		return fmt.Errorf("write image data: %w", err)
	}

	return nil
}
