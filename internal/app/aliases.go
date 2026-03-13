package app

import "strings"

var turkishNormalizer = strings.NewReplacer(
	"ç", "c",
	"ğ", "g",
	"ı", "i",
	"İ", "i",
	"ö", "o",
	"ş", "s",
	"ü", "u",
)

func normalizeToken(input string) string {
	value := strings.TrimSpace(strings.ToLower(input))
	value = turkishNormalizer.Replace(value)
	value = strings.ReplaceAll(value, "_", "-")
	value = strings.ReplaceAll(value, " ", "-")
	return value
}

var targetAliases = map[string]string{
	"next-prayer": "next-prayer",
	"imsak":       "imsak",
	"fajr":        "fajr",
	"sunrise":     "sunrise",
	"gunes":       "sunrise",
	"dhuhr":       "dhuhr",
	"ogle":        "dhuhr",
	"asr":         "asr",
	"ikindi":      "asr",
	"maghrib":     "maghrib",
	"iftar":       "maghrib",
	"aksam":       "maghrib",
	"sunset":      "sunset",
	"isha":        "isha",
	"yatsi":       "isha",
}

var canonicalTargets = []string{
	"next-prayer",
	"imsak",
	"fajr",
	"sunrise",
	"dhuhr",
	"asr",
	"maghrib",
	"sunset",
	"isha",
}

func NormalizeTarget(input string) (string, bool) {
	value, ok := targetAliases[normalizeToken(input)]
	return value, ok
}

func ValidTargets() []string {
	return append([]string(nil), canonicalTargets...)
}
