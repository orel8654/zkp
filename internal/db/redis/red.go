package redis

import (
	"context"
	"strconv"

	reds "github.com/redis/go-redis/v9"
)

type MyRedis struct {
	ctx context.Context
	r   *reds.Client
}

func (m *MyRedis) SaveVal(key string, val int64) error {
	if err := m.r.Set(m.ctx, key, val, 0).Err(); err != nil {
		return err
	}
	return nil
}

func (m *MyRedis) GetVal(key string) (int64, error) {
	val, err := m.r.Get(m.ctx, key).Result()
	if err != nil {
		return 0, err
	}
	i, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func client(address string) *reds.Client {
	return reds.NewClient(&reds.Options{
		Addr:     address,
		Password: "",
		DB:       0,
	})
}

func MyNewRedis(address string) *MyRedis {
	return &MyRedis{
		ctx: context.Background(),
		r:   client(address),
	}
}
