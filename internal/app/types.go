package app

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type Location struct {
	Name        string  `json:"name"`
	Country     string  `json:"country"`
	CountryCode string  `json:"country_code"`
	Admin1      string  `json:"admin1,omitempty"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Timezone    string  `json:"timezone"`
}

func (l Location) DisplayName() string {
	parts := []string{l.Name}
	if l.Admin1 != "" && !strings.EqualFold(l.Admin1, l.Name) {
		parts = append(parts, l.Admin1)
	}
	if l.Country != "" {
		parts = append(parts, l.Country)
	}
	return strings.Join(parts, ", ")
}

func CoordinatesName(lat, lon float64) string {
	return fmt.Sprintf("%.6f, %.6f", lat, lon)
}

type LocationSearchResponse struct {
	Query   string     `json:"query"`
	Count   int        `json:"count"`
	Results []Location `json:"results"`
}

type TimesRequest struct {
	Query       string
	CountryCode string
	Date        string
	Latitude    *float64
	Longitude   *float64
}

type CountdownRequest struct {
	Query       string
	CountryCode string
	Target      string
	At          *time.Time
	Latitude    *float64
	Longitude   *float64
}

type TimesResponse struct {
	LocationName  string  `json:"location_name"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	Timezone      string  `json:"timezone"`
	Date          string  `json:"date"`
	ImsakAt       string  `json:"imsak_at"`
	FajrAt        string  `json:"fajr_at"`
	SunriseAt     string  `json:"sunrise_at"`
	DhuhrAt       string  `json:"dhuhr_at"`
	AsrAt         string  `json:"asr_at"`
	MaghribAt     string  `json:"maghrib_at"`
	SunsetAt      string  `json:"sunset_at"`
	IshaAt        string  `json:"isha_at"`
	MethodID      int     `json:"method_id"`
	MethodName    string  `json:"method_name"`
	Source        string  `json:"source"`
	RamadanActive bool    `json:"ramadan_active"`
}

type CountdownResponse struct {
	TimesResponse
	Target           string `json:"target"`
	TargetAt         string `json:"target_at"`
	SecondsRemaining int64  `json:"seconds_remaining"`
	MinutesRemaining int64  `json:"minutes_remaining"`
}

type DaySchedule struct {
	Latitude      float64
	Longitude     float64
	Timezone      string
	Date          time.Time
	ImsakAt       time.Time
	FajrAt        time.Time
	SunriseAt     time.Time
	DhuhrAt       time.Time
	AsrAt         time.Time
	MaghribAt     time.Time
	SunsetAt      time.Time
	IshaAt        time.Time
	MethodID      int
	MethodName    string
	Source        string
	RamadanActive bool
}

type LocationResolver interface {
	Search(ctx context.Context, query, countryCode string, limit int) ([]Location, error)
}

type PrayerTimeProvider interface {
	GetByCoordinates(ctx context.Context, latitude, longitude float64, date time.Time) (DaySchedule, error)
}
