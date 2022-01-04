package retry

import (
	"fmt"
	"time"
)

type MaxRetryError struct {
	maxRetries	int
	backoff	time.Duration
}

func (m *MaxRetryError) Error() string {
	return fmt.Sprintf("function was not successful before max retries, failed after %d attempts", m.maxRetries)
}
