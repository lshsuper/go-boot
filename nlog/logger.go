package nlog

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"
)

//常量定义
const (
	defaultChSize     int = 100
	defaultTimeFormat     = "2006/01/02 15:04:05"
)

//LogLevel 日志层级枚举
type LogLevel int

const (
	_     LogLevel = iota << 1
	DEBUG          //dev
	INFO           //基础信息型
	ERROR          //错误型
)

//BString 枚举格式化函数
func (e LogLevel) BString() string {
	switch e {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case ERROR:
		return "ERROR"
	default:
		panic("LogLevel 404")
	}
}

//SString 枚举格式化函数
func (e LogLevel) SString() string {
	switch e {
	case DEBUG:
		return "debug"
	case INFO:
		return "info"
	case ERROR:
		return "error"
	default:
		panic("LogLevel 404")
	}
}

//Logger 日志接口
type Logger interface {
	Info(msg ...interface{})
	Error(msg ...interface{})
	Debug(msg ...interface{})
}

//logMsg 消息体
type logMsg struct {
	Msg    string `json:"msg"`
	FsPos  string `json:"fs-pos"`
	FunPos string `json:"fun-pos"`
	T      string `json:"time"`
	Level  string `json:"level"`
}

func (l logMsg) JsonFormat(by *[]byte) {
	*by, _ = json.Marshal(l)
}

func (l logMsg) StrFormat(by *[]byte) {

	*by = append(*by, "[time]:"+l.T...)
	*by = append(*by, "|[level]:"+l.Level...)
	*by = append(*by, "|[pos]:"+l.FsPos+":"+l.FunPos...)
	*by = append(*by, "|[msg]:"+l.Msg...)

}

//releaseLogMsg 释放消息对象
func (base *baseLogger) releaseLogMsg(m logMsg) {
	m.Msg = ""
	m.FunPos = ""
	m.T = ""
	m.FsPos = ""
	m.Level = ""
	base.msgPool.Put(m)
	return
}

//WriterHandler 写入事件类型
type WriterHandler func(level LogLevel, by []byte)

//baseLogger 基础日志实例
type baseLogger struct {
	Logger
	errorBuf []byte
	infoBuf  []byte
	debugBuf []byte

	errorCh  chan logMsg
	infoCh   chan logMsg
	debugCh  chan logMsg
	writer   WriterHandler
	callSkip int
	msgPool  *sync.Pool
	format   LogFormat
}

type LogFormat int

const (
	_ LogFormat = iota
	StrFormat
	JsonFormat
)

//Init 初始化构建
func (base *baseLogger) Init() {

	base.infoCh = make(chan logMsg, defaultChSize)
	base.debugCh = make(chan logMsg, defaultChSize)
	base.errorCh = make(chan logMsg, defaultChSize)
	base.msgPool = &sync.Pool{
		New: func() interface{} {
			return logMsg{}
		},
	}
	base.SetFormat(StrFormat)
	go base.msgConsumer()

}

func (base *baseLogger) SetFormat(f LogFormat) {
	base.format = f
}

//SetCallSkip 设置位置调用跨越层级
func (base *baseLogger) SetCallSkip(skip int) {
	base.callSkip = skip
}

//getPos 定位异常位置
func (base *baseLogger) getPos() (fsPos, funPos string) {
	var (
		ok bool
		pc uintptr
	)
	pc, fsPos, _, ok = runtime.Caller(base.callSkip)
	if !ok {
		return
	}
	f := runtime.FuncForPC(pc)
	if f == nil {
		return
	}
	funPos = f.Name()
	return

}

//timeFormat 时间格式化
func (base *baseLogger) timeFormat() string {
	return time.Now().Format(defaultTimeFormat)
}

//SetWriterHandler 设置消息写入事件
func (base *baseLogger) SetWriterHandler(fn func(level LogLevel, by []byte)) {
	base.writer = fn
}

//getLogMsg 构建消息体
func (base *baseLogger) buildLogMsg(level LogLevel, msg ...interface{}) logMsg {

	l := base.msgPool.Get().(logMsg)

	if base.format == JsonFormat {
		l.Msg = fmt.Sprint(msg...)
	} else {
		l.Msg = fmt.Sprintln(msg...)
	}

	l.FsPos, l.FunPos = base.getPos()
	l.T = base.timeFormat()
	l.Level = level.SString()
	return l
}

//putMsg 入队消息
func (base *baseLogger) putMsg(level LogLevel, msg ...interface{}) {

	m := base.buildLogMsg(level, msg...)
	switch level {
	case DEBUG:
		base.debugCh <- m
	case INFO:
		base.infoCh <- m
	case ERROR:
		base.errorCh <- m
	}

}

//recover 捕获异常
func (base *baseLogger) recover(fn func()) {

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	fn()

}

//msgConsumer 消费消息
func (base *baseLogger) msgConsumer() {

	for {
		base.recover(func() {
			select {
			case msg := <-base.infoCh:
				//格式化消息体
				base.msgFormat(&base.infoBuf, msg)
				//写入or展示消息
				base.writer(INFO, base.infoBuf)
				//回收对象
				base.releaseLogMsg(msg)
			case msg := <-base.debugCh:
				//格式化消息体
				base.msgFormat(&base.debugBuf, msg)
				//写入or展示消息
				base.writer(DEBUG, base.debugBuf)
				//回收对象
				base.releaseLogMsg(msg)
			case msg := <-base.errorCh:
				//格式化消息体
				base.msgFormat(&base.errorBuf, msg)
				//写入or展示消息
				base.writer(ERROR, base.errorBuf)
				//回收对象
				base.releaseLogMsg(msg)
			}

		})

	}
}

//msgFormat 消息格式化
func (base *baseLogger) msgFormat(by *[]byte, msg logMsg) {
	*by = (*by)[:0]
	if base.format == JsonFormat {
		msg.JsonFormat(by)
		*by = append(*by, "\n"...)
		return
	}

	msg.StrFormat(by)
	return

}
