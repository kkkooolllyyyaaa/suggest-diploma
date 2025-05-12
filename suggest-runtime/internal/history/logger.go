package history

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	maxRequestsPerUser = 5
)

type QueryLogger struct {
	client *redis.Client
}

func NewQueryLogger(redisAddr string) QueryLogger {
	return QueryLogger{
		client: redis.NewClient(&redis.Options{
			Addr: redisAddr,
		}),
	}
}

func (rl *QueryLogger) LogRequest(userID string, query string) error {
	ctx := context.Background()
	now := time.Now()

	err := rl.client.ZAdd(ctx, userID, redis.Z{
		Score:  timeToScore(now),
		Member: query,
	}).Err()
	if err != nil {
		return fmt.Errorf("failed to log query: %v", err)
	}

	// Оставляем только последние maxRequestsPerUser записей
	err = rl.client.ZRemRangeByRank(ctx, userID, 0, -(maxRequestsPerUser + 1)).Err()
	if err != nil {
		return fmt.Errorf("failed to trim requests set: %v", err)
	}

	return nil
}

func (rl *QueryLogger) GetUserRequests(userID string) ([]QueryTimestamp, error) {
	ctx := context.Background()

	// Получаем запросы в порядке от новых к старым
	results, err := rl.client.ZRevRangeWithScores(ctx, userID, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get user queries: %v", err)
	}

	var requests []QueryTimestamp
	for _, result := range results {
		requests = append(requests, QueryTimestamp{
			Query: result.Member.(string),
			Time:  scoreToTime(result.Score),
		})
	}

	return requests, nil
}

var timeOffset = time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC)

func timeToScore(ts time.Time) float64 {
	return ts.Sub(timeOffset).Seconds()
}

func scoreToTime(score float64) time.Time {
	return timeOffset.Add(time.Duration(score * 1000 * 1000 * 1000))
}
