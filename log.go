package husky

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogConfig struct {
	Level string `toml:"level"`
	Ansi  bool   `toml:"ansi"`
}

var _LogIns *zap.Logger

func InitLog(config *LogConfig) {
	level, err := zapcore.ParseLevel(config.Level)
	if err != nil {
		level = zapcore.DebugLevel
	}
	atom := zap.NewAtomicLevelAt(level)
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	if config.Ansi {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	_LogIns = zap.New(
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.Lock(os.Stdout),
			atom,
		),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)
}

func Debug(ctx context.Context, args ...interface{}) {
	traceId := ctx.Value("TraceId").(string)
	_LogIns.Sugar().With(zap.String("TraceId", traceId)).Debugln(args...)
}

func Info(ctx context.Context, args ...interface{}) {
	traceId := ctx.Value("TraceId").(string)
	_LogIns.Sugar().With(zap.String("TraceId", traceId)).Infoln(args...)
}

func Warn(ctx context.Context, args ...interface{}) {
	traceId := ctx.Value("TraceId").(string)
	_LogIns.Sugar().With(zap.String("TraceId", traceId)).Warnln(args...)
}

func Error(ctx context.Context, args ...interface{}) {
	traceId := ctx.Value("TraceId").(string)
	_LogIns.Sugar().With(zap.String("TraceId", traceId)).Errorln(args...)
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	traceId := ctx.Value("TraceId").(string)
	_LogIns.Sugar().With(zap.String("TraceId", traceId)).Debugf(format, args...)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	traceId := ctx.Value("TraceId").(string)
	_LogIns.Sugar().With(zap.String("TraceId", traceId)).Infof(format, args...)
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	traceId := ctx.Value("TraceId").(string)
	_LogIns.Sugar().With(zap.String("TraceId", traceId)).Warnf(format, args...)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	traceId := ctx.Value("TraceId").(string)
	_LogIns.Sugar().With(zap.String("TraceId", traceId)).Errorf(format, args...)
}
