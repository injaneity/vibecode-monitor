// Package claude provides usage tracking for Claude Code by parsing local session files.
package claude

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// UsageData contains computed usage statistics.
type UsageData struct {
	// 5-hour cycle stats
	CyclePrompts   int
	CycleStartTime time.Time

	// Weekly stats
	WeeklySonnetHours float64
	WeeklyOpusHours   float64
	WeeklyPrompts     int
	WeeklyStartTime   time.Time

	// Reset times
	CycleResetIn  time.Duration
	WeeklyResetIn time.Duration

	// Tier info
	Tier     TierLimits
	TierName string

	// Metadata
	LastUpdated   time.Time
	SessionsCount int
}

// Tracker manages usage calculation.
type Tracker struct {
	tier     TierLimits
	tierName string
}

// NewTracker creates a tracker with the specified tier.
func NewTracker(tierName string) *Tracker {
	return &Tracker{
		tier:     GetTierLimits(tierName),
		tierName: tierName,
	}
}

// Calculate computes current usage statistics.
func (t *Tracker) Calculate() (*UsageData, error) {
	// Time boundaries
	now := time.Now()
	weekStart := getWeekStart(now)
	cycleStart := get5hCycleStart(now)

	// Find all sessions
	sessionPaths, err := FindAllSessions()
	if err != nil {
		return nil, fmt.Errorf("finding sessions: %w", err)
	}

	usage := &UsageData{
		CycleStartTime:  cycleStart,
		WeeklyStartTime: weekStart,
		Tier:            t.tier,
		TierName:        t.tierName,
		LastUpdated:     now,
	}

	// Parse and aggregate sessions
	for _, path := range sessionPaths {
		session, err := ParseJSONLFile(path)
		if err != nil || session == nil {
			continue
		}

		// Skip sessions with no real activity
		if session.DurationHours <= 0 && session.PromptCount == 0 {
			continue
		}

		usage.SessionsCount++

		// Check if session is in current 5h cycle
		if session.StartTime.After(cycleStart) || session.StartTime.Equal(cycleStart) {
			usage.CyclePrompts += session.PromptCount
		}

		// Check if session is in current week
		if session.StartTime.After(weekStart) || session.StartTime.Equal(weekStart) {
			usage.WeeklyPrompts += session.PromptCount

			// Calculate model-specific hours
			totalResponses := session.SonnetResponses + session.OpusResponses
			if totalResponses > 0 {
				sonnetRatio := float64(session.SonnetResponses) / float64(totalResponses)
				opusRatio := float64(session.OpusResponses) / float64(totalResponses)
				usage.WeeklySonnetHours += session.DurationHours * sonnetRatio
				usage.WeeklyOpusHours += session.DurationHours * opusRatio
			} else {
				// Default to Sonnet if no model info
				usage.WeeklySonnetHours += session.DurationHours
			}
		}
	}

	// Calculate reset times
	usage.CycleResetIn = cycleStart.Add(5 * time.Hour).Sub(now)
	if usage.CycleResetIn < 0 {
		usage.CycleResetIn = 0
	}

	nextMonday := weekStart.Add(7 * 24 * time.Hour)
	usage.WeeklyResetIn = nextMonday.Sub(now)
	if usage.WeeklyResetIn < 0 {
		usage.WeeklyResetIn = 0
	}

	return usage, nil
}

// getWeekStart returns Monday 00:00:00 of the current week.
func getWeekStart(now time.Time) time.Time {
	daysSinceMonday := int(now.Weekday()) - 1
	if daysSinceMonday < 0 {
		daysSinceMonday = 6 // Sunday
	}
	monday := now.AddDate(0, 0, -daysSinceMonday)
	return time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, now.Location())
}

// get5hCycleStart returns the start of the current 5-hour cycle.
func get5hCycleStart(now time.Time) time.Time {
	// Calculate hours since Unix epoch
	hoursSinceEpoch := float64(now.Unix()) / 3600
	cycleNumber := int(hoursSinceEpoch / 5)
	cycleStartUnix := int64(cycleNumber) * 5 * 3600
	return time.Unix(cycleStartUnix, 0)
}

// TotalWeeklyHours returns combined Sonnet + Opus hours.
func (u *UsageData) TotalWeeklyHours() float64 {
	return u.WeeklySonnetHours + u.WeeklyOpusHours
}

// WeeklyPercentage returns usage as percentage of weekly limit.
func (u *UsageData) WeeklyPercentage() float64 {
	total := u.TotalWeeklyHours()
	max := u.Tier.GetTotalWeeklyMax()
	if max <= 0 {
		return 0
	}
	return (total / max) * 100
}

// FormatResetTime formats a duration as "Xh Ym".
func FormatResetTime(d time.Duration) string {
	if d <= 0 {
		return "resetting..."
	}

	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60

	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}

// DetectTier attempts to read tier from Claude credentials file.
func DetectTier() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	credPath := filepath.Join(home, ".claude", ".credentials.json")
	data, err := os.ReadFile(credPath)
	if err != nil {
		return "", err
	}

	var creds struct {
		RateLimitTier string `json:"rate_limit_tier"`
		// Also check other possible field names
		Tier string `json:"tier"`
		Plan string `json:"plan"`
	}

	if err := json.Unmarshal(data, &creds); err != nil {
		return "", err
	}

	// Try different field names
	if creds.RateLimitTier != "" {
		return creds.RateLimitTier, nil
	}
	if creds.Tier != "" {
		return creds.Tier, nil
	}
	if creds.Plan != "" {
		return creds.Plan, nil
	}

	return "", fmt.Errorf("no tier found in credentials")
}
