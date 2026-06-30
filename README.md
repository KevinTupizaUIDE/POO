# Sistema de Gestión de Libros Electrónicos

**Autor:** Kevin Tupiza  
**Fecha:** 24/05/2026 (Actualizado)  
**Descripción:** Aprendizaje Autónomo 1 - Selección de Sistemas de Gestión Empresarial (Libros Electrónicos)

---

Este proyecto es una aplicación web y API REST escrita en **Go** diseñada para gestionar un catálogo de libros electrónicos. La aplicación expone endpoints interactivos para consultas y filtros basados en programación funcional, permite la persistencia de datos mediante el ORM **GORM** en **PostgreSQL**, e incorpora una interfaz gráfica minimalista en colores negro y azul eléctrico de primer nivel.

## Características Principales

*   **Interfaz Gráfica Premium (Web Dashboard)**:
    *   Diseño minimalista moderno con paleta en negro profundo (`#030712`) y azul eléctrico.
    *   Glassmorphism en componentes y micro-animaciones en tarjetas de libros.
    *   Demostración interactiva client-side de la encapsulación interna y mutaciones (Getter/Setter).
*   **Filtros Funcionales Puros (Higher-Order Functions)**:
    *   Uso de predicados de filtrado dinámico para género (`FilterByGenre`), precio límite (`FilterByPriceLessThan`) y disponibilidad activa (`FilterAvailable`).
*   **Persistencia y Auto-gestión con GORM**:
    *   Verificación automática de base de datos e inserción inteligente de registros semilla (Seeding).
    *   **Resiliencia (Fallback)**: Si no se detecta la base de datos PostgreSQL activa, la aplicación cambia de modo automáticamente y sirve la información desde un catálogo local en memoria para evitar colapsos.
*   **API REST y Rutas CRUD**:
    *   `/api/catalogo`: Obtención de catálogo completo y filtrado dinámico por parámetros query (`?genre=X`, `?max_price=Y`, `?available=true`).
    *   `/api/catalogo/buscar`: Búsqueda específica de libros por su llave primaria (ID) mediante la consulta `?id=X`.
*   **Testing de Calidad**:
    *   Pruebas unitarias para validar la encapsulación de variables y funciones de orden superior.

---

## Estructura del Proyecto

```text
project/
│
├── books/
│   ├── books.go         # Entidad Book, encapsulación de estado y predicados funcionales
│   └── books_test.go    # Pruebas unitarias de encapsulación y filtros funcionales
│
├── config/
│   └── config.go        # Configuración por variables de entorno (.env)
│
├── database/
│   └── db.go            # Inicialización de PostgreSQL mediante GORM y semillado
│
├── public/              # Archivos estáticos de la interfaz web
│   ├── index.html       # Estructura e interfaz gráfica del dashboard
│   ├── style.css        # Hoja de estilos con tema azul/negro y efectos visuales
│   └── app.js           # Lógica JavaScript para fetching de API y control del DOM
│
├── .env.example         # Plantilla de variables de entorno
├── .env                 # Credenciales locales de base de datos (ignorado en git)
├── .gitignore           # Archivos ignorados por Git
├── go.mod               # Módulo Go
├── go.sum               # Verificación de dependencias
├── README.md            # Documentación del proyecto (este archivo)
└── main.go              # Servidor HTTP, controlador de API y logger asíncrono
```

---

## Requisitos Previos

1.  **Go** (versión 1.21 o superior).
2.  **PostgreSQL** (versión 17 o superior) en ejecución local.

---

## Configuración y Ejecución

### 1. Variables de Entorno
Copia el archivo `.env.example` y renómbralo a `.env`:
```bash
cp .env.example .env
```
Abre el archivo `.env` e ingresa tus credenciales de PostgreSQL:
```env
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=tu_contraseña_aquí
DB_NAME=ebooks_db
DB_PORT=5432
APP_ENV=development
```

### 2. Ejecutar la Aplicación
Descarga las dependencias e inicia el servidor web:
```bash
go run main.go
```

Al iniciar, la terminal te indicará la conexión a PostgreSQL y levantará la interfaz web:
```text
==========================================================
 SERVIDOR WEB INICIADO: Abre http://localhost:8080 en tu navegador
==========================================================
```
Abre tu navegador de preferencia e ingresa a: **http://localhost:8080**

---

## Ejecución de Pruebas Unitarias
Para verificar que la encapsulación del estado interno del libro y el filtrado funcional operen correctamente, ejecuta en la terminal:
```bash
go test ./books -v
```

**Salida exitosa esperada:**
```text
=== RUN   TestEncapsulacionEstadoInterno
--- PASS: TestEncapsulacionEstadoInterno (0.00s)
=== RUN   TestFiltradoFuncional
--- PASS: TestFiltradoFuncional (0.00s)
PASS
ok  	sistema_ebooks/books	1.012s
```
