package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	supaReq "github.com/icosmos-space/iadmin/server/plugin/supabase/model/request"
	"github.com/icosmos-space/iadmin/server/plugin/supabase/plugin"
)

type user struct{}

type adminUsersResponse struct {
	Users []map[string]any `json:"users"`
}

func (s *user) GetUserList(info supaReq.SupabaseUserSearch) ([]map[string]any, int64, error) {
	if info.Page <= 0 {
		info.Page = 1
	}
	if info.PageSize <= 0 {
		info.PageSize = 10
	}
	if info.PageSize > 100 {
		info.PageSize = 100
	}

	u, key, err := getSupabaseConfig()
	if err != nil {
		return nil, 0, err
	}

	qs := url.Values{}
	qs.Set("page", fmt.Sprintf("%d", info.Page))
	qs.Set("per_page", fmt.Sprintf("%d", info.PageSize))
	listURL := strings.TrimRight(u, "/") + "/auth/v1/admin/users?" + qs.Encode()

	req, err := http.NewRequest(http.MethodGet, listURL, nil)
	if err != nil {
		return nil, 0, err
	}
	setHeaders(req, key)

	client := &http.Client{Timeout: 12 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, 0, fmt.Errorf("supabase 用户列表请求失败: %s", resp.Status)
	}

	var payload adminUsersResponse
	if err = json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, 0, err
	}

	if info.Keyword != "" {
		payload.Users = filterUsersByKeyword(payload.Users, info.Keyword)
	}

	// Supabase 当前列表接口未直接返回 total，这里在分页维度先返回当前页条数。
	return payload.Users, int64(len(payload.Users)), nil
}

func (s *user) UpdateUserPassword(data supaReq.UpdateSupabaseUserPassword) error {
	u, key, err := getSupabaseConfig()
	if err != nil {
		return err
	}
	if strings.TrimSpace(data.UserID) == "" {
		return errors.New("用户ID不能为空")
	}

	bodyMap := map[string]any{
		"password": data.NewPassword,
	}
	body, _ := json.Marshal(bodyMap)
	targetURL := strings.TrimRight(u, "/") + "/auth/v1/admin/users/" + url.PathEscape(strings.TrimSpace(data.UserID))

	req, err := http.NewRequest(http.MethodPut, targetURL, bytes.NewReader(body))
	if err != nil {
		return err
	}
	setHeaders(req, key)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 12 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("supabase 修改密码失败: %s", resp.Status)
	}
	return nil
}

func getSupabaseConfig() (string, string, error) {
	u := strings.TrimSpace(plugin.Config.URL)
	key := strings.TrimSpace(plugin.Config.ServiceRoleKey)
	if u == "" {
		return "", "", errors.New("supabase.url 未配置")
	}
	if key == "" {
		return "", "", errors.New("supabase.service-role-key 未配置")
	}
	return u, key, nil
}

func setHeaders(req *http.Request, key string) {
	req.Header.Set("apikey", key)
	req.Header.Set("Authorization", "Bearer "+key)
}

func filterUsersByKeyword(users []map[string]any, keyword string) []map[string]any {
	kw := strings.ToLower(strings.TrimSpace(keyword))
	if kw == "" {
		return users
	}
	out := make([]map[string]any, 0, len(users))
	for _, u := range users {
		email := strings.ToLower(fmt.Sprintf("%v", u["email"]))
		phone := strings.ToLower(fmt.Sprintf("%v", u["phone"]))
		id := strings.ToLower(fmt.Sprintf("%v", u["id"]))
		if strings.Contains(email, kw) || strings.Contains(phone, kw) || strings.Contains(id, kw) {
			out = append(out, u)
		}
	}
	return out
}
