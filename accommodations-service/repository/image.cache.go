package repository

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"time"

	"github.com/go-redis/redis"
	"go.opentelemetry.io/otel/trace"
)

type ImageCache struct {
	cli    *redis.Client
	logger *log.Logger
	tracer trace.Tracer
}

func NewCache(logger *log.Logger, tracer trace.Tracer) *ImageCache {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisAddress := fmt.Sprintf("%s:%s", redisHost, redisPort)

	client := redis.NewClient(&redis.Options{
		Addr: redisAddress,
	})

	return &ImageCache{
		cli:    client,
		logger: logger,
		tracer: tracer,
	}
}

func (ic *ImageCache) Ping(ctx context.Context) {
	ctx, span := ic.tracer.Start(ctx, "ImageCache.Ping")
	defer span.End()
	val, _ := ic.cli.Ping().Result()
	ic.logger.Println(val)
}

func (ic *ImageCache) Post(ctx context.Context, file multipart.File, key string) error {
	ctx, span := ic.tracer.Start(ctx, "ImageCache.Post")
	defer span.End()
	data, err := io.ReadAll(file)
	if err != nil {
		ic.logger.Println("Convert error:", err)
		return err
	}
	err = ic.cli.Set(key, data, 30*time.Second).Err()

	return err
}

func (ic *ImageCache) Create(ctx context.Context, data []byte, key string) error {
	ctx, span := ic.tracer.Start(ctx, "ImageCache.Create")
	defer span.End()
	err := ic.cli.Set(key, data, 30*time.Second).Err()
	return err
}

func (ic *ImageCache) Get(ctx context.Context, key string) ([]byte, error) {
	ctx, span := ic.tracer.Start(ctx, "ImageCache.Get")
	defer span.End()
	data, err := ic.cli.Get(key).Bytes()
	if err != nil {
		ic.logger.Println("Error in opening file from cache:", err)
		return nil, err
	}
	ic.logger.Println("Cache hit")
	return data, nil
}

func (ic *ImageCache) GetAll(ctx context.Context, keys []string) ([][]byte, error) {
	ctx, span := ic.tracer.Start(ctx, "ImageCache.GetAll")
	defer span.End()
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
