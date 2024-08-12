package repository

import (
	"Marketplace/internal/entities"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Masterminds/squirrel"
	"time"
)

func (ad *Repository) CreateAdvert(ctx context.Context, advert *entities.Advert) (*entities.Advert, error) {
	const op = "db.CreateAdvert"
	var lastId int

	query := "INSERT INTO adverts (user_id, header, text, address, image_url, price, datetime) " +
		"VALUES ($1,$2, $3, $4, $5, $6, $7) RETURNING id"
	err := ad.db.QueryRowContext(
		ctx, query, advert.UserId, advert.Header, advert.Text, advert.Address, advert.ImageURL, advert.Price, advert.Datetime).Scan(&lastId)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	advert.Id = lastId

	advertData, _ := json.Marshal(advert)
	ad.rdb.Set(ctx, fmt.Sprintf("advert:%d", lastId), advertData, 10*time.Minute)

	return advert, nil
}

func (ad *Repository) GetAdvert(ctx context.Context, adId int) (*entities.Advert, error) {
	const op = "db.GetAdvert"
	advert := &entities.Advert{}

	cachedAdvert, err := ad.rdb.Get(ctx, fmt.Sprintf("user:%d", adId)).Result()
	if err == nil {
		if err := json.Unmarshal([]byte(cachedAdvert), &advert); err == nil {
			return advert, nil
		}
	}

	query := "SELECT * FROM adverts WHERE id = $1"
	err = ad.db.QueryRowContext(ctx, query, adId).Scan(
		&advert.Id, &advert.UserId, &advert.Header, &advert.Text, &advert.Address, &advert.ImageURL, &advert.Price, &advert.Datetime)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	advertData, _ := json.Marshal(advert)
	ad.rdb.Set(ctx, fmt.Sprintf("advert:%d", adId), advertData, 10*time.Minute)

	return advert, nil
}

func (ad *Repository) UpdateAdvert(ctx context.Context, advert *entities.Advert) (*entities.Advert, error) {
	const op = "db.UpdateAdvert"
	queryBuilder := squirrel.Update("adverts").Where(squirrel.Eq{"id": advert.Id})

	if advert.Header != "" {
		queryBuilder = queryBuilder.Set("header", advert.Header)
	}
	if advert.Text != "" {
		queryBuilder = queryBuilder.Set("text", advert.Text)
	}
	if advert.ImageURL != "" {
		queryBuilder = queryBuilder.Set("image_url", advert.ImageURL)
	}
	if advert.Address != "" {
		queryBuilder = queryBuilder.Set("address", advert.Address)
	}
	if advert.Price != 0 {
		queryBuilder = queryBuilder.Set("price", advert.Price)
	}
	if !advert.Datetime.IsZero() {
		queryBuilder = queryBuilder.Set("datetime", advert.Datetime)
	}
	queryBuilder = queryBuilder.Set("by_this_user", advert.ByThisUser)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = ad.db.Exec(query, args...)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	advertResult, err := ad.GetAdvert(ctx, advert.Id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	advertData, _ := json.Marshal(advertResult)
	ad.rdb.Set(ctx, fmt.Sprintf("user:%d", advert.Id), advertData, 10*time.Minute)

	return ad.GetAdvert(ctx, advert.Id)
}

func (ad *Repository) DeleteAdvert(ctx context.Context, advertId int) error {
	const op = "db.DeleteAdvert"
	query := "DELETE FROM adverts WHERE id = $1"
	_, err := ad.db.ExecContext(ctx, query, advertId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	ad.rdb.Del(ctx, fmt.Sprintf("advert:%d", advertId))

	return nil
}
