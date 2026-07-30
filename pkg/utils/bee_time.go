package utils

import (
	"fmt"
	"time"
)

type BeeTime time.Time

func (ct BeeTime) MarshalJSON() ([]byte, error) {
	t := time.Time(ct)
	formatted := t.Format("2006-01-02 15:04:05")

	// 预分配容量：len(formatted) + 2 给双引号留空间
	buf := make([]byte, 0, len(formatted)+2)
	// 直接追加 `"格式化字符串"`
	buf = fmt.Appendf(buf, `"%s"`, formatted)
	return buf, nil
}
