package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
	"tinyUrl/config"
)

var (
	logger    *zap.Logger
	sugLogger *zap.SugaredLogger
)

func GetLogger() *zap.SugaredLogger {
	if sugLogger == nil {
		initLog()
	}
	return sugLogger
}

//初始化zap日志
//log.level : "debug", "info", "warn", "error", "dpanic", "panic", and "fatal"

func initLog() {
	logLevel := config.Base.Log.Level
	logConsole := config.Base.Log.Console
	logPath := config.Base.Log.Dir

	_, err := os.Stat(logPath)
	if err != nil && os.IsNotExist(err) {
		os.MkdirAll(logPath, os.ModePerm)
	}

	encoderCfg := zapcore.EncoderConfig{
		TimeKey:       "ts",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stack",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeTime:    TimeEncoder,
		//EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})
	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)

	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logPath + "/server.log",
		MaxSize:    100, // megabytes
		MaxBackups: 3,
		MaxAge:     30, // days
	})

	fileEncoder := zapcore.NewJSONEncoder(encoderCfg)
	consoleEncoder := zapcore.NewConsoleEncoder(encoderCfg)

	log_level := zap.NewAtomicLevel()
	log_level.UnmarshalText([]byte(logLevel))
	fmt.Println("zap current log_level::", log_level.Level())

	var core zapcore.Core

	//console为true则开启控制台，控制台打印所有级别的日志
	if logConsole {
		core = zapcore.NewTee(
			zapcore.NewCore(fileEncoder, w, log_level),
			zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
			zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
		)
	} else {
		core = zapcore.NewTee(
			zapcore.NewCore(fileEncoder, w, log_level),
		)
	}
	//出错的时候打印堆栈
	logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	defer logger.Sync()
	sugLogger = logger.Sugar()

}

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}
