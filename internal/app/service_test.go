package app

import (
	"context"
	"testing"
	"time"
)

type fakeResolver struct {
	results []Location
}

func (f fakeResolver) Search(ctx context.Context, query, countryCode string, limit int) ([]Location, error) {
	return f.results, nil
}

type fakeProvider struct {
	schedules map[string]DaySchedule
}

func (f fakeProvider) GetByCoordinates(ctx context.Context, latitude, longitude float64, date time.Time) (DaySchedule, error) {
	return f.schedules[date.Format("2006-01-02")], nil
}

type fixedClock struct {
	now time.Time
}

func (f fixedClock) Now() time.Time {
	return f.now
}

type spyResolver struct {
	results     []Location
	searchCalls int
}

func (s *spyResolver) Search(ctx context.Context, query, countryCode string, limit int) ([]Location, error) {
	s.searchCalls++
	return s.results, nil
}

type spyProvider struct {
	schedules map[string]DaySchedule
	calls     int
}

func (s *spyProvider) GetByCoordinates(ctx context.Context, latitude, longitude float64, date time.Time) (DaySchedule, error) {
	s.calls++
	if s.schedules == nil {
		return DaySchedule{}, nil
	}

	return s.schedules[date.Format("2006-01-02")], nil
}

func TestServiceGetCountdownUsesNextDayAfterTargetPasses(t *testing.T) {
	t.Parallel()

	tz := mustLocation(t, "Europe/Istanbul")
	resolver := fakeResolver{
		results: []Location{{
			Name:        "Istanbul",
			Country:     "Türkiye",
			CountryCode: "TR",
			Admin1:      "Istanbul",
			Latitude:    41.01384,
			Longitude:   28.94966,
			Timezone:    "Europe/Istanbul",
		}},
	}
	provider := fakeProvider{
		schedules: map[string]DaySchedule{
			"2026-03-07": scheduleForDay(tz, 2026, 3, 7, 19, 9),
			"2026-03-08": scheduleForDay(tz, 2026, 3, 8, 19, 10),
		},
	}
	clock := fixedClock{now: time.Date(2026, 3, 7, 20, 0, 0, 0, tz)}
	service := NewService(resolver, provider, clock)

	response, err := service.GetCountdown(context.Background(), CountdownRequest{
		Query:  "Istanbul",
		Target: "iftar",
	})
	if err != nil {
		t.Fatalf("GetCountdown returned error: %v", err)
	}

	if response.Target != "maghrib" {
		t.Fatalf("Target = %q, want maghrib", response.Target)
	}
	if response.Date != "2026-03-08" {
		t.Fatalf("Date = %q, want 2026-03-08", response.Date)
	}

	wantSeconds := int64((23*time.Hour + 10*time.Minute).Seconds())
	if response.SecondsRemaining != wantSeconds {
		t.Fatalf("SecondsRemaining = %d, want %d", response.SecondsRemaining, wantSeconds)
	}
}

func TestServiceGetCountdownNextPrayerUsesTomorrowImsakAfterIsha(t *testing.T) {
	t.Parallel()

	tz := mustLocation(t, "Europe/Istanbul")
	service := NewService(
		fakeResolver{
			results: []Location{{
				Name:        "Istanbul",
				Country:     "Türkiye",
				CountryCode: "TR",
				Admin1:      "Istanbul",
				Latitude:    41.01384,
				Longitude:   28.94966,
				Timezone:    "Europe/Istanbul",
			}},
		},
		fakeProvider{
			schedules: map[string]DaySchedule{
				"2026-03-07": scheduleForDay(tz, 2026, 3, 7, 19, 9),
				"2026-03-08": scheduleForDay(tz, 2026, 3, 8, 19, 10),
			},
		},
		fixedClock{now: time.Date(2026, 3, 7, 21, 0, 0, 0, tz)},
	)

	response, err := service.GetCountdown(context.Background(), CountdownRequest{
		Query:  "Istanbul",
		Target: "next-prayer",
	})
	if err != nil {
		t.Fatalf("GetCountdown returned error: %v", err)
	}

	if response.Target != "imsak" {
		t.Fatalf("Target = %q, want imsak", response.Target)
	}
	if response.Date != "2026-03-08" {
		t.Fatalf("Date = %q, want 2026-03-08", response.Date)
	}
	if response.TargetAt != "2026-03-08T05:46:00+03:00" {
		t.Fatalf("TargetAt = %q, want 2026-03-08T05:46:00+03:00", response.TargetAt)
	}
}

