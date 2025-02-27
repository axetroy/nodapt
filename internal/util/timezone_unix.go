//go:build unix

package util

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// GetTimezone 获取系统时区
func getTimezone() (string, error) {

	if timezone, err := getTimeZoneFromLink(); err == nil {
		return timezone, nil
	}

	if timezone, err := getTimeZoneFromTimeDateCtl(); err == nil {
		return timezone, nil
	}

	if timezone, err := getTimeZoneFromFile(); err == nil {
		return timezone, nil
	}

	if timezone, err := getTimeZoneFromDate(); err == nil {
		return timezone, nil
	}

	if timezone, err := getTimeZoneFromEnv(); err == nil {
		return timezone, nil
	}

	return "", errors.New("can not get timezone")
}

func getTimeZoneFromLink() (string, error) {
	cmd := exec.Command("readlink", "/etc/localtime")
	output, err := cmd.Output()

	if err != nil {
		return "", err
	}

	// 提取时区名称
	tzPath := strings.TrimSpace(string(output))
	parts := strings.Split(tzPath, "/")

	if len(parts) >= 2 {
		timeZone := parts[len(parts)-2] + "/" + parts[len(parts)-1]
		return timeZone, nil
	}

	return "", errors.New("can not get timezone from /etc/localtime")
}

// getTimeZone 通过 `timedatectl` 命令获取系统时区
func getTimeZoneFromTimeDateCtl() (string, error) {
	cmd := exec.Command("timedatectl")

	// 捕获标准输出
	var out bytes.Buffer
	cmd.Stdout = &out

	// 运行命令
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("运行 timedatectl 失败: %v", err)
	}

	// 解析输出
	output := out.String()
	re := regexp.MustCompile(`Time zone:\s+([^\s]+)`)
	match := re.FindStringSubmatch(output)

	if len(match) > 1 {
		return match[1], nil
	}

	return "", fmt.Errorf("无法解析时区信息")
}

// getTimeZoneFromFile 读取 /etc/timezone 获取系统时区
func getTimeZoneFromFile() (string, error) {
	// 读取 /etc/timezone 文件内容
	data, err := os.ReadFile("/etc/timezone")
	if err != nil {
		return "", fmt.Errorf("无法读取 /etc/timezone: %v", err)
	}

	// 处理内容，去除可能的换行符
	timezone := strings.TrimSpace(string(data))

	if timezone == "" {
		return "", fmt.Errorf("时区文件为空")
	}

	return timezone, nil
}

// getTimeZoneFromDate 使用 `date` 命令获取时区
func getTimeZoneFromDate() (string, error) {
	cmd := exec.Command("date", "+%Z %z") // 获取时区缩写和 UTC 偏移
	var out bytes.Buffer
	cmd.Stdout = &out

	// 运行命令
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("运行 date 命令失败: %v", err)
	}

	// 处理输出
	timezone := strings.TrimSpace(out.String())

	if timezone == "" {
		return "", fmt.Errorf("无法获取时区信息")
	}

	return timezone, nil
}

// getTimeZoneFromEnv 通过环境变量 TZ 获取系统时区
func getTimeZoneFromEnv() (string, error) {
	timezone := os.Getenv("TZ")

	if timezone == "" {
		return "", fmt.Errorf("环境变量 TZ 未设置")
	}

	return timezone, nil
}
