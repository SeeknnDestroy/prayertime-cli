package cli

import (
	"fmt"
	"sort"
	"strings"

	"github.com/SeeknnDestroy/prayertime-cli/internal/app"
)

type timesFieldDefinition struct {
	Key          string
	Aliases      []string
	Label        string
	DetailedOnly bool
	ValidOrder   int
	Value        func(app.TimesResponse) any
}

type countdownFieldDefinition struct {
	Key          string
	Aliases      []string
	Label        string
	DetailedOnly bool
	ValidOrder   int
	Value        func(app.CountdownResponse) any
}

var cliFieldNormalizer = strings.NewReplacer(
	"ç", "c",
	"ğ", "g",
	"ı", "i",
	"İ", "i",
	"ö", "o",
	"ş", "s",
	"ü", "u",
)

var timesFieldDefinitions = []timesFieldDefinition{
	{Key: "location_name", Aliases: []string{"location-name"}, Label: "Location", ValidOrder: 0, Value: func(response app.TimesResponse) any { return response.LocationName }},
	{Key: "timezone", Label: "Timezone", ValidOrder: 3, Value: func(response app.TimesResponse) any { return response.Timezone }},
	{Key: "date", Label: "Date", ValidOrder: 4, Value: func(response app.TimesResponse) any { return response.Date }},
	{Key: "imsak_at", Aliases: []string{"imsak"}, Label: "Imsak", ValidOrder: 5, Value: func(response app.TimesResponse) any { return response.ImsakAt }},
	{Key: "fajr_at", Aliases: []string{"fajr"}, Label: "Fajr", ValidOrder: 6, Value: func(response app.TimesResponse) any { return response.FajrAt }},
	{Key: "sunrise_at", Aliases: []string{"sunrise", "gunes"}, Label: "Sunrise", ValidOrder: 7, Value: func(response app.TimesResponse) any { return response.SunriseAt }},
	{Key: "dhuhr_at", Aliases: []string{"dhuhr", "ogle"}, Label: "Dhuhr", ValidOrder: 8, Value: func(response app.TimesResponse) any { return response.DhuhrAt }},
	{Key: "asr_at", Aliases: []string{"asr", "ikindi"}, Label: "Asr", ValidOrder: 9, Value: func(response app.TimesResponse) any { return response.AsrAt }},
	{Key: "maghrib_at", Aliases: []string{"maghrib", "iftar", "aksam"}, Label: "Maghrib", ValidOrder: 10, Value: func(response app.TimesResponse) any { return response.MaghribAt }},
	{Key: "sunset_at", Aliases: []string{"sunset"}, Label: "Sunset", ValidOrder: 11, Value: func(response app.TimesResponse) any { return response.SunsetAt }},
	{Key: "isha_at", Aliases: []string{"isha", "yatsi"}, Label: "Isha", ValidOrder: 12, Value: func(response app.TimesResponse) any { return response.IshaAt }},
	{Key: "ramadan_active", Aliases: []string{"ramadan", "ramadan-active"}, Label: "Ramadan active", ValidOrder: 16, Value: func(response app.TimesResponse) any { return response.RamadanActive }},
	{Key: "latitude", Label: "Latitude", DetailedOnly: true, ValidOrder: 1, Value: func(response app.TimesResponse) any { return response.Latitude }},
	{Key: "longitude", Label: "Longitude", DetailedOnly: true, ValidOrder: 2, Value: func(response app.TimesResponse) any { return response.Longitude }},
	{Key: "method_id", Aliases: []string{"method-id"}, Label: "Method ID", DetailedOnly: true, ValidOrder: 13, Value: func(response app.TimesResponse) any { return response.MethodID }},
	{Key: "method_name", Aliases: []string{"method-name"}, Label: "Method", DetailedOnly: true, ValidOrder: 14, Value: func(response app.TimesResponse) any { return response.MethodName }},
	{Key: "source", Label: "Source", DetailedOnly: true, ValidOrder: 15, Value: func(response app.TimesResponse) any { return response.Source }},
}

