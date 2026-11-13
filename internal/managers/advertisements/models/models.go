package models

type Photo struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	SortOrder int    `json:"sort_order"`
	CreatedAt string `json:"created_at,omitempty"`
}

type Advertisement struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	UserID      string  `json:"user_id,omitempty"`
	Photos      []Photo `json:"photos,omitempty"`
	CreatedAt   string  `json:"created_at,omitempty"`
	UpdatedAt   string  `json:"updated_at,omitempty"`
}

type AdvertisementsList struct {
	Items []Advertisement `json:"items"`
}

type AdvertisementID struct {
	ID string `json:"id"`
}
