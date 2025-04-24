package mapping

import "github.com/Prrost/assignment1proto/proto/inventory"

type Product struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Amount    int64  `json:"amount"`
	Available bool   `json:"available"`
}

func ToProduct(product *inventory.Product) *Product {
	return &Product{
		ID:        int(product.Id),
		Name:      product.Name,
		Amount:    product.Amount,
		Available: product.Available,
	}
}

func ToProducts(products []*inventory.Product) []*Product {
	result := make([]*Product, 0, len(products))
	for _, product := range products {
		result = append(result, ToProduct(product))
	}
	return result
}
