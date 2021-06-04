package log

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UidLog struct {
	uuid string
}

func NewUidLog(ctx context.Context) UidLog {
	id := ctx.Value("uuid")
	if uid, ok := id.(string); ok {
		return UidLog{uuid: fmt.Sprintf("[%s] ", uid)}
	}
	return UidLog{uuid: fmt.Sprintf("[%s] ", uuid.New().String())}
}

func (t UidLog) GetUUID() string {
	return t.uuid[1 : len(t.uuid)-1]
}

func (t UidLog) Debug(args ...interface{}) {
	args = append([]interface{}{t.uuid}, args...)
	zap.S().Debug(args...)
}

// Info uses fmt.Sprint to construct and log a message.
func (t UidLog) Info(args ...interface{}) {
	args = append([]interface{}{t.uuid}, args...)
	zap.S().Info(args...)
}

// Warn uses fmt.Sprint to construct and log a message.
func (t UidLog) Warn(args ...interface{}) {
	args = append([]interface{}{t.uuid}, args...)
	zap.S().Warn(args...)
}

// Error uses fmt.Sprint to construct and log a message.
func (t UidLog) Error(args ...interface{}) {
	args = append([]interface{}{t.uuid}, args...)
	zap.S().Error(args...)
}

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
func (t UidLog) Fatal(args ...interface{}) {
	args = append([]interface{}{t.uuid}, args...)
	zap.S().Fatal(args...)
}

// Debugf uses fmt.Sprintf to log a templated message.
func (t UidLog) Debugf(template string, args ...interface{}) {
	zap.S().Debugf(t.uuid+template, args...)
}

// Infof uses fmt.Sprintf to log a templated message.
func (t UidLog) Infof(template string, args ...interface{}) {
	zap.S().Infof(t.uuid+template, args...)
}

// Warnf uses fmt.Sprintf to log a templated message.
func (t UidLog) Warnf(template string, args ...interface{}) {
	zap.S().Warnf(t.uuid+template, args...)
}

// Errorf uses fmt.Sprintf to log a templated message.
func (t UidLog) Errorf(template string, args ...interface{}) {
	zap.S().Errorf(t.uuid+template, args...)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func (t UidLog) Fatalf(template string, args ...interface{}) {
	zap.S().Fatalf(t.uuid+template, args...)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func (t UidLog) Printf(template string, args ...interface{}) {
	zap.S().Infof(t.uuid+template, args...)
}
