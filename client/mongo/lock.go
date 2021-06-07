package mongoClient

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/hack0072008/go-libs/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var pool *sync.Pool

type Mutex struct {
	log     log.UidLog
	ctx     context.Context
	name    string
	timeout int64
}

/*
GpuLock
*/
func NewHostGpuLock(ctx context.Context, hostId string, timeout time.Duration) Mutex {
	name := "GpuLock_" + hostId
	return CreateLock(ctx, name, timeout)
}

/*
NicLock
*/
func NewHostNicLock(ctx context.Context, hostId string, timeout time.Duration) Mutex {
	name := "NicLock_" + hostId
	return CreateLock(ctx, name, timeout)
}

/*
创建一个分布式锁 timeout 锁的超时时间 name 锁名
*/
func CreateLock(ctx context.Context, name string, timeout time.Duration) Mutex {
	if pool == nil {
		gen := func() interface{} {
			client := GetClient()
			if client == nil {
				log.Warn("mongoClient get client is nil")
				return nil
			}
			return client.Database("uhost").Collection("UpdateLock")
		}
		pool = &sync.Pool{New: gen}
	}
	return Mutex{
		log:     log.NewUidLog(ctx),
		name:    name,
		ctx:     ctx,
		timeout: int64(timeout.Seconds()),
	}
}

/*
尝试获取锁
*/
func (m Mutex) GetLock() (bool, error) {
	col, ok := pool.Get().(*mongo.Collection)
	if !ok || col == nil {
		m.log.Warn("mongoClient get client error")
		return false, errors.New("mongoClient get client error")
	}
	defer pool.Put(col)
	filter := bson.M{
		"id":          m.name,
		"exec_status": "free",
		"uuid":        m.log.GetUUID(),
	}
	doc := bson.M{"$set": bson.M{
		"id": m.name,
		"$or": []bson.M{
			{"exec_status": "buy"},
			{"last_exec_time": bson.M{"$lt": time.Now().Unix() - m.timeout}},
		},
	}}
	opt := options.FindOneAndUpdate()
	opt.SetUpsert(true)
	result := col.FindOneAndUpdate(m.ctx, filter, doc, opt)
	if result != nil && result.Err() != nil {
		if result.Err() != mongo.ErrNoDocuments {
			return false, nil
		}
		m.log.Info("insert lock")
	}
	return true, nil
}

/**
 * @param interval 如果没有获取到，间隔interval时间段后再去尝试获取
 * @param timeout 如果超过timeout时间段后还没获取到锁，返回超时错误
 */
func (m Mutex) Lock(interval, timeout time.Duration) error {
	tick := time.After(timeout)
	for {
		select {
		case <-tick:
			m.log.Error("lock timeout")
			return errors.New("lock timeout")
		default:
			ok, err := m.GetLock()
			if err != nil {
				m.log.Error(err)
				return err
			}
			if ok {
				return nil
			}
			m.log.Warn("try lock failed")
			time.Sleep(interval)
		}
	}
}

func (m Mutex) Unlock() error {
	col, ok := pool.Get().(*mongo.Collection)
	if !ok || col == nil {
		m.log.Warn("mongoClient get client error")
		return errors.New("mongoClient get client error")
	}
	defer pool.Put(col)
	filter := bson.M{
		"id":   m.name,
		"uuid": m.log.GetUUID(),
	}
	opt := options.FindOneAndDelete()
	result := col.FindOneAndDelete(m.ctx, filter, opt)
	if result != nil && result.Err() != nil {
		m.log.Error(result.Err().Error())
		return result.Err()
	}
	return nil
}
