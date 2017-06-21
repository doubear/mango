package mango

import (
	"os"
	"sync"
)

type environment struct {
	values map[string]string
	mutex  sync.RWMutex
}

var (
	env = new(environment)
)

//Env retrieves value from environment variables.
func Env(n string, d string) string {
	env.mutex.RLock()
	defer env.mutex.RUnlock()

	if v, ok := env.values[n]; ok {
		return v
	}

	if v := os.Getenv(n); v != "" {
		return v
	}

	return d
}

//SetEnv stores value to env.
func SetEnv(n, v string) {
	env.mutex.Lock()
	defer env.mutex.Unlock()

	env.values[n] = v
}
