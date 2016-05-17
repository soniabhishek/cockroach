package config

type ConfigNotFoundError struct {
	ConfigurationKey configKey
}

func (e *ConfigNotFoundError) Error() string {
	return e.ConfigurationKey.String() + " configuration Value not found"
}

func newConfigurationError(c configKey) error {
	return &ConfigNotFoundError{c}
}
