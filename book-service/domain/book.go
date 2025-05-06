package domain

// Book — основная сущность книги
type Book struct {
	ID                int64
	Title             string
	Author            string
	ISBN              string
	TotalQuantity     int32
	AvailableQuantity int32
}
