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
