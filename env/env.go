package env

import "sync"

const (
	WORKSPACE       = "workspace"
	SERVICE_CONTEXT = "service_context"
)

var envVars = map[string]any{}

var mutex = sync.Mutex{}

func GetEnvVar[T any](key string) T {
	val := envVars[key]
	if val == nil {
		var zero T
		return zero
	}
	return val.(T)
}
func SetEnvVar[T any](key string, value T) T {
	mutex.Lock()
	oldValue := envVars[key]
	envVars[key] = value
	mutex.Unlock()
	if oldValue == nil {
		var zero T
		return zero
	}
	return oldValue.(T)
}
