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
		expectedIsError   bool
	}{
		{4, true},  // MaxRetryFail 3
		{2, false}, // Success pre-max retries4
		{3, false}, // Success on max retries3
	}

	inputPolicy := &BackoffPolicy{
		MaxRetries:        3,
		MaxBackoff:        0,
		BackoffMultiplier: 2,
		MaxRandomJitter:   1000,
		InitialDelay:      500,
		IsLogging:         false,
	}

	for _, testCase := range testTable {
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
		err := Retry(retryableFunction, inputPolicy)
		if err != nil {
			// Check if an error is expected
			if testCase.expectedIsError == false {
				t.Errorf("error not expected: %s", err)
			}
			// Check if error is due to MaxRetries limit
			var e *maxRetryError
			if errors.As(err, &e) {
				// Check if the error was thrown incorrectly
				if err.(*maxRetryError).maxRetries != inputPolicy.MaxRetries {
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
		// Check if an error is expected
		if testCase.expectedIsError == true {
			t.Errorf("expected error")
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
		IsLogging:         false,
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

	actualTestPolicy := NewCustomBackoffPolicy(10, 6000, 2, 500, 500, false)
	if reflect.DeepEqual(actualTestPolicy, expectedTestPolicy) == false {
		t.Errorf("got: %+v, expected: %+v", *actualTestPolicy, *expectedTestPolicy)
	}
}

func TestNewBackoffPolicy(t *testing.T) {
	expectedTestPolicy := &BackoffPolicy{
		MaxRetries:        0,
		MaxBackoff:        0,
		BackoffMultiplier: 2,
		MaxRandomJitter:   1000,
		InitialDelay:      500,
		IsLogging:         false,
	}

	actualTestPolicy := NewBackoffPolicy()
	if reflect.DeepEqual(actualTestPolicy, expectedTestPolicy) == false {
		t.Errorf("got: %+v, expected: %+v", *actualTestPolicy, *expectedTestPolicy)
	}
}
