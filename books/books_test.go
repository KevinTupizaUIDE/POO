package books

import (
	"testing"
)

// Prueba Unitaria: Verifica que la encapsulación funcione correctamente
func TestEncapsulacionEstadoInterno(t *testing.T) {
	libro := Book{Title: "Go Programming"}
	
	// Prueba del valor por defecto
	if libro.GetEstadoInterno() != "Sin auditar" {
		t.Errorf("Se esperaba 'Sin auditar', se obtuvo '%s'", libro.GetEstadoInterno())
	}

	// Prueba del Setter
	libro.SetEstadoInterno("Revisado")
	if libro.GetEstadoInterno() != "Revisado" {
		t.Errorf("Se esperaba 'Revisado', se obtuvo '%s'", libro.GetEstadoInterno())
	}
}

// Prueba Funcional: Verifica la función de orden superior (Filter)
func TestFiltradoFuncional(t *testing.T) {
	catalogo := []Book{
		{Title: "Libro 1", Genre: "Ficción"},
		{Title: "Libro 2", Genre: "Tecnología"},
	}

	resultado := Filter(catalogo, FilterByGenre("Tecnología"))

	if len(resultado) != 1 {
		t.Fatalf("Se esperaba 1 resultado, se obtuvieron %d", len(resultado))
	}
	if resultado[0].Title != "Libro 2" {
		t.Errorf("El filtrado no retornó el libro correcto")
	}
}
