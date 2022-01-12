// Package retry implements exponential calculateBackoff algorithms to successfully retry functions.
//
// Use Retry function with a NewCustomBackoffPolicy to re-execute any retryable function with custom parameters.
// Alternatively you can use NewBackoffPolicy to re-execute with generic parameters for ease of use.
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
// If it does not succeed before BackoffPolicy.MaxRetries is reached then a maxRetryError is returned.
func Retry(function function, policy *BackoffPolicy) error {
	retryAttempt := 1
	var backoffGrowthRate float32 = 1
	rand.Seed(time.Now().Unix())

	for {
		// If the function is successful return error is nil
		if err := function(); err == nil {
			if retryAttempt == 1 {
				return nil
			}
			isLogging(policy, "function was successful on attempt: %d\n", retryAttempt)
			return nil
		}

		// If the function is not successful within MaxRetries return maxRetryError
		if retryAttempt == policy.MaxRetries {
			return &maxRetryError{maxRetries: policy.MaxRetries}
		}
		isLogging(policy, "function was unsuccessful on attempt: %d\n", retryAttempt)

		// Sleep
		time.Sleep(calculateBackoff(backoffGrowthRate, policy))

		// Increase the calculateBackoff & increment the retry counter
		backoffGrowthRate *= policy.BackoffMultiplier
		retryAttempt++
	}
}

// calculateBackoff returns the next backoff interval that Retry should sleep for depending on the policy variables
func calculateBackoff(backoffGrowthRate float32, policy *BackoffPolicy) time.Duration {
	var currentBackoff time.Duration

	// Convert int to float32
	initialDelay := float32(policy.InitialDelay)
	maxRandomJitter := rand.Float32() * float32(policy.MaxRandomJitter)

	// Add random jitter to the currentBackoff time
	if maxRandomJitter == 0 {
		currentBackoff = time.Duration(initialDelay*backoffGrowthRate) * time.Millisecond
	} else {
		currentBackoff = time.Duration((maxRandomJitter+initialDelay)*backoffGrowthRate) * time.Millisecond
	}

	// Limit currentBackoff to the maximum value set in config
	maxBackoff := time.Duration(policy.MaxBackoff) * time.Millisecond
	if currentBackoff > maxBackoff && maxBackoff != 0 {
		currentBackoff = maxBackoff
	}

	isLogging(policy, "currentBackoff: %d", currentBackoff/time.Millisecond)
	return currentBackoff
}

// isLogging logs format with parameters if logging has been selected in the policy
func isLogging(policy *BackoffPolicy, format string, params ...interface{}) {
	if policy.IsLogging {
		log.Printf(format, params...)
	}
}
