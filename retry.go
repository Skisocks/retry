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
	var backoffGrowthRate int32 = 1
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
func calculateBackoff(backoffGrowthRate int32, policy *BackoffPolicy) time.Duration {
	var backoff time.Duration

	// Add random jitter to the backoff time
	if policy.MaxRandomJitter == 0 {
		backoff = time.Duration(policy.InitialDelay*backoffGrowthRate) * time.Millisecond
	} else {
		backoff = time.Duration((rand.Int31n(policy.MaxRandomJitter)+policy.InitialDelay)*backoffGrowthRate) * time.Millisecond
	}

	// Limit backoff to the maximum value set in config
	maxBackoff := time.Duration(policy.MaxBackoff) * time.Millisecond
	if backoff > maxBackoff && maxBackoff != 0 {
		backoff = maxBackoff
	}

	isLogging(policy, "backoff: %d", backoff/time.Millisecond)
	return backoff
}

// isLogging logs format with parameters if logging has been selected in the policy
func isLogging(policy *BackoffPolicy, format string, params ...interface{}) {
	if policy.IsLogging == true {
		log.Printf(format, params...)
	}
}
