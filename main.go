package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"weibo-group-sender/auth"
	"weibo-group-sender/config"
	"weibo-group-sender/weibo"
)

func main() {
	fmt.Println("===========================================")
	fmt.Println("       微博群消息发送工具")
	fmt.Println("===========================================\n")

	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("加载配置失败: %v\n", err)
		os.Exit(1)
	}

	// 检查是否有 Cookie
	cookieStr := cfg.GetCookieString()
	if cookieStr == "" {
		fmt.Println("未检测到登录信息，需要先登录微博...")
		fmt.Println("\n即将打开浏览器，请完成登录...")
		fmt.Print("按回车键继续...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')

		// 执行自动登录
		result, err := auth.AutoLogin()
		if err != nil || !result.Success {
			fmt.Printf("\n登录失败: %v\n", err)
			os.Exit(1)
		}

		// 保存 Cookie
		cfg.Cookies = result.Cookies
		if err := cfg.Save(); err != nil {
			fmt.Printf("保存配置失败: %v\n", err)
			os.Exit(1)
		}

		cookieStr = cfg.GetCookieString()
		fmt.Println("\n登录信息已保存！")
	} else {
		fmt.Println("✓ 已检测到登录信息")
	}

	// 创建发送器
	sender := weibo.NewSender(cookieStr)

	// 主流程
	reader := bufio.NewReader(os.Stdin)

	// 1. 搜索群组
	fmt.Println("\n===========================================")
	fmt.Print("请输入群组名称关键词: ")
	keyword, _ := reader.ReadString('\n')
	keyword = strings.TrimSpace(keyword)

	if keyword == "" {
		fmt.Println("关键词不能为空")
		os.Exit(1)
	}

	fmt.Println("\n正在搜索群组...")
	groups, err := sender.SearchGroups(keyword, cfg.Source)
	if err != nil {
		fmt.Printf("搜索失败: %v\n", err)

		// 如果是认证错误，提示重新登录
		if strings.Contains(err.Error(), "401") || strings.Contains(err.Error(), "403") {
			fmt.Println("\n登录信息已过期，请重新运行程序并登录")
		}
		os.Exit(1)
	}

	if len(groups) == 0 {
		fmt.Println("未找到匹配的群组")
		os.Exit(0)
	}

	// 2. 显示搜索结果并让用户选择
	fmt.Printf("\n找到 %d 个群组：\n", len(groups))
	fmt.Println("-------------------------------------------")
	for i, group := range groups {
		fmt.Printf("%d. %s (ID: %d)\n", i+1, group.GroupName, group.GID)
	}
	fmt.Println("-------------------------------------------")

	fmt.Print("\n请选择要发送的群组（多个用逗号分隔，如: 1,2,3 或输入 all 选择全部）: ")
	selectionStr, _ := reader.ReadString('\n')
	selectionStr = strings.TrimSpace(selectionStr)

	if selectionStr == "" {
		fmt.Println("未选择任何群组")
		os.Exit(0)
	}

	// 解析选择
	var selectedGroups []weibo.GroupInfo
	if strings.ToLower(selectionStr) == "all" {
		selectedGroups = groups
	} else {
		selections := strings.Split(selectionStr, ",")
		for _, sel := range selections {
			sel = strings.TrimSpace(sel)
			index, err := strconv.Atoi(sel)
			if err != nil || index < 1 || index > len(groups) {
				fmt.Printf("无效的选择: %s\n", sel)
				continue
			}
			selectedGroups = append(selectedGroups, groups[index-1])
		}
	}

	if len(selectedGroups) == 0 {
		fmt.Println("未选择有效的群组")
		os.Exit(0)
	}

	// 3. 输入消息内容
	fmt.Printf("\n已选择 %d 个群组：\n", len(selectedGroups))
	for _, group := range selectedGroups {
		fmt.Printf("  - %s\n", group.GroupName)
	}

	fmt.Print("\n请输入要发送的消息内容: ")
	content, _ := reader.ReadString('\n')
	content = strings.TrimSpace(content)

	if content == "" {
		fmt.Println("消息内容不能为空")
		os.Exit(0)
	}

	// 4. 批量发送消息
	fmt.Println("\n===========================================")
	fmt.Printf("开始向 %d 个群组发送消息...\n", len(selectedGroups))
	fmt.Println("===========================================\n")

	successCount := 0
	failCount := 0

	for i, group := range selectedGroups {
		fmt.Printf("[%d/%d] 正在向 '%s' 发送消息...", i+1, len(selectedGroups), group.GroupName)

		groupIDStr := fmt.Sprintf("%d", group.GID)
		err := sender.SendSimpleMessage(groupIDStr, content, cfg.Source)

		if err != nil {
			fmt.Printf(" ✗ 失败\n")
			fmt.Printf("      错误: %v\n", err)
			failCount++
		} else {
			fmt.Printf(" ✓ 成功\n")
			successCount++
		}

		// 如果不是最后一个群组，等待一段时间
		if i < len(selectedGroups)-1 {
			waitTime := time.Duration(cfg.SendDelay) * time.Second
			fmt.Printf("      等待 %v 后继续...\n\n", waitTime)
			time.Sleep(waitTime)
		}
	}

	// 5. 显示结果统计
	fmt.Println("\n===========================================")
	fmt.Println("发送完成！")
	fmt.Printf("成功: %d 个群组\n", successCount)
	fmt.Printf("失败: %d 个群组\n", failCount)
	fmt.Println("===========================================")
}
