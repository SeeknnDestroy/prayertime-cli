package aladhan

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/SeeknnDestroy/prayertime-cli/internal/app"
	"github.com/SeeknnDestroy/prayertime-cli/internal/providers/httpx"
)

const (
	defaultBaseURL = "https://api.aladhan.com"
	methodID       = 13
	sourceName     = "aladhan:method=13"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(httpClient *http.Client, baseURL string) *Client {
	if baseURL == "" {
		baseURL = defaultBaseURL
	}

	return &Client{
		baseURL:    baseURL,
		httpClient: httpClient,
	}
}

func (c *Client) GetByCoordinates(ctx context.Context, latitude, longitude float64, date time.Time) (app.DaySchedule, error) {
	endpoint, err := url.Parse(c.baseURL + "/v1/timings/" + date.Format("02-01-2006"))
	if err != nil {
		return app.DaySchedule{}, app.NewInternalError("failed to build AlAdhan URL", c.baseURL, "", err)
	}

	params := endpoint.Query()
	params.Set("latitude", fmt.Sprintf("%.6f", latitude))
	params.Set("longitude", fmt.Sprintf("%.6f", longitude))
	params.Set("method", fmt.Sprintf("%d", methodID))
	endpoint.RawQuery = params.Encode()

	var payload timingsResponse
	if err := httpx.GetJSON(ctx, c.httpClient, endpoint.String(), &payload); err != nil {
		return app.DaySchedule{}, err
	}

	location, err := time.LoadLocation(payload.Data.Meta.Timezone)
	if err != nil {
		return app.DaySchedule{}, app.NewInternalError("failed to load AlAdhan timezone", payload.Data.Meta.Timezone, "", err)
	}

	dateValue, err := time.ParseInLocation("02-01-2006", payload.Data.Date.Gregorian.Date, location)
	if err != nil {
		return app.DaySchedule{}, app.NewInternalError("failed to parse AlAdhan date", payload.Data.Date.Gregorian.Date, "", err)
	}

	imsakAt, err := parseClock(dateValue, payload.Data.Timings.Imsak, location)
	if err != nil {
		return app.DaySchedule{}, err
	}
	fajrAt, err := parseClock(dateValue, payload.Data.Timings.Fajr, location)
	if err != nil {
		return app.DaySchedule{}, err
	}
	sunriseAt, err := parseClock(dateValue, payload.Data.Timings.Sunrise, location)
	if err != nil {
		return app.DaySchedule{}, err
	}
	dhuhrAt, err := parseClock(dateValue, payload.Data.Timings.Dhuhr, location)
	if err != nil {
		return app.DaySchedule{}, err
	}
	asrAt, err := parseClock(dateValue, payload.Data.Timings.Asr, location)
	if err != nil {
		return app.DaySchedule{}, err
	}
	maghribAt, err := parseClock(dateValue, payload.Data.Timings.Maghrib, location)
	if err != nil {
		return app.DaySchedule{}, err
	}
	sunsetAt, err := parseClock(dateValue, payload.Data.Timings.Sunset, location)
	if err != nil {
		return app.DaySchedule{}, err
	}
	ishaAt, err := parseClock(dateValue, payload.Data.Timings.Isha, location)
	if err != nil {
		return app.DaySchedule{}, err
	}

	return app.DaySchedule{
		Latitude:      payload.Data.Meta.Latitude,
		Longitude:     payload.Data.Meta.Longitude,
		Timezone:      payload.Data.Meta.Timezone,
		Date:          dateValue,
		ImsakAt:       imsakAt,
		FajrAt:        fajrAt,
		SunriseAt:     sunriseAt,
		DhuhrAt:       dhuhrAt,
		AsrAt:         asrAt,
		MaghribAt:     maghribAt,
		SunsetAt:      sunsetAt,
		IshaAt:        ishaAt,
		MethodID:      payload.Data.Meta.Method.ID,
		MethodName:    payload.Data.Meta.Method.Name,
		Source:        sourceName,
		RamadanActive: payload.Data.Date.Hijri.Month.Number == 9,
	}, nil
}

func parseClock(date time.Time, value string, location *time.Location) (time.Time, error) {
	clock := strings.TrimSpace(strings.Split(value, " ")[0])
	parsed, err := time.ParseInLocation("15:04", clock, location)
	if err != nil {
		return time.Time{}, app.NewInternalError("failed to parse prayer time value", value, "", err)
	}

	return time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		parsed.Hour(),
		parsed.Minute(),
		0,
		0,
		location,
	), nil
}

type timingsResponse struct {
	Data struct {
		Timings struct {
			Imsak   string `json:"Imsak"`
			Fajr    string `json:"Fajr"`
			Sunrise string `json:"Sunrise"`
			Dhuhr   string `json:"Dhuhr"`
			Asr     string `json:"Asr"`
			Maghrib string `json:"Maghrib"`
			Sunset  string `json:"Sunset"`
			Isha    string `json:"Isha"`
		} `json:"timings"`
		Date struct {
			Gregorian struct {
				Date string `json:"date"`
			} `json:"gregorian"`
			Hijri struct {
				Month struct {
					Number int `json:"number"`
				} `json:"month"`
			} `json:"hijri"`
		} `json:"date"`
		Meta struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
			Timezone  string  `json:"timezone"`
			Method    struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"method"`
		} `json:"meta"`
	} `json:"data"`
}
