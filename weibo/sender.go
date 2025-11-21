package weibo

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// SendMessageRequest 发送消息请求参数
type SendMessageRequest struct {
	SetTimeout  string // setTimeout
	Content     string // content 消息内容
	ID          string // id 群组ID
	MediaType   string // media_type
	Annotations string // annotations JSON字符串
	IsEncoded   string // is_encoded
	Source      string // source
}

// SendMessageResponse 发送消息响应
type SendMessageResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Sender 微博消息发送器
type Sender struct {
	cookieString string
}

// NewSender 创建新的发送器
func NewSender(cookieString string) *Sender {
	return &Sender{
		cookieString: cookieString,
	}
}

// SendGroupMessage 发送群组消息
func (s *Sender) SendGroupMessage(req *SendMessageRequest) error {
	apiURL := "https://api.weibo.com/webim/groupchat/send_message.json"

	// 构建表单数据
	formData := url.Values{}
	formData.Set("setTimeout", req.SetTimeout)
	formData.Set("content", req.Content)
	formData.Set("id", req.ID)
	formData.Set("media_type", req.MediaType)
	formData.Set("annotations", req.Annotations)
	formData.Set("is_encoded", req.IsEncoded)
	formData.Set("source", req.Source)

	// 创建请求
	httpReq, err := http.NewRequest("POST", apiURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	httpReq.Header.Set("Accept", "application/json, text/plain, */*")
	httpReq.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	httpReq.Header.Set("Cache-Control", "no-cache")
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	httpReq.Header.Set("Cookie", s.cookieString)
	httpReq.Header.Set("Origin", "https://api.weibo.com")
	httpReq.Header.Set("Pragma", "no-cache")
	httpReq.Header.Set("Priority", "u=1, i")
	httpReq.Header.Set("Referer", "https://api.weibo.com/chat")
	httpReq.Header.Set("Sec-Ch-Ua", `"Chromium";v="142", "Google Chrome";v="142", "Not_A Brand";v="99"`)
	httpReq.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	httpReq.Header.Set("Sec-Ch-Ua-Platform", `"macOS"`)
	httpReq.Header.Set("Sec-Fetch-Dest", "empty")
	httpReq.Header.Set("Sec-Fetch-Mode", "cors")
	httpReq.Header.Set("Sec-Fetch-Site", "same-origin")
	httpReq.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取响应失败: %w", err)
	}

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	fmt.Printf("消息发送成功！响应: %s\n", string(body))
	return nil
}

// SendSimpleMessage 发送简单消息（使用默认参数）
func (s *Sender) SendSimpleMessage(groupID, content string, source string) error {
	req := &SendMessageRequest{
		SetTimeout:  "50",
		Content:     content,
		ID:          groupID,
		MediaType:   "0",
		Annotations: `{"webchat":1,"clientid":"mx8y819bfek0ds5fpi18ltpgulg1goq"}`,
		IsEncoded:   "0",
		Source:      source,
	}
	return s.SendGroupMessage(req)
}

// GroupInfo 群组信息
type GroupInfo struct {
	GID                  int64    `json:"gid"`
	GroupName            string   `json:"groupname"`
	GroupAvatarURLs      []string `json:"group_avatar_urls"`
	ProfileImageURL      string   `json:"profile_image_url"`
	AvatarLarge          string   `json:"avatar_large"`
	RoundProfileImageURL string   `json:"round_profile_image_url"`
	RoundAvatarLarge     string   `json:"round_avatar_large"`
}

// SearchResponse 搜索响应
type SearchResponse struct {
	Contacts struct {
		Num      int `json:"num"`
		TotalNum int `json:"total_num"`
	} `json:"contacts"`
	Groups struct {
		GroupMemberNum  int         `json:"group_member_num"`
		GroupResultList []GroupInfo `json:"group_result_list"`
		GroupNum        int         `json:"group_num"`
		TotalNum        int         `json:"total_num"`
	} `json:"groups"`
}

// SearchGroups 搜索群组
func (s *Sender) SearchGroups(keyword string, source string) ([]GroupInfo, error) {
	// 构建 URL
	timestamp := time.Now().UnixMilli()
	apiURL := fmt.Sprintf(
		"https://api.weibo.com/webim/2/direct_messages/messageboxsearch.json?types=contact,group&key=%s&pagecount=20&source=%s&t=%d",
		url.QueryEscape(keyword),
		source,
		timestamp,
	)

	// 创建请求
	httpReq, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头
	httpReq.Header.Set("Accept", "application/json, text/plain, */*")
	httpReq.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	httpReq.Header.Set("Cache-Control", "no-cache")
	httpReq.Header.Set("Cookie", s.cookieString)
	httpReq.Header.Set("Pragma", "no-cache")
	httpReq.Header.Set("Priority", "u=1, i")
	httpReq.Header.Set("Referer", "https://api.weibo.com/chat")
	httpReq.Header.Set("Sec-Ch-Ua", `"Chromium";v="142", "Google Chrome";v="142", "Not_A Brand";v="99"`)
	httpReq.Header.Set("Sec-Ch-Ua-Mobile", "?0")
	httpReq.Header.Set("Sec-Ch-Ua-Platform", `"macOS"`)
	httpReq.Header.Set("Sec-Fetch-Dest", "empty")
	httpReq.Header.Set("Sec-Fetch-Mode", "cors")
	httpReq.Header.Set("Sec-Fetch-Site", "same-origin")
	httpReq.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 检查状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	// 解析 JSON 响应
	var searchResp SearchResponse
	if err := json.Unmarshal(body, &searchResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w, 响应内容: %s", err, string(body))
	}

	return searchResp.Groups.GroupResultList, nil
}
