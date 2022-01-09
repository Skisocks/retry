package retry

import (
	"bytes"
	"errors"
	"log"
	"os"
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
	testCases := []struct {
		inputPolicy     *BackoffPolicy
		expectedBackoff []time.Duration
	}{
		{ // Base test
			&BackoffPolicy{
				MaxRetries:        0,
				MaxBackoff:        6000,
				BackoffMultiplier: 2,
				MaxRandomJitter:   0,
				InitialDelay:      500,
				IsLogging:         false,
			},
			[]time.Duration{500, 1000, 2000, 4000, 6000, 6000, 6000},
		},
		{ // Test InitialDelay
			&BackoffPolicy{
				MaxRetries:        0,
				MaxBackoff:        16000,
				BackoffMultiplier: 2,
				MaxRandomJitter:   0,
				InitialDelay:      1000,
				IsLogging:         false,
			},
			[]time.Duration{1000, 2000, 4000, 8000, 16000, 16000, 16000},
		},
		{ // Test BackoffMultiplier
			&BackoffPolicy{
				MaxRetries:        0,
				MaxBackoff:        27000,
				BackoffMultiplier: 3,
				MaxRandomJitter:   0,
				InitialDelay:      1000,
				IsLogging:         false,
			},
			[]time.Duration{1000, 3000, 9000, 27000, 27000, 27000, 27000},
		},

		{ // Test no BackoffMultiplier
			&BackoffPolicy{
				MaxRetries:        0,
				MaxBackoff:        27000,
				BackoffMultiplier: 0,
				MaxRandomJitter:   0,
				InitialDelay:      1000,
				IsLogging:         false,
			},
			[]time.Duration{1000, 2000, 3000, 4000, 5000, 6000, 7000},
		},
	}

	for _, testCase := range testCases {
		var backoffGrowthRate int32 = 1

		for _, expectedBackoff := range testCase.expectedBackoff {
			expectedBackoff = expectedBackoff * time.Millisecond

			actualBackoff := calculateBackoff(backoffGrowthRate, testCase.inputPolicy)
			if actualBackoff != expectedBackoff {
				t.Errorf("got: %d, expected: %d", actualBackoff/time.Millisecond, expectedBackoff/time.Millisecond)
			}
			backoffGrowthRate *= testCase.inputPolicy.BackoffMultiplier
		}

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

func TestIsLogging(t *testing.T) {
	loggingOnPolicy := &BackoffPolicy{
		MaxRetries:        10,
		MaxBackoff:        6000,
		BackoffMultiplier: 2,
		MaxRandomJitter:   0,
		InitialDelay:      500,
		IsLogging:         true,
	}

	loggingOffPolicy := &BackoffPolicy{
		MaxRetries:        10,
		MaxBackoff:        6000,
		BackoffMultiplier: 2,
		MaxRandomJitter:   0,
		InitialDelay:      500,
		IsLogging:         false,
	}

	testTable := []struct {
		testPolicy *BackoffPolicy
	}{
		{loggingOnPolicy},
		{loggingOffPolicy}, // Success on max retries3
	}

	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	for _, testCase := range testTable {
		buf := new(bytes.Buffer)
		log.SetOutput(buf)

		defer func() {
			log.SetOutput(os.Stderr)
		}()
		isLogging(testCase.testPolicy, "%d %d %d", 1, 2, 3)
		result := buf.String()

		if testCase.testPolicy.IsLogging == true && result != "1 2 3\n" {
			t.Errorf("expected (1 2 3) got %s", result)
			continue
		}

		if testCase.testPolicy.IsLogging == false && result != "" {
			t.Errorf("expected () got %s", result)
		}
	}
}

func BenchmarkCalculateBackoff(b *testing.B) {
	benchmarkPolicy := &BackoffPolicy{
		MaxRetries:        10,
		MaxBackoff:        0,
		BackoffMultiplier: 2,
		MaxRandomJitter:   1000,
		InitialDelay:      500,
		IsLogging:         false,
	}

	for i := 0; i < b.N; i++ {
		var backoffGrowthRate int32 = 1
		for i := 0; i <= benchmarkPolicy.MaxRetries; i++ {
			calculateBackoff(backoffGrowthRate, benchmarkPolicy)
			backoffGrowthRate *= benchmarkPolicy.BackoffMultiplier
		}
	}
}

func ExampleRetry() {
	// A function that may fail returning an error
	retryableFunction := func() error { return nil }

	if err := Retry(retryableFunction, NewBackoffPolicy()); err != nil {
		// Handle error
		return
	}
	// Output:
}

func ExampleRetry_second() {
	// A custom backoff policy
	myPolicy, _ := NewCustomBackoffPolicy(10, 0, 2, 1000, 1000, false)

	// A function that may fail returning an error
	retryableFunction := func() error { return nil }

	if err := Retry(retryableFunction, myPolicy); err != nil {
		// Handle error
		return
	}
	// Output:
}
