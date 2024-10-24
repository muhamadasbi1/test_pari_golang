package models

type Item struct {
	ID           int     `json:"id"`
	CategoryID   int     `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	Price        float64 `json:"price"`
	CreatedAt    string  `json:"created_at"`
}
