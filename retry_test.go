package retry

import (
	"errors"
	"log"
	"testing"
)

func TestRetry(t *testing.T) {
	testTable := []struct {
		maxRetries int
		functionSuccessOn int
	}{
		{3, 4}, // MaxRetryFail
		{4, 3}, // Success
		{0, 3}, // Success (infinite retries)
	}

	for testNumber, table := range testTable {
		log.Printf("STARTING TEST %d", testNumber+1)

		// Create retryable function that succeeds on specific iteration
		var iteration = 0
		retryableFunction := func() error {
			iteration++
			if iteration == table.functionSuccessOn {
				return nil
			}
			return errors.New("error")
		}

		// Test Retry()
		if err := Retry(retryableFunction, table.maxRetries); err != nil {
			// Check if error is due to maxRetries limit
			if serr, ok := err.(*maxRetryError); ok {
				// Check if the error was thrown incorrectly
				if serr.maxRetries != table.maxRetries {
					t.Errorf("maxRetryError thrown incorrectly: %s", serr)
					continue
				}
				// If correct log and continue
				log.Printf("max retry error: %s", serr)
				continue
			}
			// Check if any other errors are thrown
			t.Errorf("unknown error: %s", err)
		}
	}
}