package retry

import (
	"log"
	"math/rand"
	"time"
)
type Function func() error

var maxRetries int
var backoffMultiplier int32

// Retry
// set maxRetries to 0 for infinite retries
func Retry(retryableFunction Function, maxRetries int) error {
	var retryAttempt = 1
	var backoffMultiplier int32 = 1
	rand.Seed(time.Now().Unix())

	for {
		// If the function is successful return error is nil
		if err := retryableFunction(); err == nil {
			if retryAttempt == 1 {
				return nil
			} else {
				log.Printf("function was sucsessful on attempt: %d\n", retryAttempt)
				return nil
			}
		}

		// Todo: Check if the error is retryable

		// If the function is not successful within maxRetries return MaxRetryError
		if retryAttempt == maxRetries {
			return &MaxRetryError{maxRetries: maxRetries}
		}

		log.Printf("function was unsuccessful on attempt: %d\n", retryAttempt)

		// Sleep for a random time between 1-2sec
		backoff := time.Duration((rand.Int31n(1000) + 1000 ) * backoffMultiplier) * time.Millisecond
		time.Sleep(backoff)

		// Exponentially increase the backoff & increment the retry counter
		backoffMultiplier *= 2
		retryAttempt++
	}
}

