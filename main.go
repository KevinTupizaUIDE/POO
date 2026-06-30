//Comentarios de línea
/*
@nombre: Kevin Tupiza
@fecha: 24/05/2026
@descripción: Aprendizaje Autónomo 1 Selección de Sistemas de Gestión empresarial (Libros Electrónicos)
*/

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sistema_ebooks/books"
	"sistema_ebooks/config"
	"sistema_ebooks/database"
	"strconv"
)

var dbConn *database.DBConnection
var catalog []books.Book
var logChan = make(chan string, 100)

func main() {
	fmt.Println("==========================================================")
	fmt.Println(" INICIALIZANDO: SISTEMA DE GESTIÓN DE LIBROS ELECTRÓNICOS")
	fmt.Println("==========================================================")

	// 1. Carga de configuraciones de entorno (Etapa 1 - Planeación)
	cfg := config.LoadConfig()
	fmt.Printf("[Config] Entorno de ejecución actual: %s\n", cfg.AppEnv)

	// 2. Inicialización de la capa de persistencia (Estructura de GORM incorporada)
	dbConn = database.InitDB(cfg)

	// 3. Dataset semilla en memoria para demostrar el funcionamiento del paradigma funcional
	catalog = []books.Book{
		{ID: 1, Title: "The Go Programming Language", Author: "Donovan & Kernighan", Genre: "Tecnología", Price: 39.99, IsAvailable: true},
		{ID: 2, Title: "Clean Architecture", Author: "Robert C. Martin", Genre: "Tecnología", Price: 45.50, IsAvailable: true},
		{ID: 3, Title: "Don Quijote de la Mancha", Author: "Miguel de Cervantes", Genre: "Clásicos", Price: 15.99, IsAvailable: false},
		{ID: 4, Title: "Cien Años de Soledad", Author: "Gabriel García Márquez", Genre: "Ficción", Price: 22.00, IsAvailable: true},
		{ID: 5, Title: "Go in Action", Author: "William Kennedy", Genre: "Tecnología", Price: 29.99, IsAvailable: false},
	}

	// 3.5. Configuración de Goroutines para Lector de Logs
	go func() {
		for logMsg := range logChan {
			fmt.Printf("\n[HTTP Log] %s\n", logMsg)
		}
	}()

	// 4. Registro de Endpoints del Servidor HTTP
	http.HandleFunc("/api/catalogo", obtenerCatalogoHandler)          // Consulta y filtrado funcional
	http.HandleFunc("/api/catalogo/buscar", obtenerLibroPorIDHandler) // Nuevo requerimiento del profesor (ID)

	// Servir archivos estáticos del frontend de la carpeta "./public" en "/"
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	// 5. Iniciar el Servidor HTTP (Bloqueante, reemplaza al menú de consola)
	fmt.Println("\n==========================================================")
	fmt.Println(" SERVIDOR WEB INICIADO: Abre http://localhost:8080 en tu navegador")
	fmt.Println("==========================================================")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("[HTTP Error] Servidor falló al arrancar: %v\n", err)
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

// Controlador para obtener todo el catálogo con filtrado funcional (Requisito CRUD/Funcional)
func obtenerCatalogoHandler(w http.ResponseWriter, r *http.Request) {
	var libros []books.Book
	var dbError error

	if dbConn != nil && dbConn.Instance != nil {
		// Obtener catálogo de la base de datos
		resultado := dbConn.Instance.Find(&libros)
		dbError = resultado.Error
	} else {
		// Fallback resiliente en memoria
		libros = catalog
	}

	if dbError != nil {
		logChan <- fmt.Sprintf("Error 500: Error al cargar el catálogo de la BD: %v", dbError)
		http.Error(w, `{"error": "No se pudo cargar el catálogo"}`, http.StatusInternalServerError)
		return
	}

	// Filtros funcionales (Paradigma Funcional - Cátedra)
	
	// Filtro por género exacto
	if genre := r.URL.Query().Get("genre"); genre != "" {
		libros = books.Filter(libros, books.FilterByGenre(genre))
	}
	
	// Filtro por precio menor o igual
	if maxPriceStr := r.URL.Query().Get("max_price"); maxPriceStr != "" {
		if maxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			libros = books.Filter(libros, books.FilterByPriceLessThan(maxPrice))
		}
	}
	
	// Filtro por disponibilidad
	if availableStr := r.URL.Query().Get("available"); availableStr != "" {
		if available, err := strconv.ParseBool(availableStr); err == nil && available {
			libros = books.Filter(libros, books.FilterAvailable())
		}
	}

	logChan <- fmt.Sprintf("Búsqueda exitosa: %d libros retornados", len(libros))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(libros)
}

// Controlador para obtener un libro específico por su ID (Requisito CRUD)
func obtenerLibroPorIDHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id") // Se espera una petición como: /api/catalogo/buscar?id=1
	
	var libro books.Book
	var dbError error

	if dbConn != nil && dbConn.Instance != nil {
		// Uso de GORM para buscar por Primary Key (Equivalente al SELECT WHERE ID = ?)
		resultado := dbConn.Instance.First(&libro, id)
		dbError = resultado.Error
	} else {
		// Fallback in-memory para resiliencia (modo simulado)
		idInt, convErr := strconv.Atoi(id)
		if convErr != nil {
			dbError = convErr
		} else {
			encontrado := false
			for _, b := range catalog {
				if int(b.ID) == idInt {
					libro = b
					encontrado = true
					break
				}
			}
			if !encontrado {
				dbError = fmt.Errorf("libro no encontrado")
			}
		}
	}
	
	if dbError != nil {
		logChan <- fmt.Sprintf("Error 404: Búsqueda fallida para el ID %s", id)
		http.Error(w, `{"error": "Libro no encontrado"}`, http.StatusNotFound)
		return
	}

	logChan <- fmt.Sprintf("Búsqueda exitosa: Libro ID %s retornado", id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(libro)
}
