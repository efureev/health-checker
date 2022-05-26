package checkers

import (
	"context"
	"fmt"
	"strings"

	"github.com/efureev/health-checker"
	"github.com/go-redis/redis/v8"
)

type RedisConfig struct {
	Host     string
	Password string
	Port     int
	DB       int
}

func getRedisInfo(rdb *redis.Client) (map[string]string, error) {
	cmd := rdb.Info(context.Background(), `server`)
	if cmd.Err() != nil {
		return nil, cmd.Err()
	}

	str := cmd.String()

	_, afterStr, _ := strings.Cut(str, "\r\n")
	strList := strings.Split(afterStr, "\r\n")

	var result = make(map[string]string)
	for _, line := range strList {
		if !strings.Contains(line, `:`) {
			continue
		}

		lineArr := strings.Split(line, `:`)
		result[lineArr[0]] = lineArr[1]
	}

	return result, nil
}

func CheckingRedisFn(config RedisConfig) checker.CmdFn {
	return func(result *checker.Result, log checker.ILogger) {
		rdb := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf(`%s:%d`, config.Host, config.Port),
			Password: config.Password, // no password set
			DB:       config.DB,       // use default DB
		})

		redisInfo, err := getRedisInfo(rdb)
		if err != nil {
			result.AddError(err)
			return
		}

		result.Info.Version = fmt.Sprintf(`v%s`, redisInfo[`redis_version`])
		return
	}
}
