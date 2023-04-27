package redisclient

import (
	"backend/utils"
	"context"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var client *redis.Client

func init() {

	// load the environment variables
	err := godotenv.Load()
	if err != nil {
		utils.Logger.Infof(".env file not found, using OS ENV variables. Err: %s", err)
	}

	//get port from env
	redispassword, exists := os.LookupEnv("REDIS_PASSWORD")

	if !exists {
		utils.Logger.Error("REDIS_PASSWORD not defined in env")
	}

	//get the host from env
	redishost, exists := os.LookupEnv("REDIS_HOST")
	if !exists {
		utils.Logger.Error("REDIS_HOST not defined in env, defaulting to localhost")
		redishost = "localhost"
	}
	opt := &redis.Options{
		Addr:     redishost + ":6379",
		Password: redispassword,
		DB:       0,
	}
	client = redis.NewClient(opt)

	_, err = client.Ping(context.TODO()).Result()
	if err != nil {
		utils.Logger.Panicf("Error connecting to redis. Err: %s", err)
		return
	}

	utils.Logger.Info("Connected to redis")
}

func Set(key string, value interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := client.Set(ctx, key, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func Get(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	val, err := client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return val, nil
}

func SetJson(key string, path string, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := client.Do(ctx, "JSON.SET", key, path, value).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetJson(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	val, err := client.Do(ctx, "JSON.GET", key).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}
	return val.(string), nil
}
