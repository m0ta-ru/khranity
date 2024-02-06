package log

func Debug(msg string, fields ...Field) {
	logger.Debug(msg, fields...)
}

func Debugf(format string, args ...interface{}) {
	logger.Sugar().Debugf(format, args...)
}

func Info(msg string, fields ...Field) {
	logger.Info(msg, fields...)
}

func Infof(format string, args ...interface{}) {
	logger.Sugar().Infof(format, args...)
}

func Warn(msg string, fields ...Field) {
	logger.Warn(msg, fields...)
}

func Warnf(format string, args ...interface{}) {
	logger.Sugar().Warnf(format, args...)
}

func Error(msg string, fields ...Field) {
	logger.Error(msg, fields...)
}

func Errorf(format string, args ...interface{}) {
	logger.Sugar().Errorf(format, args...)
}

func Panic(msg string, fields ...Field) {
	logger.Panic(msg, fields...)
}

func Panicf(format string, args ...interface{}) {
	logger.Sugar().Panicf(format, args...)
}

func Fatal(msg string, fields ...Field) {
	logger.Fatal(msg, fields...)
}

func Fatalf(format string, args ...interface{}) {
	logger.Sugar().Fatalf(format, args...)
}
