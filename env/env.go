package env

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

const appPrefix = "PUREST_"

const dbPrefix = appPrefix + "DB_"
const dbDriver = dbPrefix + "DRIVER"
const dbURL = dbPrefix + "URL"
const dbUser = dbPrefix + "USER"
const dbSchema = dbPrefix + "SCHEMA"

const httpPrefix = appPrefix + "HTTP_"
const httpPort = httpPrefix + "PORT"

const logPrefix = appPrefix + "LOG_"
const logLevel = logPrefix + "LEVEL"
const logToConsole = logPrefix + "TO_CONSOLE"
const logToFile = logPrefix + "TO_FILE"
const logDirectory = logPrefix + "DIRECTORY"
const logFilename = logPrefix + "FILENAME"
const logMaxSize = logPrefix + "MAX_SIZE"
const logMaxBackups = logPrefix + "MAX_BACKUPS"
const logMaxAge = logPrefix + "MAX_AGE"
const logRequests = logPrefix + "REQUESTS"

// Load environment. For details see https://github.com/joho/godotenv
func init() {
	log.Info().Msg("loading env ...")
	env := GetAppEnv().String()
	log.Info().Msgf("loaded `%s` env", env)

	godotenv.Load(".env." + env + ".local")
	if "test" != env {
		godotenv.Load(".env.local")
	}
	godotenv.Load(".env." + env)
	godotenv.Load() // The Original .env
}

func getEnvOrPanic(key string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		panic(fmt.Sprintf("Env var '%s' is not set", key))
	}
	log.Info().Msgf("ENV: %s=%s", key, value)
	return value
}

func getBoolEnvOrPanic(key string) bool {
	vs := getEnvOrPanic(key)
	v, err := strconv.ParseBool(vs)
	if err != nil {
		panic(fmt.Sprintf("Env var '%s' value '%s' is not a boolean (true or false): %v", key, vs, err))
	}
	return v
}

func getIntEnvOrPanic(key string) int {
	vs := getEnvOrPanic(key)
	v, err := strconv.Atoi(vs)
	if err != nil {
		panic(fmt.Sprintf("Env var '%s' value '%s' is not an integer number: %v", key, vs, err))
	}
	return v
}

// AppEnv ...
type AppEnv int8

// AppEnv ...
const (
	Development AppEnv = iota
	Test
	Production
)

func (ae AppEnv) String() string {
	switch ae {
	case Test:
		return "test"
	case Production:
		return "production"
	default:
		return "development"
	}
}

// ParseAppEnv ...
func ParseAppEnv(appEnvStr string) AppEnv {
	switch appEnvStr {
	case Test.String():
		return Test
	case Production.String():
		return Production
	default:
		return Development
	}
}

// GetAppEnv ...
func GetAppEnv() AppEnv {
	return ParseAppEnv(strings.ToLower(strings.TrimSpace(os.Getenv(appPrefix + "ENV"))))
}

// GetDbDriver ...
func GetDbDriver() string {
	return getEnvOrPanic(dbDriver)
}

// GetDbURL ...
func GetDbURL() string {
	return getEnvOrPanic(dbURL)
}

// GetDbUser ...
func GetDbUser() string {
	return getEnvOrPanic(dbUser)
}

// GetDbSchema ...
func GetDbSchema() string {
	return getEnvOrPanic(dbSchema)
}

// GetHTTPPort ...
func GetHTTPPort() string {
	return getEnvOrPanic(httpPort)
}

// GetLogLevel ...
func GetLogLevel() string {
	return getEnvOrPanic(logLevel)
}

// GetLogToConsole ...
func GetLogToConsole() bool {
	return getBoolEnvOrPanic(logToConsole)
}

// GetLogToFile ...
func GetLogToFile() bool {
	return getBoolEnvOrPanic(logToFile)
}

// GetLogDirectory ...
func GetLogDirectory() string {
	return getEnvOrPanic(logDirectory)
}

// GetLogFilename ...
func GetLogFilename() string {
	return getEnvOrPanic(logFilename)
}

// GetLogMaxSize ...
func GetLogMaxSize() int {
	return getIntEnvOrPanic(logMaxSize)
}

// GetLogMaxBackups ...
func GetLogMaxBackups() int {
	return getIntEnvOrPanic(logMaxBackups)
}

// GetLogMaxAge ...
func GetLogMaxAge() int {
	return getIntEnvOrPanic(logMaxAge)
}

// GetLogRequests ...
func GetLogRequests() bool {
	return getBoolEnvOrPanic(logRequests)
}
