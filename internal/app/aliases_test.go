package app

import "testing"

func TestNormalizeTargetAliases(t *testing.T) {
	t.Parallel()

	cases := map[string]string{
		"iftar":   "maghrib",
		"akşam":   "maghrib",
		"öğle":    "dhuhr",
		"güneş":   "sunrise",
		"yatsı":   "isha",
		"imsak":   "imsak",
		"maghrib": "maghrib",
	}

	for input, want := range cases {
		got, ok := NormalizeTarget(input)
		if !ok {
			t.Fatalf("expected alias %q to resolve", input)
		}
		if got != want {
			t.Fatalf("NormalizeTarget(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestNormalizeTimesFieldAliases(t *testing.T) {
	t.Parallel()

	cases := map[string]string{
		"iftar":          "maghrib_at",
		"öğle":           "dhuhr_at",
		"yatsı":          "isha_at",
		"ramadan_active": "ramadan_active",
		"ramadan-active": "ramadan_active",
	}

	for input, want := range cases {
		got, ok := NormalizeTimesField(input)
		if !ok {
			t.Fatalf("expected field alias %q to resolve", input)
		}
		if got != want {
			t.Fatalf("NormalizeTimesField(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestNormalizeCountdownFieldAliases(t *testing.T) {
	t.Parallel()

	cases := map[string]string{
		"seconds":           "seconds_remaining",
		"seconds_remaining": "seconds_remaining",
		"minutes-remaining": "minutes_remaining",
		"target-at":         "target_at",
		"maghrib":           "maghrib_at",
	}

	for input, want := range cases {
		got, ok := NormalizeCountdownField(input)
		if !ok {
			t.Fatalf("expected field alias %q to resolve", input)
		}
		if got != want {
			t.Fatalf("NormalizeCountdownField(%q) = %q, want %q", input, got, want)
		}
	}
}
