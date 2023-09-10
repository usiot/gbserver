package logger

import (
	"fmt"

	"github.com/ghettovoice/gosip/log"
	"github.com/inysc/hog"
)

type GoSipLogger struct{ hog.Event }

var Gsl log.Logger = GoSipLogger{}

func (l GoSipLogger) e(lvl uint8) hog.Event {
	if l.Event == nil {
		return lg.NewEvent(4, lvl, true)
	}
	return l.Event
}

func (l GoSipLogger) Print(args ...interface{})                 { l.e(hog.INFO).Msg(fmt.Sprint(args...)) }
func (l GoSipLogger) Printf(format string, args ...interface{}) { l.e(hog.INFO).Msgf(format, args...) }
func (l GoSipLogger) Trace(args ...interface{})                 { l.e(hog.TRACE).Msg(fmt.Sprint(args...)) }
func (l GoSipLogger) Tracef(format string, args ...interface{}) { l.e(hog.TRACE).Msgf(format, args...) }
func (l GoSipLogger) Debug(args ...interface{})                 { l.e(hog.DEBUG).Msg(fmt.Sprint(args...)) }
func (l GoSipLogger) Debugf(format string, args ...interface{}) { l.e(hog.DEBUG).Msgf(format, args...) }
func (l GoSipLogger) Info(args ...interface{})                  { l.e(hog.INFO).Msg(fmt.Sprint(args...)) }
func (l GoSipLogger) Infof(format string, args ...interface{})  { l.e(hog.INFO).Msgf(format, args...) }
func (l GoSipLogger) Warn(args ...interface{})                  { l.e(hog.WARN).Msg(fmt.Sprint(args...)) }
func (l GoSipLogger) Warnf(format string, args ...interface{})  { l.e(hog.WARN).Msgf(format, args...) }
func (l GoSipLogger) Error(args ...interface{})                 { l.e(hog.ERROR).Msg(fmt.Sprint(args...)) }
func (l GoSipLogger) Errorf(format string, args ...interface{}) { l.e(hog.ERROR).Msgf(format, args...) }
func (l GoSipLogger) Fatal(args ...interface{})                 { l.e(hog.FATAL).Msg(fmt.Sprint(args...)) }
func (l GoSipLogger) Fatalf(format string, args ...interface{}) { l.e(hog.FATAL).Msgf(format, args...) }
func (l GoSipLogger) Panic(args ...interface{})                 { l.e(hog.PANIC).Msg(fmt.Sprint(args...)) }
func (l GoSipLogger) Panicf(format string, args ...interface{}) { l.e(hog.PANIC).Msgf(format, args...) }

func (l GoSipLogger) WithPrefix(prefix string) log.Logger {
	return GoSipLogger{l.e(hog.INFO).String(prefix, " ")}
}
func (l GoSipLogger) Prefix() string { return "" }

func (l GoSipLogger) WithFields(fields map[string]interface{}) log.Logger {
	e := l.e(hog.INFO)
	for k, v := range fields {
		e.Any(k, v)
	}
	return l
}

func (l GoSipLogger) Fields() log.Fields    { return nil }
func (l GoSipLogger) SetLevel(level uint32) {}
