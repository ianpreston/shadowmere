package kenny

import (
	"io"
	"strings"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	CRITICAL
)

type Logger struct {
	writer io.Writer
	format string
	level  LogLevel
}

func New(writer io.Writer, format string, level LogLevel) *Logger {
	return &Logger{
		writer: writer,
		format: format,
		level:  level,
	}
}

func (lg *Logger) Debug(line string) {
	lg.Log(line, DEBUG)
}

func (lg *Logger) Info(line string) {
	lg.Log(line, INFO)
}

func (lg *Logger) Warn(line string) {
	lg.Log(line, WARN)
}

func (lg *Logger) Error(line string) {
	lg.Log(line, ERROR)
}

func (lg *Logger) ErrorErr(err error) {
	if err == nil {
		return
	}
	lg.Log(err.Error(), ERROR)
}

func (lg *Logger) Critical(line string) {
	lg.Log(line, CRITICAL)
}

func (lg *Logger) CriticalErr(err error) {
	if err == nil {
		return
	}
	lg.Log(err.Error(), CRITICAL)
}

// TODO - Does this belong in a logging package?
func (lg *Logger) Fatal(err error) {
	if err == nil {
		return
	}
	lg.Log(err.Error(), CRITICAL)
	panic(err)
}

func (lg *Logger) Log(line string, level LogLevel) {
	if level < lg.level {
		return
	}
	lg.write(lg.formatLine(line, level))
}

func (lg *Logger) formatLine(message string, level LogLevel) string {
	time := time.Now().Format(time.Stamp)
	levelStr := levelToString(level)

	// TODO - Less fugly
	stmt := strings.ToLower(lg.format) + "\n"
	stmt = strings.Replace(stmt, "{{time}}", time, -1)
	stmt = strings.Replace(stmt, "{{level}}", levelStr, -1)
	stmt = strings.Replace(stmt, "{{message}}", message, -1)
	return stmt
}

func (lg *Logger) write(raw string) error {
	_, err := lg.writer.Write([]byte(raw))
	return err
}

func levelToString(level LogLevel) string {
	switch level {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case CRITICAL:
		return "CRITICAL"
	}
	return ""
}

func stringToLevel(level string) LogLevel {
	level = strings.ToUpper(level)
	level = strings.TrimSpace(level)
	switch level {
	case "DEBUG":
		return DEBUG
	case "INFO":
		return INFO
	case "WARN":
		return WARN
	case "ERROR":
		return ERROR
	case "CRITICAL":
		return CRITICAL
	}
	return DEBUG
}
