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

func (m *MyRedis) SaveVal(key string, val ...int64) error {
	for i, n := range val {
		if err := m.r.Set(m.ctx, key+strconv.Itoa(i+1), n, 0).Err(); err != nil {
			return err
		}
	}
	return nil
}

func (m *MyRedis) GetVal(key ...string) (map[string]int64, error) {
	result := map[string]int64{}
	for _, n := range key {
		val, err := m.r.Get(m.ctx, n).Result()
		if err != nil {
			return result, err
		}
		convNum, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return result, err
		}
		result[n] = convNum
	}
	return result, nil
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
