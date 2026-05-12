package log

import (
	"fmt"
	"strings"
)

func ColorWrap(content any, colors ...COLOR_ENUM) string {
	// 合并颜色
	colorNumbers := make([]string, 0, len(colors)+1)
	// 重置颜色
	colorNumbers = append(colorNumbers, "0")
	for _, color := range colors {
		number := strings.Replace(string(color), "\x1b[", "", -1)
		number = strings.Replace(number, "m", "", -1)
		colorNumbers = append(colorNumbers, number)
	}
	return fmt.Sprintf("\x1b[%sm%s%s", strings.Join(colorNumbers, ";"), fmt.Sprint(content), COLOR_CTRL_RESET)
}
