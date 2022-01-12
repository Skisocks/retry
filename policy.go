package retry

// BackoffPolicy is a backoff implementation that exponentially increases the delay between retry attempts
// whilst adding random jitter to each delay. The backoff time is calculated using the formula:
// Backoff interval =
//		(RandomJitterRange + InitialDelay) * BackoffGrowthRate
// If no random jitter is required then the formula is:
// Backoff interval =
// 		InitialDelay * BackoffGrowthRate
type BackoffPolicy struct {
	MaxRetries        int     // Maximum number of attempts that Retry will make.
	MaxBackoff        int     // Maximum amount of time in milliseconds that Retry will backoff for.
	BackoffMultiplier float32 // Multiplier used to increase the backoff delay exponentially.
	MaxRandomJitter   int32   // Maximum value of random jitter added to the delay in milliseconds. 0 is no jitter.
	InitialDelay      int32   // Initial delay in milliseconds
	IsLogging         bool    // Defines whether logging should occur during retries
}

const (
	DefaultMaxRetries        int     = 0
	DefaultMaxBackoff        int     = 0
	DefaultBackoffMultiplier float32 = 2
	DefaultMaxRandomJitter   int32   = 1000
	DefaultInitialDelay      int32   = 500
	DefaultIsLogging         bool    = false
)

// NewBackoffPolicy returns a BackoffPolicy with default parameters
func NewBackoffPolicy() *BackoffPolicy {
	return &BackoffPolicy{
		MaxRetries:        DefaultMaxRetries,
		MaxBackoff:        DefaultMaxBackoff,
		BackoffMultiplier: DefaultBackoffMultiplier,
		MaxRandomJitter:   DefaultMaxRandomJitter,
		InitialDelay:      DefaultInitialDelay,
		IsLogging:         DefaultIsLogging,
	}
}

// NewCustomBackoffPolicy returns a BackoffPolicy with custom parameters
func NewCustomBackoffPolicy(
	maxRetries int,
	maxBackoff int,
	backoffMultiplier float32,
	maxRandomJitter int32,
	initialDelay int32,
	isLogging bool,
) (*BackoffPolicy, error) {
	if maxRetries < 0 {
		return nil, &inputError{err: "maxRetries cannot be negative"}
	}
	if maxBackoff < 0 {
		return nil, &inputError{err: "maxBackoff cannot be negative"}
	}
	if backoffMultiplier <= 0 {
		return nil, &inputError{err: "backoff multiplier must be a positive integer"}
	}
	if maxRandomJitter < 0 {
		return nil, &inputError{err: "maxRandomJitter cannot be negative"}
	}
	if initialDelay < 0 {
		return nil, &inputError{err: "initialDelay cannot be negative"}
	}

	return &BackoffPolicy{
		MaxRetries:        maxRetries,
		MaxBackoff:        maxBackoff,
		BackoffMultiplier: backoffMultiplier,
		MaxRandomJitter:   maxRandomJitter,
		InitialDelay:      initialDelay,
		IsLogging:         isLogging,
	}, nil
}