var countdownFieldDefinitions = []countdownFieldDefinition{
	{Key: "location_name", Aliases: []string{"location-name"}, Label: "Location", ValidOrder: 0, Value: func(response app.CountdownResponse) any { return response.LocationName }},
	{Key: "timezone", Label: "Timezone", ValidOrder: 3, Value: func(response app.CountdownResponse) any { return response.Timezone }},
	{Key: "date", Label: "Date", ValidOrder: 4, Value: func(response app.CountdownResponse) any { return response.Date }},
	{Key: "target", Label: "Target", ValidOrder: 17, Value: func(response app.CountdownResponse) any { return response.Target }},
	{Key: "target_at", Aliases: []string{"target-at"}, Label: "Target at", ValidOrder: 18, Value: func(response app.CountdownResponse) any { return response.TargetAt }},
	{Key: "seconds_remaining", Aliases: []string{"seconds", "seconds-remaining"}, Label: "Seconds remaining", ValidOrder: 19, Value: func(response app.CountdownResponse) any { return response.SecondsRemaining }},
	{Key: "minutes_remaining", Aliases: []string{"minutes", "minutes-remaining"}, Label: "Minutes remaining", ValidOrder: 20, Value: func(response app.CountdownResponse) any { return response.MinutesRemaining }},
	{Key: "imsak_at", Aliases: []string{"imsak"}, Label: "Imsak", DetailedOnly: true, ValidOrder: 5, Value: func(response app.CountdownResponse) any { return response.ImsakAt }},
	{Key: "fajr_at", Aliases: []string{"fajr"}, Label: "Fajr", DetailedOnly: true, ValidOrder: 6, Value: func(response app.CountdownResponse) any { return response.FajrAt }},
	{Key: "sunrise_at", Aliases: []string{"sunrise", "gunes"}, Label: "Sunrise", DetailedOnly: true, ValidOrder: 7, Value: func(response app.CountdownResponse) any { return response.SunriseAt }},
	{Key: "dhuhr_at", Aliases: []string{"dhuhr", "ogle"}, Label: "Dhuhr", DetailedOnly: true, ValidOrder: 8, Value: func(response app.CountdownResponse) any { return response.DhuhrAt }},
	{Key: "asr_at", Aliases: []string{"asr", "ikindi"}, Label: "Asr", DetailedOnly: true, ValidOrder: 9, Value: func(response app.CountdownResponse) any { return response.AsrAt }},
	{Key: "maghrib_at", Aliases: []string{"maghrib", "iftar", "aksam"}, Label: "Maghrib", DetailedOnly: true, ValidOrder: 10, Value: func(response app.CountdownResponse) any { return response.MaghribAt }},
	{Key: "sunset_at", Aliases: []string{"sunset"}, Label: "Sunset", DetailedOnly: true, ValidOrder: 11, Value: func(response app.CountdownResponse) any { return response.SunsetAt }},
	{Key: "isha_at", Aliases: []string{"isha", "yatsi"}, Label: "Isha", DetailedOnly: true, ValidOrder: 12, Value: func(response app.CountdownResponse) any { return response.IshaAt }},
	{Key: "ramadan_active", Aliases: []string{"ramadan", "ramadan-active"}, Label: "Ramadan active", DetailedOnly: true, ValidOrder: 16, Value: func(response app.CountdownResponse) any { return response.RamadanActive }},
	{Key: "latitude", Label: "Latitude", DetailedOnly: true, ValidOrder: 1, Value: func(response app.CountdownResponse) any { return response.Latitude }},
	{Key: "longitude", Label: "Longitude", DetailedOnly: true, ValidOrder: 2, Value: func(response app.CountdownResponse) any { return response.Longitude }},
	{Key: "method_id", Aliases: []string{"method-id"}, Label: "Method ID", DetailedOnly: true, ValidOrder: 13, Value: func(response app.CountdownResponse) any { return response.MethodID }},
	{Key: "method_name", Aliases: []string{"method-name"}, Label: "Method", DetailedOnly: true, ValidOrder: 14, Value: func(response app.CountdownResponse) any { return response.MethodName }},
	{Key: "source", Label: "Source", DetailedOnly: true, ValidOrder: 15, Value: func(response app.CountdownResponse) any { return response.Source }},
}

var timesFieldLookup = buildTimesFieldLookup(timesFieldDefinitions)
var countdownFieldLookup = buildCountdownFieldLookup(countdownFieldDefinitions)
var orderedTimesFieldKeys = orderedTimesFields(timesFieldDefinitions)
var orderedCountdownFieldKeys = orderedCountdownFields(countdownFieldDefinitions)

func normalizeTimesField(field string) (timesFieldDefinition, bool) {
	definition, ok := timesFieldLookup[normalizeCLIFieldToken(field)]
	return definition, ok
}

func normalizeCountdownField(field string) (countdownFieldDefinition, bool) {
	definition, ok := countdownFieldLookup[normalizeCLIFieldToken(field)]
	return definition, ok
}

func validTimesFields() []string {
	return append([]string(nil), orderedTimesFieldKeys...)
}

func validCountdownFields() []string {
	return append([]string(nil), orderedCountdownFieldKeys...)
}

