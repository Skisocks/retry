package retry

/*
Policy contains the back off parameters that this package implements.

MaxRetries: The maximum number of attempts that Retry will make.
MaxBackoff: The maximum amount of time in milliseconds that Retry will backoff for.
BackoffMultiplier: The multiplier used to increase the backoff delay exponentially.
MaxRandomJitter: The maximum value of random jitter added to the delay in milliseconds. 0 is no jitter.
InitialDelay: The initial delay in milliseconds/
*/
type Policy struct {
	MaxRetries        int   // Maximum number of attempts
	MaxBackoff        int   // Maximum backoff time in milliseconds. 0 is no maximum backoff
	BackoffMultiplier int32 // Multiplier added to delay between attempts
	MaxRandomJitter   int32 // Maximum value of random jitter added to the delay in milliseconds. 0 is no jitter.
	InitialDelay      int32 // Delay in milliseconds
}

// NewPolicy returns a backoff Policy
func NewPolicy(
	maxRetries int,
	maxBackoff int,
	backoffMultiplier int32,
	maxRandomJitter int32,
	initialDelay int32,
) *Policy {
	return &Policy{
		MaxRetries:        maxRetries,
		MaxBackoff:        maxBackoff,
		BackoffMultiplier: backoffMultiplier,
		MaxRandomJitter:   maxRandomJitter,
		InitialDelay:      initialDelay,
	}
}
