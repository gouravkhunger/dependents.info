package utils

import (
	"errors"
	"math/rand/v2"
	"time"
)

type PermanentError struct {
	Err error
}

func (e PermanentError) Error() string { return e.Err.Error() }
func (e PermanentError) Unwrap() error { return e.Err }

func RetryWithBackoff(maxRetries int, base time.Duration, fn func() error) error {
	var lastErr error
	for attempt := range maxRetries + 1 {
		lastErr = fn()
		if lastErr == nil {
			return nil
		}
		var perm PermanentError
		if errors.As(lastErr, &perm) {
			return lastErr
		}
		if attempt < maxRetries {
			backoff := base * time.Duration(1<<uint(attempt))
			jitter := time.Duration(rand.Int64N(int64(backoff/2) + 1))
			time.Sleep(backoff + jitter)
		}
	}
	return lastErr
}
