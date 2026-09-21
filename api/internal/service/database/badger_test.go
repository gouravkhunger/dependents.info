package database

import (
	"errors"
	"testing"
	"time"
)

func TestIsLockErr(t *testing.T) {
	if isLockErr(nil) {
		t.Fatal("nil")
	}
	if isLockErr(errors.New("no such file")) {
		t.Fatal("other error")
	}
	err := errors.New(`Cannot acquire directory lock on "/tmp/x".  Another process is using this Badger database.`)
	if !isLockErr(err) {
		t.Fatal("lock error")
	}
}

func TestNewBadgerService_WaitsForLock(t *testing.T) {
	origA, origW := openAttempts, openWait
	openAttempts, openWait = 40, 50*time.Millisecond
	defer func() {
		openAttempts, openWait = origA, origW
	}()

	dir := t.TempDir()
	first := NewBadgerService(dir)

	got := make(chan *BadgerService, 1)
	go func() {
		got <- NewBadgerService(dir)
	}()

	time.Sleep(150 * time.Millisecond)
	if err := first.Close(); err != nil {
		t.Fatal(err)
	}

	select {
	case second := <-got:
		if err := second.Close(); err != nil {
			t.Fatal(err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("second open did not succeed after lock released")
	}
}
