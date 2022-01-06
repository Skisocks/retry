package retry

import (
	"fmt"
	"time"
)

// maxRetryError is an error returned when Retry does not succeed within BackoffPolicy.MaxRetries
type maxRetryError struct {
	maxRetries int
	backoff    time.Duration
}

func (m *maxRetryError) Error() string {
	return fmt.Sprintf("function was not successful before max retries, failed after %d attempts", m.maxRetries)
}
