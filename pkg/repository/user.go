package repository

import (
	"Marketplace/internal/entities"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// TODO check hashed password

func (ad *Repository) CreateUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	const op = "db.CreateUser"
	var lastId int

	flag := ad.CheckUserExists(user.Login)
	if !flag {
		return nil, fmt.Errorf("user already exists")
	}

	query := "INSERT INTO users (login, password) VALUES ($1,$2) RETURNING id"
	err := ad.db.QueryRowContext(ctx, query, user.Login, user.Password).Scan(&lastId)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	fmt.Println(lastId)

	userResp := &entities.User{
		Id:       lastId,
		Login:    user.Login,
		Password: user.Password,
	}

	userData, _ := json.Marshal(userResp)
	ad.rdb.Set(ctx, fmt.Sprintf("user:%d", lastId), userData, 10*time.Minute)

	return userResp, nil
}

func (ad *Repository) LoginUser(ctx context.Context, user *entities.LoginReqUser) (*entities.User, error) {
	const op = "db.LoginUser"
	var password string
	var id int

	query := "SELECT id, password FROM users WHERE login = $1"
	err := ad.db.QueryRowContext(ctx, query, user.Login).Scan(&id, &password)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	userResp := &entities.User{
		Id:       id,
		Login:    user.Login,
		Password: password,
	}

	userData, _ := json.Marshal(userResp)
	ad.rdb.Set(ctx, fmt.Sprintf("user:%d", id), userData, 10*time.Minute)

	return userResp, nil
}

func (ad *Repository) CheckUserExists(login string) bool {
	_, err := ad.GetUserByLogin(login)
	if err != nil {
		return true
	}

	return false
}

func (ad *Repository) GetUserByLogin(login string) (*entities.User, error) {
	const op = "db.getUserByLogin"
	var id int
	var password string

	query := "SELECT id, password FROM users WHERE login = $1"
	err := ad.db.QueryRow(query, login).Scan(&id, &password)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &entities.User{
		Id:       id,
		Login:    login,
		Password: password,
	}, nil
}
