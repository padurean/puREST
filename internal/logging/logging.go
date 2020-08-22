package logging

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/padurean/purest/internal/env"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Config ...
type Config struct {
	// Level can be one of the values supported by zerolog (https://github.com/rs/zerolog)
	// i.e. from highest to lowest:
	// panic, fatal, error, warn, info, debug, trace
	Level string

	// Enable console logging
	ConsoleLoggingEnabled bool

	// FileLoggingEnabled makes the framework log to a file
	// the fields below can be skipped if this value is false!
	FileLoggingEnabled bool

	// Directory to log to to when filelogging is enabled
	Directory string

	// Filename is the name of the logfile which will be placed inside the directory
	Filename string

	// MaxSize the max size in MB of the logfile before it's rolled
	MaxSize int

	// MaxBackups the max number of rolled files to keep
	MaxBackups int

	// MaxAge the max age in days to keep a logfile
	MaxAge int
}

// Logger ...
type Logger struct {
	*zerolog.Logger
}

// Simple ...
func Simple(r *http.Request) *Logger {
	return SimpleFromCtx(r.Context())
}

// SimpleFromCtx ...
func SimpleFromCtx(ctx context.Context) *Logger {
	logger, ok := ctx.Value("logger").(*Logger)
	if !ok {
		panic(fmt.Sprintf("logger not found on request context"))
	}
	return logger
}

// Detailed ...
func Detailed(r *http.Request) *zerolog.Logger {
	return hlog.FromRequest(r)
}

// DetailedFromCtx ...
func DetailedFromCtx(ctx context.Context) *zerolog.Logger {
	return log.Ctx(ctx)
}

// FromEnv ...
func FromEnv() *Logger {
	return FromConfig(Config{
		Level:                 env.GetLogLevel(),
		ConsoleLoggingEnabled: env.GetLogToConsole(),
		FileLoggingEnabled:    env.GetLogToFile(),
		Directory:             env.GetLogDirectory(),
		Filename:              env.GetLogFilename(),
		MaxSize:               env.GetLogMaxSize(),
		MaxBackups:            env.GetLogMaxBackups(),
		MaxAge:                env.GetLogMaxAge(),
	})
}

// FromConfig sets up the logging framework and creates a new logger
//
// In production, the container logs will be collected and file logging should be disabled. However,
// during development it's nicer to see logs as text and optionally write to a file when debugging
// problems in the containerized pipeline
//
// The output log file will be located at /var/log/service-xyz/service-xyz.log and
// will be rolled according to configuration set.
func FromConfig(config Config) *Logger {
	var writers []io.Writer

	if config.ConsoleLoggingEnabled {
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr})
	}
	if config.FileLoggingEnabled {
		writers = append(writers, newRollingFile(config))
	}
	mw := io.MultiWriter(writers...)

	logLevel, err := zerolog.ParseLevel(config.Level)
	if err != nil {
		logLevel = zerolog.InfoLevel
		log.Error().Msgf("error parsing log level '%s' from env var: %v", config.Level, err)
	}
	zerolog.SetGlobalLevel(logLevel)

	logger := zerolog.New(mw).Level(logLevel).With().Timestamp().Logger()

	logger.Info().Msgf(
		"logging configured: level:%s, logToConsole:%t, logToFile:%t, logsDir:%s, logFilename:%s, maxFileSize:%dMB, maxBackups:%d, maxAge:%ddays",
		logger.GetLevel(),
		config.ConsoleLoggingEnabled,
		config.FileLoggingEnabled,
		config.Directory,
		config.Filename,
		config.MaxSize,
		config.MaxBackups,
		config.MaxAge,
	)

	return &Logger{
		Logger: &logger,
	}
}

func newRollingFile(config Config) io.Writer {
	if err := os.MkdirAll(config.Directory, 0744); err != nil {
		log.Error().Err(err).Str("path", config.Directory).Msg("can't create log directory")
		return nil
	}

	return &lumberjack.Logger{
		Filename:   path.Join(config.Directory, config.Filename),
		MaxBackups: config.MaxBackups, // files
		MaxSize:    config.MaxSize,    // megabytes
		MaxAge:     config.MaxAge,     // days
	}
}
