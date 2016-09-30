package klog

import (
    "github.com/jcelliott/lumber"
    "grpc_center/module/env"
)

type LoggerInterface interface {
    AddLoggers(...lumber.Logger)
    Trace(string, ...interface{})
    Debug(string, ...interface{})
    Info(string, ...interface{})
    Warn(string, ...interface{})
    Error(string, ...interface{})
    Fatal(string, ...interface{})
}

var Logger LoggerInterface

func NewLogger(needDebug bool, filePath string, maxLine int, backups int) LoggerInterface {


    logFile, err := lumber.NewFileLogger(
        filePath,
        lumber.INFO,
        lumber.ROTATE,
        maxLine,
        backups,
        256,
    )

    if err != nil {
        panic(err)
    }

    var logger LoggerInterface
    if needDebug {
        logger = &DebugLogger{
            multiLogger: lumber.NewMultiLogger(),
        }
        logConsole := lumber.NewConsoleLogger(lumber.TRACE)
        logger.AddLoggers(logFile, logConsole)
    } else {
        logger = lumber.NewMultiLogger()
        logger.AddLoggers(logFile)
    }
    return logger
}

func init() {
    envConf := env.Conf.Env
    Logger = NewLogger(env.IsDev, envConf.Log_file, envConf.Log_max_line, envConf.Log_backups)
}

