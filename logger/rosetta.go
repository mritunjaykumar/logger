package logger

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"time"
)

type LogMessage struct {
	ClientIP             string
	CorrelationId        string
	StartTime            time.Time
	EndTime              time.Time
	LatencyNanoSeconds   int64
	LoggerContext        string
	Method               string
	Path                 string
	Protocol             string
	Query                string
	Status               int
	UserAgent            string
	Message              string
	AdditionalProperties map[string]interface{}
}

func New() *LogMessage {
	logMsg := &LogMessage{}
	logMsg.AdditionalProperties = make(map[string]interface{})
	return logMsg
}

// InfoMessage logs log message with INFO level
func InfoMessage(logMessage *LogMessage) {
	infoMessage(logMessage)
}

// ErrorMessage logs log message with ERROR level
func ErrorMessage(logMessage *LogMessage) {
	errorMessage(logMessage)
}

// FatalMessage logs log message with FATAL level
func FatalMessage(logMessage *LogMessage) {
	fatalMessage(logMessage)
}

// WarnMessage logs log message with WARN level
func WarnMessage(logMessage *LogMessage) {
	warnMessage(logMessage)
}

// DebugMessage logs log message with DEBUG level
func DebugMessage(logMessage *LogMessage) {
	debugMessage(logMessage)
}

func (l *LogMessage) SerializeFields(skipGlobalTags bool) string {
	var fields []string
	if l.LoggerContext != "" {
		fields = append(fields, fmt.Sprintf("%v=\"%v\"", loggerContext, l.LoggerContext))
	}
	if l.Status != 0 {
		fields = append(fields, fmt.Sprintf("%v=%v", status, l.Status))
	}
	if l.Method != "" {
		fields = append(fields, fmt.Sprintf("%v=\"%v\"", method, l.Method))
	}
	if l.Protocol != "" {
		fields = append(fields, fmt.Sprintf("%v=\"%v\"", protocol, l.Protocol))
	}
	if l.Path != "" {
		fields = append(fields, fmt.Sprintf("%v=\"%v\"", path, l.Path))
	}
	if l.Query != "" {
		fields = append(fields, fmt.Sprintf("%v=\"%v\"", query, l.Query))
	}
	if l.ClientIP != "" {
		fields = append(fields, fmt.Sprintf("%v=\"%v\"", clientIp, l.ClientIP))
	}
	if l.UserAgent != "" {
		fields = append(fields, fmt.Sprintf("%v=\"%v\"", userAgent, l.UserAgent))
	}
	if !l.StartTime.IsZero() {
		fields = append(fields, fmt.Sprintf("%v=\"%v\"", startTime, l.StartTime.Format(UtcTimeFormat)))
	}
	if !l.EndTime.IsZero() {
		fields = append(fields, fmt.Sprintf("%v=\"%v\"", endTime, l.EndTime.Format(UtcTimeFormat)))
	}
	if l.LatencyNanoSeconds != 0 {
		fields = append(fields, fmt.Sprintf("%v=\"%v\"", latencyUnit, ns))
		fields = append(fields, fmt.Sprintf("%v=%v", latency, l.LatencyNanoSeconds))
	}

	keys := make([]string, 0, len(l.AdditionalProperties))
	for k := range l.AdditionalProperties {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		if reflect.TypeOf(l.AdditionalProperties[key]) == nil {
			fields = append(fields, fmt.Sprintf("%v=\"%v\"", key, nil))
		} else if reflect.TypeOf(l.AdditionalProperties[key]).Kind() == reflect.String {
			fields = append(fields, fmt.Sprintf("%v=\"%v\"", key, l.AdditionalProperties[key]))
		} else {
			fields = append(fields, fmt.Sprintf("%v=%v", key, l.AdditionalProperties[key]))
		}
	}

	if !skipGlobalTags {
		for k, v := range getGlobalTags() {
			fields = append(fields, fmt.Sprintf("%v=\"%v\"", k, v))
		}
	}

	return strings.Join(fields, " ")
}
