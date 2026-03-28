package api

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/icosmos-space/iadmin/server/model/common/response"
	commonResp "github.com/icosmos-space/iadmin/server/model/common/response"
	"github.com/icosmos-space/iadmin/server/plugin/aiassistant/model"
	aiReq "github.com/icosmos-space/iadmin/server/plugin/aiassistant/model/request"
)

type assistant struct{}

func (a *assistant) GetConfig(c *gin.Context) {
	cfg, err := serviceAssistant.GetConfig()
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(cfg, c)
}

func (a *assistant) UpdateConfig(c *gin.Context) {
	var in aiReq.UpdateConfigReq
	if err := c.ShouldBindJSON(&in); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := serviceAssistant.UpdateConfig(in); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("updated", c)
}

func (a *assistant) CreatePrompt(c *gin.Context) {
	var prompt model.AIAssistantPrompt
	if err := c.ShouldBindJSON(&prompt); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := serviceAssistant.CreatePrompt(&prompt); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("created", c)
}

func (a *assistant) UpdatePrompt(c *gin.Context) {
	var prompt model.AIAssistantPrompt
	if err := c.ShouldBindJSON(&prompt); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := serviceAssistant.UpdatePrompt(prompt); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("updated", c)
}

func (a *assistant) DeletePrompt(c *gin.Context) {
	id := c.Query("ID")
	if id == "" {
		response.FailWithMessage("ID is required", c)
		return
	}
	if err := serviceAssistant.DeletePrompt(id); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("deleted", c)
}

func (a *assistant) DeletePromptByIDs(c *gin.Context) {
	ids := c.QueryArray("IDs[]")
	if err := serviceAssistant.DeletePromptByIDs(ids); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("deleted", c)
}

func (a *assistant) FindPrompt(c *gin.Context) {
	id := c.Query("ID")
	if id == "" {
		response.FailWithMessage("ID is required", c)
		return
	}
	prompt, err := serviceAssistant.GetPrompt(id)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(prompt, c)
}

func (a *assistant) GetPromptList(c *gin.Context) {
	var pageInfo aiReq.PromptSearch
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := serviceAssistant.GetPromptList(pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(commonResp.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "ok", c)
}

func (a *assistant) GetEnabledPrompts(c *gin.Context) {
	list, err := serviceAssistant.GetEnabledPromptList(20)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(list, c)
}

func (a *assistant) Chat(c *gin.Context) {
	var in aiReq.ChatReq
	if err := c.ShouldBindJSON(&in); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if in.Stream {
		//设置SSE必要Header
		c.Writer.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("X-Accel-Buffering", "no")
		c.Writer.WriteHeader(http.StatusOK)

		flusher, ok := c.Writer.(http.Flusher)
		if !ok {
			response.FailWithMessage("streaming unsupported", c)
			return
		}
		// 立刻发一个 meta，便于前端确认流已建立
		c.SSEvent("meta", gin.H{"ok": true, "ts": time.Now().Unix()})
		flusher.Flush()
	}

	resp, err := serviceAssistant.RequestChat(c.Request.Context(), in)
	if err != nil {
		c.SSEvent("error", gin.H{"message": err.Error()})
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		response.FailWithMessage(extractUpstreamError(string(body), resp.Status), c)
		return
	}

	// 如果请求是流式，直接使用 SSE 代理
	if in.Stream {
		a.proxyStream(c, resp.Body)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 尝试解析 AI 响应为 JSON
	var aiResponse map[string]interface{}
	if err := json.Unmarshal(body, &aiResponse); err != nil {
		// 解析失败，返回原始内容
		response.OkWithData(gin.H{"content": string(body)}, c)
		return
	}

	// 包装成标准响应格式
	response.OkWithData(aiResponse, c)
}

func (a *assistant) proxyStream(c *gin.Context, body io.Reader) {
	// 设置 SSE 必要的响应头
	c.Writer.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		response.FailWithMessage("streaming unsupported", c)
		return
	}

	// 发送初始 meta 事件，让前端确认连接已建立
	c.SSEvent("meta", gin.H{"ok": true, "ts": time.Now().Unix()})
	flusher.Flush()

	// 读取并转发 AI 的 SSE 流
	buffer := make([]byte, 4096)
	for {
		select {
		case <-c.Request.Context().Done():
			return
		default:
		}

		n, err := body.Read(buffer)
		if n > 0 {
			// 直接转发 AI 的原始 SSE 数据
			if _, wErr := c.Writer.Write(buffer[:n]); wErr != nil {
				return
			}
			flusher.Flush()
		}
		if err == nil {
			continue
		}
		if errors.Is(err, io.EOF) {
			// 发送结束事件
			c.SSEvent("done", gin.H{"ok": true})
			flusher.Flush()
			return
		}
		// 其他错误
		c.SSEvent("error", gin.H{"msg": err.Error()})
		flusher.Flush()
		return
	}
}

func extractUpstreamError(raw, statusText string) string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return statusText
	}

	var payload map[string]any
	if err := json.Unmarshal([]byte(trimmed), &payload); err == nil {
		if msg := firstString(
			payload["message"],
			payload["msg"],
			payload["error_description"],
			payload["error"],
			payload["detail"],
		); msg != "" {
			return msg
		}
		if errObj, ok := payload["error"].(map[string]any); ok {
			if msg := firstString(errObj["message"], errObj["msg"], errObj["detail"]); msg != "" {
				return msg
			}
		}
	}

	if len(trimmed) > 500 {
		return trimmed[:500]
	}
	return trimmed
}

func firstString(values ...any) string {
	for _, v := range values {
		s := strings.TrimSpace(toString(v))
		if s != "" {
			return s
		}
	}
	return ""
}

func toString(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case []byte:
		return string(t)
	default:
		return ""
	}
}
