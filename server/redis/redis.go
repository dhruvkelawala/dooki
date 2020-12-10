package sider

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"os"
	"time"
)

var rc *redis.Client
var c = context.Background()

func init() {
	rUrl, ok := os.LookupEnv("REDIS_URL")
	if ok {
		opt, err := redis.ParseURL(rUrl)
		if err != nil {
			log.Println(err)
		}

		rc = redis.NewClient(opt)
	}

	rc = redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		Password: "",
	})
}


func SetData(key string, message interface{}, expiry time.Duration) error {
	return rc.Set(c, key, message, expiry).Err()
}


func GetDataJson(key string) ([]byte, error) {
	res, err := rc.Get(c, key).Result()
	if err != nil {
		return nil, err
	}


	return []byte(res),nil
}

func IsRoomInDb(roomID string) bool {
	if _, err := GetDataJson(roomID); err != nil {
		return false
	}

	return true
}

func PublishToChannel(channel string, message interface{}) error {
	return rc.Publish(c, channel, message).Err()
}

func SubscribeToChannel(channel string) *redis.PubSub {
	return rc.Subscribe(c, channel)
}

func UnsubscribeChannel(sub *redis.PubSub, channel string) error {
	if err := sub.Unsubscribe(c, channel); err != nil {
		return err
	}

	return nil
}