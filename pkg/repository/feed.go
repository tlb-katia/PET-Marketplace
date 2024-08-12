package repository

import (
	"Marketplace/internal/entities"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/go-redis/redis/v8"
	"time"
)

func (ad *Repository) GetSorted(ctx context.Context, filter *entities.Filter) (*entities.AdvList, error) {
	const op = "repository.GetSorted"
	rowsSlice := make([]entities.Advert, 0)

	cacheKey := GenerateCacheKey(filter)

	cachedResult, err := ad.rdb.Get(ctx, cacheKey).Result()

	if errors.Is(err, redis.Nil) {
		queryBuilder := sq.Select("*").From("adverts")
		if filter.AscendingOrder && filter.ByPrice {
			queryBuilder.OrderBy("price")
		}
		if filter.MinPrice != 0.0 {
			queryBuilder.Where("price >= ?", filter.MinPrice)
		}
		if filter.MaxPrice != 0.0 {
			queryBuilder.Where("price <= ?", filter.MaxPrice)
		}

		query, args, err := queryBuilder.ToSql()
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		rows, err := ad.db.QueryContext(ctx, query, args...)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		for rows.Next() {
			adv := &entities.Advert{}
			rows.Scan(&adv.Id, &adv.UserId, &adv.Header, &adv.Text, &adv.Address,
				&adv.ImageURL, &adv.Price, &adv.Datetime)
			//fmt.Println(adv)
			rowsSlice = append(rowsSlice, *adv)
		}
		resultData, _ := json.Marshal(rowsSlice)
		err = ad.rdb.Set(ctx, cacheKey, resultData, 10*time.Minute).Err()
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("%s: failed to get data from Redis: %w", op, err)
	} else {
		err = json.Unmarshal([]byte(cachedResult), &rowsSlice)
		if err != nil {
			return nil, fmt.Errorf("%s: failed to unmarshal cached data: %w", op, err)
		}
	}

	return &entities.AdvList{List: rowsSlice}, nil
}

func GenerateCacheKey(filter *entities.Filter) string {
	hash := sha1.New()
	hash.Write([]byte(fmt.Sprintf("%v", filter)))
	return hex.EncodeToString(hash.Sum(nil))
}
