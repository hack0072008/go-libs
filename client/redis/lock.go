package redis

import "time"

const lockerPrefix = "locker"

// 通过 redis 加分布式锁
func (c *Client) Lock(key string, expireTime ...int) error {
	realKey := c.realKey(c.lockerKey(key))
	expire := 0 * time.Second
	if len(expireTime) > 0 {
		expire = time.Duration(expireTime[0]) * time.Second
	}
	if err := c.redisClient.SetNX(realKey, time.Now().Unix(), expire).Err(); err != nil {
		return err
	}
	return nil
}

// 解锁
func (c *Client) Unlock(key string) error {
	realKey := c.realKey(c.lockerKey(key))
	if err := c.redisClient.Del(realKey).Err(); err != nil {
		return err
	}
	return nil
}

func (c *Client) lockerKey(part string) string {
	return lockerPrefix + "." + part
}
