package config

import "os"

type Environment struct {
	vars map[string]string
}

// InitializeWithDefaults will search for the environment variables with the associated
// keys from the vars map parameter. If that environment variable doesn't exist, it will
// populate with the preset default from the vars map
// Returns Environment struct with populated Vars map
func InitializeEnvironmentWithDefaults(vars map[string]string) Environment {
	env := Environment{
		vars: map[string]string{},
	}

	for key, val := range vars {
		env.vars[key] = val
		if v := os.Getenv(key); v != "" {
			env.vars[key] = v;
		}
	}

	return env
}

// GetVariable retrieves the stored variable for the associated key
// Returns an empty string if it does not exist
func (e Environment) GetVariable(key string) string {
	if val, ok := e.vars[key]; ok {
		return val
	}

	return ""
}