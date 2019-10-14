package kit

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Initizalize the Mu-Kit environment
func Init() {
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
}