func TestServiceGetCountdownUsesScheduleTimezoneForAt(t *testing.T) {
	t.Parallel()

	tz := mustLocation(t, "Europe/Istanbul")
	service := NewService(
		fakeResolver{
			results: []Location{{
				Name:        "Istanbul",
				Country:     "Türkiye",
				CountryCode: "TR",
				Admin1:      "Istanbul",
				Latitude:    41.01384,
				Longitude:   28.94966,
				Timezone:    "Europe/Istanbul",
			}},
		},
		fakeProvider{
			schedules: map[string]DaySchedule{
				"2026-03-08": scheduleForDay(tz, 2026, 3, 8, 19, 10),
			},
		},
		fixedClock{now: time.Date(2026, 3, 7, 18, 0, 0, 0, time.UTC)},
	)

	at := time.Date(2026, 3, 7, 22, 30, 0, 0, time.UTC)
	response, err := service.GetCountdown(context.Background(), CountdownRequest{
		Query:  "Istanbul",
		Target: "next-prayer",
		At:     &at,
	})
	if err != nil {
		t.Fatalf("GetCountdown returned error: %v", err)
	}

	if response.Date != "2026-03-08" {
		t.Fatalf("Date = %q, want 2026-03-08", response.Date)
	}
	if response.Target != "imsak" {
		t.Fatalf("Target = %q, want imsak", response.Target)
	}
}

func TestServiceGetCountdownReturnsInternalErrorWhenTargetCannotBeResolved(t *testing.T) {
	t.Parallel()

	tz := mustLocation(t, "Europe/Istanbul")
	service := NewService(
		fakeResolver{
			results: []Location{{
				Name:        "Istanbul",
				Country:     "Türkiye",
				CountryCode: "TR",
				Admin1:      "Istanbul",
				Latitude:    41.01384,
				Longitude:   28.94966,
				Timezone:    "Europe/Istanbul",
			}},
		},
		fakeProvider{
			schedules: map[string]DaySchedule{
				"2026-03-07": {
					Latitude:  41.01384,
					Longitude: 28.94966,
					Timezone:  tz.String(),
					Date:      time.Date(2026, 3, 7, 0, 0, 0, 0, tz),
				},
				"2026-03-08": {
					Latitude:  41.01384,
					Longitude: 28.94966,
					Timezone:  tz.String(),
					Date:      time.Date(2026, 3, 8, 0, 0, 0, 0, tz),
				},
			},
		},
		fixedClock{now: time.Date(2026, 3, 7, 21, 0, 0, 0, tz)},
	)

	_, err := service.GetCountdown(context.Background(), CountdownRequest{
		Query:  "Istanbul",
		Target: "next-prayer",
	})
	if err == nil {
		t.Fatal("expected internal error")
	}

	cliErr := AsCLIError(err)
	if cliErr.ExitCode != ExitFailure {
		t.Fatalf("ExitCode = %d, want %d", cliErr.ExitCode, ExitFailure)
	}
	if cliErr.Message != "failed to resolve countdown target from schedule" {
		t.Fatalf("Message = %q, want failed target resolution message", cliErr.Message)
	}
}

func TestServiceGetTimesRejectsPartialCoordinates(t *testing.T) {
	t.Parallel()

	service := NewService(fakeResolver{}, fakeProvider{}, fixedClock{})
	latitude := 41.01384

	_, err := service.GetTimes(context.Background(), TimesRequest{Latitude: &latitude})
	if err == nil {
		t.Fatal("expected usage error")
	}

	cliErr := AsCLIError(err)
	if cliErr.ExitCode != ExitUsage {
		t.Fatalf("ExitCode = %d, want %d", cliErr.ExitCode, ExitUsage)
	}
}

