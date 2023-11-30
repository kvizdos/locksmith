package fingerprint

import (
	"github.com/kvizdos/locksmith/observability"
)

var UseSafetyScore SafetyScore = REGULAR_SAFETY

type SafetyScore int

const (
	HIGH_SAFETY    SafetyScore = 90
	REGULAR_SAFETY SafetyScore = 80
	LOW_SAFETY     SafetyScore = 65
)

type Fingerprint struct {
	// Screen size + color depth
	// Doesn't change frequently.
	// High score
	Screen string `json:"screen"`
	// Current user Time Zone
	// Could periodically change.
	// Low score.
	Timezone string `json:"timezone"`
	// CPU devices. Does not change
	// frequently.
	// High score.
	HardwareConcurrency int `json:"hardwareConcurrency"`
	// Device memory. Does not change
	// frequently.
	// High Score
	DeviceMemory int `json:"deviceMemory"`
	// Canvas Fingerprint
	// Does not change frequently.
	// High Score
	Canvas string `json:"canvas"`
	// User-set language.
	// Low Score
	Language string `json:"lang"`
	// WebGL Fingerprint
	// High Score
	WebGL string `json:"webgl"`
	// Touch info
	// Low score
	Touch string `json:"touch"`
	// Does the device use a battery?
	// Doesn't change frequently.
	// Medium score
	Battery bool `json:"battery"`
	// Device platform
	// Medium score
	Platform string `json:"platform"`
	// User Agent
	// Changes frequently with updates.
	// Low Score
	UserAgent string `json:"userAgent"`
	// Window Size
	// Changes very frequently.
	// Low score
	WindowSize string `json:"windowSize"`
	// Do Not Track status
	// Can change frequently.
	// Low score
	DoNotTrack bool `json:"dnt"`
	// Connected Devices.
	// Does not change much.
	// Medium Score
	Devices string `json:"devices"`
	// Audio Fingerprint
	// Low Score
	Audio string `json:"audio"`

	// Non-User Defined; defined at Login
	// by the Server
	IPAddress string `json:"-"`
	// Quick Hash of Everything Together
	// If this matches 100%, its faster
	// sometimes than doing each one individually.
	QuickHash string `json:"-"`
}

func (f Fingerprint) ValidateRequest() bool {
	if f.Timezone == "" {
		return false
	}
	if len(f.Canvas) != 64 {
		return false
	}
	if len(f.WebGL) != 64 {
		return false
	}
	if len(f.Screen) != 64 {
		return false
	}
	if len(f.Touch) != 64 {
		return false
	}
	if len(f.UserAgent) != 64 {
		return false
	}
	if len(f.WindowSize) != 64 {
		return false
	}
	if len(f.Devices) != 64 {
		return false
	}
	if len(f.Audio) != 64 {
		return false
	}
	if f.Language == "" {
		return false
	}
	if f.Platform == "" {
		return false
	}

	// An extra layer of protection
	// against IP or QuickHash being set
	// by the User
	if f.QuickHash != "" {
		return false
	}
	if f.IPAddress != "" {
		return false
	}

	return true
}

type fingerprintScore int

const (
	ULTRA  fingerprintScore = 20
	HIGH   fingerprintScore = 15
	MEDIUM fingerprintScore = 10
	LOW    fingerprintScore = 5
)

func fingerprintScoreToString(score fingerprintScore) string {
	switch score {
	case ULTRA:
		return "ultra"
	case HIGH:
		return "high"
	case MEDIUM:
		return "medium"
	case LOW:
		return "low"
	default:
		return "low"
	}
}

var fingerprintScores = map[string]fingerprintScore{
	"canvas": ULTRA,
	"webGL":  ULTRA,

	"screen":              HIGH,
	"hardwareConcurrency": HIGH,
	"deviceMemory":        HIGH,
	"battery":             HIGH,

	"touch":    MEDIUM,
	"platform": MEDIUM,
	"devices":  MEDIUM,

	"timezone":   LOW,
	"language":   LOW,
	"userAgent":  LOW,
	"windowSize": LOW,
	"doNotTrack": LOW,
	"ipAddress":  LOW,
	"audio":      LOW,
}

