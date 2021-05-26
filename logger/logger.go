package logger

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	application   = "astra"
	correlationId = "correlation-id"
	clientIp      = "client-ip"
	endTime       = "end-time"
	latency       = "latency"
	latencyUnit   = "latency-unit"
	loggerContext = "rosetta-context"
	method        = "method"
	nilLogMessage = "rosetta is called with nil log message"
	ns            = "ns"
	path          = "path"
	protocol      = "protocol"
	query         = "query"
	startTime     = "start-time"
	status        = "status"
	timeStamp     = "timestamp"
	userAgent     = "user-agent"
	UtcTimeFormat = "2006-01-02T15:04:05.000000Z0700"

	// Supported log levels
	LogLevel     = "LOG_LEVEL"
	DebugLevel   = "DEBUG"
	InfoLevel    = "INFO"
	WarnLevel    = "WARN"
	WarningLevel = "WARNING"
	ErrorLevel   = "ERROR"
	FatalLevel   = "FATAL"

	LoggerEnvironment = "LOGGER_ENVIRONMENT"
	development       = "DEVELOPMENT"
	dev               = "DEV"
	logOutputFile     = "LOG_OUTPUT_FILE"
)

var (
	zapLogger         *zap.Logger            // zap logger instance based on the zapLogger environment and other config settings
	logEnv            string                 // logger environment (DEV or non-dev (PROD, STAGING or anything else)
	logLvl            = zap.NewAtomicLevel() // Dynamic log level
	initZapLoggerOnce sync.Once
	NoStacktrace      string
)

// UTC time encode
func utcTimeEncode(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.UTC().Format(UtcTimeFormat))
}

// Init initializes rosetta zapLogger.
// It uses following environment variables to override any configuration
// 		- LOGGER_ENVIRONMENT. If this has value of "DEVELOPMENT" or "DEV", it defaults to
//							  to NewDevelopmentConfig with console colored plain-text logging (no JSON).
//							  Otherwise, it will default to NewProductionConfig with JSON formatted logging.
//		- LOG_OUTPUT_FILE. If it's not empty, it will create a log file with that name and start writing logs
// 						   to log file.
//		- LOG_LEVEL. Supported log levels are DEBUG, INFO, WARN, ERROR, PANIC and FATAL
// Make sure we are creating ONLY one instance of zapLogger.
func GetZapLogger() *zap.Logger {
	initZapLoggerOnce.Do(func() {
		buildZapLogger("")
	})
	return zapLogger
}

func buildZapLogger(memoryOutputPathName string) {
	const callerSkipOffset = 3
	zapConfig := getConfigBasedOnLoggerEnvironment()

	logLvl = zapConfig.Level // Initial log-level

	// override log-level if LOG_LEVEL env variable is set
	setLogLevelFromEnvironment()

	zapConfig.EncoderConfig.EncodeTime = utcTimeEncode
	zapConfig.EncoderConfig.TimeKey = timeStamp
	zapConfig.EncoderConfig.EncodeDuration = zapcore.MillisDurationEncoder
	setFileOutput(&zapConfig)

	if memoryOutputPathName != "" {
		// Redirect all messages to the MemorySink.
		zapConfig.OutputPaths = []string{fmt.Sprintf("%s://", memoryOutputPathName)}
	}

	zapConfig.Sampling = nil
	var err error
	if zapLogger, err = zapConfig.Build(zap.AddCallerSkip(callerSkipOffset)); err != nil {
		panic(err)
	}
}

func getConfigBasedOnLoggerEnvironment() zap.Config {
	logEnv = os.Getenv(LoggerEnvironment)
	var zapConfig zap.Config
	if logEnv == development || logEnv == dev {
		zapConfig = zap.NewDevelopmentConfig()
		zapConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		zapConfig = zap.NewProductionConfig()
	}
	return zapConfig
}

func setLogLevelFromEnvironment() {
	// We are ignoring error returned by the below function call
	setLogLevel(os.Getenv(LogLevel))
}

// AddStacktrace configures the Logger to record a stack trace for all messages at or above a given level.
func addStackTrace(logLevel string) {
	fmt.Println(fmt.Sprintf("value of NoStacktrace is [%v]", NoStacktrace))
	nst, err := strconv.ParseBool(NoStacktrace)
	if err != nil {
		fmt.Println(errors.New(fmt.Sprintf("cannot parse bool value for [%v]", NoStacktrace)))
	}

	if nst {
		zapLogger = GetZapLogger().WithOptions()
		return
	}

	switch logLevel {
	case DebugLevel:
		zapLogger = GetZapLogger().WithOptions(zap.AddStacktrace(zap.DebugLevel))
	case InfoLevel:
		zapLogger = GetZapLogger().WithOptions(zap.AddStacktrace(zap.InfoLevel))
	case WarnLevel, WarningLevel:
		zapLogger = GetZapLogger().WithOptions(zap.AddStacktrace(zap.WarnLevel))
	case ErrorLevel:
		zapLogger = GetZapLogger().WithOptions(zap.AddStacktrace(zap.ErrorLevel))
	default:
		fmt.Println(errors.New(fmt.Sprintf("Cannot add stack trace for level %v", logLevel)))
	}
}

