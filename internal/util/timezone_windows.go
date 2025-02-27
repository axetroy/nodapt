//go:build windows

package util

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"golang.org/x/sys/windows/registry"
)

func getTimezone() (string, error) {
	if timezone, err := getTimezoneViaRegistry(); err == nil {
		return timezone, nil
	}

	if timezone, err := getTimezoneViaCommand(); err == nil {
		return timezone, nil
	}

	return "", errors.New("无法获取 Windows 时区")
}

func getTimezoneViaRegistry() (string, error) {
	// 通过 Windows 注册表获取
	key, err := registry.OpenKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\TimeZoneInformation`, registry.QUERY_VALUE)

	if err == nil {
		defer key.Close()
		windowsTZ, _, err := key.GetStringValue("TimeZoneKeyName")
		if err == nil {
			return convertWindowsToIana(windowsTZ)
		}
	}

	return "", errors.New("无法获取 Windows 时区")
}

func getTimezoneViaCommand() (string, error) {
	// 通过 tzutil /g 获取 Windows 时区名
	cmd := exec.Command("tzutil", "/g")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err == nil {
		windowsTZ := strings.TrimSpace(out.String())
		return convertWindowsToIana(windowsTZ)
	}

	return "", errors.New("无法获取 Windows 时区")
}

// convertWindowsToIana 将 Windows 时区转换为 IANA 时区
func convertWindowsToIana(windowsTZ string) (string, error) {
	windowsToIana := map[string]string{
		"China Standard Time":              "Asia/Shanghai",
		"Pacific Standard Time":            "America/Los_Angeles",
		"Eastern Standard Time":            "America/New_York",
		"Central European Standard Time":   "Europe/Berlin",
		"Greenwich Standard Time":          "Europe/London",
		"India Standard Time":              "Asia/Kolkata",
		"Japan Standard Time":              "Asia/Tokyo",
		"Australian Eastern Standard Time": "Australia/Sydney",
		"Central Standard Time":            "America/Chicago",
		"Mountain Standard Time":           "America/Denver",
		"Alaskan Standard Time":            "America/Anchorage",
		"Hawaiian Standard Time":           "Pacific/Honolulu",
		"Atlantic Standard Time":           "America/Halifax",
		"Newfoundland Standard Time":       "America/St_Johns",
		"Arabian Standard Time":            "Asia/Riyadh",
		"Israel Standard Time":             "Asia/Jerusalem",
		"Russian Standard Time":            "Europe/Moscow",
		"South Africa Standard Time":       "Africa/Johannesburg",
		"Singapore Standard Time":          "Asia/Singapore",
		"Taipei Standard Time":             "Asia/Taipei",
		"Korea Standard Time":              "Asia/Seoul",
		"Western European Standard Time":   "Europe/Lisbon",
		"Eastern European Standard Time":   "Europe/Athens",
		// 其他时区可以在这里添加
	}

	if ianaTZ, ok := windowsToIana[windowsTZ]; ok {
		return ianaTZ, nil
	}
	return "", fmt.Errorf("未知 Windows 时区: %s", windowsTZ)
}
