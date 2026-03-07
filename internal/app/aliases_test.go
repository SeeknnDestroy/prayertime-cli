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

func TestNormalizeFieldAliases(t *testing.T) {
	t.Parallel()

	cases := map[string]string{
		"iftar": "maghrib_at",
		"öğle":  "dhuhr_at",
		"yatsı": "isha_at",
	}

	for input, want := range cases {
		got, ok := NormalizeField(input)
		if !ok {
			t.Fatalf("expected field alias %q to resolve", input)
		}
		if got != want {
			t.Fatalf("NormalizeField(%q) = %q, want %q", input, got, want)
		}
	}
}
