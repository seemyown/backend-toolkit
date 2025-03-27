package logging

// Пакет для настройки логирования

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	zl "github.com/rs/zerolog"
)

func createLogger(logFile string) zl.Logger {
	if err := os.MkdirAll(filepath.Dir(logFile), os.ModePerm); err != nil {
		panic(err)
	}

	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	consoleWriter := zl.ConsoleWriter{
		Out: os.Stdout, TimeFormat: time.RFC3339, NoColor: true,
		FormatLevel: func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
		},
		FormatMessage: func(i interface{}) string {
			return fmt.Sprintf("%s", i)
		},
		FormatFieldName: func(i interface{}) string {
			return fmt.Sprintf("%s:", i)
		},
		FormatFieldValue: func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("%s", i))
		},
	}

	multiWriter := zl.MultiLevelWriter(consoleWriter, file)

	return zl.New(multiWriter).With().Timestamp().Logger()
}

type Logger struct {
	logger zl.Logger
	name   string
}

type Config struct {
	FileName string
	Path     string
	Name     string
}

func New(config Config) *Logger {
	var path string
	if config.Path != "" {
		path = config.Path
	} else {
		path = "./logs"
	}

	path = filepath.Join(path, fmt.Sprintf("%s.log", config.FileName))
	return &Logger{
		logger: createLogger(path),
		name:   config.Name,
	}
}

func (l *Logger) NewSubLogger(component string) *Logger {
	ctx := l.logger.With().Str("component", component)
	return &Logger{
		logger: ctx.Logger(),
		name:   l.name,
	}
}

func (l *Logger) withCaller() zl.Logger {
	pc, _, _, ok := runtime.Caller(2) // 2 - чтобы получить метод, вызвавший логгер
	if !ok {
		return l.logger
	}

	funcName := runtime.FuncForPC(pc).Name()
	parts := strings.Split(funcName, "/")
	fullName := parts[len(parts)-1] // Оставляем только последнее имя функции
	shortName := l.name

	return l.logger.With().Str("caller", fmt.Sprintf("%s(%s)", shortName, fullName)).Logger()
}

func (l *Logger) Info(format string, args ...interface{}) {
	caller := l.withCaller()
	caller.Info().Msgf(format, args...)
}
func (l *Logger) Warn(format string, args ...interface{}) {
	caller := l.withCaller()
	caller.Warn().Msgf(format, args...)
}
func (l *Logger) Error(err error, format string, args ...interface{}) {
	caller := l.withCaller()
	caller.Error().Err(err).Msgf(format, args...)
}
func (l *Logger) Debug(format string, args ...interface{}) {
	caller := l.withCaller()
	caller.Debug().Msgf(format, args...)
}
func (l *Logger) Trace(format string, args ...interface{}) {
	caller := l.withCaller()
	caller.Trace().Msgf(format, args...)
}
func (l *Logger) Fatal(err error, format string, args ...interface{}) {
	caller := l.withCaller()
	caller.Fatal().Err(err).Msgf(format, args...)
}
