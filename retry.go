// Package retry implements exponential backoff algorithms to successfully retry functions.
//
// Use Retry function with a NewCustomBackoffPolicy to re-execute any retryable function with default parameters.
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
			log.Printf("function was sucsessful on attempt: %d\n", retryAttempt)
			return nil
		}

		// If the function is not successful within MaxRetries return maxRetryError
		if retryAttempt == policy.MaxRetries {
			return &maxRetryError{maxRetries: policy.MaxRetries}
		}
		log.Printf("function was unsuccessful on attempt: %d\n", retryAttempt)

		// Sleep
		backoff(backoffGrowthRate, policy)

		// Increase the backoff & increment the retry counter
		backoffGrowthRate *= policy.BackoffMultiplier
		retryAttempt++
	}
}

// backoff causes Retry to sleep for a period depending on the config settings
func backoff(backoffGrowthRate int32, cfg *BackoffPolicy) {
	var backoff time.Duration

	// Add random jitter to the backoff time
	if cfg.MaxRandomJitter == 0 {
		backoff = time.Duration(cfg.InitialDelay*backoffGrowthRate) * time.Millisecond
	} else {
		backoff = time.Duration((rand.Int31n(cfg.MaxRandomJitter)+cfg.InitialDelay)*backoffGrowthRate) * time.Millisecond
	}

	// Limit backoff to the maximum value set in config
	maxBackoff := time.Duration(cfg.MaxBackoff) * time.Millisecond
	if backoff > maxBackoff && maxBackoff != 0 {
		backoff = maxBackoff
	}

	log.Printf("backoff: %d", backoff/time.Millisecond)
	time.Sleep(backoff)
}
