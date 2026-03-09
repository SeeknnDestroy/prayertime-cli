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

var timesFieldAliases = map[string]string{
	"location-name":  "location_name",
	"location_name":  "location_name",
	"latitude":       "latitude",
	"longitude":      "longitude",
	"timezone":       "timezone",
	"date":           "date",
	"imsak":          "imsak_at",
	"imsak-at":       "imsak_at",
	"imsak_at":       "imsak_at",
	"fajr":           "fajr_at",
	"fajr-at":        "fajr_at",
	"fajr_at":        "fajr_at",
	"sunrise":        "sunrise_at",
	"gunes":          "sunrise_at",
	"sunrise-at":     "sunrise_at",
	"sunrise_at":     "sunrise_at",
	"dhuhr":          "dhuhr_at",
	"ogle":           "dhuhr_at",
	"dhuhr-at":       "dhuhr_at",
	"dhuhr_at":       "dhuhr_at",
	"asr":            "asr_at",
	"ikindi":         "asr_at",
	"asr-at":         "asr_at",
	"asr_at":         "asr_at",
	"maghrib":        "maghrib_at",
	"iftar":          "maghrib_at",
	"aksam":          "maghrib_at",
	"maghrib-at":     "maghrib_at",
	"maghrib_at":     "maghrib_at",
	"sunset":         "sunset_at",
	"sunset-at":      "sunset_at",
	"sunset_at":      "sunset_at",
	"isha":           "isha_at",
	"yatsi":          "isha_at",
	"isha-at":        "isha_at",
	"isha_at":        "isha_at",
	"method-id":      "method_id",
	"method_id":      "method_id",
	"method-name":    "method_name",
	"method_name":    "method_name",
	"source":         "source",
	"ramadan":        "ramadan_active",
	"ramadan-active": "ramadan_active",
}

var canonicalTimesFields = []string{
	"location_name",
	"latitude",
	"longitude",
	"timezone",
	"date",
	"imsak_at",
	"fajr_at",
	"sunrise_at",
	"dhuhr_at",
	"asr_at",
	"maghrib_at",
	"sunset_at",
	"isha_at",
	"method_id",
	"method_name",
	"source",
	"ramadan_active",
}

var countdownFieldAliases = cloneAliases(timesFieldAliases)

var canonicalCountdownFields = []string{
	"location_name",
	"latitude",
	"longitude",
	"timezone",
	"date",
	"imsak_at",
	"fajr_at",
	"sunrise_at",
	"dhuhr_at",
	"asr_at",
	"maghrib_at",
	"sunset_at",
	"isha_at",
	"method_id",
	"method_name",
	"source",
	"ramadan_active",
	"target",
	"target_at",
	"seconds_remaining",
	"minutes_remaining",
}

func init() {
	countdownFieldAliases["target"] = "target"
	countdownFieldAliases["target-at"] = "target_at"
	countdownFieldAliases["target_at"] = "target_at"
	countdownFieldAliases["seconds"] = "seconds_remaining"
	countdownFieldAliases["seconds-remaining"] = "seconds_remaining"
	countdownFieldAliases["seconds_remaining"] = "seconds_remaining"
	countdownFieldAliases["minutes"] = "minutes_remaining"
	countdownFieldAliases["minutes-remaining"] = "minutes_remaining"
	countdownFieldAliases["minutes_remaining"] = "minutes_remaining"
}

func NormalizeTarget(input string) (string, bool) {
	value, ok := targetAliases[normalizeToken(input)]
	return value, ok
}

func NormalizeTimesField(input string) (string, bool) {
	value, ok := timesFieldAliases[normalizeToken(input)]
	return value, ok
}

func NormalizeCountdownField(input string) (string, bool) {
	value, ok := countdownFieldAliases[normalizeToken(input)]
	return value, ok
}

func ValidTargets() []string {
	return append([]string(nil), canonicalTargets...)
}

func ValidTimesFields() []string {
	return append([]string(nil), canonicalTimesFields...)
}

func ValidCountdownFields() []string {
	return append([]string(nil), canonicalCountdownFields...)
}

func cloneAliases(input map[string]string) map[string]string {
	cloned := make(map[string]string, len(input))
	for key, value := range input {
		cloned[key] = value
	}
	return cloned
}
