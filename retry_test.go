package retry

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"
	"time"
)

func TestRetry(t *testing.T) {
	testCases := []struct {
		name              string
		functionSuccessOn int
		expectedIsError   bool
	}{
		{name: "max retry error", functionSuccessOn: 4, expectedIsError: true},
		{name: "success pre-max retry", functionSuccessOn: 2},
		{name: "success on max retry", functionSuccessOn: 3},
	}

	inputPolicy := &BackoffPolicy{
		MaxRetries:        3,
		MaxBackoff:        0,
		BackoffMultiplier: 2,
		MaxRandomJitter:   0,
		InitialDelay:      0,
		IsLogging:         false,
	}

	for _, testCase := range testCases {
		// Create retryable function that succeeds on specific iteration
		iteration := 0
		retryableFunction := func() error {
			iteration++
			if iteration == testCase.functionSuccessOn {
				return nil
			}
			return errors.New("error")
		}

		testName := fmt.Sprintf("%s test", testCase.name)

		t.Run(testName, func(t *testing.T) {
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
						return
					}
					// If correct continue
					return
				}
				// Check if it was an unknown error
				t.Errorf("unknown error: %s", err)
				return
			}
			// Check if an error is expected
			if testCase.expectedIsError == true {
				t.Errorf("expected error")
			}
		})
	}
}

// TestCalculateBackoff tests calculateBackoff without random jitter
func TestCalculateBackoff(t *testing.T) {
	testCases := []struct {
		name            string
		inputPolicy     *BackoffPolicy
		expectedBackoff []time.Duration
	}{
		{ // Base test
			name: "base",
			inputPolicy: &BackoffPolicy{
				MaxRetries:        0,
				MaxBackoff:        6000,
				BackoffMultiplier: 2,
				MaxRandomJitter:   0,
				InitialDelay:      500,
				IsLogging:         false,
			},
			expectedBackoff: []time.Duration{500, 1000, 2000, 4000, 6000, 6000, 6000},
		},
		{
			name: "initialDelay",
			inputPolicy: &BackoffPolicy{
				MaxRetries:        0,
				MaxBackoff:        16000,
				BackoffMultiplier: 2,
				MaxRandomJitter:   0,
				InitialDelay:      1000,
				IsLogging:         false,
			},
			expectedBackoff: []time.Duration{1000, 2000, 4000, 8000, 16000, 16000, 16000},
		},
		{
			name: "backoffMultiplier",
			inputPolicy: &BackoffPolicy{
				MaxRetries:        0,
				MaxBackoff:        27000,
				BackoffMultiplier: 3,
				MaxRandomJitter:   0,
				InitialDelay:      1000,
				IsLogging:         false,
			},
			expectedBackoff: []time.Duration{1000, 3000, 9000, 27000, 27000, 27000, 27000},
		},

		{
			name: "no backoffMultiplier",
			inputPolicy: &BackoffPolicy{
				MaxRetries:        0,
				MaxBackoff:        27000,
				BackoffMultiplier: 1,
				MaxRandomJitter:   0,
				InitialDelay:      1000,
				IsLogging:         false,
			},
			expectedBackoff: []time.Duration{1000, 1000, 1000, 1000, 1000, 1000, 1000},
		},
	}

	for _, testCase := range testCases {
		var backoffGrowthRate int32 = 1

		testName := fmt.Sprintf("%s test", testCase.name)

		t.Run(testName, func(t *testing.T) {
			for _, expectedBackoff := range testCase.expectedBackoff {
				expectedBackoff = expectedBackoff * time.Millisecond

				actualBackoff := calculateBackoff(backoffGrowthRate, testCase.inputPolicy)
				if actualBackoff != expectedBackoff {
					t.Errorf("got: %d, expected: %d", actualBackoff/time.Millisecond, expectedBackoff/time.Millisecond)
				}
				backoffGrowthRate *= testCase.inputPolicy.BackoffMultiplier
			}
		})
	}
}

// TestCalculateBackoff2 tests calculateBackoff with random jitter
func TestCalculateBackoff2(t *testing.T) {
	inputPolicy := &BackoffPolicy{
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
		for backoffNumber, expectedValue := range expectedBackoffRange {
			expectedBackoffRange[backoffNumber] = expectedValue * time.Millisecond
		}
	}

	var backoffGrowthRate int32 = 1

	for _, expectedBackoff := range expectedJitterResults {
		actualBackoff := calculateBackoff(backoffGrowthRate, inputPolicy)
		if actualBackoff < expectedBackoff[0] && actualBackoff > expectedBackoff[1] {
			t.Errorf("got: %d, expected: %d", actualBackoff, expectedBackoff)
		}

		backoffGrowthRate *= inputPolicy.BackoffMultiplier
	}
}

func TestIsLogging(t *testing.T) {
	testCases := []struct {
		name       string
		testPolicy *BackoffPolicy
	}{
		{
			name: "with logging",
			testPolicy: &BackoffPolicy{
				MaxRetries:        10,
				MaxBackoff:        6000,
				BackoffMultiplier: 2,
				MaxRandomJitter:   0,
				InitialDelay:      500,
				IsLogging:         true,
			},
		},
		{
			name: "without logging",
			testPolicy: &BackoffPolicy{
				MaxRetries:        10,
				MaxBackoff:        6000,
				BackoffMultiplier: 2,
				MaxRandomJitter:   0,
				InitialDelay:      500,
				IsLogging:         false,
			},
		},
	}

	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	for _, testCase := range testCases {
		testName := fmt.Sprintf("%s test", testCase.name)

		t.Run(testName, func(t *testing.T) {
			buf := new(bytes.Buffer)
			log.SetOutput(buf)

			defer func() {
				log.SetOutput(os.Stderr)
			}()
			isLogging(testCase.testPolicy, "%d %d %d", 1, 2, 3)
			result := buf.String()

			if testCase.testPolicy.IsLogging == true && result != "1 2 3\n" {
				t.Errorf("expected (1 2 3) got %s", result)
				return
			}

			if testCase.testPolicy.IsLogging == false && result != "" {
				t.Errorf("expected () got %s", result)
			}
		})
	}
}

func BenchmarkRetry(b *testing.B) {
	benchmarkPolicy := &BackoffPolicy{
		MaxRetries:        10,
		MaxBackoff:        0,
		BackoffMultiplier: 1,
		MaxRandomJitter:   0,
		InitialDelay:      0,
		IsLogging:         false,
	}

	for i := 0; i < b.N; i++ {
		retryableFunction := func() error {
			return nil
		}

		_ = Retry(retryableFunction, benchmarkPolicy)
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
	myPolicy, err := NewCustomBackoffPolicy(10, 0, 2, 1000, 1000, false)
	if err != nil {
		// Handle error
		return
	}

	// A function that may fail returning an error
	retryableFunction := func() error { return nil }

	if err := Retry(retryableFunction, myPolicy); err != nil {
		// Handle error
		return
	}
	// Output:
}
