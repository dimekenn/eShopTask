package models

type Product struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Count       int    `json:"count"`
}

type Cart struct {
	UserId   int        `json:"user_id"`
	Id       int        `json:"id"`
	Products []*CartProduct `json:"product_ids"`
}

type CartProduct struct {
	Id        int `json:"id"`
	ProductId int `json:"product_id"`
	Count     int `json:"count"`
}

type DBCart struct {
	UserId   int        `json:"user_id"`
	Id       int        `json:"id"`
	Products []string `json:"product_ids"`
}

type GetCartRes struct {
	Id int `json:"id"`
	ProductId int `json:"product_id"`
	CartId int `json:"cart_id"`
	Count int `json:"count"`
}

type Msg struct {
	Msg string `json:"msg"`
}