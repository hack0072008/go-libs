package redis

const counterPrefix = "counter"

/*
计数器自增
 */
func (c *Client) Incr(key string) (int64, error) {
	realKey := c.realKey(c.CounterKey(key))
	value, err := c.redisClient.Incr(realKey).Result()
	if err != nil {
		return 0, err
	}
	return value, nil
}

func (c *Client) IncrBy(key string, count int) (int64, error) {
	realKey := c.realKey(c.CounterKey(key))
	value, err := c.redisClient.IncrBy(realKey, int64(count)).Result()
	if err != nil {
		return 0, err
	}
	return value, nil
}

/*
计数器自减
 */
func (c *Client) Decr(key string) (int64, error) {
	realKey := c.realKey(c.CounterKey(key))
	value, err := c.redisClient.Decr(realKey).Result()
	if err != nil {
		return 0, err
	}
	return value, nil
}

func (c *Client) DecrBy(key string, count int) (int64, error) {
	realKey := c.realKey(c.CounterKey(key))
	value, err := c.redisClient.DecrBy(realKey, int64(count)).Result()
	if err != nil {
		return 0, err
	}
	return value, nil
}

/*
获取计数器当前值
 */
func (c *Client) GetCounter(key string) (int64, error) {
	return c.IncrBy(key, 0)
}

func (c *Client) CounterKey(part string) string {
	return counterPrefix + "." + part
}
