package retry

import (
	"errors"
	"log"
	"testing"
)

func TestRetry(t *testing.T) {
	testTable := []struct {
		functionSuccessOn int
		Pass	bool
	}{
		{4, false}, // MaxRetryFail
		{3, true}, // Success pre-max retries
		{3, true}, // Success on max retries
		{3, true}, // Success (infinite retries)
	}

	config := Config{
		maxRetries:        5,
		maxBackoff:        0,
		backoffMultiplier: 2,
		maxRandomJitter:   1000,
		initialDelay: 1000,
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
		if err := Retry(retryableFunction, config); err != nil {
			// Check if error is due to maxRetries limit
			if serr, ok := err.(*maxRetryError); ok {
				// Check if the error was thrown incorrectly
				if serr.maxRetries != config.maxRetries {
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
