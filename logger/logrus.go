package logger

import (
	"fmt"
)

type Fields map[string]interface{}

func (fields Fields) CloneWith(name string, value interface{}) Fields {
	newFields := Fields{}

	for k, v := range fields {
		newFields[k] = v
	}

	newFields[name] = value

	return newFields
}

func (fields Fields) CloneWithAll(otherFields Fields) Fields {
	newFields := Fields{}

	for k, v := range fields {
		newFields[k] = v
	}

	for k, v := range otherFields {
		newFields[k] = v
	}

	return newFields
}

func (fields Fields) ToMap() map[string]interface{} {
	fieldsAsMap := make(map[string]interface{})

	for k, v := range fields {
		fieldsAsMap[k] = v
	}

	return fieldsAsMap
}

type entry struct {
	value Fields
}

func (e *entry) Info(msg string) {
	infoMessage(e.storeFields(msg))
}

func (e *entry) Infof(format string, args ...interface{}) {
	infoMessage(e.storeFields(fmt.Sprintf(format, args...)))
}

func (e *entry) Debug(msg string) {
	debugMessage(e.storeFields(msg))
}

func (e *entry) Debugf(format string, args ...interface{}) {
	debugMessage(e.storeFields(fmt.Sprintf(format, args...)))
}

func (e *entry) Error(msg string) {
	errorMessage(e.storeFields(msg))
}

func (e *entry) Errorf(format string, args ...interface{}) {
	errorMessage(e.storeFields(fmt.Sprintf(format, args...)))
}

func (e *entry) Warn(msg string) {
	warnMessage(e.storeFields(msg))
}

func (e *entry) Warnf(format string, args ...interface{}) {
	warnMessage(e.storeFields(fmt.Sprintf(format, args...)))
}

func (e *entry) Warning(msg string) {
	warnMessage(e.storeFields(msg))
}

func (e *entry) Warningf(format string, args ...interface{}) {
	warnMessage(e.storeFields(fmt.Sprintf(format, args...)))
}

func (e *entry) Print(msg string) {
	infoMessage(e.storeFields(msg))
}

func (e *entry) Printf(format string, args ...interface{}) {
	infoMessage(e.storeFields(fmt.Sprintf(format, args...)))
}

func (e *entry) Fatal(msg string) {
	fatalMessage(e.storeFields(msg))
}

func (e *entry) Fatalf(format string, args ...interface{}) {
	fatalMessage(e.storeFields(fmt.Sprintf(format, args...)))
}

func (e *entry) WithField(key string, value interface{}) *entry {
	e.value[key] = value
	return e
}

func (e *entry) WithFields(fields Fields) *entry {
	for k, v := range fields {
		e.value[k] = v
	}

	return e
}

func (e *entry) WithError(err error) *entry {
	const errorFieldKey = "error"

	if err != nil {
		e.value[errorFieldKey] = err.Error()
	}

	return e
}

func (e *entry) storeFields(msg string) *LogMessage {
	logMessage := &LogMessage{
		Message:              msg,
		AdditionalProperties: make(map[string]interface{}),
	}

	for key, val := range e.value {
		logMessage.AdditionalProperties[key] = val
	}

	return logMessage
}

func WithField(key string, value interface{}) *entry {
	return &entry{
		value: Fields{
			key: value,
		},
	}
}

func WithFields(fields Fields) *entry {
	newEntry := &entry{
		value: make(Fields),
	}

	for k, v := range fields {
		newEntry.value[k] = v
	}

	return newEntry
}

func WithError(err error) *entry {
	newEntry := &entry{
		value: make(Fields),
	}

	return newEntry.WithError(err)
}

func Info(args ...interface{}) {
	infoMessage(&LogMessage{Message: fmt.Sprint(args...)})
}

func Infof(format string, args ...interface{}) {
	infoMessage(&LogMessage{Message: fmt.Sprintf(format, args...)})
}

func Print(args ...interface{}) {
	infoMessage(&LogMessage{Message: fmt.Sprint(args...)})
}

func Printf(format string, args ...interface{}) {
	infoMessage(&LogMessage{Message: fmt.Sprintf(format, args...)})
}

func Warn(args ...interface{}) {
	warnMessage(&LogMessage{Message: fmt.Sprint(args...)})
}

func Warning(args ...interface{}) {
	warnMessage(&LogMessage{Message: fmt.Sprint(args...)})
}

func Warnf(format string, args ...interface{}) {
	warnMessage(&LogMessage{Message: fmt.Sprintf(format, args...)})
}

func Warningf(format string, args ...interface{}) {
	warnMessage(&LogMessage{Message: fmt.Sprintf(format, args...)})
}

func Error(args ...interface{}) {
	errorMessage(&LogMessage{Message: fmt.Sprint(args...)})
}

func Errorf(format string, args ...interface{}) {
	errorMessage(&LogMessage{Message: fmt.Sprintf(format, args...)})
}

func Fatal(args ...interface{}) {
	fatalMessage(&LogMessage{Message: fmt.Sprint(args...)})
}

func Fatalf(format string, args ...interface{}) {
	fatalMessage(&LogMessage{Message: fmt.Sprintf(format, args...)})
}

func Debug(args ...interface{}) {
	debugMessage(&LogMessage{Message: fmt.Sprint(args...)})
}

func Debugf(format string, args ...interface{}) {
	debugMessage(&LogMessage{Message: fmt.Sprintf(format, args...)})
}

func SetLevel(level string) error {
	return setLogLevel(level)
}

func GetLevel() string {
	return getLogLevel().String()
}

// AddStacktrace configures the Logger to record a stack trace for all messages at or above a given level.
func AddStackTrace(logLevel string) {
	addStackTrace(logLevel)
}
