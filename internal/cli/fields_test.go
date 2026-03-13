package cli

import (
	"reflect"
	"testing"
)

func TestNormalizeTimesFieldAliases(t *testing.T) {
	t.Parallel()

	cases := map[string]string{
		"iftar":          "maghrib_at",
		"öğle":           "dhuhr_at",
		"yatsı":          "isha_at",
		"ramadan_active": "ramadan_active",
		"ramadan-active": "ramadan_active",
		"method-id":      "method_id",
	}

	for input, want := range cases {
		got, ok := normalizeTimesField(input)
		if !ok {
			t.Fatalf("expected field alias %q to resolve", input)
		}
		if got.Key != want {
			t.Fatalf("normalizeTimesField(%q) = %q, want %q", input, got.Key, want)
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
		got, ok := normalizeCountdownField(input)
		if !ok {
			t.Fatalf("expected field alias %q to resolve", input)
		}
		if got.Key != want {
			t.Fatalf("normalizeCountdownField(%q) = %q, want %q", input, got.Key, want)
		}
	}
}

func TestValidTimesFieldsPreserveCanonicalOrder(t *testing.T) {
	t.Parallel()

	want := []string{
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

	if got := validTimesFields(); !reflect.DeepEqual(got, want) {
		t.Fatalf("validTimesFields() = %#v, want %#v", got, want)
	}
}

func TestValidCountdownFieldsPreserveCanonicalOrder(t *testing.T) {
	t.Parallel()

	want := []string{
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

	if got := validCountdownFields(); !reflect.DeepEqual(got, want) {
		t.Fatalf("validCountdownFields() = %#v, want %#v", got, want)
	}
}