func setLogLevel(level string) error {
	switch level {
	case DebugLevel:
		logLvl.SetLevel(zapcore.DebugLevel)
	case InfoLevel:
		logLvl.SetLevel(zapcore.InfoLevel)
	case WarnLevel, WarningLevel:
		logLvl.SetLevel(zapcore.WarnLevel)
	case ErrorLevel:
		logLvl.SetLevel(zapcore.ErrorLevel)
	case FatalLevel:
		logLvl.SetLevel(zapcore.FatalLevel)
	default:
		return errors.New(fmt.Sprintf("unknown log level %v, so log level in not set", level))
	}

	return nil
}

func getLogLevel() zap.AtomicLevel {
	return logLvl
}

// setFileOutput sets the log output file if it has some value for env variable "LOG_OUTPUT_FILE"
func setFileOutput(config *zap.Config) {
	outputFile := os.Getenv(logOutputFile)
	if outputFile == "" {
		return
	}
	config.OutputPaths = append(config.OutputPaths, outputFile)
}

// getGlobalTags provides global tags added to the logs
func getGlobalTags() map[string]string {
	// ADD additional custom tags to the logs
	globalTags := make(map[string]string)

	globalTags["application"] = application
	tempComponent := os.Args[0] // this might provide value like "/go/bin/usersapi"

	// Get just the app name and not the whole path. For example: out of "/go/bin/usersapi", just get "usersapi"
	globalTags["component"] = tempComponent[strings.LastIndex(tempComponent, "/")+1:]
	return globalTags
}

// zap info wrapper
func infoMessage(logMessage *LogMessage) {
	callZapLogger(logMessage, GetZapLogger().Info)
}

// errorMessage wraps zap "Error" function
func errorMessage(logMessage *LogMessage) {
	callZapLogger(logMessage, GetZapLogger().Error)
}

// fatalMessage wraps zap "Fatal" function
func fatalMessage(logMessage *LogMessage) {
	callZapLogger(logMessage, GetZapLogger().Fatal)
}

// warnMessage wraps zap "Warn" function
func warnMessage(logMessage *LogMessage) {
	callZapLogger(logMessage, GetZapLogger().Warn)
}

// debugMessage wraps zap "Debug" function
func debugMessage(logMessage *LogMessage) {
	callZapLogger(logMessage, GetZapLogger().Debug)
}

// callZapLogger calls the zap logger functions.
func callZapLogger(logMessage *LogMessage, logCaller func(msg string, fields ...zap.Field)) {
	if logMessage == nil {
		logCaller = GetZapLogger().Error
		logCaller(nilLogMessage)
	} else {
		if logEnv == development || logEnv == dev {
			logCaller(fmt.Sprintf("%v %v", logMessage.Message, logMessage.SerializeFields(true)))
		} else {
			fields := logMessage.getZapFields()
			logCaller(logMessage.Message, fields...)
		}
	}
	GetZapLogger().Sync()
}

func (l *LogMessage) getZapFields() []zap.Field {
	var fields []zap.Field
	if l.LoggerContext != "" {
		fields = append(fields, zap.String(loggerContext, l.LoggerContext))
	}
	if l.Status != 0 {
		fields = append(fields, zap.Int(status, l.Status))
	}
	if l.Method != "" {
		fields = append(fields, zap.String(method, l.Method))
	}
	if l.Protocol != "" {
		fields = append(fields, zap.String(protocol, l.Protocol))
	}
	if l.Path != "" {
		fields = append(fields, zap.String(path, l.Path))
	}
	if l.Query != "" {
		fields = append(fields, zap.String(query, l.Query))
	}
	if l.ClientIP != "" {
		fields = append(fields, zap.String(clientIp, l.ClientIP))
	}
	if l.UserAgent != "" {
		fields = append(fields, zap.String(userAgent, l.UserAgent))
	}
	if !l.StartTime.IsZero() {
		fields = append(fields, zap.String(startTime, l.StartTime.Format(UtcTimeFormat)))
	}
	if !l.EndTime.IsZero() {
		fields = append(fields, zap.String(endTime, l.EndTime.Format(UtcTimeFormat)))
	}
	if l.LatencyNanoSeconds != 0 {
		fields = append(fields, zap.String(latencyUnit, ns))
		fields = append(fields, zap.Int64(latency, l.LatencyNanoSeconds))
	}
	for key, val := range l.AdditionalProperties {
		fields = append(fields, zap.Any(key, val))
	}

	for k, v := range getGlobalTags() {
		fields = append(fields, zap.String(k, v))
	}

	return fields
}
