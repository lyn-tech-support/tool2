package auth

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

// LoginResult 登录结果
type LoginResult struct {
	Cookies map[string]string
	Success bool
	Error   error
}

// AutoLogin 自动登录微博并获取 Cookie
func AutoLogin() (*LoginResult, error) {
	result := &LoginResult{
		Cookies: make(map[string]string),
	}

	// 创建 Chrome 实例
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false), // 显示浏览器窗口，方便用户登录
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/142.0.0.0 Safari/537.36"),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	// 设置超时
	ctx, cancel = context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	// 执行登录流程
	err := chromedp.Run(ctx,
		// 访问微博登录页面
		chromedp.Navigate("https://weibo.com/login.php"),

		// 等待用户手动登录
		chromedp.Sleep(2*time.Second),
	)

	if err != nil {
		result.Error = fmt.Errorf("导航到登录页面失败: %w", err)
		return result, result.Error
	}

	fmt.Println("\n===========================================")
	fmt.Println("请在打开的浏览器窗口中完成登录")
	fmt.Println("登录成功后，程序将自动获取 Cookie")
	fmt.Println("===========================================\n")

	// 等待登录成功（检测是否跳转到主页或聊天页面）
	err = chromedp.Run(ctx,
		// 等待登录成功的标志
		chromedp.WaitVisible(`body`, chromedp.ByQuery),
	)

	if err != nil {
		result.Error = fmt.Errorf("等待页面加载失败: %w", err)
		return result, result.Error
	}

	// 轮询检查是否登录成功
	maxAttempts := 150 // 5分钟，每2秒检查一次
	for i := 0; i < maxAttempts; i++ {
		var currentURL string
		err = chromedp.Run(ctx,
			chromedp.Location(&currentURL),
		)

		if err != nil {
			result.Error = fmt.Errorf("获取当前URL失败: %w", err)
			return result, result.Error
		}

		// 检查是否已经登录（URL不再是登录页面）
		if currentURL != "https://weibo.com/login.php" &&
			!contains(currentURL, "login") {
			fmt.Println("检测到登录成功！正在获取 Cookie...")
			break
		}

		if i == maxAttempts-1 {
			result.Error = fmt.Errorf("登录超时，请重试")
			return result, result.Error
		}

		time.Sleep(2 * time.Second)
	}

	// 访问聊天页面以确保获取完整的 Cookie
	err = chromedp.Run(ctx,
		chromedp.Navigate("https://api.weibo.com/chat"),
		chromedp.Sleep(3*time.Second),
	)

	if err != nil {
		result.Error = fmt.Errorf("访问聊天页面失败: %w", err)
		return result, result.Error
	}

	// 获取 Cookie
	var cookies []*network.Cookie
	err = chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			cookies, err = network.GetCookies().Do(ctx)
			return err
		}),
	)

	if err != nil {
		result.Error = fmt.Errorf("获取 Cookie 失败: %w", err)
		return result, result.Error
	}

	// 提取需要的 Cookie
	requiredCookies := []string{"SCF", "SUB", "SUBP", "ALF", "_s_tentry", "Apache", "SINAGLOBAL", "ULV"}
	for _, cookie := range cookies {
		for _, required := range requiredCookies {
			if cookie.Name == required {
				result.Cookies[cookie.Name] = cookie.Value
				break
			}
		}
	}

	if len(result.Cookies) == 0 {
		result.Error = fmt.Errorf("未能获取到有效的 Cookie，请确保已成功登录")
		return result, result.Error
	}

	result.Success = true
	fmt.Printf("\n成功获取 %d 个 Cookie\n", len(result.Cookies))

	return result, nil
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
