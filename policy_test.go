package retry

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNewCustomBackoffPolicy(t *testing.T) {
	testCases := []struct {
		testName               string
		inputMaxRetries        int
		inputMaxBackoff        int
		inputBackoffMultiplier int32
		inputMaxRandomJitter   int32
		inputInitialDelay      int32
		inputIsLogging         bool
		expectedPolicy         *BackoffPolicy
		errIsExpected          bool
	}{
		{
			"base",
			10, 6000, 2, 500, 500, false,
			&BackoffPolicy{
				MaxRetries:        10,
				MaxBackoff:        6000,
				BackoffMultiplier: 2,
				MaxRandomJitter:   500,
				InitialDelay:      500,
				IsLogging:         false,
			},
			false,
		},
		{
			"neg MaxRetries",
			-1, 6000, 2, 500, 500, false,
			nil,
			true,
		},
		{
			"neg maxBackoff",
			10, -1, 2, 500, 500, false,
			nil,
			true,
		},
		{
			"neg BackoffMultiplier",
			10, 6000, -1, 500, 500, false,
			nil,
			true,
		},
		{
			"neg RandomJitter",
			10, 6000, 2, -1, 500, false,
			nil,
			true,
		},
		{
			"neg InitialDelay",
			10, 6000, 2, 500, -1, false,
			nil,
			true,
		},
	}

	for _, testCase := range testCases {
		testName := fmt.Sprintf("%s test", testCase.testName)

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
				if testCase.errIsExpected == false {
					t.Errorf("got: %+v, expected: No Error", err)
				}
				return
			}

			if reflect.DeepEqual(actualTestPolicy, testCase.expectedPolicy) == false {
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
	if reflect.DeepEqual(actualTestPolicy, expectedTestPolicy) == false {
		t.Errorf("got: %+v, expected: %+v", *actualTestPolicy, *expectedTestPolicy)
	}
}
