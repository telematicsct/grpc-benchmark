package env

import (
	"os"
)

//GetString returns the value by looking up the environment.
//Returns the fallback string if the key doesn't exist
func GetString(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
