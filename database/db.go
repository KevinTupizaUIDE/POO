package database

import (
	"fmt"
	"log"
	"sistema_ebooks/books"
	"sistema_ebooks/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DBConnection encapsula la instancia de la base de datos
type DBConnection struct {
	Instance *gorm.DB
}

// InitDB inicializa la conexión con GORM utilizando la configuración cargada
func InitDB(cfg *config.Config) *DBConnection {
	// 1. Conexión temporal a la base de datos por defecto 'postgres' para verificar/crear la BD destino
	postgresDsn := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=%s sslmode=disable TimeZone=UTC",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBPort)

	dbDefault, err := gorm.Open(postgres.Open(postgresDsn), &gorm.Config{})
	if err != nil {
		log.Printf("Aviso: No se pudo conectar a la base de datos 'postgres' de control: %v. Intentando conexión directa...\n", err)
	} else {
		var exists bool
		row := dbDefault.Raw("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = ?)", cfg.DBName).Row()
		if err := row.Scan(&exists); err == nil && !exists {
			log.Printf("Base de datos '%s' no encontrada. Creándola...\n", cfg.DBName)
			if err := dbDefault.Exec(fmt.Sprintf("CREATE DATABASE %s", cfg.DBName)).Error; err != nil {
				log.Printf("Error al crear la base de datos '%s': %v\n", cfg.DBName, err)
			} else {
				log.Printf("Base de datos '%s' creada con éxito.\n", cfg.DBName)
			}
		}
		sqlDB, err := dbDefault.DB()
		if err == nil {
			sqlDB.Close()
		}
	}

	// 2. Construcción de la cadena de conexión final (DSN)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Error al conectar a la base de datos final: %v\n", err)
		log.Println("Arquitectura de GORM configurada en modo simulado por fallo de conexión.")
		return &DBConnection{Instance: nil}
	}

	log.Println("¡Conexión a la base de datos PostgreSQL establecida con éxito mediante GORM!")

	// 3. Ejecutar Migraciones de Esquema Automáticas (AutoMigrate)
	if err := db.AutoMigrate(&books.Book{}); err != nil {
		log.Printf("Error durante AutoMigrate: %v\n", err)
	} else {
		log.Println("Migración de esquema de la tabla 'books' completada.")
	}

	// 4. Alimentar el catálogo inicial si la tabla se encuentra vacía (Seed)
	var count int64
	db.Model(&books.Book{}).Count(&count)
	if count == 0 {
		log.Println("La tabla 'books' está vacía. Insertando datos semilla...")
		seedBooks := []books.Book{
			{Title: "The Go Programming Language", Author: "Donovan & Kernighan", Genre: "Tecnología", Price: 39.99, IsAvailable: true},
			{Title: "Clean Architecture", Author: "Robert C. Martin", Genre: "Tecnología", Price: 45.50, IsAvailable: true},
			{Title: "Don Quijote de la Mancha", Author: "Miguel de Cervantes", Genre: "Clásicos", Price: 15.99, IsAvailable: false},
			{Title: "Cien Años de Soledad", Author: "Gabriel García Márquez", Genre: "Ficción", Price: 22.00, IsAvailable: true},
			{Title: "Go in Action", Author: "William Kennedy", Genre: "Tecnología", Price: 29.99, IsAvailable: false},
		}
		if err := db.Create(&seedBooks).Error; err != nil {
			log.Printf("Error al insertar los registros semilla: %v\n", err)
		} else {
			log.Println("Registros semilla insertados con éxito.")
		}
	}

	return &DBConnection{Instance: db}
}
