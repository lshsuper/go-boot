package nlog

import (
	"os"
)

type consoleLogger struct {
	baseLogger
}

func NewConsoleLogger() *consoleLogger {

	logger := &consoleLogger{}
	logger.Init()
	logger.SetWriterHandler(logger.writer)
	return logger

}

func (logger *consoleLogger) Info(msg ...interface{}) {
	logger.putMsg(INFO, msg...)
}

func (logger *consoleLogger) Debug(msg ...interface{}) {
	logger.putMsg(DEBUG, msg...)
}
func (logger *consoleLogger) Error(msg ...interface{}) {
	logger.putMsg(ERROR, msg...)
}

func (logger *consoleLogger) writer(level LogLevel, by []byte) {

	switch level {
	case DEBUG:
		_, _ = os.Stdin.Write(by)
	case INFO:
		_, _ = os.Stdout.Write(by)
	case ERROR:
		_, _ = os.Stderr.Write(by)
	default:
		panic("err: level not found")
	}

}