func TestServiceGetTimesReturnsAmbiguousLocation(t *testing.T) {
	t.Parallel()

	resolver := fakeResolver{
		results: []Location{
			{Name: "Springfield", Country: "United States", Admin1: "Illinois", Latitude: 39.78, Longitude: -89.64, Timezone: "America/Chicago"},
			{Name: "Springfield", Country: "United States", Admin1: "Missouri", Latitude: 37.20, Longitude: -93.29, Timezone: "America/Chicago"},
			{Name: "Springfield", Country: "United States", Admin1: "Missouri", Latitude: 37.20, Longitude: -93.29, Timezone: "America/Chicago"},
		},
	}
	service := NewService(resolver, fakeProvider{}, fixedClock{})

	_, err := service.GetTimes(context.Background(), TimesRequest{Query: "Springfield"})
	if err == nil {
		t.Fatal("expected ambiguous location error")
	}

	cliErr := AsCLIError(err)
	if cliErr.ExitCode != ExitNotFound {
		t.Fatalf("ExitCode = %d, want %d", cliErr.ExitCode, ExitNotFound)
	}
	if cliErr.ErrorType != "ambiguous_location" {
		t.Fatalf("ErrorType = %q, want ambiguous_location", cliErr.ErrorType)
	}
	if cliErr.Details == nil {
		t.Fatal("Details = nil, want candidates")
	}
	if len(cliErr.Details.Candidates) != 2 {
		t.Fatalf("len(Details.Candidates) = %d, want 2", len(cliErr.Details.Candidates))
	}
}

func TestServiceGetCountdownRejectsInvalidTargetBeforeLookup(t *testing.T) {
	t.Parallel()

	resolver := &spyResolver{
		results: []Location{{
			Name:        "Istanbul",
			Country:     "Türkiye",
			CountryCode: "TR",
			Admin1:      "Istanbul",
			Latitude:    41.01384,
			Longitude:   28.94966,
			Timezone:    "Europe/Istanbul",
		}},
	}
	provider := &spyProvider{}
	service := NewService(resolver, provider, fixedClock{})

	_, err := service.GetCountdown(context.Background(), CountdownRequest{
		Query:  "Istanbul",
		Target: "bogus",
	})
	if err == nil {
		t.Fatal("expected usage error")
	}

	cliErr := AsCLIError(err)
	if cliErr.ExitCode != ExitUsage {
		t.Fatalf("ExitCode = %d, want %d", cliErr.ExitCode, ExitUsage)
	}
	if cliErr.Details == nil || len(cliErr.Details.ValidTargets) == 0 {
		t.Fatal("expected valid target details")
	}
	if resolver.searchCalls != 0 {
		t.Fatalf("resolver.searchCalls = %d, want 0", resolver.searchCalls)
	}
	if provider.calls != 0 {
		t.Fatalf("provider.calls = %d, want 0", provider.calls)
	}
}

func scheduleForDay(tz *time.Location, year int, month time.Month, day int, maghribHour int, maghribMinute int) DaySchedule {
	return DaySchedule{
		Latitude:      41.01384,
		Longitude:     28.94966,
		Timezone:      tz.String(),
		Date:          time.Date(year, month, day, 0, 0, 0, 0, tz),
		ImsakAt:       time.Date(year, month, day, 5, 46, 0, 0, tz),
		FajrAt:        time.Date(year, month, day, 5, 56, 0, 0, tz),
		SunriseAt:     time.Date(year, month, day, 7, 21, 0, 0, tz),
		DhuhrAt:       time.Date(year, month, day, 13, 20, 0, 0, tz),
		AsrAt:         time.Date(year, month, day, 16, 31, 0, 0, tz),
		MaghribAt:     time.Date(year, month, day, maghribHour, maghribMinute, 0, 0, tz),
		SunsetAt:      time.Date(year, month, day, maghribHour, maghribMinute, 0, 0, tz),
		IshaAt:        time.Date(year, month, day, 20, 28, 0, 0, tz),
		MethodID:      13,
		MethodName:    "Diyanet İşleri Başkanlığı, Turkey (experimental)",
		Source:        "aladhan:method=13",
		RamadanActive: true,
	}
}

func mustLocation(t *testing.T, name string) *time.Location {
	t.Helper()

	location, err := time.LoadLocation(name)
	if err != nil {
		t.Fatalf("LoadLocation(%q): %v", name, err)
	}

	return location
}
