package retry

import (
	"errors"
	"log"
	"testing"
)

func TestRetry(t *testing.T) {
	testTable := []struct {
		functionSuccessOn int
		policy            *BackoffPolicy
		shouldTestPass    bool
	}{
		{4, testingPolicy(3), false}, // MaxRetryFail
		{3, testingPolicy(4), true},  // Success pre-max retries
		{3, testingPolicy(3), true},  // Success on max retries
		{3, testingPolicy(0), true},  // Success (infinite retries)
	}

	for testNumber, table := range testTable {
		log.Printf("STARTING TEST %d", testNumber+1)

		// Create retryable function that succeeds on specific iteration
		iteration := 0
		retryableFunction := func() error {
			iteration++
			if iteration == table.functionSuccessOn {
				return nil
			}
			return errors.New("error")
		}

		// Test Retry()
		err := Retry(retryableFunction, table.policy)
		if err != nil {
			// Check if the test should have passed
			if table.shouldTestPass == true {
				t.Errorf("test should have passed: %s", err)
			}
			// Check if error is due to MaxRetries limit
			var e *maxRetryError
			if errors.As(err, &e) {
				// Check if the error was thrown incorrectly
				if err.(*maxRetryError).maxRetries != table.policy.MaxRetries {
					t.Errorf("maxRetryError thrown incorrectly: %s", err)
					continue
				}
				// If correct log and continue
				log.Printf("max retry error: %s", err)
				continue
			}
			// Check if it was an unknown error
			t.Errorf("unknown error: %s", err)
			continue
		}
		// Check if the test should have failed
		if table.shouldTestPass == false {
			t.Errorf("test should have failed")
		}
	}
}

func testingPolicy(maxRetries int) *BackoffPolicy {
	policy := NewBackoffPolicy()
	policy.MaxRetries = maxRetries
	return policy
}
