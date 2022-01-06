package retry

/*
BackoffPolicy contains the back off parameters that this package implements.

MaxRetries: The maximum number of attempts that Retry will make.
MaxBackoff: The maximum amount of time in milliseconds that Retry will backoff for.
BackoffMultiplier: The multiplier used to increase the backoff delay exponentially.
MaxRandomJitter: The maximum value of random jitter added to the delay in milliseconds. 0 is no jitter.
InitialDelay: The initial delay in milliseconds/
*/
type BackoffPolicy struct {
	MaxRetries        int   // Maximum number of attempts
	MaxBackoff        int   // Maximum backoff time in milliseconds. 0 is no maximum backoff
	BackoffMultiplier int32 // Multiplier added to delay between attempts
	MaxRandomJitter   int32 // Maximum value of random jitter added to the delay in milliseconds. 0 is no jitter.
	InitialDelay      int32 // Delay in milliseconds
}

const (
	DefaultMaxRetries        int   = 5
	DefaultMaxBackoff        int   = 0
	DefaultBackoffMultiplier int32 = 2
	DefaultMaxRandomJitter   int32 = 1000
	DefaultInitialDelay      int32 = 1000
)

// NewBackoffPolicy returns a BackoffPolicy with default parameters
func NewBackoffPolicy() *BackoffPolicy {
	return &BackoffPolicy{
		MaxRetries:        DefaultMaxRetries,
		MaxBackoff:        DefaultMaxBackoff,
		BackoffMultiplier: DefaultBackoffMultiplier,
		MaxRandomJitter:   DefaultMaxRandomJitter,
		InitialDelay:      DefaultInitialDelay,
	}
}

func NewCustomBackoffPolicy(
	maxRetries int,
	maxBackoff int,
	backoffMultiplier int32,
	maxRandomJitter int32,
	initialDelay int32,
) *BackoffPolicy {
	return &BackoffPolicy{
		MaxRetries:        maxRetries,
		MaxBackoff:        maxBackoff,
		BackoffMultiplier: backoffMultiplier,
		MaxRandomJitter:   maxRandomJitter,
		InitialDelay:      initialDelay,
	}
}
