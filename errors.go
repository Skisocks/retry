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

// inputError is an error returned when the input of a function is invalid
type inputError struct {
	err string
}

func (i *inputError) Error() string {
	return fmt.Sprintf("input error: %s", i.err)
}
