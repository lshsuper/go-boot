package nlog

import (
	"io"
	"os"
	"path"
)

type fileLogger struct {
	baseLogger
	errWriter   io.Writer
	infoWriter  io.Writer
	debugWriter io.Writer
}
type FileLoggerConfig struct {
	ErrorLogPath string
	InfoLogPath  string
	DebugLogPath string
}

func NewFileLogger(conf FileLoggerConfig) *fileLogger {

	logger := &fileLogger{}
	logger.Init()
	logger.SetWriterHandler(logger.writer)

	os.MkdirAll(path.Dir(conf.InfoLogPath), os.ModePerm)
	os.MkdirAll(path.Dir(conf.InfoLogPath), os.ModePerm)
	os.MkdirAll(path.Dir(conf.DebugLogPath), os.ModePerm)

	logger.errWriter, _ = os.OpenFile(conf.ErrorLogPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModePerm)
	logger.infoWriter, _ = os.OpenFile(conf.InfoLogPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModePerm)
	logger.debugWriter, _ = os.OpenFile(conf.DebugLogPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModePerm)

	return logger

}

func (logger *fileLogger) Info(msg ...interface{}) {
	logger.putMsg(INFO, msg...)
}

func (logger *fileLogger) Debug(msg ...interface{}) {
	logger.putMsg(DEBUG, msg...)
}
func (logger *fileLogger) Error(msg ...interface{}) {
	logger.putMsg(ERROR, msg...)
}

func (logger *fileLogger) writer(level LogLevel, by []byte) {

	switch level {
	case DEBUG:
		logger.debugWriter.Write(by)
	case INFO:
		logger.infoWriter.Write(by)
	case ERROR:
		logger.errWriter.Write(by)
	default:
		panic("err: level not found")
	}

}
