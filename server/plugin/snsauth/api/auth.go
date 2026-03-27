package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/icosmos-space/iadmin/server/model/common/response"
	snsReq "github.com/icosmos-space/iadmin/server/plugin/snsauth/model/request"
	"github.com/icosmos-space/iadmin/server/utils"
)

type auth struct{}

func (a *auth) GetProviderList(c *gin.Context) {
	list, err := serviceAuth.ListProviders()
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(list, c)
}

func (a *auth) GetEnabledProviderList(c *gin.Context) {
	list, err := serviceAuth.ListEnabledProviders()
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(list, c)
}

func (a *auth) UpdateProviderConfig(c *gin.Context) {
	var in snsReq.UpdateProviderConfigReq
	if err := c.ShouldBindJSON(&in); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := serviceAuth.UpdateProviderConfig(in); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

func (a *auth) GetLoginURL(c *gin.Context) {
	var q snsReq.BuildURLReq
	if err := c.ShouldBindQuery(&q); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	url, err := serviceAuth.BuildAuthURL(q.Provider, "login", 0)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(gin.H{"url": url}, c)
}

func (a *auth) GetBindURL(c *gin.Context) {
	var q snsReq.BuildURLReq
	if err := c.ShouldBindQuery(&q); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	url, err := serviceAuth.BuildAuthURL(q.Provider, "bind", utils.GetUserID(c))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(gin.H{"url": url}, c)
}

func (a *auth) Callback(c *gin.Context) {
	provider := c.Param("provider")
	code := c.Query("code")
	state := c.Query("state")
	result, err := serviceAuth.HandleCallback(provider, code, state)
	if err != nil {
		a.renderCallbackHTML(c, false, provider, nil, err.Error())
		return
	}
	a.renderCallbackHTML(c, true, provider, result, "")
}

func (a *auth) GetMyBindings(c *gin.Context) {
	list, err := serviceAuth.GetMyBindings(utils.GetUserID(c))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(list, c)
}

func (a *auth) Unbind(c *gin.Context) {
	provider := c.Query("provider")
	if provider == "" {
		response.FailWithMessage("provider不能为空", c)
		return
	}
	if err := serviceAuth.Unbind(utils.GetUserID(c), provider); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("解绑成功", c)
}

func (a *auth) TelegramLogin(c *gin.Context) {
	var in snsReq.TelegramAuthReq
	if err := c.ShouldBindJSON(&in); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	res, err := serviceAuth.HandleTelegramAuth(in, "login", 0)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(res, c)
}

func (a *auth) TelegramBind(c *gin.Context) {
	var in snsReq.TelegramAuthReq
	if err := c.ShouldBindJSON(&in); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	res, err := serviceAuth.HandleTelegramAuth(in, "bind", utils.GetUserID(c))
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(res, c)
}

func (a *auth) renderCallbackHTML(c *gin.Context, ok bool, provider string, data any, errMsg string) {
	payload := fmt.Sprintf(`{"type":"SNS_AUTH_RESULT","ok":%t,"provider":"%s","data":%s,"error":"%s"}`,
		ok, provider, mustJSON(data), escapeJS(errMsg))
	html := "<!doctype html><html><body><script>" +
		"if(window.opener){window.opener.postMessage(" + payload + ", '*');}" +
		"window.close();" +
		"</script>处理中...</body></html>"
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(html))
}

func mustJSON(v any) string {
	if v == nil {
		return "null"
	}
	b, _ := json.Marshal(v)
	return string(b)
}

func escapeJS(s string) string {
	b, _ := json.Marshal(s)
	return string(b[1 : len(b)-1])
}
