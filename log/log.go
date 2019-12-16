package log

import (
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	zlog "github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

type Config struct {
	Level string `default:"info"`
}

func Init(name string, config *Config) {

	level := zerolog.InfoLevel
	var err error
	switch strings.ToLower(config.Level) {
	case "panic":
		level = zerolog.PanicLevel
	case "fatal":
		level = zerolog.FatalLevel
	case "error":
		level = zerolog.ErrorLevel
	case "warn":
		level = zerolog.WarnLevel
	case "info":
		level = zerolog.InfoLevel
	case "debug":
		level = zerolog.DebugLevel
	case "trace":
		level = zerolog.TraceLevel
	default:
		err = fmt.Errorf("Unknown log level %s", config.Level)
	}

	zerolog.SetGlobalLevel(level)
	zerolog.TimestampFieldName = "t"
	zerolog.LevelFieldName = "l"
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.MessageFieldName = "m"
	zlog.Logger = zerolog.New(os.Stdout).With().Str("service", name).Timestamp().Logger()

	if err != nil {
		log.Error().Err(err).Send()
	}
}
