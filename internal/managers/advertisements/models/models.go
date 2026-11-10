package models

type Advertisement struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Price       int    `json:"price"`
}

type AdvertisementsList struct {
	Items []Advertisement `json:"items"`
}
