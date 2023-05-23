package stack

import (
	"github.com/go-redis/redis"
)

// RedisStack is a Redis backed Stack
type RedisStack struct {
	r   *redis.Client
	key string
	cap int64
}

// NewRedisStack returns a new redis-backed stack which may be used to perform
// stack operations. Returns error if unable to ping database.
// url: Redis URL
// password: Redis Password
// dbName: Redis database name
// key: Key inside Redis to store stack in
// cap: Capacity of stack
func NewRedisStack(url, password string, dbName int, key string, cap int64) (db Stack, err error) {
	redisOpts := redis.Options{
		Addr:     url,
		Password: password,
		DB:       dbName,
	}
	rs := RedisStack{
		r:   redis.NewClient(&redisOpts),
		key: key,
		cap: cap - 1, // Convert index 1 to index 0
	}
	_, err = rs.r.Ping().Result()
	db = rs
	return
}

// Push pushes message into Stack
func (s RedisStack) Push(msg string) (err error) {
	pRes := s.r.LPush(s.key, msg)
	err = pRes.Err()
	if err != nil {
		return
	}
	tRes := s.r.LTrim(s.key, 0, s.cap)
	err = tRes.Err()
	return
}

// Read reads all messages from stack
func (s RedisStack) Read() ([]string, error) {
	res := s.r.LRange(s.key, 0, s.cap)
	return res.Val(), res.Err()
}
