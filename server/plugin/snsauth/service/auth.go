package service

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	neturl "net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/officialAccount"
	"github.com/icosmos-space/iadmin/server/global"
	"github.com/icosmos-space/iadmin/server/model/system"
	"github.com/icosmos-space/iadmin/server/plugin/snsauth/model"
	snsReq "github.com/icosmos-space/iadmin/server/plugin/snsauth/model/request"
	snsRes "github.com/icosmos-space/iadmin/server/plugin/snsauth/model/response"
	systemService "github.com/icosmos-space/iadmin/server/service/system"
	"github.com/icosmos-space/iadmin/server/utils"
	"golang.org/x/oauth2"
)

type auth struct{}

type authState struct {
	Mode      string
	Provider  string
	UserID    uint
	ExpiredAt time.Time
}

var (
	stateStore = make(map[string]authState)
	stateMux   sync.Mutex
)

func (s *auth) InitDefaultProviders() error {
	def := []model.SnsProviderConfig{
		{Provider: "github", AuthURL: "https://github.com/login/oauth/authorize", TokenURL: "https://github.com/login/oauth/access_token", UserInfoURL: "https://api.github.com/user"},
		{Provider: "feishu", AuthURL: "https://accounts.feishu.cn/open-apis/authen/v1/authorize", TokenURL: "https://open.feishu.cn/open-apis/authen/v1/oidc/access_token", UserInfoURL: "https://open.feishu.cn/open-apis/authen/v1/user_info"},
		{Provider: "wechat", AuthURL: "https://open.weixin.qq.com/connect/oauth2/authorize", TokenURL: "https://api.weixin.qq.com/sns/oauth2/access_token", UserInfoURL: "https://api.weixin.qq.com/sns/userinfo", Scopes: "snsapi_userinfo"},
		{Provider: "telegram"},
	}
	for i := range def {
		row := def[i]
		if err := global.IADMIN_DB.Where("provider = ?", row.Provider).FirstOrCreate(&row).Error; err != nil {
			return err
		}
	}
	return nil
}

func (s *auth) ListProviders() ([]model.SnsProviderConfig, error) {
	var list []model.SnsProviderConfig
	err := global.IADMIN_DB.Order("id asc").Find(&list).Error
	return list, err
}

func (s *auth) ListEnabledProviders() ([]model.SnsProviderConfig, error) {
	var list []model.SnsProviderConfig
	err := global.IADMIN_DB.Where("enabled = ?", true).Order("id asc").Find(&list).Error
	return list, err
}

func (s *auth) UpdateProviderConfig(in snsReq.UpdateProviderConfigReq) error {
	var row model.SnsProviderConfig
	if err := global.IADMIN_DB.Where("provider = ?", in.Provider).First(&row).Error; err != nil {
		return err
	}
	row.Enabled = in.Enabled
	row.ClientID = strings.TrimSpace(in.ClientID)
	row.ClientSecret = strings.TrimSpace(in.ClientSecret)
	row.RedirectURL = strings.TrimSpace(in.RedirectURL)
	row.Scopes = strings.TrimSpace(in.Scopes)
	row.AuthURL = strings.TrimSpace(in.AuthURL)
	row.TokenURL = strings.TrimSpace(in.TokenURL)
	row.UserInfoURL = strings.TrimSpace(in.UserInfoURL)
	if err := global.IADMIN_DB.Save(&row).Error; err != nil {
		return err
	}
	s.writeBackConfigFile(row)
	return nil
}

