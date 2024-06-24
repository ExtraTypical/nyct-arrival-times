package format

import (
	"fmt"
	"time"
)

func Duration(d time.Duration) string {
	minutes := int(d.Minutes())
	if minutes == 0 {
		return "now"
	}
	return fmt.Sprintf("%d min", minutes)
}
