package service

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/goccy/go-json"
	"github.com/icosmos-space/iadmin/server/global"
	"github.com/icosmos-space/iadmin/server/plugin/aiassistant/model"
	aiReq "github.com/icosmos-space/iadmin/server/plugin/aiassistant/model/request"
	"gorm.io/gorm"
)

type assistant struct{}

type openAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIChatRequest struct {
	Model       string          `json:"model"`
	Messages    []openAIMessage `json:"messages"`
	Stream      bool            `json:"stream"`
	Temperature float64         `json:"temperature,omitempty"`
}

func (s *assistant) InitDefaults() error {
	if err := s.ensureDefaultConfig(); err != nil {
		return err
	}
	return s.ensureDefaultPrompts()
}

func (s *assistant) GetConfig() (model.AIAssistantConfig, error) {
	cfg, err := s.getOrCreateDefaultConfig()
	if err != nil {
		return model.AIAssistantConfig{}, err
	}
	return cfg, nil
}

func (s *assistant) UpdateConfig(in aiReq.UpdateConfigReq) error {
	cfg, err := s.getOrCreateDefaultConfig()
	if err != nil {
		return err
	}

	cfg.BaseURL = strings.TrimSpace(in.BaseURL)
	cfg.Token = strings.TrimSpace(in.Token)
	cfg.Model = strings.TrimSpace(in.Model)
	cfg.ChatPath = strings.TrimSpace(in.ChatPath)
	cfg.Temperature = in.Temperature
	cfg.TimeoutSec = in.TimeoutSec
	cfg.Enabled = in.Enabled

	if cfg.Model == "" {
		cfg.Model = "gpt-4o-mini"
	}
	if cfg.ChatPath == "" {
		cfg.ChatPath = "/v1/chat/completions"
	}
	if cfg.TimeoutSec < 10 || cfg.TimeoutSec > 300 {
		cfg.TimeoutSec = 60
	}
	if cfg.Temperature < 0 {
		cfg.Temperature = 0
	}
	if cfg.Temperature > 2 {
		cfg.Temperature = 2
	}

	return global.IADMIN_DB.Save(&cfg).Error
}

func (s *assistant) CreatePrompt(prompt *model.AIAssistantPrompt) error {
	prompt.Title = strings.TrimSpace(prompt.Title)
	prompt.Content = strings.TrimSpace(prompt.Content)
	if prompt.Title == "" {
		return errors.New("title cannot be empty")
	}
	if prompt.Content == "" {
		return errors.New("content cannot be empty")
	}
	return global.IADMIN_DB.Create(prompt).Error
}

func (s *assistant) UpdatePrompt(prompt model.AIAssistantPrompt) error {
	if prompt.ID == 0 {
		return errors.New("id is required")
	}
	prompt.Title = strings.TrimSpace(prompt.Title)
	prompt.Content = strings.TrimSpace(prompt.Content)
	if prompt.Title == "" {
		return errors.New("title cannot be empty")
	}
	if prompt.Content == "" {
		return errors.New("content cannot be empty")
	}
	return global.IADMIN_DB.Model(&model.AIAssistantPrompt{}).Where("id = ?", prompt.ID).Updates(map[string]any{
		"title":   prompt.Title,
		"content": prompt.Content,
		"sort":    prompt.Sort,
		"enabled": prompt.Enabled,
	}).Error
}

func (s *assistant) DeletePrompt(id string) error {
	return global.IADMIN_DB.Delete(&model.AIAssistantPrompt{}, "id = ?", id).Error
}

func (s *assistant) DeletePromptByIDs(ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	return global.IADMIN_DB.Delete(&[]model.AIAssistantPrompt{}, "id in ?", ids).Error
}

func (s *assistant) GetPrompt(id string) (model.AIAssistantPrompt, error) {
	var prompt model.AIAssistantPrompt
	err := global.IADMIN_DB.Where("id = ?", id).First(&prompt).Error
	return prompt, err
}

