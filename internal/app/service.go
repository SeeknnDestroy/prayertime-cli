package app

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type Service struct {
	resolver LocationResolver
	provider PrayerTimeProvider
	clock    Clock
}

func NewService(resolver LocationResolver, provider PrayerTimeProvider, clock Clock) *Service {
	return &Service{
		resolver: resolver,
		provider: provider,
		clock:    clock,
	}
}

func (s *Service) SearchLocations(ctx context.Context, query, countryCode string, limit int) (LocationSearchResponse, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return LocationSearchResponse{}, NewUsageError("missing required flag: --query", "", "Run 'prayertime-cli locations search --help' for usage.")
	}
	if limit <= 0 {
		return LocationSearchResponse{}, NewUsageError("limit must be greater than zero", fmt.Sprintf("%d", limit), "Use a positive --limit value.")
	}

	results, err := s.resolver.Search(ctx, query, countryCode, limit)
	if err != nil {
		return LocationSearchResponse{}, err
	}
	if len(results) == 0 {
		return LocationSearchResponse{}, NewNotFoundError(
			fmt.Sprintf("no locations matched %q", query),
			query,
			fmt.Sprintf("Run 'prayertime-cli locations search --query %q --country-code <code>' to refine the search.", query),
		)
	}

	serializedResults := make([]LocationSearchResult, 0, len(results))
	for _, result := range results {
		serializedResults = append(serializedResults, LocationSearchResult{
			Name:        result.Name,
			DisplayName: result.DisplayName(),
			Country:     result.Country,
			CountryCode: result.CountryCode,
			Admin1:      result.Admin1,
			Latitude:    result.Latitude,
			Longitude:   result.Longitude,
			Timezone:    result.Timezone,
		})
	}

	return LocationSearchResponse{
		Query:   query,
		Count:   len(serializedResults),
		Results: serializedResults,
	}, nil
}

func (s *Service) GetTimes(ctx context.Context, req TimesRequest) (TimesResponse, error) {
	location, err := s.resolveLocation(ctx, req.Query, req.CountryCode, req.Latitude, req.Longitude)
	if err != nil {
		return TimesResponse{}, err
	}

	asOf := s.clock.Now()
	schedule, err := s.fetchSchedule(ctx, location, req.Date, asOf)
	if err != nil {
		return TimesResponse{}, err
	}

	return s.toTimesResponse(location, schedule), nil
}

func (s *Service) GetCountdown(ctx context.Context, req CountdownRequest) (CountdownResponse, error) {
	target, ok := NormalizeTarget(req.Target)
	if !ok {
		return CountdownResponse{}, NewUsageError(
			fmt.Sprintf("unsupported target %q", req.Target),
			req.Target,
			"Use one of: next-prayer, imsak, fajr, sunrise, dhuhr, asr, maghrib, sunset, isha, iftar.",
		).WithDetails(ErrorDetails{ValidTargets: ValidTargets()})
	}

	location, err := s.resolveLocation(ctx, req.Query, req.CountryCode, req.Latitude, req.Longitude)
	if err != nil {
		return CountdownResponse{}, err
	}

	asOf := s.clock.Now()
	if req.At != nil {
		asOf = *req.At
	}

	schedule, err := s.fetchSchedule(ctx, location, "today", asOf)
	if err != nil {
		return CountdownResponse{}, err
	}

	localNow, err := localTimeInTimezone(schedule.Timezone, asOf)
	if err != nil {
		return CountdownResponse{}, err
	}

	schedule, selectedTarget, targetTime, err := s.resolveCountdownTarget(ctx, location, target, schedule, localNow)
	if err != nil {
		return CountdownResponse{}, err
	}

	secondsRemaining := int64(targetTime.Sub(localNow).Seconds())
	base := s.toTimesResponse(location, schedule)
	return CountdownResponse{
		TimesResponse:    base,
		Target:           selectedTarget,
		TargetAt:         targetTime.Format(time.RFC3339),
		SecondsRemaining: secondsRemaining,
		MinutesRemaining: secondsRemaining / 60,
	}, nil
}

