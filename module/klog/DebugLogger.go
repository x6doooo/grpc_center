package klog

import (
    "github.com/jcelliott/lumber"
    "fmt"
    "runtime"
    "path/filepath"
)

type DebugLogger struct {
    multiLogger *lumber.MultiLogger
}

func getCallerInfo() string {
    _, filename, lineNumber, _ := runtime.Caller(2)
    filename = filepath.Base(filename)
    return fmt.Sprintf("<%s:%d>", filename, lineNumber)
}

func (me *DebugLogger) AddLoggers(newLogs ...lumber.Logger) {
    me.multiLogger.AddLoggers(newLogs...)
}

func (me *DebugLogger) Trace(format string, v ...interface{}) {
    info := getCallerInfo()
    me.multiLogger.Trace(info + " " + format, v...)
}
func (me *DebugLogger) Debug(format string, v ...interface{}) {
    info := getCallerInfo()
    me.multiLogger.Debug(info + " " + format, v...)
}
func (me *DebugLogger) Info(format string, v ...interface{}) {
    info := getCallerInfo()
    me.multiLogger.Info(info + " " + format, v...)
}
func (me *DebugLogger) Warn(format string, v ...interface{}) {
    info := getCallerInfo()
    me.multiLogger.Warn(info + " " + format, v...)
}
func (me *DebugLogger) Error(format string, v ...interface{}) {
    info := getCallerInfo()
    me.multiLogger.Error(info + " " + format, v...)
}
func (me *DebugLogger) Fatal(format string, v ...interface{}) {
    info := getCallerInfo()
    me.multiLogger.Fatal(info + " " + format, v...)
}

