package main

type Delivery struct {
	OrderId        int            `json:"order_id"`
	TableId        int            `json:"table_id"`
	Items          []int          `json:"items"`
	Priority       int            `json:"priority"`
	MaxWait        int            `json:"max_wait"`
	PickUpTime     int64          `json:"pick_up_time"`
	CookingTime    int            `json:"cooking_time"`
	CookingDetails []FoodDelivery `json:"cooking_details"`
}


type FoodDelivery struct {
	FoodId int `json:"food_id"`
	CookId int `json:"cook_id"`
}
