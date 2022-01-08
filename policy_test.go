package retry

import (
	"reflect"
	"testing"
)

func TestNewCustomBackoffPolicy(t *testing.T) {
	expectedTestPolicy := &BackoffPolicy{
		MaxRetries:        10,
		MaxBackoff:        6000,
		BackoffMultiplier: 2,
		MaxRandomJitter:   500,
		InitialDelay:      500,
		IsLogging:         false,
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
