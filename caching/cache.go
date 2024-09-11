package caching

import (
	"context"
	"encoding/json"
	"proxy/client"
	"proxy/model"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepo interface {
	SaveCache(c context.Context, request string, respone interface{}, config model.ConfigSettings) error
	CheackCache(c context.Context, request string) (interface{}, error)
	UpdateCache(c context.Context, request string, config model.ConfigSettings) error
}

type redisImpl struct {
	Rdb *redis.Client
}

func NewRedisRepo(rdb *redis.Client) RedisRepo {
	return &redisImpl{
		Rdb: rdb,
	}
}

func ConnectRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func (R *redisImpl) SaveCache(c context.Context, request string, respone interface{}, config model.ConfigSettings) error {
	err := R.Rdb.Set(c, request, respone, time.Duration(config.TimeDuration)).Err()
	return err
}

func (R *redisImpl) CheackCache(c context.Context, request string) (interface{}, error) {
	val, err := R.Rdb.Get(c, request).Result()
	if err != nil {
		return nil, err
	}

	var resp interface{}
	err = json.Unmarshal([]byte(val), &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (R *redisImpl) UpdateCache(c context.Context, request string, config model.ConfigSettings) error {
	for {
		time.Sleep(time.Duration(config.RefreshTime))
		keys, err := R.Rdb.Keys(c, "*").Result()
		if err != nil{
			return nil
		}

		for _, req := range keys{
			resp, err := client.ForwardBackend(req)
			if err != nil{
				return nil
			}
			err = R.SaveCache(c, req, resp, config)
			if err != nil{
				return nil
			}
		}
	}
}
