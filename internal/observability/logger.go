package observability

import (
	"os"
	"time"

	"github.com/rs/zerolog"
)

func NewLogger(env string) zerolog.Logger {
	l := zerolog.New(os.Stdout).With().Timestamp().Logger()
	if env == "development" {
		zerolog.TimeFieldFormat = time.RFC3339
	}
	return l
}

