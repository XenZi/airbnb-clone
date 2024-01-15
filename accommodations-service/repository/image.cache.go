package repository

import (
	"fmt"
	"github.com/go-redis/redis"
	"io"
	"log"
	"mime/multipart"
	"os"
	"time"
)

type ImageCache struct {
	cli    *redis.Client
	logger *log.Logger
}

func NewCache(logger *log.Logger) *ImageCache {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisAddress := fmt.Sprintf("%s:%s", redisHost, redisPort)

	client := redis.NewClient(&redis.Options{
		Addr: redisAddress,
	})

	return &ImageCache{
		cli:    client,
		logger: logger,
	}
}

func (ic *ImageCache) Ping() {
	val, _ := ic.cli.Ping().Result()
	ic.logger.Println(val)
}

func (ic *ImageCache) Post(file multipart.File, key string) error {
	data, err := io.ReadAll(file)
	if err != nil {
		ic.logger.Println("Convert error:", err)
		return err
	}
	err = ic.cli.Set(key, data, 30*time.Second).Err()

	return err
}

func (ic *ImageCache) Create(data []byte, key string) error {
	err := ic.cli.Set(key, data, 30*time.Second).Err()
	return err
}

func (ic *ImageCache) Get(key string) ([]byte, error) {
	data, err := ic.cli.Get(key).Bytes()
	if err != nil {
		ic.logger.Println("Error in opening file from cache:", err)
		return nil, err
	}
	ic.logger.Println("Cache hit")
	return data, nil
}

func (ic *ImageCache) GetAll(keys []string) ([][]byte, error) {
	var values [][]byte
	for _, key := range keys {
		value, err := ic.cli.Get(key).Bytes()
		if err != nil {
			ic.logger.Println("Error in opening file from cache:", err)
			return nil, err
		}
		values = append(values, value)
	}

	ic.logger.Println("Cache hit")
	return values, nil
}
