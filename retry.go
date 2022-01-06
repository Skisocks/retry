package retry

import (
	"log"
	"math/rand"
	"time"
)

type function func() error

// Retry calls a function and re-executes it if it fails
func Retry(function function, policy *Policy) error {
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

		// If the function is not successful within maxRetries return maxRetryError
		if retryAttempt == policy.maxRetries {
			return &maxRetryError{maxRetries: policy.maxRetries}
		}
		log.Printf("function was unsuccessful on attempt: %d\n", retryAttempt)

		// Sleep
		backoff(backoffGrowthRate, policy)

		// Exponentially increase the backoff & increment the retry counter
		backoffGrowthRate *= policy.backoffMultiplier
		retryAttempt++
	}
}

// backoff causes the Retry function to sleep for a period depending on the config settings
func backoff(backoffMultiplier int32, cfg *Policy) {
	var backoff time.Duration

	// Add random jitter to the backoff time
	if cfg.maxRandomJitter == 0 {
		backoff = time.Duration(cfg.initialDelay*backoffMultiplier) * time.Millisecond
	} else {
		backoff = time.Duration((rand.Int31n(cfg.maxRandomJitter)+cfg.initialDelay)*backoffMultiplier) * time.Millisecond
	}

	// Limit backoff to the maximum value set in config
	maxBackoff := time.Duration(cfg.maxBackoff) * time.Millisecond
	if backoff > maxBackoff && maxBackoff != 0 {
		backoff = maxBackoff
	}

	log.Printf("backoff: %d", backoff/time.Millisecond)
	time.Sleep(backoff)
}
