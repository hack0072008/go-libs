package log

import (
	"os"
	"strings"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger(path string, level string) {
	logLevel := zapcore.InfoLevel
	switch strings.ToLower(level) {
	case "debug":
		logLevel = zapcore.DebugLevel
		break
	case "warn":
		logLevel = zapcore.WarnLevel
		break
	case "error":
		logLevel = zapcore.ErrorLevel
		break
	case "fatal":
		logLevel = zapcore.FatalLevel
		break
	default:

	}

	sync := zapcore.AddSync(os.Stdout)
	if len(path) > 0 {
		lj := &lumberjack.Logger{
			Filename:   path,
			MaxSize:    200,
			MaxAge:     30,
			MaxBackups: 10,
		}
		sync = zapcore.AddSync(lj)
	}

	core := zapcore.NewCore(GetEncoder(), sync, logLevel)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	zap.ReplaceGlobals(logger)
}

func GetEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func Debug(args ...interface{}) {
	zap.S().Debug(args...)
}

// Info uses fmt.Sprint to construct and log a message.
func Info(args ...interface{}) {
	zap.S().Info(args...)
}

// Warn uses fmt.Sprint to construct and log a message.
func Warn(args ...interface{}) {
	zap.S().Warn(args...)
}

// Error uses fmt.Sprint to construct and log a message.
func Error(args ...interface{}) {
	zap.S().Error(args...)
}

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
func Fatal(args ...interface{}) {
	zap.S().Fatal(args...)
}

// Debugf uses fmt.Sprintf to log a templated message.
func Debugf(template string, args ...interface{}) {
	zap.S().Debugf(template, args...)
}

// Infof uses fmt.Sprintf to log a templated message.
func Infof(template string, args ...interface{}) {
	zap.S().Infof(template, args...)
}

// Warnf uses fmt.Sprintf to log a templated message.
func Warnf(template string, args ...interface{}) {
	zap.S().Warnf(template, args...)
}

// Errorf uses fmt.Sprintf to log a templated message.
func Errorf(template string, args ...interface{}) {
	zap.S().Errorf(template, args...)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func Fatalf(template string, args ...interface{}) {
	zap.S().Fatalf(template, args...)
}
