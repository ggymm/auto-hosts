package log

import (
	"github.com/rs/zerolog"
)

var log zerolog.Logger

func Info() *zerolog.Event {
	return log.Info()
}

func Error() *zerolog.Event {
	return log.Error().Stack()
}

func Fatal() *zerolog.Event {
	return log.Fatal().Stack()
}
