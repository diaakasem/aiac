package libaiac

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/adrg/xdg"
)

// BackendType is a const type used for identifying backends, a.k.a LLM providers.
type BackendType string

const (
	// BackendOpenAI represents the OpenAI LLM provider.
	BackendOpenAI BackendType = "openai"

	// BackendBedrock represents the Amazon Bedrock LLM provider.
	BackendBedrock BackendType = "bedrock"

	// BackendOllama represents the Ollama LLM provider.
	BackendOllama BackendType = "ollama"
)

// Config holds the configuration for aiac.
type Config struct {
	// Backends is the map of named backends that can be used to generate
	// IaC templates.
	Backends map[string]BackendConfig `toml:"backends"`

	// DefaultBackend is the name of the default backend to use when one is
	// not specifically selected.
	DefaultBackend string `toml:"default_backend"`
}

// BackendConfig holds backend-specific configuration.
type BackendConfig struct {
	// Type is the type of the backend (generally the name of an LLM provider)
	Type BackendType `toml:"type"`

	// AWSProfile is used by Amazon Bedrock. It is the name of the AWS profile
	// in the credentials file to use.
	AWSProfile string `toml:"aws_profile"`

	// AWSRegion is used by Amazon Bedrock. It is the name of the region where
	// the models to use are hosted.
	AWSRegion string `toml:"aws_region"`

	// APIKey is an API key used for authentication. It is used by backends such
	// as OpenAI.
	APIKey string `toml:"api_key"`

	// APIVersion allows setting a specific API version to use. It is accepted
	// by the OpenAI backend.
	APIVersion string `toml:"api_version"`

	// URL allows setting a custom URL for a backend's API. It is accepted by
	// backends such as OpenAI and Ollama.
	URL string `toml:"url"`

	// DefaultModel is the name of the model to use by default when a specific
	// one is not selected.
	DefaultModel string `toml:"default_model"`

	// ExtraHeaders allows setting extra HTTP headers whenever aiac sends
	// requests to the backend. Bedrock backends do not support this setting.
	ExtraHeaders map[string]string `toml:"extra_headers"`
}

// LoadConfig loads an aiac configuration file from the provided path, which
// must be a TOML file. If path is an empty string, the default path will be
// checked based on the XDG specification. On Unix-like operating systems, this
// will be ~/.config/aiac/aiac.toml.
func LoadConfig(path string) (conf Config, err error) {
	if path == "" {
		path, err = xdg.ConfigFile("aiac/aiac.toml")
		if err != nil {
			return conf, fmt.Errorf("failed getting default config path: %w", err)
		}
	}

	_, err = toml.DecodeFile(path, &conf)
	if err != nil {
		return conf, fmt.Errorf("failed loading configuration: %w", err)
	}

	// If any of the config values are env vars, replace them
	conf = replaceEnvVars(conf)

	return conf, nil
}

// replaceEnvVars replaces any environment variables in the config with their
// actual values.
func replaceEnvVars(conf Config) Config {
	for backendName, backendConfig := range conf.Backends {
		if backendConfig.APIKey != "" {
			backendConfig.APIKey = replaceEnvVar(backendConfig.APIKey)
		}

		if backendConfig.AWSProfile != "" {
			backendConfig.AWSProfile = replaceEnvVar(backendConfig.AWSProfile)
		}

		if backendConfig.AWSRegion != "" {
			backendConfig.AWSRegion = replaceEnvVar(backendConfig.AWSRegion)
		}

		if backendConfig.URL != "" {
			backendConfig.URL = replaceEnvVar(backendConfig.URL)
		}

		if backendConfig.DefaultModel != "" {
			backendConfig.DefaultModel = replaceEnvVar(backendConfig.DefaultModel)
		}

		if backendConfig.APIVersion != "" {
			backendConfig.APIVersion = replaceEnvVar(backendConfig.APIVersion)
		}

		conf.Backends[backendName] = backendConfig
	}

	return conf
}

// replaceEnvVar replaces an environment variable in a string with its actual
// value.
func replaceEnvVar(s string) string {
	return os.ExpandEnv(s)
}
