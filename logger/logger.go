package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
	"time"
)

const (
	ns            = "ns"
	utcTimeFormat = "2006-01-02T15:04:05Z0700"
)

var (
	// Logger creates zap logger based on the config
	logger *zap.Logger
)

type LogMessage struct {
	ClientIP           string
	StartTime          time.Time
	EndTime            time.Time
	LatencyNanoSeconds int64
	LoggerContext      string
	Method             string
	Path               string
	Protocol           string
	Query              string
	Status             int
	UserAgent          string
	Message            string
	Tags               map[string]interface{}
}

// UTC time encode
func utcTimeEncode(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.UTC().Format(utcTimeFormat))
}
func init() {
	zapConfig := zap.NewProductionConfig()
	zapConfig.EncoderConfig.EncodeTime = utcTimeEncode
	zapConfig.EncoderConfig.TimeKey = "timestamp"
	// ADD additional custom tags to the logs
	zapConfig.InitialFields = GetGlobalTags()
	zapConfig.Sampling = nil
	var err error
	if logger, err = zapConfig.Build(); err != nil {
		panic(err)
	}
}

// GetGlobalTags provides global tags added to the logs
func GetGlobalTags() map[string]interface{} {
	// ADD additional custom tags to the logs
	globalTags := make(map[string]interface{})
	globalTags["application"] = "astra"
	tempComponent := os.Args[0] // this might provide value like "/go/bin/usersapi"
	// Get just the app name and not the whole path. For example: out of "/go/bin/usersapi", just get "usersapi"
	globalTags["component"] = tempComponent[strings.LastIndex(tempComponent, "/")+1:]
	return globalTags
}
func (l *LogMessage) getZapFields() []zap.Field {
	var fields []zap.Field
	if l.LoggerContext != "" {
		fields = append(fields, zap.String("logger-context", l.LoggerContext))
	}
	if l.Status != 0 {
		fields = append(fields, zap.Int("status", l.Status))
	}
	if l.Method != "" {
		fields = append(fields, zap.String("method", l.Method))
	}
	if l.Protocol != "" {
		fields = append(fields, zap.String("protocol", l.Protocol))
	}
	if l.Path != "" {
		fields = append(fields, zap.String("path", l.Path))
	}
	if l.Query != "" {
		fields = append(fields, zap.String("query", l.Query))
	}
	if l.ClientIP != "" {
		fields = append(fields, zap.String("client-ip", l.ClientIP))
	}
	if l.UserAgent != "" {
		fields = append(fields, zap.String("user-agent", l.UserAgent))
	}
	if !l.StartTime.IsZero() {
		fields = append(fields, zap.String("start-time", l.StartTime.Format(utcTimeFormat)))
	}
	if !l.EndTime.IsZero() {
		fields = append(fields, zap.String("end-time", l.EndTime.Format(utcTimeFormat)))
	}
	if l.LatencyNanoSeconds != 0 {
		fields = append(fields, zap.String("latency-unit", ns))
		fields = append(fields, zap.Int64("latency", l.LatencyNanoSeconds))
	}
	for key, val := range l.Tags {
		switch v := val.(type) {
		case bool:
			fields = append(fields, zap.Bool(key, v))
		case int:
			fields = append(fields, zap.Int(key, v))
		case int64:
			fields = append(fields, zap.Int64(key, v))
		case float64:
			fields = append(fields, zap.Float64(key, v))
		case string:
			fields = append(fields, zap.String(key, v))
		default:
			logger.Error("Encounter unsupported type [" + fmt.Sprintf("%T", val) + "]")
			// TODO: what to do here?
		}
	}
	return fields
}

// Log wraps zap "Info" function
func Info(logMessage *LogMessage) {
	message := logMessage.Message
	fields := logMessage.getZapFields()
	//fields = []zap.Field{}
	//fields = append(fields, zap.String("foo", "bar"))
	logger.Info(message, fields...) //fields ...zap.Field
	logger.Sync()
}

// Error wraps zap "Error" function
func Error(logMessage *LogMessage) {
	message := logMessage.Message
	fields := logMessage.getZapFields()
	logger.Error(message, fields...) //fields ...zap.Field
	logger.Sync()
}

// Warn wraps zap "Warn" function
func Warn(logMessage *LogMessage) {
	message := logMessage.Message
	fields := logMessage.getZapFields()
	logger.Warn(message, fields...) //fields ...zap.Field
	logger.Sync()
}

// Debug wraps zap "Debug" function
func Debug(logMessage *LogMessage) {
	message := logMessage.Message
	fields := logMessage.getZapFields()
	logger.Debug(message, fields...) //fields ...zap.Field
	logger.Sync()
}
