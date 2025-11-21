package config

import (
	"encoding/json"
	"os"
)

// Config 配置结构
type Config struct {
	Cookies   map[string]string `json:"cookies"`    // 存储 Cookie 键值对
	Source    string            `json:"source"`     // source 参数
	SendDelay int               `json:"send_delay"` // 发送延迟（秒），默认 2 秒
}

const configFileName = "config.json"

// GetConfigPath 获取配置文件路径
func GetConfigPath() (string, error) {
	return configFileName, nil
}

// Load 加载配置
func Load() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	// 如果配置文件不存在，返回空配置
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &Config{
			Cookies:   make(map[string]string),
			Source:    "209678993",
			SendDelay: 2,
		}, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	if cfg.Cookies == nil {
		cfg.Cookies = make(map[string]string)
	}
	if cfg.Source == "" {
		cfg.Source = "209678993"
	}
	if cfg.SendDelay <= 0 {
		cfg.SendDelay = 2
	}

	return &cfg, nil
}

// Save 保存配置
func (c *Config) Save() error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// GetCookieString 获取 Cookie 字符串
func (c *Config) GetCookieString() string {
	if len(c.Cookies) == 0 {
		return ""
	}

	cookieStr := ""
	for key, value := range c.Cookies {
		if cookieStr != "" {
			cookieStr += "; "
		}
		cookieStr += key + "=" + value
	}
	return cookieStr
}

// SetCookiesFromString 从字符串设置 Cookies
func (c *Config) SetCookiesFromString(cookieStr string) {
	c.Cookies = make(map[string]string)
	// 简单解析 Cookie 字符串
	cookies := splitCookies(cookieStr)
	for _, cookie := range cookies {
		if key, value, ok := parseCookie(cookie); ok {
			c.Cookies[key] = value
		}
	}
}

func splitCookies(cookieStr string) []string {
	var cookies []string
	var current string
	for _, char := range cookieStr {
		if char == ';' {
			if current != "" {
				cookies = append(cookies, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}
	if current != "" {
		cookies = append(cookies, current)
	}
	return cookies
}

func parseCookie(cookie string) (string, string, bool) {
	for i, char := range cookie {
		if char == '=' {
			key := cookie[:i]
			value := cookie[i+1:]
			// 去除空格
			key = trimSpace(key)
			value = trimSpace(value)
			return key, value, true
		}
	}
	return "", "", false
}

func trimSpace(s string) string {
	start := 0
	end := len(s)
	for start < end && s[start] == ' ' {
		start++
	}
	for end > start && s[end-1] == ' ' {
		end--
	}
	return s[start:end]
}
