package cli

import (
	"fmt"
	"time"
)

func FormatTimeSince(sinceStr string) string {
	if sinceStr == "" {
		return "Unknown"
	}

	// Parse the ISO 8601 date string
	since, err := time.Parse(time.RFC3339, sinceStr)
	if err != nil {
		return sinceStr // Return original if parsing fails
	}

	duration := time.Since(since)
	days := int(duration.Hours() / 24)

	if days < 1 {
		hours := int(duration.Hours())
		if hours < 1 {
			return "Less than an hour ago"
		}
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	} else if days < 30 {
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	} else if days < 365 {
		months := days / 30
		if months == 1 {
			return "1 month ago"
		}
		return fmt.Sprintf("%d months ago", months)
	} else {
		years := days / 365
		remainingDays := days % 365
		if remainingDays < 30 {
			if years == 1 {
				return "1 year ago"
			}
			return fmt.Sprintf("%d years ago", years)
		} else {
			months := remainingDays / 30
			yearText := "year"
			monthText := "month"

			if years > 1 {
				yearText = "years"
			}
			if months > 1 {
				monthText = "months"
			}

			return fmt.Sprintf("%d %s, %d %s ago", years, yearText, months, monthText)
		}
	}
}

func FormatTime(sinceStr string) string {
	if sinceStr == "" {
		return "Unknown"
	}

	// Parse the ISO 8601 date string
	since, err := time.Parse(time.RFC3339, sinceStr)
	if err != nil {
		return sinceStr // Return original if parsing fails
	}

	return since.Format("2006-01-02 15:04")
}
