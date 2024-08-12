package entities

import "time"

type Entity struct {
}

type User struct {
	Id       int    `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginReqUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Advert struct {
	Id         int       `json:"id"`
	UserId     int       `json:"user_id"`
	Header     string    `json:"header"`
	Text       string    `json:"text"`
	ImageURL   string    `json:"image_url"`
	Address    string    `json:"address"`
	Price      float64   `json:"price"`
	Datetime   time.Time `json:"datetime"`
	ByThisUser bool      `json:"by_this_user"`
}

type AdvList struct {
	List []Advert `json:"feed"`
}

type Filter struct {
	MinPrice       float64 `json:"min_price"`
	MaxPrice       float64 `json:"max_price"`
	ByPrice        bool    `json:"by_price"`
	AscendingOrder bool    `json:"ascending_order"`
}
