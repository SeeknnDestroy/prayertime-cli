package httpx

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/SeeknnDestroy/prayertime-cli/internal/app"
)

var waitBeforeRetry = sleepWithContext

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
				if waitErr := waitBeforeRetry(ctx, backoff); waitErr != nil {
					return newNetworkCLIError("request cancelled while waiting to retry", url, waitErr)
				}
				backoff *= 2
				continue
			}
			if isTimeoutError(err) {
				return app.NewNetworkTimeoutError("upstream request timed out after retries", url, "Retry the command in a few seconds.", err)
			}
			return app.NewNetworkError("upstream request failed before a response was received", url, "Check network connectivity or upstream TLS/proxy settings and retry.", err)
		}

		if shouldRetry(nil, resp.StatusCode) && attempt < 2 {
			lastErr = fmt.Errorf("upstream status %d", resp.StatusCode)
			closeBody(resp.Body)
			if waitErr := waitBeforeRetry(ctx, backoff); waitErr != nil {
				return newNetworkCLIError("request cancelled while waiting to retry", url, waitErr)
			}
			backoff *= 2
			continue
		}

		if resp.StatusCode >= 400 {
			body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
			closeBody(resp.Body)
			return app.NewInternalError(
				fmt.Sprintf("upstream request failed with status %d", resp.StatusCode),
				url,
				string(body),
				lastErr,
			)
		}

		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(target); err != nil {
			closeBody(resp.Body)
			return app.NewInternalError("failed to decode upstream response", url, "", err)
		}
		closeBody(resp.Body)

		return nil
	}

	return newNetworkCLIError("upstream request failed before a response was received", url, lastErr)
}

func shouldRetry(err error, statusCode int) bool {
	if err != nil {
		return isTimeoutError(err)
	}

	return statusCode == http.StatusTooManyRequests || statusCode >= http.StatusInternalServerError
}

func isTimeoutError(err error) bool {
	if err == nil {
		return false
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}

	var netErr net.Error
	return errors.As(err, &netErr) && netErr.Timeout()
}

func newNetworkCLIError(message, url string, cause error) error {
	if isTimeoutError(cause) {
		return app.NewNetworkTimeoutError(message, url, "Retry the command in a few seconds.", cause)
	}

	return app.NewNetworkError(message, url, "Check network connectivity or upstream TLS/proxy settings and retry.", cause)
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

func closeBody(closer io.Closer) {
	_ = closer.Close()
}
