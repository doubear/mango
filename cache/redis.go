package cache

import "gopkg.in/redis.v5"
import "github.com/go-mango/mango"
import "time"

//RedisCacher redis cache driver.
type RedisCacher struct {
	client *redis.Client
}

//RedisOption configures redis client.
type RedisOption redis.Options

//Get retrieves value from cacher.
func (c *RedisCacher) Get(id string) interface{} {
	var v interface{}
	err := c.client.Get(id).Scan(v)
	if err != nil {
		mango.Warn(err.Error())
		return nil
	}

	return v
}

//Set stores value into cacher.
func (c *RedisCacher) Set(id string, value interface{}, ttl time.Duration) {
	c.client.Set(id, value, ttl)
}

//Del deletes cached value by id.
func (c *RedisCacher) Del(id string) {
	c.client.Del(id)
}

//Push pushs value to queue.
func (c *RedisCacher) Push(id string, value interface{}) {
	c.client.LPush(id, value)
}

//Pop pops value from queue.
func (c *RedisCacher) Pop(id string) interface{} {
	var v interface{}
	err := c.client.RPop(id).Scan(v)
	if err != nil {
		return nil
	}

	return v
}

//Flush clear all data that stored in current db.
func (c *RedisCacher) Flush() {
	c.client.FlushDb()
}

//GC ignored func.
func (c *RedisCacher) GC() {

}

//Redis create redis cache driver instance.
func Redis(opt RedisOption) mango.Cacher {
	r := redis.Options(opt)
	client := redis.NewClient(&r)
	return &RedisCacher{client}
}
