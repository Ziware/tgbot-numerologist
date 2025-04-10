package database

import (
	"context"
	"encoding/json"
	"fmt"

	"tgbot-numerologist/objects"

	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client = nil

func InitRDB(host, port string) {
	rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", host, port),
	})
}

func SaveProfileToRedis(profile *objects.Profile) error {
	ctx := context.Background()

	data, err := json.Marshal(profile)
	if err != nil {
		return err
	}

	err = rdb.Set(ctx, profile.Username, data, 0).Err()
	return err
}

func GetProfileFromRedis(username string) (*objects.Profile, error) {
	ctx := context.Background()

	data, err := rdb.Get(ctx, username).Result()
	if err != nil {
		return nil, err
	}

	var profile objects.Profile
	err = json.Unmarshal([]byte(data), &profile)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}

func GetChatId(username string) (int64, error) {
	profile, err := GetProfileFromRedis(username)
	if err != nil {
		return 0, err
	}
	return profile.ChatID, nil
}

func UserExistsInRedis(username string) (bool, error) {
	ctx := context.Background()

	n, err := rdb.Exists(ctx, username).Result()
	if err != nil {
		return false, err
	}

	return n > 0, nil
}
