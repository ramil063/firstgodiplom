package logger

import "go.uber.org/zap"

var logLevel = "INFO"

// Log будет доступен всему коду как синглтон.
// Никакой код навыка, кроме функции Initialize, не должен модифицировать эту переменную.
// По умолчанию установлен no-op-логер, который не выводит никаких сообщений.
var Log *zap.Logger = zap.NewNop()

func WriteInfoLog(message string) {
	Log.Info(message)
}

func WriteDebugLog(message string) {
	Log.Debug(message)
}

func WriteErrorLog(message string) {
	Log.Error(message)
}
