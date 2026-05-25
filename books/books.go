package books

// Book representa la entidad estructural del libro electrónico (Inmutable por diseño funcional)
type Book struct {
	ID          uint
	Title       string
	Author      string
	Genre       string
	Price       float64
	IsAvailable bool
}

// BookPredicate define el tipo de función pura de orden superior para realizar filtrados dinámicos
type BookPredicate func(Book) bool

// Filter es una función pura de orden superior (Higher-Order Function)
// Recibe una colección y un predicado funcional, retornando un nuevo slice filtrado sin efectos secundarios.
func Filter(books []Book, predicate BookPredicate) []Book {
	var filtered []Book
	for _, book := range books {
		if predicate(book) {
			filtered = append(filtered, book)
		}
	}
	return filtered
}

// --- GENERADORES DE PREDICADOS (Clisures / Funciones Currificadas) ---

// FilterByGenre retorna un predicado para filtrar libros por su género de forma exacta
func FilterByGenre(genre string) BookPredicate {
	return func(b Book) bool {
		return b.Genre == genre
	}
}

// FilterByPriceLessThan retorna un predicado para verificar si el precio es menor o igual a un límite
func FilterByPriceLessThan(maxPrice float64) BookPredicate {
	return func(b Book) bool {
		return b.Price <= maxPrice
	}
}

// FilterAvailable retorna un predicado que verifica si el libro cuenta con stock o disponibilidad activa
func FilterAvailable() BookPredicate {
	return func(b Book) bool {
		return b.IsAvailable
	}
}