// This returns a "trust value" of 0-1 (low-high)
func (f Fingerprint) ScoreAgainst(against Fingerprint) int {
	if f.QuickHash == against.QuickHash {
		return 1
	}

	score := 100

	// Ultra level stuff
	if f.Canvas != against.Canvas {
		givenScore := fingerprintScores["canvas"]
		score -= int(givenScore)
		observability.FingerprintHitCounter.WithLabelValues(fingerprintScoreToString(givenScore)).Inc()
	}

	if f.WebGL != against.WebGL {
		givenScore := fingerprintScores["webGL"]
		score -= int(givenScore)
		observability.FingerprintHitCounter.WithLabelValues(fingerprintScoreToString(givenScore)).Inc()
	}

	// High Level stuff
	if f.Screen != against.Screen {
		givenScore := fingerprintScores["screen"]
		score -= int(givenScore)
		observability.FingerprintHitCounter.WithLabelValues(fingerprintScoreToString(givenScore)).Inc()
	}

	if f.HardwareConcurrency != against.HardwareConcurrency {
		givenScore := fingerprintScores["hardwareConcurrency"]
		score -= int(givenScore)
		observability.FingerprintHitCounter.WithLabelValues(fingerprintScoreToString(givenScore)).Inc()
	}

	if f.DeviceMemory != against.DeviceMemory {
		givenScore := fingerprintScores["deviceMemory"]
		score -= int(givenScore)
		observability.FingerprintHitCounter.WithLabelValues(fingerprintScoreToString(givenScore)).Inc()
	}

	// Medium Scoring
	if f.Touch != against.Touch {
		givenScore := fingerprintScores["touch"]
		score -= int(givenScore)
		observability.FingerprintHitCounter.WithLabelValues(fingerprintScoreToString(givenScore)).Inc()
	}

	if f.Battery != against.Battery {
		givenScore := fingerprintScores["battery"]
		score -= int(givenScore)
		observability.FingerprintHitCounter.WithLabelValues(fingerprintScoreToString(givenScore)).Inc()
	}

	if f.Platform != against.Platform {
		givenScore := fingerprintScores["platform"]
		score -= int(givenScore)
		observability.FingerprintHitCounter.WithLabelValues(fingerprintScoreToString(givenScore)).Inc()
	}

	// Low Scoring things
	if f.Timezone != against.Timezone {
		givenScore := fingerprintScores["timezone"]
		score -= int(givenScore)
		observability.FingerprintHitCounter.WithLabelValues(fingerprintScoreToString(givenScore)).Inc()
	}

	if f.Language != against.Language {
		givenScore := fingerprintScores["language"]
		score -= int(givenScore)
		observability.FingerprintHitCounter.WithLabelValues(fingerprintScoreToString(givenScore)).Inc()
	}

	if f.UserAgent != against.UserAgent {
		givenScore := fingerprintScores["userAgent"]
		score -= int(givenScore)
		observability.FingerprintHitCounter.WithLabelValues(fingerprintScoreToString(givenScore)).Inc()
	}

	if f.WindowSize != against.WindowSize {
		givenScore := fingerprintScores["windowSize"]
		score -= int(givenScore)
		observability.FingerprintHitCounter.WithLabelValues(fingerprintScoreToString(givenScore)).Inc()
	}

	if f.DoNotTrack != against.DoNotTrack {
		givenScore := fingerprintScores["doNotTrack"]
		score -= int(givenScore)
		observability.FingerprintHitCounter.WithLabelValues(fingerprintScoreToString(givenScore)).Inc()
	}

	if f.IPAddress != against.IPAddress {
		givenScore := fingerprintScores["ipAddress"]
		score -= int(givenScore)
		observability.FingerprintHitCounter.WithLabelValues(fingerprintScoreToString(givenScore)).Inc()
	}

	if f.Audio != against.Audio {
		givenScore := fingerprintScores["audio"]
		score -= int(givenScore)
		observability.FingerprintHitCounter.WithLabelValues(fingerprintScoreToString(givenScore)).Inc()
	}

	observability.FingerprintScore.Add(float64(score))
	observability.FingerprintEvaluations.Inc()

	return score
}

func IsPassingScore(checkScore int) bool {
	return checkScore >= int(UseSafetyScore)
}
