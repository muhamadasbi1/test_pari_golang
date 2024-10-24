package dto

type ItemCreateRequest struct {
	CategoryID  int     `json:"category_id" validate:"required"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description,omitempty"`
	Price       float64 `json:"price" validate:"required,numeric"`
}

type ItemUpdateRequest struct {
	ID          int     `json:"id" validate:"required"`
	CategoryID  int     `json:"category_id" validate:"required"`
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description,omitempty"`
	Price       float64 `json:"price" validate:"required,numeric"`
}