func (s *auth) BuildAuthURL(provider, mode string, userID uint) (string, error) {
	row, err := s.getProvider(provider)
	if err != nil {
		return "", err
	}
	if !row.Enabled {
		return "", errors.New("该平台未启用")
	}
	if row.AuthURL == "" || row.TokenURL == "" {
		return "", errors.New("请先配置授权地址和Token地址")
	}
	if provider == "telegram" {
		return "", errors.New("telegram请使用登录组件模式，不支持标准OAuth跳转")
	}
	st := s.newState(mode, provider, userID)
	if provider == "wechat" {
		u, _ := neturl.Parse(row.AuthURL)
		q := u.Query()
		q.Set("appid", row.ClientID)
		q.Set("redirect_uri", row.RedirectURL)
		scope := row.Scopes
		if strings.TrimSpace(scope) == "" {
			scope = "snsapi_userinfo"
		}
		q.Set("response_type", "code")
		q.Set("scope", scope)
		q.Set("state", st)
		u.RawQuery = q.Encode()
		return u.String() + "#wechat_redirect", nil
	}
	conf := oauth2.Config{
		ClientID:     row.ClientID,
		ClientSecret: row.ClientSecret,
		RedirectURL:  row.RedirectURL,
		Scopes:       splitScopes(row.Scopes),
		Endpoint: oauth2.Endpoint{
			AuthURL:  row.AuthURL,
			TokenURL: row.TokenURL,
		},
	}
	return conf.AuthCodeURL(st, oauth2.AccessTypeOnline), nil
}

func (s *auth) HandleCallback(provider, code, state string) (*snsRes.SnsAuthResult, error) {
	st, err := s.consumeState(state, provider)
	if err != nil {
		return nil, err
	}

	if provider == "wechat" {
		userInfo, err := s.wechatFetchUserInfo(code)
		if err != nil {
			return nil, err
		}
		openID := pickString(userInfo, "openid")
		if openID == "" {
			return nil, errors.New("微信 openid 为空")
		}
		unionID := pickString(userInfo, "unionid")
		if st.Mode == "bind" {
			if st.UserID == 0 {
				return nil, errors.New("无效绑定状态")
			}
			return s.bindAndBuildResult(st.UserID, provider, openID, unionID, userInfo)
		}
		return s.loginByBind(provider, openID)
	}

	row, err := s.getProvider(provider)
	if err != nil {
		return nil, err
	}
	conf := oauth2.Config{
		ClientID:     row.ClientID,
		ClientSecret: row.ClientSecret,
		RedirectURL:  row.RedirectURL,
		Scopes:       splitScopes(row.Scopes),
		Endpoint: oauth2.Endpoint{
			AuthURL:  row.AuthURL,
			TokenURL: row.TokenURL,
		},
	}
	token, err := conf.Exchange(context.Background(), code)
	if err != nil {
		return nil, err
	}
	userInfo, err := s.fetchUserInfo(row.UserInfoURL, token.AccessToken)
	if err != nil {
		return nil, err
	}
	openID := pickString(userInfo, "id", "open_id", "sub", "user_id")
	if openID == "" {
		return nil, errors.New("第三方用户标识为空")
	}
	unionID := pickString(userInfo, "union_id", "unionId")

	if st.Mode == "bind" {
		if st.UserID == 0 {
			return nil, errors.New("无效绑定状态")
		}
		return s.bindAndBuildResult(st.UserID, provider, openID, unionID, userInfo)
	}
	return s.loginByBind(provider, openID)
}

func (s *auth) HandleTelegramAuth(in snsReq.TelegramAuthReq, mode string, userID uint) (*snsRes.SnsAuthResult, error) {
	row, err := s.getProvider("telegram")
	if err != nil {
		return nil, err
	}
	if !row.Enabled {
		return nil, errors.New("telegram未启用")
	}
	if strings.TrimSpace(row.ClientSecret) == "" {
		return nil, errors.New("请先配置telegram bot token（clientSecret）")
	}
	if err = verifyTelegramAuth(in, row.ClientSecret); err != nil {
		return nil, err
	}
	openID := strconv.FormatInt(in.ID, 10)
	meta := map[string]any{
		"id":         in.ID,
		"first_name": in.FirstName,
		"last_name":  in.LastName,
		"username":   in.Username,
		"photo_url":  in.PhotoURL,
		"auth_date":  in.AuthDate,
	}
	if mode == "bind" {
		if userID == 0 {
			return nil, errors.New("无效绑定状态")
		}
		return s.bindAndBuildResult(userID, "telegram", openID, "", meta)
	}
	return s.loginByBind("telegram", openID)
}

