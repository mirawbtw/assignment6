package domain

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float32 `json:"price"` // Измените float64 на float32
	Stock       int32   `json:"stock"` // Измените int на int32
}