func (s *Service) resolveCountdownTarget(ctx context.Context, location Location, target string, schedule DaySchedule, localNow time.Time) (DaySchedule, string, time.Time, error) {
	selectedTarget, targetTime := selectTarget(target, schedule, localNow)
	if countdownTargetResolved(targetTime, localNow) {
		return schedule, selectedTarget, targetTime, nil
	}

	nextDate := localNow.Add(24 * time.Hour).Format("2006-01-02")
	nextSchedule, err := s.fetchSchedule(ctx, location, nextDate, localNow)
	if err != nil {
		return DaySchedule{}, "", time.Time{}, err
	}

	selectedTarget, targetTime = selectTarget(target, nextSchedule, localNow)
	if !countdownTargetResolved(targetTime, localNow) {
		return DaySchedule{}, "", time.Time{}, NewInternalError("failed to resolve countdown target from schedule", target, "", nil)
	}

	return nextSchedule, selectedTarget, targetTime, nil
}

func (s *Service) resolveLocation(ctx context.Context, query, countryCode string, latitude, longitude *float64) (Location, error) {
	hasQuery := strings.TrimSpace(query) != ""
	hasLat := latitude != nil
	hasLon := longitude != nil

	if hasQuery && (hasLat || hasLon) {
		return Location{}, NewUsageError(
			"use either --query or --lat/--lon, not both",
			query,
			"Choose a place query or explicit coordinates.",
		)
	}

	if hasLat != hasLon {
		return Location{}, NewUsageError(
			"--lat and --lon must be provided together",
			"",
			"Provide both --lat and --lon for coordinate-based lookup.",
		)
	}

	if !hasQuery && !hasLat {
		return Location{}, NewUsageError(
			"missing location input",
			"",
			"Provide --query <place> or both --lat and --lon.",
		).WithDetails(ErrorDetails{
			RequiredOneOf: [][]string{{"query"}, {"lat", "lon"}},
		})
	}

	if hasLat && hasLon {
		return Location{
			Name:      CoordinatesName(*latitude, *longitude),
			Latitude:  *latitude,
			Longitude: *longitude,
		}, nil
	}

	results, err := s.resolver.Search(ctx, query, countryCode, 5)
	if err != nil {
		return Location{}, err
	}
	if len(results) == 0 {
		return Location{}, NewNotFoundError(
			fmt.Sprintf("no locations matched %q", query),
			query,
			fmt.Sprintf("Run 'prayertime-cli locations search --query %q --json' to inspect candidates.", query),
		)
	}

	queryToken := normalizeToken(query)
	exactMatches := make([]Location, 0, len(results))
	for _, result := range results {
		if normalizeToken(result.Name) == queryToken {
			exactMatches = append(exactMatches, result)
		}
	}

	if len(exactMatches) == 1 {
		return exactMatches[0], nil
	}
	if len(results) == 1 {
		return results[0], nil
	}

	candidates := exactMatches
	if len(candidates) == 0 {
		candidates = results
	}
	candidates = dedupeLocations(candidates)

	details := make([]LocationCandidate, 0, min(3, len(candidates)))
	for i := 0; i < len(candidates) && i < 3; i++ {
		details = append(details, candidates[i].Candidate())
	}

	return Location{}, NewAmbiguousError(
		fmt.Sprintf("location query %q is ambiguous", query),
		query,
		fmt.Sprintf("Run 'prayertime-cli locations search --query %q --json' to inspect candidates or use --lat/--lon.", query),
	).WithDetails(ErrorDetails{Candidates: details})
}

func (s *Service) fetchSchedule(ctx context.Context, location Location, dateInput string, asOf time.Time) (DaySchedule, error) {
	explicitDate := strings.TrimSpace(dateInput)
	if explicitDate == "" {
		explicitDate = "today"
	}

	requestDate, explicit, err := resolveDate(location.Timezone, explicitDate, asOf)
	if err != nil {
		return DaySchedule{}, err
	}

	schedule, err := s.provider.GetByCoordinates(ctx, location.Latitude, location.Longitude, requestDate)
	if err != nil {
		return DaySchedule{}, err
	}

	if explicit {
		return schedule, nil
	}

	localTZ, err := time.LoadLocation(schedule.Timezone)
	if err != nil {
		return DaySchedule{}, NewInternalError("failed to load location timezone", schedule.Timezone, "", err)
	}

	expectedDate := dateOnly(asOf.In(localTZ))
	if !sameDay(schedule.Date, expectedDate) {
		return s.provider.GetByCoordinates(ctx, location.Latitude, location.Longitude, expectedDate)
	}

	return schedule, nil
}

