package deepseek

import (
	"bytes"
	"fmt"
	"io"
	"time"
	"github.com/gofireflyio/aiac/v5/libaiac/openai"
	"github.com/ido50/requests"
)

// DeepseekBackend is the default URI endpoint for the Deepseek API
const DeepseekBackend = "https://api.deepseek.com/v1"

// DefaultTimeout is the default timeout for requests to the Deepseek API
const DefaultTimeout = 60 * time.Second

// New creates a new instance of the Deepseek client, which uses the OpenAI
// client implementation since Deepseek's API is compatible with OpenAI's.
func New(opts *openai.Options) (*openai.OpenAI, error) {
	if opts == nil {
		opts = &openai.Options{}
	}

	if opts.URL == "" {
		opts.URL = DeepseekBackend
	}

	// Configure the HTTP client with a timeout and debug logging
	httpClient := requests.NewClient(opts.URL).
		Timeout(DefaultTimeout).
		ErrorHandler(func(status int, contentType string, body io.Reader) error {
			var buf bytes.Buffer
			if _, err := io.Copy(&buf, body); err != nil {
				return fmt.Errorf("failed to read error response: %w", err)
			}
			return fmt.Errorf("Deepseek API error (status %d): %s", status, buf.String())
		})

	opts.HTTPClient = httpClient

	return openai.New(opts)
}