func validateTimesField(field string) error {
	if field == "" {
		return nil
	}

	if _, ok := normalizeTimesField(field); ok {
		return nil
	}

	return newTimesFieldUsageError(field)
}

func validateCountdownField(field string) error {
	if field == "" {
		return nil
	}

	if _, ok := normalizeCountdownField(field); ok {
		return nil
	}

	return newCountdownFieldUsageError(field)
}

func timesFieldEntry(response app.TimesResponse, field string) (fieldEntry, error) {
	definition, ok := normalizeTimesField(field)
	if !ok {
		return fieldEntry{}, newTimesFieldUsageError(field)
	}

	return fieldEntry{
		Key:   definition.Key,
		Label: definition.Label,
		Value: definition.Value(response),
	}, nil
}

func countdownFieldEntry(response app.CountdownResponse, field string) (fieldEntry, error) {
	definition, ok := normalizeCountdownField(field)
	if !ok {
		return fieldEntry{}, newCountdownFieldUsageError(field)
	}

	return fieldEntry{
		Key:   definition.Key,
		Label: definition.Label,
		Value: definition.Value(response),
	}, nil
}

func timesEntries(response app.TimesResponse, view viewMode) []fieldEntry {
	entries := make([]fieldEntry, 0, len(timesFieldDefinitions))
	for _, definition := range timesFieldDefinitions {
		if definition.DetailedOnly && view != viewDetailed {
			continue
		}

		entries = append(entries, fieldEntry{
			Key:   definition.Key,
			Label: definition.Label,
			Value: definition.Value(response),
		})
	}
	return entries
}

func countdownEntries(response app.CountdownResponse, view viewMode) []fieldEntry {
	entries := make([]fieldEntry, 0, len(countdownFieldDefinitions))
	for _, definition := range countdownFieldDefinitions {
		if definition.DetailedOnly && view != viewDetailed {
			continue
		}

		entries = append(entries, fieldEntry{
			Key:   definition.Key,
			Label: definition.Label,
			Value: definition.Value(response),
		})
	}
	return entries
}

func normalizeCLIFieldToken(input string) string {
	value := strings.TrimSpace(strings.ToLower(input))
	value = cliFieldNormalizer.Replace(value)
	value = strings.ReplaceAll(value, "_", "-")
	value = strings.ReplaceAll(value, " ", "-")
	return value
}

func newTimesFieldUsageError(field string) error {
	return app.NewUsageError(
		fmt.Sprintf("unsupported field %q", field),
		field,
		"Use a canonical prayer-times field such as maghrib_at, timezone, method_name, or source.",
	).WithDetails(app.ErrorDetails{ValidFields: validTimesFields()})
}

func newCountdownFieldUsageError(field string) error {
	return app.NewUsageError(
		fmt.Sprintf("unsupported field %q", field),
		field,
		"Use a canonical countdown field such as seconds_remaining, target_at, maghrib_at, or method_name.",
	).WithDetails(app.ErrorDetails{ValidFields: validCountdownFields()})
}

func buildTimesFieldLookup(definitions []timesFieldDefinition) map[string]timesFieldDefinition {
	lookup := make(map[string]timesFieldDefinition, len(definitions)*2)
	for _, definition := range definitions {
		lookup[normalizeCLIFieldToken(definition.Key)] = definition
		for _, alias := range definition.Aliases {
			lookup[normalizeCLIFieldToken(alias)] = definition
		}
	}
	return lookup
}

func buildCountdownFieldLookup(definitions []countdownFieldDefinition) map[string]countdownFieldDefinition {
	lookup := make(map[string]countdownFieldDefinition, len(definitions)*2)
	for _, definition := range definitions {
		lookup[normalizeCLIFieldToken(definition.Key)] = definition
		for _, alias := range definition.Aliases {
			lookup[normalizeCLIFieldToken(alias)] = definition
		}
	}
	return lookup
}

func orderedTimesFields(definitions []timesFieldDefinition) []string {
	ordered := append([]timesFieldDefinition(nil), definitions...)
	sort.Slice(ordered, func(i, j int) bool {
		return ordered[i].ValidOrder < ordered[j].ValidOrder
	})

	keys := make([]string, 0, len(ordered))
	for _, definition := range ordered {
		keys = append(keys, definition.Key)
	}
	return keys
}

func orderedCountdownFields(definitions []countdownFieldDefinition) []string {
	ordered := append([]countdownFieldDefinition(nil), definitions...)
	sort.Slice(ordered, func(i, j int) bool {
		return ordered[i].ValidOrder < ordered[j].ValidOrder
	})

	keys := make([]string, 0, len(ordered))
	for _, definition := range ordered {
		keys = append(keys, definition.Key)
	}
	return keys
}
