package main

import (
	"testing"
	"weibo-group-sender/config"
)

// TestConfigLoad 测试配置加载
func TestConfigLoad(t *testing.T) {
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("加载配置失败: %v", err)
	}

	if cfg == nil {
		t.Fatal("配置对象为空")
	}

	if cfg.Cookies == nil {
		t.Fatal("Cookies 映射为空")
	}

	if cfg.Source == "" {
		t.Fatal("Source 为空")
	}

	t.Logf("配置加载成功，Source: %s", cfg.Source)
}

// TestCookieStringParsing 测试 Cookie 字符串解析
func TestCookieStringParsing(t *testing.T) {
	cfg := &config.Config{
		Cookies: make(map[string]string),
		Source:  "209678993",
	}

	// 测试设置 Cookie
	testCookieStr := "SCF=test1; SUB=test2; SUBP=test3"
	cfg.SetCookiesFromString(testCookieStr)

	if len(cfg.Cookies) != 3 {
		t.Fatalf("期望解析 3 个 Cookie，实际得到 %d 个", len(cfg.Cookies))
	}

	if cfg.Cookies["SCF"] != "test1" {
		t.Errorf("SCF Cookie 值不正确，期望 'test1'，实际 '%s'", cfg.Cookies["SCF"])
	}

	if cfg.Cookies["SUB"] != "test2" {
		t.Errorf("SUB Cookie 值不正确，期望 'test2'，实际 '%s'", cfg.Cookies["SUB"])
	}

	if cfg.Cookies["SUBP"] != "test3" {
		t.Errorf("SUBP Cookie 值不正确，期望 'test3'，实际 '%s'", cfg.Cookies["SUBP"])
	}

	// 测试获取 Cookie 字符串
	cookieStr := cfg.GetCookieString()
	if cookieStr == "" {
		t.Fatal("Cookie 字符串为空")
	}

	t.Logf("Cookie 字符串: %s", cookieStr)
}

// TestCookieStringGeneration 测试 Cookie 字符串生成
func TestCookieStringGeneration(t *testing.T) {
	cfg := &config.Config{
		Cookies: map[string]string{
			"SCF":  "value1",
			"SUB":  "value2",
			"SUBP": "value3",
		},
		Source: "209678993",
	}

	cookieStr := cfg.GetCookieString()
	if cookieStr == "" {
		t.Fatal("生成的 Cookie 字符串为空")
	}

	// 验证包含所有键
	if !contains(cookieStr, "SCF=value1") {
		t.Error("Cookie 字符串不包含 SCF")
	}
	if !contains(cookieStr, "SUB=value2") {
		t.Error("Cookie 字符串不包含 SUB")
	}
	if !contains(cookieStr, "SUBP=value3") {
		t.Error("Cookie 字符串不包含 SUBP")
	}

	t.Logf("生成的 Cookie 字符串: %s", cookieStr)
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

