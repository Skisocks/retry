package retry

type Config struct {
	maxRetries        int   // Maximum number of attempts
	maxBackoff        int   // Maximum backoff time in milliseconds. 0 is no maximum backoff
	backoffMultiplier int32 // Multiplier added to delay between attempts
	maxRandomJitter   int32 // Maximum value of random jitter added to the delay in milliseconds. 0 is no jitter.
	initialDelay      int32 // Delay in milliseconds
}

func NewConfig(
	maxRetries int,
	maxBackoff int,
	backoffMultiplier int32,
	maxRandomJitter int32,
	initialDelay int32,
) *Config {
	return &Config{
		maxRetries:        maxRetries,
		maxBackoff:        maxBackoff,
		backoffMultiplier: backoffMultiplier,
		maxRandomJitter:   maxRandomJitter,
		initialDelay:      initialDelay,
	}
}
