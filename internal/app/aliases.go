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

var fieldAliases = map[string]string{
	"location-name": "location_name",
	"location_name": "location_name",
	"latitude":      "latitude",
	"longitude":     "longitude",
	"timezone":      "timezone",
	"date":          "date",
	"imsak":         "imsak_at",
	"imsak-at":      "imsak_at",
	"imsak_at":      "imsak_at",
	"fajr":          "fajr_at",
	"fajr-at":       "fajr_at",
	"fajr_at":       "fajr_at",
	"sunrise":       "sunrise_at",
	"gunes":         "sunrise_at",
	"sunrise-at":    "sunrise_at",
	"sunrise_at":    "sunrise_at",
	"dhuhr":         "dhuhr_at",
	"ogle":          "dhuhr_at",
	"dhuhr-at":      "dhuhr_at",
	"dhuhr_at":      "dhuhr_at",
	"asr":           "asr_at",
	"ikindi":        "asr_at",
	"asr-at":        "asr_at",
	"asr_at":        "asr_at",
	"maghrib":       "maghrib_at",
	"iftar":         "maghrib_at",
	"aksam":         "maghrib_at",
	"maghrib-at":    "maghrib_at",
	"maghrib_at":    "maghrib_at",
	"sunset":        "sunset_at",
	"sunset-at":     "sunset_at",
	"sunset_at":     "sunset_at",
	"isha":          "isha_at",
	"yatsi":         "isha_at",
	"isha-at":       "isha_at",
	"isha_at":       "isha_at",
	"method-id":     "method_id",
	"method_id":     "method_id",
	"method-name":   "method_name",
	"method_name":   "method_name",
	"source":        "source",
	"ramadan":       "ramadan_active",
}

func NormalizeTarget(input string) (string, bool) {
	value, ok := targetAliases[normalizeToken(input)]
	return value, ok
}

func NormalizeField(input string) (string, bool) {
	value, ok := fieldAliases[normalizeToken(input)]
	return value, ok
}
