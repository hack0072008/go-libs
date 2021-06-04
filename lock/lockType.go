package tryLock

import (
	"context"
	"time"
)

func NewHostGpuLock(ctx context.Context, hostId string, timeout time.Duration) Mutex {
	name := "GpuLock_" + hostId
	return NewTryLock(ctx, name, timeout)
}

func NewHostNicLock(ctx context.Context, hostId string, timeout time.Duration) Mutex {
	name := "NicLock_" + hostId
	return NewTryLock(ctx, name, timeout)
}
