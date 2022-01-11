// Package retry implements exponential calculateBackoff algorithms to successfully retry functions.
//
// Use Retry function with a NewCustomPolicy to re-execute any retryable function with custom parameters.
// Alternatively you can use NewPolicy to re-execute with generic parameters for ease of use.
//
// This package currently cannot be used with channels.
package retry

import (
	"log"
	"math/rand"
	"time"
)

type function func() error

// Retry calls a function and re-executes it if it fails.
// If it does not succeed before settings.MaxRetries is reached then a maxRetryError is returned.
func Retry(function function, policy *Policy) error {
	retryAttempt := 1
	var backoffGrowthRate float32 = 1
	rand.Seed(time.Now().Unix())

	for {
		// If the function is successful return error is nil
		if err := function(); err == nil {
			if retryAttempt == 1 {
				return nil
			}
			isLogging(*policy.Settings, "function was successful on attempt: %d\n", retryAttempt)
			return nil
		}

		// If the function is not successful within MaxRetries return maxRetryError
		if retryAttempt == policy.Settings.MaxRetries {
			return &maxRetryError{maxRetries: policy.Settings.MaxRetries}
		}
		isLogging(*policy.Settings, "function was unsuccessful on attempt: %d\n", retryAttempt)

		// Sleep
		time.Sleep(calculateBackoff(backoffGrowthRate, *policy.Settings))

		// Increase the calculateBackoff & increment the retry counter
		backoffGrowthRate *= policy.Settings.BackoffMultiplier
		retryAttempt++
	}
}

// calculateBackoff returns the next backoff interval that Retry should sleep for depending on the policy settings
func calculateBackoff(backoffGrowthRate float32, settings settings) time.Duration {
	var currentBackoff time.Duration

	// Convert int to float32
	initialDelay := float32(settings.InitialDelay)
	maxRandomJitter := rand.Float32() * float32(settings.MaxRandomJitter)

	// Add random jitter to the currentBackoff time
	if maxRandomJitter == 0 {
		currentBackoff = time.Duration(initialDelay*backoffGrowthRate) * time.Millisecond
	} else {
		currentBackoff = time.Duration((maxRandomJitter+initialDelay)*backoffGrowthRate) * time.Millisecond
	}

	// Limit currentBackoff to the maximum value set in config
	maxBackoff := time.Duration(settings.MaxBackoff) * time.Millisecond
	if currentBackoff > maxBackoff && maxBackoff != 0 {
		currentBackoff = maxBackoff
	}

	isLogging(settings, "currentBackoff: %d", currentBackoff/time.Millisecond)
	return currentBackoff
}

// isLogging logs format with parameters if logging has been selected in the policy
func isLogging(policy settings, format string, params ...interface{}) {
	if policy.IsLogging {
		log.Printf(format, params...)
	}
}
