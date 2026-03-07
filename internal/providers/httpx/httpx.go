package httpx

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/SeeknnDestroy/prayertime-cli/internal/app"
)

func GetJSON(ctx context.Context, client *http.Client, url string, target any) error {
	var lastErr error
	backoff := 250 * time.Millisecond

	for attempt := 0; attempt < 3; attempt++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return app.NewInternalError("failed to create upstream request", url, "", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			if shouldRetry(err, 0) && attempt < 2 {
				if waitErr := sleepWithContext(ctx, backoff); waitErr != nil {
					return app.NewNetworkError("request cancelled while waiting to retry", url, "Retry the command in a few seconds.", waitErr)
				}
				backoff *= 2
				continue
			}
			return app.NewNetworkError("upstream request timed out after retries", url, "Retry the command in a few seconds.", err)
		}

		if shouldRetry(nil, resp.StatusCode) && attempt < 2 {
			lastErr = fmt.Errorf("upstream status %d", resp.StatusCode)
			resp.Body.Close()
			if waitErr := sleepWithContext(ctx, backoff); waitErr != nil {
				return app.NewNetworkError("request cancelled while waiting to retry", url, "Retry the command in a few seconds.", waitErr)
			}
			backoff *= 2
			continue
		}

		if resp.StatusCode >= 400 {
			body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
			resp.Body.Close()
			return app.NewInternalError(
				fmt.Sprintf("upstream request failed with status %d", resp.StatusCode),
				url,
				string(body),
				lastErr,
			)
		}

		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(target); err != nil {
			resp.Body.Close()
			return app.NewInternalError("failed to decode upstream response", url, "", err)
		}
		resp.Body.Close()

		return nil
	}

	return app.NewNetworkError("upstream request timed out after retries", url, "Retry the command in a few seconds.", lastErr)
}

func shouldRetry(err error, statusCode int) bool {
	if err != nil {
		if netErr, ok := err.(net.Error); ok && (netErr.Timeout() || netErr.Temporary()) {
			return true
		}
		return true
	}

	return statusCode == http.StatusTooManyRequests || statusCode >= http.StatusInternalServerError
}

func sleepWithContext(ctx context.Context, duration time.Duration) error {
	timer := time.NewTimer(duration)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}
