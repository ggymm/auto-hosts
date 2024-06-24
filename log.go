package main

import (
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/natefinch/lumberjack.v2"
)

var log = NewLog()

const timeFormat = "2006-01-02 15:04:05.000"

func console() io.Writer {
	return zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: timeFormat,
	}
}

func rollingFile(name string) io.Writer {
	return &lumberjack.Logger{
		Filename:   name,
		MaxSize:    256, // megabytes
		MaxAge:     30,  // days
		MaxBackups: 128, // files
	}
}

type Log struct {
	log zerolog.Logger
}

func NewLog() *Log {
	return &Log{}
}

func (l *Log) Init(name ...string) {
	file := "app.log"
	if len(name) > 0 {
		file = name[0] + ".log"
	}
	// 设置日志文件路径
	writers := io.MultiWriter(
		rollingFile(file),
		io.MultiWriter(console()),
	)

	zerolog.TimeFieldFormat = timeFormat
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	l.log = zerolog.New(writers).With().Caller().Timestamp().Logger()
}

func (l *Log) Info() *zerolog.Event {
	return l.log.Info()
}

func (l *Log) Error() *zerolog.Event {
	return l.log.Error().Stack()
}

func (l *Log) Fatal() *zerolog.Event {
	return l.log.Fatal().Stack()
}
