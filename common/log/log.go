package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is the global logger
var Logger *zap.Logger

// SugaredLogger is the global sugared logger, configured with AddCallerSkip
var SugaredLogger *zap.SugaredLogger

func init() {
	zcfg := zap.NewProductionConfig()
	zcfg.Encoding = "console"
	zcfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	zcfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	zcfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	zcfg.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	zcfg.Level.SetLevel(zapcore.DebugLevel)

	log, err := zcfg.Build(zap.AddStacktrace(zapcore.ErrorLevel))
	if err != nil {
		panic(err)
	}
	zap.RedirectStdLog(log)

	Logger = log
	SugaredLogger = Logger.WithOptions(zap.AddCallerSkip(1)).Sugar()
}

func Debug(v ...interface{}) {
	SugaredLogger.Debug(v...)
}

func Info(v ...interface{}) {
	SugaredLogger.Info(v...)
}

func Warn(v ...interface{}) {
	SugaredLogger.Warn(v...)
}

func Error(v ...interface{}) {
	SugaredLogger.Error(v...)
}

func Fatal(v ...interface{}) {
	SugaredLogger.Fatal(v...)
}

func Debugf(tmpl string, v ...interface{}) {
	SugaredLogger.Debugf(tmpl, v...)
}

func Infof(tmpl string, v ...interface{}) {
	SugaredLogger.Infof(tmpl, v...)
}

func Warnf(tmpl string, v ...interface{}) {
	SugaredLogger.Warnf(tmpl, v...)
}

func Errorf(tmpl string, v ...interface{}) {
	SugaredLogger.Errorf(tmpl, v...)
}

func Fatalf(tmpl string, v ...interface{}) {
	SugaredLogger.Fatalf(tmpl, v...)
}
