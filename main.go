//Comentarios de línea
/*
@nombre: Kevin Tupiza
@fecha: 24/05/2026
@descripción: Aprendizaje Autónomo 1 Selección de Sistemas de Gestión empresarial (Libros Electrónicos)
*/

package main

import (
	"fmt"
	"os"
	"sistema_ebooks/books"
	"sistema_ebooks/config"
	"sistema_ebooks/database"
	"strconv"
)

func main() {
	fmt.Println("==========================================================")
	fmt.Println(" INICIALIZANDO: SISTEMA DE GESTIÓN DE LIBROS ELECTRÓNICOS")
	fmt.Println("==========================================================")

	// 1. Carga de configuraciones de entorno (Etapa 1 - Planeación)
	cfg := config.LoadConfig()
	fmt.Printf("[Config] Entorno de ejecución actual: %s\n", cfg.AppEnv)

	// 2. Inicialización de la capa de persistencia (Estructura de GORM incorporada)
	_ = database.InitDB(cfg)

	// 3. Dataset semilla en memoria para demostrar el funcionamiento del paradigma funcional
	catalog := []books.Book{
		{ID: 1, Title: "The Go Programming Language", Author: "Donovan & Kernighan", Genre: "Tecnología", Price: 39.99, IsAvailable: true},
		{ID: 2, Title: "Clean Architecture", Author: "Robert C. Martin", Genre: "Tecnología", Price: 45.50, IsAvailable: true},
		{ID: 3, Title: "Don Quijote de la Mancha", Author: "Miguel de Cervantes", Genre: "Clásicos", Price: 15.99, IsAvailable: false},
		{ID: 4, Title: "Cien Años de Soledad", Author: "Gabriel García Márquez", Genre: "Ficción", Price: 22.00, IsAvailable: true},
		{ID: 5, Title: "Go in Action", Author: "William Kennedy", Genre: "Tecnología", Price: 29.99, IsAvailable: false},
	}

	// 4. Estructura de Control de Flujo Iterativo (Menú solicitado por la cátedra)
	for {
		fmt.Println("\n--- MENÚ OPERACIONAL DEL SISTEMA (ETAPA 1) ---")
		fmt.Println("1. Listar Catálogo Completo")
		fmt.Println("2. Filtrar Libros de 'Tecnología' (Programación Funcional)")
		fmt.Println("3. Filtrar Libros por Precio Económico (<= $30.00)")
		fmt.Println("4. Mostrar Solo Libros Disponibles para Descarga")
		fmt.Println("5. Salir de la Aplicación")
		fmt.Print("Seleccione una opción: ")

		var optionStr string
		fmt.Scanln(&optionStr)

		option, err := strconv.Atoi(optionStr)
		if err != nil {
			fmt.Println("[Error] Por favor, ingrese un número de opción válido.")
			continue
		}

		switch option {
		case 1:
			fmt.Println("\n--- CATÁLOGO COMPLETO ---")
			displayBooks(catalog)
		case 2:
			fmt.Println("\n--- LIBROS DE TECNOLOGÍA RETORNADOS POR FILTRO FUNCIONAL ---")
			// Uso de función de orden superior pasando el predicado de género
			techBooks := books.Filter(catalog, books.FilterByGenre("Tecnología"))
			displayBooks(techBooks)
		case 3:
			fmt.Println("\n--- LIBROS ECONÓMICOS (<= $30.00) RETORNADOS POR FILTRO FUNCIONAL ---")
			// Uso de función de orden superior pasando el predicado de precio umbral
			cheapBooks := books.Filter(catalog, books.FilterByPriceLessThan(30.00))
			displayBooks(cheapBooks)
		case 4:
			fmt.Println("\n--- LIBROS DISPONIBLES INMEDIATAMENTE ---")
			// Uso de función de orden superior pasando el predicado de disponibilidad activa
			availableBooks := books.Filter(catalog, books.FilterAvailable())
			displayBooks(availableBooks)
		case 5:
			fmt.Println("\n[Sistema] Finalizando ejecución del software. Planificación completada con éxito.")
			os.Exit(0)
		default:
			fmt.Println("[Error] Opción fuera de rango. Intente nuevamente.")
		}
	}
}

// Función auxiliar para imprimir las colecciones resultantes en consola de forma limpia
func displayBooks(bList []books.Book) {
	if len(bList) == 0 {
		fmt.Println("No se encontraron registros que coincidan con el criterio seleccionado.")
		return
	}
	fmt.Printf("%-4s | %-30s | %-22s | %-12s | %-8s | %-10s\n", "ID", "TÍTULO", "AUTOR", "GÉNERO", "PRECIO", "DISPONIBLE")
	fmt.Println("------------------------------------------------------------------------------------------------------")
	for _, b := range bList {
		disp := "No"
		if b.IsAvailable {
			disp = "Sí"
		}
		fmt.Printf("%-4d | %-30s | %-22s | %-12s | $%-7.2f | %-10s\n", b.ID, b.Title, b.Author, b.Genre, b.Price, disp)
	}
}
