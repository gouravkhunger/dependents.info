package utils

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRetryWithBackoff_SuccessFirst(t *testing.T) {
	calls := 0
	err := RetryWithBackoff(context.Background(), 3, 10*time.Millisecond, func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if calls != 1 {
		t.Errorf("expected 1 call, got %d", calls)
	}
}

func TestRetryWithBackoff_SuccessAfterRetries(t *testing.T) {
	calls := 0
	err := RetryWithBackoff(context.Background(), 3, 10*time.Millisecond, func() error {
		calls++
		if calls < 3 {
			return errors.New("transient")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if calls != 3 {
		t.Errorf("expected 3 calls, got %d", calls)
	}
}

func TestRetryWithBackoff_AllFail(t *testing.T) {
	calls := 0
	want := errors.New("persistent error")
	err := RetryWithBackoff(context.Background(), 2, 10*time.Millisecond, func() error {
		calls++
		return want
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != want.Error() {
		t.Errorf("expected %q, got %q", want, err)
	}
	if calls != 3 {
		t.Errorf("expected 3 calls (initial + 2 retries), got %d", calls)
	}
}

func TestRetryWithBackoff_PermanentError(t *testing.T) {
	calls := 0
	err := RetryWithBackoff(context.Background(), 3, 10*time.Millisecond, func() error {
		calls++
		return PermanentError{Err: errors.New("not found")}
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if calls != 1 {
		t.Errorf("expected 1 call, got %d", calls)
	}
	var perm PermanentError
	if !errors.As(err, &perm) {
		t.Errorf("expected PermanentError, got %T", err)
	}
}

func TestRetryWithBackoff_ZeroRetries(t *testing.T) {
	calls := 0
	err := RetryWithBackoff(context.Background(), 0, 10*time.Millisecond, func() error {
		calls++
		return errors.New("fail")
	})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if calls != 1 {
		t.Errorf("expected 1 call, got %d", calls)
	}
}

func TestRetryWithBackoff_ContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	calls := 0
	err := RetryWithBackoff(ctx, 5, 50*time.Millisecond, func() error {
		calls++
		if calls == 1 {
			cancel()
		}
		return errors.New("transient")
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if calls != 1 {
		t.Errorf("expected 1 call, got %d", calls)
	}
}