func (s *auth) GetMyBindings(userID uint) ([]model.SnsUserBind, error) {
	var list []model.SnsUserBind
	err := global.IADMIN_DB.Where("user_id = ?", userID).Find(&list).Error
	return list, err
}

func (s *auth) Unbind(userID uint, provider string) error {
	return global.IADMIN_DB.Where("user_id = ? AND provider = ?", userID, provider).Delete(&model.SnsUserBind{}).Error
}

func (s *auth) bindAndBuildResult(userID uint, provider, openID, unionID string, meta map[string]any) (*snsRes.SnsAuthResult, error) {
	bind := model.SnsUserBind{
		UserID:   userID,
		Provider: provider,
		OpenID:   openID,
		UnionID:  unionID,
		Meta:     meta,
	}
	var old model.SnsUserBind
	err := global.IADMIN_DB.Where("provider = ? AND open_id = ?", provider, openID).First(&old).Error
	if err == nil {
		old.UserID = userID
		old.UnionID = unionID
		old.Meta = meta
		if err = global.IADMIN_DB.Save(&old).Error; err != nil {
			return nil, err
		}
	} else {
		if err = global.IADMIN_DB.Create(&bind).Error; err != nil {
			return nil, err
		}
	}
	return &snsRes.SnsAuthResult{
		Mode:     "bind",
		Provider: provider,
		Msg:      "绑定成功",
	}, nil
}

func (s *auth) loginByBind(provider, openID string) (*snsRes.SnsAuthResult, error) {
	var bind model.SnsUserBind
	if err := global.IADMIN_DB.Where("provider = ? AND open_id = ?", provider, openID).First(&bind).Error; err != nil {
		return nil, errors.New("该SNS账号尚未绑定系统用户")
	}
	var user system.SysUser
	if err := global.IADMIN_DB.Preload("Authorities").Preload("Authority").First(&user, bind.UserID).Error; err != nil {
		return nil, err
	}
	systemService.MenuServiceApp.UserAuthorityDefaultRouter(&user)
	token, _, err := utils.LoginToken(&user)
	if err != nil {
		return nil, err
	}
	return &snsRes.SnsAuthResult{
		Mode:     "login",
		Provider: provider,
		Token:    token,
		User:     user,
		Msg:      "登录成功",
	}, nil
}

