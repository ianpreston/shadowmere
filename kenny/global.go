package kenny

import (
    "os"
    "io"
)

const DEFAULT_FORMAT = "{{Time}} - {{Level}} - {{Message}}"
var globalLogger = New(os.Stdout, DEFAULT_FORMAT, DEBUG)

func GlobalConfig(writer io.Writer, format string, level string) {
    globalLogger = New(writer, format, stringToLevel(level))
}

func Debug(line string) {
    globalLogger.Debug(line)
}

func Info(line string) {
    globalLogger.Info(line)
}

func Warn(line string) {
    globalLogger.Warn(line)
}

func Error(line string) {
    globalLogger.Error(line)
}

func ErrorErr(err error) {
    globalLogger.ErrorErr(err)
}

func Critical(line string) {
    globalLogger.Critical(line)
}

func CriticalErr(err error) {
    globalLogger.CriticalErr(err)
}

func Fatal(err error) {
    globalLogger.Fatal(err)
}
