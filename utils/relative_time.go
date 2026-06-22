package utils

import (
	"fmt"
	"time"
)

func FormatRelativeDuration(now, target time.Time) string {
	totalSeconds := int64(now.Sub(target) / time.Second)
	if totalSeconds == 0 {
		return "0s"
	}

	isPast := totalSeconds > 0
	if totalSeconds < 0 {
		totalSeconds = -totalSeconds
	}

	days := totalSeconds / (24 * 60 * 60)
	totalSeconds %= 24 * 60 * 60
	hours := totalSeconds / (60 * 60)
	totalSeconds %= 60 * 60
	minutes := totalSeconds / 60
	seconds := totalSeconds % 60

	switch {
	case days > 0:
		return formatRelativeTimeUnit(isPast, "%dd", days)
	case hours > 0:
		return formatRelativeTimeUnit(isPast, "%dh", hours)
	case minutes > 0:
		return formatRelativeTimeUnit(isPast, "%dm", minutes)
	default:
		return formatRelativeTimeUnit(isPast, "%ds", seconds)
	}
}

func formatRelativeTimeUnit(isPast bool, unitFormat string, value int64) string {
	unit := fmt.Sprintf(unitFormat, value)
	if isPast {
		return unit + " ago"
	}

	return "in " + unit
}
