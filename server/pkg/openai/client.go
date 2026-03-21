package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap"
	"spiritFruit/pkg/logger"
)

var (
	globalHTTPClient *http.Client
	once             sync.Once
)

func initGlobalHTTPClient() {
	apiTransport := &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:        100,
		TLSHandshakeTimeout: 10 * time.Second,
	}
	globalHTTPClient = &http.Client{
		Timeout:   300 * time.Second,
		Transport: apiTransport,
	}
}

// NewProvider 工厂方法：根据配置返回对应的实现
func NewProvider(cfg Config) Provider {
	once.Do(func() { initGlobalHTTPClient() })
	switch cfg.Provider {
	case "getgoapi":
		return &GetGoAPIClient{Config: cfg, client: globalHTTPClient}
	case "gemini":
		return &GeminiClient{Config: cfg, client: globalHTTPClient}
	case "doubao", "volces":
		return &DoubaoClient{Config: cfg, client: globalHTTPClient}
	case "vertex", "gcp":
		return &VertexClient{Config: cfg, client: globalHTTPClient}
	case "openai":
		fallthrough
	default:
		return &OpenAIClient{Config: cfg, client: globalHTTPClient}
	}
}

// doRequest 通用泛型请求 (保留你的风格)
func doRequest[T any](client *http.Client, method, url string, headers map[string]string, payload interface{}) (T, error) {
	var result T
	var reqBody io.Reader

	if payload != nil {
		jsonBytes, _ := json.Marshal(payload)
		reqBody = bytes.NewBuffer(jsonBytes)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return result, err
	}

	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		logger.Error("AI API Error", zap.String("url", url), zap.String("body", string(bodyBytes)))
		return result, fmt.Errorf("API error: %d", resp.StatusCode)
	}

	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return result, fmt.Errorf("unmarshal failed: %v", err)
	}

	return result, nil
}