func (s *assistant) GetPromptList(info aiReq.PromptSearch) ([]model.AIAssistantPrompt, int64, error) {
	if info.Page <= 0 {
		info.Page = 1
	}
	if info.PageSize <= 0 {
		info.PageSize = 10
	}
	if info.PageSize > 100 {
		info.PageSize = 100
	}

	db := global.IADMIN_DB.Model(&model.AIAssistantPrompt{})
	if kw := strings.TrimSpace(info.Keyword); kw != "" {
		like := "%" + kw + "%"
		db = db.Where("title LIKE ? OR content LIKE ?", like, like)
	}
	if info.Enabled != nil {
		db = db.Where("enabled = ?", *info.Enabled)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var list []model.AIAssistantPrompt
	err := db.Order("sort asc, id asc").Offset((info.Page - 1) * info.PageSize).Limit(info.PageSize).Find(&list).Error
	return list, total, err
}

func (s *assistant) GetEnabledPromptList(limit int) ([]model.AIAssistantPrompt, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	var list []model.AIAssistantPrompt
	err := global.IADMIN_DB.Where("enabled = ?", true).Order("sort asc, id asc").Limit(limit).Find(&list).Error
	return list, err
}

func (s *assistant) RequestChat(ctx context.Context, in aiReq.ChatReq) (*http.Response, error) {
	cfg, err := s.getOrCreateDefaultConfig()
	if err != nil {
		return nil, err
	}
	if !cfg.Enabled {
		return nil, errors.New("AI assistant is disabled")
	}
	if strings.TrimSpace(cfg.BaseURL) == "" {
		return nil, errors.New("please configure AI base URL")
	}
	if strings.TrimSpace(cfg.Token) == "" {
		return nil, errors.New("please configure AI token")
	}

	messages := buildMessages(in)
	if len(messages) == 0 {
		return nil, errors.New("message cannot be empty")
	}

	modelName := strings.TrimSpace(in.Model)
	if modelName == "" {
		modelName = strings.TrimSpace(cfg.Model)
	}
	if modelName == "" {
		modelName = "gpt-4o-mini"
	}

	chatReq := openAIChatRequest{
		Model:    modelName,
		Messages: messages,
		Stream:   in.Stream,
	}
	if cfg.Temperature > 0 {
		chatReq.Temperature = cfg.Temperature
	}
	payload, err := json.Marshal(chatReq)
	if err != nil {
		return nil, err
	}

	target := buildTargetURL(cfg.BaseURL, cfg.ChatPath)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, target, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+cfg.Token)
	if in.Stream {
		httpReq.Header.Set("Accept", "text/event-stream, application/json")
	} else {
		httpReq.Header.Set("Accept", "application/json")
	}

	timeout := cfg.TimeoutSec
	if timeout < 10 || timeout > 300 {
		timeout = 60
	}
	client := &http.Client{Timeout: time.Duration(timeout) * time.Second}
	return client.Do(httpReq)
}

func buildMessages(in aiReq.ChatReq) []openAIMessage {
	list := make([]openAIMessage, 0, len(in.Messages)+1)
	for _, msg := range in.Messages {
		content := strings.TrimSpace(msg.Content)
		if content == "" {
			continue
		}
		role := normalizeRole(msg.Role)
		list = append(list, openAIMessage{Role: role, Content: content})
	}
	if len(list) > 0 {
		return list
	}

	question := strings.TrimSpace(in.Message)
	if question == "" {
		question = strings.TrimSpace(in.Question)
	}
	if question == "" {
		question = strings.TrimSpace(in.Prompt)
	}
	if question == "" {
		return nil
	}
	return []openAIMessage{{Role: "user", Content: question}}
}

func normalizeRole(role string) string {
	r := strings.ToLower(strings.TrimSpace(role))
	switch r {
	case "system", "assistant", "tool", "user":
		return r
	default:
		return "user"
	}
}

func buildTargetURL(baseURL, chatPath string) string {
	path := strings.TrimSpace(chatPath)
	if path == "" {
		path = "/v1/chat/completions"
	}
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path
	}

	base := strings.TrimSpace(baseURL)
	base = strings.TrimRight(base, "/")
	if base == "" {
		return path
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return base + path
}

func (s *assistant) ensureDefaultConfig() error {
	_, err := s.getOrCreateDefaultConfig()
	return err
}

func (s *assistant) getOrCreateDefaultConfig() (model.AIAssistantConfig, error) {
	var cfg model.AIAssistantConfig
	err := global.IADMIN_DB.Where("name = ?", "default").First(&cfg).Error
	if err == nil {
		return cfg, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return model.AIAssistantConfig{}, err
	}
	cfg = model.AIAssistantConfig{
		Name:        "default",
		BaseURL:     "https://api.openai.com",
		Token:       "",
		Model:       "gpt-4o-mini",
		ChatPath:    "/v1/chat/completions",
		Temperature: 0.7,
		TimeoutSec:  60,
		Enabled:     false,
	}
	if err = global.IADMIN_DB.Create(&cfg).Error; err != nil {
		return model.AIAssistantConfig{}, err
	}
	return cfg, nil
}

func (s *assistant) ensureDefaultPrompts() error {
	defaults := []model.AIAssistantPrompt{
		{Title: "总结工作重点", Content: "帮我总结今天的工作重点", Sort: 1, Enabled: true},
		{Title: "排查前端报错", Content: "帮我排查一个前端报错", Sort: 2, Enabled: true},
		{Title: "SQL 优化建议", Content: "给我一个 SQL 优化建议", Sort: 3, Enabled: true},
	}
	for i := range defaults {
		row := defaults[i]
		var found model.AIAssistantPrompt
		err := global.IADMIN_DB.Where("title = ?", row.Title).First(&found).Error
		if err == nil {
			continue
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err = global.IADMIN_DB.Create(&row).Error; err != nil {
				return err
			}
			continue
		}
		return err
	}
	return nil
}
