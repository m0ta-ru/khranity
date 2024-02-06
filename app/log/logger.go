package log

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"khranity/app/config"
)

// Logger is a logger :)
type Logger struct {
	*zap.Logger
}

var (
	logger *Logger
	once   sync.Once
)

// Get inits logger with log level from environment. Once.
func Get() *Logger {
	once.Do(func() {
		cfg := config.Get()
		pe := zap.NewProductionEncoderConfig()
		je := zap.NewProductionEncoderConfig()
		level := zapcore.Level(zap.InfoLevel)
		switch cfg.LogLevel {
		case "debug":
			level = zapcore.Level(zap.DebugLevel)
		case "info":
			level = zapcore.Level(zap.InfoLevel)
		case "warn", "warning":
			level = zapcore.Level(zap.WarnLevel)
		case "err", "error":
			level = zapcore.Level(zap.ErrorLevel)
		case "fatal":
			level = zapcore.Level(zap.FatalLevel)
		case "panic":
			level = zapcore.Level(zap.PanicLevel)
		default:
			level = zapcore.Level(zap.InfoLevel)
		}

		pe.EncodeTime = zapcore.ISO8601TimeEncoder
		pe.EncodeLevel = zapcore.CapitalColorLevelEncoder
		pe.EncodeCaller = zapcore.FullCallerEncoder

		je.EncodeTime = zapcore.ISO8601TimeEncoder
		je.EncodeLevel = zapcore.CapitalLevelEncoder
		je.EncodeCaller = zapcore.FullCallerEncoder

		consoleEncoder := zapcore.NewConsoleEncoder(pe)
		jsonEncoder := zapcore.NewJSONEncoder(je)
		errorEncoder := zapcore.NewJSONEncoder(je)

		path := filepath.Join(".", cfg.LogFolder)

		infoWS := zapcore.AddSync(&lumberjack.Logger{
			Filename:   fmt.Sprintf("%s/%s-info.log", path, cfg.LogPrefix), //"./logs/info.log",
			MaxSize:    500,                                                // megabytes
			MaxBackups: 3,
			MaxAge:     30, // days
		})
		errWS := zapcore.AddSync(&lumberjack.Logger{
			Filename:   fmt.Sprintf("%s/%s-error.log", path, cfg.LogPrefix), //"./logs/error.log",
			MaxSize:    500,                                                 // megabytes
			MaxBackups: 3,
			MaxAge:     1, // days
		})

		core := zapcore.NewTee(
			zapcore.NewCore(errorEncoder, errWS, zap.ErrorLevel),
			zapcore.NewCore(jsonEncoder, infoWS, level),
			zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
		)
		zapLogger := zap.New(core, zap.AddCaller())
		logger = &Logger{zapLogger}
	})
	return logger
}

func Start() {
	logger = Get()
	defer logger.Logger.Sync()
}

func GetDefault() *Logger {
	once.Do(func() {
		//cfg 	:= config.Get()
		pe := zap.NewProductionEncoderConfig()
		//je 		:= zap.NewProductionEncoderConfig()
		level := zapcore.Level(zap.InfoLevel)

		pe.EncodeTime = zapcore.ISO8601TimeEncoder
		pe.EncodeLevel = zapcore.CapitalColorLevelEncoder
		pe.EncodeCaller = zapcore.FullCallerEncoder

		// je.EncodeTime 	= zapcore.ISO8601TimeEncoder
		// je.EncodeLevel 	= zapcore.CapitalLevelEncoder
		// je.EncodeCaller = zapcore.FullCallerEncoder

		consoleEncoder := zapcore.NewConsoleEncoder(pe)
		// jsonEncoder		:= zapcore.NewJSONEncoder(je)
		// errorEncoder	:= zapcore.NewJSONEncoder(je)

		// path := filepath.Join(".", cfg.LogFolder)

		// infoWS := zapcore.AddSync(&lumberjack.Logger{
		// 	Filename:   fmt.Sprintf("%s/%s-info.log", path, cfg.LogPrefix),		//"./logs/info.log",
		// 	MaxSize:    500, // megabytes
		// 	MaxBackups: 3,
		// 	MaxAge:     30, // days
		// })
		// errWS := zapcore.AddSync(&lumberjack.Logger{
		// 	Filename:   fmt.Sprintf("%s/%s-error.log", path, cfg.LogPrefix), 	//"./logs/error.log",
		// 	MaxSize:    500, // megabytes
		// 	MaxBackups: 3,
		// 	MaxAge:     1, // days
		// })

		core := zapcore.NewTee(
			// zapcore.NewCore(errorEncoder, errWS, zap.ErrorLevel),
			// zapcore.NewCore(jsonEncoder, infoWS, level),
			zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
		)
		zapLogger := zap.New(core, zap.AddCaller())
		logger = &Logger{zapLogger}
	})
	return logger
}
