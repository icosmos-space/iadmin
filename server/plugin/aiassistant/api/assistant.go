package api

import (
	"bufio"
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

	var flusher http.Flusher
	if in.Stream {
		var ok bool
		flusher, ok = a.initSSE(c)
		if !ok {
			return
		}
	}

	resp, err := serviceAssistant.RequestChat(c.Request.Context(), in)
	if err != nil {
		if in.Stream {
			a.emitSSE(c, flusher, "error", gin.H{"message": err.Error()})
			a.emitSSE(c, flusher, "done", gin.H{"ok": true})
			return
		}
		response.FailWithMessage(err.Error(), c)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		msg := extractUpstreamError(string(body), resp.Status)
		if in.Stream {
			a.emitSSE(c, flusher, "error", gin.H{"message": msg})
			a.emitSSE(c, flusher, "done", gin.H{"ok": true})
			return
		}
		response.FailWithMessage(msg, c)
		return
	}

	if in.Stream {
		a.proxyStream(c, resp.Body, flusher)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	var aiResponse map[string]any
	if err := json.Unmarshal(body, &aiResponse); err != nil {
		response.OkWithData(gin.H{"content": string(body)}, c)
		return
	}

	response.OkWithData(aiResponse, c)
}

func (a *assistant) initSSE(c *gin.Context) (http.Flusher, bool) {
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		response.FailWithMessage("streaming unsupported", c)
		return nil, false
	}

	c.Writer.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")
	c.Writer.WriteHeader(http.StatusOK)

	a.emitSSE(c, flusher, "meta", gin.H{"ok": true, "ts": time.Now().Unix()})
	return flusher, true
}

func (a *assistant) emitSSE(c *gin.Context, flusher http.Flusher, event string, payload any) {
	c.SSEvent(event, payload)
	flusher.Flush()
}

func (a *assistant) proxyStream(c *gin.Context, body io.Reader, flusher http.Flusher) {
	reader := bufio.NewReader(body)
	eventName := ""
	dataLines := make([]string, 0, 4)

	flushEvent := func() (done bool) {
		if len(dataLines) == 0 {
			eventName = ""
			return false
		}

		data := strings.TrimSpace(strings.Join(dataLines, "\n"))
		event := strings.ToLower(strings.TrimSpace(eventName))
		eventName = ""
		dataLines = dataLines[:0]

		if data == "" {
			return false
		}
		if data == "[DONE]" || event == "done" {
			a.emitSSE(c, flusher, "done", gin.H{"ok": true})
			return true
		}
		if event == "error" {
			a.emitSSE(c, flusher, "error", gin.H{"message": extractUpstreamError(data, "upstream error")})
			a.emitSSE(c, flusher, "done", gin.H{"ok": true})
			return true
		}

		var payload any
		if err := json.Unmarshal([]byte(data), &payload); err == nil {
			if payloadMap, ok := payload.(map[string]any); ok {
				if msg := extractPayloadError(payloadMap); msg != "" {
					a.emitSSE(c, flusher, "error", gin.H{"message": msg})
					a.emitSSE(c, flusher, "done", gin.H{"ok": true})
					return true
				}
				if delta := extractPayloadDelta(payloadMap); delta != "" {
					a.emitSSE(c, flusher, "delta", gin.H{"content": delta})
				}
				if isPayloadDone(payloadMap) {
					a.emitSSE(c, flusher, "done", gin.H{"ok": true})
					return true
				}
				return false
			}
			if delta := extractPayloadDelta(payload); delta != "" {
				a.emitSSE(c, flusher, "delta", gin.H{"content": delta})
			}
			return false
		}

		a.emitSSE(c, flusher, "delta", gin.H{"content": data})
		return false
	}

	for {
		select {
		case <-c.Request.Context().Done():
			return
		default:
		}

		line, err := reader.ReadString('\n')
		if len(line) > 0 {
			line = strings.TrimRight(line, "\r\n")
			switch {
			case line == "":
				if flushEvent() {
					return
				}
			case strings.HasPrefix(line, ":"):
				continue
			case strings.HasPrefix(line, "event:"):
				eventName = strings.TrimSpace(line[len("event:"):])
			case strings.HasPrefix(line, "data:"):
				dataLines = append(dataLines, strings.TrimSpace(line[len("data:"):]))
			default:
				dataLines = append(dataLines, strings.TrimSpace(line))
			}
		}

		if err == nil {
			continue
		}
		if errors.Is(err, io.EOF) {
			if flushEvent() {
				return
			}
			a.emitSSE(c, flusher, "done", gin.H{"ok": true})
			return
		}

		a.emitSSE(c, flusher, "error", gin.H{"message": err.Error()})
		a.emitSSE(c, flusher, "done", gin.H{"ok": true})
		return
	}
}

func extractPayloadError(payload map[string]any) string {
	if msg := firstString(
		payload["message"],
		payload["msg"],
		payload["error_description"],
		payload["detail"],
	); msg != "" {
		return msg
	}
	if errObj, ok := payload["error"].(map[string]any); ok {
		if msg := firstString(errObj["message"], errObj["msg"], errObj["detail"]); msg != "" {
			return msg
		}
	}
	if errText := strings.TrimSpace(toString(payload["error"])); errText != "" {
		return errText
	}
	return ""
}

func extractPayloadDelta(payload any) string {
	switch t := payload.(type) {
	case string:
		return t
	case map[string]any:
		if choices, ok := t["choices"].([]any); ok {
			parts := make([]string, 0, len(choices))
			for _, row := range choices {
				choice, ok := row.(map[string]any)
				if !ok {
					continue
				}
				if deltaObj, ok := choice["delta"].(map[string]any); ok {
					if content := toString(deltaObj["content"]); content != "" {
						parts = append(parts, content)
					}
				}
				if messageObj, ok := choice["message"].(map[string]any); ok {
					if content := toString(messageObj["content"]); content != "" {
						parts = append(parts, content)
					}
				}
				if text := toString(choice["text"]); text != "" {
					parts = append(parts, text)
				}
			}
			if len(parts) > 0 {
				return strings.Join(parts, "")
			}
		}

		if deltaObj, ok := t["delta"].(map[string]any); ok {
			if content := toString(deltaObj["content"]); content != "" {
				return content
			}
		}
		if content := toString(t["content"]); content != "" {
			return content
		}
		if output := toString(t["output_text"]); output != "" {
			return output
		}
	}
	return ""
}

func isPayloadDone(payload map[string]any) bool {
	if toBool(payload["done"]) || toBool(payload["is_end"]) || toBool(payload["finish"]) {
		return true
	}

	choices, ok := payload["choices"].([]any)
	if !ok {
		return false
	}
	for _, row := range choices {
		choice, ok := row.(map[string]any)
		if !ok {
			continue
		}
		if strings.TrimSpace(toString(choice["finish_reason"])) != "" {
			return true
		}
	}
	return false
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

func toBool(v any) bool {
	switch t := v.(type) {
	case bool:
		return t
	case string:
		s := strings.TrimSpace(strings.ToLower(t))
		return s == "true" || s == "1" || s == "yes"
	case float64:
		return t != 0
	case float32:
		return t != 0
	case int:
		return t != 0
	case int64:
		return t != 0
	case int32:
		return t != 0
	case uint:
		return t != 0
	case uint64:
		return t != 0
	case uint32:
		return t != 0
	default:
		return false
	}
}
