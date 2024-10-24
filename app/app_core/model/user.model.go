package models

type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`   // Include password field
	CreatedAt string `json:"created_at"` // Add timestamps for created_at
	UpdatedAt string `json:"updated_at"` // Add timestamps for updated_at
}
