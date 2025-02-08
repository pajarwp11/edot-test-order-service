package cart

type Cart struct {
	UserId   int     `json:"user_id"`
	Products Product `json:"products"`
}

type Product struct {
	ProductId   int
	ProductName string
	Quantity    int
	Price       int
	TotalPrice  int
}

type UpsertRequest struct {
	UserId      int    `json:"user_id" validate:"required"`
	ProductId   int    `json:"product_id" validate:"required"`
	ProductName string `json:"product_name" validate:"required"`
	Quantity    int    `json:"quantity" validate:"required"`
	Price       int    `json:"price" validate:"required"`
	TotalPrice  int    `json:"total_price" validate:"required"`
}