func (s *Service) toTimesResponse(location Location, schedule DaySchedule) TimesResponse {
	name := location.DisplayName()
	if name == "" {
		name = CoordinatesName(location.Latitude, location.Longitude)
	}

	return TimesResponse{
		LocationName:  name,
		Latitude:      schedule.Latitude,
		Longitude:     schedule.Longitude,
		Timezone:      schedule.Timezone,
		Date:          schedule.Date.Format("2006-01-02"),
		ImsakAt:       schedule.ImsakAt.Format(time.RFC3339),
		FajrAt:        schedule.FajrAt.Format(time.RFC3339),
		SunriseAt:     schedule.SunriseAt.Format(time.RFC3339),
		DhuhrAt:       schedule.DhuhrAt.Format(time.RFC3339),
		AsrAt:         schedule.AsrAt.Format(time.RFC3339),
		MaghribAt:     schedule.MaghribAt.Format(time.RFC3339),
		SunsetAt:      schedule.SunsetAt.Format(time.RFC3339),
		IshaAt:        schedule.IshaAt.Format(time.RFC3339),
		MethodID:      schedule.MethodID,
		MethodName:    schedule.MethodName,
		Source:        schedule.Source,
		RamadanActive: schedule.RamadanActive,
	}
}

func resolveDate(timezone, input string, asOf time.Time) (time.Time, bool, error) {
	if input != "today" {
		parsed, err := time.Parse("2006-01-02", input)
		if err != nil {
			return time.Time{}, false, NewUsageError(
				fmt.Sprintf("invalid date %q", input),
				input,
				"Use --date YYYY-MM-DD or --date today.",
			)
		}

		return dateOnly(parsed), true, nil
	}

	location := time.UTC
	if timezone != "" {
		tz, err := time.LoadLocation(timezone)
		if err != nil {
			return time.Time{}, false, NewInternalError("failed to load location timezone", timezone, "", err)
		}
		location = tz
	}

	return dateOnly(asOf.In(location)), false, nil
}

func dateOnly(value time.Time) time.Time {
	year, month, day := value.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func localTimeInTimezone(timezone string, asOf time.Time) (time.Time, error) {
	localTZ, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, NewInternalError("failed to load location timezone", timezone, "", err)
	}

	return asOf.In(localTZ), nil
}

func countdownTargetResolved(targetTime, now time.Time) bool {
	return !targetTime.IsZero() && !targetTime.Before(now)
}

func sameDay(left, right time.Time) bool {
	ly, lm, ld := left.Date()
	ry, rm, rd := right.Date()
	return ly == ry && lm == rm && ld == rd
}

func selectTarget(target string, schedule DaySchedule, now time.Time) (string, time.Time) {
	targets := []struct {
		Name string
		At   time.Time
	}{
		{Name: "imsak", At: schedule.ImsakAt},
		{Name: "fajr", At: schedule.FajrAt},
		{Name: "dhuhr", At: schedule.DhuhrAt},
		{Name: "asr", At: schedule.AsrAt},
		{Name: "maghrib", At: schedule.MaghribAt},
		{Name: "isha", At: schedule.IshaAt},
	}

	if target != "next-prayer" {
		for _, item := range targets {
			if item.Name == target {
				return item.Name, item.At
			}
		}

		switch target {
		case "sunrise":
			return "sunrise", schedule.SunriseAt
		case "sunset":
			return "sunset", schedule.SunsetAt
		default:
			return "", time.Time{}
		}
	}

	for _, item := range targets {
		if item.At.After(now) {
			return item.Name, item.At
		}
	}

	return "", time.Time{}
}

func min(left, right int) int {
	if left < right {
		return left
	}
	return right
}

func dedupeLocations(locations []Location) []Location {
	seen := make(map[string]struct{}, len(locations))
	deduped := make([]Location, 0, len(locations))
	for _, location := range locations {
		key := fmt.Sprintf(
			"%s|%s|%.6f|%.6f|%s",
			location.DisplayName(),
			location.CountryCode,
			location.Latitude,
			location.Longitude,
			location.Timezone,
		)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		deduped = append(deduped, location)
	}
	return deduped
}
