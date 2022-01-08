package retry

import (
	"errors"
	"log"
	"reflect"
	"testing"
	"time"
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

	for testNumber, testCase := range testTable {
		log.Printf("STARTING TEST %d", testNumber+1)

		// Create retryable function that succeeds on specific iteration
		iteration := 0
		retryableFunction := func() error {
			iteration++
			if iteration == testCase.functionSuccessOn {
				return nil
			}
			return errors.New("error")
		}

		// Test Retry
		err := Retry(retryableFunction, testCase.policy)
		if err != nil {
			// Check if the test should have passed
			if testCase.shouldTestPass == true {
				t.Errorf("test should have passed: %s", err)
			}
			// Check if error is due to MaxRetries limit
			var e *maxRetryError
			if errors.As(err, &e) {
				// Check if the error was thrown incorrectly
				if err.(*maxRetryError).maxRetries != testCase.policy.MaxRetries {
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
		if testCase.shouldTestPass == false {
			t.Errorf("test should have failed")
		}
	}
}

// TestCalculateBackoff tests calculateBackoff without random jitter
func TestCalculateBackoff(t *testing.T) {
	testPolicy := &BackoffPolicy{
		MaxRetries:        10,
		MaxBackoff:        6000,
		BackoffMultiplier: 2,
		MaxRandomJitter:   0,
		InitialDelay:      500,
	}

	expectedResults := []time.Duration{500, 1000, 2000, 4000, 6000, 6000, 6000}
	for testNumber, expectedBackoff := range expectedResults {
		expectedResults[testNumber] = expectedBackoff * time.Millisecond
	}

	var backoffGrowthRate int32 = 1

	for _, expectedBackoff := range expectedResults {
		actualBackoff := calculateBackoff(backoffGrowthRate, testPolicy)
		if actualBackoff != expectedBackoff {
			t.Errorf("got: %d, expected: %d", actualBackoff, expectedBackoff)
		}
		backoffGrowthRate *= testPolicy.BackoffMultiplier
	}
}

// TestCalculateBackoff2 tests calculateBackoff with random jitter
func TestCalculateBackoff2(t *testing.T) {
	testPolicy := &BackoffPolicy{
		MaxRetries:        10,
		MaxBackoff:        6000,
		BackoffMultiplier: 2,
		MaxRandomJitter:   500,
		InitialDelay:      500,
	}

	expectedJitterResults := [][]time.Duration{
		{500, 1000}, // {minBackoff, maxBackoff}
		{1000, 2000},
		{2000, 4000},
		{4000, 8000},
		{8000, 10000},
		{10000, 10000},
		{10000, 10000},
	}
	for _, expectedBackoffRange := range expectedJitterResults {
		for testNumber, expectedValue := range expectedBackoffRange {
			expectedBackoffRange[testNumber] = expectedValue * time.Millisecond
		}
	}

	var backoffGrowthRate int32 = 1

	for _, expectedBackoff := range expectedJitterResults {
		actualBackoff := calculateBackoff(backoffGrowthRate, testPolicy)
		if actualBackoff < expectedBackoff[0] && actualBackoff > expectedBackoff[1] {
			t.Errorf("got: %d, expected: %d", actualBackoff, expectedBackoff)
		}

		backoffGrowthRate *= testPolicy.BackoffMultiplier
	}
}

func TestNewCustomBackoffPolicy(t *testing.T) {
	expectedTestPolicy := &BackoffPolicy{
		MaxRetries:        10,
		MaxBackoff:        6000,
		BackoffMultiplier: 2,
		MaxRandomJitter:   500,
		InitialDelay:      500,
	}

	actualTestPolicy := NewCustomBackoffPolicy(10, 6000, 2, 500, 500)
	if reflect.DeepEqual(actualTestPolicy, expectedTestPolicy) == false {
		t.Errorf("got: %d, expected: %d", *actualTestPolicy, *expectedTestPolicy)
	}
}

func TestNewBackoffPolicy(t *testing.T) {
	expectedTestPolicy := &BackoffPolicy{
		MaxRetries:        0,
		MaxBackoff:        0,
		BackoffMultiplier: 2,
		MaxRandomJitter:   1000,
		InitialDelay:      500,
	}

	actualTestPolicy := NewBackoffPolicy()
	if reflect.DeepEqual(actualTestPolicy, expectedTestPolicy) == false {
		t.Errorf("got: %d, expected: %d", *actualTestPolicy, *expectedTestPolicy)
	}
}

func testingPolicy(maxRetries int) *BackoffPolicy {
	policy := NewBackoffPolicy()
	policy.MaxRetries = maxRetries
	return policy
}
