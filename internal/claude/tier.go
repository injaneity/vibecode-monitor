// Package claude provides usage tracking for Claude Code by parsing local session files.
package claude

// TierLimits defines the usage limits for a subscription tier.
type TierLimits struct {
	Tier            string  // Tier identifier
	Cycle5hMin      int     // Minimum prompts per 5h cycle
	Cycle5hMax      int     // Maximum prompts per 5h cycle
	WeeklySonnetMin float64 // Minimum weekly Sonnet hours
	WeeklySonnetMax float64 // Maximum weekly Sonnet hours
	WeeklyOpusMin   float64 // Minimum weekly Opus hours (0 if not available)
	WeeklyOpusMax   float64 // Maximum weekly Opus hours (0 if not available)
}

// Predefined tier limits based on Claude's actual limits.
var Tiers = map[string]TierLimits{
	"free": {
		Tier:            "free",
		Cycle5hMin:      10,
		Cycle5hMax:      40,
		WeeklySonnetMin: 40,
		WeeklySonnetMax: 80,
	},
	"pro": {
		Tier:            "pro",
		Cycle5hMin:      10,
		Cycle5hMax:      40,
		WeeklySonnetMin: 40,
		WeeklySonnetMax: 80,
	},
	"max_5x": {
		Tier:            "max_5x",
		Cycle5hMin:      50,
		Cycle5hMax:      200,
		WeeklySonnetMin: 140,
		WeeklySonnetMax: 280,
		WeeklyOpusMin:   15,
		WeeklyOpusMax:   35,
	},
	"max_20x": {
		Tier:            "max_20x",
		Cycle5hMin:      200,
		Cycle5hMax:      800,
		WeeklySonnetMin: 240,
		WeeklySonnetMax: 480,
		WeeklyOpusMin:   24,
		WeeklyOpusMax:   40,
	},
}

// TierFromRateLimitTier maps OAuth rate_limit_tier values to our tier names.
var TierFromRateLimitTier = map[string]string{
	"free":       "free",
	"pro":        "pro",
	"max5":       "max_5x",
	"max5x":      "max_5x",
	"max_5x":     "max_5x",
	"max20":      "max_20x",
	"max20x":     "max_20x",
	"max_20x":    "max_20x",
	"team":       "max_5x",  // Team defaults to max_5x level
	"enterprise": "max_20x", // Enterprise defaults to max_20x level
}

// GetTierLimits returns the limits for a given tier, defaulting to "pro" if unknown.
func GetTierLimits(tier string) TierLimits {
	if limits, ok := Tiers[tier]; ok {
		return limits
	}
	// Check if it's a rate_limit_tier mapping
	if mappedTier, ok := TierFromRateLimitTier[tier]; ok {
		if limits, ok := Tiers[mappedTier]; ok {
			return limits
		}
	}
	return Tiers["pro"]
}

// HasOpus returns true if the tier includes Opus access.
func (t TierLimits) HasOpus() bool {
	return t.WeeklyOpusMax > 0
}

// GetTotalWeeklyMax returns the total weekly hours limit (Sonnet + Opus).
func (t TierLimits) GetTotalWeeklyMax() float64 {
	return t.WeeklySonnetMax + t.WeeklyOpusMax
}