func (s *auth) fetchUserInfo(userInfoURL, accessToken string) (map[string]any, error) {
	if userInfoURL == "" {
		return nil, errors.New("请配置用户信息地址")
	}
	req, _ := http.NewRequest(http.MethodGet, userInfoURL, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")
	client := http.Client{Timeout: 12 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var raw map[string]any
	if err = json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}
	// 兼容 data 包裹结构
	if data, ok := raw["data"].(map[string]any); ok {
		return data, nil
	}
	return raw, nil
}

func (s *auth) wechatFetchUserInfo(code string) (map[string]any, error) {
	row, err := s.getProvider("wechat")
	if err != nil {
		return nil, err
	}
	appID := strings.TrimSpace(row.ClientID)
	secret := strings.TrimSpace(row.ClientSecret)
	if appID == "" || secret == "" {
		return nil, errors.New("请先配置微信 appid/clientSecret")
	}
	app, err := officialAccount.NewOfficialAccount(&officialAccount.UserConfig{
		AppID:  appID,
		Secret: secret,
		OAuth: officialAccount.OAuth{
			Callback: row.RedirectURL,
			Scopes:   splitScopes(row.Scopes),
		},
	})
	if err != nil {
		return nil, err
	}
	oauthUser, err := app.OAuth.UserFromCode(code)
	if err != nil {
		return nil, err
	}

	user := map[string]any{
		"openid":      oauthUser.GetOpenID(),
		"id":          oauthUser.GetID(),
		"nickname":    oauthUser.GetNickname(),
		"name":        oauthUser.GetName(),
		"avatar":      oauthUser.GetAvatar(),
		"email":       oauthUser.GetEmail(),
		"mobile":      oauthUser.GetMobile(),
		"accessToken": oauthUser.GetAccessToken(),
	}
	if user["openid"] == "" {
		user["openid"] = oauthUser.GetID()
	}
	if raw, e := oauthUser.GetRaw(); e == nil && raw != nil {
		for k, v := range *raw {
			user[k] = v
		}
	}
	if strings.TrimSpace(fmt.Sprintf("%v", user["openid"])) == "" {
		return nil, errors.New("微信用户标识为空")
	}
	return user, nil
}

func (s *auth) getProvider(provider string) (model.SnsProviderConfig, error) {
	var row model.SnsProviderConfig
	err := global.IADMIN_DB.Where("provider = ?", provider).First(&row).Error
	return row, err
}

func (s *auth) newState(mode, provider string, userID uint) string {
	b := make([]byte, 24)
	_, _ = rand.Read(b)
	st := base64.RawURLEncoding.EncodeToString(b)
	stateMux.Lock()
	stateStore[st] = authState{
		Mode:      mode,
		Provider:  provider,
		UserID:    userID,
		ExpiredAt: time.Now().Add(5 * time.Minute),
	}
	stateMux.Unlock()
	return st
}

func (s *auth) consumeState(state, provider string) (authState, error) {
	stateMux.Lock()
	defer stateMux.Unlock()
	st, ok := stateStore[state]
	if !ok {
		return authState{}, errors.New("state无效")
	}
	delete(stateStore, state)
	if st.Provider != provider {
		return authState{}, errors.New("state与平台不匹配")
	}
	if time.Now().After(st.ExpiredAt) {
		return authState{}, errors.New("state已过期")
	}
	return st, nil
}

func (s *auth) writeBackConfigFile(row model.SnsProviderConfig) {
	key := "sns-auth." + row.Provider
	global.IADMIN_VP.Set(key+".enabled", row.Enabled)
	global.IADMIN_VP.Set(key+".client-id", row.ClientID)
	global.IADMIN_VP.Set(key+".client-secret", row.ClientSecret)
	global.IADMIN_VP.Set(key+".redirect-url", row.RedirectURL)
	global.IADMIN_VP.Set(key+".scopes", row.Scopes)
	global.IADMIN_VP.Set(key+".auth-url", row.AuthURL)
	global.IADMIN_VP.Set(key+".token-url", row.TokenURL)
	global.IADMIN_VP.Set(key+".user-info-url", row.UserInfoURL)
	_ = global.IADMIN_VP.WriteConfig()
}

func splitScopes(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	return strings.Fields(strings.ReplaceAll(s, ",", " "))
}

func pickString(m map[string]any, keys ...string) string {
	for _, k := range keys {
		if v, ok := m[k]; ok {
			return strings.TrimSpace(fmt.Sprintf("%v", v))
		}
	}
	return ""
}

func verifyTelegramAuth(in snsReq.TelegramAuthReq, botToken string) error {
	if in.AuthDate <= 0 {
		return errors.New("telegram auth_date 无效")
	}
	// 10分钟过期
	if time.Now().Unix()-in.AuthDate > 600 {
		return errors.New("telegram授权数据已过期")
	}
	data := make([]string, 0, 8)
	data = append(data, "auth_date="+strconv.FormatInt(in.AuthDate, 10))
	data = append(data, "id="+strconv.FormatInt(in.ID, 10))
	if in.FirstName != "" {
		data = append(data, "first_name="+in.FirstName)
	}
	if in.LastName != "" {
		data = append(data, "last_name="+in.LastName)
	}
	if in.PhotoURL != "" {
		data = append(data, "photo_url="+in.PhotoURL)
	}
	if in.Username != "" {
		data = append(data, "username="+in.Username)
	}
	sort.Strings(data)
	checkString := strings.Join(data, "\n")

	secret := sha256.Sum256([]byte(botToken))
	h := hmac.New(sha256.New, secret[:])
	_, _ = h.Write([]byte(checkString))
	sign := hex.EncodeToString(h.Sum(nil))
	if !hmac.Equal([]byte(sign), []byte(strings.ToLower(in.Hash))) {
		return errors.New("telegram签名校验失败")
	}
	return nil
}
