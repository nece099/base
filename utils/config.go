package utils

import (
	"strings"

	"github.com/nece099/base/datasource"
)

type ServerConfig struct {
	Platform   string
	NodeID     string
	LogLevel   string
	ScreenLog  bool
	DataSource []datasource.DataSource
}

var serverConfig = &ServerConfig{}

func PurifyConfig(content string) string {

	lines := strings.Split(content, "\n")

	res := ""

	for _, line := range lines {
		trimline := strings.TrimSpace(line)
		idx := strings.Index(trimline, "//")
		if idx == -1 { // 无注释,返回trim过后的行
			res = res + trimline + "\n"
			continue
		}

		if idx != 0 {
			res = res + trimline + "\n"
			continue
		}

		// 截取注释符号之前的字符串
		// 忽略注释
	}

	return res
}
