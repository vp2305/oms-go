package common

import "syscall"

func EnvString(key, fallback string) string {
	if val, exists := syscall.Getenv(key); exists {
		return val
	}
	return fallback
}
