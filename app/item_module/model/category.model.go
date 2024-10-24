package models

type Category struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"` // Add timestamps for created_at
	UpdatedAt string `json:"updated_at"` // Add timestamps for updated_at
}
