package fingerprint

import "testing"

func TestFingerprintHitsOnThreeLowScores(t *testing.T) {
	againstFingerprint := Fingerprint{
		Screen:              "483d4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e",
		Timezone:            "America/California",
		HardwareConcurrency: 8,
		DeviceMemory:        4,
		Canvas:              "483d4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e",
		Language:            "Es_US",
		WebGL:               "483d4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e",
		Touch:               "483d4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e",
		Battery:             false,
		Platform:            "483d4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e",
		UserAgent:           "483d4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e",
		WindowSize:          "483d4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e",
		DoNotTrack:          true,
		Devices:             "483d4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e",
		IPAddress:           "192.158.24.21",
		QuickHash:           "483d4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e",
	}

	newFingerprint := Fingerprint{
		Screen:              "483d4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e",
		Timezone:            "America/New_York",
		HardwareConcurrency: 8,
		DeviceMemory:        4,
		Canvas:              "483d4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e",
		Language:            "En_US",
		WebGL:               "483d4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e",
		Touch:               "483d4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e",
		Battery:             false,
		Platform:            "483d4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e",
		UserAgent:           "483d4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e",
		WindowSize:          "483d4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e",
		DoNotTrack:          false,
		Devices:             "483d4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e",
		IPAddress:           "192.158.24.21",
		QuickHash:           "diff4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e",
	}

	score := newFingerprint.ScoreAgainst(againstFingerprint)

	if score != 85 {
		t.Errorf("expected 85, got %d", score)
	}
}

func TestFingerprintHitsOnThreeLowAndOneMediumScores(t *testing.T) {
	againstFingerprint := Fingerprint{
		Screen:              "483d4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e",
		Timezone:            "America/California",
		HardwareConcurrency: 8,
		DeviceMemory:        4,
		Canvas:              "483d4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e",
		Language:            "Es_US",
		WebGL:               "483d4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e",
		Touch:               "483d4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e",
		Battery:             false,
		Platform:            "483d4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e",
		UserAgent:           "483d4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e",
		WindowSize:          "483d4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e",
		DoNotTrack:          true,
		Devices:             "483d4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e",
		IPAddress:           "192.158.24.21",
		QuickHash:           "483d4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e",
	}

	newFingerprint := Fingerprint{
		Screen:              "483d4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e",
		Timezone:            "America/New_York",
		HardwareConcurrency: 8,
		DeviceMemory:        4,
		Canvas:              "483d4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e",
		Language:            "En_US",
		WebGL:               "483d4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e",
		Touch:               "483d4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e",
		Battery:             false,
		Platform:            "MacOS",
		UserAgent:           "483d4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e",
		WindowSize:          "483d4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e",
		DoNotTrack:          false,
		Devices:             "483d4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e",
		IPAddress:           "192.158.24.21",
		QuickHash:           "diff4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e",
	}

	score := newFingerprint.ScoreAgainst(againstFingerprint)

	if score != 75 {
		t.Errorf("expected 75, got %d", score)
	}
}

func TestFingerprintHitsOnOneUltraScore(t *testing.T) {
	againstFingerprint := Fingerprint{
		Screen:              "483d4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e",
		Timezone:            "America/California",
		HardwareConcurrency: 8,
		DeviceMemory:        4,
		Canvas:              "483d4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e",
		Language:            "Es_US",
		WebGL:               "483d4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e",
		Touch:               "483d4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e",
		Battery:             false,
		Platform:            "483d4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e",
		UserAgent:           "483d4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e",
		WindowSize:          "483d4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e",
		DoNotTrack:          true,
		Devices:             "483d4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e",
		IPAddress:           "192.158.24.21",
		QuickHash:           "z83d4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e",
	}

	newFingerprint := Fingerprint{
		Screen:              "483d4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e",
		Timezone:            "America/California",
		HardwareConcurrency: 8,
		DeviceMemory:        4,
		Canvas:              "483d4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e",
		Language:            "Es_US",
		WebGL:               "diff4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e",
		Touch:               "483d4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e",
		Battery:             false,
		Platform:            "483d4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e",
		UserAgent:           "483d4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e",
		WindowSize:          "483d4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e",
		DoNotTrack:          true,
		Devices:             "483d4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e",
		IPAddress:           "192.158.24.21",
		QuickHash:           "483d4549c44f24e47049254bf3c60b923f193830294fda89507f2b17f892da5e",
	}

	score := newFingerprint.ScoreAgainst(againstFingerprint)

	if score != 80 {
		t.Errorf("expected 80, got %d", score)
	}
}

func TestIsPassingScore(t *testing.T) {
	UseSafetyScore = REGULAR_SAFETY

	shouldFail := IsPassingScore(75)

	if shouldFail {
		t.Error("expected score of 75 to fail!")
	}

	shouldPass := IsPassingScore(85)

	if !shouldPass {
		t.Error("expected score of 85 to pass!")
	}
}
