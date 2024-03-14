package log

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Field = zapcore.Field

func Int(key string, value int) Field {
	return zap.Int(key, value)
}

func Int8(key string, value int8) Field {
	return zap.Int8(key, value)
}

func Int16(key string, value int16) Field {
	return zap.Int16(key, value)
}

func Int32(key string, value int32) Field {
	return zap.Int32(key, value)
}

func Int64(key string, value int64) Field {
	return zap.Int64(key, value)
}

func Float32(key string, value float32) Field {
	return zap.Float32(key, value)
}

func Float64(key string, value float64) Field {
	return zap.Float64(key, value)
}

func Duration(key string, value time.Duration) Field {
	return zap.Duration(key, value)
}

func ByteString(key string, value []byte) Field {
	return zap.ByteString(key, value)
}

func String(key string, value string) Field {
	return zap.String(key, value)
}

func Object(key string, value zapcore.ObjectMarshaler) Field {
	return zap.Object(key, value)
}

func Strings(key string, value []string) Field {
	return zap.Strings(key, value)
}
