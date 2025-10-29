package lock

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type Lock struct {
	Path string
	file *os.File
}

// Acquire creates and locks a PID file to prevent concurrent runs
func Acquire(lockPath string) (*Lock, error) {
	dir := filepath.Dir(lockPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create lock directory: %w", err)
	}

	// Check if lock already exists
	if _, err := os.Stat(lockPath); err == nil {
		// Lock exists, check if process is still running
		pidBytes, err := os.ReadFile(lockPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read lock file: %w", err)
		}
		pid, err := strconv.Atoi(string(pidBytes))
		if err == nil {
			// Check if process is running
			proc, err := os.FindProcess(pid)
			if err == nil {
				err = proc.Signal(os.Signal(nil)) // Signal(0) checks if process exists
				if err == nil {
					return nil, fmt.Errorf("another instance is already running (PID %d)", pid)
				}
			}
		}
		// Stale lock, remove it
		os.Remove(lockPath)
	}

	// Create lock file with current PID
	file, err := os.OpenFile(lockPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire lock: %w", err)
	}

	pid := os.Getpid()
	if _, err := file.WriteString(fmt.Sprintf("%d\n", pid)); err != nil {
		file.Close()
		os.Remove(lockPath)
		return nil, fmt.Errorf("failed to write PID to lock: %w", err)
	}
	if err := file.Sync(); err != nil {
		file.Close()
		os.Remove(lockPath)
		return nil, fmt.Errorf("failed to sync lock: %w", err)
	}

	return &Lock{Path: lockPath, file: file}, nil
}

// Release removes the lock file
func (l *Lock) Release() error {
	if l == nil {
		return nil
	}
	if l.file != nil {
		l.file.Close()
	}
	if err := os.Remove(l.Path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to release lock: %w", err)
	}
	return nil
}

// LockInfo returns the PID in the lock file, or 0 if not found
func LockInfo(lockPath string) (int, time.Time, error) {
	info, err := os.Stat(lockPath)
	if err != nil {
		return 0, time.Time{}, err
	}

	pidBytes, err := os.ReadFile(lockPath)
	if err != nil {
		return 0, info.ModTime(), err
	}

	pid, err := strconv.Atoi(string(pidBytes))
	if err != nil {
		return 0, info.ModTime(), nil
	}

	return pid, info.ModTime(), nil
}
