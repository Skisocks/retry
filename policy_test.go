package retry

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNewCustomBackoffPolicy(t *testing.T) {
	testCases := []struct {
		name string

		inputMaxRetries        int
		inputMaxBackoff        int
		inputBackoffMultiplier float32
		inputMaxRandomJitter   int32
		inputInitialDelay      int32
		inputIsLogging         bool

		expectedPolicy *BackoffPolicy
		errIsExpected  bool
	}{
		{
			name:            "base",
			inputMaxRetries: 10, inputMaxBackoff: 6000, inputBackoffMultiplier: 2, inputMaxRandomJitter: 500, inputInitialDelay: 500,
			expectedPolicy: &BackoffPolicy{
				MaxRetries:        10,
				MaxBackoff:        6000,
				BackoffMultiplier: 2,
				MaxRandomJitter:   500,
				InitialDelay:      500,
				IsLogging:         false,
			},
		},
		{
			name:            "neg MaxRetries",
			inputMaxRetries: -1, inputMaxBackoff: 6000, inputBackoffMultiplier: 2, inputMaxRandomJitter: 500, inputInitialDelay: 500,
			errIsExpected: true,
		},
		{
			name:            "neg maxBackoff",
			inputMaxRetries: 10, inputMaxBackoff: -1, inputBackoffMultiplier: 2, inputMaxRandomJitter: 500, inputInitialDelay: 500,
			errIsExpected: true,
		},
		{
			name:            "neg BackoffMultiplier",
			inputMaxRetries: 10, inputMaxBackoff: 6000, inputBackoffMultiplier: -1, inputMaxRandomJitter: 500, inputInitialDelay: 500,
			errIsExpected: true,
		},
		{
			name:            "neg RandomJitter",
			inputMaxRetries: 10, inputMaxBackoff: 6000, inputBackoffMultiplier: 2, inputMaxRandomJitter: -1, inputInitialDelay: 500,
			errIsExpected: true,
		},
		{
			name:            "neg InitialDelay",
			inputMaxRetries: 10, inputMaxBackoff: 6000, inputBackoffMultiplier: 2, inputMaxRandomJitter: 500, inputInitialDelay: -1,
			errIsExpected: true,
		},
	}

	for _, testCase := range testCases {
		testName := fmt.Sprintf("%s test", testCase.name)

		t.Run(testName, func(t *testing.T) {
			actualTestPolicy, err := NewCustomBackoffPolicy(
				testCase.inputMaxRetries,
				testCase.inputMaxBackoff,
				testCase.inputBackoffMultiplier,
				testCase.inputMaxRandomJitter,
				testCase.inputInitialDelay,
				testCase.inputIsLogging,
			)
			if err != nil {
				if !testCase.errIsExpected {
					t.Errorf("got: %+v, expected: No Error", err)
				}
				return
			}

			if !reflect.DeepEqual(actualTestPolicy, testCase.expectedPolicy) {
				t.Errorf("got: %+v, expected: %+v", *actualTestPolicy, *testCase.expectedPolicy)
			}
		})
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
	if !reflect.DeepEqual(actualTestPolicy, expectedTestPolicy) {
		t.Errorf("got: %+v, expected: %+v", *actualTestPolicy, *expectedTestPolicy)
	}
}
